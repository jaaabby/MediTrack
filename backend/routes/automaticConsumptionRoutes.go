package routes

import (
	"meditrack/controllers"
	"meditrack/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupAutomaticConsumptionRoutes configura las rutas de consumo automático
func SetupAutomaticConsumptionRoutes(router *gin.Engine, automaticConsumptionController *controllers.AutomaticConsumptionController, secretKey string, db *gorm.DB) {
	automaticConsumption := router.Group("/automatic-consumption")
	automaticConsumption.Use(middleware.AuthMiddleware(secretKey, db))
	{
		// Ejecutar consumo automático manualmente (solo para administradores)
		automaticConsumption.POST("/process", automaticConsumptionController.ProcessAutomaticConsumption)
	}
}

