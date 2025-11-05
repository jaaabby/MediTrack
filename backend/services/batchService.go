package services

import (
	"errors"
	"fmt"
	"log"
	"meditrack/mailer"
	"meditrack/models"
	"time"

	"gorm.io/gorm"
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

// SetMedicalSupplyService establece el servicio de suministros médicos
func (s *BatchService) SetMedicalSupplyService(medicalSupplyService *MedicalSupplyService) {
	s.MedicalSupplyService = medicalSupplyService
}

// SetBatchHistoryService establece el servicio de historial de lotes
func (s *BatchService) SetBatchHistoryService(batchHistoryService *BatchHistoryService) {
	s.BatchHistoryService = batchHistoryService
}

// CreateBatch crea un nuevo lote con QR único
func (s *BatchService) CreateBatch(batch *models.Batch) error {

	batch.ID = 0
	batch.QRCode = ""

	// Crear el lote
	if err := s.DB.Create(batch).Error; err != nil {
		return fmt.Errorf("error creating batch: %v", err)
	}

	// Generar QR después de obtener ID
	batch.QRCode = fmt.Sprintf("BATCH_%d", batch.ID)

	// Actualizar con QR
	if err := s.DB.Model(batch).Update("qr_code", batch.QRCode).Error; err != nil {
		// No fallar si el QR no se puede actualizar
		log.Printf("Warning: Could not update QR code for batch %d: %v", batch.ID, err)
	}

	// Crear configuración de proveedor si no existe (con valor por defecto de 90 días)
	if batch.Supplier != "" {
		s.ensureSupplierConfig(batch.Supplier, batch.ID)
	}

	return nil
}

// ensureSupplierConfig verifica y crea configuración de proveedor si no existe (con valor por defecto de 90 días)
func (s *BatchService) ensureSupplierConfig(supplierName string, batchID int) {
	s.ensureSupplierConfigWithDays(supplierName, batchID, 90)
}

// ensureSupplierConfigWithDays verifica y crea configuración de proveedor si no existe con días de alerta específicos
func (s *BatchService) ensureSupplierConfigWithDays(supplierName string, batchID int, expirationAlertDays int) {
	if supplierName == "" {
		return
	}

	// Si expirationAlertDays es 0 o menor, usar 90 por defecto
	if expirationAlertDays <= 0 {
		expirationAlertDays = 90
	}

	var existingSupplierConfig models.SupplierConfig
	if err := s.DB.Where("supplier_name = ?", supplierName).First(&existingSupplierConfig).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No existe configuración para este proveedor, crear una con el valor proporcionado
			supplierConfig := models.SupplierConfig{
				SupplierName:       supplierName,
				ExpirationAlertDays: expirationAlertDays,
				Notes:              fmt.Sprintf("Configuración creada automáticamente al crear lote %d", batchID),
			}
			if err := s.DB.Create(&supplierConfig).Error; err != nil {
				// No fallar si no se puede crear la configuración, solo log
				log.Printf("Advertencia: No se pudo crear configuración de proveedor para '%s': %v", supplierName, err)
			} else {
				log.Printf("✅ Configuración de proveedor creada automáticamente para '%s' con %d días de alerta", supplierName, supplierConfig.ExpirationAlertDays)
			}
		} else {
			// Error diferente, solo log
			log.Printf("Advertencia: Error verificando configuración de proveedor para '%s': %v", supplierName, err)
		}
	}
	// Si ya existe, no hacer nada (mantener configuración existente)
}

// CreateBatchWithIndividualSupplies crea un lote junto con sus insumos individuales
// expirationAlertDays: días de alerta para crear configuración del proveedor (si no existe). Si es 0, usa 90 por defecto.
func (s *BatchService) CreateBatchWithIndividualSupplies(batch *models.Batch, supplyCode *models.SupplyCode, individualCount int, expirationAlertDays int) (*models.Batch, []models.MedicalSupply, error) {
	var individualSupplies []models.MedicalSupply

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Asegurar que el ID sea 0 para permitir auto-incremento
		batch.ID = 0
		batch.QRCode = ""

		// 1. Crear el lote primero sin QR
		if err := tx.Create(batch).Error; err != nil {
			return fmt.Errorf("error creando lote: %v", err)
		}

		// 2. Generar QR único para el lote después de tener el ID
		qrCode, err := s.QRService.generateUniqueQRCode("BATCH", batch.ID)
		if err != nil {
			return fmt.Errorf("error generando QR del lote: %v", err)
		}

		// 3. Actualizar el lote con el QR generado
		if err := tx.Model(batch).Update("qr_code", qrCode).Error; err != nil {
			return fmt.Errorf("error actualizando QR del lote: %v", err)
		}
		batch.QRCode = qrCode

		// 4. Crear o actualizar supply_code
		var existingSupplyCode models.SupplyCode
		if err := tx.Where("code = ?", supplyCode.Code).First(&existingSupplyCode).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// No existe, crear nuevo
				if err := tx.Create(supplyCode).Error; err != nil {
					return fmt.Errorf("error creando código de insumo: %v", err)
				}
			} else {
				return fmt.Errorf("error buscando código de insumo: %v", err)
			}
		} else {
			// Existe, actualizar datos si es necesario
			existingSupplyCode.Name = supplyCode.Name
			existingSupplyCode.CodeSupplier = supplyCode.CodeSupplier
			existingSupplyCode.CriticalStock = supplyCode.CriticalStock
			if err := tx.Save(&existingSupplyCode).Error; err != nil {
				return fmt.Errorf("error actualizando código de insumo: %v", err)
			}
		}

		// 4.5. Crear configuración de proveedor si no existe
		var existingSupplierConfig models.SupplierConfig
		if err := tx.Where("supplier_name = ?", batch.Supplier).First(&existingSupplierConfig).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// No existe configuración para este proveedor, crear una con el valor proporcionado
				// Si expirationAlertDays es 0 o menor, usar 90 por defecto
				alertDays := expirationAlertDays
				if alertDays <= 0 {
					alertDays = 90 // Valor por defecto: 90 días (3 meses)
				}
				supplierConfig := models.SupplierConfig{
					SupplierName:       batch.Supplier,
					ExpirationAlertDays: alertDays,
					Notes:              fmt.Sprintf("Configuración creada automáticamente al crear lote %d", batch.ID),
				}
				if err := tx.Create(&supplierConfig).Error; err != nil {
					// No fallar si no se puede crear la configuración, solo log
					log.Printf("Advertencia: No se pudo crear configuración de proveedor para '%s': %v", batch.Supplier, err)
				} else {
					log.Printf("✅ Configuración de proveedor creada automáticamente para '%s' con %d días de alerta", batch.Supplier, supplierConfig.ExpirationAlertDays)
				}
			} else {
				// Error diferente, solo log
				log.Printf("Advertencia: Error verificando configuración de proveedor para '%s': %v", batch.Supplier, err)
			}
		}
		// Si ya existe, no hacer nada (mantener configuración existente)

		// 5. Crear insumos individuales
		if s.MedicalSupplyService != nil {
			supplies, err := s.MedicalSupplyService.CreateMultipleIndividualSuppliesTx(
				tx,
				batch.ID,
				supplyCode.Code,
				individualCount,
				batch.StoreID, // Pasar el storeID para establecer LocationID
			)
			if err != nil {
				return fmt.Errorf("error creando insumos individuales: %v", err)
			}
			individualSupplies = supplies
		} else {
			// Fallback: crear manualmente si el servicio no está disponible
			for i := 0; i < individualCount; i++ {
				qrCode, err := s.QRService.generateUniqueQRCode("SUPPLY", 0)
				if err != nil {
					return fmt.Errorf("error generando QR para insumo %d: %v", i+1, err)
				}

				supply := models.MedicalSupply{
					Code:         supplyCode.Code,
					QRCode:       qrCode,
					BatchID:      batch.ID,
					LocationType: models.SupplyLocationStore,
					LocationID:   batch.StoreID,
					Status:       models.StatusAvailable,
				}

				if err := tx.Create(&supply).Error; err != nil {
					return fmt.Errorf("error creando insumo individual %d: %v", i+1, err)
				}

				individualSupplies = append(individualSupplies, supply)
			}
		}

		// 5.5. Crear entrada en store_inventory_summary para el nuevo lote
		var existingStoreSummary models.StoreInventorySummary
		if err := tx.Where("store_id = ? AND batch_id = ?", batch.StoreID, batch.ID).First(&existingStoreSummary).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Crear resumen de inventario de bodega para el nuevo lote
				storeSummary := models.StoreInventorySummary{
					StoreID:            batch.StoreID,
					BatchID:            batch.ID,
					SupplyCode:         supplyCode.Code,
					SurgeryID:          batch.SurgeryID,
					OriginalAmount:     batch.Amount,
					CurrentInStore:     batch.Amount, // Todos los insumos están inicialmente en bodega
					TotalTransferredOut: 0,
					TotalReturnedIn:     0,
					TotalConsumedInStore: 0,
				}
				if err := tx.Create(&storeSummary).Error; err != nil {
					return fmt.Errorf("error creando resumen de inventario de bodega: %v", err)
				}
				log.Printf("✅ Resumen de inventario de bodega creado: StoreID=%d, BatchID=%d, CurrentInStore=%d", 
					batch.StoreID, batch.ID, batch.Amount)
			} else {
				return fmt.Errorf("error verificando resumen de inventario de bodega: %v", err)
			}
		} else {
			// Ya existe, actualizar los valores iniciales
			existingStoreSummary.CurrentInStore = batch.Amount
			existingStoreSummary.OriginalAmount = batch.Amount
			if err := tx.Save(&existingStoreSummary).Error; err != nil {
				return fmt.Errorf("error actualizando resumen de inventario de bodega: %v", err)
			}
			log.Printf("✅ Resumen de inventario de bodega actualizado: StoreID=%d, BatchID=%d, CurrentInStore=%d", 
				batch.StoreID, batch.ID, batch.Amount)
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	// 5. CREAR HISTORIAL DE LOTE (después de la transacción exitosa)
	if s.BatchHistoryService != nil {
		go func() {
			// Ejecutar en goroutine para no bloquear la respuesta
			if err := s.BatchHistoryService.RegisterBatchCreation(batch.ID, "12345678-9"); err != nil {
				// Log del error pero no fallar toda la operación
				fmt.Printf("Advertencia: Error creando historial de lote %d: %v\n", batch.ID, err)
			}
		}()
	}

	// 5.5. Asegurar que existe configuración del proveedor (después de la transacción exitosa)
	// Nota: Esto ya se hizo dentro de la transacción, pero lo mantenemos aquí por si acaso
	// Si expirationAlertDays es 0 o menor, usar 90 por defecto
	alertDays := expirationAlertDays
	if alertDays <= 0 {
		alertDays = 90
	}
	if batch.Supplier != "" {
		s.ensureSupplierConfigWithDays(batch.Supplier, batch.ID, alertDays)
	}

	// 6. CREAR HISTORIAL INICIAL PARA CADA INSUMO (opcional)
	if len(individualSupplies) > 0 {
		go func() {
			// Ejecutar en goroutine para no bloquear
			for _, supply := range individualSupplies {
				// Crear historial inicial de "creado" para cada insumo
				history := models.SupplyHistory{
					DateTime:        time.Now(),
					Status:          "creado",
					DestinationType: "almacen",
					DestinationID:   batch.StoreID,
					MedicalSupplyID: supply.ID,
					UserRUT:         "12345678-9",
				}

				if err := s.DB.Create(&history).Error; err != nil {
					fmt.Printf("Advertencia: Error creando historial inicial para insumo %d: %v\n", supply.ID, err)
				}
			}
		}()
	}

	return batch, individualSupplies, nil
}

func (s *BatchService) GetBatchByID(id int) (*models.Batch, error) {
	var batch models.Batch
	if err := s.DB.First(&batch, id).Error; err != nil {
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

// GetBatchesWithFilters obtiene lotes con filtros opcionales
func (s *BatchService) GetBatchesWithFilters(surgeryID *int, storeID *int, supplier string) ([]models.Batch, error) {
	var batches []models.Batch
	query := s.DB.Model(&models.Batch{})

	// Aplicar filtros si se proporcionan
	if surgeryID != nil {
		query = query.Where("surgery_id = ?", *surgeryID)
	}

	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}

	if supplier != "" {
		query = query.Where("supplier ILIKE ?", "%"+supplier+"%")
	}

	if err := query.Find(&batches).Error; err != nil {
		return nil, err
	}

	return batches, nil
}

// GetBatchWithSupplyInfo obtiene un lote con información de sus insumos
func (s *BatchService) GetBatchWithSupplyInfo(id int) (map[string]interface{}, error) {
	var batch models.Batch
	if err := s.DB.First(&batch, id).Error; err != nil {
		return nil, fmt.Errorf("lote no encontrado: %v", err)
	}

	// Contar insumos totales y consumidos
	var totalSupplies, consumedSupplies int64
	s.DB.Model(&models.MedicalSupply{}).Where("batch_id = ?", id).Count(&totalSupplies)

	subquery := s.DB.Model(&models.SupplyHistory{}).
		Select("medical_supply_id").
		Where("status = ?", "consumido")
	s.DB.Model(&models.MedicalSupply{}).
		Where("batch_id = ? AND id IN (?)", id, subquery).
		Count(&consumedSupplies)

	// Obtener información del supply code
	var supplyCode models.SupplyCode
	s.DB.Joins("JOIN medical_supply ON medical_supply.code = supply_code.code").
		Where("medical_supply.batch_id = ?", id).
		First(&supplyCode)

	result := map[string]interface{}{
		"batch":                     batch,
		"supply_code":               supplyCode,
		"total_individual_supplies": totalSupplies,
		"consumed_supplies":         consumedSupplies,
		"available_supplies":        totalSupplies - consumedSupplies,
		"batch_amount":              batch.Amount,
		"amounts_synchronized":      batch.Amount == int(totalSupplies-consumedSupplies),
		"consumption_percentage":    float64(consumedSupplies) / float64(totalSupplies) * 100,
	}

	return result, nil
}

func (s *BatchService) UpdateBatch(id int, newBatch *models.Batch) (*models.Batch, error) {
	var batch models.Batch
	if err := s.DB.First(&batch, id).Error; err != nil {
		return nil, err
	}

	// Guardar valores anteriores para el historial
	previousBatch := batch
	amountChanged := batch.Amount != newBatch.Amount

	batch.ExpirationDate = newBatch.ExpirationDate
	batch.Amount = newBatch.Amount
	batch.Supplier = newBatch.Supplier
	batch.StoreID = newBatch.StoreID

	if err := s.DB.Save(&batch).Error; err != nil {
		return nil, err
	}

	// Si cambió la cantidad, actualizar el resumen de inventario para mantener consistencia
	if amountChanged {
		var storeSummary models.StoreInventorySummary
		if err := s.DB.Where("batch_id = ?", batch.ID).First(&storeSummary).Error; err == nil {
			// Si existe el resumen, actualizar el OriginalAmount si es necesario
			// El CurrentInStore se mantiene según los insumos reales, pero ajustamos OriginalAmount
			// Solo si el nuevo amount es mayor que el original y hay diferencia
			if newBatch.Amount > storeSummary.OriginalAmount {
				// Si se aumentó la cantidad del lote, ajustar el original y el current
				diff := newBatch.Amount - storeSummary.OriginalAmount
				storeSummary.OriginalAmount = newBatch.Amount
				storeSummary.CurrentInStore += diff
				if err := s.DB.Save(&storeSummary).Error; err != nil {
					fmt.Printf("⚠️ Error actualizando resumen de bodega tras cambio de cantidad: %v\n", err)
				}
			} else if newBatch.Amount < storeSummary.OriginalAmount {
				// Si se disminuyó la cantidad del lote, ajustar el original
				// No ajustamos CurrentInStore porque refleja el stock real
				storeSummary.OriginalAmount = newBatch.Amount
				if storeSummary.CurrentInStore > newBatch.Amount {
					// Si el stock actual es mayor que el nuevo amount, ajustar
					storeSummary.CurrentInStore = newBatch.Amount
				}
				if err := s.DB.Save(&storeSummary).Error; err != nil {
					fmt.Printf("⚠️ Error actualizando resumen de bodega tras cambio de cantidad: %v\n", err)
				}
			}
		}
	}

	// Registrar en el historial (RUT hardcodeado por ahora)
	userRUT := "12345678-9" // RUT hardcodeado temporalmente
	if s.BatchHistoryService != nil {
		if err := s.BatchHistoryService.RegisterBatchUpdate(batch.ID, userRUT, &previousBatch, &batch); err != nil {
			// Solo log del error, no fallar la actualización
			fmt.Printf("Error al registrar en historial: %v\n", err)
		}
	} else {
		fmt.Printf("BatchHistoryService no inicializado, no se registra historial\n")
	}

	// Verificar si se debe enviar correo por stock bajo usando el stock crítico del tipo de insumo
	var supplyCode models.SupplyCode
	if err := s.DB.Joins("JOIN medical_supply ON medical_supply.code = supply_code.code").
		Where("medical_supply.batch_id = ?", batch.ID).
		First(&supplyCode).Error; err == nil {

		// Usar el critical_stock del tipo de insumo
		if batch.Amount > 0 && batch.Amount <= supplyCode.CriticalStock {
			if err := s.sendLowStockAlert(batch); err != nil {
				// Solo log del error, no fallar la actualización
				fmt.Printf("Error al enviar alerta de stock bajo: %v\n", err)
			}
		}
	} else {
		fmt.Printf("Error obteniendo información del código de insumo para verificar stock crítico: %v\n", err)
	}

	// Verificar si se debe enviar correo por vencimiento próximo (usando configuración del proveedor)
	alertDays := s.getExpirationAlertDays(batch.Supplier)
	expirationThreshold := time.Now().AddDate(0, 0, alertDays)
	if batch.ExpirationDate.Before(expirationThreshold) && batch.ExpirationDate.After(time.Now()) {
		if err := s.sendExpirationAlert(batch); err != nil {
			// Solo log del error, no fallar la actualización
			fmt.Printf("Error al enviar alerta de vencimiento: %v\n", err)
		}
	}

	return &batch, nil
}

// UpdateBatchAmount actualiza solo la cantidad del lote (usado internamente por el sistema de consumo)
func (s *BatchService) UpdateBatchAmount(id int, newAmount int) error {
	return s.DB.Model(&models.Batch{}).Where("id = ?", id).Update("amount", newAmount).Error
}

// CheckLowStockAlert verifica y envía alertas de stock bajo usando el stock crítico del tipo de insumo
func (s *BatchService) CheckLowStockAlert(batchID int, threshold int) error {
	var batch models.Batch
	if err := s.DB.First(&batch, batchID).Error; err != nil {
		return fmt.Errorf("lote no encontrado: %v", err)
	}

	// Obtener el stock crítico del tipo de insumo
	var supplyCode models.SupplyCode
	if err := s.DB.Joins("JOIN medical_supply ON medical_supply.code = supply_code.code").
		Where("medical_supply.batch_id = ?", batchID).
		First(&supplyCode).Error; err != nil {
		return fmt.Errorf("error al obtener información del código de insumo: %v", err)
	}

	// Usar el critical_stock del tipo de insumo en lugar del threshold
	if batch.Amount > 0 && batch.Amount <= supplyCode.CriticalStock {
		return s.sendLowStockAlert(batch)
	}

	return nil
}

// CheckExpirationAlert verifica y envía alertas de vencimiento próximo
// Si daysThreshold es 0, usa la configuración del proveedor
func (s *BatchService) CheckExpirationAlert(batchID int, daysThreshold int) error {
	var batch models.Batch
	if err := s.DB.First(&batch, batchID).Error; err != nil {
		return fmt.Errorf("lote no encontrado: %v", err)
	}

	// Si no se especifica threshold, usar configuración del proveedor
	if daysThreshold <= 0 {
		daysThreshold = s.getExpirationAlertDays(batch.Supplier)
	}

	daysUntilExpiration := int(time.Until(batch.ExpirationDate).Hours() / 24)
	if daysUntilExpiration <= daysThreshold && daysUntilExpiration > 0 {
		return s.sendExpirationAlert(batch)
	}

	return nil
}

func (s *BatchService) DeleteBatch(id int) error {
	// Usar transacción para eliminar de manera segura
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Eliminar registros de medical_supply (supply_history se elimina automáticamente por CASCADE)
		if err := tx.Where("batch_id = ?", id).Delete(&models.MedicalSupply{}).Error; err != nil {
			return fmt.Errorf("error al eliminar insumos médicos del lote: %v", err)
		}

		// 2. Eliminar el lote (el trigger log_batch_delete se ejecutará automáticamente)
		if err := tx.Delete(&models.Batch{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetBatchQRInfo obtiene información completa de un lote por su QR
func (s *BatchService) GetBatchQRInfo(qrCode string) (map[string]interface{}, error) {
	var batch models.Batch
	if err := s.DB.Where("qr_code = ?", qrCode).First(&batch).Error; err != nil {
		return nil, fmt.Errorf("lote no encontrado con QR %s: %v", qrCode, err)
	}

	return s.GetBatchWithSupplyInfo(batch.ID)
}

// SyncAllBatchAmounts sincroniza todas las cantidades de lotes
func (s *BatchService) SyncAllBatchAmounts() error {
	if s.MedicalSupplyService != nil {
		return s.MedicalSupplyService.SyncBatchAmounts()
	}

	// Fallback si no está disponible el servicio
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var batches []models.Batch
		if err := tx.Find(&batches).Error; err != nil {
			return fmt.Errorf("error obteniendo lotes: %v", err)
		}

		for _, batch := range batches {
			var availableCount int64
			subquery := tx.Model(&models.SupplyHistory{}).
				Select("medical_supply_id").
				Where("status = ?", "consumido")

			if err := tx.Model(&models.MedicalSupply{}).
				Where("batch_id = ? AND id NOT IN (?)", batch.ID, subquery).
				Count(&availableCount).Error; err != nil {
				return fmt.Errorf("error contando productos disponibles para lote %d: %v", batch.ID, err)
			}

			if int64(batch.Amount) != availableCount {
				if err := tx.Model(&batch).Update("amount", availableCount).Error; err != nil {
					return fmt.Errorf("error actualizando cantidad del lote %d: %v", batch.ID, err)
				}
			}
		}

		return nil
	})
}

// GetBatchesNeedingSync obtiene lotes que necesitan sincronización
func (s *BatchService) GetBatchesNeedingSync() ([]map[string]interface{}, error) {
	query := `
		SELECT 
			b.id,
			b.qr_code,
			b.amount as batch_amount,
			b.supplier,
			b.expiration_date,
			COUNT(ms.id) as total_supplies,
			COUNT(CASE WHEN consumed.supply_id IS NOT NULL THEN 1 END) as consumed_supplies,
			(COUNT(ms.id) - COUNT(CASE WHEN consumed.supply_id IS NOT NULL THEN 1 END)) as available_supplies
		FROM batch b
		LEFT JOIN medical_supply ms ON b.id = ms.batch_id
		LEFT JOIN (
			SELECT DISTINCT sh.medical_supply_id as supply_id
			FROM supply_history sh
			WHERE sh.status = 'consumido'
		) consumed ON ms.id = consumed.supply_id
		GROUP BY b.id, b.qr_code, b.amount, b.supplier, b.expiration_date
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

// CheckLowStockForSupplyType verifica stock bajo para todos los lotes de un tipo específico de insumo
func (s *BatchService) CheckLowStockForSupplyType(supplyCode string) error {
	// Obtener información del código de insumo
	var sc models.SupplyCode
	if err := s.DB.Where("code = ?", supplyCode).First(&sc).Error; err != nil {
		return fmt.Errorf("código de insumo no encontrado: %v", err)
	}

	// Obtener todos los lotes activos de este tipo de insumo
	var batches []models.Batch
	query := `
		SELECT DISTINCT b.* 
		FROM batch b
		JOIN medical_supply ms ON b.id = ms.batch_id
		WHERE ms.code = ? AND b.amount > 0 AND b.amount <= ?
	`

	if err := s.DB.Raw(query, supplyCode, sc.CriticalStock).Scan(&batches).Error; err != nil {
		return fmt.Errorf("error obteniendo lotes con stock bajo: %v", err)
	}

	// Enviar alerta para cada lote que esté en stock crítico
	alertsSent := 0
	for _, batch := range batches {
		if err := s.sendLowStockAlert(batch); err != nil {
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

// CheckAllBatchesLowStock verifica todos los lotes para alertas de stock bajo basándose en su stock crítico
func (s *BatchService) CheckAllBatchesLowStock() error {
	// Query para obtener lotes que están en stock crítico o por debajo
	query := `
		SELECT DISTINCT b.*, sc.critical_stock, sc.name as supply_name
		FROM batch b
		JOIN medical_supply ms ON b.id = ms.batch_id
		JOIN supply_code sc ON ms.code = sc.code
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
		if err := s.sendLowStockAlert(batchInfo.Batch); err != nil {
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

// StartAutomaticLowStockChecker inicia el proceso automático de verificación de stock bajo (ejecutar como goroutine)
func (s *BatchService) StartAutomaticLowStockChecker() {
	ticker := time.NewTicker(24 * time.Hour) // Verificar cada 24 horas
	defer ticker.Stop()

	fmt.Println("🔄 Iniciado verificador automático de stock bajo")

	for range ticker.C {
		fmt.Println("🔍 Ejecutando verificación automática de stock bajo...")
		if err := s.CheckAllBatchesLowStock(); err != nil {
			fmt.Printf("❌ Error en verificación automática de stock bajo: %v\n", err)
		}
	}
}

// GetLowStockSummary obtiene un resumen de todos los insumos con stock bajo
func (s *BatchService) GetLowStockSummary() ([]map[string]interface{}, error) {
	query := `
		SELECT 
			sc.code,
			sc.name,
			sc.critical_stock,
			SUM(b.amount) as total_current_stock,
			COUNT(DISTINCT b.id) as batch_count,
			MIN(b.expiration_date) as nearest_expiration,
			ARRAY_AGG(DISTINCT b.supplier) as suppliers
		FROM supply_code sc
		JOIN medical_supply ms ON sc.code = ms.code
		JOIN batch b ON ms.batch_id = b.id
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
		var suppliers string // PostgreSQL array como string

		err := rows.Scan(&code, &name, &criticalStock, &totalCurrentStock, &batchCount, &nearestExpiration, &suppliers)
		if err != nil {
			return nil, fmt.Errorf("error escaneando resultado: %v", err)
		}

		stockPercentage := float64(totalCurrentStock) / float64(criticalStock) * 100
		isUrgent := totalCurrentStock <= criticalStock/2 // Urgente si está al 50% o menos del stock crítico

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

// CheckAllBatchesExpiration verifica todos los lotes para alertas de vencimiento
func (s *BatchService) CheckAllBatchesExpiration() error {

	// Obtener lotes que no han vencido aún
	var activeBatches []models.Batch
	if err := s.DB.Where("expiration_date > ?", time.Now()).Find(&activeBatches).Error; err != nil {
		return fmt.Errorf("error obteniendo lotes activos para verificación: %v", err)
	}

	alertsSent := 0
	errors := 0

	// Verificar lotes activos
	for _, batch := range activeBatches {
		// Verificar si está próximo a vencer usando configuración del proveedor (0 = usar configuración)
		if err := s.CheckExpirationAlert(batch.ID, 0); err != nil {
			log.Printf("Error verificando alerta de vencimiento para lote %d: %v", batch.ID, err)
			errors++
			continue
		}
		alertsSent++
	}

	return nil
}

// sendLowStockAlert envía correo de alerta de stock bajo
func (s *BatchService) sendLowStockAlert(batch models.Batch) error {
	// Obtener información completa del insumo desde supply_code a través de medical_supply
	var supplyCode models.SupplyCode
	if err := s.DB.Joins("JOIN medical_supply ON medical_supply.code = supply_code.code").
		Where("medical_supply.batch_id = ?", batch.ID).
		First(&supplyCode).Error; err != nil {
		return fmt.Errorf("error al obtener información del código de insumo: %v", err)
	}

	// Crear datos para la plantilla
	data := map[string]interface{}{
		"BatchID":       batch.ID,
		"Code":          supplyCode.Code,
		"Name":          supplyCode.Name,
		"CurrentStock":  batch.Amount,
		"CriticalStock": supplyCode.CriticalStock,
		"Date":          time.Now().Format("02/01/2006"),
	}

	// Crear solicitud de correo
	req := mailer.NewRequest([]string{"vergara.javiera12@gmail.com"}, "Alerta: Stock Bajo en Lote")

	// Enviar correo usando la plantilla de stock bajo
	templatePath := "mailer/templates/low_stock.html"
	return req.SendMailSkipTLS(templatePath, data)
}

// getExpirationAlertDays obtiene los días de alerta de vencimiento para un proveedor
// Si no existe configuración, retorna 90 días (3 meses) por defecto
func (s *BatchService) getExpirationAlertDays(supplierName string) int {
	var config models.SupplierConfig
	if err := s.DB.Where("supplier_name = ?", supplierName).First(&config).Error; err != nil {
		// Si no existe configuración, usar default de 90 días (3 meses)
		return 90
	}
	return config.ExpirationAlertDays
}

// sendExpirationAlert envía correo de alerta de vencimiento próximo
func (s *BatchService) sendExpirationAlert(batch models.Batch) error {
	// Obtener información completa del insumo desde supply_code a través de medical_supply
	var supplyCode models.SupplyCode
	if err := s.DB.Joins("JOIN medical_supply ON medical_supply.code = supply_code.code").
		Where("medical_supply.batch_id = ?", batch.ID).
		First(&supplyCode).Error; err != nil {
		return fmt.Errorf("error al obtener información del código de vencimiento: %v", err)
	}

	// Crear datos para la plantilla
	data := map[string]interface{}{
		"BatchID":        batch.ID,
		"Code":           supplyCode.Code,
		"Name":           supplyCode.Name,
		"ExpirationDate": batch.ExpirationDate.Format("02/01/2006"),
		"DaysUntilExp":   int(time.Until(batch.ExpirationDate).Hours() / 24),
		"Date":           time.Now().Format("02/01/2006"),
	}

	// Crear solicitud de correo
	req := mailer.NewRequest([]string{"vergara.javiera12@gmail.com"}, "Alerta: Lote Próximo a Vencer")

	// Enviar correo usando la plantilla de vencimiento
	templatePath := "mailer/templates/expiration_warning.html"
	return req.SendMailSkipTLS(templatePath, data)
}
