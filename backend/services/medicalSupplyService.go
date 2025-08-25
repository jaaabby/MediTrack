package services

import (
	"fmt"
	"meditrack/models"

	"gorm.io/gorm"
)

type MedicalSupplyService struct {
	DB        *gorm.DB
	QRService *QRService
}

// CreateMultipleIndividualSuppliesTx crea múltiples insumos individuales usando una transacción existente
func (s *MedicalSupplyService) CreateMultipleIndividualSuppliesTx(tx *gorm.DB, batchID int, code int, quantity int) ([]models.MedicalSupply, error) {
	var supplies []models.MedicalSupply
	for i := 0; i < quantity; i++ {
		supply := models.MedicalSupply{
			Code:    code,
			BatchID: batchID,
		}
		// Generar QR único para cada insumo individual si hay QRService
		if s.QRService != nil {
			qrCode, err := s.QRService.GenerateUniqueQRCode("SUPPLY")
			if err != nil {
				return nil, fmt.Errorf("error generando QR para insumo %d: %v", i+1, err)
			}
			supply.QRCode = qrCode
		}
		if err := tx.Create(&supply).Error; err != nil {
			return nil, fmt.Errorf("error creando insumo %d: %v", i+1, err)
		}
		supplies = append(supplies, supply)
	}
	return supplies, nil
}

func NewMedicalSupplyService(db *gorm.DB, qrService *QRService) *MedicalSupplyService {
	return &MedicalSupplyService{
		DB:        db,
		QRService: qrService,
	}
}

// ===== FUNCIONALIDADES BÁSICAS (de la versión anterior) =====

func (s *MedicalSupplyService) CreateMedicalSupply(supply *models.MedicalSupply) error {
	// Si hay QRService disponible, generar QR automáticamente
	if s.QRService != nil && supply.QRCode == "" {
		qrCode, err := s.QRService.GenerateUniqueQRCode("SUPPLY")
		if err != nil {
			return fmt.Errorf("error generando código QR: %v", err)
		}
		supply.QRCode = qrCode
	}

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

	// Actualizar campos pero mantener el QR code original si existe
	supply.Code = newSupply.Code
	supply.BatchID = newSupply.BatchID
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

// ===== INVENTARIO BÁSICO (de la versión anterior - RESTAURADO) =====

// GetInventoryList obtiene el inventario básico combinando datos de batch y supplyCode
// VERSIÓN ANTERIOR RESTAURADA - funcionalidad simple y confiable
func (s *MedicalSupplyService) GetInventoryList() ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	// Query simple de la versión anterior que funcionaba
	query := `
		SELECT DISTINCT ON (b.id)
			b.id as batch_id,
			b.expiration_date,
			b.amount,
			b.supplier,
			b.store_id,
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
		var storeID int
		var supplyCode *int
		var supplyName *string

		err := rows.Scan(&batchID, &expirationDate, &amount, &supplier, &storeID, &supplyCode, &supplyName)
		if err != nil {
			return nil, err
		}

		item["batch_id"] = batchID
		item["expiration_date"] = expirationDate
		item["amount"] = amount
		item["supplier"] = supplier
		item["store_id"] = storeID
		item["code"] = supplyCode
		item["name"] = supplyName

		result = append(result, item)
	}

	return result, nil
}

// ===== FUNCIONALIDADES AVANZADAS (de la versión actual - MANTENIDAS) =====

// GetInventoryListAdvanced obtiene el inventario con información avanzada de productos consumidos
// Para usar cuando necesites funcionalidades más complejas
func (s *MedicalSupplyService) GetInventoryListAdvanced() ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	query := `
		SELECT DISTINCT
			b.id as batch_id,
			b.qr_code as batch_qr_code,
			b.expiration_date,
			b.amount as current_amount,
			b.supplier,
			sc.code as supply_code,
			sc.name as supply_name,
			COUNT(ms.id) as total_individual_supplies,
			COUNT(CASE WHEN consumed_supplies.supply_id IS NOT NULL THEN 1 END) as consumed_supplies,
			(COUNT(ms.id) - COUNT(CASE WHEN consumed_supplies.supply_id IS NOT NULL THEN 1 END)) as available_supplies
		FROM batch b
		LEFT JOIN medical_supply ms ON b.id = ms.batch_id
		LEFT JOIN supply_code sc ON ms.code = sc.code
		LEFT JOIN (
			SELECT DISTINCT sh.medical_supply_id as supply_id
			FROM supply_history sh
			WHERE sh.status = 'consumido'
		) consumed_supplies ON ms.id = consumed_supplies.supply_id
		GROUP BY b.id, b.qr_code, b.expiration_date, b.amount, b.supplier, sc.code, sc.name
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
		var batchQRCode *string
		var expirationDate *string
		var currentAmount int
		var supplier string
		var supplyCode *int
		var supplyName *string
		var totalIndividualSupplies int
		var consumedSupplies int
		var availableSupplies int

		err := rows.Scan(
			&batchID, &batchQRCode, &expirationDate, &currentAmount, &supplier,
			&supplyCode, &supplyName, &totalIndividualSupplies,
			&consumedSupplies, &availableSupplies,
		)
		if err != nil {
			return nil, err
		}

		item["batch_id"] = batchID
		item["batch_qr_code"] = batchQRCode
		item["expiration_date"] = expirationDate
		item["current_amount"] = currentAmount
		item["supplier"] = supplier
		item["code"] = supplyCode
		item["name"] = supplyName
		item["total_individual_supplies"] = totalIndividualSupplies
		item["consumed_supplies"] = consumedSupplies
		item["available_supplies"] = availableSupplies
		item["sync_status"] = map[string]interface{}{
			"batch_amount_matches": currentAmount == availableSupplies,
			"needs_sync":           currentAmount != availableSupplies,
		}

		result = append(result, item)
	}

	return result, nil
}

// CreateMultipleIndividualSupplies crea múltiples insumos individuales para un lote
func (s *MedicalSupplyService) CreateMultipleIndividualSupplies(batchID int, code int, quantity int) ([]models.MedicalSupply, error) {
	var supplies []models.MedicalSupply

	// Usar transacción para asegurar consistencia
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < quantity; i++ {
			supply := models.MedicalSupply{
				Code:    code,
				BatchID: batchID,
			}

			// Generar QR único para cada insumo individual si hay QRService
			if s.QRService != nil {
				qrCode, err := s.QRService.GenerateUniqueQRCode("SUPPLY")
				if err != nil {
					return fmt.Errorf("error generando QR para insumo %d: %v", i+1, err)
				}
				supply.QRCode = qrCode
			}

			if err := tx.Create(&supply).Error; err != nil {
				return fmt.Errorf("error creando insumo %d: %v", i+1, err)
			}

			supplies = append(supplies, supply)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return supplies, nil
}

// ConsumeSupplyByQR marca un insumo como consumido y actualiza la cantidad del lote
func (s *MedicalSupplyService) ConsumeSupplyByQR(qrCode string, userRUT string, destinationType string, destinationID int) (*models.MedicalSupply, error) {
	var supply models.MedicalSupply

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Buscar el insumo por QR
		if err := tx.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
			return fmt.Errorf("insumo no encontrado con QR %s: %v", qrCode, err)
		}

		// Verificar que el insumo no haya sido consumido previamente
		var historyCount int64
		tx.Model(&models.SupplyHistory{}).Where("medical_supply_id = ? AND status = ?", supply.ID, "consumido").Count(&historyCount)
		if historyCount > 0 {
			return fmt.Errorf("el insumo con QR %s ya ha sido consumido", qrCode)
		}

		// Crear historial de consumo
		history := models.SupplyHistory{
			DateTime:        models.CurrentTime(),
			Status:          "consumido",
			DestinationType: destinationType,
			DestinationID:   destinationID,
			MedicalSupplyID: supply.ID,
			UserRUT:         userRUT,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("error creando historial de consumo: %v", err)
		}

		// Actualizar cantidad del lote (restar 1)
		if err := tx.Model(&models.Batch{}).Where("id = ?", supply.BatchID).Update("amount", gorm.Expr("amount - 1")).Error; err != nil {
			return fmt.Errorf("error actualizando cantidad del lote: %v", err)
		}

		// Verificar que la cantidad del lote no sea negativa
		var batch models.Batch
		if err := tx.Where("id = ?", supply.BatchID).First(&batch).Error; err != nil {
			return fmt.Errorf("error obteniendo lote: %v", err)
		}

		if batch.Amount < 0 {
			return fmt.Errorf("no hay suficiente stock en el lote para consumir este insumo")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &supply, nil
}

// GetMedicalSupplyByQR obtiene un insumo médico por su código QR
func (s *MedicalSupplyService) GetMedicalSupplyByQR(qrCode string) (*models.MedicalSupply, error) {
	var supply models.MedicalSupply
	if err := s.DB.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
		return nil, err
	}
	return &supply, nil
}

// GetSupplyWithBatchInfo obtiene información completa de un insumo con datos del lote
func (s *MedicalSupplyService) GetSupplyWithBatchInfo(qrCode string) (map[string]interface{}, error) {
	query := `
		SELECT 
			ms.id as supply_id,
			ms.code as supply_code,
			ms.qr_code,
			ms.batch_id,
			sc.name as supply_name,
			sc.code_supplier,
			b.expiration_date,
			b.amount as batch_remaining_amount,
			b.supplier,
			st.name as store_name,
			st.type as store_type
		FROM medical_supply ms
		JOIN supply_code sc ON ms.code = sc.code
		JOIN batch b ON ms.batch_id = b.id
		JOIN store st ON b.store_id = st.id
		WHERE ms.qr_code = ?
	`

	row := s.DB.Raw(query, qrCode).Row()

	var supplyID int
	var supplyCode int
	var qr string
	var batchID int
	var supplyName string
	var codeSupplier int
	var expirationDate string
	var batchRemainingAmount int
	var supplier string
	var storeName string
	var storeType string

	err := row.Scan(
		&supplyID, &supplyCode, &qr, &batchID, &supplyName, &codeSupplier,
		&expirationDate, &batchRemainingAmount, &supplier, &storeName, &storeType,
	)

	if err != nil {
		return nil, fmt.Errorf("insumo no encontrado: %v", err)
	}

	// Verificar si ya fue consumido
	var consumedCount int64
	s.DB.Model(&models.SupplyHistory{}).Where("medical_supply_id = ? AND status = ?", supplyID, "consumido").Count(&consumedCount)

	result := map[string]interface{}{
		"supply_id":              supplyID,
		"supply_code":            supplyCode,
		"qr_code":                qr,
		"batch_id":               batchID,
		"supply_name":            supplyName,
		"code_supplier":          codeSupplier,
		"expiration_date":        expirationDate,
		"batch_remaining_amount": batchRemainingAmount,
		"supplier":               supplier,
		"store_name":             storeName,
		"store_type":             storeType,
		"is_consumed":            consumedCount > 0,
		"available_for_use":      consumedCount == 0 && batchRemainingAmount > 0,
	}

	return result, nil
}

// SyncBatchAmounts sincroniza las cantidades de los lotes con los productos individuales disponibles
func (s *MedicalSupplyService) SyncBatchAmounts() error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Obtener todos los lotes
		var batches []models.Batch
		if err := tx.Find(&batches).Error; err != nil {
			return fmt.Errorf("error obteniendo lotes: %v", err)
		}

		for _, batch := range batches {
			// Contar productos individuales disponibles (no consumidos)
			var availableCount int64
			subquery := tx.Model(&models.SupplyHistory{}).
				Select("medical_supply_id").
				Where("status = ?", "consumido")

			if err := tx.Model(&models.MedicalSupply{}).
				Where("batch_id = ? AND id NOT IN (?)", batch.ID, subquery).
				Count(&availableCount).Error; err != nil {
				return fmt.Errorf("error contando productos disponibles para lote %d: %v", batch.ID, err)
			}

			// Actualizar cantidad del lote si es diferente
			if int64(batch.Amount) != availableCount {
				if err := tx.Model(&batch).Update("amount", availableCount).Error; err != nil {
					return fmt.Errorf("error actualizando cantidad del lote %d: %v", batch.ID, err)
				}
			}
		}

		return nil
	})
}

// GetIndividualSuppliesByCode obtiene todos los insumos individuales de un supply code específico
func (s *MedicalSupplyService) GetIndividualSuppliesByCode(code int) ([]models.MedicalSupply, error) {
	var supplies []models.MedicalSupply
	if err := s.DB.Where("code = ?", code).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

// GetAvailableSuppliesByBatch obtiene todos los insumos individuales disponibles de un lote
func (s *MedicalSupplyService) GetAvailableSuppliesByBatch(batchID int) ([]models.MedicalSupply, error) {
	var supplies []models.MedicalSupply

	subquery := s.DB.Model(&models.SupplyHistory{}).
		Select("medical_supply_id").
		Where("status = ?", "consumido")

	if err := s.DB.Where("batch_id = ? AND id NOT IN (?)", batchID, subquery).Find(&supplies).Error; err != nil {
		return nil, err
	}

	return supplies, nil
}
