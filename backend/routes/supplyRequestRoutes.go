package routes

import (
	"meditrack/controllers"

	"github.com/gin-gonic/gin"
)

// SetupSupplyRequestRoutes configura las rutas para solicitudes de insumo con trazabilidad QR
func SetupSupplyRequestRoutes(router *gin.Engine, supplyRequestController *controllers.SupplyRequestController) {
	// Grupo de rutas para solicitudes de insumo
	supplyRequestGroup := router.Group("/api/v1/supply-requests")
	{
		// CRUD básico de solicitudes
		supplyRequestGroup.POST("", supplyRequestController.CreateSupplyRequest)
		supplyRequestGroup.GET("", supplyRequestController.GetAllSupplyRequests)
		supplyRequestGroup.GET("/:id", supplyRequestController.GetSupplyRequestByID)
		supplyRequestGroup.PUT("/:id/approve", supplyRequestController.ApproveSupplyRequest)
		supplyRequestGroup.PUT("/:id/reject", supplyRequestController.RejectSupplyRequest)
		supplyRequestGroup.PUT("/:id/complete", supplyRequestController.CompleteSupplyRequest)
		supplyRequestGroup.DELETE("/:id", supplyRequestController.DeleteSupplyRequest)

		// Rutas por pabellón
		supplyRequestGroup.GET("/pavilion/:pavilion_id", supplyRequestController.GetSupplyRequestsByPavilion)

		// Estadísticas
		supplyRequestGroup.GET("/stats", supplyRequestController.GetSupplyRequestStats)
	}

	// Grupo de rutas para asignación de QR
	qrAssignmentGroup := router.Group("/api/v1/qr-assignments")
	{
		// Asignación individual y en lote
		qrAssignmentGroup.POST("", supplyRequestController.AssignQRToRequest)
		qrAssignmentGroup.POST("/bulk", supplyRequestController.BulkAssignQRs)

		// Gestión de estado de asignaciones
		qrAssignmentGroup.PUT("/:qr_code/deliver", supplyRequestController.MarkQRAsDelivered)
	}

	// Grupo de rutas para trazabilidad QR
	traceabilityGroup := router.Group("/api/v1/traceability")
	{
		// Trazabilidad específica
		traceabilityGroup.GET("/qr/:qr_code", supplyRequestController.GetQRTraceability)
	}
}
