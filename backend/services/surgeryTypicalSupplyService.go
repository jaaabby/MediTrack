package services

import (
	"fmt"
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
	if err := s.DB.Create(typicalSupply).Error; err != nil {
		return err
	}
	
	// Actualizar store_inventory_summary para lotes existentes que tienen este supply_code y no tienen surgery_id
	s.syncSurgeryIDForSupplyCode(typicalSupply.SupplyCode, typicalSupply.SurgeryID)
	
	return nil
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

	if err := s.DB.Model(&existingTypicalSupply).Omit("id", "surgery_id", "supply_code", "created_at", "updated_at").Updates(typicalSupply).Error; err != nil {
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

// GetAllTypicalSupplies obtiene todos los insumos típicos
func (s *SurgeryTypicalSupplyService) GetAllTypicalSupplies() ([]models.SurgeryTypicalSupply, error) {
	var typicalSupplies []models.SurgeryTypicalSupply
	if err := s.DB.Find(&typicalSupplies).Error; err != nil {
		return nil, err
	}
	return typicalSupplies, nil
}

// GetTypicalSuppliesCount obtiene el conteo total de insumos típicos
func (s *SurgeryTypicalSupplyService) GetTypicalSuppliesCount() (int64, error) {
	var count int64
	if err := s.DB.Model(&models.SurgeryTypicalSupply{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// BulkCreateSurgeryTypicalSupplies crea múltiples insumos típicos para una cirugía
func (s *SurgeryTypicalSupplyService) BulkCreateSurgeryTypicalSupplies(surgeryID int, typicalSupplies []models.SurgeryTypicalSupply) error {
	// Verificar que todos los insumos típicos pertenezcan a la misma cirugía
	for i := range typicalSupplies {
		typicalSupplies[i].SurgeryID = surgeryID
	}
	if err := s.DB.Create(&typicalSupplies).Error; err != nil {
		return err
	}
	
	// Actualizar store_inventory_summary para todos los supply_code creados
	supplyCodes := make(map[int]bool)
	for _, ts := range typicalSupplies {
		supplyCodes[ts.SupplyCode] = true
	}
	for supplyCode := range supplyCodes {
		s.syncSurgeryIDForSupplyCode(supplyCode, surgeryID)
	}
	
	return nil
}

// syncSurgeryIDForSupplyCode actualiza el surgery_id en store_inventory_summary
// para lotes que tienen este supply_code y no tienen surgery_id asignado
func (s *SurgeryTypicalSupplyService) syncSurgeryIDForSupplyCode(supplyCode int, surgeryID int) {
	// Actualizar store_inventory_summary donde supply_code coincide y surgery_id es NULL
	result := s.DB.Model(&models.StoreInventorySummary{}).
		Where("supply_code = ? AND surgery_id IS NULL", supplyCode).
		Update("surgery_id", surgeryID)
	
	if result.Error != nil {
		// Log error pero no fallar
		fmt.Printf("Advertencia: Error actualizando surgery_id para supply_code %d: %v\n", supplyCode, result.Error)
	} else if result.RowsAffected > 0 {
		fmt.Printf("✅ Actualizado surgery_id=%d para %d resúmenes de inventario con supply_code=%d\n", 
			surgeryID, result.RowsAffected, supplyCode)
	}
}

