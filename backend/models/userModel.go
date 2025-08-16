package models

// User representa un usuario del sistema
type User struct {
	RUT             string `json:"rut" gorm:"primaryKey"`
	Name            string `json:"name" db:"name"`
	Email           string `json:"email" db:"email"`
	Password        string `json:"password" db:"password"`
	Role            string `json:"role" db:"role"`
	MedicalCenterID int    `json:"medical_center_id" db:"medical_center_id"`
}

func (u User) TableName() string {
	return "user"
}
