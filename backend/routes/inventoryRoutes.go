package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupInventoryRoutes configura las rutas de inventario
func SetupInventoryRoutes(router *gin.RouterGroup, inventoryService *services.InventoryService) {
	inventoryController := controllers.NewInventoryController(inventoryService)
	inventory := router.Group("/inventory")
	{
		// Inventario de bodega
		inventory.GET("/store", inventoryController.GetStoreInventory)

		// Inventario de pabellón
		inventory.GET("/pavilion/:pavilion_id", inventoryController.GetPavilionInventory)

		// Historial de movimientos
		inventory.GET("/movements", inventoryController.GetMovementHistory)

		// Resumen y sincronización
		inventory.GET("/summary", inventoryController.GetInventorySummary)
		inventory.POST("/sync", inventoryController.SyncInventory)

		// Inventario por tipo de cirugía
		inventory.GET("/by-surgery", inventoryController.GetInventoryBySurgeryType)

		// Reportes
		inventory.GET("/reports/transfers", inventoryController.GetTransferReport)
	}
}
