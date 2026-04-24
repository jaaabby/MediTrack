package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupSupplyHistoryRoutes configura las rutas de supply history
func SetupSupplyHistoryRoutes(router *gin.RouterGroup, supplyHistoryService services.SupplyHistoryService) {
	supplyHistoryController := controllers.NewSupplyHistoryController(supplyHistoryService)
	supplyHistory := router.Group("/supply-history")
	{
		supplyHistory.POST("/", supplyHistoryController.CreateSupplyHistory)
		supplyHistory.GET("/", supplyHistoryController.GetAllSupplyHistory)
		supplyHistory.GET("/with-details", supplyHistoryController.GetAllSupplyHistoriesWithDetails)
		supplyHistory.GET("/consumption-stats", supplyHistoryController.GetConsumptionStatsBySurgery)
		supplyHistory.GET("/by-batch/:batchId", supplyHistoryController.GetSupplyHistoryByBatch)
		supplyHistory.GET("/:id", supplyHistoryController.GetSupplyHistoryByID)
		supplyHistory.PUT("/:id", supplyHistoryController.UpdateSupplyHistory)
		supplyHistory.DELETE("/:id", supplyHistoryController.DeleteSupplyHistory)
	}
}
