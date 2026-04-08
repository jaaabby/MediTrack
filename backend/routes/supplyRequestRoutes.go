package routes

import (
	"meditrack/controllers"
	"meditrack/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupSupplyRequestRoutes configura las rutas para solicitudes de insumo con trazabilidad QR
func SetupSupplyRequestRoutes(router *gin.Engine, supplyRequestController *controllers.SupplyRequestController, secretKey string, db *gorm.DB) {
	// Grupo de rutas para solicitudes de insumo
	supplyRequestGroup := router.Group("/api/v1/supply-requests")
	supplyRequestGroup.Use(middleware.AuthMiddleware(secretKey, db)) // Aplicar autenticación a todas las rutas
	{
		// CRUD básico de solicitudes
		supplyRequestGroup.POST("", supplyRequestController.CreateSupplyRequest)
		supplyRequestGroup.GET("", supplyRequestController.GetAllSupplyRequests)
		supplyRequestGroup.GET("/:id", supplyRequestController.GetSupplyRequestByID)
		supplyRequestGroup.GET("/:id/items", supplyRequestController.GetSupplyRequestItemsController)
		supplyRequestGroup.PUT("/:id/approve", supplyRequestController.ApproveSupplyRequest)
		supplyRequestGroup.PUT("/:id/reject", supplyRequestController.RejectSupplyRequest)
		supplyRequestGroup.PUT("/:id/complete", supplyRequestController.CompleteSupplyRequest)
		supplyRequestGroup.DELETE("/:id", supplyRequestController.DeleteSupplyRequest)

		// Rutas para flujo de aprobación con Pavedad
		supplyRequestGroup.GET("/pending-pavedad", supplyRequestController.GetPendingRequestsForPavedad)
		supplyRequestGroup.PUT("/:id/assign", supplyRequestController.AssignRequestToWarehouseManager)
		supplyRequestGroup.GET("/assigned/:rut", supplyRequestController.GetAssignedRequestsForWarehouseManager)

		// Rutas para revisión individual de items
		supplyRequestGroup.PUT("/items/:itemId/review", supplyRequestController.ReviewSupplyRequestItem)

		// Ruta para reenviar solicitudes devueltas
		supplyRequestGroup.PUT("/:id/resubmit", supplyRequestController.ResubmitReturnedRequest)

		// Ruta para configurar autorización de retiro (solo encargado de bodega)
		supplyRequestGroup.PUT("/:id/configure-pickup", supplyRequestController.ConfigurePickupAuthorization)

		// Rutas por pabellón
		supplyRequestGroup.GET("/pavilion/:pavilion_id", supplyRequestController.GetSupplyRequestsByPavilion)

		// Estadísticas
		supplyRequestGroup.GET("/stats", supplyRequestController.GetSupplyRequestStats)
	}

	// Grupo de rutas para asignación de QR
	qrAssignmentGroup := router.Group("/api/v1/qr-assignments")
	qrAssignmentGroup.Use(middleware.AuthMiddleware(secretKey, db)) // Aplicar autenticación
	{
		// Asignación individual y en lote
		qrAssignmentGroup.POST("", supplyRequestController.AssignQRToRequest)
		qrAssignmentGroup.POST("/bulk", supplyRequestController.BulkAssignQRs)

		// Gestión de estado de asignaciones
		qrAssignmentGroup.PUT("/:qr_code/deliver", supplyRequestController.MarkQRAsDelivered)
	}

	// Grupo de rutas para trazabilidad QR
	traceabilityGroup := router.Group("/api/v1/traceability")
	traceabilityGroup.Use(middleware.AuthMiddleware(secretKey, db)) // Aplicar autenticación
	{
		// Trazabilidad específica
		traceabilityGroup.GET("/qr/:qr_code", supplyRequestController.GetQRTraceability)
	}
}
