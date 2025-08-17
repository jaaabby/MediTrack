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
	supplyMovementRepo := repository.NewSupplyMovementRepository()
	operatingRoomRepo := repository.NewOperatingRoomRepository()
	locationRepo := repository.NewLocationRepository()
	supplyRouteRepo := repository.NewSupplyRouteRepository()
	doctorSupplyRouteRepo := repository.NewDoctorSupplyRouteRepository()
	supplyRouteMedicalSupplyRepo := repository.NewSupplyRouteMedicalSupplyRepository()

	// Crear servicios
	userService := services.NewUserService(userRepo)
	medicalSupplyService := services.NewMedicalSupplyService(medicalSupplyRepo)
	supplyMovementService := services.NewSupplyMovementService(supplyMovementRepo)
	operatingRoomService := services.NewOperatingRoomService(operatingRoomRepo)
	locationService := services.NewLocationService(locationRepo)
	supplyRouteService := services.NewSupplyRouteService(supplyRouteRepo)
	doctorSupplyRouteService := services.NewDoctorSupplyRouteService(doctorSupplyRouteRepo)
	supplyRouteMedicalSupplyService := services.NewSupplyRouteMedicalSupplyService(supplyRouteMedicalSupplyRepo)

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
		supplyMovementService,
		operatingRoomService,
		userService, // doctorService ahora es userService
		locationService,
		supplyRouteService,
		doctorSupplyRouteService,
		supplyRouteMedicalSupplyService,
	)

	// Iniciar servidor
	log.Printf("Servidor iniciando en puerto %d", cfg.Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), router); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
