package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupSurgeryTypicalSupplyRoutes configura las rutas de insumos típicos por cirugía
func SetupSurgeryTypicalSupplyRoutes(router *gin.RouterGroup, typicalSupplyService *services.SurgeryTypicalSupplyService) {
	typicalSupplyController := controllers.NewSurgeryTypicalSupplyController(typicalSupplyService)
	typicalSupplies := router.Group("/surgery-typical-supplies")
	{
		// Rutas específicas primero (antes de :id para evitar conflictos)
		typicalSupplies.GET("/count", typicalSupplyController.GetTypicalSuppliesCount)
		typicalSupplies.GET("/surgery/:surgery_id", typicalSupplyController.GetTypicalSuppliesBySurgeryID)
		typicalSupplies.GET("/supply/:supply_code", typicalSupplyController.GetSurgeriesBySupplyCode)
		typicalSupplies.POST("/surgery/:surgery_id/bulk", typicalSupplyController.BulkCreateSurgeryTypicalSupplies)
		
		// CRUD básico
		typicalSupplies.POST("/", typicalSupplyController.CreateSurgeryTypicalSupply)
		typicalSupplies.GET("/:id", typicalSupplyController.GetSurgeryTypicalSupplyByID)
		typicalSupplies.PUT("/:id", typicalSupplyController.UpdateSurgeryTypicalSupply)
		typicalSupplies.DELETE("/:id", typicalSupplyController.DeleteSurgeryTypicalSupply)
	}
}

