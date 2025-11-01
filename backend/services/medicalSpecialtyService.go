package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type MedicalSpecialtyService struct {
	DB *gorm.DB
}

func NewMedicalSpecialtyService(db *gorm.DB) *MedicalSpecialtyService {
	return &MedicalSpecialtyService{DB: db}
}

// CreateMedicalSpecialty crea una nueva especialidad médica
func (s *MedicalSpecialtyService) CreateMedicalSpecialty(specialty *models.MedicalSpecialty) error {
	return s.DB.Create(specialty).Error
}

// GetMedicalSpecialtyByID obtiene una especialidad médica por ID
func (s *MedicalSpecialtyService) GetMedicalSpecialtyByID(id int) (*models.MedicalSpecialty, error) {
	var specialty models.MedicalSpecialty
	if err := s.DB.Preload("Surgeries").Preload("Doctors").First(&specialty, id).Error; err != nil {
		return nil, err
	}
	return &specialty, nil
}

// GetAllMedicalSpecialties obtiene todas las especialidades médicas
func (s *MedicalSpecialtyService) GetAllMedicalSpecialties() ([]models.MedicalSpecialty, error) {
	var specialties []models.MedicalSpecialty
	if err := s.DB.Where("is_active = ?", true).Order("name ASC").Find(&specialties).Error; err != nil {
		return nil, err
	}
	return specialties, nil
}

// GetMedicalSpecialtiesPaginated obtiene especialidades médicas con paginación
func (s *MedicalSpecialtyService) GetMedicalSpecialtiesPaginated(page int, pageSize int, search *string) ([]models.MedicalSpecialty, int64, error) {
	var specialties []models.MedicalSpecialty
	var total int64

	query := s.DB.Model(&models.MedicalSpecialty{}).Where("is_active = ?", true)

	// Aplicar búsqueda si se proporciona
	if search != nil && *search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ? OR code ILIKE ?", "%"+*search+"%", "%"+*search+"%", "%"+*search+"%")
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	offset := (page - 1) * pageSize
	if err := query.Order("name ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&specialties).Error; err != nil {
		return nil, 0, err
	}

	return specialties, total, nil
}

// UpdateMedicalSpecialty actualiza una especialidad médica
func (s *MedicalSpecialtyService) UpdateMedicalSpecialty(id int, specialty *models.MedicalSpecialty) (*models.MedicalSpecialty, error) {
	var existingSpecialty models.MedicalSpecialty
	if err := s.DB.First(&existingSpecialty, id).Error; err != nil {
		return nil, err
	}

	existingSpecialty.Name = specialty.Name
	existingSpecialty.Description = specialty.Description
	existingSpecialty.Code = specialty.Code
	existingSpecialty.IsActive = specialty.IsActive

	if err := s.DB.Save(&existingSpecialty).Error; err != nil {
		return nil, err
	}

	return &existingSpecialty, nil
}

// DeleteMedicalSpecialty elimina (desactiva) una especialidad médica
func (s *MedicalSpecialtyService) DeleteMedicalSpecialty(id int) error {
	// Verificar si hay cirugías o doctores asociados
	var surgeryCount int64
	var doctorCount int64
	s.DB.Model(&models.Surgery{}).Where("specialty_id = ?", id).Count(&surgeryCount)
	s.DB.Model(&models.DoctorInfo{}).Where("specialty_id = ?", id).Count(&doctorCount)

	if surgeryCount > 0 || doctorCount > 0 {
		// Si hay asociaciones, solo desactivar
		return s.DB.Model(&models.MedicalSpecialty{}).Where("id = ?", id).Update("is_active", false).Error
	}

	return s.DB.Delete(&models.MedicalSpecialty{}, id).Error
}

// SearchMedicalSpecialtiesByName busca especialidades médicas por nombre
func (s *MedicalSpecialtyService) SearchMedicalSpecialtiesByName(name string) ([]models.MedicalSpecialty, error) {
	var specialties []models.MedicalSpecialty
	if err := s.DB.Where("name ILIKE ? AND is_active = ?", "%"+name+"%", true).
		Order("name ASC").
		Find(&specialties).Error; err != nil {
		return nil, err
	}
	return specialties, nil
}

// GetMedicalSpecialtyByCode obtiene una especialidad médica por código
func (s *MedicalSpecialtyService) GetMedicalSpecialtyByCode(code string) (*models.MedicalSpecialty, error) {
	var specialty models.MedicalSpecialty
	if err := s.DB.Where("code = ? AND is_active = ?", code, true).First(&specialty).Error; err != nil {
		return nil, err
	}
	return &specialty, nil
}

