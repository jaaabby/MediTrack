package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"meditrack/config"
	"meditrack/controllers"
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

	// Crear servicios para el nuevo sistema de inventario
	transferService := services.NewSupplyTransferService(db)
	inventoryService := services.NewInventoryService(db)
	surgeryService := services.NewSurgeryService(db)

	// Crear servicios de configuración médica
	specialtyService := services.NewMedicalSpecialtyService(db)
	typicalSupplyService := services.NewSurgeryTypicalSupplyService(db)
	doctorInfoService := services.NewDoctorInfoService(db)

	// Configurar Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// Configurar rutas principales
	routes.SetupRoutes(
		router,
		db,
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
		transferService,
		inventoryService,
		surgeryService,
		specialtyService,
		typicalSupplyService,
		doctorInfoService,
		secretKey,
	)

	// Inicializar servicio y controlador de SupplyRequest
	supplyRequestService := services.NewSupplyRequestService(db)
	supplyRequestController := controllers.NewSupplyRequestController(supplyRequestService, qrService)

	// Inicializar servicio y controlador de Cart
	cartService := services.NewCartService(db)
	cartController := controllers.NewCartController(cartService)

	// Registrar rutas de supply requests y trazabilidad QR
	routes.SetupSupplyRequestRoutes(router, supplyRequestController)

	// Registrar rutas de carritos
	routes.SetupCartRoutes(router, cartController)

	// Iniciar el verificador automático de retornos a bodega en una goroutine
	go medicalSupplyService.StartAutomaticReturnChecker()
	log.Println("✅ Iniciado verificador automático de retornos a bodega")

	// Iniciar el verificador automático de stock bajo en una goroutine
	go batchService.StartAutomaticLowStockChecker()
	log.Println("✅ Iniciado verificador automático de stock bajo")

	// Iniciar servidores HTTP y HTTPS
	log.Printf("Servidor iniciando en puerto %d (HTTP)", cfg.Server.Port)

	// Verificar si existen los certificados SSL
	certFile := "certs/server.crt"
	keyFile := "certs/server.key"
	httpsPort := 8443

	// Verificar certificados y rutas en diferentes ubicaciones (Docker vs local)
	certPaths := []string{
		certFile,
		"/root/certs/server.crt",
		"./certs/server.crt",
	}
	keyPaths := []string{
		keyFile,
		"/root/certs/server.key",
		"./certs/server.key",
	}

	var foundCert, foundKey string
	for _, path := range certPaths {
		if _, err := os.Stat(path); err == nil {
			foundCert = path
			break
		}
	}
	for _, path := range keyPaths {
		if _, err := os.Stat(path); err == nil {
			foundKey = path
			break
		}
	}

	// Iniciar servidor HTTP siempre
	go func() {
		log.Println("🌐 Iniciando servidor HTTP en puerto", cfg.Server.Port)
		if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
			log.Fatalf("Error al iniciar servidor HTTP: %v", err)
		}
	}()

	// Iniciar servidor HTTPS si hay certificados
	if foundCert != "" && foundKey != "" {
		log.Println("✅ Certificados SSL encontrados")
		log.Printf("🔒 Iniciando servidor HTTPS en puerto %d", httpsPort)
		log.Printf("   Certificado: %s", foundCert)
		log.Printf("   Clave: %s", foundKey)

		// Usar http.ListenAndServeTLS para HTTPS en puerto separado
		go func() {
			httpsServer := &http.Server{
				Addr:    fmt.Sprintf(":%d", httpsPort),
				Handler: router,
			}
			if err := httpsServer.ListenAndServeTLS(foundCert, foundKey); err != nil {
				log.Printf("⚠️  Error al iniciar servidor HTTPS: %v (continuando solo con HTTP)", err)
			}
		}()
	} else {
		log.Println("⚠️  Certificados SSL no encontrados, solo usando HTTP")
		log.Println("   Para habilitar HTTPS, ejecuta: scripts/generate-certs.sh (o .bat)")
	}

	// Mantener el proceso corriendo
	select {}
}
