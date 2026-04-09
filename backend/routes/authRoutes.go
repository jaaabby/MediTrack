package routes

import (
	"meditrack/controllers"
	"meditrack/middleware"
	"meditrack/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupAuthRoutes configura las rutas de autenticación
func SetupAuthRoutes(router *gin.RouterGroup, userService services.UserService, secretKey string, db *gorm.DB, webauthnService *services.WebAuthnService) {
	authController := controllers.NewAuthController(userService, secretKey)
	webauthnController := controllers.NewWebAuthnController(webauthnService, userService, secretKey)

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

		// Passkey login (público – discoverable, sin necesidad de email)
		auth.POST("/passkey/login/begin", webauthnController.BeginLogin)
		auth.POST("/passkey/login/finish", webauthnController.FinishLogin)

		// Rutas protegidas
		auth.Use(middleware.AuthMiddleware(secretKey, db))
		{
			auth.GET("/profile", authController.GetProfile)
			auth.PUT("/change-password", authController.ChangePassword)
			auth.PUT("/first-time-password-change", authController.FirstTimePasswordChange)
			auth.POST("/logout-all-devices", authController.LogoutAllDevices)

			// Rutas TOTP protegidas (requieren JWT completo)
			auth.GET("/totp/setup", authController.SetupTOTP)
			auth.POST("/totp/activate", authController.ActivateTOTP)
			auth.DELETE("/totp", authController.DisableTOTP)

			// Passkey: registro y gestión (requieren JWT completo)
			auth.POST("/passkey/register/begin", webauthnController.BeginRegistration)
			auth.POST("/passkey/register/finish", webauthnController.FinishRegistration)
			auth.GET("/passkey/credentials", webauthnController.ListPasskeys)
			auth.DELETE("/passkey/credentials/:id", webauthnController.DeletePasskey)
		}
	}
}
