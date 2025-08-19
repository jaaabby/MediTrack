package models

import (
	"time"
)

// Batch representa un lote de insumos médicos
type Batch struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ExpirationDate time.Time `json:"expiration_date" db:"expiration_date" gorm:"not null"`
	Amount         int       `json:"amount" db:"amount" gorm:"not null"`
	Supplier       string    `json:"supplier" db:"supplier" gorm:"not null"`
	StoreID        int       `json:"store_id" db:"store_id" gorm:"not null"`
}

func (b Batch) TableName() string {
	return "batch"
}
