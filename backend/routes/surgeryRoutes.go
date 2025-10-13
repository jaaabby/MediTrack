package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupSurgeryRoutes configura las rutas de tipos de cirugía
func SetupSurgeryRoutes(router *gin.RouterGroup, surgeryService *services.SurgeryService) {
	surgeryController := controllers.NewSurgeryController(surgeryService)
	surgeries := router.Group("/surgeries")
	{
		// CRUD básico
		surgeries.POST("/", surgeryController.CreateSurgery)
		surgeries.GET("/", surgeryController.GetSurgeriesPaginated)

		// Rutas específicas ANTES de las rutas con parámetros
		surgeries.GET("/all", surgeryController.GetAllSurgeries)
		surgeries.GET("/search", surgeryController.SearchSurgeries)

		// Rutas con parámetros al final
		surgeries.GET("/:id", surgeryController.GetSurgeryByID)
		surgeries.PUT("/:id", surgeryController.UpdateSurgery)
		surgeries.DELETE("/:id", surgeryController.DeleteSurgery)
	}
}
