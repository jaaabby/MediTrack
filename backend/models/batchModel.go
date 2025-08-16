package models

import (
	"time"
)

// Batch representa un lote de insumos médicos
type Batch struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ExpirationDate time.Time `json:"expiration_date" db:"expiration_date"`
	Amount         int       `json:"amount" db:"amount"`
	Supplier       string    `json:"supplier" db:"supplier"`
	StoreID        int       `json:"store_id" db:"store_id"`
}

func (b Batch) TableName() string {
	return "batch"
}
