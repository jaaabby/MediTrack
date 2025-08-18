package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type MedicalSupplyService struct {
	DB *gorm.DB
}

// InventoryItem representa un item del inventario con información completa
type InventoryItem struct {
	ID              int    `json:"id"`
	Code            int    `json:"code"`
	CodeSupplier    int    `json:"code_supplier"`
	Name            string `json:"name"`
	BatchID         int    `json:"batch_id"`
	ExpirationDate  string `json:"expiration_date"`
	Amount          int    `json:"amount"`
	Supplier        string `json:"supplier"`
	StoreName       string `json:"store_name"`
	StoreType       string `json:"store_type"`
	MedicalCenterID int    `json:"medical_center_id"`
}

func NewMedicalSupplyService(db *gorm.DB) *MedicalSupplyService {
	return &MedicalSupplyService{DB: db}
}

func (s *MedicalSupplyService) CreateMedicalSupply(supply *models.MedicalSupply) error {
	return s.DB.Create(supply).Error
}

func (s *MedicalSupplyService) GetMedicalSupplyByID(id int) (*models.MedicalSupply, error) {
	var supply models.MedicalSupply
	if err := s.DB.First(&supply, id).Error; err != nil {
		return nil, err
	}
	return &supply, nil
}

func (s *MedicalSupplyService) GetAllMedicalSupplies() ([]models.MedicalSupply, error) {
	var supplies []models.MedicalSupply
	if err := s.DB.Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (s *MedicalSupplyService) UpdateMedicalSupply(id int, newSupply *models.MedicalSupply) (*models.MedicalSupply, error) {
	var supply models.MedicalSupply
	if err := s.DB.First(&supply, id).Error; err != nil {
		return nil, err
	}
	supply.Name = newSupply.Name
	supply.BatchID = newSupply.BatchID

	if err := s.DB.Save(&supply).Error; err != nil {
		return nil, err
	}
	return &supply, nil
}

func (s *MedicalSupplyService) DeleteMedicalSupply(id int) error {
	if err := s.DB.Delete(&models.MedicalSupply{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *MedicalSupplyService) GetInventoryWithDetails() ([]InventoryItem, error) {
	var items []InventoryItem
	if err := s.DB.Table("medical_supply").
		Joins("LEFT JOIN batch ON medical_supply.batch_id = batch.id").
		Joins("LEFT JOIN store ON batch.store_id = store.id").
		Joins("LEFT JOIN medical_center ON store.medical_center_id = medical_center.id").
		Select(`
			medical_supply.id,
			medical_supply.code,
			medical_supply.code_supplier,
			medical_supply.name,
			batch.expiration_date,
			batch.amount,
			batch.supplier,
			store.name as store_name,
			store.type as store_type,
			medical_center.name as medical_center_name
		`).
		Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
