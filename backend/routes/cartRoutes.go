package routes

import (
	"meditrack/controllers"
	"meditrack/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

// SetupCartRoutes configura las rutas para los carritos
func SetupCartRoutes(router *gin.Engine, cartController *controllers.CartController) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "default-secret-key-change-in-production"
	}

	cartRoutes := router.Group("/api/carts")
	cartRoutes.Use(middleware.AuthMiddleware(secretKey))
	{
		// Obtener todos los carritos con paginación
		cartRoutes.GET("", cartController.GetAllCarts)

		// Obtener carrito por ID
		cartRoutes.GET("/:id", cartController.GetCartByID)

		// Obtener detalles completos de un carrito
		cartRoutes.GET("/:id/details", cartController.GetCartDetails)

		// Obtener carrito por solicitud
		cartRoutes.GET("/request/:requestId", cartController.GetCartByRequestID)

		// Obtener carrito por código QR
		cartRoutes.GET("/qr/:qrCode", cartController.GetCartByQRCode)

		// Crear carrito para una solicitud (manual)
		cartRoutes.POST("/request/:requestId", cartController.CreateCartForRequest)

		// Agregar item al carrito
		cartRoutes.POST("/:id/items", cartController.AddItemToCart)

		// Remover item del carrito
		cartRoutes.DELETE("/:id/items/:itemId", cartController.RemoveItemFromCart)

		// Marcar item como utilizado
		cartRoutes.POST("/:id/items/:itemId/use", cartController.MarkItemAsUsed)

		// Marcar item para devolución
		cartRoutes.POST("/:id/items/:itemId/return", cartController.MarkItemForReturn)

		// Operación múltiple de items (marcar varios como usados/devueltos en un solo paso)
		cartRoutes.POST("/:id/items/batch-operation", cartController.BatchOperationItems)

		// Cerrar carrito
		cartRoutes.POST("/:id/close", cartController.CloseCart)

		// Transferir carrito al pabellón
		cartRoutes.POST("/:id/transfer-to-pavilion", cartController.TransferCartToPavilion)
	}
}
