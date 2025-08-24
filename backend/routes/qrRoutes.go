package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupQRRoutes configura las rutas de códigos QR con funcionalidades mejoradas
func SetupQRRoutes(router *gin.RouterGroup, qrService services.QRService, medicalSupplyService services.MedicalSupplyService) {
	qrController := controllers.NewQRController(qrService)
	qrController.SetMedicalSupplyService(medicalSupplyService)

	qr := router.Group("/qr")
	{
		// === RUTAS BÁSICAS DE QR ===

		// Escanear un código QR y obtener toda su información (ENFOCADO EN INSUMOS INDIVIDUALES)
		qr.GET("/scan/:qrcode", qrController.ScanQR)

		// Validar si un código QR es válido
		qr.GET("/validate/:qrcode", qrController.ValidateQR)

		// Obtener historial de un insumo por código QR
		qr.GET("/history/:qrcode", qrController.GetSupplyHistory)

		// === GENERACIÓN DE QR CODES ===

		// Generar códigos QR con imagen
		qr.POST("/generate/batch", qrController.GenerateBatchQR)
		qr.POST("/generate/supply", qrController.GenerateSupplyQR)

		// === IMÁGENES QR ===

		// Servir imágenes QR (para mostrar en la UI)
		qr.GET("/image/:qrcode", qrController.GetQRImage)

		// Descargar imágenes QR con diferentes resoluciones
		// Query param: resolution=normal|high (default: normal)
		qr.GET("/download/:qrcode", qrController.DownloadQRImage)

		// === FUNCIONALIDADES DE CONSUMO ===

		// Consumir un insumo por QR (actualiza automáticamente las cantidades del lote)
		qr.POST("/consume", qrController.ConsumeSupply)

		// Consumir un insumo individual específico (NUEVA FUNCIONALIDAD ENFOCADA EN INSUMOS INDIVIDUALES)
		qr.POST("/consume/individual", qrController.ConsumeIndividualSupply)

		// Consumir múltiples insumos en lote
		qr.POST("/consume/bulk", qrController.BulkConsumeSupplies)

		// Verificar disponibilidad de un insumo para consumo
		qr.GET("/verify/:qrcode", qrController.VerifySupplyAvailability)

		// === INFORMACIÓN DETALLADA DE INSUMOS INDIVIDUALES ===

		// Obtener información detallada de un insumo con datos del lote
		qr.GET("/details/:qrcode", qrController.GetSupplyDetails)

		// === ADMINISTRACIÓN Y ESTADÍSTICAS ===

		// Sincronizar cantidades de lotes con productos individuales
		qr.POST("/sync/batch-amounts", qrController.SyncBatchAmounts)

		// Obtener estadísticas generales de uso de QR codes
		qr.GET("/stats", qrController.GetQRStats)

		// === RUTAS ADICIONALES PARA MEJOR ORGANIZACIÓN ===

		// Grupo de rutas específicas para insumos individuales
		individual := qr.Group("/individual")
		{
			// Consumir insumo individual
			individual.POST("/consume", qrController.ConsumeIndividualSupply)

			// Obtener información de insumo individual
			individual.GET("/:qrcode", qrController.ScanQR)

			// Verificar disponibilidad de insumo individual
			individual.GET("/:qrcode/availability", qrController.VerifySupplyAvailability)

			// Obtener historial de insumo individual
			individual.GET("/:qrcode/history", qrController.GetSupplyHistory)
		}

		// Grupo de rutas específicas para lotes
		batch := qr.Group("/batch")
		{
			// Obtener insumos individuales de un lote
			batch.GET("/:batch_id/supplies", qrController.GetIndividualSuppliesByBatch)

			// Generar QR de lote
			batch.POST("/generate", qrController.GenerateBatchQR)
		}

		// Grupo de rutas para imágenes
		images := qr.Group("/images")
		{
			// Servir imagen QR
			images.GET("/:qrcode", qrController.GetQRImage)

			// Descargar imagen QR
			images.GET("/:qrcode/download", qrController.DownloadQRImage)
		}

		// Grupo de rutas para administración
		admin := qr.Group("/admin")
		{
			// Estadísticas generales
			admin.GET("/stats", qrController.GetQRStats)

			// Sincronizar cantidades
			admin.POST("/sync", qrController.SyncBatchAmounts)
		}
	}
}
