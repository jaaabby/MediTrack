package models

// Pavilion representa un pabellón
type Pavilion struct {
	ID              int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string `json:"name" gorm:"not null"`
	MedicalCenterID int    `json:"medical_center_id" gorm:"not null"`
}

func (p Pavilion) TableName() string {
	return "pavilion"
}
