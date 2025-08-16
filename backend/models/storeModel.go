package models

// Store representa una bodega
type Store struct {
	ID              int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Type            string `json:"type" db:"type"`
	MedicalCenterID int    `json:"medical_center_id" db:"medical_center_id"`
}

func (s Store) TableName() string {
	return "store"
}
