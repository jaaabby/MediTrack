package routes

import (
	"meditrack/controllers"
	"meditrack/middleware"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupSupplyTransferRoutes configura las rutas de transferencias
func SetupSupplyTransferRoutes(router *gin.RouterGroup, transferService *services.SupplyTransferService, secretKey string) {
	transferController := controllers.NewSupplyTransferController(transferService)
	transfers := router.Group("/transfers")
	transfers.Use(middleware.AuthMiddleware(secretKey)) // Aplicar autenticación a todas las rutas
	{
		// Crear transferencias
		transfers.POST("/to-pavilion/", transferController.TransferToPavilion)
		transfers.POST("/return-to-store/", transferController.ReturnToStore)

		// Confirmar y gestionar transferencias
		transfers.POST("/:code/confirm/", transferController.ConfirmReception)
		transfers.POST("/:code/cancel/", transferController.CancelTransfer)

		// Consultar transferencias
		transfers.GET("/:code/", transferController.GetTransferByCode)
		transfers.GET("/", transferController.GetTransfers)
	}
}
