package models

import (
	"time"
)

// SupplyHistory representa el historial de un insumo médico
type SupplyHistory struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	DateTime        time.Time `json:"date_time" db:"date_time" gorm:"not null"`
	Status          string    `json:"status" db:"status" gorm:"not null"`
	DestinationType string    `json:"destination_type" db:"destination_type" gorm:"not null"`
	DestinationID   int       `json:"destination_id" db:"destination_id" gorm:"not null"`
	MedicalSupplyID int       `json:"medical_supply_id" db:"medical_supply_id" gorm:"not null"`
	UserRUT         string    `json:"user_rut" db:"user_rut" gorm:"not null"`
}

// Constantes para los tipos de destino
const (
	DestinationTypePavilion = "pavilion"
	DestinationTypeStore    = "store"
)

func (s SupplyHistory) TableName() string {
	return "supply_history"
}

// CurrentTime retorna el tiempo actual - útil para testing y consistencia
func CurrentTime() time.Time {
	return time.Now()
}

// IsConsumed verifica si el estado indica que el insumo fue consumido
func (s SupplyHistory) IsConsumed() bool {
	return s.Status == "consumido"
}

// IsAvailable verifica si el estado indica que el insumo está disponible
func (s SupplyHistory) IsAvailable() bool {
	return s.Status == "disponible" || s.Status == "recepcionado"
}

// GetDestinationDescription retorna una descripción del destino basada en el tipo
func (s SupplyHistory) GetDestinationDescription() string {
	switch s.DestinationType {
	case DestinationTypePavilion:
		return "Pabellón"
	case DestinationTypeStore:
		return "Almacén"
	default:
		return "Destino desconocido"
	}
}

// SupplyHistoryWithDestination extiende SupplyHistory con información del destino
type SupplyHistoryWithDestination struct {
	SupplyHistory
	DestinationName   *string `json:"destination_name,omitempty"`
	MedicalCenterName *string `json:"medical_center_name,omitempty"`
	UserName          *string `json:"user_name,omitempty"`
}
