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
		// Rutas básicas CRUD
		medicalSupplies.POST("/", medicalSupplyController.CreateMedicalSupply)
		medicalSupplies.GET("/", medicalSupplyController.GetAllMedicalSupplies)
		medicalSupplies.GET("/:id", medicalSupplyController.GetMedicalSupplyByID)
		medicalSupplies.PUT("/:id", medicalSupplyController.UpdateMedicalSupply)
		medicalSupplies.DELETE("/:id", medicalSupplyController.DeleteMedicalSupply)

		// Rutas mejoradas para QR y gestión individual
		medicalSupplies.GET("/qr/:qrcode", medicalSupplyController.GetMedicalSupplyByQR)
		medicalSupplies.GET("/details/:qrcode", medicalSupplyController.GetSupplyWithBatchInfo)
		medicalSupplies.GET("/inventory", medicalSupplyController.GetInventoryList)
		medicalSupplies.GET("/code/:code", medicalSupplyController.GetIndividualSuppliesByCode)
		medicalSupplies.GET("/batch/:batch_id/available", medicalSupplyController.GetAvailableSuppliesByBatch)

		// Rutas para crear múltiples insumos y consumo
		medicalSupplies.POST("/create-multiple", medicalSupplyController.CreateMultipleSupplies)
		medicalSupplies.POST("/consume", medicalSupplyController.ConsumeSupply)
		medicalSupplies.POST("/sync-amounts", medicalSupplyController.SyncBatchAmounts)
	}
}
