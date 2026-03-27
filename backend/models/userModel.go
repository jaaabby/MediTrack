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
	Phone                  *string           `json:"phone,omitempty" db:"phone"`
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
	CreatedAt              int64             `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt              int64             `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

// UserResponse representa la respuesta del usuario sin contraseña
type UserResponse struct {
	RUT                string            `json:"rut"`
	Name               string            `json:"name"`
	Email              string            `json:"email"`
	Phone              *string           `json:"phone,omitempty"`
	Role               string            `json:"role"`
	MedicalCenterID    int               `json:"medical_center_id"`
	MedicalCenter      *MedicalCenter    `json:"medical_center,omitempty"`
	PavilionID         *int              `json:"pavilion_id,omitempty"`
	SpecialtyID        *int              `json:"specialty_id,omitempty"`
	Specialty          *MedicalSpecialty `json:"specialty,omitempty"`
	IsActive           bool              `json:"is_active"`
	MustChangePassword bool              `json:"must_change_password"`
	CreatedAt          int64             `json:"created_at"`
	UpdatedAt          int64             `json:"updated_at"`
}

// ToResponse convierte un User a UserResponse
func (u User) ToResponse() UserResponse {
	return UserResponse{
		RUT:                u.RUT,
		Name:               u.Name,
		Email:              u.Email,
		Phone:              u.Phone,
		Role:               u.Role,
		MedicalCenterID:    u.MedicalCenterID,
		MedicalCenter:      u.MedicalCenter,
		PavilionID:         u.PavilionID,
		SpecialtyID:        u.SpecialtyID,
		Specialty:          u.Specialty,
		IsActive:           u.IsActive,
		MustChangePassword: u.MustChangePassword,
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
	}
}

// OtpSession representa una sesión de verificación OTP para el login
type OtpSession struct {
	ID        string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserRUT   string `json:"user_rut" gorm:"not null"`
	Code      string `json:"-" gorm:"not null"`
	ExpiresAt int64  `json:"expires_at" gorm:"not null"`
	Used      bool   `json:"used" gorm:"default:false"`
	Attempts  int    `json:"attempts" gorm:"default:0"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
}

// MaxOTPAttempts es el número máximo de intentos fallidos antes de invalidar la sesión OTP
const MaxOTPAttempts = 5

func (OtpSession) TableName() string {
	return "otp_session"
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
