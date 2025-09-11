package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupMedicalSupplyRoutes configura las rutas de insumos médicos híbridas
func SetupMedicalSupplyRoutes(router *gin.RouterGroup, medicalSupplyService services.MedicalSupplyService) {
	medicalSupplyController := controllers.NewMedicalSupplyController(medicalSupplyService)
	medicalSupplies := router.Group("/medical-supplies")
	{
		// ===== RUTAS BÁSICAS CRUD =====
		medicalSupplies.POST("/", medicalSupplyController.CreateMedicalSupply)
		medicalSupplies.GET("/", medicalSupplyController.GetAllMedicalSupplies)
		medicalSupplies.GET("/:id", medicalSupplyController.GetMedicalSupplyByID)
		medicalSupplies.PUT("/:id", medicalSupplyController.UpdateMedicalSupply)
		medicalSupplies.DELETE("/:id", medicalSupplyController.DeleteMedicalSupply)

		// ===== INVENTARIO (DOBLE ENDPOINT PARA COMPATIBILIDAD) =====

		// Inventario básico (versión anterior - FUNCIONAL)
		medicalSupplies.GET("/inventory", medicalSupplyController.GetInventoryList)

		// Inventario avanzado (versión nueva - para funcionalidades complejas)
		medicalSupplies.GET("/inventory/advanced", medicalSupplyController.GetInventoryListAdvanced)

		// ===== RUTAS QR Y FUNCIONALIDADES AVANZADAS =====
		medicalSupplies.GET("/qr/:qrcode", medicalSupplyController.GetMedicalSupplyByQR)
		medicalSupplies.GET("/details/:qrcode", medicalSupplyController.GetSupplyWithBatchInfo)
		medicalSupplies.GET("/code/:code", medicalSupplyController.GetIndividualSuppliesByCode)
		medicalSupplies.GET("/batch/:batch_id/available", medicalSupplyController.GetAvailableSuppliesByBatch)

		// Rutas para crear múltiples insumos y consumo
		medicalSupplies.POST("/create-multiple", medicalSupplyController.CreateMultipleSupplies)
		medicalSupplies.POST("/consume", medicalSupplyController.ConsumeSupply)
		medicalSupplies.POST("/sync-amounts", medicalSupplyController.SyncBatchAmounts)

		// ===== RUTAS DE ALERTA PARA INSUMOS NO CONSUMIDOS =====
		medicalSupplies.GET("/unconsumed", medicalSupplyController.GetUnconsumedSupplies)
		medicalSupplies.POST("/check-unconsumed", medicalSupplyController.CheckUnconsumedSupplies)
	}
}
