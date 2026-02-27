package models

import (
	"time"

	"gorm.io/gorm"
)

// Batch representa un lote de insumos médicos
type Batch struct {
	ID                  int       `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	ExpirationDate      time.Time `json:"expiration_date" gorm:"not null"`
	Amount              int       `json:"amount" gorm:"not null"`
	SupplierID          int       `json:"supplier_id" gorm:"column:supplier_id;not null"` // FK → supplier_config.id
	Supplier            string    `json:"supplier,omitempty" gorm:"-"`                    // campo transiente (no va a BD)
	StoreID             int       `json:"store_id" gorm:"not null"`
	QRCode              string    `json:"qr_code,omitempty" gorm:"unique"`
	ExpirationAlertDays int       `json:"expiration_alert_days" gorm:"not null;default:90"`

	// Campos para gestión de inventario
	SurgeryID    *int   `json:"surgery_id,omitempty" gorm:"default:null"`
	LocationType string `json:"location_type" gorm:"not null;default:'store'"`
	LocationID   int    `json:"location_id" gorm:"not null"`

	// Relación (solo para preloads; no se carga automáticamente)
	SupplierConfig *SupplierConfig `json:"supplier_config,omitempty" gorm:"foreignKey:SupplierID;references:ID"`
}

// AfterFind puebla el campo transiente Supplier desde SupplierConfig si está cargado
func (b *Batch) AfterFind(tx *gorm.DB) error {
	if b.SupplierConfig != nil && b.Supplier == "" {
		b.Supplier = b.SupplierConfig.SupplierName
	}
	return nil
}

// Constantes para tipos de ubicación de lotes
const (
	BatchLocationStore    = "store"
	BatchLocationPavilion = "pavilion"
)

func (b Batch) TableName() string {
	return "batch"
}
