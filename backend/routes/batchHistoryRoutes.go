package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupBatchHistoryRoutes configura las rutas de historial de lotes
func SetupBatchHistoryRoutes(router *gin.RouterGroup, batchHistoryService services.BatchHistoryService) {
	batchHistoryController := controllers.NewBatchHistoryController(batchHistoryService)
	batchHistory := router.Group("/batch-history")
	{
		batchHistory.POST("/", batchHistoryController.CreateBatchHistory)
		batchHistory.GET("/", batchHistoryController.GetAllBatchHistory)
		batchHistory.GET("/:id", batchHistoryController.GetBatchHistoryByID)
		batchHistory.PUT("/:id", batchHistoryController.UpdateBatchHistory)
		batchHistory.DELETE("/:id", batchHistoryController.DeleteBatchHistory)
		batchHistory.GET("/search/:batchNumber", batchHistoryController.SearchBatchHistoryByBatchNumber)
		batchHistory.GET("/details", batchHistoryController.GetAllBatchHistoriesWithDetails)
	}
}
