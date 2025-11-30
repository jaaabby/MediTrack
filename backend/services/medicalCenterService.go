package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type MedicalCenterService struct {
	DB *gorm.DB
}

func NewMedicalCenterService(db *gorm.DB) *MedicalCenterService {
	return &MedicalCenterService{DB: db}
}

func (s *MedicalCenterService) CreateMedicalCenter(center *models.MedicalCenter) error {
	return s.DB.Create(center).Error
}

func (s *MedicalCenterService) GetMedicalCenterByID(id int) (*models.MedicalCenter, error) {
	var center models.MedicalCenter
	if err := s.DB.First(&center, id).Error; err != nil {
		return nil, err
	}
	return &center, nil
}

func (s *MedicalCenterService) GetAllMedicalCenters() ([]models.MedicalCenter, error) {
	var centers []models.MedicalCenter
	if err := s.DB.Find(&centers).Error; err != nil {
		return nil, err
	}
	return centers, nil
}

func (s *MedicalCenterService) UpdateMedicalCenter(id int, newCenter *models.MedicalCenter) (*models.MedicalCenter, error) {
	var center models.MedicalCenter
	if err := s.DB.First(&center, id).Error; err != nil {
		return nil, err
	}

	// Actualizar campos omitiendo ID
	if err := s.DB.Model(&center).Omit("id").Updates(newCenter).Error; err != nil {
		return nil, err
	}

	return &center, nil
}

func (s *MedicalCenterService) DeleteMedicalCenter(id int) error {
	return s.DB.Delete(&models.MedicalCenter{}, id).Error
}
