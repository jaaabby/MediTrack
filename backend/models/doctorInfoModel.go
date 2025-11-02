package models

import (
	"time"
)

// DoctorInfo representa información extendida de un doctor
type DoctorInfo struct {
	UserRUT              string    `json:"user_rut" gorm:"primaryKey"`
	MedicalLicense       string    `json:"medical_license" gorm:"unique"`
	LicenseExpirationDate *time.Time `json:"license_expiration_date"`
	Specialization       string    `json:"specialization"`
	SpecialtyID          *int      `json:"specialty_id"`
	YearsOfExperience    *int      `json:"years_of_experience"`
	Phone                string    `json:"phone"`
	EmergencyContact     string    `json:"emergency_contact"`
	EmergencyPhone       string    `json:"emergency_phone"`
	IsAvailable          bool      `json:"is_available" gorm:"default:true"`
	Notes                string    `json:"notes" gorm:"type:text"`
	CreatedAt            time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relaciones
	User      User             `json:"user,omitempty" gorm:"foreignKey:UserRUT;references:RUT"`
	Specialty *MedicalSpecialty `json:"specialty,omitempty" gorm:"foreignKey:SpecialtyID"`
}

func (d DoctorInfo) TableName() string {
	return "doctor_info"
}

