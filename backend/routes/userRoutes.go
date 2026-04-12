package routes

import (
	"meditrack/controllers"
	"meditrack/middleware"
	"meditrack/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupUserRoutes configura las rutas de usuarios
func SetupUserRoutes(router *gin.RouterGroup, userService services.UserService, secretKey string, db *gorm.DB) {
	userController := controllers.NewUserController(userService)

	users := router.Group("/users")
	{
		// Aplicar middleware de autenticación a todas las rutas de usuarios
		users.Use(middleware.AuthMiddleware(secretKey, db))

		// Rutas para administradores
		adminRoutes := users.Group("/")
		adminRoutes.Use(middleware.RequireAdmin())
		{
			adminRoutes.POST("/", userController.CreateUser)
			adminRoutes.GET("/", userController.GetAllUsers)
			adminRoutes.DELETE("/:id", userController.DeleteUser)
			adminRoutes.PUT("/:id/deactivate", userController.DeactivateUser)
			adminRoutes.PUT("/:id/activate", userController.ActivateUser)
		}

		// Ruta accesible para Admin y Pavedad (para obtener usuarios por rol)
		users.GET("/by-role", userController.GetUsersByRole)

		// Rutas para administradores y encargados de bodega
		adminStoreRoutes := users.Group("/")
		adminStoreRoutes.Use(middleware.RequireAdminOrStoreManager())
		{
			adminStoreRoutes.PUT("/:id", userController.UpdateUser)
			adminStoreRoutes.GET("/search", userController.SearchUsers)
		}

		// Rutas para todos los usuarios autenticados
		users.GET("/:id", userController.GetUserByID)
		users.GET("/profile", userController.GetUserProfileByEmail)
	}
}
