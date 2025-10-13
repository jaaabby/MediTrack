package models

import "time"

// StoreInventorySummary representa el resumen de inventario de una bodega
type StoreInventorySummary struct {
	ID                   int        `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID              int        `json:"store_id" gorm:"not null;index"`
	BatchID              int        `json:"batch_id" gorm:"not null;index;unique"` // Un lote solo puede estar en una bodega
	SupplyCode           int        `json:"supply_code" gorm:"not null"`
	SurgeryID            *int       `json:"surgery_id,omitempty"`                              // Tipo de cirugía asociado
	OriginalAmount       int        `json:"original_amount" gorm:"not null"`                   // Cantidad original del lote
	CurrentInStore       int        `json:"current_in_store" gorm:"not null"`                  // Stock actual en bodega
	TotalTransferredOut  int        `json:"total_transferred_out" gorm:"not null;default:0"`   // Total enviado a pabellones
	TotalReturnedIn      int        `json:"total_returned_in" gorm:"not null;default:0"`       // Total devuelto de pabellones
	TotalConsumedInStore int        `json:"total_consumed_in_store" gorm:"not null;default:0"` // Total consumido desde bodega
	LastTransferOutDate  *time.Time `json:"last_transfer_out_date,omitempty"`                  // Última fecha de envío
	LastReturnInDate     *time.Time `json:"last_return_in_date,omitempty"`                     // Última fecha de devolución
	LastConsumedDate     *time.Time `json:"last_consumed_date,omitempty"`                      // Última fecha de consumo
	CreatedAt            time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (sis StoreInventorySummary) TableName() string {
	return "store_inventory_summary"
}

// IsLowStock verifica si el stock está bajo (menos del 20% del original)
func (sis StoreInventorySummary) IsLowStock() bool {
	if sis.OriginalAmount == 0 {
		return false
	}
	threshold := float64(sis.OriginalAmount) * 0.2
	return float64(sis.CurrentInStore) < threshold
}

// GetTransferRate calcula el porcentaje de transferencias realizadas
func (sis StoreInventorySummary) GetTransferRate() float64 {
	if sis.OriginalAmount == 0 {
		return 0
	}
	return (float64(sis.TotalTransferredOut) / float64(sis.OriginalAmount)) * 100
}

// GetNetStock calcula el stock neto (original - transferido + devuelto - consumido)
func (sis StoreInventorySummary) GetNetStock() int {
	return sis.OriginalAmount - sis.TotalTransferredOut + sis.TotalReturnedIn - sis.TotalConsumedInStore
}

// StoreInventorySummaryWithDetails extiende StoreInventorySummary con información adicional
type StoreInventorySummaryWithDetails struct {
	StoreInventorySummary
	StoreName         string  `json:"store_name"`
	SupplyName        string  `json:"supply_name"`
	SurgeryName       *string `json:"surgery_name,omitempty"`
	BatchSupplier     string  `json:"batch_supplier"`
	ExpirationDate    string  `json:"expiration_date"`
	MedicalCenterID   int     `json:"medical_center_id"`
	MedicalCenterName *string `json:"medical_center_name,omitempty"`
}
