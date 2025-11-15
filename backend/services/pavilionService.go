package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type PavilionService struct {
	DB *gorm.DB
}

func NewPavilionService(db *gorm.DB) *PavilionService {
	return &PavilionService{DB: db}
}

func (s *PavilionService) CreatePavilion(pavilion *models.Pavilion) error {
	return s.DB.Create(pavilion).Error
}

func (s *PavilionService) GetPavilionByID(id int) (*models.Pavilion, error) {
	var pavilion models.Pavilion
	if err := s.DB.First(&pavilion, id).Error; err != nil {
		return nil, err
	}
	return &pavilion, nil
}

func (s *PavilionService) GetAllPavilions() ([]models.Pavilion, error) {
	var pavilions []models.Pavilion
	if err := s.DB.Preload("MedicalCenter").Find(&pavilions).Error; err != nil {
		return nil, err
	}
	return pavilions, nil
}

func (s *PavilionService) UpdatePavilion(id int, newPavilion *models.Pavilion) (*models.Pavilion, error) {
	var pavilion models.Pavilion
	if err := s.DB.First(&pavilion, id).Error; err != nil {
		return nil, err
	}

	// Actualizar campos omitiendo ID
	if err := s.DB.Model(&pavilion).Omit("id").Updates(newPavilion).Error; err != nil {
		return nil, err
	}

	return &pavilion, nil
}

func (s *PavilionService) DeletePavilion(id int) error {
	return s.DB.Delete(&models.Pavilion{}, id).Error
}
