package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type MedicalSupplyService struct {
	DB        *gorm.DB
	QRService *QRService
}

func NewMedicalSupplyService(db *gorm.DB, qrService *QRService) *MedicalSupplyService {
	return &MedicalSupplyService{
		DB:        db,
		QRService: qrService,
	}
}

func (s *MedicalSupplyService) CreateMedicalSupply(supply *models.MedicalSupply) error {
	// Generar código QR único para el insumo médico
	qrCode, err := s.QRService.GenerateMedicalSupplyQRCode()
	if err != nil {
		return err
	}
	supply.QRCode = qrCode

	return s.DB.Create(supply).Error
}

func (s *MedicalSupplyService) GetMedicalSupplyByID(id int) (*models.MedicalSupply, error) {
	var supply models.MedicalSupply
	if err := s.DB.First(&supply, id).Error; err != nil {
		return nil, err
	}
	return &supply, nil
}

func (s *MedicalSupplyService) GetMedicalSupplyByQRCode(qrCode string) (*models.MedicalSupply, error) {
	var supply models.MedicalSupply
	if err := s.DB.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
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

	// Actualizar campos pero mantener el QR code original
	supply.Code = newSupply.Code
	// No actualizamos supply.QRCode para mantener la trazabilidad

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
		SELECT 
			b.id as batch_id,
			b.qr_code as batch_qr_code,
			sc.code as supply_code,
			sc.name as supply_name,
			b.expiration_date,
			b.amount,
			b.supplier,
			COUNT(ms.id) as individual_supplies_count
		FROM batch b
		JOIN supply_code sc ON b.id = sc.batch_id
		LEFT JOIN medical_supply ms ON sc.code = ms.code
		GROUP BY b.id, sc.code, sc.name, b.expiration_date, b.amount, b.supplier, b.qr_code
		ORDER BY b.id
	`

	rows, err := s.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item map[string]interface{} = make(map[string]interface{})
		var batchID int
		var batchQRCode string
		var supplyCode int
		var supplyName string
		var expirationDate string
		var amount int
		var supplier string
		var individualSuppliesCount int

		err := rows.Scan(&batchID, &batchQRCode, &supplyCode, &supplyName, &expirationDate, &amount, &supplier, &individualSuppliesCount)
		if err != nil {
			return nil, err
		}

		item["batch_id"] = batchID
		item["batch_qr_code"] = batchQRCode
		item["code"] = supplyCode
		item["name"] = supplyName
		item["expiration_date"] = expirationDate
		item["amount"] = amount
		item["supplier"] = supplier
		item["individual_supplies_count"] = individualSuppliesCount

		result = append(result, item)
	}

	return result, nil
}

// GetIndividualSuppliesByCode obtiene todos los insumos individuales de un supply code específico
func (s *MedicalSupplyService) GetIndividualSuppliesByCode(code int) ([]models.MedicalSupply, error) {
	var supplies []models.MedicalSupply
	if err := s.DB.Where("code = ?", code).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}
