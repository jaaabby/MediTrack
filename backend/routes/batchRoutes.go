package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupBatchRoutes configura las rutas de batch
func SetupBatchRoutes(router *gin.RouterGroup, batchService services.BatchService) {
	batchController := controllers.NewBatchController(batchService)
	batches := router.Group("/batches")
	{
		batches.POST("/", batchController.CreateBatch)
		batches.GET("/", batchController.GetAllBatches)
		batches.GET("/:id", batchController.GetBatchByID)
		batches.PUT("/:id", batchController.UpdateBatch)
		batches.DELETE("/:id", batchController.DeleteBatch)
	}
}
