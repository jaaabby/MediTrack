package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configura los headers de CORS para Gin usando gin-contrib/cors
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		// Permite acceso desde cualquier origen en desarrollo
		AllowAllOrigins: true,
		// Para producción, especifica los orígenes permitidos:
		// AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173", "http://192.168.100.65:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization", "X-Requested-With",
			"X-Session-ID", "X-Device-Info", "X-Browser-Info",
		},
		AllowCredentials: false, // Debe ser false cuando AllowAllOrigins es true
	})
}
