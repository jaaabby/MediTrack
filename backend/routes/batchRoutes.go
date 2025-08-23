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
		// Rutas básicas CRUD
		batches.POST("/", batchController.CreateBatch)
		batches.GET("/", batchController.GetAllBatches)
		batches.GET("/:id", batchController.GetBatchByID)
		batches.PUT("/:id", batchController.UpdateBatch)
		batches.DELETE("/:id", batchController.DeleteBatch)

		// Rutas mejoradas para manejo completo
		batches.POST("/create-with-supplies", batchController.CreateBatchWithIndividualSupplies)
		batches.GET("/:id/with-supplies", batchController.GetBatchWithSupplyInfo)
		batches.GET("/qr/:qrcode", batchController.GetBatchByQR)
		batches.GET("/sync/needed", batchController.GetBatchesNeedingSync)

		// Rutas para actualización y mantenimiento
		batches.PATCH("/:id/amount", batchController.UpdateBatchAmount)
		batches.POST("/sync/all", batchController.SyncAllBatchAmounts)

		// Rutas para alertas
		batches.POST("/:id/check-low-stock", batchController.CheckLowStockAlert)
		batches.POST("/:id/check-expiration", batchController.CheckExpirationAlert)
	}
}
