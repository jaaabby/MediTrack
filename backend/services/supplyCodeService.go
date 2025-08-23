package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type SupplyCodeService struct {
	DB *gorm.DB
}

func NewSupplyCodeService(db *gorm.DB) *SupplyCodeService {
	return &SupplyCodeService{DB: db}
}

func (s *SupplyCodeService) CreateSupplyCode(supplyCode *models.SupplyCode) error {
	return s.DB.Create(supplyCode).Error
}

func (s *SupplyCodeService) GetSupplyCodeByID(id int) (*models.SupplyCode, error) {
	var supplyCode models.SupplyCode
	if err := s.DB.First(&supplyCode, id).Error; err != nil {
		return nil, err
	}
	return &supplyCode, nil
}

func (s *SupplyCodeService) GetAllSupplyCodes() ([]models.SupplyCode, error) {
	var supplyCodes []models.SupplyCode
	if err := s.DB.Find(&supplyCodes).Error; err != nil {
		return nil, err
	}
	return supplyCodes, nil
}

func (s *SupplyCodeService) UpdateSupplyCode(id int, newSupplyCode *models.SupplyCode) (*models.SupplyCode, error) {
	var supplyCode models.SupplyCode
	if err := s.DB.First(&supplyCode, id).Error; err != nil {
		return nil, err
	}
	supplyCode.Name = newSupplyCode.Name
	supplyCode.CodeSupplier = newSupplyCode.CodeSupplier

	if err := s.DB.Save(&supplyCode).Error; err != nil {
		return nil, err
	}
	return &supplyCode, nil
}

func (s *SupplyCodeService) DeleteSupplyCode(id int) error {
	if err := s.DB.Delete(&models.SupplyCode{}, id).Error; err != nil {
		return err
	}
	return nil
}
