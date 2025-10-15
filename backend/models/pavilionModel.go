package models

// Pavilion representa un pabellón
type Pavilion struct {
	ID              int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string         `json:"name" gorm:"not null"`
	MedicalCenterID int            `json:"medical_center_id" gorm:"not null"`
	MedicalCenter   *MedicalCenter `json:"medical_center,omitempty" gorm:"foreignKey:MedicalCenterID"`
}

func (p Pavilion) TableName() string {
	return "pavilion"
}
