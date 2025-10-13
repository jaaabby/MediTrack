package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type SurgeryService struct {
	DB *gorm.DB
}

func NewSurgeryService(db *gorm.DB) *SurgeryService {
	return &SurgeryService{DB: db}
}

// CreateSurgery crea un nuevo tipo de cirugía
func (s *SurgeryService) CreateSurgery(surgery *models.Surgery) error {
	return s.DB.Create(surgery).Error
}

// GetSurgeryByID obtiene un tipo de cirugía por ID
func (s *SurgeryService) GetSurgeryByID(id int) (*models.Surgery, error) {
	var surgery models.Surgery
	if err := s.DB.First(&surgery, id).Error; err != nil {
		return nil, err
	}
	return &surgery, nil
}

// GetAllSurgeries obtiene todos los tipos de cirugía
func (s *SurgeryService) GetAllSurgeries() ([]models.Surgery, error) {
	var surgeries []models.Surgery
	if err := s.DB.Order("name ASC").Find(&surgeries).Error; err != nil {
		return nil, err
	}
	return surgeries, nil
}

// GetSurgeriesPaginated obtiene tipos de cirugía con paginación
func (s *SurgeryService) GetSurgeriesPaginated(page int, pageSize int, search *string) ([]models.Surgery, int64, error) {
	var surgeries []models.Surgery
	var total int64

	query := s.DB.Model(&models.Surgery{})

	// Aplicar búsqueda si se proporciona
	if search != nil && *search != "" {
		query = query.Where("name ILIKE ?", "%"+*search+"%")
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
		Find(&surgeries).Error; err != nil {
		return nil, 0, err
	}

	return surgeries, total, nil
}

// UpdateSurgery actualiza un tipo de cirugía
func (s *SurgeryService) UpdateSurgery(id int, surgery *models.Surgery) (*models.Surgery, error) {
	var existingSurgery models.Surgery
	if err := s.DB.First(&existingSurgery, id).Error; err != nil {
		return nil, err
	}

	existingSurgery.Name = surgery.Name
	existingSurgery.Duration = surgery.Duration

	if err := s.DB.Save(&existingSurgery).Error; err != nil {
		return nil, err
	}

	return &existingSurgery, nil
}

// DeleteSurgery elimina un tipo de cirugía
func (s *SurgeryService) DeleteSurgery(id int) error {
	// Verificar si hay lotes asociados
	var count int64
	s.DB.Model(&models.Batch{}).Where("surgery_id = ?", id).Count(&count)

	if count > 0 {
		// Si hay lotes asociados, establecer surgery_id a NULL en lugar de eliminar
		return s.DB.Model(&models.Batch{}).Where("surgery_id = ?", id).Update("surgery_id", nil).Error
	}

	return s.DB.Delete(&models.Surgery{}, id).Error
}

// SearchSurgeriesByName busca tipos de cirugía por nombre
func (s *SurgeryService) SearchSurgeriesByName(name string) ([]models.Surgery, error) {
	var surgeries []models.Surgery
	if err := s.DB.Where("name ILIKE ?", "%"+name+"%").
		Order("name ASC").
		Find(&surgeries).Error; err != nil {
		return nil, err
	}
	return surgeries, nil
}
