package models

import (
	"time"
)

// SurgeryTypicalSupply representa la relación entre una cirugía y sus insumos típicos
type SurgeryTypicalSupply struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	SurgeryID      int       `json:"surgery_id" gorm:"not null"`
	SupplyCode      int       `json:"supply_code" gorm:"not null"`
	TypicalQuantity int      `json:"typical_quantity" gorm:"default:1"`
	IsRequired      bool      `json:"is_required" gorm:"default:false"`
	Notes           string    `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relaciones
	Surgery    Surgery    `json:"surgery,omitempty" gorm:"foreignKey:SurgeryID"`
	SupplyCodeInfo SupplyCode `json:"supply_code_info,omitempty" gorm:"foreignKey:SupplyCode;references:Code"`
}

func (s SurgeryTypicalSupply) TableName() string {
	return "surgery_typical_supply"
}

