package models

import "time"

// PavilionInventorySummary representa el resumen de inventario de un pabellón
type PavilionInventorySummary struct {
	ID               int        `json:"id" gorm:"primaryKey;autoIncrement"`
	PavilionID       int        `json:"pavilion_id" gorm:"not null;index"`
	BatchID          int        `json:"batch_id" gorm:"not null;index"`
	SupplyCode       int        `json:"supply_code" gorm:"not null"`
	TotalReceived    int        `json:"total_received" gorm:"not null;default:0"`    // Total recibido históricamente
	CurrentAvailable int        `json:"current_available" gorm:"not null;default:0"` // Stock disponible actual
	TotalConsumed    int        `json:"total_consumed" gorm:"not null;default:0"`    // Total consumido
	TotalReturned    int        `json:"total_returned" gorm:"not null;default:0"`    // Total devuelto a bodega
	LastReceivedDate *time.Time `json:"last_received_date,omitempty"`                // Última fecha de recepción
	LastConsumedDate *time.Time `json:"last_consumed_date,omitempty"`                // Última fecha de consumo
	LastReturnedDate *time.Time `json:"last_returned_date,omitempty"`                // Última fecha de devolución
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (pis PavilionInventorySummary) TableName() string {
	return "pavilion_inventory_summary"
}

// IsLowStock verifica si el stock está bajo (menos del 20% del total recibido)
func (pis PavilionInventorySummary) IsLowStock() bool {
	if pis.TotalReceived == 0 {
		return false
	}
	threshold := float64(pis.TotalReceived) * 0.2
	return float64(pis.CurrentAvailable) < threshold
}

// GetConsumptionRate calcula el porcentaje de consumo
func (pis PavilionInventorySummary) GetConsumptionRate() float64 {
	if pis.TotalReceived == 0 {
		return 0
	}
	return (float64(pis.TotalConsumed) / float64(pis.TotalReceived)) * 100
}

// PavilionInventorySummaryWithDetails extiende PavilionInventorySummary con información adicional
type PavilionInventorySummaryWithDetails struct {
	PavilionInventorySummary
	PavilionName      string  `json:"pavilion_name"`
	SupplyName        string  `json:"supply_name"`
	BatchSupplier     string  `json:"batch_supplier"`
	ExpirationDate    string  `json:"expiration_date"`
	MedicalCenterID   int     `json:"medical_center_id"`
	MedicalCenterName *string `json:"medical_center_name,omitempty"`
	InTransit         bool    `json:"in_transit" gorm:"-"`
}
