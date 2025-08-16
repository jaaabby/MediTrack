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
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Conectar a la base de datos
	db, err := config.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer config.Close(db)

	// Crear repositorios
	userRepo := repository.NewUserRepository(db)
	medicalSupplyRepo := repository.NewMedicalSupplyRepository(db)
	supplyMovementRepo := repository.NewSupplyMovementRepository(db)
	operatingRoomRepo := repository.NewOperatingRoomRepository(db)
	doctorRepo := repository.NewDoctorRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	supplyRouteRepo := repository.NewSupplyRouteRepository(db)
	doctorSupplyRouteRepo := repository.NewDoctorSupplyRouteRepository(db)
	supplyRouteMedicalSupplyRepo := repository.NewSupplyRouteMedicalSupplyRepository(db)

	// Crear servicios
	userService := services.NewUserService(userRepo)
	medicalSupplyService := services.NewMedicalSupplyService(medicalSupplyRepo)
	supplyMovementService := services.NewSupplyMovementService(supplyMovementRepo)
	operatingRoomService := services.NewOperatingRoomService(operatingRoomRepo)
	doctorService := services.NewDoctorService(doctorRepo)
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
		doctorService,
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
