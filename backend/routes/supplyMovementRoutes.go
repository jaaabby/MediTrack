package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupSupplyMovementRoutes configura las rutas de movimientos de insumos
func SetupSupplyMovementRoutes(router *gin.RouterGroup, supplyMovementService services.SupplyMovementService) {
	supplyMovementController := controllers.NewSupplyMovementController(supplyMovementService)

	supplyMovements := router.Group("/supply-movements")
	{
		supplyMovements.POST("/", supplyMovementController.CreateSupplyMovement)
		supplyMovements.GET("/", supplyMovementController.GetAllSupplyMovements)
		supplyMovements.GET("/:id", supplyMovementController.GetSupplyMovementByID)
		supplyMovements.GET("/status", supplyMovementController.GetSupplyMovementsByStatus)
		supplyMovements.PUT("/:id", supplyMovementController.UpdateSupplyMovement)
		supplyMovements.DELETE("/:id", supplyMovementController.DeleteSupplyMovement)
	}
}
