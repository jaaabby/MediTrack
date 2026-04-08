package middleware

import (
	"net/http"

	"meditrack/config"
	"meditrack/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthMiddleware verifica la autenticación del usuario y valida el token_version contra la DB
func AuthMiddleware(secretKey string, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token de autorización requerido",
			})
			c.Abort()
			return
		}

		tokenString, err := config.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Formato de token inválido",
			})
			c.Abort()
			return
		}

		claims, err := config.ValidateToken(tokenString, secretKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token inválido o expirado",
			})
			c.Abort()
			return
		}

		// Verificar token_version contra la base de datos para detectar cierres de sesión globales
		var user models.User
		if err := db.Select("token_version").First(&user, "rut = ?", claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Usuario no encontrado",
			})
			c.Abort()
			return
		}
		if claims.TokenVersion != user.TokenVersion {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Sesión cerrada remotamente. Por favor, inicia sesión nuevamente.",
			})
			c.Abort()
			return
		}

		// Agregar claims al contexto
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// RequireRole verifica que el usuario tenga el rol requerido
func RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Usuario no autenticado",
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		hasRequiredRole := false

		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				hasRequiredRole = true
				break
			}
		}

		if !hasRequiredRole {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Acceso denegado: rol insuficiente",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin verifica que el usuario sea administrador
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin)
}

// RequirePavilion verifica que el usuario sea de pabellón
func RequirePavilion() gin.HandlerFunc {
	return RequireRole(models.RolePavilion)
}

// RequireStoreManager verifica que el usuario sea encargado de bodega
func RequireStoreManager() gin.HandlerFunc {
	return RequireRole(models.RoleStoreManager)
}

// RequireAdminOrPavilion verifica que el usuario sea admin o de pabellón
func RequireAdminOrPavilion() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin, models.RolePavilion)
}

// RequireAdminOrStoreManager verifica que el usuario sea admin o encargado de bodega
func RequireAdminOrStoreManager() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin, models.RoleStoreManager)
}

