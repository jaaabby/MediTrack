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
	if err := s.DB.Create(pavilion).Error; err != nil {
		return err
	}
	return nil
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
	if err := s.DB.Find(&pavilions).Error; err != nil {
		return nil, err
	}
	return pavilions, nil
}

func (s *PavilionService) UpdatePavilion(id int, newPavilion *models.Pavilion) (*models.Pavilion, error) {
	var pavilion models.Pavilion
	if err := s.DB.First(&pavilion, id).Error; err != nil {
		return nil, err
	}

	pavilion.Name = newPavilion.Name
	pavilion.MedicalCenterID = newPavilion.MedicalCenterID

	if err := s.DB.Save(&pavilion).Error; err != nil {
		return nil, err
	}
	return &pavilion, nil
}

func (s *PavilionService) DeletePavilion(id int) error {
	if err := s.DB.Delete(&models.Pavilion{}, id).Error; err != nil {
		return err
	}
	return nil
}
