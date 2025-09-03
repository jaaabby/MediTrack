package routes

import (
	"meditrack/controllers"
	"meditrack/middleware"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configura las rutas de usuarios
func SetupUserRoutes(router *gin.RouterGroup, userService services.UserService, secretKey string) {
	userController := controllers.NewUserController(userService)

	users := router.Group("/users")
	{
		// Aplicar middleware de autenticación a todas las rutas de usuarios
		users.Use(middleware.AuthMiddleware(secretKey))

		// Rutas para administradores
		adminRoutes := users.Group("/")
		adminRoutes.Use(middleware.RequireAdmin())
		{
			adminRoutes.POST("/", userController.CreateUser)
			adminRoutes.GET("/", userController.GetAllUsers)
			adminRoutes.GET("/by-role", userController.GetUsersByRole)
			adminRoutes.DELETE("/:id", userController.DeleteUser)
			adminRoutes.PUT("/:id/deactivate", userController.DeactivateUser)
			adminRoutes.PUT("/:id/activate", userController.ActivateUser)
		}

		// Rutas para administradores y encargados de bodega
		adminStoreRoutes := users.Group("/")
		adminStoreRoutes.Use(middleware.RequireAdminOrStoreManager())
		{
			adminStoreRoutes.PUT("/:id", userController.UpdateUser)
		}

		// Rutas para todos los usuarios autenticados
		users.GET("/:id", userController.GetUserByID)
	}
}
