package controllers

import (
	"meditrack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QRController struct {
	qrService services.QRService
}

func NewQRController(qrService services.QRService) *QRController {
	return &QRController{qrService: qrService}
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

// GetSupplyHistory obtiene el historial de un insumo por código QR
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
			Error:   "Error al obtener historial: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Historial obtenido exitosamente",
		Data:    history,
	})
}

// GenerateBatchQR genera un nuevo código QR para lote con imagen
func (c *QRController) GenerateBatchQR(ctx *gin.Context) {
	qrResponse, err := c.qrService.GenerateBatchQRCode()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al generar código QR: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Código QR de lote generado exitosamente",
		Data:    qrResponse,
	})
}

// GenerateSupplyQR genera un nuevo código QR para insumo médico con imagen
func (c *QRController) GenerateSupplyQR(ctx *gin.Context) {
	qrResponse, err := c.qrService.GenerateMedicalSupplyQRCode()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al generar código QR: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Código QR de insumo generado exitosamente",
		Data:    qrResponse,
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
			Error:   "Error al generar imagen QR: " + err.Error(),
		})
		return
	}

	// Configurar headers para imagen PNG
	ctx.Header("Content-Type", "image/png")
	ctx.Header("Content-Disposition", "inline; filename=\"qr_"+qrCode+".png\"")
	ctx.Header("Cache-Control", "public, max-age=3600") // Cache por 1 hora

	// Enviar imagen
	ctx.Data(http.StatusOK, "image/png", imageBytes)
}

// DownloadQRImage permite descargar la imagen QR
func (c *QRController) DownloadQRImage(ctx *gin.Context) {
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
			Error:   "Error al generar imagen QR: " + err.Error(),
		})
		return
	}

	// Configurar headers para descarga
	ctx.Header("Content-Type", "image/png")
	ctx.Header("Content-Disposition", "attachment; filename=\"qr_"+qrCode+".png\"")

	// Enviar imagen para descarga
	ctx.Data(http.StatusOK, "image/png", imageBytes)
}
