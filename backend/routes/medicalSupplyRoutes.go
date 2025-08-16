package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupMedicalSupplyRoutes configura las rutas de insumos médicos
func SetupMedicalSupplyRoutes(router *gin.RouterGroup, medicalSupplyService services.MedicalSupplyService) {
	medicalSupplyController := controllers.NewMedicalSupplyController(medicalSupplyService)

	medicalSupplies := router.Group("/medical-supplies")
	{
		medicalSupplies.POST("/", medicalSupplyController.CreateMedicalSupply)
		medicalSupplies.GET("/", medicalSupplyController.GetAllMedicalSupplies)
		medicalSupplies.GET("/:id", medicalSupplyController.GetMedicalSupplyByID)
		medicalSupplies.GET("/qr", medicalSupplyController.GetMedicalSupplyByQRCode)
		medicalSupplies.PUT("/:id", medicalSupplyController.UpdateMedicalSupply)
		medicalSupplies.DELETE("/:id", medicalSupplyController.DeleteMedicalSupply)
	}
}
