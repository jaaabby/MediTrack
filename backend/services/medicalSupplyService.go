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
			b.supplier,
			b.expiration_date,
			s.name as store_name,
			s.id as store_id,
			ms.updated_at as received_at
		FROM medical_supply ms
		JOIN supply_code sc ON ms.code = sc.code
		JOIN batch b ON ms.batch_id = b.id
		JOIN store s ON b.store_id = s.id
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
		)

		if err := rows.Scan(&id, &qrCode, &status, &supplyName, &supplyCode, &batchID, &supplier, &expirationDate, &storeName, &storeID, &receivedAt); err != nil {
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

		now := time.Now()
		oldStatus := supply.Status
		oldLocationType := supply.LocationType
		oldLocationID := supply.LocationID

		// Cambiar estado y ubicación del insumo
		supply.Status = models.StatusAvailable
		supply.LocationType = models.SupplyLocationStore
		supply.LocationID = store.ID
		supply.InTransit = false
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

		// Incrementar el resumen de bodega
		var storeSummary models.StoreInventorySummary
		if err := tx.Where("store_id = ? AND batch_id = ?", store.ID, batch.ID).
			First(&storeSummary).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Crear resumen si no existe
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

		// Crear registro en supply_history
		returnType := "manual"
		if isAutomatic {
			returnType = "automatico"
		}

		// Usar el mensaje que viene en notes (ya contiene las horas laborales correctas)
		finalNotes := fmt.Sprintf("Retorno a bodega (%s): %s", returnType, notes)

		historyEntry := models.SupplyHistory{
			DateTime:         now,
			Status:           models.StatusAvailable,
			DestinationType:  models.DestinationTypeStore,
			DestinationID:    store.ID,
			MedicalSupplyID:  supply.ID,
			UserRUT:          userRUT,
			Notes:            finalNotes,
			OriginType:       &oldLocationType,
			OriginID:         &oldLocationID,
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

// ProcessAutomaticReturns procesa automáticamente los retornos de insumos que llevan 8 horas laborales sin consumir
func (s *MedicalSupplyService) ProcessAutomaticReturns() error {
	supplies, err := s.GetSuppliesForReturn()
	if err != nil {
		return fmt.Errorf("error obteniendo insumos para retorno: %v", err)
	}

	returnedCount := 0
	errorCount := 0
	var errors []string

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
		}
	}

	if returnedCount > 0 {
		fmt.Printf("📦 Procesamiento automático completado: %d insumos retornados a bodega\n", returnedCount)
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
