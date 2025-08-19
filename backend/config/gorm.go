package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectGORM establece conexión con PostgreSQL usando GORM
func ConnectGORM(cfg DatabaseConfig) (*gorm.DB, error) {
	// Construir string de conexión para GORM
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode,
	)

	// Abrir conexión con GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error al conectar con GORM: %w", err)
	}

	// Obtener la conexión SQL subyacente para configurar el pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error al obtener conexión SQL: %w", err)
	}

	// Configurar pool de conexiones
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	return db, nil
}
