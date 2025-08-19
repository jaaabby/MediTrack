package routes

import (
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(
	router *gin.Engine,
	userService services.UserService,
	medicalSupplyService services.MedicalSupplyService,
	medicalCenterService services.MedicalCenterService,
	batchService services.BatchService,
	pavilionService services.PavilionService,
	storeService services.StoreService,
	supplyHistoryService services.SupplyHistoryService,
	supplyCodeService services.SupplyCodeService,
	qrService services.QRService,
) {
	// API v1
	v1 := router.Group("/api/v1")
	{
		// Configurar rutas de usuarios
		SetupUserRoutes(v1, userService)

		// Configurar rutas de insumos médicos
		SetupMedicalSupplyRoutes(v1, medicalSupplyService)

		// Configurar rutas de centros médicos
		SetupMedicalCenterRoutes(v1, medicalCenterService)

		// Configurar rutas de lotes
		SetupBatchRoutes(v1, batchService)

		// Configurar rutas de pabellones
		SetupPavilionRoutes(v1, pavilionService)

		// Configurar rutas de bodegas
		SetupStoreRoutes(v1, storeService)

		// Configurar rutas de historial de insumos
		SetupSupplyHistoryRoutes(v1, supplyHistoryService)

		// Configurar rutas de códigos de insumos
		SetupSupplyCodeRoutes(v1, supplyCodeService)

		// Configurar rutas de códigos QR
		SetupQRRoutes(v1, qrService)
	}

	// Ruta de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "MediTrack API",
			"version": "1.0.0",
		})
	})

	// Ruta raíz
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Bienvenido a MediTrack API",
			"version": "1.0.0",
			"docs":    "/api/v1/docs",
		})
	})
}
