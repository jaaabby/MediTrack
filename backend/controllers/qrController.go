package controllers

import (
	"fmt"
	"meditrack/services"
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

// ScanQR obtiene toda la información de un código QR escaneado
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

	// Obtener información del QR usando el método existente
	result, err := c.qrService.ScanQRCode(decodedQR)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Código QR no encontrado: " + err.Error(),
		})
		return
	}

	// PRIORIZAR INSUMOS INDIVIDUALES
	qrType := determineQRType(decodedQR)

	// Convertir QRInfo a map para añadir información contextual
	resultMap := map[string]interface{}{
		"type":                 result["type"],
		"id":                   result["id"],
		"qr_code":              result["qr_code"],
		"batch_info":           result["batch_info"],
		"supply_info":          result["supply_info"],
		"supply_code":          result["supply_code"],
		"history":              result["history"],
		"is_consumed":          result["is_consumed"],
		"can_consume":          result["can_consume"],
		"batch_status":         result["batch_status"],
		"scan_timestamp":       time.Now(),
		"qr_type":              qrType,
		"is_individual_supply": qrType == "SUPPLY",
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

	// Procesar consumo
	// Use ConsumeSupplyByQR with QRConsumptionRequest
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
			// available_count and consumed_count not calculated here, as supplies is []MedicalSupply
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

	// Not implemented in QRService, return error
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

	// Obtener resolución solicitada (por defecto: normal)

	ctx.JSON(http.StatusNotImplemented, response.Response{
		Success: false,
		Error:   "Funcionalidad de imagen QR no implementada en QRService",
	})
	return
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
