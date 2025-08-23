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

func (s *SupplyHistoryService) GetAllSupplyHistories() ([]models.SupplyHistory, error) {
	var histories []models.SupplyHistory
	if err := s.DB.Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

// UpdateSupplyHistory actualiza un SupplyHistory por ID
func (s *SupplyHistoryService) UpdateSupplyHistory(id int, history *models.SupplyHistory) error {
	var existing models.SupplyHistory
	if err := s.DB.First(&existing, id).Error; err != nil {
		return err
	}
	// Actualiza los campos necesarios
	if err := s.DB.Model(&existing).Updates(history).Error; err != nil {
		return err
	}
	return nil
}
