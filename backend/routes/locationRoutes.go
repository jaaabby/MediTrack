package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupLocationRoutes configura las rutas de ubicaciones
func SetupLocationRoutes(router *gin.RouterGroup, locationService services.LocationService) {
	locationController := controllers.NewLocationController(locationService)

	locations := router.Group("/locations")
	{
		locations.POST("/", locationController.CreateLocation)
		locations.GET("/", locationController.GetAllLocations)
		locations.GET("/:id", locationController.GetLocationByID)
		locations.PUT("/:id", locationController.UpdateLocation)
		locations.DELETE("/:id", locationController.DeleteLocation)
	}
}
