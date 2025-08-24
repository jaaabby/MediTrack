package models

import (
	"time"
)

// Batch representa un lote de insumos médicos
type Batch struct {
	ID             int       `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	ExpirationDate time.Time `json:"expiration_date" gorm:"not null"`
	Amount         int       `json:"amount" gorm:"not null"`
	Supplier       string    `json:"supplier" gorm:"not null"`
	StoreID        int       `json:"store_id" gorm:"not null"`
	QRCode         string    `json:"qr_code,omitempty" gorm:"unique"` // Removido "not null"
}

func (b Batch) TableName() string {
	return "batch"
}
