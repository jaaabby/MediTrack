package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupQRRoutes configura las rutas de códigos QR
func SetupQRRoutes(router *gin.RouterGroup, qrService services.QRService) {
	qrController := controllers.NewQRController(qrService)

	qr := router.Group("/qr")
	{
		// Escanear un código QR y obtener toda su información
		qr.GET("/scan/:qrcode", qrController.ScanQR)

		// Validar si un código QR es válido
		qr.GET("/validate/:qrcode", qrController.ValidateQR)

		// Obtener historial de un insumo por código QR
		qr.GET("/history/:qrcode", qrController.GetSupplyHistory)

		// Generar códigos QR con imagen
		qr.POST("/generate/batch", qrController.GenerateBatchQR)
		qr.POST("/generate/supply", qrController.GenerateSupplyQR)

		// Servir imágenes QR
		qr.GET("/image/:qrcode", qrController.GetQRImage)
		qr.GET("/download/:qrcode", qrController.DownloadQRImage)
	}
}
