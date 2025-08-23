package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupSupplyHistoryRoutes configura las rutas de supply history
func SetupSupplyHistoryRoutes(router *gin.RouterGroup, supplyHistoryService services.SupplyHistoryService) {
	supplyHistoryController := controllers.NewSupplyHistoryController(supplyHistoryService)
	histories := router.Group("/supply-histories")
	{
		histories.POST("/", supplyHistoryController.CreateSupplyHistory)
		histories.GET("/", supplyHistoryController.GetAllSupplyHistories)
		histories.GET("/:id", supplyHistoryController.GetSupplyHistoryByID)
		histories.DELETE("/:id", supplyHistoryController.DeleteSupplyHistory)
	}
}
