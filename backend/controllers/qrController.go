package controllers

import (
	"fmt"
	"meditrack/services"
	"net/http"
	"strconv"

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

// ScanQR obtiene toda la información de un código QR escaneado
func (c *QRController) ScanQR(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	info, err := c.qrService.GetQRInfo(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Código QR no encontrado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Información del código QR obtenida exitosamente",
		Data:    info,
	})
}

// ValidateQR valida si un código QR es válido
func (c *QRController) ValidateQR(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	isValid, qrType, err := c.qrService.ValidateQRCode(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
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
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	history, err := c.qrService.GetSupplyHistory(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Historial obtenido exitosamente",
		Data:    history,
	})
}

// GenerateBatchQR genera un código QR para un lote con imagen
func (c *QRController) GenerateBatchQR(ctx *gin.Context) {
	result, err := c.qrService.GenerateBatchQRCode()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error generando código QR de lote: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Código QR de lote generado exitosamente",
		Data:    result,
	})
}

// GenerateSupplyQR genera un código QR para un insumo médico con imagen
func (c *QRController) GenerateSupplyQR(ctx *gin.Context) {
	result, err := c.qrService.GenerateMedicalSupplyQRCode()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error generando código QR de insumo: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Código QR de insumo generado exitosamente",
		Data:    result,
	})
}

// ConsumeSupply consume un insumo por su código QR y actualiza automáticamente las cantidades del lote
func (c *QRController) ConsumeSupply(ctx *gin.Context) {
	var request services.QRConsumptionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de consumo inválidos: " + err.Error(),
		})
		return
	}

	// Validar destination_type
	if request.DestinationType != "pavilion" && request.DestinationType != "store" {
		ctx.JSON(http.StatusBadRequest, Response{
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

		ctx.JSON(statusCode, Response{
			Success: false,
			Error:   err.Error(),
			Data:    result,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: result.Message,
		Data:    result,
	})
}

// GetQRImage sirve la imagen QR como PNG
func (c *QRController) GetQRImage(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	imageBytes, err := c.qrService.GetQRImage(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Configurar headers para la imagen
	ctx.Header("Content-Type", "image/png")
	ctx.Header("Content-Length", strconv.Itoa(len(imageBytes)))
	ctx.Header("Cache-Control", "public, max-age=3600") // Cache por 1 hora

	ctx.Data(http.StatusOK, "image/png", imageBytes)
}

// DownloadQRImage permite descargar la imagen QR con un nombre específico
func (c *QRController) DownloadQRImage(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	// Obtener resolución solicitada (por defecto: normal)
	resolution := ctx.DefaultQuery("resolution", "normal")

	var imageBytes []byte
	var err error

	if resolution == "high" {
		imageBytes, err = c.qrService.GetQRImageHighRes(qrCode)
	} else {
		imageBytes, err = c.qrService.GetQRImage(qrCode)
	}

	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Configurar headers para descarga
	filename := qrCode + "_qr.png"
	if resolution == "high" {
		filename = qrCode + "_qr_hd.png"
	}

	ctx.Header("Content-Type", "image/png")
	ctx.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	ctx.Header("Content-Length", strconv.Itoa(len(imageBytes)))

	ctx.Data(http.StatusOK, "image/png", imageBytes)
}

// GetSupplyDetails obtiene información detallada de un insumo por su código QR
func (c *QRController) GetSupplyDetails(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	details, err := c.medicalSupplyService.GetSupplyWithBatchInfo(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Detalles del insumo obtenidos exitosamente",
		Data:    details,
	})
}

// SyncBatchAmounts sincroniza las cantidades de los lotes con los productos individuales
func (c *QRController) SyncBatchAmounts(ctx *gin.Context) {
	err := c.medicalSupplyService.SyncBatchAmounts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error sincronizando cantidades de lotes: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Cantidades de lotes sincronizadas exitosamente",
	})
}

// GetQRStats obtiene estadísticas generales de uso de QR codes
func (c *QRController) GetQRStats(ctx *gin.Context) {
	stats, err := c.qrService.GetQRStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error obteniendo estadísticas: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Estadísticas obtenidas exitosamente",
		Data:    stats,
	})
}

// VerifySupplyAvailability verifica si un insumo está disponible para consumo
func (c *QRController) VerifySupplyAvailability(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	info, err := c.qrService.GetQRInfo(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Solo verificar disponibilidad para insumos individuales
	if info.Type != "medical_supply" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "La verificación de disponibilidad solo aplica para insumos individuales",
		})
		return
	}

	response := map[string]interface{}{
		"qr_code":      qrCode,
		"is_available": info.CanConsume,
		"is_consumed":  info.IsConsumed,
		"supply_info":  info.SupplyInfo,
		"batch_status": info.BatchStatus,
	}

	if !info.CanConsume {
		response["reason"] = "Insumo ya consumido o lote sin stock"
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Verificación de disponibilidad completada",
		Data:    response,
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
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	// Validar destination_type
	if request.DestinationType != "pavilion" && request.DestinationType != "store" {
		ctx.JSON(http.StatusBadRequest, Response{
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

	response := map[string]interface{}{
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

	ctx.JSON(statusCode, Response{
		Success: successCount > 0,
		Message: fmt.Sprintf("Procesados %d de %d códigos QR exitosamente", successCount, len(request.QRCodes)),
		Data:    response,
	})
}
