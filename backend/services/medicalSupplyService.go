package services

import (
	"errors"
	"fmt"
	"log"
	"meditrack/mailer"
	"meditrack/models"
	"meditrack/pkg"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

type MedicalSupplyService struct {
	DB        *gorm.DB
	QRService *QRService
}

// CreateMultipleIndividualSuppliesTx crea múltiples insumos individuales usando una transacción existente
func (s *MedicalSupplyService) CreateMultipleIndividualSuppliesTx(tx *gorm.DB, batchID int, code int, quantity int, storeID int) ([]models.MedicalSupply, error) {
	var supplies []models.MedicalSupply
	for i := 0; i < quantity; i++ {
		supply := models.MedicalSupply{
			Code:         code,
			BatchID:      batchID,
			LocationType: models.SupplyLocationStore,
			LocationID:   storeID,
			Status:       models.StatusAvailable,
			// QRCode se generará después de crear el insumo para tener el ID
		}

		// Crear el insumo primero para obtener el ID
		if err := tx.Create(&supply).Error; err != nil {
			return nil, fmt.Errorf("error creando insumo %d: %v", i+1, err)
		}

		// Generar QR único usando el ID real del insumo
		if s.QRService != nil {
			qrCode, err := s.QRService.generateUniqueQRCode("SUPPLY", supply.ID)
			if err != nil {
				return nil, fmt.Errorf("error generando QR para insumo %d: %v", i+1, err)
			}

			// Actualizar el insumo con el QR generado
			if err := tx.Model(&supply).Update("qr_code", qrCode).Error; err != nil {
				return nil, fmt.Errorf("error actualizando QR del insumo %d: %v", i+1, err)
			}
			supply.QRCode = qrCode
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
	// Crear el insumo primero para obtener el ID
	if err := s.DB.Create(supply).Error; err != nil {
		return err
	}

	// Si hay QRService disponible y no hay QR, generar uno usando el ID real
	if s.QRService != nil && supply.QRCode == "" {
		qrCode, err := s.QRService.generateUniqueQRCode("SUPPLY", supply.ID)
		if err != nil {
			return fmt.Errorf("error generando código QR: %v", err)
		}

		// Actualizar el insumo con el QR generado
		if err := s.DB.Model(supply).Update("qr_code", qrCode).Error; err != nil {
			return fmt.Errorf("error actualizando QR del insumo: %v", err)
		}
		supply.QRCode = qrCode
	}

	return nil
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

	// Guardar el estado anterior para verificar si cambió a "Recepcionado"
	previousStatus := supply.Status

	// Actualizar campos omitiendo ID, QRCode (para mantener trazabilidad) y UpdatedAt
	if err := s.DB.Model(&supply).Omit("id", "qr_code", "updated_at").Updates(newSupply).Error; err != nil {
		return nil, err
	}

	// Verificar si el cambio es a estado "Recepcionado" y programar alerta
	if newSupply.Status == models.StatusReceived && previousStatus != models.StatusReceived {
		// Programar verificación de alerta después de 1 minuto (para pruebas)
		go s.scheduleUnconsumedSupplyAlert(supply.ID)
		fmt.Printf("Insumo %d cambió a estado 'Recepcionado'. Alerta programada para 1 minuto.\n", supply.ID)
	}

	return &supply, nil
}

func (s *MedicalSupplyService) DeleteMedicalSupply(id int) error {
	return s.DB.Delete(&models.MedicalSupply{}, id).Error
}

// ===== INVENTARIO BÁSICO (de la versión anterior - RESTAURADO) =====

// GetInventoryList obtiene el inventario básico combinando datos de batch y supplyCode
// GetInventoryList obtiene el inventario basado en store_inventory_summary para consistencia
// Esto asegura que los cambios en cantidad y estado se reflejen correctamente en todas las vistas
func (s *MedicalSupplyService) GetInventoryList() ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	// Query que usa store_inventory_summary como fuente de verdad para consistencia
	// y además entrega información de ubicación y timestamps para el frontend.
	query := `
		SELECT DISTINCT ON (b.id)
			b.id AS batch_id,
			b.expiration_date,
			b.expiration_alert_days,
			COALESCE(
				(SELECT SUM(sis2.current_in_store)
				 FROM store_inventory_summary sis2
				 WHERE sis2.batch_id = b.id),
				b.amount
			) AS amount,
			COALESCE(
				(SELECT MAX(sis2.original_amount)
				 FROM store_inventory_summary sis2
				 WHERE sis2.batch_id = b.id),
				b.amount
			) AS original_amount,
			supc.supplier_name AS supplier,
			b.store_id,
			sc.code AS supply_code,
			sc.name AS supply_name,
			COALESCE(sc.critical_stock, 1) AS critical_stock,
			COALESCE(ms.location_type, 'store') AS location_type,
			COALESCE(ms.location_id, b.store_id) AS location_id,
			CASE
				WHEN COALESCE(ms.location_type, 'store') = 'store'
					THEN st.name
				WHEN COALESCE(ms.location_type, 'store') = 'pavilion'
					THEN pv.name
				ELSE st.name
			END AS location_name,
			MIN(ms.updated_at) OVER (PARTITION BY b.id) AS created_at,
			MAX(ms.updated_at) OVER (PARTITION BY b.id) AS updated_at
		FROM batch b
		LEFT JOIN medical_supply ms ON b.id = ms.batch_id
		LEFT JOIN supply_code sc ON ms.code = sc.code
		LEFT JOIN store st ON b.store_id = st.id
		LEFT JOIN pavilion pv ON b.location_type = 'pavilion' AND b.location_id = pv.id
		LEFT JOIN supplier_config supc ON b.supplier_id = supc.id
		ORDER BY b.id, sc.code
	`

	rows, err := s.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			item                = make(map[string]interface{})
			batchID             int
			expirationDate      string
			expirationAlertDays int
			amount              int
			originalAmount      int
			supplier            string
			storeID             int
			supplyCode          *int
			supplyName          *string
			criticalStock       int
			locationType        string
			locationID          int
			locationName        *string
			createdAt           *time.Time
			updatedAt           *time.Time
		)

		err := rows.Scan(
			&batchID,
			&expirationDate,
			&expirationAlertDays,
			&amount,
			&originalAmount,
			&supplier,
			&storeID,
			&supplyCode,
			&supplyName,
			&criticalStock,
			&locationType,
			&locationID,
			&locationName,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		item["batch_id"] = batchID
		item["expiration_date"] = expirationDate
		item["expiration_alert_days"] = expirationAlertDays
		item["amount"] = amount
		item["original_amount"] = originalAmount
		item["supplier"] = supplier
		item["store_id"] = storeID
		item["code"] = supplyCode
		item["name"] = supplyName
		item["critical_stock"] = criticalStock
		item["location_type"] = locationType
		item["location_id"] = locationID
		item["location_name"] = locationName
		item["created_at"] = createdAt
		item["updated_at"] = updatedAt

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
			supc.supplier_name AS supplier,
			sc.code as supply_code,
			sc.name as supply_name,
			COUNT(ms.id) as total_individual_supplies,
			COUNT(CASE WHEN consumed_supplies.supply_id IS NOT NULL THEN 1 END) as consumed_supplies,
			(COUNT(ms.id) - COUNT(CASE WHEN consumed_supplies.supply_id IS NOT NULL THEN 1 END)) as available_supplies
		FROM batch b
		LEFT JOIN medical_supply ms ON b.id = ms.batch_id
		LEFT JOIN supply_code sc ON ms.code = sc.code
		LEFT JOIN supplier_config supc ON b.supplier_id = supc.id
		LEFT JOIN (
			SELECT DISTINCT sh.medical_supply_id as supply_id
			FROM supply_history sh
			WHERE sh.status = 'consumido'
		) consumed_supplies ON ms.id = consumed_supplies.supply_id
		GROUP BY b.id, b.qr_code, b.expiration_date, b.amount, supc.supplier_name, sc.code, sc.name
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
				qrCode, err := s.QRService.generateUniqueQRCode("SUPPLY", 0)
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

// ConsumeSupplyByQR marca un insumo como consumido y actualiza la cantidad del lote y los resúmenes de inventario
func (s *MedicalSupplyService) ConsumeSupplyByQR(qrCode string, userRUT string, destinationType string, destinationID int) (*models.MedicalSupply, error) {
	var supply models.MedicalSupply

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Buscar el insumo por QR
		if err := tx.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
			return fmt.Errorf("insumo no encontrado con QR %s: %v", qrCode, err)
		}

		// Verificar que el insumo no haya sido consumido previamente
		if supply.Status == models.StatusConsumed {
			return fmt.Errorf("el insumo con QR %s ya ha sido consumido", qrCode)
		}

		// Obtener información del lote antes de la actualización
		var batch models.Batch
		if err := tx.Where("id = ?", supply.BatchID).First(&batch).Error; err != nil {
			return fmt.Errorf("error obteniendo lote: %v", err)
		}

		// Verificar que hay stock disponible
		if batch.Amount <= 0 {
			return fmt.Errorf("no hay stock disponible en el lote %d", batch.ID)
		}

		// Actualizar el estado del insumo a consumido
		// El trigger automáticamente creará el registro en supply_history
		if err := tx.Model(&supply).Update("status", models.StatusConsumed).Error; err != nil {
			return fmt.Errorf("error actualizando estado del insumo: %v", err)
		}

		// Actualizar cantidad del lote (restar 1)
		newAmount := batch.Amount - 1
		if err := tx.Model(&batch).Update("amount", newAmount).Error; err != nil {
			return fmt.Errorf("error actualizando cantidad del lote: %v", err)
		}

		// Actualizar resúmenes de inventario según la ubicación del insumo
		now := time.Now()
		if supply.LocationType == models.SupplyLocationStore {
			// Insumo consumido desde bodega - actualizar store_inventory_summary
			var storeSummary models.StoreInventorySummary
			if err := tx.Where("store_id = ? AND batch_id = ?", supply.LocationID, supply.BatchID).
				First(&storeSummary).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Si no existe el resumen, calcular el stock real en bodega
					var realCount int64
					tx.Model(&models.MedicalSupply{}).
						Where("batch_id = ? AND location_type = ? AND location_id = ? AND status != ?",
							supply.BatchID, models.SupplyLocationStore, supply.LocationID, models.StatusConsumed).
						Count(&realCount)

					// Crear resumen con valores calculados
					storeSummary = models.StoreInventorySummary{
						StoreID:              supply.LocationID,
						BatchID:              supply.BatchID,
						SupplyCode:           supply.Code,
						SurgeryID:            batch.SurgeryID,
						OriginalAmount:       int(realCount) + 1, // Cantidad antes del consumo
						CurrentInStore:       int(realCount),     // Stock actual en bodega
						TotalConsumedInStore: 1,
						LastConsumedDate:     &now,
					}
					if err := tx.Create(&storeSummary).Error; err != nil {
						return fmt.Errorf("error creando resumen de bodega: %v", err)
					}
				} else {
					return fmt.Errorf("error obteniendo resumen de bodega: %v", err)
				}
			} else {
				// Actualizar resumen existente
				storeSummary.CurrentInStore--
				storeSummary.TotalConsumedInStore++
				storeSummary.LastConsumedDate = &now
				if err := tx.Save(&storeSummary).Error; err != nil {
					return fmt.Errorf("error actualizando resumen de bodega: %v", err)
				}
			}
		} else if supply.LocationType == models.SupplyLocationPavilion {
			// Insumo consumido desde pabellón - actualizar pavilion_inventory_summary
			var pavilionSummary models.PavilionInventorySummary
			if err := tx.Where("pavilion_id = ? AND batch_id = ?", supply.LocationID, supply.BatchID).
				First(&pavilionSummary).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Si no existe el resumen, crearlo
					pavilionSummary = models.PavilionInventorySummary{
						PavilionID:       supply.LocationID,
						BatchID:          supply.BatchID,
						SupplyCode:       supply.Code,
						TotalReceived:    1,
						CurrentAvailable: 0, // Ya no hay disponible
						TotalConsumed:    1,
						LastConsumedDate: &now,
					}
					if err := tx.Create(&pavilionSummary).Error; err != nil {
						return fmt.Errorf("error creando resumen de pabellón: %v", err)
					}
				} else {
					return fmt.Errorf("error obteniendo resumen de pabellón: %v", err)
				}
			} else {
				// Actualizar resumen existente
				pavilionSummary.CurrentAvailable--
				pavilionSummary.TotalConsumed++
				pavilionSummary.LastConsumedDate = &now
				if err := tx.Save(&pavilionSummary).Error; err != nil {
					return fmt.Errorf("error actualizando resumen de pabellón: %v", err)
				}
			}
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
			supc.supplier_name AS supplier,
			st.name as store_name,
			st.type as store_type
		FROM medical_supply ms
		JOIN supply_code sc ON ms.code = sc.code
		JOIN batch b ON ms.batch_id = b.id
		JOIN store st ON b.store_id = st.id
		JOIN supplier_config supc ON b.supplier_id = supc.id
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

	// Obtener el insumo para verificar su estado
	var supply models.MedicalSupply
	if err := s.DB.First(&supply, supplyID).Error; err != nil {
		return nil, fmt.Errorf("insumo no encontrado: %v", err)
	}

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
		"status":                 supply.Status,
		"is_consumed":            supply.IsConsumed(),
		"available_for_use":      supply.CanBeConsumed() && batchRemainingAmount > 0,
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
			if err := tx.Model(&models.MedicalSupply{}).
				Where("batch_id = ? AND status != ?", batch.ID, models.StatusConsumed).
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

	if err := s.DB.Where("batch_id = ? AND status != ?", batchID, models.StatusConsumed).Find(&supplies).Error; err != nil {
		return nil, err
	}

	return supplies, nil
}

// ===== FUNCIONALIDADES DE ALERTA PARA INSUMOS NO CONSUMIDOS =====

// scheduleUnconsumedSupplyAlert programa una verificación de alerta para un insumo recepcionado
func (s *MedicalSupplyService) scheduleUnconsumedSupplyAlert(supplyID int) {
	// Esperar 1 minuto para pruebas (cambiar a 12 * time.Hour en producción)
	time.Sleep(1 * time.Minute)

	// Verificar si el insumo sigue en estado "Recepcionado"
	supply, err := s.GetMedicalSupplyByID(supplyID)
	if err != nil {
		fmt.Printf("Error obteniendo insumo %d para verificación de alerta: %v\n", supplyID, err)
		return
	}

	// Si sigue en estado "Recepcionado", enviar alerta
	if supply.Status == models.StatusReceived {
		s.sendUnconsumedSupplyAlert(supply)
	}
}

// sendUnconsumedSupplyAlert envía una alerta por correo para un insumo no consumido
func (s *MedicalSupplyService) sendUnconsumedSupplyAlert(supply *models.MedicalSupply) {
	// Obtener información adicional del insumo
	supplyInfo, err := s.GetSupplyWithBatchInfo(supply.QRCode)
	if err != nil {
		fmt.Printf("Error obteniendo información del insumo %d: %v\n", supply.ID, err)
		return
	}

	// Calcular horas transcurridas desde la recepción
	hoursElapsed := int(time.Since(supply.UpdatedAt).Hours())

	// Preparar datos para el correo
	emailData := map[string]interface{}{
		"SupplyID":     supply.ID,
		"SupplyName":   supplyInfo["supply_name"],
		"SupplyCode":   supplyInfo["supply_code"],
		"QRCode":       supply.QRCode,
		"BatchID":      supplyInfo["batch_id"],
		"ReceivedAt":   supply.UpdatedAt.Format("2006-01-02 15:04:05"),
		"HoursElapsed": hoursElapsed,
		"Date":         time.Now().Format("2006-01-02 15:04:05"),
	}

	// Configurar el correo
	// Leer email de destino desde variable de entorno
	alertEmail := os.Getenv("ALERT_EMAIL")
	if alertEmail == "" {
		fmt.Printf("ALERT_EMAIL no configurado, no se enviará alerta para insumo %d\n", supply.ID)
		return
	}
	recipients := []string{alertEmail}

	request := mailer.NewRequest(recipients, "Alerta: Insumo Médico No Consumido - MediTrack")

	// Enviar el correo
	// Obtener el directorio actual y construir la ruta absoluta
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error obteniendo directorio de trabajo: %v\n", err)
		return
	}

	templatePath := filepath.Join(wd, "mailer", "templates", "unconsumed_supply_alert.html")
	fmt.Printf("Buscando plantilla en: %s\n", templatePath)

	if err := request.SendMailSkipTLS(templatePath, emailData); err != nil {
		fmt.Printf("Error enviando alerta de insumo no consumido: %v\n", err)
	} else {
		fmt.Printf("Alerta enviada para insumo %d que no ha sido consumido\n", supply.ID)
	}
}

// CheckUnconsumedSupplies verifica todos los insumos en estado "Recepcionado" y envía alertas si es necesario
func (s *MedicalSupplyService) CheckUnconsumedSupplies() error {
	// Obtener todos los insumos en estado "Recepcionado" que tengan más de 1 minuto (para pruebas)
	var supplies []models.MedicalSupply
	cutoffTime := time.Now().Add(-1 * time.Minute)

	if err := s.DB.Where("status = ? AND updated_at < ?", models.StatusReceived, cutoffTime).Find(&supplies).Error; err != nil {
		return fmt.Errorf("error obteniendo insumos recepcionados: %v", err)
	}

	// Enviar alerta para cada insumo encontrado
	for _, supply := range supplies {
		s.sendUnconsumedSupplyAlert(&supply)
	}

	return nil
}

// GetUnconsumedSupplies obtiene todos los insumos que están en estado "Recepcionado" por más de 1 minuto (para pruebas)
func (s *MedicalSupplyService) GetUnconsumedSupplies() ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	cutoffTime := time.Now().Add(-1 * time.Minute)

	query := `
		SELECT 
			ms.id as supply_id,
			ms.code as supply_code,
			ms.qr_code,
			ms.batch_id,
			ms.updated_at as received_at,
			sc.name as supply_name,
			b.expiration_date,
			supc.supplier_name AS supplier,
			st.name as store_name,
			EXTRACT(EPOCH FROM (NOW() - ms.updated_at))/3600 as hours_elapsed
		FROM medical_supply ms
		JOIN supply_code sc ON ms.code = sc.code
		JOIN batch b ON ms.batch_id = b.id
		JOIN store st ON b.store_id = st.id
		JOIN supplier_config supc ON b.supplier_id = supc.id
		WHERE ms.status = ? AND ms.updated_at < ?
		ORDER BY ms.updated_at ASC
	`

	rows, err := s.DB.Raw(query, models.StatusReceived, cutoffTime).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item map[string]interface{} = make(map[string]interface{})
		var supplyID int
		var supplyCode int
		var qrCode string
		var batchID int
		var receivedAt time.Time
		var supplyName string
		var expirationDate string
		var supplier string
		var storeName string
		var hoursElapsed float64

		err := rows.Scan(
			&supplyID, &supplyCode, &qrCode, &batchID, &receivedAt,
			&supplyName, &expirationDate, &supplier, &storeName, &hoursElapsed,
		)
		if err != nil {
			return nil, err
		}

		item["supply_id"] = supplyID
		item["supply_code"] = supplyCode
		item["qr_code"] = qrCode
		item["batch_id"] = batchID
		item["received_at"] = receivedAt.Format("2006-01-02 15:04:05")
		item["supply_name"] = supplyName
		item["expiration_date"] = expirationDate
		item["supplier"] = supplier
		item["store_name"] = storeName
		item["hours_elapsed"] = int(hoursElapsed)

		result = append(result, item)
	}

	return result, nil
}

// ===== FUNCIONALIDADES DE RETORNO A BODEGA =====

// GetSuppliesForReturn obtiene insumos que deben regresar a bodega (8 horas laborales sin consumir)
func (s *MedicalSupplyService) GetSuppliesForReturn() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// Buscar insumos recepcionados que están en un pabellón
	// Usar ms.updated_at como fecha de recepción (más confiable que el historial que puede tener problemas de zona horaria)
	query := `
		SELECT 
			ms.id,
			ms.qr_code,
			ms.status,
			sc.name as supply_name,
			sc.code as supply_code,
			b.id as batch_id,
			supc.supplier_name AS supplier,
			b.expiration_date,
			s.name as store_name,
			s.id as store_id,
			ms.updated_at as received_at,
			ms.location_id as pavilion_id,
			COALESCE(p.name, CONCAT('Pabellón ', ms.location_id::text)) as pavilion_name,
			b.amount as batch_amount,
			sc.code_supplier as supply_code_supplier
		FROM medical_supply ms
		JOIN supply_code sc ON ms.code = sc.code
		JOIN batch b ON ms.batch_id = b.id
		JOIN supplier_config supc ON b.supplier_id = supc.id
		JOIN store s ON b.store_id = s.id
		LEFT JOIN pavilion p ON p.id = ms.location_id
		WHERE ms.status = ? 
		AND ms.location_type = ?
		AND ms.location_id > 0
		ORDER BY ms.updated_at ASC
	`

	fmt.Printf("🔍 Buscando insumos recepcionados en pabellones...\n")
	rows, err := s.DB.Raw(query, models.StatusReceived, models.SupplyLocationPavilion).Rows()
	if err != nil {
		return nil, fmt.Errorf("error consultando insumos para retorno: %v", err)
	}
	defer rows.Close()

	// Importar el paquete de horas laborales
	// Calcular horas laborales para cada insumo
	config := pkg.DefaultBusinessHoursConfig()
	count := 0

	for rows.Next() {
		count++
		var (
			id                 int
			qrCode             string
			status             string
			supplyName         string
			supplyCode         int
			batchID            int
			supplier           string
			expirationDate     time.Time
			storeName          string
			storeID            int
			receivedAt         time.Time
			pavilionID         int
			pavilionName       string
			batchAmount        int
			supplyCodeSupplier int
		)

		if err := rows.Scan(&id, &qrCode, &status, &supplyName, &supplyCode, &batchID, &supplier, &expirationDate, &storeName, &storeID, &receivedAt, &pavilionID, &pavilionName, &batchAmount, &supplyCodeSupplier); err != nil {
			fmt.Printf("⚠️ Error escaneando fila %d: %v\n", count, err)
			continue
		}

		fmt.Printf("📦 Insumo encontrado: ID=%d, QR=%s, Status=%s, LocationType=%s, LocationID=%d, ReceivedAt=%s\n",
			id, qrCode, status, "pavilion", 0, receivedAt.Format("2006-01-02 15:04:05"))

		// Calcular horas laborales transcurridas desde la recepción
		now := time.Now()

		// Verificar si la recepción fue fuera del horario laboral (después de las 17:00 o antes de las 8:00)
		receivedHour := receivedAt.Hour()
		isOutsideBusinessHours := receivedHour >= config.EndHour || receivedHour < config.StartHour

		var businessHoursElapsed float64
		if isOutsideBusinessHours {
			// Si la recepción fue fuera del horario laboral, calcular horas laborales desde el próximo día
			// Si fue después de las 17:00, empezar a contar desde las 8:00 del día siguiente
			// Si fue antes de las 8:00, empezar desde las 8:00 del mismo día
			nextBusinessDayStart := receivedAt
			if receivedHour >= config.EndHour {
				// Si fue después de las 17:00, empezar desde las 8:00 del día siguiente
				nextBusinessDayStart = time.Date(receivedAt.Year(), receivedAt.Month(), receivedAt.Day(), config.StartHour, 0, 0, 0, receivedAt.Location()).AddDate(0, 0, 1)
			} else if receivedHour < config.StartHour {
				// Si fue antes de las 8:00, empezar desde las 8:00 del mismo día
				nextBusinessDayStart = time.Date(receivedAt.Year(), receivedAt.Month(), receivedAt.Day(), config.StartHour, 0, 0, 0, receivedAt.Location())
			}
			businessHoursElapsed = pkg.CalculateBusinessHours(nextBusinessDayStart, now, config)
			fmt.Printf("  ℹ️ Recepción fuera del horario laboral (hora: %d), calculando desde próximo horario laboral: %.4f horas laborales\n",
				receivedHour, businessHoursElapsed)
		} else {
			// Si fue dentro del horario laboral, calcular horas laborales normalmente
			businessHoursElapsed = pkg.CalculateBusinessHours(receivedAt, now, config)
		}

		// Usar 8 horas laborales como umbral
		shouldReturn := businessHoursElapsed >= 8.0

		fmt.Printf("  ⏱️ Horas laborales transcurridas: %.4f (umbral: %.2f), Debe retornar: %v\n",
			businessHoursElapsed, 8.0, shouldReturn)

		item := map[string]interface{}{
			"supply_id":              id,
			"qr_code":                qrCode,
			"status":                 status,
			"supply_name":            supplyName,
			"supply_code":            supplyCode,
			"batch_id":               batchID,
			"supplier":               supplier,
			"expiration_date":        expirationDate.Format("2006-01-02"),
			"store_name":             storeName,
			"store_id":               storeID,
			"received_at":            receivedAt.Format("2006-01-02 15:04:05"),
			"business_hours_elapsed": businessHoursElapsed,
			"should_return":          shouldReturn,
			"pavilion_id":            pavilionID,
			"pavilion_name":          pavilionName,
			"batch_amount":           batchAmount,
			"supply_code_supplier":   supplyCodeSupplier,
		}

		// Solo agregar si debe retornarse (8 horas laborales o más)
		if shouldReturn {
			results = append(results, item)
			fmt.Printf("✅ Insumo agregado para retorno: QR=%s\n", qrCode)
		} else {
			fmt.Printf("⏳ Insumo aún no debe retornar: QR=%s (%.4f < %.2f horas laborales)\n",
				qrCode, businessHoursElapsed, 8.0)
		}
	}

	fmt.Printf("📊 Total filas encontradas: %d, Total insumos para retorno: %d\n", count, len(results))
	return results, nil
}

// NotifyPavilionForReturn envía un correo electrónico al pabellón donde está el insumo, adjuntando su PDF/QR,
// para que el personal realice la devolución física a bodega.
// pdfBytes puede ser nil; en ese caso se envía sólo el cuerpo del correo sin adjunto.
func (s *MedicalSupplyService) NotifyPavilionForReturn(qrCode string, pdfBytes []byte) ([]string, error) {
	// 1. Buscar el insumo por QR
	var supply models.MedicalSupply
	if err := s.DB.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
		return nil, fmt.Errorf("insumo no encontrado con QR %s", qrCode)
	}

	// 2. Verificar que está en un pabellón
	if supply.LocationType != models.SupplyLocationPavilion {
		return nil, fmt.Errorf("el insumo no se encuentra en un pabellón (ubicación: %s)", supply.LocationType)
	}

	pavilionID := supply.LocationID

	// 3. Obtener nombre del insumo
	var supplyCode models.SupplyCode
	s.DB.Where("code = ?", supply.Code).First(&supplyCode)

	// 4. Obtener datos del lote y bodega
	var batch models.Batch
	s.DB.First(&batch, supply.BatchID)

	var store models.Store
	s.DB.First(&store, batch.StoreID)

	// 5. Obtener nombre del pabellón
	var pavilion models.Pavilion
	s.DB.First(&pavilion, pavilionID)

	pavilionName := pavilion.Name
	if pavilionName == "" {
		pavilionName = fmt.Sprintf("Pabellón %d", pavilionID)
	}

	// 6. Buscar usuarios del pabellón con email
	var users []models.User
	s.DB.Where(
		"pavilion_id = ? AND role = ? AND is_active = true AND email != ''",
		pavilionID, models.RolePavilion,
	).Find(&users)

	if len(users) == 0 {
		return nil, fmt.Errorf("no se encontraron usuarios con correo en el pabellón '%s' (ID %d)", pavilionName, pavilionID)
	}

	emails := make([]string, 0, len(users))
	for _, u := range users {
		if u.Email != "" {
			emails = append(emails, u.Email)
		}
	}
	if len(emails) == 0 {
		return nil, fmt.Errorf("los usuarios del pabellón '%s' no tienen correo registrado", pavilionName)
	}

	// 7. Preparar datos del correo
	emailData := struct {
		PavilionName   string
		SupplyName     string
		SupplyCode     string
		QRCode         string
		BatchID        int
		Supplier       string
		ExpirationDate string
		ReceivedAt     string
		StoreName      string
		Date           string
	}{
		PavilionName:   pavilionName,
		SupplyName:     supplyCode.Name,
		SupplyCode:     fmt.Sprintf("%d", supply.Code),
		QRCode:         qrCode,
		BatchID:        supply.BatchID,
		Supplier:       batch.Supplier,
		ExpirationDate: batch.ExpirationDate.Format("02/01/2006"),
		ReceivedAt:     supply.UpdatedAt.Format("02/01/2006 15:04"),
		StoreName:      store.Name,
		Date:           time.Now().Format("02/01/2006 15:04:05"),
	}

	// 8. Enviar correo
	subject := fmt.Sprintf("⚠️ Solicitud de Devolución de Insumo — %s", pavilionName)
	req := mailer.NewRequest(emails, subject)

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo directorio de trabajo: %v", err)
	}
	templatePath := filepath.Join(wd, "mailer", "templates", "pavilion_return_request.html")

	if len(pdfBytes) > 0 {
		filename := fmt.Sprintf("insumo-%s.pdf", strings.ReplaceAll(qrCode, "_", "-"))
		if err := req.SendMailWithAttachment(templatePath, emailData, filename, pdfBytes); err != nil {
			return nil, fmt.Errorf("error enviando correo con adjunto: %v", err)
		}
	} else {
		if err := req.SendMailSkipTLS(templatePath, emailData); err != nil {
			return nil, fmt.Errorf("error enviando correo: %v", err)
		}
	}

	fmt.Printf("✉️ Notificación de devolución enviada a %v para insumo QR=%s en pabellón '%s'\n",
		emails, qrCode, pavilionName)
	return emails, nil
}

// ReturnSupplyToStore regresa un insumo a bodega
func (s *MedicalSupplyService) ReturnSupplyToStore(supplyID int, userRUT string, notes string, isAutomatic bool) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Obtener el insumo
		var supply models.MedicalSupply
		if err := tx.First(&supply, supplyID).Error; err != nil {
			return fmt.Errorf("insumo no encontrado: %v", err)
		}

		// Verificar que pueda ser regresado
		// Permitir devolver insumos consumidos automáticamente
		if supply.Status == models.StatusConsumed {
			// Verificar si fue consumido automáticamente
			var lastConsumptionHistory models.SupplyHistory
			if err := tx.Where("medical_supply_id = ? AND status = ?", supply.ID, models.StatusConsumed).
				Order("date_time DESC").
				First(&lastConsumptionHistory).Error; err == nil {
				// Si las notas contienen el prefijo de consumo automático, permitir devolución
				if strings.Contains(lastConsumptionHistory.Notes, "[CONSUMO_AUTOMATICO]") {
					// Permitir devolución de insumo consumido automáticamente
					log.Printf("🔄 Permitiendo devolución de insumo %s consumido automáticamente", supply.QRCode)
				} else {
					return fmt.Errorf("no se puede regresar un insumo consumido manualmente")
				}
			} else {
				return fmt.Errorf("no se puede regresar un insumo consumido")
			}
		}

		// Obtener información del lote y bodega
		var batch models.Batch
		if err := tx.First(&batch, supply.BatchID).Error; err != nil {
			return fmt.Errorf("error obteniendo información del lote: %v", err)
		}

		var store models.Store
		if err := tx.First(&store, batch.StoreID).Error; err != nil {
			return fmt.Errorf("error obteniendo información de la bodega: %v", err)
		}

		now := time.Now()
		oldStatus := supply.Status
		oldLocationType := supply.LocationType
		oldLocationID := supply.LocationID

		// Cambiar estado a "en_camino_a_bodega" (NO a disponible todavía)
		// La ubicación se actualizará cuando se confirme la llegada
		supply.Status = models.StatusEnRouteToStore
		supply.InTransit = true
		// Mantener la ubicación actual hasta que se confirme la llegada
		if err := tx.Save(&supply).Error; err != nil {
			return fmt.Errorf("error actualizando estado del insumo: %v", err)
		}

		// Si el insumo estaba en un pabellón, decrementar el resumen del pabellón
		if oldLocationType == models.SupplyLocationPavilion && oldLocationID > 0 {
			var pavilionSummary models.PavilionInventorySummary
			if err := tx.Where("pavilion_id = ? AND batch_id = ?", oldLocationID, batch.ID).
				First(&pavilionSummary).Error; err == nil {
				pavilionSummary.CurrentAvailable--
				if pavilionSummary.CurrentAvailable < 0 {
					pavilionSummary.CurrentAvailable = 0
				}
				pavilionSummary.TotalReturned++
				pavilionSummary.LastReturnedDate = &now
				if err := tx.Save(&pavilionSummary).Error; err != nil {
					return fmt.Errorf("error actualizando resumen de pabellón: %v", err)
				}
				fmt.Printf("✅ Resumen de pabellón actualizado: PavilionID=%d, CurrentAvailable=%d\n",
					oldLocationID, pavilionSummary.CurrentAvailable)
			}
		}

		// NO incrementar el resumen de bodega todavía - se hará en ConfirmArrivalToStore
		// NO incrementar la cantidad del lote todavía - se hará en ConfirmArrivalToStore

		// Crear registro en supply_history con estado "en_camino_a_bodega"
		returnType := "manual"
		if isAutomatic {
			returnType = "automatico"
		}

		// Usar el mensaje que viene en notes (ya contiene las horas laborales correctas)
		finalNotes := fmt.Sprintf("Insumo en camino a bodega (%s): %s", returnType, notes)

		historyEntry := models.SupplyHistory{
			DateTime:        now,
			Status:          models.StatusEnRouteToStore,
			DestinationType: models.DestinationTypeStore,
			DestinationID:   store.ID,
			MedicalSupplyID: supply.ID,
			UserRUT:         userRUT,
			Notes:           finalNotes,
			OriginType:      &oldLocationType,
			OriginID:        &oldLocationID,
		}

		if err := tx.Create(&historyEntry).Error; err != nil {
			return fmt.Errorf("error creando historial: %v", err)
		}

		// Si el insumo está en un carrito, actualizar la asignación QR
		var assignment models.SupplyRequestQRAssignment
		if err := tx.Where("qr_code = ? AND status IN (?, ?, ?)",
			supply.QRCode,
			models.AssignmentStatusAssigned,
			models.AssignmentStatusDelivered,
			models.AssignmentStatusConsumed).
			Order("assigned_date DESC").
			First(&assignment).Error; err == nil {
			// Actualizar el estado de la asignación a "returned"
			assignment.Status = models.AssignmentStatusReturned
			if assignment.Notes != "" {
				assignment.Notes += "\n" + notes
			} else {
				assignment.Notes = notes
			}
			if err := tx.Save(&assignment).Error; err != nil {
				return fmt.Errorf("error actualizando asignación QR: %v", err)
			}
			fmt.Printf("✅ Asignación QR actualizada para insumo %s\n", supply.QRCode)
		}

		// Log para depuración
		fmt.Printf("✅ Insumo %s regresado a bodega %s (estado: %s -> %s)\n",
			supply.QRCode, store.Name, oldStatus, supply.Status)

		return nil
	})
}

// ReturnSupplyToStoreByQR regresa un insumo a bodega usando su código QR
func (s *MedicalSupplyService) ReturnSupplyToStoreByQR(qrCode string, userRUT string, notes string, isAutomatic bool) error {
	var supply models.MedicalSupply
	if err := s.DB.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
		return fmt.Errorf("insumo con QR %s no encontrado", qrCode)
	}

	return s.ReturnSupplyToStore(supply.ID, userRUT, notes, isAutomatic)
}

// ProcessAutomaticReturns procesa automáticamente los retornos de insumos que llevan 8 horas laborales sin consumir
func (s *MedicalSupplyService) ProcessAutomaticReturns() error {
	supplies, err := s.GetSuppliesForReturn()
	if err != nil {
		return fmt.Errorf("error obteniendo insumos para retorno: %v", err)
	}

	returnedCount := 0
	errorCount := 0
	var errors []string
	var returnedSupplies []map[string]interface{} // Para el correo

	for _, supply := range supplies {
		supplyID := supply["supply_id"].(int)
		qrCode := supply["qr_code"].(string)
		businessHoursElapsed := supply["business_hours_elapsed"].(float64)

		// Ya están filtrados por GetSuppliesForReturn (8 horas laborales), pero verificamos por seguridad
		if businessHoursElapsed >= 8.0 {
			// Usuario del sistema para operaciones automáticas (usar SYSTEM-INIT que existe en la BD)
			systemUserRUT := "SYSTEM-INIT"
			notes := fmt.Sprintf("Retorno automático después de %.1f horas laborales sin consumo", businessHoursElapsed)

			err := s.ReturnSupplyToStore(supplyID, systemUserRUT, notes, true)
			if err != nil {
				errorMsg := fmt.Sprintf("Error retornando insumo %s: %v", qrCode, err)
				fmt.Printf("❌ %s\n", errorMsg)
				errors = append(errors, errorMsg)
				errorCount++
				continue
			}

			returnedCount++
			fmt.Printf("✅ Insumo %s retornado automáticamente a bodega (%.1f horas laborales)\n", qrCode, businessHoursElapsed)

			// Guardar información para el correo
			returnedSupplies = append(returnedSupplies, map[string]interface{}{
				"SupplyName":  supply["supply_name"],
				"QRCode":      qrCode,
				"TimeElapsed": fmt.Sprintf("%.1f horas laborales", businessHoursElapsed),
				"StoreName":   supply["store_name"],
			})
		}
	}

	if returnedCount > 0 {
		fmt.Printf("📦 Procesamiento automático completado: %d insumos retornados a bodega\n", returnedCount)

		// Enviar correo con el resumen (en goroutine para no bloquear)
		go func() {
			if err := s.sendAutomaticReturnSummaryEmail(returnedCount, errorCount, returnedSupplies); err != nil {
				fmt.Printf("⚠️ Error enviando correo de resumen de retornos automáticos: %v\n", err)
			} else {
				fmt.Printf("📧 Correo de resumen enviado exitosamente\n")
			}
		}()
	}

	if errorCount > 0 {
		if returnedCount == 0 {
			// Si todos fallaron, retornar error
			return fmt.Errorf("todos los insumos fallaron al procesarse: %v", errors)
		}
		// Si algunos fallaron, solo loguear pero no retornar error
		fmt.Printf("⚠️ Advertencia: %d insumos fallaron al procesarse\n", errorCount)
	}

	return nil
}

// sendAutomaticReturnSummaryEmail envía un correo con el resumen de retornos automáticos
func (s *MedicalSupplyService) sendAutomaticReturnSummaryEmail(returnedCount int, errorCount int, returnedSupplies []map[string]interface{}) error {
	// Leer email de destino desde variable de entorno
	alertEmail := os.Getenv("ALERT_EMAIL")
	if alertEmail == "" {
		return fmt.Errorf("ALERT_EMAIL no configurado, no se puede enviar resumen de retornos automáticos")
	}
	recipients := []string{alertEmail}

	// Preparar datos para el correo
	emailData := map[string]interface{}{
		"ReturnedCount":    returnedCount,
		"ErrorCount":       errorCount,
		"ProcessDate":      time.Now().Format("02/01/2006 15:04:05"),
		"ReturnedSupplies": returnedSupplies,
		"Date":             time.Now().Format("02/01/2006 15:04:05"),
	}

	request := mailer.NewRequest(recipients, fmt.Sprintf("Resumen de Retornos Automáticos - %d insumo(s) retornado(s)", returnedCount))

	// Obtener el directorio actual y construir la ruta absoluta
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error obteniendo directorio de trabajo: %v", err)
	}

	templatePath := filepath.Join(wd, "mailer", "templates", "automatic_return_summary.html")
	fmt.Printf("📧 Enviando correo de resumen de retornos automáticos a %s\n", alertEmail)

	if err := request.SendMailSkipTLS(templatePath, emailData); err != nil {
		return fmt.Errorf("error enviando correo de resumen: %v", err)
	}

	return nil
}

// StartAutomaticReturnChecker inicia el proceso automático de verificación de retornos (ejecutar como goroutine)
func (s *MedicalSupplyService) StartAutomaticReturnChecker() {
	ticker := time.NewTicker(24 * time.Hour) // Verificar cada 24 horas
	defer ticker.Stop()

	fmt.Println("🔄 Iniciado verificador automático de retornos a bodega")

	for range ticker.C {
		fmt.Println("🔍 Ejecutando verificación automática de retornos...")
		if err := s.ProcessAutomaticReturns(); err != nil {
			fmt.Printf("❌ Error en verificación automática: %v\n", err)
		}
	}
}

// ===== ALERTAS DE STOCK BAJO PARA INSUMOS INDIVIDUALES =====

// CheckLowStockForIndividualSupplies verifica y envía alertas de stock bajo para insumos individuales
// Alerta cuando quede exactamente 1 insumo disponible por código
func (s *MedicalSupplyService) CheckLowStockForIndividualSupplies(supplyCode int) error {
	// Contar insumos disponibles por código
	var availableCount int64
	if err := s.DB.Model(&models.MedicalSupply{}).
		Where("code = ? AND status = ?", supplyCode, models.StatusAvailable).
		Count(&availableCount).Error; err != nil {
		return fmt.Errorf("error contando insumos disponibles: %v", err)
	}

	// Alertar solo cuando quede exactamente 1 insumo
	if availableCount == 1 {
		return s.sendLowStockAlertForSupply(supplyCode, int(availableCount))
	}

	return nil
}

// CheckAllIndividualSuppliesLowStock verifica todos los códigos de insumos para alertas de stock bajo
func (s *MedicalSupplyService) CheckAllIndividualSuppliesLowStock() error {
	// Query para obtener códigos con exactamente 1 insumo disponible
	query := `
		SELECT 
			ms.code,
			sc.name as supply_name,
			COUNT(*) as available_count
		FROM medical_supply ms
		JOIN supply_code sc ON ms.code = sc.code
		WHERE ms.status = ?
		GROUP BY ms.code, sc.name
		HAVING COUNT(*) = 1
		ORDER BY ms.code
	`

	type SupplyStockInfo struct {
		Code           int    `json:"code"`
		SupplyName     string `json:"supply_name"`
		AvailableCount int    `json:"available_count"`
	}

	var suppliesWithLowStock []SupplyStockInfo
	if err := s.DB.Raw(query, models.StatusAvailable).Scan(&suppliesWithLowStock).Error; err != nil {
		return fmt.Errorf("error obteniendo insumos con stock bajo: %v", err)
	}

	alertsSent := 0
	errors := 0

	for _, supplyInfo := range suppliesWithLowStock {
		if err := s.sendLowStockAlertForSupply(supplyInfo.Code, supplyInfo.AvailableCount); err != nil {
			log.Printf("Error enviando alerta de stock bajo para insumo %d (%s): %v",
				supplyInfo.Code, supplyInfo.SupplyName, err)
			errors++
			continue
		}
		alertsSent++
		log.Printf("Alerta enviada para insumo %d (%s): Queda 1 insumo disponible",
			supplyInfo.Code, supplyInfo.SupplyName)
	}

	if alertsSent > 0 {
		log.Printf("Verificación de stock bajo de insumos completada: %d alertas enviadas, %d errores", alertsSent, errors)
	}

	return nil
}

// sendLowStockAlertForSupply envía correo de alerta de stock bajo para insumos individuales
func (s *MedicalSupplyService) sendLowStockAlertForSupply(supplyCode int, availableCount int) error {
	// Obtener información del código de insumo
	var supplyCodeInfo models.SupplyCode
	if err := s.DB.Where("code = ?", supplyCode).First(&supplyCodeInfo).Error; err != nil {
		return fmt.Errorf("error al obtener información del código de insumo: %v", err)
	}

	// Crear datos para la plantilla
	data := map[string]interface{}{
		"Code":          supplyCodeInfo.Code,
		"Name":          supplyCodeInfo.Name,
		"CurrentStock":  availableCount,
		"CriticalStock": 1, // Para insumos individuales, el crítico es 1
		"Date":          time.Now().Format("02/01/2006"),
		"IsIndividual":  true, // Indicar que es alerta para insumos individuales
	}

	// Crear solicitud de correo
	req := mailer.NewRequest([]string{"vergara.javiera12@gmail.com"}, "Alerta: Stock Crítico - Queda 1 Insumo")

	// Usar la misma plantilla de stock bajo (o crear una específica)
	templatePath := "mailer/templates/low_stock.html"
	return req.SendMailSkipTLS(templatePath, data)
}

// ConfirmArrivalToStore confirma la llegada de un insumo a bodega y lo cambia a estado disponible.
// Además, recepciona automáticamente otros insumos del mismo carrito que estén en camino a bodega.
func (s *MedicalSupplyService) ConfirmArrivalToStore(qrCode string, userRUT string, notes string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Obtener nombre del usuario
		var user models.User
		var userName string
		if err := tx.Where("rut = ?", userRUT).First(&user).Error; err == nil {
			userName = user.Name
		} else {
			userName = userRUT // Fallback
		}

		// Obtener el insumo por QR
		var supply models.MedicalSupply
		if err := tx.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
			return fmt.Errorf("insumo con QR %s no encontrado", qrCode)
		}

		// Verificar que esté en estado "en_camino_a_bodega"
		if supply.Status != models.StatusEnRouteToStore {
			return fmt.Errorf("el insumo no está en camino a bodega (estado actual: %s)", supply.Status)
		}

		// Obtener información del lote y bodega
		var batch models.Batch
		if err := tx.First(&batch, supply.BatchID).Error; err != nil {
			return fmt.Errorf("error obteniendo información del lote: %v", err)
		}

		var store models.Store
		if err := tx.First(&store, batch.StoreID).Error; err != nil {
			return fmt.Errorf("error obteniendo información de la bodega: %v", err)
		}

		// ============================
		// 1) Actualizar transferencia de devolución principal (QR escaneado)
		// ============================

		var transfer models.SupplyTransfer
		if err := tx.Where("qr_code = ? AND status = ?", qrCode, models.TransferStatusInTransit).
			First(&transfer).Error; err != nil {
			// Si no se encuentra por QR code, intentar por medical_supply_id
			if err := tx.Where("medical_supply_id = ? AND status = ? AND destination_type = ?",
				supply.ID, models.TransferStatusInTransit, models.TransferLocationStore).
				Order("send_date DESC").First(&transfer).Error; err != nil {
				fmt.Printf("⚠️ No se encontró transferencia en tránsito para QR %s (puede que ya haya sido recibida)\n", qrCode)
			} else {
				// Transferencia encontrada por medical_supply_id - actualizar a recibida
				now := time.Now()
				transfer.Status = models.TransferStatusReceived
				transfer.ReceiveDate = &now
				transfer.ReceivedBy = &userRUT

				var u models.User
				if err := tx.Where("rut = ?", userRUT).First(&u).Error; err == nil {
					transfer.ReceivedByName = &u.Name
				}
				if notes != "" {
					if transfer.Notes != "" {
						transfer.Notes = transfer.Notes + "\n" + notes
					} else {
						transfer.Notes = notes
					}
				}
				if err := tx.Save(&transfer).Error; err != nil {
					return fmt.Errorf("error actualizando transferencia: %v", err)
				}
				fmt.Printf("✅ Transferencia de devolución actualizada (por medical_supply_id): Code=%s, Status=%s\n",
					transfer.TransferCode, transfer.Status)
			}
		} else {
			// Transferencia encontrada por QR code - actualizar a recibida
			now := time.Now()
			transfer.Status = models.TransferStatusReceived
			transfer.ReceiveDate = &now
			transfer.ReceivedBy = &userRUT

			var u models.User
			if err := tx.Where("rut = ?", userRUT).First(&u).Error; err == nil {
				transfer.ReceivedByName = &u.Name
			}
			if notes != "" {
				if transfer.Notes != "" {
					transfer.Notes = transfer.Notes + "\n" + notes
				} else {
					transfer.Notes = notes
				}
			}
			if err := tx.Save(&transfer).Error; err != nil {
				return fmt.Errorf("error actualizando transferencia: %v", err)
			}
			fmt.Printf("✅ Transferencia de devolución actualizada (por QR code): Code=%s, Status=%s\n",
				transfer.TransferCode, transfer.Status)
		}

		// ============================
		// 2) Actualizar inventario y lote para el insumo principal
		// ============================

		var storeSummary models.StoreInventorySummary
		if err := tx.Where("store_id = ? AND batch_id = ?", store.ID, batch.ID).
			First(&storeSummary).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Crear resumen si no existe
				now := time.Now()
				storeSummary = models.StoreInventorySummary{
					StoreID:          store.ID,
					BatchID:          batch.ID,
					SupplyCode:       supply.Code,
					OriginalAmount:   batch.Amount,
					CurrentInStore:   1,
					TotalReturnedIn:  1,
					LastReturnInDate: &now,
				}
				if err := tx.Create(&storeSummary).Error; err != nil {
					return fmt.Errorf("error creando resumen de bodega: %v", err)
				}
			} else {
				return fmt.Errorf("error obteniendo resumen de bodega: %v", err)
			}
		} else {
			// Actualizar resumen existente
			now := time.Now()
			storeSummary.CurrentInStore++
			storeSummary.TotalReturnedIn++
			storeSummary.LastReturnInDate = &now
			if err := tx.Save(&storeSummary).Error; err != nil {
				return fmt.Errorf("error actualizando resumen de bodega: %v", err)
			}
		}

		// Incrementar cantidad del lote
		batch.Amount++
		if err := tx.Save(&batch).Error; err != nil {
			return fmt.Errorf("error actualizando lote: %v", err)
		}

		// ============================
		// 3) Marcar asignaciones antiguas como 'returned' y desactivar items de carrito (QR principal)
		// ============================

		var oldAssignments []models.SupplyRequestQRAssignment
		if err := tx.Where("qr_code = ? AND status NOT IN (?, ?)", qrCode, models.AssignmentStatusReturned, models.AssignmentStatusConsumed).
			Find(&oldAssignments).Error; err == nil {
			for _, assignment := range oldAssignments {
				assignment.Status = models.AssignmentStatusReturned
				if notes != "" {
					if assignment.Notes != "" {
						assignment.Notes = assignment.Notes + "\n[DEVUELTO A BODEGA] " + notes
					} else {
						assignment.Notes = "[DEVUELTO A BODEGA] " + notes
					}
				}
				if err := tx.Save(&assignment).Error; err != nil {
					fmt.Printf("⚠️ Error actualizando asignación a 'returned': %v\n", err)
				}

				// Desactivar items de carrito asociados a esta asignación
				now := time.Now()
				userNamePtr := userName
				if err := tx.Model(&models.SupplyCartItem{}).
					Where("supply_request_qr_assignment_id = ? AND is_active = ?", assignment.ID, true).
					Updates(map[string]interface{}{
						"is_active":       false,
						"removed_at":      &now,
						"removed_by":      &userRUT,
						"removed_by_name": &userNamePtr,
						"notes":           "Insumo devuelto a bodega",
					}).Error; err != nil {
					fmt.Printf("⚠️ Error desactivando items de carrito: %v\n", err)
				}
			}
			fmt.Printf("✅ Marcadas %d asignaciones antiguas como 'returned' para QR %s\n", len(oldAssignments), qrCode)
		}

		// ============================
		// 4) Cambiar estado del insumo principal y registrar historial
		// ============================

		oldStatus := supply.Status
		supply.Status = models.StatusAvailable
		supply.LocationType = models.SupplyLocationStore
		supply.LocationID = store.ID
		supply.InTransit = false
		if err := tx.Save(&supply).Error; err != nil {
			return fmt.Errorf("error actualizando estado del insumo: %v", err)
		}

		finalNotes := fmt.Sprintf("Llegada confirmada a bodega %s. %s", store.Name, notes)
		if notes == "" {
			finalNotes = fmt.Sprintf("Llegada confirmada a bodega %s", store.Name)
		}

		historyEntry := models.SupplyHistory{
			DateTime:        time.Now(),
			Status:          models.StatusAvailable,
			DestinationType: models.DestinationTypeStore,
			DestinationID:   store.ID,
			MedicalSupplyID: supply.ID,
			UserRUT:         userRUT,
			Notes:           finalNotes,
		}

		if err := tx.Create(&historyEntry).Error; err != nil {
			return fmt.Errorf("error creando historial: %v", err)
		}

		fmt.Printf("✅ Llegada confirmada para insumo %s a bodega %s (estado: %s -> %s)\n",
			supply.QRCode, store.Name, oldStatus, supply.Status)

		// ============================
		// 5) Recepción automática de otros insumos del mismo carrito
		// ============================

		var assignmentForCart models.SupplyRequestQRAssignment
		if err := tx.Where("medical_supply_id = ?", supply.ID).
			Order("assigned_date DESC").
			First(&assignmentForCart).Error; err == nil {

			var cartItem models.SupplyCartItem
			if err := tx.Where("supply_request_qr_assignment_id = ? AND is_active = ?", assignmentForCart.ID, true).
				First(&cartItem).Error; err == nil {

				var cartItems []models.SupplyCartItem
				if err := tx.Where("supply_cart_id = ? AND is_active = ?", cartItem.SupplyCartID, true).
					Preload("SupplyRequestQRAssignment").
					Preload("SupplyRequestQRAssignment.MedicalSupply").
					Find(&cartItems).Error; err == nil {

					autoReceived := 0

					for _, item := range cartItems {
						// Saltar el insumo principal (ya procesado)
						if item.SupplyRequestQRAssignment.MedicalSupplyID == supply.ID {
							continue
						}

						otherSupply := item.SupplyRequestQRAssignment.MedicalSupply

						// Solo procesar insumos que estén efectivamente en camino a bodega
						if otherSupply.Status != models.StatusEnRouteToStore ||
							!otherSupply.InTransit ||
							otherSupply.LocationType != models.SupplyLocationStore {
							continue
						}

						// === 5.1) Actualizar transferencia de devolución para este insumo ===
						var otherTransfer models.SupplyTransfer
						if err := tx.Where("qr_code = ? AND status = ?", otherSupply.QRCode, models.TransferStatusInTransit).
							First(&otherTransfer).Error; err != nil {
							if err := tx.Where("medical_supply_id = ? AND status = ? AND destination_type = ?",
								otherSupply.ID, models.TransferStatusInTransit, models.TransferLocationStore).
								Order("send_date DESC").First(&otherTransfer).Error; err != nil {
								fmt.Printf("⚠️ No se encontró transferencia en tránsito para QR %s (carrito): %v\n", otherSupply.QRCode, err)
							}
						}

						if otherTransfer.ID != 0 {
							now := time.Now()
							otherTransfer.Status = models.TransferStatusReceived
							otherTransfer.ReceiveDate = &now
							otherTransfer.ReceivedBy = &userRUT

							var tu models.User
							if err := tx.Where("rut = ?", userRUT).First(&tu).Error; err == nil {
								otherTransfer.ReceivedByName = &tu.Name
							}

							autoNotes := notes
							if autoNotes == "" {
								autoNotes = "Recepción automática de devolución desde carrito"
							} else {
								autoNotes = autoNotes + " (recepción automática desde carrito)"
							}

							if otherTransfer.Notes != "" {
								otherTransfer.Notes = otherTransfer.Notes + "\n" + autoNotes
							} else {
								otherTransfer.Notes = autoNotes
							}

							if err := tx.Save(&otherTransfer).Error; err != nil {
								fmt.Printf("⚠️ Error actualizando transferencia de devolución (carrito) para QR %s: %v\n", otherSupply.QRCode, err)
							}
						}

						// === 5.2) Actualizar inventario y lote para este insumo ===
						var otherBatch models.Batch
						if err := tx.First(&otherBatch, otherSupply.BatchID).Error; err != nil {
							fmt.Printf("⚠️ Error obteniendo lote para insumo %s: %v\n", otherSupply.QRCode, err)
							continue
						}

						var otherStore models.Store
						if err := tx.First(&otherStore, otherBatch.StoreID).Error; err != nil {
							fmt.Printf("⚠️ Error obteniendo bodega para insumo %s: %v\n", otherSupply.QRCode, err)
							continue
						}

						var otherStoreSummary models.StoreInventorySummary
						if err := tx.Where("store_id = ? AND batch_id = ?", otherStore.ID, otherBatch.ID).
							First(&otherStoreSummary).Error; err != nil {
							if errors.Is(err, gorm.ErrRecordNotFound) {
								now := time.Now()
								otherStoreSummary = models.StoreInventorySummary{
									StoreID:          otherStore.ID,
									BatchID:          otherBatch.ID,
									SupplyCode:       otherSupply.Code,
									OriginalAmount:   otherBatch.Amount,
									CurrentInStore:   1,
									TotalReturnedIn:  1,
									LastReturnInDate: &now,
								}
								if err := tx.Create(&otherStoreSummary).Error; err != nil {
									fmt.Printf("⚠️ Error creando resumen de bodega para insumo %s: %v\n", otherSupply.QRCode, err)
									continue
								}
							} else {
								fmt.Printf("⚠️ Error obteniendo resumen de bodega para insumo %s: %v\n", otherSupply.QRCode, err)
								continue
							}
						} else {
							now := time.Now()
							otherStoreSummary.CurrentInStore++
							otherStoreSummary.TotalReturnedIn++
							otherStoreSummary.LastReturnInDate = &now
							if err := tx.Save(&otherStoreSummary).Error; err != nil {
								fmt.Printf("⚠️ Error actualizando resumen de bodega para insumo %s: %v\n", otherSupply.QRCode, err)
								continue
							}
						}

						otherBatch.Amount++
						if err := tx.Save(&otherBatch).Error; err != nil {
							fmt.Printf("⚠️ Error actualizando lote para insumo %s: %v\n", otherSupply.QRCode, err)
							continue
						}

						// === 5.3) Marcar asignaciones antiguas y desactivar items de carrito para este QR ===
						var otherAssignments []models.SupplyRequestQRAssignment
						if err := tx.Where("qr_code = ? AND status NOT IN (?, ?)", otherSupply.QRCode, models.AssignmentStatusReturned, models.AssignmentStatusConsumed).
							Find(&otherAssignments).Error; err == nil {
							for _, a := range otherAssignments {
								a.Status = models.AssignmentStatusReturned
								if notes != "" {
									if a.Notes != "" {
										a.Notes = a.Notes + "\n[DEVUELTO A BODEGA] " + notes
									} else {
										a.Notes = "[DEVUELTO A BODEGA] " + notes
									}
								}
								if err := tx.Save(&a).Error; err != nil {
									fmt.Printf("⚠️ Error actualizando asignación a 'returned' para QR %s: %v\n", otherSupply.QRCode, err)
								}

								now := time.Now()
								userNamePtr := userName
								if err := tx.Model(&models.SupplyCartItem{}).
									Where("supply_request_qr_assignment_id = ? AND is_active = ?", a.ID, true).
									Updates(map[string]interface{}{
										"is_active":       false,
										"removed_at":      &now,
										"removed_by":      &userRUT,
										"removed_by_name": &userNamePtr,
										"notes":           "Insumo devuelto a bodega (recepción automática de carrito)",
									}).Error; err != nil {
									fmt.Printf("⚠️ Error desactivando items de carrito para QR %s: %v\n", otherSupply.QRCode, err)
								}
							}
						}

						// === 5.4) Actualizar insumo y registrar historial ===
						oldOtherStatus := otherSupply.Status
						otherSupply.Status = models.StatusAvailable
						otherSupply.LocationType = models.SupplyLocationStore
						otherSupply.LocationID = otherStore.ID
						otherSupply.InTransit = false
						if err := tx.Save(&otherSupply).Error; err != nil {
							fmt.Printf("⚠️ Error actualizando insumo %s: %v\n", otherSupply.QRCode, err)
							continue
						}

						autoFinalNotes := fmt.Sprintf("Llegada confirmada a bodega %s (recepción automática desde carrito).", otherStore.Name)
						if notes != "" {
							autoFinalNotes = fmt.Sprintf("Llegada confirmada a bodega %s (recepción automática desde carrito). %s", otherStore.Name, notes)
						}

						history := models.SupplyHistory{
							DateTime:        time.Now(),
							Status:          models.StatusAvailable,
							DestinationType: models.DestinationTypeStore,
							DestinationID:   otherStore.ID,
							MedicalSupplyID: otherSupply.ID,
							UserRUT:         userRUT,
							Notes:           autoFinalNotes,
						}
						if err := tx.Create(&history).Error; err != nil {
							fmt.Printf("⚠️ Error creando historial para insumo %s: %v\n", otherSupply.QRCode, err)
							continue
						}

						autoReceived++
						fmt.Printf("✅ Llegada confirmada automáticamente para insumo %s (carrito) estado: %s -> %s\n",
							otherSupply.QRCode, oldOtherStatus, otherSupply.Status)
					}

					if autoReceived > 0 {
						fmt.Printf("✅ Recepción automática de %d insumos adicionales del mismo carrito para QR base %s\n", autoReceived, qrCode)
					}
				}
			}
		}

		return nil
	})
}
