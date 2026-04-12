package routes

import (
	"meditrack/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(
	router *gin.Engine,
	db *gorm.DB,
	userService services.UserService,
	medicalSupplyService services.MedicalSupplyService,
	medicalCenterService services.MedicalCenterService,
	batchService services.BatchService,
	pavilionService services.PavilionService,
	storeService services.StoreService,
	supplyHistoryService services.SupplyHistoryService,
	supplyCodeService services.SupplyCodeService,
	qrService services.QRService,
	batchHistoryService services.BatchHistoryService,
	transferService *services.SupplyTransferService,
	inventoryService *services.InventoryService,
	surgeryService *services.SurgeryService,
	specialtyService *services.MedicalSpecialtyService,
	typicalSupplyService *services.SurgeryTypicalSupplyService,
	doctorInfoService *services.DoctorInfoService,
	secretKey string,
	webauthnService *services.WebAuthnService,
) {
	// API v1
	v1 := router.Group("/api/v1")
	{
		// Configurar rutas de autenticación
		SetupAuthRoutes(v1, userService, secretKey, db, webauthnService)

		// Configurar rutas de usuarios
		SetupUserRoutes(v1, userService, secretKey, db)

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
		// Obtener cartService desde main.go o crear aquí
		cartService := services.NewCartService(db)
		SetupQRRoutes(v1, qrService, medicalSupplyService, cartService)

		// Configurar rutas de historial de lotes
		SetupBatchHistoryRoutes(v1, batchHistoryService)

		// Configurar rutas de transferencias
		SetupSupplyTransferRoutes(v1, transferService, secretKey, db)

		// Configurar rutas de inventario
		SetupInventoryRoutes(v1, inventoryService)

		// Configurar rutas de tipos de cirugía
		SetupSurgeryRoutes(v1, surgeryService)

		// Configurar rutas de especialidades médicas
		SetupMedicalSpecialtyRoutes(v1, specialtyService)

		// Configurar rutas de insumos típicos por cirugía
		SetupSurgeryTypicalSupplyRoutes(v1, typicalSupplyService)

		// Configurar rutas de información de doctores
		SetupDoctorInfoRoutes(v1, doctorInfoService)

		// Configurar rutas de configuración de proveedores
		supplierConfigService := services.NewSupplierConfigService(db)
		SetupSupplierConfigRoutes(v1, *supplierConfigService)
	}

	// Ruta de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "Meditrack API",
			"version": "1.0.0",
		})
	})

	// Ruta raíz
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Bienvenido a Meditrack API",
			"version": "1.0.0",
			"docs":    "/api/v1/docs",
		})
	})
}
