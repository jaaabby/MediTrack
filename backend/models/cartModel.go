package models

import (
	"time"
)

// SupplyCart representa un carrito de insumos generado automáticamente al aprobar una solicitud
type SupplyCart struct {
	ID              int        `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplyRequestID int        `json:"supply_request_id" gorm:"not null;unique"`
	CartNumber      string     `json:"cart_number" gorm:"unique;not null"`
	Status          string     `json:"status" gorm:"not null;default:active"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy       string     `json:"created_by" gorm:"not null"`
	CreatedByName   string     `json:"created_by_name" gorm:"not null"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	ClosedAt        *time.Time `json:"closed_at"`
	ClosedBy        *string    `json:"closed_by"`
	ClosedByName    *string    `json:"closed_by_name"`
	Notes           string     `json:"notes" gorm:"type:text"`

	// Relaciones
	SupplyRequest SupplyRequest    `json:"supply_request,omitempty" gorm:"foreignKey:SupplyRequestID"`
	Items         []SupplyCartItem `json:"items,omitempty" gorm:"foreignKey:SupplyCartID"`
}

// SupplyCartItem representa un item individual dentro de un carrito
type SupplyCartItem struct {
	ID                          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplyCartID                int        `json:"supply_cart_id" gorm:"not null"`
	SupplyRequestQRAssignmentID int        `json:"supply_request_qr_assignment_id" gorm:"not null"`
	AddedAt                     time.Time  `json:"added_at" gorm:"autoCreateTime"`
	AddedBy                     string     `json:"added_by" gorm:"not null"`
	AddedByName                 string     `json:"added_by_name" gorm:"not null"`
	RemovedAt                   *time.Time `json:"removed_at"`
	RemovedBy                   *string    `json:"removed_by"`
	RemovedByName               *string    `json:"removed_by_name"`
	IsActive                    bool       `json:"is_active" gorm:"not null;default:true"`
	Notes                       string     `json:"notes" gorm:"type:text"`

	// Relaciones
	SupplyCart                SupplyCart                `json:"supply_cart,omitempty" gorm:"foreignKey:SupplyCartID"`
	SupplyRequestQRAssignment SupplyRequestQRAssignment `json:"supply_request_qr_assignment,omitempty" gorm:"foreignKey:SupplyRequestQRAssignmentID"`
}

// Constantes para Status de SupplyCart
const (
	CartStatusActive    = "active"
	CartStatusClosed    = "closed"
	CartStatusCancelled = "cancelled"
)

// TableName especifica el nombre de la tabla para SupplyCart
func (s SupplyCart) TableName() string {
	return "supply_cart"
}

// TableName especifica el nombre de la tabla para SupplyCartItem
func (s SupplyCartItem) TableName() string {
	return "supply_cart_item"
}

// GetStatusLabel retorna una etiqueta legible para el estado del carrito
func (s *SupplyCart) GetStatusLabel() string {
	switch s.Status {
	case CartStatusActive:
		return "Activo"
	case CartStatusClosed:
		return "Cerrado"
	case CartStatusCancelled:
		return "Cancelado"
	default:
		return "Desconocido"
	}
}

// IsActive verifica si el carrito está activo
func (s *SupplyCart) IsActive() bool {
	return s.Status == CartStatusActive
}

// CanAddItems verifica si se pueden agregar items al carrito
func (s *SupplyCart) CanAddItems() bool {
	return s.Status == CartStatusActive
}

// GetActiveItemsCount retorna el número de items activos en el carrito
func (s *SupplyCart) GetActiveItemsCount() int {
	count := 0
	for _, item := range s.Items {
		if item.IsActive {
			count++
		}
	}
	return count
}

// GenerateCartNumber genera un número de carrito único
func GenerateCartNumber() string {
	now := time.Now()
	return "CART-" + now.Format("20060102150405")
}

// SupplyCartDetailView representa la vista detallada de un carrito con información agregada
type SupplyCartDetailView struct {
	CartID           int        `json:"cart_id"`
	CartNumber       string     `json:"cart_number"`
	SupplyRequestID  int        `json:"supply_request_id"`
	RequestNumber    string     `json:"request_number"`
	CartStatus       string     `json:"cart_status"`
	CartCreatedAt    time.Time  `json:"cart_created_at"`
	CreatedBy        string     `json:"created_by"`
	CreatedByName    string     `json:"created_by_name"`
	CartUpdatedAt    time.Time  `json:"cart_updated_at"`
	ClosedAt         *time.Time `json:"closed_at"`
	ClosedBy         *string    `json:"closed_by"`
	ClosedByName     *string    `json:"closed_by_name"`
	CartNotes        string     `json:"cart_notes"`
	ActiveItemsCount int        `json:"active_items_count"`
	TotalItemsCount  int        `json:"total_items_count"`
	RequestStatus    string     `json:"request_status"`
	RequestedByName  string     `json:"requested_by_name"`
	SurgeryDatetime  time.Time  `json:"surgery_datetime"`
	PavilionID       int        `json:"pavilion_id"`
}

// TableName especifica el nombre de la vista
func (s SupplyCartDetailView) TableName() string {
	return "v_supply_cart_details"
}
