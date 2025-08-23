package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

func SetupBatchHistoryRoutes(router *gin.RouterGroup, batchHistoryService services.BatchHistoryService) {
	batchHistoryController := controllers.NewBatchHistoryController(batchHistoryService)
	histories := router.Group("/batch-histories")
	{
		histories.POST("/", batchHistoryController.CreateBatchHistory)
		histories.GET("/", batchHistoryController.GetAllBatchHistories)
		histories.GET("/:id", batchHistoryController.GetBatchHistoryByID)
		histories.DELETE("/:id", batchHistoryController.DeleteBatchHistory)
		histories.GET("/details", batchHistoryController.GetAllBatchHistoriesWithDetails)
		histories.GET("/details/paginated", batchHistoryController.GetAllBatchHistoriesWithDetailsPaginated)

		// Nuevas rutas para búsqueda por número de lote
		histories.GET("/batch/:batchNumber", batchHistoryController.GetBatchHistoryByBatchNumber)
		histories.GET("/search/:batchNumber", batchHistoryController.SearchBatchHistoryByBatchNumber)
	}
}
