package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupMedicalSupplyRoutes configura las rutas de medical supply
func SetupMedicalSupplyRoutes(router *gin.RouterGroup, medicalSupplyService services.MedicalSupplyService) {
	medicalSupplyController := controllers.NewMedicalSupplyController(medicalSupplyService)
	supplies := router.Group("/medical-supplies")
	{
		supplies.POST("/", medicalSupplyController.CreateMedicalSupply)
		supplies.GET("/", medicalSupplyController.GetAllMedicalSupplies)
		supplies.GET("/:id", medicalSupplyController.GetMedicalSupplyByID)
		supplies.PUT("/:id", medicalSupplyController.UpdateMedicalSupply)
		supplies.DELETE("/:id", medicalSupplyController.DeleteMedicalSupply)
	}
}
