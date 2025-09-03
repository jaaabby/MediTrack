package main

import (
	"fmt"
	"log"
	"os"

	"meditrack/config"
	"meditrack/mailer"
	"meditrack/middleware"
	"meditrack/routes"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Obtener secret key para JWT
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "default-secret-key-change-in-production"
		log.Println("ADVERTENCIA: Usando secret key por defecto. Configura JWT_SECRET_KEY en producción.")
	}

	// Inicializar el mailer
	mailer.Setup()

	// Conectar a la base de datos con GORM
	db, err := config.ConnectGORM(cfg.Database)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos con GORM: %v", err)
	}

	// Crear servicio QR primero (sin dependencias)
	qrService := services.NewQRService(db)

	// Crear servicios con dependencias de QR
	userService := services.NewUserService(db)
	medicalSupplyService := services.NewMedicalSupplyService(db, qrService)
	medicalCenterService := services.NewMedicalCenterService(db)
	batchHistoryService := services.NewBatchHistoryService(db)
	batchService := services.NewBatchService(db, qrService, medicalSupplyService, batchHistoryService)
	pavilionService := services.NewPavilionService(db)
	storeService := services.NewStoreService(db)
	supplyHistoryService := services.NewSupplyHistoryService(db)
	supplyCodeService := services.NewSupplyCodeService(db)

	batchService.SetBatchHistoryService(batchHistoryService)
	// Inicializar BatchHistoryService en BatchService
	// Configurar el servicio de lotes con el servicio de suministros médicos
	batchService.SetMedicalSupplyService(medicalSupplyService)
	// Configurar Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// Configurar rutas
	routes.SetupRoutes(
		router,
		*userService,
		*medicalSupplyService,
		*medicalCenterService,
		*batchService,
		*pavilionService,
		*storeService,
		*supplyHistoryService,
		*supplyCodeService,
		*qrService,
		*batchHistoryService,
		secretKey,
	)

	// Iniciar servidor correctamente con Gin
	log.Printf("Servidor iniciando en puerto %d", cfg.Server.Port)
	if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
