package models

// Pavilion representa un pabellón
type Pavilion struct {
	ID              int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string `json:"name"`
	MedicalCenterID int    `json:"medical_center_id"`
}

func (p Pavilion) TableName() string {
	return "pavilion"
}
