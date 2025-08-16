package models

// MedicalSupply representa un insumo médico
type MedicalSupply struct {
	ID      int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Code    int    `json:"code" db:"code"`
	Name    string `json:"name" db:"name"`
	BatchID int    `json:"batch_id" db:"batch_id"`
}

func (m MedicalSupply) TableName() string {
	return "medical_supply"
}
