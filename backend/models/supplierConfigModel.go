package models

import "time"

// SupplierConfig representa la configuración de alertas de vencimiento por proveedor
type SupplierConfig struct {
	SupplierName       string    `json:"supplier_name" gorm:"primaryKey"`
	ExpirationAlertDays int       `json:"expiration_alert_days" gorm:"not null;default:90;check:expiration_alert_days > 0"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Notes              string    `json:"notes,omitempty" gorm:"type:text"`
}

func (s SupplierConfig) TableName() string {
	return "supplier_config"
}

