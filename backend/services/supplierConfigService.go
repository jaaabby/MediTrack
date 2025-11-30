package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type SupplierConfigService struct {
	DB *gorm.DB
}

func NewSupplierConfigService(db *gorm.DB) *SupplierConfigService {
	return &SupplierConfigService{DB: db}
}

// CreateSupplierConfig crea una nueva configuración de proveedor
func (s *SupplierConfigService) CreateSupplierConfig(config *models.SupplierConfig) error {
	return s.DB.Create(config).Error
}

// GetSupplierConfig obtiene la configuración de un proveedor por nombre
func (s *SupplierConfigService) GetSupplierConfig(supplierName string) (*models.SupplierConfig, error) {
	var config models.SupplierConfig
	if err := s.DB.Where("supplier_name = ?", supplierName).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// GetAllSupplierConfigs obtiene todas las configuraciones de proveedores
func (s *SupplierConfigService) GetAllSupplierConfigs() ([]models.SupplierConfig, error) {
	var configs []models.SupplierConfig
	if err := s.DB.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// UpdateSupplierConfig actualiza la configuración de un proveedor
func (s *SupplierConfigService) UpdateSupplierConfig(supplierName string, config *models.SupplierConfig) (*models.SupplierConfig, error) {
	var existingConfig models.SupplierConfig
	if err := s.DB.Where("supplier_name = ?", supplierName).First(&existingConfig).Error; err != nil {
		return nil, err
	}

	// Actualizar campos omitiendo SupplierName (primary key) y timestamps
	if err := s.DB.Model(&existingConfig).Omit("supplier_name", "created_at", "updated_at").Updates(config).Error; err != nil {
		return nil, err
	}

	return &existingConfig, nil
}

// DeleteSupplierConfig elimina la configuración de un proveedor
func (s *SupplierConfigService) DeleteSupplierConfig(supplierName string) error {
	return s.DB.Where("supplier_name = ?", supplierName).Delete(&models.SupplierConfig{}).Error
}

