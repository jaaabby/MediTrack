package routes

import (
	"meditrack/controllers"
	"meditrack/middleware"
	"meditrack/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupAuthRoutes configura las rutas de autenticación
func SetupAuthRoutes(router *gin.RouterGroup, userService services.UserService, secretKey string, db *gorm.DB) {
	authController := controllers.NewAuthController(userService, secretKey)

	auth := router.Group("/auth")
	{
		// Rutas públicas
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
		auth.POST("/forgot-password", authController.RequestPasswordReset)
		auth.POST("/reset-password", authController.ResetPassword)
		auth.POST("/validate-reset-token", authController.ValidateResetToken)

		// Rutas protegidas
		auth.Use(middleware.AuthMiddleware(secretKey, db))
		{
			auth.GET("/profile", authController.GetProfile)
			auth.PUT("/change-password", authController.ChangePassword)
			auth.PUT("/first-time-password-change", authController.FirstTimePasswordChange)
			auth.POST("/logout-all-devices", authController.LogoutAllDevices)
		}
	}
}
