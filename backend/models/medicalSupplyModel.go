package models

// MedicalSupply representa un insumo médico individual
type MedicalSupply struct {
	ID     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Code   int    `json:"code" db:"code" gorm:"not null"`
	QRCode string `json:"qr_code" db:"qr_code" gorm:"unique;not null"` // Código QR único para cada insumo individual
}

func (m MedicalSupply) TableName() string {
	return "medical_supply"
}
