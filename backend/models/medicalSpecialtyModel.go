package models

import (
	"time"
)

// MedicalSpecialty representa una especialidad médica
type MedicalSpecialty struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"not null;unique"`
	Description string    `json:"description" gorm:"type:text"`
	Code        *string   `json:"code" gorm:"unique"`
	IsActive    *bool     `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relaciones
	Surgeries []Surgery `json:"surgeries,omitempty" gorm:"foreignKey:SpecialtyID"`
	Doctors   []User    `json:"doctors,omitempty" gorm:"foreignKey:SpecialtyID;where:role = 'doctor'"`
}

func (m MedicalSpecialty) TableName() string {
	return "medical_specialty"
}
