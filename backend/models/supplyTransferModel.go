package models

import "time"

// SupplyTransfer representa una transferencia de insumos entre ubicaciones
type SupplyTransfer struct {
	ID              int        `json:"id" gorm:"primaryKey;autoIncrement"`
	TransferCode    string     `json:"transfer_code" gorm:"unique;not null"` // Código único de transferencia
	QRCode          string     `json:"qr_code" gorm:"not null;index"`        // QR del producto transferido
	MedicalSupplyID int        `json:"medical_supply_id" gorm:"not null"`
	OriginType      string     `json:"origin_type" gorm:"not null"`                 // 'store' o 'pavilion'
	OriginID        int        `json:"origin_id" gorm:"not null"`                   // ID de bodega o pabellón origen
	DestinationType string     `json:"destination_type" gorm:"not null"`            // 'store' o 'pavilion'
	DestinationID   int        `json:"destination_id" gorm:"not null"`              // ID de bodega o pabellón destino
	SentBy          string     `json:"sent_by" gorm:"not null"`                     // RUT de quien envía
	SentByName      string     `json:"sent_by_name" gorm:"not null"`                // Nombre de quien envía
	ReceivedBy      *string    `json:"received_by,omitempty"`                       // RUT de quien recibe
	ReceivedByName  *string    `json:"received_by_name,omitempty"`                  // Nombre de quien recibe
	Status          string     `json:"status" gorm:"not null;default:'pendiente'"`  // Estado de la transferencia
	TransferReason  string     `json:"transfer_reason" gorm:"type:text"`            // Motivo de la transferencia
	SendDate        time.Time  `json:"send_date" gorm:"not null"`                   // Fecha de envío
	ReceiveDate     *time.Time `json:"receive_date,omitempty"`                      // Fecha de recepción
	Notes           string     `json:"notes" gorm:"type:text"`                      // Notas adicionales
	RejectionReason *string    `json:"rejection_reason,omitempty" gorm:"type:text"` // Motivo de rechazo si aplica
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// Constantes para estados de transferencia
const (
	TransferStatusPending   = "pendiente"
	TransferStatusInTransit = "en_transito"
	TransferStatusReceived  = "recibido"
	TransferStatusRejected  = "rechazado"
	TransferStatusCancelled = "cancelado"
)

// Constantes para tipos de ubicación en transferencias
const (
	TransferLocationStore    = "store"
	TransferLocationPavilion = "pavilion"
)

func (st SupplyTransfer) TableName() string {
	return "supply_transfer"
}

// IsPending verifica si la transferencia está pendiente
func (st SupplyTransfer) IsPending() bool {
	return st.Status == TransferStatusPending || st.Status == TransferStatusInTransit
}

// IsCompleted verifica si la transferencia está completada
func (st SupplyTransfer) IsCompleted() bool {
	return st.Status == TransferStatusReceived
}

// CanBeCancelled verifica si la transferencia puede ser cancelada
func (st SupplyTransfer) CanBeCancelled() bool {
	return st.Status == TransferStatusPending || st.Status == TransferStatusInTransit
}

// SupplyTransferWithDetails extiende SupplyTransfer con información adicional
type SupplyTransferWithDetails struct {
	SupplyTransfer
	SupplyName        string  `json:"supply_name,omitempty"`
	SupplyCode        int     `json:"supply_code,omitempty"`
	BatchNumber       int     `json:"batch_number,omitempty"`
	OriginName        string  `json:"origin_name,omitempty"`
	DestinationName   string  `json:"destination_name,omitempty"`
	MedicalCenterName *string `json:"medical_center_name,omitempty"`
}
