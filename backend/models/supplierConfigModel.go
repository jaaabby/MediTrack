package models

import "time"

// SupplierConfig representa los datos de un proveedor
type SupplierConfig struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplierName string    `json:"supplier_name" gorm:"uniqueIndex;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Notes        string    `json:"notes,omitempty" gorm:"type:text"`
}

func (s SupplierConfig) TableName() string {
	return "supplier_config"
}
