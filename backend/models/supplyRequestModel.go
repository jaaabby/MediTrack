package models

import (
	"time"
)

// SupplyRequest representa una solicitud de insumo con trazabilidad QR
type SupplyRequest struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	RequestNumber   string    `json:"request_number" gorm:"unique;not null"`
	PavilionID      int       `json:"pavilion_id" gorm:"not null"`
	RequestedBy     string    `json:"requested_by" gorm:"not null"` // RUT del solicitante
	RequestedByName string    `json:"requested_by_name" gorm:"not null"`
	RequestDate     time.Time `json:"request_date" gorm:"not null"`
	SurgeryDatetime time.Time `json:"surgery_datetime" gorm:"not null"`
	Status          string    `json:"status" gorm:"not null;default:pendiente_pavedad"`
	Notes           string    `json:"notes" gorm:"type:text"`
	// Campos de asignación por Pavedad
	AssignedTo            *string    `json:"assigned_to"`
	AssignedToName        *string    `json:"assigned_to_name"`
	AssignedDate          *time.Time `json:"assigned_date"`
	AssignedByPavedad     *string    `json:"assigned_by_pavedad"`
	AssignedByPavedadName *string    `json:"assigned_by_pavedad_name"`
	PavedadNotes          *string    `json:"pavedad_notes" gorm:"type:text"`
	// Campos de aprobación/rechazo
	ApprovedBy      *string    `json:"approved_by"`
	ApprovedByName  *string    `json:"approved_by_name"`
	ApprovalDate    *time.Time `json:"approval_date"`
	CompletedDate   *time.Time `json:"completed_date"`
	MedicalCenterID int        `json:"medical_center_id" gorm:"not null"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// SupplyRequestItem representa un item individual dentro de una solicitud
type SupplyRequestItem struct {
	ID                int    `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplyRequestID   int    `json:"supply_request_id" gorm:"not null"`
	SupplyCode        int    `json:"supply_code" gorm:"not null"`
	SupplyName        string `json:"supply_name" gorm:"not null"`
	QuantityRequested int    `json:"quantity_requested" gorm:"not null"`
	QuantityApproved  *int   `json:"quantity_approved"`
	QuantityDelivered int    `json:"quantity_delivered" gorm:"default:0"`
	IsPediatric       bool   `json:"is_pediatric" gorm:"default:false"`
	// Campos para gestión individual por bodega
	ItemStatus     string     `json:"item_status" gorm:"default:pendiente"` // pendiente, aceptado, rechazado, devuelto
	ItemNotes      *string    `json:"item_notes" gorm:"type:text"`
	ReviewedBy     *string    `json:"reviewed_by"` // RUT del que revisó
	ReviewedByName *string    `json:"reviewed_by_name"`
	ReviewedAt     *time.Time `json:"reviewed_at"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relaciones
	SupplyRequest  SupplyRequest `json:"supply_request,omitempty" gorm:"foreignKey:SupplyRequestID"`
	SupplyCodeInfo SupplyCode    `json:"supply_code_info,omitempty" gorm:"foreignKey:SupplyCode;references:Code"`
}

// SupplyRequestQRAssignment vincula códigos QR específicos con items de solicitud
type SupplyRequestQRAssignment struct {
	ID                  int        `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplyRequestID     int        `json:"supply_request_id" gorm:"not null"`
	SupplyRequestItemID int        `json:"supply_request_item_id" gorm:"not null"`
	MedicalSupplyID     int        `json:"medical_supply_id" gorm:"not null"`
	AssignedDate        time.Time  `json:"assigned_date" gorm:"not null"`
	AssignedBy          string     `json:"assigned_by" gorm:"not null"`
	AssignedByName      string     `json:"assigned_by_name" gorm:"not null"`
	DeliveredDate       *time.Time `json:"delivered_date"`
	DeliveredBy         *string    `json:"delivered_by"`
	DeliveredByName     *string    `json:"delivered_by_name"`
	Status              string     `json:"status" gorm:"not null;default:assigned"`
	Notes               string     `json:"notes" gorm:"type:text"`
	CreatedAt           time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relaciones
	SupplyRequest     SupplyRequest     `json:"supply_request,omitempty" gorm:"foreignKey:SupplyRequestID"`
	SupplyRequestItem SupplyRequestItem `json:"supply_request_item,omitempty" gorm:"foreignKey:SupplyRequestItemID"`
	MedicalSupply     MedicalSupply     `json:"medical_supply,omitempty" gorm:"foreignKey:MedicalSupplyID"`
}

// Constantes para Status de SupplyRequest
const (
	RequestStatusPendingPavedad = "pendiente_pavedad" // Doctor crea solicitud
	RequestStatusAssignedStore  = "asignado_bodega"   // Pavedad asigna a encargado de bodega
	RequestStatusInProcess      = "en_proceso"        // Encargado de bodega está procesando
	RequestStatusApproved       = "aprobado"          // Encargado de bodega aprueba
	RequestStatusRejected       = "rechazado"         // Encargado de bodega rechaza
	RequestStatusCompleted      = "completado"        // Solicitud completada
	RequestStatusCancelled      = "cancelado"         // Solicitud cancelada

	// Alias para compatibilidad con código existente
	RequestStatusPending = RequestStatusPendingPavedad
)

// Constantes para Status de SupplyRequestQRAssignment
const (
	AssignmentStatusAssigned  = "assigned"
	AssignmentStatusDelivered = "delivered"
	AssignmentStatusConsumed  = "consumed"
	AssignmentStatusReturned  = "returned"
	AssignmentStatusLost      = "lost"
)

// TableName especifica el nombre de la tabla para SupplyRequest
func (s SupplyRequest) TableName() string {
	return "supply_request"
}

// TableName especifica el nombre de la tabla para SupplyRequestItem
func (s SupplyRequestItem) TableName() string {
	return "supply_request_item"
}

// TableName especifica el nombre de la tabla para SupplyRequestQRAssignment
func (s SupplyRequestQRAssignment) TableName() string {
	return "supply_request_qr_assignment"
}

// GetStatusLabel retorna una etiqueta legible para el estado
func (s *SupplyRequest) GetStatusLabel() string {
	switch s.Status {
	case RequestStatusPending:
		return "Pendiente"
	case RequestStatusApproved:
		return "Aprobada"
	case RequestStatusRejected:
		return "Rechazada"
	case RequestStatusInProcess:
		return "En Proceso"
	case RequestStatusCompleted:
		return "Completada"
	case RequestStatusCancelled:
		return "Cancelada"
	default:
		return "Desconocido"
	}
}

// IsEditable verifica si la solicitud puede ser editada
func (s *SupplyRequest) IsEditable() bool {
	return s.Status == RequestStatusPending
}

// CanBeApproved verifica si la solicitud puede ser aprobada
func (s *SupplyRequest) CanBeApproved() bool {
	return s.Status == RequestStatusPending
}

// CanBeProcessed verifica si la solicitud puede ser procesada (asignar QRs)
func (s *SupplyRequest) CanBeProcessed() bool {
	return s.Status == RequestStatusApproved
}

// GetHoursUntilSurgery retorna las horas restantes hasta la cirugía
func (s *SupplyRequest) GetHoursUntilSurgery() float64 {
	return s.SurgeryDatetime.Sub(time.Now()).Hours()
}

// IsSurgeryOverdue verifica si la cirugía ya pasó
func (s *SupplyRequest) IsSurgeryOverdue() bool {
	return time.Now().After(s.SurgeryDatetime)
}

// GenerateRequestNumber genera un número de solicitud único
func GenerateRequestNumber() string {
	now := time.Now()
	return "SOL-" + now.Format("20060102150405")
}

// GetTotalItemsRequested calcula el total de items solicitados
func (s *SupplyRequest) GetTotalItemsRequested(items []SupplyRequestItem) int {
	total := 0
	for _, item := range items {
		total += item.QuantityRequested
	}
	return total
}

// GetTotalItemsApproved calcula el total de items aprobados
func (s *SupplyRequest) GetTotalItemsApproved(items []SupplyRequestItem) int {
	total := 0
	for _, item := range items {
		if item.QuantityApproved != nil {
			total += *item.QuantityApproved
		}
	}
	return total
}

// GetTotalItemsDelivered calcula el total de items entregados
func (s *SupplyRequest) GetTotalItemsDelivered(items []SupplyRequestItem) int {
	total := 0
	for _, item := range items {
		total += item.QuantityDelivered
	}
	return total
}
