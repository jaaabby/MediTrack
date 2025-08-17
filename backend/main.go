package main

import (
	"fmt"
	"log"
	"net/http"

	"meditrack/config"
	"meditrack/middleware"
	"meditrack/repository"
	"meditrack/routes"
	"meditrack/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Conectar a la base de datos con GORM
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos con GORM: %v", err)
	}

	// Crear repositorios
	userRepo := repository.NewUserRepository(db)
	medicalSupplyRepo := repository.NewMedicalSupplyRepository(db)
	medicalCenterRepo := repository.NewMedicalCenterRepository(db)
	batchRepo := repository.NewBatchRepository(db)
	storeRepo := repository.NewStoreRepository(db)
	supplyHistoryRepo := repository.NewSupplyHistoryRepository(db)

	// Crear servicios
	userService := services.NewUserService(userRepo)
	medicalSupplyService := services.NewMedicalSupplyService(medicalSupplyRepo)
	medicalCenterService := services.NewMedicalCenterService(medicalCenterRepo)
	batchService := services.NewBatchService(batchRepo)
	pavilionService := services.NewPavilionService(db)
	storeService := services.NewStoreService(storeRepo)
	supplyHistoryService := services.NewSupplyHistoryService(supplyHistoryRepo)

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
		userService,
		medicalSupplyService,
		medicalCenterService,
		batchService,
		pavilionService,
		storeService,
		supplyHistoryService,
		userService,
	)

	// Iniciar servidor
	log.Printf("Servidor iniciando en puerto %d", cfg.Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), router); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
