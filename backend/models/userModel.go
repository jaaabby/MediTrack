package models

// User representa un usuario del sistema
type User struct {
	RUT             string `json:"rut" gorm:"primaryKey"`
	Name            string `json:"name" db:"name" gorm:"not null"`
	Email           string `json:"email" db:"email"`
	Password        string `json:"password" db:"password" gorm:"not null"`
	Role            string `json:"role" db:"role" gorm:"not null"`
	MedicalCenterID int    `json:"medical_center_id" db:"medical_center_id" gorm:"not null"`
}

func (u User) TableName() string {
	return "user"
}
