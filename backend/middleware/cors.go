package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configura los headers de CORS para Gin
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Permitir todos los orígenes
		c.Header("Access-Control-Allow-Origin", "*")

		// Permitir métodos HTTP
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Permitir headers
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Permitir credenciales
		c.Header("Access-Control-Allow-Credentials", "true")

		// Manejar preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
