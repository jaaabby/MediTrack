package services

import (
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

	return nil
}

// CreateBatchWithIndividualSupplies crea un lote junto con sus insumos individuales
func (s *BatchService) CreateBatchWithIndividualSupplies(batch *models.Batch, supplyCode *models.SupplyCode, individualCount int) (*models.Batch, []models.MedicalSupply, error) {
	var individualSupplies []models.MedicalSupply

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Generar QR único para el lote
		if batch.QRCode == "" {
			qrCode, err := s.QRService.GenerateUniqueQRCode("BATCH")
			if err != nil {
				return fmt.Errorf("error generando QR del lote: %v", err)
			}
			batch.QRCode = qrCode
		}

		// 2. Crear el lote
		if err := tx.Create(batch).Error; err != nil {
			return fmt.Errorf("error creando lote: %v", err)
		}

		// 3. Crear o actualizar supply_code
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
			if err := tx.Save(&existingSupplyCode).Error; err != nil {
				return fmt.Errorf("error actualizando código de insumo: %v", err)
			}
		}

		// 4. Crear insumos individuales
		if s.MedicalSupplyService != nil {
			supplies, err := s.MedicalSupplyService.CreateMultipleIndividualSuppliesTx(
				tx,
				batch.ID,
				supplyCode.Code,
				individualCount,
			)
			if err != nil {
				return fmt.Errorf("error creando insumos individuales: %v", err)
			}
			individualSupplies = supplies
		} else {
			// Fallback: crear manualmente si el servicio no está disponible
			for i := 0; i < individualCount; i++ {
				qrCode, err := s.QRService.GenerateUniqueQRCode("SUPPLY")
				if err != nil {
					return fmt.Errorf("error generando QR para insumo %d: %v", i+1, err)
				}

				supply := models.MedicalSupply{
					Code:    supplyCode.Code,
					QRCode:  qrCode,
					BatchID: batch.ID,
				}

				if err := tx.Create(&supply).Error; err != nil {
					return fmt.Errorf("error creando insumo individual %d: %v", i+1, err)
				}

				individualSupplies = append(individualSupplies, supply)
			}
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

	batch.ExpirationDate = newBatch.ExpirationDate
	batch.Amount = newBatch.Amount
	batch.Supplier = newBatch.Supplier
	batch.StoreID = newBatch.StoreID

	if err := s.DB.Save(&batch).Error; err != nil {
		return nil, err
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

	// Verificar si se debe enviar correo por stock bajo
	if batch.Amount < 10 {
		if err := s.sendLowStockAlert(batch); err != nil {
			// Solo log del error, no fallar la actualización
			fmt.Printf("Error al enviar alerta de stock bajo: %v\n", err)
		}
	}

	// Verificar si se debe enviar correo por vencimiento próximo (30 días)
	expirationThreshold := time.Now().AddDate(0, 0, 30)
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

// CheckLowStockAlert verifica y envía alertas de stock bajo
func (s *BatchService) CheckLowStockAlert(batchID int, threshold int) error {
	if threshold <= 0 {
		threshold = 5 // Threshold por defecto
	}

	var batch models.Batch
	if err := s.DB.First(&batch, batchID).Error; err != nil {
		return fmt.Errorf("lote no encontrado: %v", err)
	}

	if batch.Amount > 0 && batch.Amount <= threshold {
		return s.sendLowStockAlert(batch)
	}

	return nil
}

// CheckExpirationAlert verifica y envía alertas de vencimiento próximo
func (s *BatchService) CheckExpirationAlert(batchID int, daysThreshold int) error {
	if daysThreshold <= 0 {
		daysThreshold = 30 // 30 días por defecto
	}

	var batch models.Batch
	if err := s.DB.First(&batch, batchID).Error; err != nil {
		return fmt.Errorf("lote no encontrado: %v", err)
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

// startDailyScheduler inicia el scheduler diario para verificar vencimientos
func (s *BatchService) startDailyScheduler() {
	ticker := time.NewTicker(24 * time.Hour) // Verificar cada 24 horas
	defer ticker.Stop()

	// Ejecutar inmediatamente al iniciar
	s.CheckAllBatchesExpiration()

	for {
		select {
		case <-ticker.C:
			s.CheckAllBatchesExpiration()
		}
	}
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
		// Verificar si está próximo a vencer (30 días por defecto)
		if err := s.CheckExpirationAlert(batch.ID, 30); err != nil {
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
		"BatchID":      batch.ID,
		"Code":         supplyCode.Code,
		"Name":         supplyCode.Name,
		"CurrentStock": batch.Amount,
		"Date":         time.Now().Format("02/01/2006"),
	}

	// Crear solicitud de correo
	req := mailer.NewRequest([]string{"vergara.javiera12@gmail.com"}, "Alerta: Stock Bajo en Lote")

	// Enviar correo usando la plantilla de stock bajo
	templatePath := "mailer/templates/low_stock.html"
	return req.SendMailSkipTLS(templatePath, data)
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
