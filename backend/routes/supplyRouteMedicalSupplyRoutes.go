package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupSupplyRouteMedicalSupplyRoutes configura las rutas de la relación ruta de suministro-insumo médico
func SetupSupplyRouteMedicalSupplyRoutes(router *gin.RouterGroup, supplyRouteMedicalSupplyService services.SupplyRouteMedicalSupplyService) {
	supplyRouteMedicalSupplyController := controllers.NewSupplyRouteMedicalSupplyController(supplyRouteMedicalSupplyService)

	supplyRouteMedicalSupplies := router.Group("/supply-route-medical-supplies")
	{
		supplyRouteMedicalSupplies.POST("/", supplyRouteMedicalSupplyController.CreateSupplyRouteMedicalSupply)
		supplyRouteMedicalSupplies.GET("/", supplyRouteMedicalSupplyController.GetAllSupplyRouteMedicalSupplies)
		supplyRouteMedicalSupplies.GET("/:id", supplyRouteMedicalSupplyController.GetSupplyRouteMedicalSupplyByID)
		supplyRouteMedicalSupplies.GET("/supply-route", supplyRouteMedicalSupplyController.GetSupplyRouteMedicalSuppliesBySupplyRouteID)
		supplyRouteMedicalSupplies.GET("/medical-supply", supplyRouteMedicalSupplyController.GetSupplyRouteMedicalSuppliesByMedicalSupplyID)
		supplyRouteMedicalSupplies.PUT("/:id", supplyRouteMedicalSupplyController.UpdateSupplyRouteMedicalSupply)
		supplyRouteMedicalSupplies.DELETE("/:id", supplyRouteMedicalSupplyController.DeleteSupplyRouteMedicalSupply)
	}
}
