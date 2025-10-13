package models

import (
	"time"
)

// Batch representa un lote de insumos médicos
type Batch struct {
	ID             int       `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	ExpirationDate time.Time `json:"expiration_date" gorm:"not null"`
	Amount         int       `json:"amount" gorm:"not null"`
	Supplier       string    `json:"supplier" gorm:"not null"`
	StoreID        int       `json:"store_id" gorm:"not null"`
	QRCode         string    `json:"qr_code,omitempty" gorm:"unique"` // Removido "not null"

	// Nuevos campos para gestión de inventario
	SurgeryID    *int   `json:"surgery_id,omitempty" gorm:"default:null"`      // Tipo de cirugía asociado
	LocationType string `json:"location_type" gorm:"not null;default:'store'"` // 'store' o 'pavilion'
	LocationID   int    `json:"location_id" gorm:"not null"`                   // ID de la ubicación (store_id o pavilion_id)
}

// Constantes para tipos de ubicación de lotes
const (
	BatchLocationStore    = "store"
	BatchLocationPavilion = "pavilion"
)

func (b Batch) TableName() string {
	return "batch"
}
