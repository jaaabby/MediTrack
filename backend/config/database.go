package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Connect establece conexión con PostgreSQL
func Connect(cfg DatabaseConfig) (*sql.DB, error) {
	// Construir string de conexión
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	// Abrir conexión
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión: %w", err)
	}

	// Verificar conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error al verificar conexión: %w", err)
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

// Close cierra la conexión a la base de datos
func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}


