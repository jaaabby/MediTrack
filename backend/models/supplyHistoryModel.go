package models

import (
	"time"
)

// SupplyHistory representa el historial de un insumo médico
type SupplyHistory struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	DateTime        time.Time `json:"date_time" db:"date_time"`
	Status          string    `json:"status" db:"status"`
	DestinationType string    `json:"destination_type" db:"destination_type"`
	DestinationID   int       `json:"destination_id" db:"destination_id"`
	MedicalSupplyID int       `json:"medical_supply_id" db:"medical_supply_id"`
	UserRUT         string    `json:"user_rut" db:"user_rut"`
}

const (
	DestinationTypePavilion = "pavilion"
	DestinationTypeStore    = "store"
)

func (s SupplyHistory) TableName() string {
	return "supply_history"
}
