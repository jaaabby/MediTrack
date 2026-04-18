package services

import (
	"errors"
	"fmt"
	"log"
	"meditrack/mailer"
	"meditrack/models"
	"os"
	"meditrack/pkg/clinicconfig"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Constantes
const (
	DefaultUserRUT               = "12345678-9"
	DefaultExpirationAlertDays   = 90
	DefaultLowStockCheckInterval = 24 * time.Hour
)

type BatchService struct {
	DB                   *gorm.DB
	QRService            *QRService
	MedicalSupplyService *MedicalSupplyService
	BatchHistoryService  *BatchHistoryService
}

func NewBatchService(db *gorm.DB, qrService *QRService, medicalSupplyService *MedicalSupplyService, batchHistoryService *BatchHistoryService) *BatchService {
	return &BatchService{
		DB:                   db,
		QRService:            qrService,
		MedicalSupplyService: medicalSupplyService,
		BatchHistoryService:  batchHistoryService,
	}
}

func (s *BatchService) SetMedicalSupplyService(medicalSupplyService *MedicalSupplyService) {
	s.MedicalSupplyService = medicalSupplyService
}

func (s *BatchService) SetBatchHistoryService(batchHistoryService *BatchHistoryService) {
	s.BatchHistoryService = batchHistoryService
}

// ========================
// CRUD BÁSICO
// ========================

func (s *BatchService) CreateBatch(batch *models.Batch) error {
	batch.ID = 0
	batch.QRCode = ""

	// La ubicación inicial siempre es la bodega de origen
	if batch.LocationID == 0 {
		batch.LocationType = models.BatchLocationStore
		batch.LocationID = batch.StoreID
	}

	// Resolver nombre de proveedor a ID antes de guardar en BD
	if batch.Supplier != "" && batch.SupplierID == 0 {
		supplierID, err := s.ensureSupplierConfig(batch.Supplier)
		if err != nil {
			return fmt.Errorf("error resolviendo proveedor '%s': %v", batch.Supplier, err)
		}
		batch.SupplierID = supplierID
	}

	if err := s.DB.Create(batch).Error; err != nil {
		return fmt.Errorf("error creating batch: %v", err)
	}

	batch.QRCode = fmt.Sprintf("BATCH_%d", batch.ID)
	if err := s.DB.Model(batch).Update("qr_code", batch.QRCode).Error; err != nil {
		log.Printf("Warning: Could not update QR code for batch %d: %v", batch.ID, err)
	}

	return nil
}

func (s *BatchService) GetBatchByID(id int) (*models.Batch, error) {
	var batch models.Batch
	if err := s.DB.Preload("SupplierConfig").First(&batch, id).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

func (s *BatchService) GetAllBatches() ([]models.Batch, error) {
	var batches []models.Batch
	if err := s.DB.Find(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

func (s *BatchService) GetBatchesWithFilters(surgeryID *int, storeID *int, supplier string) ([]models.Batch, error) {
	query := s.DB.Model(&models.Batch{}).Preload("SupplierConfig")

	if surgeryID != nil {
		query = query.Where("surgery_id = ?", *surgeryID)
	}
	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}
	if supplier != "" {
		query = query.Joins("JOIN supplier_config ON supplier_config.id = batch.supplier_id").
			Where("supplier_config.supplier_name ILIKE ?", "%"+supplier+"%")
	}

	var batches []models.Batch
	if err := query.Find(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

func (s *BatchService) GetBatchWithSupplyInfo(id int) (map[string]interface{}, error) {
	var batch models.Batch
	if err := s.DB.First(&batch, id).Error; err != nil {
		return nil, fmt.Errorf("lote no encontrado: %v", err)
	}

	var totalSupplies, consumedSupplies int64
	if err := s.DB.Model(&models.MedicalSupply{}).Where("batch_id = ?", id).Count(&totalSupplies).Error; err != nil {
		return nil, fmt.Errorf("error contando insumos del lote: %v", err)
	}

	subquery := s.DB.Model(&models.SupplyHistory{}).
		Select("medical_supply_id").
		Where("status = ?", "consumido")
	if err := s.DB.Model(&models.MedicalSupply{}).
		Where("batch_id = ? AND id IN (?)", id, subquery).
		Count(&consumedSupplies).Error; err != nil {
		return nil, fmt.Errorf("error contando insumos consumidos: %v", err)
	}

	var supplyCode models.SupplyCode
	if err := s.DB.Where("code = ?", batch.SupplyCode).First(&supplyCode).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo código de insumo: %v", err)
	}

	availableSupplies := totalSupplies - consumedSupplies
	consumptionPercentage := float64(0)
	if totalSupplies > 0 {
		consumptionPercentage = float64(consumedSupplies) / float64(totalSupplies) * 100
	}

	return map[string]interface{}{
		"batch":                     batch,
		"supply_code":               supplyCode,
		"total_individual_supplies": totalSupplies,
		"consumed_supplies":         consumedSupplies,
		"available_supplies":        availableSupplies,
		"batch_amount":              batch.Amount,
		"amounts_synchronized":      batch.Amount == int(availableSupplies),
		"consumption_percentage":    consumptionPercentage,
	}, nil
}

func (s *BatchService) UpdateBatch(id int, newBatch *models.Batch) (*models.Batch, error) {
	var updatedBatch *models.Batch
	var previousBatchCopy *models.Batch

	// Usar transacción para mantener consistencia entre lote, insumos e inventario
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var batch models.Batch
		if err := tx.First(&batch, id).Error; err != nil {
			return err
		}

		previousBatch := batch
		amountChanged := batch.Amount != newBatch.Amount

		// Actualizar campos del lote (sin tocar id/qr_code/timestamps)
		if err := tx.Model(&batch).Omit("id", "qr_code", "created_at", "updated_at").Updates(newBatch).Error; err != nil {
			return err
		}

		// Si cambió la cantidad del lote, ajustar inventario relacionado
		if amountChanged {
			diff := batch.Amount - previousBatch.Amount

			// Si la cantidad aumenta, crear nuevos insumos individuales y mantener resúmenes coherentes
			if diff > 0 && s.MedicalSupplyService != nil {
				// Crear insumos individuales adicionales en la misma bodega
				newSupplies, err := s.MedicalSupplyService.CreateMultipleIndividualSuppliesTx(
					tx,
					batch.ID,
					diff,
					batch.StoreID,
				)
				if err != nil {
					return fmt.Errorf("error creando insumos individuales adicionales para lote %d: %v", batch.ID, err)
				}

				// Actualizar resumen de inventario de bodega usando la transacción
				if err := s.updateStoreSummaryOnAmountChangeTx(tx, &batch, &previousBatch); err != nil {
					return err
				}

				// Crear historial para los nuevos insumos de forma asíncrona (fuera de la transacción)
				go s.createHistoryAsync(batch.ID, newSupplies)
			} else {
				// Si la cantidad disminuye o no tenemos servicio de insumos,
				// solo ajustar el resumen de bodega para mantener contadores coherentes.
				if err := s.updateStoreSummaryOnAmountChangeTx(tx, &batch, &previousBatch); err != nil {
					return err
				}
			}
		}

		updatedBatch = &batch
		// Guardar copia previa para historial fuera de la transacción
		tmp := previousBatch
		previousBatchCopy = &tmp
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Registrar historial de cambios de lote
	if s.BatchHistoryService != nil && previousBatchCopy != nil && updatedBatch != nil {
		if err := s.BatchHistoryService.RegisterBatchUpdate(updatedBatch.ID, DefaultUserRUT, previousBatchCopy, updatedBatch); err != nil {
			log.Printf("Error registrando historial: %v", err)
		}
	}

	// Verificar y enviar alertas según nueva configuración del lote
	if updatedBatch != nil {
		s.checkAndSendAlerts(updatedBatch)
	}

	return updatedBatch, nil
}

func (s *BatchService) UpdateBatchAmount(id int, newAmount int) error {
	return s.DB.Model(&models.Batch{}).Where("id = ?", id).Update("amount", newAmount).Error
}

func (s *BatchService) DeleteBatch(id int) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var medicalSupplyIDs []int
		if err := tx.Model(&models.MedicalSupply{}).
			Where("batch_id = ?", id).
			Pluck("id", &medicalSupplyIDs).Error; err != nil {
			return fmt.Errorf("error obteniendo IDs de insumos: %v", err)
		}

		if len(medicalSupplyIDs) > 0 {
			if err := tx.Where("supply_request_qr_assignment_id IN (SELECT id FROM supply_request_qr_assignment WHERE medical_supply_id IN (?))",
				medicalSupplyIDs).Delete(&models.SupplyCartItem{}).Error; err != nil {
				return fmt.Errorf("error eliminando items de carrito: %v", err)
			}

			if err := tx.Where("medical_supply_id IN (?)", medicalSupplyIDs).
				Delete(&models.SupplyRequestQRAssignment{}).Error; err != nil {
				return fmt.Errorf("error eliminando asignaciones QR: %v", err)
			}

			if err := tx.Where("medical_supply_id IN (?)", medicalSupplyIDs).
				Delete(&models.SupplyTransfer{}).Error; err != nil {
				return fmt.Errorf("error eliminando transferencias: %v", err)
			}
		}

		if err := tx.Where("batch_id = ?", id).Delete(&models.StoreInventorySummary{}).Error; err != nil {
			return fmt.Errorf("error eliminando resumen de bodega: %v", err)
		}

		if err := tx.Where("batch_id = ?", id).Delete(&models.PavilionInventorySummary{}).Error; err != nil {
			return fmt.Errorf("error eliminando resumen de pabellón: %v", err)
		}

		if err := tx.Where("batch_id = ?", id).Delete(&models.QRScanEvent{}).Error; err != nil {
			return fmt.Errorf("error eliminando eventos QR: %v", err)
		}

		if err := tx.Where("batch_id = ?", id).Delete(&models.MedicalSupply{}).Error; err != nil {
			return fmt.Errorf("error eliminando insumos médicos: %v", err)
		}

		return tx.Delete(&models.Batch{}, id).Error
	})
}

// ========================
// CREACIÓN AVANZADA
// ========================

func (s *BatchService) CreateBatchWithIndividualSupplies(batch *models.Batch, supplyCode *models.SupplyCode, individualCount int, expirationAlertDays int) (*models.Batch, []models.MedicalSupply, error) {
	var individualSupplies []models.MedicalSupply

	if expirationAlertDays <= 0 {
		expirationAlertDays = DefaultExpirationAlertDays
	}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		batch.ID = 0
		batch.QRCode = ""
		if batch.ExpirationAlertDays <= 0 {
			batch.ExpirationAlertDays = expirationAlertDays
		}

		// La ubicación inicial siempre es la bodega de origen
		if batch.LocationID == 0 {
			batch.LocationType = models.BatchLocationStore
			batch.LocationID = batch.StoreID
		}

		// Resolver nombre de proveedor a ID ANTES de crear el lote (requerido por la FK)
		if batch.Supplier != "" && batch.SupplierID == 0 {
			supplierID, err := s.ensureSupplierConfigTx(tx, batch.Supplier)
			if err != nil {
				return fmt.Errorf("error resolviendo proveedor '%s': %v", batch.Supplier, err)
			}
			batch.SupplierID = supplierID
		}

		// Upsert supply_code ANTES de crear el lote (requerido por la FK batch.supply_code → supply_code.code)
		if err := s.upsertSupplyCodeTx(tx, supplyCode); err != nil {
			return err
		}
		batch.SupplyCode = supplyCode.Code

		if err := tx.Create(batch).Error; err != nil {
			return fmt.Errorf("error creando lote: %v", err)
		}

		qrCode, err := s.QRService.generateUniqueQRCode("BATCH", batch.ID)
		if err != nil {
			return fmt.Errorf("error generando QR del lote: %v", err)
		}

		if err := tx.Model(batch).Update("qr_code", qrCode).Error; err != nil {
			return fmt.Errorf("error actualizando QR del lote: %v", err)
		}
		batch.QRCode = qrCode

		supplies, err := s.createIndividualSuppliesTx(tx, batch, supplyCode, individualCount)
		if err != nil {
			return err
		}
		individualSupplies = supplies

		return s.createStoreSummaryTx(tx, batch, supplyCode)
	})

	if err != nil {
		return nil, nil, err
	}

	go s.createHistoryAsync(batch.ID, individualSupplies)
	return batch, individualSupplies, nil
}

// ========================
// HELPERS PRIVADOS
// ========================

func (s *BatchService) ensureSupplierConfig(supplierName string) (int, error) {
	if supplierName == "" {
		return 0, fmt.Errorf("nombre de proveedor vacío")
	}

	var config models.SupplierConfig
	if err := s.DB.Where("supplier_name = ?", supplierName).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			config = models.SupplierConfig{
				SupplierName: supplierName,
				Notes:        "Auto-creado automáticamente",
			}
			if err := s.DB.Create(&config).Error; err != nil {
				log.Printf("Error creando config de proveedor: %v", err)
				return 0, err
			}
		} else {
			return 0, err
		}
	}
	return config.ID, nil
}

func (s *BatchService) ensureSupplierConfigTx(tx *gorm.DB, supplierName string) (int, error) {
	if supplierName == "" {
		return 0, fmt.Errorf("nombre de proveedor vacío")
	}

	// Limpiar espacios en blanco del nombre del proveedor
	supplierName = strings.TrimSpace(supplierName)
	if supplierName == "" {
		return 0, fmt.Errorf("nombre de proveedor vacío")
	}

	var config models.SupplierConfig
	if err := tx.Where("supplier_name = ?", supplierName).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			config = models.SupplierConfig{
				SupplierName: supplierName,
				Notes:        "Auto-creado automáticamente",
			}
			if err := tx.Create(&config).Error; err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}
	return config.ID, nil
}

func (s *BatchService) upsertSupplyCodeTx(tx *gorm.DB, supplyCode *models.SupplyCode) error {
	var existing models.SupplyCode
	if err := tx.Where("code = ?", supplyCode.Code).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Create(supplyCode).Error
		}
		return err
	}

	existing.Name = supplyCode.Name
	existing.CodeSupplier = supplyCode.CodeSupplier
	existing.CriticalStock = supplyCode.CriticalStock
	return tx.Save(&existing).Error
}

func (s *BatchService) createIndividualSuppliesTx(tx *gorm.DB, batch *models.Batch, supplyCode *models.SupplyCode, count int) ([]models.MedicalSupply, error) {
	if s.MedicalSupplyService != nil {
		return s.MedicalSupplyService.CreateMultipleIndividualSuppliesTx(tx, batch.ID, count, batch.StoreID)
	}

	supplies := make([]models.MedicalSupply, 0, count)
	for i := 0; i < count; i++ {
		qrCode, err := s.QRService.generateUniqueQRCode("SUPPLY", 0)
		if err != nil {
			return nil, fmt.Errorf("error generando QR insumo %d: %v", i+1, err)
		}

		supply := models.MedicalSupply{
			QRCode:       qrCode,
			BatchID:      batch.ID,
			LocationType: models.SupplyLocationStore,
			LocationID:   batch.StoreID,
			Status:       models.StatusAvailable,
		}

		if err := tx.Create(&supply).Error; err != nil {
			return nil, fmt.Errorf("error creando insumo %d: %v", i+1, err)
		}
		supplies = append(supplies, supply)
	}
	return supplies, nil
}

func (s *BatchService) createStoreSummaryTx(tx *gorm.DB, batch *models.Batch, supplyCode *models.SupplyCode) error {
	var existing models.StoreInventorySummary
	if err := tx.Where("store_id = ? AND batch_id = ?", batch.StoreID, batch.ID).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			summary := models.StoreInventorySummary{
				StoreID:              batch.StoreID,
				BatchID:              batch.ID,
				SupplyCode:           supplyCode.Code,
				SurgeryID:            batch.SurgeryID,
				OriginalAmount:       batch.Amount,
				CurrentInStore:       batch.Amount,
				TotalTransferredOut:  0,
				TotalReturnedIn:      0,
				TotalConsumedInStore: 0,
			}
			return tx.Create(&summary).Error
		}
		return err
	}

	existing.CurrentInStore = batch.Amount
	existing.OriginalAmount = batch.Amount
	return tx.Save(&existing).Error
}

// updateStoreSummaryOnAmountChangeTx ajusta el resumen de inventario de bodega cuando cambia la cantidad del lote
func (s *BatchService) updateStoreSummaryOnAmountChangeTx(tx *gorm.DB, newBatch, previousBatch *models.Batch) error {
	var summary models.StoreInventorySummary
	if err := tx.Where("batch_id = ?", newBatch.ID).First(&summary).Error; err != nil {
		// Si no existe resumen, no forzamos su creación aquí para no alterar otros flujos
		return nil
	}

	diff := newBatch.Amount - previousBatch.Amount
	summary.OriginalAmount = newBatch.Amount
	summary.CurrentInStore += diff

	if summary.CurrentInStore < 0 {
		summary.CurrentInStore = 0
	}

	if err := tx.Save(&summary).Error; err != nil {
		log.Printf("Error actualizando resumen de bodega (tx): %v", err)
		return err
	}
	return nil
}

func (s *BatchService) createHistoryAsync(batchID int, supplies []models.MedicalSupply) {
	if s.BatchHistoryService != nil {
		if err := s.BatchHistoryService.RegisterBatchCreation(batchID, DefaultUserRUT); err != nil {
			log.Printf("Error creando historial lote %d: %v", batchID, err)
		}
	}

	for _, supply := range supplies {
		history := models.SupplyHistory{
			DateTime:        time.Now(),
			Status:          "creado",
			DestinationType: models.DestinationTypeStore,
			DestinationID:   supply.LocationID,
			MedicalSupplyID: supply.ID,
			UserRUT:         DefaultUserRUT,
		}
		if err := s.DB.Create(&history).Error; err != nil {
			log.Printf("Error creando historial insumo %d: %v", supply.ID, err)
		}
	}
}

func (s *BatchService) updateStoreSummaryOnAmountChange(newBatch, previousBatch *models.Batch) {
	var summary models.StoreInventorySummary
	if err := s.DB.Where("batch_id = ?", newBatch.ID).First(&summary).Error; err != nil {
		return
	}

	diff := newBatch.Amount - previousBatch.Amount
	summary.OriginalAmount = newBatch.Amount
	summary.CurrentInStore += diff

	if summary.CurrentInStore < 0 {
		summary.CurrentInStore = 0
	}

	if err := s.DB.Save(&summary).Error; err != nil {
		log.Printf("Error actualizando resumen de bodega: %v", err)
	}
}

// getSupplyCodeByBatchID obtiene el supply code asociado a un lote
func (s *BatchService) getSupplyCodeByBatchID(batchID int) (*models.SupplyCode, error) {
	var batch models.Batch
	if err := s.DB.Select("id, supply_code").First(&batch, batchID).Error; err != nil {
		return nil, err
	}
	var supplyCode models.SupplyCode
	if err := s.DB.First(&supplyCode, batch.SupplyCode).Error; err != nil {
		return nil, err
	}
	return &supplyCode, nil
}

func (s *BatchService) checkAndSendAlerts(batch *models.Batch) {
	supplyCode, err := s.getSupplyCodeByBatchID(batch.ID)
	if err != nil {
		log.Printf("Error obteniendo supply code: %v", err)
		return
	}

	if batch.Amount > 0 && batch.Amount <= supplyCode.CriticalStock {
		if err := s.sendAlert(*batch, *supplyCode, "low_stock"); err != nil {
			log.Printf("Error enviando alerta stock bajo: %v", err)
		}
	}

	alertDays := batch.ExpirationAlertDays
	if alertDays <= 0 {
		alertDays = DefaultExpirationAlertDays
	}
	expirationThreshold := time.Now().AddDate(0, 0, alertDays)
	if batch.ExpirationDate.Before(expirationThreshold) && batch.ExpirationDate.After(time.Now()) {
		if err := s.sendAlert(*batch, *supplyCode, "expiration"); err != nil {
			log.Printf("Error enviando alerta vencimiento: %v", err)
		}
	}
}

func (s *BatchService) sendAlert(batch models.Batch, supplyCode models.SupplyCode, alertType string) error {
	var templatePath string
	var subject string
	data := map[string]interface{}{
		"BatchID": batch.ID,
		"Code":    supplyCode.Code,
		"Name":    supplyCode.Name,
		"Date":    time.Now().Format("02/01/2006"),
	}

	switch alertType {
	case "low_stock":
		templatePath = "mailer/templates/low_stock.html"
		subject = "Alerta: Stock Bajo en Lote"
		data["CurrentStock"] = batch.Amount
		data["CriticalStock"] = supplyCode.CriticalStock
	case "expiration":
		templatePath = "mailer/templates/expiration_warning.html"
		subject = "Alerta: Lote Próximo a Vencer"
		data["ExpirationDate"] = batch.ExpirationDate.Format("02/01/2006")
		data["DaysUntilExp"] = int(time.Until(batch.ExpirationDate).Hours() / 24)
	default:
		return fmt.Errorf("tipo de alerta desconocido: %s", alertType)
	}

	alertEmail := clinicconfig.GetAlertEmail(s.DB)
	if alertEmail == "" {
		return fmt.Errorf("email de alertas no configurado en la BD ni en ALERT_EMAIL")
	}
	req := mailer.NewRequest([]string{alertEmail}, subject)
	return req.SendMailSkipTLS(templatePath, data)
}

// RESTO DE MÉTODOS PÚBLICOS (GetBatchQRInfo, SyncAllBatchAmounts, etc.) SIN CAMBIOS...
// [Continúa con los demás métodos que no requieren optimización]

func (s *BatchService) GetBatchQRInfo(qrCode string) (map[string]interface{}, error) {
	var batch models.Batch
	if err := s.DB.Where("qr_code = ?", qrCode).First(&batch).Error; err != nil {
		return nil, fmt.Errorf("lote no encontrado con QR %s: %v", qrCode, err)
	}
	return s.GetBatchWithSupplyInfo(batch.ID)
}

func (s *BatchService) SyncAllBatchAmounts() error {
	if s.MedicalSupplyService != nil {
		return s.MedicalSupplyService.SyncBatchAmounts()
	}
	// Implementación fallback...
	return nil
}

func (s *BatchService) CheckLowStockAlert(batchID int, threshold int) error {
	var batch models.Batch
	if err := s.DB.First(&batch, batchID).Error; err != nil {
		return fmt.Errorf("lote no encontrado: %v", err)
	}

	supplyCode, err := s.getSupplyCodeByBatchID(batchID)
	if err != nil {
		return fmt.Errorf("error obteniendo supply code: %v", err)
	}

	if batch.Amount > 0 && batch.Amount <= supplyCode.CriticalStock {
		return s.sendAlert(batch, *supplyCode, "low_stock")
	}
	return nil
}

func (s *BatchService) CheckExpirationAlert(batchID int, daysThreshold int) error {
	var batch models.Batch
	if err := s.DB.First(&batch, batchID).Error; err != nil {
		return fmt.Errorf("lote no encontrado: %v", err)
	}

	if daysThreshold <= 0 {
		daysThreshold = batch.ExpirationAlertDays
		if daysThreshold <= 0 {
			daysThreshold = DefaultExpirationAlertDays
		}
	}

	daysUntilExpiration := int(time.Until(batch.ExpirationDate).Hours() / 24)
	if daysUntilExpiration <= daysThreshold && daysUntilExpiration > 0 {
		supplyCode, err := s.getSupplyCodeByBatchID(batchID)
		if err != nil {
			return err
		}
		return s.sendAlert(batch, *supplyCode, "expiration")
	}
	return nil
}

// ========================
// SINCRONIZACIÓN Y VERIFICACIÓN
// ========================

func (s *BatchService) GetBatchesNeedingSync() ([]map[string]interface{}, error) {
	query := `
		SELECT 
			b.id,
			b.qr_code,
			b.amount as batch_amount,
			sc.supplier_name as supplier,
			b.expiration_date,
			COUNT(ms.id) as total_supplies,
			COUNT(CASE WHEN consumed.supply_id IS NOT NULL THEN 1 END) as consumed_supplies,
			(COUNT(ms.id) - COUNT(CASE WHEN consumed.supply_id IS NOT NULL THEN 1 END)) as available_supplies
		FROM batch b
		LEFT JOIN supplier_config sc ON b.supplier_id = sc.id
		LEFT JOIN medical_supply ms ON b.id = ms.batch_id
		LEFT JOIN (
			SELECT DISTINCT sh.medical_supply_id as supply_id
			FROM supply_history sh
			WHERE sh.status = 'consumido'
		) consumed ON ms.id = consumed.supply_id
		GROUP BY b.id, b.qr_code, b.amount, sc.supplier_name, b.expiration_date
		HAVING b.amount != (COUNT(ms.id) - COUNT(CASE WHEN consumed.supply_id IS NOT NULL THEN 1 END))
		ORDER BY b.id
	`

	rows, err := s.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var batchID int
		var qrCode string
		var batchAmount int
		var supplier string
		var expirationDate time.Time
		var totalSupplies, consumedSupplies, availableSupplies int

		err := rows.Scan(&batchID, &qrCode, &batchAmount, &supplier, &expirationDate, &totalSupplies, &consumedSupplies, &availableSupplies)
		if err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"batch_id":           batchID,
			"qr_code":            qrCode,
			"batch_amount":       batchAmount,
			"available_supplies": availableSupplies,
			"total_supplies":     totalSupplies,
			"consumed_supplies":  consumedSupplies,
			"supplier":           supplier,
			"expiration_date":    expirationDate,
			"difference":         batchAmount - availableSupplies,
		})
	}

	return results, nil
}

func (s *BatchService) CheckLowStockForSupplyType(supplyCode string) error {
	var sc models.SupplyCode
	if err := s.DB.Where("code = ?", supplyCode).First(&sc).Error; err != nil {
		return fmt.Errorf("código de insumo no encontrado: %v", err)
	}

	query := `
		SELECT DISTINCT b.* 
		FROM batch b
		WHERE b.supply_code = ? AND b.amount > 0 AND b.amount <= ?
	`

	var batches []models.Batch
	if err := s.DB.Raw(query, supplyCode, sc.CriticalStock).Scan(&batches).Error; err != nil {
		return fmt.Errorf("error obteniendo lotes con stock bajo: %v", err)
	}

	alertsSent := 0
	for _, batch := range batches {
		var supplyCodeData models.SupplyCode
		if err := s.DB.Where("code = ?", supplyCode).First(&supplyCodeData).Error; err != nil {
			log.Printf("Error obteniendo supply code: %v", err)
			continue
		}

		if err := s.sendAlert(batch, supplyCodeData, "low_stock"); err != nil {
			log.Printf("Error enviando alerta de stock bajo para lote %d: %v", batch.ID, err)
			continue
		}
		alertsSent++
	}

	if alertsSent > 0 {
		log.Printf("Se enviaron %d alertas de stock bajo para el insumo %s", alertsSent, supplyCode)
	}

	return nil
}

func (s *BatchService) CheckAllBatchesLowStock() error {
	query := `
		SELECT DISTINCT b.*, sc.critical_stock, sc.name as supply_name
		FROM batch b
		JOIN supply_code sc ON b.supply_code = sc.code
		WHERE b.amount > 0 AND b.amount <= sc.critical_stock
		ORDER BY b.id
	`

	type BatchWithCriticalStock struct {
		models.Batch
		CriticalStock int    `json:"critical_stock"`
		SupplyName    string `json:"supply_name"`
	}

	var batchesWithCriticalStock []BatchWithCriticalStock
	if err := s.DB.Raw(query).Scan(&batchesWithCriticalStock).Error; err != nil {
		return fmt.Errorf("error obteniendo lotes con stock crítico: %v", err)
	}

	alertsSent := 0
	errors := 0

	for _, batchInfo := range batchesWithCriticalStock {
		supplyCode, err := s.getSupplyCodeByBatchID(batchInfo.ID)
		if err != nil {
			log.Printf("Error obteniendo supply code para lote %d: %v", batchInfo.ID, err)
			errors++
			continue
		}

		if err := s.sendAlert(batchInfo.Batch, *supplyCode, "low_stock"); err != nil {
			log.Printf("Error enviando alerta de stock bajo para lote %d (%s): %v",
				batchInfo.ID, batchInfo.SupplyName, err)
			errors++
			continue
		}
		alertsSent++
		log.Printf("Alerta enviada para lote %d (%s): Stock actual %d, Stock crítico %d",
			batchInfo.ID, batchInfo.SupplyName, batchInfo.Amount, batchInfo.CriticalStock)
	}

	log.Printf("Verificación de stock bajo completada: %d alertas enviadas, %d errores", alertsSent, errors)
	return nil
}

func (s *BatchService) StartAutomaticLowStockChecker() {
	ticker := time.NewTicker(DefaultLowStockCheckInterval)
	defer ticker.Stop()

	log.Println("🔄 Iniciado verificador automático de stock bajo")

	for range ticker.C {
		log.Println("🔍 Ejecutando verificación automática de stock bajo...")
		if err := s.CheckAllBatchesLowStock(); err != nil {
			log.Printf("Error en verificación automática de stock bajo: %v\n", err)
		}
	}
}

func (s *BatchService) GetLowStockSummary() ([]map[string]interface{}, error) {
	query := `
		SELECT 
			sc.code,
			sc.name,
			sc.critical_stock,
			SUM(b.amount) as total_current_stock,
			COUNT(DISTINCT b.id) as batch_count,
			MIN(b.expiration_date) as nearest_expiration,
			ARRAY_AGG(DISTINCT supc.supplier_name) as suppliers
		FROM supply_code sc
		JOIN batch b ON sc.code = b.supply_code
		JOIN supplier_config supc ON b.supplier_id = supc.id
		WHERE b.amount > 0
		GROUP BY sc.code, sc.name, sc.critical_stock
		HAVING SUM(b.amount) <= sc.critical_stock
		ORDER BY (SUM(b.amount)::float / sc.critical_stock::float) ASC, sc.name
	`

	rows, err := s.DB.Raw(query).Rows()
	if err != nil {
		return nil, fmt.Errorf("error ejecutando consulta de resumen de stock bajo: %v", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var code int
		var name string
		var criticalStock int
		var totalCurrentStock int
		var batchCount int
		var nearestExpiration time.Time
		var suppliers string

		err := rows.Scan(&code, &name, &criticalStock, &totalCurrentStock, &batchCount, &nearestExpiration, &suppliers)
		if err != nil {
			return nil, fmt.Errorf("error escaneando resultado: %v", err)
		}

		stockPercentage := float64(totalCurrentStock) / float64(criticalStock) * 100
		isUrgent := totalCurrentStock <= criticalStock/2

		results = append(results, map[string]interface{}{
			"code":                code,
			"name":                name,
			"critical_stock":      criticalStock,
			"total_current_stock": totalCurrentStock,
			"batch_count":         batchCount,
			"nearest_expiration":  nearestExpiration,
			"suppliers":           suppliers,
			"stock_percentage":    stockPercentage,
			"is_urgent":           isUrgent,
			"deficit":             criticalStock - totalCurrentStock,
		})
	}

	return results, nil
}

func (s *BatchService) CheckAllBatchesExpiration() error {
	var activeBatches []models.Batch
	if err := s.DB.Where("expiration_date > ?", time.Now()).Find(&activeBatches).Error; err != nil {
		return fmt.Errorf("error obteniendo lotes activos para verificación: %v", err)
	}

	alertsSent := 0
	errors := 0

	for _, batch := range activeBatches {
		if err := s.CheckExpirationAlert(batch.ID, 0); err != nil {
			log.Printf("Error verificando alerta de vencimiento para lote %d: %v", batch.ID, err)
			errors++
			continue
		}
		alertsSent++
	}

	log.Printf("Verificación de vencimiento completada: %d verificaciones, %d errores", alertsSent, errors)
	return nil
}
