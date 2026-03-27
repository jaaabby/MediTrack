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
		// Rutas públicas (con rate limiting en endpoints sensibles)
		auth.POST("/login", middleware.LoginLimiter.Middleware(), authController.Login)
		auth.POST("/verify-otp", middleware.OTPLimiter.Middleware(), authController.VerifyOTP)
		auth.POST("/register", authController.Register)
		auth.POST("/forgot-password", middleware.ForgotPasswordLimiter.Middleware(), authController.RequestPasswordReset)
		auth.POST("/reset-password", authController.ResetPassword)
		auth.POST("/validate-reset-token", authController.ValidateResetToken)

		// Rutas protegidas
		auth.Use(middleware.AuthMiddleware(secretKey))
		{
			auth.GET("/profile", authController.GetProfile)
			auth.PUT("/change-password", authController.ChangePassword)
			auth.PUT("/first-time-password-change", authController.FirstTimePasswordChange)
		}
	}
}
