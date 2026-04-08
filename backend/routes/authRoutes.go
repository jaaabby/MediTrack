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
		auth.POST("/forgot-password", authController.RequestPasswordReset)
		auth.POST("/reset-password", authController.ResetPassword)
		auth.POST("/validate-reset-token", authController.ValidateResetToken)

		// Rutas TOTP públicas (usan pre-auth token, no el JWT completo)
		auth.POST("/totp/verify", authController.VerifyTOTP)

		// Rutas protegidas
		auth.Use(middleware.AuthMiddleware(secretKey))
		{
			auth.GET("/profile", authController.GetProfile)
			auth.PUT("/change-password", authController.ChangePassword)
			auth.PUT("/first-time-password-change", authController.FirstTimePasswordChange)

			// Rutas TOTP protegidas (requieren JWT completo)
			auth.GET("/totp/setup", authController.SetupTOTP)
			auth.POST("/totp/activate", authController.ActivateTOTP)
			auth.DELETE("/totp", authController.DisableTOTP)
		}
	}
}
