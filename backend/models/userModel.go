package models

// Constantes para los roles de usuario
const (
	RoleAdmin        = "admin"
	RolePavilion     = "pabellón"
	RoleStoreManager = "encargado de bodega"
)

// User representa un usuario del sistema
type User struct {
	RUT             string `json:"rut" gorm:"primaryKey"`
	Name            string `json:"name" db:"name" gorm:"not null"`
	Email           string `json:"email" db:"email"`
	Password        string `json:"password" db:"password" gorm:"not null"`
	Role            string `json:"role" db:"role" gorm:"not null;check:role IN ('admin', 'pabellón', 'encargado de bodega')"`
	MedicalCenterID int    `json:"medical_center_id" db:"medical_center_id" gorm:"not null"`
	IsActive        bool   `json:"is_active" db:"is_active" gorm:"default:true"`
	CreatedAt       int64  `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       int64  `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

// UserResponse representa la respuesta del usuario sin contraseña
type UserResponse struct {
	RUT             string `json:"rut"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	MedicalCenterID int    `json:"medical_center_id"`
	IsActive        bool   `json:"is_active"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

// ToResponse convierte un User a UserResponse
func (u User) ToResponse() UserResponse {
	return UserResponse{
		RUT:             u.RUT,
		Name:            u.Name,
		Email:           u.Email,
		Role:            u.Role,
		MedicalCenterID: u.MedicalCenterID,
		IsActive:        u.IsActive,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

// IsValidRole verifica si el rol es válido
func (u User) IsValidRole() bool {
	return u.Role == RoleAdmin || u.Role == RolePavilion || u.Role == RoleStoreManager
}

func (u User) TableName() string {
	return "user"
}
