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
	if err := s.DB.Create(center).Error; err != nil {
		return err
	}
	return nil
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
	center.Name = newCenter.Name
	center.Address = newCenter.Address
	center.Phone = newCenter.Phone
	center.Email = newCenter.Email

	if err := s.DB.Save(&center).Error; err != nil {
		return nil, err
	}
	return &center, nil
}

func (s *MedicalCenterService) DeleteMedicalCenter(id int) error {
	if err := s.DB.Delete(&models.MedicalCenter{}, id).Error; err != nil {
		return err
	}
	return nil
}
