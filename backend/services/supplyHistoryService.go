package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type SupplyHistoryService struct {
	DB *gorm.DB
}

func NewSupplyHistoryService(db *gorm.DB) *SupplyHistoryService {
	return &SupplyHistoryService{DB: db}
}

func (s *SupplyHistoryService) CreateSupplyHistory(history *models.SupplyHistory) error {
	return s.DB.Create(history).Error
}

func (s *SupplyHistoryService) GetSupplyHistoryByID(id int) (*models.SupplyHistory, error) {
	var history models.SupplyHistory
	if err := s.DB.First(&history, id).Error; err != nil {
		return nil, err
	}
	return &history, nil
}

func (s *SupplyHistoryService) DeleteSupplyHistory(id int) error {
	if err := s.DB.Delete(&models.SupplyHistory{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *SupplyHistoryService) UpdateSupplyHistory(id int, newHistory *models.SupplyHistory) (*models.SupplyHistory, error) {
	var history models.SupplyHistory
	if err := s.DB.First(&history, id).Error; err != nil {
		return nil, err
	}
	history.DateTime = newHistory.DateTime
	history.Status = newHistory.Status
	history.DestinationType = newHistory.DestinationType
	history.DestinationID = newHistory.DestinationID
	history.MedicalSupplyID = newHistory.MedicalSupplyID
	history.UserRUT = newHistory.UserRUT

	if err := s.DB.Save(&history).Error; err != nil {
		return nil, err
	}
	return &history, nil
}

func (s *SupplyHistoryService) GetAllSupplyHistories() ([]models.SupplyHistory, error) {
	var histories []models.SupplyHistory
	if err := s.DB.Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}
