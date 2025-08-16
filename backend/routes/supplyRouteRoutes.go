package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupSupplyRouteRoutes configura las rutas de rutas de suministro
func SetupSupplyRouteRoutes(router *gin.RouterGroup, supplyRouteService services.SupplyRouteService) {
	supplyRouteController := controllers.NewSupplyRouteController(supplyRouteService)

	supplyRoutes := router.Group("/supply-routes")
	{
		supplyRoutes.POST("/", supplyRouteController.CreateSupplyRoute)
		supplyRoutes.GET("/", supplyRouteController.GetAllSupplyRoutes)
		supplyRoutes.GET("/:id", supplyRouteController.GetSupplyRouteByID)
		supplyRoutes.GET("/patient", supplyRouteController.GetSupplyRoutesByPatientID)
		supplyRoutes.GET("/operating-room", supplyRouteController.GetSupplyRoutesByOperatingRoomID)
		supplyRoutes.PUT("/:id", supplyRouteController.UpdateSupplyRoute)
		supplyRoutes.DELETE("/:id", supplyRouteController.DeleteSupplyRoute)
	}
}
