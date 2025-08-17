package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupStoreRoutes configura las rutas de store
func SetupStoreRoutes(router *gin.RouterGroup, storeService services.StoreService) {
	storeController := controllers.NewStoreController(storeService)
	stores := router.Group("/stores")
	{
		stores.POST("/", storeController.CreateStore)
		stores.GET("/", storeController.GetAllStores)
		stores.GET("/:id", storeController.GetStoreByID)
		stores.PUT("/:id", storeController.UpdateStore)
		stores.DELETE("/:id", storeController.DeleteStore)
	}
}
