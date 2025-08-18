package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupPavilionRoutes configura las rutas de pavilion
func SetupPavilionRoutes(router *gin.RouterGroup, pavilionService services.PavilionService) {
	pavilionController := controllers.NewPavilionController(pavilionService)
	pavilions := router.Group("/pavilions")
	{
		pavilions.POST("/", pavilionController.CreatePavilion)
		pavilions.GET("/", pavilionController.GetAllPavilions)
		pavilions.GET("/:id", pavilionController.GetPavilionByID)
		pavilions.PUT("/:id", pavilionController.UpdatePavilion)
		pavilions.DELETE("/:id", pavilionController.DeletePavilion)
	}
}
