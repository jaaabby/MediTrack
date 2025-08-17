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
	supplyMovementService services.SupplyMovementService,
	operatingRoomService services.OperatingRoomService,
	doctorService services.UserService,
	locationService services.LocationService,
	supplyRouteService services.SupplyRouteService,
	doctorSupplyRouteService services.DoctorSupplyRouteService,
	supplyRouteMedicalSupplyService services.SupplyRouteMedicalSupplyService,
) {
	// API v1
	v1 := router.Group("/api/v1")
	{
		// Configurar rutas de usuarios
		SetupUserRoutes(v1, userService)

		// Configurar rutas de insumos médicos
		SetupMedicalSupplyRoutes(v1, medicalSupplyService)

		// Configurar rutas de movimientos de insumos
		SetupSupplyMovementRoutes(v1, supplyMovementService)

		// Configurar rutas de salas de operación
		SetupOperatingRoomRoutes(v1, operatingRoomService)

		// Configurar rutas de doctores
		SetupDoctorRoutes(v1, doctorService)

		// Configurar rutas de ubicaciones
		SetupLocationRoutes(v1, locationService)

		// Configurar rutas de rutas de suministro
		SetupSupplyRouteRoutes(v1, supplyRouteService)

		// Configurar rutas de relación doctor-ruta de suministro
		SetupDoctorSupplyRouteRoutes(v1, doctorSupplyRouteService)

		// Configurar rutas de relación ruta de suministro-insumo médico
		SetupSupplyRouteMedicalSupplyRoutes(v1, supplyRouteMedicalSupplyService)

		// Configurar rutas adicionales (trazabilidad, estadísticas, alertas, inventario)
		SetupAdditionalRoutes(v1)
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
