package services

import (
	"fmt"
	"meditrack/mailer"
	"meditrack/models"
	"os"
	"path/filepath"
	"time"

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

	// Actualizar campos pero mantener el QR code original si existe
	supply.Code = newSupply.Code
	supply.BatchID = newSupply.BatchID
	supply.Status = newSupply.Status
	// No actualizamos supply.QRCode para mantener la trazabilidad

	if err := s.DB.Save(&supply).Error; err != nil {
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

// ConsumeSupplyByQR marca un insumo como consumido y actualiza la cantidad del lote
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

		// Actualizar el estado del insumo a consumido
		// El trigger automáticamente creará el registro en supply_history
		if err := tx.Model(&supply).Update("status", models.StatusConsumed).Error; err != nil {
			return fmt.Errorf("error actualizando estado del insumo: %v", err)
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
	// Leer email de destino desde variable de entorno o usar por defecto
	alertEmail := os.Getenv("ALERT_EMAIL")
	if alertEmail == "" {
		alertEmail = "matias.yanez@usach.cl" // Email por defecto
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
			b.supplier,
			st.name as store_name,
			EXTRACT(EPOCH FROM (NOW() - ms.updated_at))/3600 as hours_elapsed
		FROM medical_supply ms
		JOIN supply_code sc ON ms.code = sc.code
		JOIN batch b ON ms.batch_id = b.id
		JOIN store st ON b.store_id = st.id
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

// GetSuppliesForReturn obtiene insumos que deben regresar a bodega (15 días sin consumir)
func (s *MedicalSupplyService) GetSuppliesForReturn() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// Buscar insumos recepcionados hace más de 15 días que no han sido consumidos
	query := `
		SELECT 
			ms.id,
			ms.qr_code,
			ms.status,
			sc.name as supply_name,
			sc.code as supply_code,
			b.id as batch_id,
			b.supplier,
			b.expiration_date,
			s.name as store_name,
			s.id as store_id,
			sh.date_time as received_at,
			EXTRACT(EPOCH FROM (NOW() - sh.date_time))/86400 as days_elapsed
		FROM medical_supply ms
		JOIN supply_code sc ON ms.code = sc.code
		JOIN batch b ON ms.batch_id = b.id
		JOIN store s ON b.store_id = s.id
		JOIN supply_history sh ON ms.id = sh.medical_supply_id
		WHERE ms.status = ? 
		AND sh.status = ?
		AND sh.destination_type = ?
		AND sh.date_time <= NOW() - INTERVAL '15 days'
		AND NOT EXISTS (
			SELECT 1 FROM supply_history sh2 
			WHERE sh2.medical_supply_id = ms.id 
			AND sh2.date_time > sh.date_time
		)
		ORDER BY sh.date_time ASC
	`

	rows, err := s.DB.Raw(query, models.StatusReceived, models.StatusReceived, models.DestinationTypePavilion).Rows()
	if err != nil {
		return nil, fmt.Errorf("error consultando insumos para retorno: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id             int
			qrCode         string
			status         string
			supplyName     string
			supplyCode     int
			batchID        int
			supplier       string
			expirationDate time.Time
			storeName      string
			storeID        int
			receivedAt     time.Time
			daysElapsed    float64
		)

		if err := rows.Scan(&id, &qrCode, &status, &supplyName, &supplyCode, &batchID, &supplier, &expirationDate, &storeName, &storeID, &receivedAt, &daysElapsed); err != nil {
			continue
		}

		item := map[string]interface{}{
			"supply_id":       id,
			"qr_code":         qrCode,
			"status":          status,
			"supply_name":     supplyName,
			"supply_code":     supplyCode,
			"batch_id":        batchID,
			"supplier":        supplier,
			"expiration_date": expirationDate.Format("2006-01-02"),
			"store_name":      storeName,
			"store_id":        storeID,
			"received_at":     receivedAt.Format("2006-01-02 15:04:05"),
			"days_elapsed":    int(daysElapsed),
			"should_return":   daysElapsed >= 15,
		}

		results = append(results, item)
	}

	return results, nil
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
		if supply.Status == models.StatusConsumed {
			return fmt.Errorf("no se puede regresar un insumo consumido")
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

		// Cambiar estado del insumo
		oldStatus := supply.Status
		supply.Status = models.StatusEnRouteToStore
		if err := tx.Save(&supply).Error; err != nil {
			return fmt.Errorf("error actualizando estado del insumo: %v", err)
		}

		// Crear registro en supply_history
		returnType := "manual"
		if isAutomatic {
			returnType = "automatico"
		}

		finalNotes := fmt.Sprintf("Retorno a bodega (%s): %s", returnType, notes)
		if isAutomatic {
			finalNotes = fmt.Sprintf("Retorno automático a bodega después de 15 días sin consumo. %s", notes)
		}

		historyEntry := models.SupplyHistory{
			DateTime:        time.Now(),
			Status:          models.StatusEnRouteToStore,
			DestinationType: models.DestinationTypeStore,
			DestinationID:   store.ID,
			MedicalSupplyID: supply.ID,
			UserRUT:         userRUT,
			Notes:           finalNotes,
		}

		if err := tx.Create(&historyEntry).Error; err != nil {
			return fmt.Errorf("error creando historial: %v", err)
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

// ProcessAutomaticReturns procesa automáticamente los retornos de insumos que llevan 15 días sin consumir
func (s *MedicalSupplyService) ProcessAutomaticReturns() error {
	supplies, err := s.GetSuppliesForReturn()
	if err != nil {
		return fmt.Errorf("error obteniendo insumos para retorno: %v", err)
	}

	returnedCount := 0
	for _, supply := range supplies {
		supplyID := supply["supply_id"].(int)
		qrCode := supply["qr_code"].(string)
		daysElapsed := supply["days_elapsed"].(int)

		if daysElapsed >= 15 {
			// Usuario del sistema para operaciones automáticas
			systemUserRUT := "SYSTEM-AUTO"
			notes := fmt.Sprintf("Retorno automático después de %d días sin consumo", daysElapsed)

			err := s.ReturnSupplyToStore(supplyID, systemUserRUT, notes, true)
			if err != nil {
				fmt.Printf("❌ Error retornando insumo %s: %v\n", qrCode, err)
				continue
			}

			returnedCount++
			fmt.Printf("✅ Insumo %s retornado automáticamente a bodega\n", qrCode)
		}
	}

	if returnedCount > 0 {
		fmt.Printf("📦 Procesamiento automático completado: %d insumos retornados a bodega\n", returnedCount)
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

// ConfirmArrivalToStore confirma la llegada de un insumo a bodega y lo cambia a estado disponible
func (s *MedicalSupplyService) ConfirmArrivalToStore(qrCode string, userRUT string, notes string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
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

		// Cambiar estado del insumo a disponible
		oldStatus := supply.Status
		supply.Status = models.StatusAvailable
		if err := tx.Save(&supply).Error; err != nil {
			return fmt.Errorf("error actualizando estado del insumo: %v", err)
		}

		// Crear registro en supply_history
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

		// Log para depuración
		fmt.Printf("✅ Llegada confirmada para insumo %s a bodega %s (estado: %s -> %s)\n",
			supply.QRCode, store.Name, oldStatus, supply.Status)

		return nil
	})
}
