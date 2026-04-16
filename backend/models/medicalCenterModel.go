package models

// MedicalCenter representa un centro médico
type MedicalCenter struct {
	ID         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string `json:"name" db:"name" gorm:"not null"`
	Address    string `json:"address" db:"address"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
	AlertEmail string `json:"alert_email" db:"alert_email" gorm:"default:''"`
}

func (m MedicalCenter) TableName() string {
	return "medical_center"
}
