package services

import (
	"fmt"
	"meditrack/models"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type DoctorInfoService struct {
	DB *gorm.DB
}

func NewDoctorInfoService(db *gorm.DB) *DoctorInfoService {
	return &DoctorInfoService{DB: db}
}

// CreateDoctor crea un nuevo usuario doctor
func (s *DoctorInfoService) CreateDoctor(doctor *models.User) error {
	// Asegurar que el rol sea doctor
	doctor.Role = models.RoleDoctor
	
	// Hashear la contraseña si se proporciona
	if doctor.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error al hashear contraseña: %v", err)
		}
		doctor.Password = string(hashedPassword)
	}
	
	return s.DB.Create(doctor).Error
}

// GetDoctorByRUT obtiene un doctor por RUT
func (s *DoctorInfoService) GetDoctorByRUT(rut string) (*models.User, error) {
	var doctor models.User
	if err := s.DB.Preload("MedicalCenter").Preload("Specialty").
		Where("rut = ? AND role = ?", rut, models.RoleDoctor).
		First(&doctor).Error; err != nil {
		return nil, err
	}
	return &doctor, nil
}

// GetAllDoctors obtiene todos los doctores
func (s *DoctorInfoService) GetAllDoctors() ([]models.User, error) {
	var doctors []models.User
	if err := s.DB.Preload("MedicalCenter").Preload("Specialty").
		Where("role = ? AND is_active = ?", models.RoleDoctor, true).
		Order("name ASC").
		Find(&doctors).Error; err != nil {
		return nil, err
	}
	return doctors, nil
}

// GetDoctorsBySpecialtyID obtiene todos los doctores de una especialidad
func (s *DoctorInfoService) GetDoctorsBySpecialtyID(specialtyID int) ([]models.User, error) {
	var doctors []models.User
	if err := s.DB.Preload("MedicalCenter").Preload("Specialty").
		Where("role = ? AND specialty_id = ? AND is_active = ?", models.RoleDoctor, specialtyID, true).
		Order("name ASC").
		Find(&doctors).Error; err != nil {
		return nil, err
	}
	return doctors, nil
}

// GetDoctorsPaginated obtiene doctores con paginación
func (s *DoctorInfoService) GetDoctorsPaginated(page int, pageSize int, search *string, specialtyID *int) ([]models.User, int64, error) {
	var doctors []models.User
	var total int64

	query := s.DB.Model(&models.User{}).Where("role = ? AND is_active = ?", models.RoleDoctor, true)

	// Filtrar por especialidad si se proporciona
	if specialtyID != nil {
		query = query.Where("specialty_id = ?", *specialtyID)
	}

	// Aplicar búsqueda si se proporciona
	if search != nil && *search != "" {
		searchTerm := "%" + strings.ToLower(*search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(rut) LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	offset := (page - 1) * pageSize
	if err := query.Preload("MedicalCenter").Preload("Specialty").
		Order("name ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&doctors).Error; err != nil {
		return nil, 0, err
	}

	return doctors, total, nil
}

// UpdateDoctor actualiza información de un doctor
func (s *DoctorInfoService) UpdateDoctor(rut string, doctor *models.User) (*models.User, error) {
	var existingDoctor models.User
	if err := s.DB.Where("rut = ? AND role = ?", rut, models.RoleDoctor).First(&existingDoctor).Error; err != nil {
		return nil, fmt.Errorf("doctor no encontrado: %v", err)
	}

	// Asegurar que el rol siga siendo doctor
	doctor.Role = models.RoleDoctor

	// Actualizar campos omitiendo rut, role, password, created_at y updated_at
	updates := map[string]interface{}{
		"name":             doctor.Name,
		"email":            doctor.Email,
		"medical_center_id": doctor.MedicalCenterID,
		"specialty_id":      doctor.SpecialtyID,
		"is_active":         doctor.IsActive,
	}

	if err := s.DB.Model(&existingDoctor).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Recargar con relaciones para retornar datos completos
	if err := s.DB.Preload("MedicalCenter").Preload("Specialty").
		Where("rut = ?", rut).
		First(&existingDoctor).Error; err != nil {
		return nil, err
	}

	return &existingDoctor, nil
}

// DeleteDoctor elimina un doctor (desactiva en lugar de eliminar físicamente)
func (s *DoctorInfoService) DeleteDoctor(rut string) error {
	return s.DB.Model(&models.User{}).
		Where("rut = ? AND role = ?", rut, models.RoleDoctor).
		Update("is_active", false).Error
}

// Métodos de compatibilidad con la API anterior (deprecated, pero mantenidos para compatibilidad)
// Estos métodos ahora trabajan con la tabla user directamente

// CreateDoctorInfo crea información extendida de un doctor (compatibilidad)
func (s *DoctorInfoService) CreateDoctorInfo(doctorInfo *models.DoctorInfo) error {
	// Este método ya no se usa, pero lo mantenemos para compatibilidad
	// En su lugar, se debe usar CreateDoctor
	return fmt.Errorf("este método está deprecado, use CreateDoctor en su lugar")
}

// GetDoctorInfoByRUT obtiene información de un doctor por RUT (compatibilidad)
func (s *DoctorInfoService) GetDoctorInfoByRUT(rut string) (*models.DoctorInfo, error) {
	// Este método ya no se usa, pero lo mantenemos para compatibilidad
	return nil, fmt.Errorf("este método está deprecado, use GetDoctorByRUT en su lugar")
}

// GetAllDoctorInfo obtiene información de todos los doctores (compatibilidad)
func (s *DoctorInfoService) GetAllDoctorInfo() ([]models.DoctorInfo, error) {
	// Este método ya no se usa, pero lo mantenemos para compatibilidad
	return nil, fmt.Errorf("este método está deprecado, use GetAllDoctors en su lugar")
}

// GetDoctorsBySpecialtyID ya está implementado arriba

// GetDoctorsPaginated ya está implementado arriba

// UpdateDoctorInfo actualiza información de un doctor (compatibilidad)
func (s *DoctorInfoService) UpdateDoctorInfo(rut string, doctorInfo *models.DoctorInfo) (*models.DoctorInfo, error) {
	// Este método ya no se usa, pero lo mantenemos para compatibilidad
	return nil, fmt.Errorf("este método está deprecado, use UpdateDoctor en su lugar")
}

// DeleteDoctorInfo elimina información de un doctor (compatibilidad)
func (s *DoctorInfoService) DeleteDoctorInfo(rut string) error {
	// Este método ya no se usa, pero lo mantenemos para compatibilidad
	return s.DeleteDoctor(rut)
}
