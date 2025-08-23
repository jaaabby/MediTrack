package main

import (
	"fmt"
	"log"
	"net/http"

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

	// Inicializar el mailer
	mailer.Setup()

	// Conectar a la base de datos con GORM
	db, err := config.ConnectGORM(cfg.Database)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos con GORM: %v", err)
	}

	// Crear servicios
	userService := services.NewUserService(db)
	medicalSupplyService := services.NewMedicalSupplyService(db)
	medicalCenterService := services.NewMedicalCenterService(db)
	batchService := services.NewBatchService(db)
	pavilionService := services.NewPavilionService(db)
	storeService := services.NewStoreService(db)
	supplyHistoryService := services.NewSupplyHistoryService(db)
	supplyCodeService := services.NewSupplyCodeService(db)
	batchHistoryService := services.NewBatchHistoryService(db)

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
		*batchHistoryService,
	)

	// Iniciar servidor
	log.Printf("Servidor iniciando en puerto %d", cfg.Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), router); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
