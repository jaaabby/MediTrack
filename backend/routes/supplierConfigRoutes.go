package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupSupplierConfigRoutes configura las rutas de configuración de proveedores
func SetupSupplierConfigRoutes(router *gin.RouterGroup, supplierConfigService services.SupplierConfigService) {
	supplierConfigController := controllers.NewSupplierConfigController(supplierConfigService)
	supplierConfigs := router.Group("/supplier-configs")
	{
		// Rutas CRUD
		supplierConfigs.POST("/", supplierConfigController.CreateSupplierConfig)
		supplierConfigs.GET("/", supplierConfigController.GetAllSupplierConfigs)
		supplierConfigs.GET("/:name", supplierConfigController.GetSupplierConfig)
		supplierConfigs.PUT("/:name", supplierConfigController.UpdateSupplierConfig)
		supplierConfigs.DELETE("/:name", supplierConfigController.DeleteSupplierConfig)
	}
}

