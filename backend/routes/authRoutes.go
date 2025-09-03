package routes

import (
	"meditrack/controllers"
	"meditrack/middleware"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes configura las rutas de autenticación
func SetupAuthRoutes(router *gin.RouterGroup, userService services.UserService, secretKey string) {
	authController := controllers.NewAuthController(userService, secretKey)

	auth := router.Group("/auth")
	{
		// Rutas públicas
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)

		// Rutas protegidas
		auth.Use(middleware.AuthMiddleware(secretKey))
		{
			auth.GET("/profile", authController.GetProfile)
			auth.PUT("/change-password", authController.ChangePassword)
		}
	}
}
