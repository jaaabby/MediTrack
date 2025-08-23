package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type MedicalSupplyService struct {
	DB *gorm.DB
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

	supply.Code = newSupply.Code

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

// GetInventoryList obtiene el inventario combinando datos de batch y supplyCode
func (s *MedicalSupplyService) GetInventoryList() ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	query := `
		SELECT DISTINCT ON (b.id)
			b.id as batch_id,
			b.expiration_date,
			b.amount,
			b.supplier,
			sc.code as supply_code,
			sc.name as supply_name
		FROM batch b
		LEFT JOIN medical_supply ms ON b.id = ms.batch_id
		LEFT JOIN supply_code sc ON ms.code = sc.code
		ORDER BY b.id, sc.code
	`

	rows, err := s.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item map[string]interface{} = make(map[string]interface{})
		var batchID int
		var expirationDate string
		var amount int
		var supplier string
		var supplyCode *int
		var supplyName *string

		err := rows.Scan(&batchID, &expirationDate, &amount, &supplier, &supplyCode, &supplyName)
		if err != nil {
			return nil, err
		}

		item["batch_id"] = batchID
		item["expiration_date"] = expirationDate
		item["amount"] = amount
		item["supplier"] = supplier
		item["code"] = supplyCode
		item["name"] = supplyName

		result = append(result, item)
	}

	return result, nil
}
