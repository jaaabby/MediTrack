package models

import "strings"

// Constantes para los roles de usuario
const (
	RoleAdmin        = "admin"
	RolePavilion     = "pabellón"
	RoleStoreManager = "encargado de bodega"
	RoleNurse        = "enfermera"
	RoleDoctor       = "doctor"
	RolePavedad      = "pavedad"
)

// User representa un usuario del sistema
type User struct {
	RUT                    string            `json:"rut" gorm:"primaryKey"`
	Name                   string            `json:"name" db:"name" gorm:"not null"`
	Email                  string            `json:"email" db:"email"`
	Password               string            `json:"password" db:"password" gorm:"not null"`
	Role                   string            `json:"role" db:"role" gorm:"not null;check:role IN ('admin', 'pabellón', 'encargado de bodega', 'enfermera', 'doctor', 'pavedad')"`
	MedicalCenterID        int               `json:"medical_center_id" db:"medical_center_id" gorm:"not null"`
	MedicalCenter          *MedicalCenter    `json:"medical_center,omitempty" gorm:"foreignKey:MedicalCenterID"`
	PavilionID             *int              `json:"pavilion_id" db:"pavilion_id"`
	Pavilion               *Pavilion         `json:"pavilion,omitempty" gorm:"foreignKey:PavilionID"`
	SpecialtyID            *int              `json:"specialty_id" db:"specialty_id"`
	Specialty              *MedicalSpecialty `json:"specialty,omitempty" gorm:"foreignKey:SpecialtyID"`
	IsActive               bool              `json:"is_active" db:"is_active" gorm:"default:true"`
	MustChangePassword     bool              `json:"must_change_password" db:"must_change_password" gorm:"default:false"`
	ResetPasswordToken     *string           `json:"-" db:"reset_password_token"`
	ResetPasswordExpiresAt *int64            `json:"-" db:"reset_password_expires_at"`
	FailedLoginAttempts    int               `json:"-" db:"failed_login_attempts" gorm:"default:0"`
	LockedUntil            *int64            `json:"-" db:"locked_until"`
	TokenVersion           int               `json:"-" db:"token_version" gorm:"default:1;not null"`
	TOTPSecret             *string           `json:"-" db:"totp_secret" gorm:"column:totp_secret"`
	TOTPEnabled            bool              `json:"totp_enabled" db:"totp_enabled" gorm:"column:totp_enabled;default:false"`
	CreatedAt              int64             `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt              int64             `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

// UserResponse representa la respuesta del usuario sin contraseña
type UserResponse struct {
	RUT                string            `json:"rut"`
	Name               string            `json:"name"`
	Email              string            `json:"email"`
	Role               string            `json:"role"`
	MedicalCenterID    int               `json:"medical_center_id"`
	MedicalCenter      *MedicalCenter    `json:"medical_center,omitempty"`
	PavilionID         *int              `json:"pavilion_id,omitempty"`
	SpecialtyID        *int              `json:"specialty_id,omitempty"`
	Specialty          *MedicalSpecialty `json:"specialty,omitempty"`
	IsActive           bool              `json:"is_active"`
	MustChangePassword bool              `json:"must_change_password"`
	TOTPEnabled        bool              `json:"totp_enabled"`
	CreatedAt          int64             `json:"created_at"`
	UpdatedAt          int64             `json:"updated_at"`
}

// ToResponse convierte un User a UserResponse
func (u User) ToResponse() UserResponse {
	return UserResponse{
		RUT:                u.RUT,
		Name:               u.Name,
		Email:              u.Email,
		Role:               u.Role,
		MedicalCenterID:    u.MedicalCenterID,
		MedicalCenter:      u.MedicalCenter,
		PavilionID:         u.PavilionID,
		SpecialtyID:        u.SpecialtyID,
		Specialty:          u.Specialty,
		IsActive:           u.IsActive,
		MustChangePassword: u.MustChangePassword,
		TOTPEnabled:        u.TOTPEnabled,
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
	}
}

// IsValidRole verifica si el rol es válido
func (u User) IsValidRole() bool {
	return u.Role == RoleAdmin ||
		u.Role == RolePavilion ||
		u.Role == RoleStoreManager ||
		u.Role == RoleNurse ||
		u.Role == RoleDoctor ||
		u.Role == RolePavedad
}

// IsConsignationWarehouse verifica si el usuario es de bodega de consignación basándose en el email
func (u User) IsConsignationWarehouse() bool {
	if u.Role != RoleStoreManager {
		return false
	}
	emailLower := strings.ToLower(u.Email)
	return strings.Contains(emailLower, "bodegaconsignacion") || strings.Contains(emailLower, "consignacion")
}

// IsCentralWarehouse verifica si el usuario es de bodega central basándose en el email
func (u User) IsCentralWarehouse() bool {
	if u.Role != RoleStoreManager {
		return false
	}
	emailLower := strings.ToLower(u.Email)
	return !strings.Contains(emailLower, "bodegaconsignacion") && !strings.Contains(emailLower, "consignacion")
}

func (u User) TableName() string {
	return "user"
}
