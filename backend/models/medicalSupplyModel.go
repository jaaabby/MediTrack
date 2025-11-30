package models

import "time"

// MedicalSupply representa un insumo médico individual
type MedicalSupply struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Code      int       `json:"code" db:"code" gorm:"not null"`
	QRCode    string    `json:"qr_code" db:"qr_code" gorm:"unique;not null"` // Código QR único para cada insumo individual
	BatchID   int       `json:"batch_id" db:"batch_id" gorm:"not null"`
	Status    string    `json:"status" db:"status" gorm:"not null;default:'disponible'"` // Estado actual del insumo
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`

	// Nuevos campos para gestión de inventario por ubicación
	LocationType  string     `json:"location_type" gorm:"not null;default:'store'"` // 'store' o 'pavilion'
	LocationID    int        `json:"location_id" gorm:"not null"`                   // ID de la ubicación actual
	InTransit     bool       `json:"in_transit" gorm:"default:false"`               // Si está en tránsito
	TransferDate  *time.Time `json:"transfer_date,omitempty"`                       // Fecha de última transferencia
	TransferredBy *string    `json:"transferred_by,omitempty"`                      // RUT de quien transfirió
}

// Constantes para los estados de insumos médicos
const (
	StatusAvailable         = "disponible"
	StatusPendingPickup     = "pendiente_retiro"      // Pendiente de retiro físico de bodega
	StatusEnRouteToPavilion = "en_camino_a_pabellon"
	StatusReceived          = "recepcionado"
	StatusConsumed          = "consumido"
	StatusEnRouteToStore    = "en_camino_a_bodega"
)

// Constantes para tipos de ubicación de insumos
const (
	SupplyLocationStore    = "store"
	SupplyLocationPavilion = "pavilion"
)

func (m MedicalSupply) TableName() string {
	return "medical_supply"
}

// IsConsumed verifica si el insumo está consumido
func (m MedicalSupply) IsConsumed() bool {
	return m.Status == StatusConsumed
}

// IsAvailable verifica si el insumo está disponible para uso
func (m MedicalSupply) IsAvailable() bool {
	return m.Status == StatusAvailable
}

// CanBeConsumed verifica si el insumo puede ser consumido
func (m MedicalSupply) CanBeConsumed() bool {
	return m.Status == StatusAvailable || m.Status == StatusReceived
}

// GetStatusDescription retorna una descripción del estado
func (m MedicalSupply) GetStatusDescription() string {
	switch m.Status {
	case StatusAvailable:
		return "Disponible"
	case StatusPendingPickup:
		return "Pendiente de retiro"
	case StatusEnRouteToPavilion:
		return "En camino a pabellón"
	case StatusReceived:
		return "Recepcionado"
	case StatusConsumed:
		return "Consumido"
	case StatusEnRouteToStore:
		return "En camino a bodega"
	default:
		return "Estado desconocido"
	}
}
