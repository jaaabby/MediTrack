package routes

import (
	"meditrack/controllers"
	"meditrack/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAutomaticConsumptionRoutes configura las rutas de consumo automático
func SetupAutomaticConsumptionRoutes(router *gin.Engine, automaticConsumptionController *controllers.AutomaticConsumptionController, secretKey string) {
	automaticConsumption := router.Group("/automatic-consumption")
	automaticConsumption.Use(middleware.AuthMiddleware(secretKey))
	{
		// Ejecutar consumo automático manualmente (solo para administradores)
		automaticConsumption.POST("/process", automaticConsumptionController.ProcessAutomaticConsumption)
	}
}

