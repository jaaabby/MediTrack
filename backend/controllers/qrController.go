package controllers

import (
	"fmt"
	"meditrack/models"
	"meditrack/services"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"meditrack/pkg/response"

	"github.com/gin-gonic/gin"
)

type QRController struct {
	qrService            services.QRService
	medicalSupplyService services.MedicalSupplyService
}

func NewQRController(qrService services.QRService) *QRController {
	return &QRController{qrService: qrService}
}

// SetMedicalSupplyService establece el servicio de suministros médicos
func (c *QRController) SetMedicalSupplyService(medicalSupplyService services.MedicalSupplyService) {
	c.medicalSupplyService = medicalSupplyService
}

// determineQRType determina el tipo de código QR
func determineQRType(qrCode string) string {
	upperQR := strings.ToUpper(qrCode)
	if strings.HasPrefix(upperQR, "SUPPLY_") {
		return "SUPPLY"
	} else if strings.HasPrefix(upperQR, "BATCH_") {
		return "BATCH"
	}
	return "UNKNOWN"
}

// countAvailableSupplies cuenta los insumos disponibles
func countAvailableSupplies(supplies []map[string]interface{}) int {
	count := 0
	for _, supply := range supplies {
		if available, ok := supply["is_available"].(bool); ok && available {
			count++
		}
	}
	return count
}

// countConsumedSupplies cuenta los insumos consumidos
func countConsumedSupplies(supplies []map[string]interface{}) int {
	count := 0
	for _, supply := range supplies {
		if consumed, ok := supply["is_consumed"].(bool); ok && consumed {
			count++
		}
	}
	return count
}

// =============================================
// ENDPOINT PRINCIPAL CON TRAZABILIDAD
// =============================================

// ScanQR escanea un código QR y registra automáticamente el evento
func (c *QRController) ScanQR(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Código QR requerido",
		})
		return
	}

	// Decodificar URL
	decodedQR, err := url.QueryUnescape(qrCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Código QR inválido",
		})
		return
	}

	// Crear contexto de escaneo desde la request
	scanContext := c.buildScanContext(ctx)

	// Usar el nuevo método con logging automático
	qrInfo, err := c.qrService.ScanQRWithAutoLogging(decodedQR, scanContext)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Código QR no encontrado: " + err.Error(),
		})
		return
	}

	// PRIORIZAR INSUMOS INDIVIDUALES
	qrType := determineQRType(decodedQR)

	// Construir respuesta mejorada
	resultMap := gin.H{
		"success":              true,
		"type":                 qrInfo.Type,
		"id":                   qrInfo.ID,
		"qr_code":              qrInfo.QRCode,
		"batch_info":           qrInfo.BatchInfo,
		"supply_code":          qrInfo.SupplyCode,
		"history":              qrInfo.History,
		"is_consumed":          false,
		"can_consume":          true,
		"batch_status":         "active",
		"scan_timestamp":       time.Now(),
		"qr_type":              qrType,
		"is_individual_supply": qrType == "SUPPLY",
		"traceability":         qrInfo.Traceability,
		"scan_events":          qrInfo.ScanEvents,
		"scan_statistics":      qrInfo.ScanStatistics,
	}

	// Estructura especial para supply_info con batch anidado
	if qrInfo.SupplyInfo != nil {
		supplyInfoMap := gin.H{
			"ID":           qrInfo.SupplyInfo.ID,
			"Code":         qrInfo.SupplyInfo.Code,
			"BatchID":      qrInfo.SupplyInfo.BatchID,
			"QRCode":       qrInfo.SupplyInfo.QRCode,
			"IsConsumed":   qrInfo.SupplyInfo.IsConsumed,
			"LastMovement": qrInfo.SupplyInfo.LastMovement,
			"DaysToExpire": qrInfo.SupplyInfo.DaysToExpire,
		}

		// Agregar información del batch dentro de supply_info
		if qrInfo.SupplyInfo.BatchInfo != nil {
			supplyInfoMap["batch"] = qrInfo.SupplyInfo.BatchInfo
		}

		// Agregar información del código de insumo
		if qrInfo.SupplyInfo.SupplyCode != nil {
			supplyInfoMap["SupplyCode"] = qrInfo.SupplyInfo.SupplyCode
		}

		resultMap["supply_info"] = supplyInfoMap
		resultMap["is_consumed"] = qrInfo.SupplyInfo.IsConsumed
		resultMap["can_consume"] = !qrInfo.SupplyInfo.IsConsumed
	}

	if qrType == "SUPPLY" {
		// Información contextual para insumos individuales
		resultMap["scan_priority"] = "high"
		resultMap["recommended_actions"] = []string{"consume", "view_history", "check_batch"}
		resultMap["ui_focus"] = "individual_supply"
	} else if qrType == "BATCH" {
		// Información de lote con recomendación de usar insumos individuales
		resultMap["scan_priority"] = "medium"
		resultMap["recommended_actions"] = []string{"view_batch", "scan_individual_supplies"}
		resultMap["ui_focus"] = "batch_overview"
		resultMap["user_guidance"] = "Para trazabilidad completa, escanee códigos QR de insumos individuales del lote"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resultMap,
	})
}

// buildScanContext construye el contexto de escaneo desde la request HTTP
func (c *QRController) buildScanContext(ctx *gin.Context) *services.ScanContext {
	scanContext := &services.ScanContext{
		ScanSource: models.ScanSourceWeb, // default
	}

	// Obtener información del usuario desde headers o query params
	if userRUT := ctx.GetHeader("X-User-RUT"); userRUT != "" {
		scanContext.UserRUT = &userRUT
	}
	if userRUT := ctx.Query("user_rut"); userRUT != "" {
		scanContext.UserRUT = &userRUT
	}

	if userName := ctx.GetHeader("X-User-Name"); userName != "" {
		scanContext.UserName = &userName
	}
	if userName := ctx.Query("user_name"); userName != "" {
		scanContext.UserName = &userName
	}

	// Obtener información de ubicación
	if pavilionIDStr := ctx.Query("pavilion_id"); pavilionIDStr != "" {
		if pavilionID, err := strconv.Atoi(pavilionIDStr); err == nil {
			scanContext.PavilionID = &pavilionID
		}
	}

	if medicalCenterIDStr := ctx.Query("medical_center_id"); medicalCenterIDStr != "" {
		if medicalCenterID, err := strconv.Atoi(medicalCenterIDStr); err == nil {
			scanContext.MedicalCenterID = &medicalCenterID
		}
	}

	// Obtener propósito del escaneo
	if purpose := ctx.Query("scan_purpose"); purpose != "" {
		scanContext.ScanPurpose = &purpose
	}

	// Obtener información técnica
	if userAgent := ctx.GetHeader("User-Agent"); userAgent != "" {
		scanContext.UserAgent = &userAgent
	}

	// Obtener IP address
	if clientIP := ctx.ClientIP(); clientIP != "" {
		if ip := net.ParseIP(clientIP); ip != nil {
			scanContext.IPAddress = &ip
		}
	}

	// Obtener información de sesión
	if sessionID := ctx.GetHeader("X-Session-ID"); sessionID != "" {
		scanContext.SessionID = &sessionID
	}
	if sessionID := ctx.Query("session_id"); sessionID != "" {
		scanContext.SessionID = &sessionID
	}

	// Determinar fuente de escaneo
	if source := ctx.Query("scan_source"); source != "" {
		scanContext.ScanSource = source
	} else if ctx.GetHeader("X-Mobile-App") != "" {
		scanContext.ScanSource = models.ScanSourceMobile
	}

	// Notas adicionales
	if notes := ctx.Query("notes"); notes != "" {
		scanContext.Notes = &notes
	}

	return scanContext
}

// =============================================
// ENDPOINTS DE TRAZABILIDAD AVANZADA
// =============================================

// GetCompleteTraceability obtiene la trazabilidad completa de un QR
func (c *QRController) GetCompleteTraceability(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Código QR requerido",
		})
		return
	}

	traceability, err := c.qrService.GetCompleteTraceability(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    traceability,
	})
}

// GetScanHistory obtiene el historial de escaneos de un QR
func (c *QRController) GetScanHistory(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Código QR requerido",
		})
		return
	}

	limitStr := ctx.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	history, err := c.qrService.GetScanEventHistory(qrCode, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"qr_code": qrCode,
			"history": history,
			"total":   len(history),
		},
	})
}

// GetScanStatistics obtiene estadísticas de escaneo de un QR
func (c *QRController) GetScanStatistics(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Código QR requerido",
		})
		return
	}

	stats, err := c.qrService.GetQRScanStatistics(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// RegisterManualScanEvent permite registrar manualmente un evento de escaneo
func (c *QRController) RegisterManualScanEvent(ctx *gin.Context) {
	var request struct {
		QRCode          string `json:"qr_code" binding:"required"`
		UserRUT         string `json:"user_rut,omitempty"`
		UserName        string `json:"user_name,omitempty"`
		PavilionID      *int   `json:"pavilion_id,omitempty"`
		MedicalCenterID *int   `json:"medical_center_id,omitempty"`
		ScanPurpose     string `json:"scan_purpose,omitempty"`
		Notes           string `json:"notes,omitempty"`
		ScanSource      string `json:"scan_source,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Datos inválidos: " + err.Error(),
		})
		return
	}

	// Construir contexto
	scanContext := &services.ScanContext{
		ScanSource:      models.ScanSourceAPI,
		PavilionID:      request.PavilionID,
		MedicalCenterID: request.MedicalCenterID,
	}

	if request.UserRUT != "" {
		scanContext.UserRUT = &request.UserRUT
	}
	if request.UserName != "" {
		scanContext.UserName = &request.UserName
	}
	if request.ScanPurpose != "" {
		scanContext.ScanPurpose = &request.ScanPurpose
	}
	if request.Notes != "" {
		scanContext.Notes = &request.Notes
	}
	if request.ScanSource != "" {
		scanContext.ScanSource = request.ScanSource
	}

	// Información técnica de la request
	if userAgent := ctx.GetHeader("User-Agent"); userAgent != "" {
		scanContext.UserAgent = &userAgent
	}
	if clientIP := ctx.ClientIP(); clientIP != "" {
		if ip := net.ParseIP(clientIP); ip != nil {
			scanContext.IPAddress = &ip
		}
	}

	// Registrar el escaneo
	qrInfo, err := c.qrService.ScanQRWithAutoLogging(request.QRCode, scanContext)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Evento de escaneo registrado correctamente",
		"data":    qrInfo,
	})
}

// =============================================
// ENDPOINTS ORIGINALES MANTENIDOS
// =============================================

func (c *QRController) ConsumeIndividualSupply(ctx *gin.Context) {
	var request struct {
		QRCode          string `json:"qr_code" binding:"required"`
		UserRUT         string `json:"user_rut" binding:"required"`
		DestinationType string `json:"destination_type" binding:"required"`
		DestinationID   int    `json:"destination_id" binding:"required"`
		Notes           string `json:"notes,omitempty"`
		ConsumedAt      string `json:"consumed_at,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Datos inválidos: " + err.Error(),
		})
		return
	}

	// Validar que sea un código QR de insumo individual
	if !strings.HasPrefix(strings.ToUpper(request.QRCode), "SUPPLY_") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Solo se pueden consumir insumos individuales. Use códigos QR que comiencen con SUPPLY_",
			"qr_type": "invalid",
		})
		return
	}

	// Crear contexto de escaneo para el consumo
	scanContext := &services.ScanContext{
		UserRUT:     &request.UserRUT,
		ScanSource:  models.ScanSourceWeb,
		ScanPurpose: stringPtr(models.ScanPurposeConsume),
	}

	if request.Notes != "" {
		scanContext.Notes = &request.Notes
	}

	// Determinar ubicación según el tipo de destino
	if request.DestinationType == "pavilion" {
		scanContext.PavilionID = &request.DestinationID
	}

	// Información técnica de la request
	if userAgent := ctx.GetHeader("User-Agent"); userAgent != "" {
		scanContext.UserAgent = &userAgent
	}
	if clientIP := ctx.ClientIP(); clientIP != "" {
		if ip := net.ParseIP(clientIP); ip != nil {
			scanContext.IPAddress = &ip
		}
	}

	// Registrar el escaneo primero
	qrInfo, err := c.qrService.ScanQRWithAutoLogging(request.QRCode, scanContext)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Verificar que el insumo no esté ya consumido
	if qrInfo.SupplyInfo == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No se encontró información del insumo",
		})
		return
	}

	if qrInfo.SupplyInfo.IsConsumed {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error":      "Este insumo individual ya fue consumido anteriormente",
			"error_type": "already_consumed",
			"suggestions": []string{
				"Escanee otro insumo del mismo lote",
				"Verifique el historial del insumo",
			},
		})
		return
	}

	// Procesar consumo
	consumeReq := services.QRConsumptionRequest{
		QRCode:          request.QRCode,
		UserRUT:         request.UserRUT,
		DestinationType: request.DestinationType,
		DestinationID:   request.DestinationID,
		Notes:           request.Notes,
	}
	result, err := c.qrService.ConsumeSupplyByQR(consumeReq)

	if err != nil {
		// Manejo específico de errores para insumos individuales
		if strings.Contains(err.Error(), "ya ha sido consumido") {
			ctx.JSON(http.StatusConflict, gin.H{
				"success":    false,
				"error":      "Este insumo individual ya fue consumido anteriormente",
				"error_type": "already_consumed",
				"suggestions": []string{
					"Escanee otro insumo del mismo lote",
					"Verifique el historial del insumo",
				},
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Insumo individual consumido exitosamente",
		"data":    result,
		"next_actions": []string{
			"scan_next_supply",
			"view_batch_status",
			"check_inventory",
		},
	})
}

func (c *QRController) GetIndividualSuppliesByBatch(ctx *gin.Context) {
	batchID := ctx.Param("batch_id")
	if batchID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID de lote requerido",
		})
		return
	}

	batchIDInt, err := strconv.Atoi(batchID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID de lote inválido",
		})
		return
	}

	supplies, err := c.medicalSupplyService.GetAvailableSuppliesByBatch(batchIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"batch_id":            batchIDInt,
			"individual_supplies": supplies,
			"total_count":         len(supplies),
		},
	})
}

// ValidateQR valida si un código QR es válido
func (c *QRController) ValidateQR(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	isValid, qrType, err := c.qrService.ValidateQRCode(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Código QR válido",
		Data: map[string]interface{}{
			"valid": isValid,
			"type":  qrType,
		},
	})
}

// GetSupplyHistory obtiene el historial de un insumo por su código QR
func (c *QRController) GetSupplyHistory(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	history, err := c.qrService.GetSupplyHistory(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Historial obtenido exitosamente",
		Data:    history,
	})
}

// GenerateBatchQR genera un código QR para un lote con imagen
func (c *QRController) GenerateBatchQR(ctx *gin.Context) {
	result, err := c.qrService.GenerateBatchQR()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error generando código QR de lote: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Código QR de lote generado exitosamente",
		Data:    result,
	})
}

// GenerateSupplyQR genera un código QR para un insumo médico con imagen
func (c *QRController) GenerateSupplyQR(ctx *gin.Context) {
	result, err := c.qrService.GenerateSupplyQR()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error generando código QR de insumo: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Código QR de insumo generado exitosamente",
		Data:    result,
	})
}

// ConsumeSupply consume un insumo por su código QR y actualiza automáticamente las cantidades del lote
func (c *QRController) ConsumeSupply(ctx *gin.Context) {
	var request services.QRConsumptionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos de consumo inválidos: " + err.Error(),
		})
		return
	}

	// Validar destination_type
	if request.DestinationType != "pavilion" && request.DestinationType != "store" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "destination_type debe ser 'pavilion' o 'store'",
		})
		return
	}

	result, err := c.qrService.ConsumeSupplyByQR(request)
	if err != nil {
		// Determinar el código de estado basado en el tipo de error
		statusCode := http.StatusInternalServerError
		if err.Error() == "insumo no encontrado" || err.Error() == "código QR no válido" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "el insumo ya ha sido consumido" || err.Error() == "no hay stock disponible" {
			statusCode = http.StatusConflict
		}

		ctx.JSON(statusCode, response.Response{
			Success: false,
			Error:   err.Error(),
			Data:    result,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: result.Message,
		Data:    result,
	})
}

// GetQRImage sirve la imagen QR como PNG
func (c *QRController) GetQRImage(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de imagen QR no implementada en QRService",
	})
	return
}

// DownloadQRImage permite descargar la imagen QR con un nombre específico
func (c *QRController) DownloadQRImage(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	// Obtener resolución (normal|high)
	resolution := ctx.DefaultQuery("resolution", "normal")
	// Generar imagen QR usando el servicio
	var size int
	if resolution == "high" {
		size = 512
	} else {
		size = 256
	}

	png, err := services.GenerateQRCodeImage(qrCode, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error generando imagen QR: " + err.Error(),
		})
		return
	}

	// Establecer headers para descarga
	ctx.Header("Content-Type", "image/png")
	ctx.Header("Content-Disposition", "attachment; filename=\"qr_"+qrCode+".png\"")
	ctx.Writer.Write(png)
	ctx.Status(http.StatusOK)
}

// GetSupplyDetails obtiene información detallada de un insumo por su código QR
func (c *QRController) GetSupplyDetails(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	details, err := c.medicalSupplyService.GetSupplyWithBatchInfo(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Detalles del insumo obtenidos exitosamente",
		Data:    details,
	})
}

// SyncBatchAmounts sincroniza las cantidades de los lotes con los productos individuales
func (c *QRController) SyncBatchAmounts(ctx *gin.Context) {
	err := c.medicalSupplyService.SyncBatchAmounts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error sincronizando cantidades de lotes: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Cantidades de lotes sincronizadas exitosamente",
	})
}

// GetQRStats obtiene estadísticas generales de uso de QR codes
func (c *QRController) GetQRStats(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de estadísticas QR no implementada en QRService",
	})
	return
}

// VerifySupplyAvailability verifica si un insumo está disponible para consumo
func (c *QRController) VerifySupplyAvailability(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	info, err := c.qrService.ScanQRCode(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Solo verificar disponibilidad para insumos individuales
	if info["type"] != "medical_supply" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "La verificación de disponibilidad solo aplica para insumos individuales",
		})
		return
	}

	resp := map[string]interface{}{
		"qr_code":      qrCode,
		"is_available": info["can_consume"],
		"is_consumed":  info["is_consumed"],
		"supply_info":  info["supply_info"],
		"batch_status": info["batch_status"],
	}

	if available, ok := info["can_consume"].(bool); ok && !available {
		resp["reason"] = "Insumo ya consumido o lote sin stock"
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Verificación de disponibilidad completada",
		Data:    resp,
	})
}

// BulkConsumeSupplies permite consumir múltiples insumos en una sola operación
func (c *QRController) BulkConsumeSupplies(ctx *gin.Context) {
	var request struct {
		QRCodes         []string `json:"qr_codes" binding:"required"`
		UserRUT         string   `json:"user_rut" binding:"required"`
		DestinationType string   `json:"destination_type" binding:"required"`
		DestinationID   int      `json:"destination_id" binding:"required"`
		Notes           string   `json:"notes,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	// Validar destination_type
	if request.DestinationType != "pavilion" && request.DestinationType != "store" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "destination_type debe ser 'pavilion' o 'store'",
		})
		return
	}

	var results []map[string]interface{}
	var errors []map[string]interface{}
	successCount := 0

	for _, qrCode := range request.QRCodes {
		consumeRequest := services.QRConsumptionRequest{
			QRCode:          qrCode,
			UserRUT:         request.UserRUT,
			DestinationType: request.DestinationType,
			DestinationID:   request.DestinationID,
			Notes:           request.Notes,
		}

		result, err := c.qrService.ConsumeSupplyByQR(consumeRequest)
		if err != nil {
			errors = append(errors, map[string]interface{}{
				"qr_code": qrCode,
				"error":   err.Error(),
			})
		} else {
			results = append(results, map[string]interface{}{
				"qr_code": qrCode,
				"result":  result,
			})
			successCount++
		}
	}

	resp := map[string]interface{}{
		"total_requested": len(request.QRCodes),
		"successful":      successCount,
		"failed":          len(errors),
		"results":         results,
		"errors":          errors,
	}

	statusCode := http.StatusOK
	if len(errors) > 0 && successCount == 0 {
		statusCode = http.StatusBadRequest
	} else if len(errors) > 0 {
		statusCode = http.StatusPartialContent
	}

	ctx.JSON(statusCode, response.Response{
		Success: successCount > 0,
		Message: fmt.Sprintf("Procesados %d de %d códigos QR exitosamente", successCount, len(request.QRCodes)),
		Data:    resp,
	})
}

// =============================================
// NUEVOS ENDPOINTS ADICIONALES
// =============================================

// GenerateSuppliesFromBatch genera códigos QR para insumos individuales de un lote
func (c *QRController) GenerateSuppliesFromBatch(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad no implementada",
	})
}

// GetScanAnalytics obtiene resumen de actividad de escaneos por período
func (c *QRController) GetScanAnalytics(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de analytics no implementada",
	})
}

// GetTopScannedQRs obtiene top QRs más escaneados
func (c *QRController) GetTopScannedQRs(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de analytics no implementada",
	})
}

// GetUserScanActivity obtiene actividad por usuario
func (c *QRController) GetUserScanActivity(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de analytics no implementada",
	})
}

// GetPavilionScanActivity obtiene actividad por pabellón
func (c *QRController) GetPavilionScanActivity(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de analytics no implementada",
	})
}

// GetMovementPatterns obtiene patrones de movimiento
func (c *QRController) GetMovementPatterns(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de analytics no implementada",
	})
}

// CleanupOldScanEvents limpia eventos de escaneo antiguos
func (c *QRController) CleanupOldScanEvents(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de limpieza no implementada",
	})
}

// GetSystemQRStats obtiene estadísticas del sistema de QR
func (c *QRController) GetSystemQRStats(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de stats no implementada",
	})
}

// ExportTraceabilityData exporta datos de trazabilidad
func (c *QRController) ExportTraceabilityData(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de export no implementada",
	})
}

// VerifyDataIntegrity verifica integridad de datos de trazabilidad
func (c *QRController) VerifyDataIntegrity(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de verificación no implementada",
	})
}

// HandleScanWebhook maneja webhook para sistemas externos cuando se escanea un QR
func (c *QRController) HandleScanWebhook(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de webhook no implementada",
	})
}

// HandleConsumeWebhook maneja webhook para notificaciones de consumo
func (c *QRController) HandleConsumeWebhook(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de webhook no implementada",
	})
}

// =============================================
// MÉTODOS AUXILIARES
// =============================================

// stringPtr helper para crear punteros a string
func stringPtr(s string) *string {
	return &s
}
