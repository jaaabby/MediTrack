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
	Notes           string    `json:"notes" gorm:"type:text"` // Historial completo de comentarios: solicitante, pavedad, encargado, etc.
	// Campos de asignación por Pavedad
	AssignedTo            *string    `json:"assigned_to"`
	AssignedToName        *string    `json:"assigned_to_name"`
	AssignedDate          *time.Time `json:"assigned_date"`
	AssignedByPavedad     *string    `json:"assigned_by_pavedad"`
	AssignedByPavedadName *string    `json:"assigned_by_pavedad_name"`
	PavedadNotes          *string    `json:"pavedad_notes,omitempty" gorm:"type:text"` // DEPRECATED: usar Notes para historial completo
	// Campos de aprobación/rechazo
	ApprovedBy      *string    `json:"approved_by"`
	ApprovedByName  *string    `json:"approved_by_name"`
	ApprovalDate    *time.Time `json:"approval_date"`
	CompletedDate   *time.Time `json:"completed_date"`
	MedicalCenterID int        `json:"medical_center_id" gorm:"not null"`
	// Campos de médico responsable
	SurgeonID   *string `json:"surgeon_id"`
	SurgeonName *string `json:"surgeon_name"`
	SurgeryID   *int    `json:"surgery_id"`
	SpecialtyID *int    `json:"specialty_id"`
	// Campos para control de retiro
	AllowAnyoneToPickup  bool      `json:"allow_anyone_to_pickup" gorm:"default:true"`
	AuthorizedPickupRUT  *string   `json:"authorized_pickup_rut"`
	AuthorizedPickupName *string   `json:"authorized_pickup_name"`
	CreatedAt            time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relaciones
	Surgeon   *User             `json:"surgeon,omitempty" gorm:"foreignKey:SurgeonID;references:RUT"`
	Surgery   *Surgery          `json:"surgery,omitempty" gorm:"foreignKey:SurgeryID"`
	Specialty *MedicalSpecialty `json:"specialty,omitempty" gorm:"foreignKey:SpecialtyID"`
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
	QRCode              string     `json:"qr_code" gorm:"not null"` // Código QR del insumo asignado
	AssignedDate        time.Time  `json:"assigned_date" gorm:"not null"`
	AssignedBy          string     `json:"assigned_by" gorm:"not null"`
	AssignedByName      string     `json:"assigned_by_name" gorm:"not null"`
	DeliveredDate       *time.Time `json:"delivered_date"`
	DeliveredBy         *string    `json:"delivered_by"`
	DeliveredByName     *string    `json:"delivered_by_name"`
	Status              string     `json:"status" gorm:"not null;default:assigned"`
	Notes               string     `json:"notes" gorm:"type:text"`
	// Campos para notificaciones de insumos no consumidos
	LastNotificationSent *time.Time `json:"last_notification_sent"`
	NotificationCount    int        `json:"notification_count" gorm:"default:0"`
	CreatedAt            time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relaciones
	SupplyRequest     SupplyRequest     `json:"supply_request,omitempty" gorm:"foreignKey:SupplyRequestID"`
	SupplyRequestItem SupplyRequestItem `json:"supply_request_item,omitempty" gorm:"foreignKey:SupplyRequestItemID"`
	MedicalSupply     MedicalSupply     `json:"medical_supply,omitempty" gorm:"foreignKey:MedicalSupplyID"`
	Cart              *SupplyCart       `json:"cart,omitempty" gorm:"-"` // Carrito asociado (no es un campo de BD, se carga manualmente)
}

// Constantes para Status de SupplyRequest
const (
	RequestStatusPendingPavedad      = "pendiente_pavedad"     // Doctor crea solicitud
	RequestStatusAssignedStore       = "asignado_bodega"       // Pavedad asigna a encargado de bodega
	RequestStatusInProcess           = "en_proceso"            // Encargado de bodega está procesando
	RequestStatusApproved            = "aprobado"              // Encargado de bodega aprueba
	RequestStatusRejected            = "rechazado"             // Encargado de bodega rechaza
	RequestStatusCompleted           = "completado"            // Solicitud completada
	RequestStatusCancelled           = "cancelado"             // Solicitud cancelada
	RequestStatusReturnedToRequester = "devuelto"              // Encargado devuelve items al solicitante para modificar
	RequestStatusReturnedToStore     = "devuelto_al_encargado" // Doctor reenvía solicitud devuelta al encargado

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
	return time.Until(s.SurgeryDatetime).Hours()
}

// GetDaysUntilSurgery retorna los días restantes hasta la cirugía
func (s *SupplyRequest) GetDaysUntilSurgery() float64 {
	return time.Until(s.SurgeryDatetime).Hours() / 24.0
}

// IsSurgeryOverdue verifica si la cirugía ya pasó
func (s *SupplyRequest) IsSurgeryOverdue() bool {
	return time.Now().After(s.SurgeryDatetime)
}

// IsUrgent verifica si la cirugía es urgente (menos de 48 horas)
func (s *SupplyRequest) IsUrgent() bool {
	hoursUntil := s.GetHoursUntilSurgery()
	return hoursUntil > 0 && hoursUntil <= 48
}

// IsEmergency verifica si la cirugía es de emergencia (menos de 12 horas)
func (s *SupplyRequest) IsEmergency() bool {
	hoursUntil := s.GetHoursUntilSurgery()
	return hoursUntil > 0 && hoursUntil <= 12
}

// IsNotProgrammed verifica si la cirugía no está programada con suficiente anticipación (menos de 3 días)
func (s *SupplyRequest) IsNotProgrammed() bool {
	daysUntil := s.GetDaysUntilSurgery()
	return daysUntil > 0 && daysUntil < 3
}

// GetUrgencyLevel retorna el nivel de urgencia: "emergency", "urgent", "normal", "low"
func (s *SupplyRequest) GetUrgencyLevel() string {
	hoursUntil := s.GetHoursUntilSurgery()
	if hoursUntil < 0 {
		return "completed"
	}
	if hoursUntil <= 12 {
		return "emergency"
	}
	if hoursUntil <= 48 {
		return "urgent"
	}
	if hoursUntil <= 72 {
		return "normal"
	}
	return "low"
}

// GetDaysUntilSurgeryFromRequest calcula los días desde la solicitud hasta la cirugía
func (s *SupplyRequest) GetDaysUntilSurgeryFromRequest() float64 {
	return s.SurgeryDatetime.Sub(s.RequestDate).Hours() / 24.0
}

// HasMinimumAdvanceNotice verifica si tiene la anticipación mínima requerida (3 días por defecto)
func (s *SupplyRequest) HasMinimumAdvanceNotice(minDays float64) bool {
	if minDays <= 0 {
		minDays = 3.0 // Por defecto 3 días
	}
	daysUntil := s.GetDaysUntilSurgeryFromRequest()
	return daysUntil >= minDays
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
