package models

// MedicalSupply representa un insumo médico
type MedicalSupply struct {
	ID   int `json:"id" gorm:"primaryKey;autoIncrement"`
	Code int `json:"code" db:"code" gorm:"not null"`
}

func (m MedicalSupply) TableName() string {
	return "medical_supply"
}
