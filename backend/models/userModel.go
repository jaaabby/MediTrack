package models

// Constantes para los roles de usuario
const (
	RoleAdmin        = "admin"
	RolePavilion     = "pabellón"
	RoleStoreManager = "encargado de bodega"
	RoleNurse        = "enfermera"
	RoleDoctor       = "doctor"
	RolePavedad      = "pavedad"
	RoleConsignation = "consignación"
)

// User representa un usuario del sistema
type User struct {
	RUT             string            `json:"rut" gorm:"primaryKey"`
	Name            string            `json:"name" db:"name" gorm:"not null"`
	Email           string            `json:"email" db:"email"`
	Password        string            `json:"password" db:"password" gorm:"not null"`
	Role            string            `json:"role" db:"role" gorm:"not null;check:role IN ('admin', 'pabellón', 'encargado de bodega', 'enfermera', 'doctor', 'pavedad', 'consignación')"`
	MedicalCenterID int               `json:"medical_center_id" db:"medical_center_id" gorm:"not null"`
	MedicalCenter   *MedicalCenter    `json:"medical_center,omitempty" gorm:"foreignKey:MedicalCenterID"`
	SpecialtyID     *int              `json:"specialty_id" db:"specialty_id"`
	Specialty       *MedicalSpecialty `json:"specialty,omitempty" gorm:"foreignKey:SpecialtyID"`
	IsActive        bool              `json:"is_active" db:"is_active" gorm:"default:true"`
	CreatedAt       int64             `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       int64             `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

// UserResponse representa la respuesta del usuario sin contraseña
type UserResponse struct {
	RUT             string            `json:"rut"`
	Name            string            `json:"name"`
	Email           string            `json:"email"`
	Role            string            `json:"role"`
	MedicalCenterID int               `json:"medical_center_id"`
	MedicalCenter   *MedicalCenter    `json:"medical_center,omitempty"`
	SpecialtyID     *int              `json:"specialty_id"`
	Specialty       *MedicalSpecialty `json:"specialty,omitempty"`
	IsActive        bool              `json:"is_active"`
	CreatedAt       int64             `json:"created_at"`
	UpdatedAt       int64             `json:"updated_at"`
}

// ToResponse convierte un User a UserResponse
func (u User) ToResponse() UserResponse {
	return UserResponse{
		RUT:             u.RUT,
		Name:            u.Name,
		Email:           u.Email,
		Role:            u.Role,
		MedicalCenterID: u.MedicalCenterID,
		MedicalCenter:   u.MedicalCenter,
		SpecialtyID:     u.SpecialtyID,
		Specialty:       u.Specialty,
		IsActive:        u.IsActive,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

// IsValidRole verifica si el rol es válido
func (u User) IsValidRole() bool {
	return u.Role == RoleAdmin ||
		u.Role == RolePavilion ||
		u.Role == RoleStoreManager ||
		u.Role == RoleNurse ||
		u.Role == RoleDoctor ||
		u.Role == RolePavedad ||
		u.Role == RoleConsignation
}

func (u User) TableName() string {
	return "user"
}
