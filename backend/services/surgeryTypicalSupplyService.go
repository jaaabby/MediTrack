package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type SurgeryTypicalSupplyService struct {
	DB *gorm.DB
}

func NewSurgeryTypicalSupplyService(db *gorm.DB) *SurgeryTypicalSupplyService {
	return &SurgeryTypicalSupplyService{DB: db}
}

// CreateSurgeryTypicalSupply crea un nuevo insumo típico para una cirugía
func (s *SurgeryTypicalSupplyService) CreateSurgeryTypicalSupply(typicalSupply *models.SurgeryTypicalSupply) error {
	return s.DB.Create(typicalSupply).Error
}

// GetSurgeryTypicalSupplyByID obtiene un insumo típico por ID
func (s *SurgeryTypicalSupplyService) GetSurgeryTypicalSupplyByID(id int) (*models.SurgeryTypicalSupply, error) {
	var typicalSupply models.SurgeryTypicalSupply
	if err := s.DB.Preload("Surgery").Preload("SupplyCodeInfo").First(&typicalSupply, id).Error; err != nil {
		return nil, err
	}
	return &typicalSupply, nil
}

// GetTypicalSuppliesBySurgeryID obtiene todos los insumos típicos de una cirugía
func (s *SurgeryTypicalSupplyService) GetTypicalSuppliesBySurgeryID(surgeryID int) ([]models.SurgeryTypicalSupply, error) {
	var typicalSupplies []models.SurgeryTypicalSupply
	if err := s.DB.Where("surgery_id = ?", surgeryID).
		Preload("SupplyCodeInfo").
		Order("is_required DESC, supply_code ASC").
		Find(&typicalSupplies).Error; err != nil {
		return nil, err
	}
	return typicalSupplies, nil
}

// GetSurgeriesBySupplyCode obtiene todas las cirugías que requieren un insumo específico
func (s *SurgeryTypicalSupplyService) GetSurgeriesBySupplyCode(supplyCode int) ([]models.SurgeryTypicalSupply, error) {
	var typicalSupplies []models.SurgeryTypicalSupply
	if err := s.DB.Where("supply_code = ?", supplyCode).
		Preload("Surgery").
		Preload("Surgery.Specialty").
		Order("surgery_id ASC").
		Find(&typicalSupplies).Error; err != nil {
		return nil, err
	}
	return typicalSupplies, nil
}

// UpdateSurgeryTypicalSupply actualiza un insumo típico
func (s *SurgeryTypicalSupplyService) UpdateSurgeryTypicalSupply(id int, typicalSupply *models.SurgeryTypicalSupply) (*models.SurgeryTypicalSupply, error) {
	var existingTypicalSupply models.SurgeryTypicalSupply
	if err := s.DB.First(&existingTypicalSupply, id).Error; err != nil {
		return nil, err
	}

	existingTypicalSupply.TypicalQuantity = typicalSupply.TypicalQuantity
	existingTypicalSupply.IsRequired = typicalSupply.IsRequired
	existingTypicalSupply.Notes = typicalSupply.Notes

	if err := s.DB.Save(&existingTypicalSupply).Error; err != nil {
		return nil, err
	}

	return &existingTypicalSupply, nil
}

// DeleteSurgeryTypicalSupply elimina un insumo típico de una cirugía
func (s *SurgeryTypicalSupplyService) DeleteSurgeryTypicalSupply(id int) error {
	return s.DB.Delete(&models.SurgeryTypicalSupply{}, id).Error
}

// DeleteTypicalSuppliesBySurgeryID elimina todos los insumos típicos de una cirugía
func (s *SurgeryTypicalSupplyService) DeleteTypicalSuppliesBySurgeryID(surgeryID int) error {
	return s.DB.Where("surgery_id = ?", surgeryID).Delete(&models.SurgeryTypicalSupply{}).Error
}

// BulkCreateSurgeryTypicalSupplies crea múltiples insumos típicos para una cirugía
func (s *SurgeryTypicalSupplyService) BulkCreateSurgeryTypicalSupplies(surgeryID int, typicalSupplies []models.SurgeryTypicalSupply) error {
	// Verificar que todos los insumos típicos pertenezcan a la misma cirugía
	for i := range typicalSupplies {
		typicalSupplies[i].SurgeryID = surgeryID
	}
	return s.DB.Create(&typicalSupplies).Error
}

