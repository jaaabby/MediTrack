package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

func SetupSupplyCodeRoutes(router *gin.RouterGroup, supplyCodeService services.SupplyCodeService) {
	supplyCodeController := controllers.NewSupplyCodeController(supplyCodeService)
	supplyCodes := router.Group("/supply-codes")
	{
		supplyCodes.POST("/", supplyCodeController.CreateSupplyCode)
		supplyCodes.GET("/", supplyCodeController.GetAllSupplyCodes)
		supplyCodes.GET("/:id", supplyCodeController.GetSupplyCodeByID)
		supplyCodes.PUT("/:id", supplyCodeController.UpdateSupplyCode)
		supplyCodes.DELETE("/:id", supplyCodeController.DeleteSupplyCode)
	}
}
