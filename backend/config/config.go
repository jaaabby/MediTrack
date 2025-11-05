package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Logging  LoggingConfig
	CORS     CORSConfig
	Redis    RedisConfig
}

// DatabaseConfig configuración de la base de datos
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// ServerConfig configuración del servidor
type ServerConfig struct {
	Port int
	Env  string
}

// JWTConfig configuración de JWT
type JWTConfig struct {
	Secret     string
	Expiration string
}

// LoggingConfig configuración de logging
type LoggingConfig struct {
	Level  string
	Format string
}

// CORSConfig configuración de CORS
type CORSConfig struct {
	AllowedOrigins []string
}

// RedisConfig configuración de Redis
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	// Detectar entorno
	env := os.Getenv("ENV")

	// Solo cargar .env si NO estamos en producción
	if env != "production" {
		if err := godotenv.Load(); err != nil {
			logrus.Warn("No se pudo cargar archivo .env, usando variables de entorno del sistema")
		} else {
			logrus.Info("Archivo .env cargado (desarrollo local)")
		}
	} else {
		logrus.Info("Modo producción: usando variables de entorno de Railway")
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "meditrack"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnvAsInt("PORT", 8080),
			Env:  getEnv("ENV", "development"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "default_jwt_secret"),
			Expiration: getEnv("JWT_EXPIRATION", "24h"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}

	return config, nil
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt obtiene una variable de entorno como entero o retorna un valor por defecto
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsSlice obtiene una variable de entorno como slice o retorna un valor por defecto
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		slice := strings.Split(value, ",")
		for i, item := range slice {
			slice[i] = strings.TrimSpace(item)
		}
		return slice
	}
	return defaultValue
}
