package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"meditrack/models"
	"time"

	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type QRService struct {
	DB *gorm.DB
}

func NewQRService(db *gorm.DB) *QRService {
	return &QRService{DB: db}
}

// QRGenerationResponse representa la respuesta al generar un QR
type QRGenerationResponse struct {
	QRCode    string `json:"qr_code"`
	Type      string `json:"type"`
	ImageData string `json:"image_data,omitempty"` // Base64 encoded image
	ImageURL  string `json:"image_url,omitempty"`  // URL para descargar la imagen
}

// QRConsumptionRequest representa la solicitud de consumo de un QR
type QRConsumptionRequest struct {
	QRCode          string `json:"qr_code" binding:"required"`
	UserRUT         string `json:"user_rut" binding:"required"`
	DestinationType string `json:"destination_type" binding:"required"` // "pavilion" o "store"
	DestinationID   int    `json:"destination_id" binding:"required"`
	Notes           string `json:"notes,omitempty"`
}

// QRConsumptionResponse representa la respuesta del consumo de un QR
type QRConsumptionResponse struct {
	Success            bool                   `json:"success"`
	Message            string                 `json:"message"`
	ConsumedSupply     *models.MedicalSupply  `json:"consumed_supply,omitempty"`
	UpdatedBatch       map[string]interface{} `json:"updated_batch,omitempty"`
	RemainingAmount    int                    `json:"remaining_amount"`
	ConsumptionHistory map[string]interface{} `json:"consumption_history,omitempty"`
}

// QRInfo representa la información completa de un código QR escaneado
type QRInfo struct {
	Type       string                    `json:"type"` // "batch" o "medical_supply"
	ID         int                       `json:"id"`
	QRCode     string                    `json:"qr_code"`
	BatchInfo  *models.Batch             `json:"batch_info,omitempty"`
	SupplyInfo *MedicalSupplyWithDetails `json:"supply_info,omitempty"`
	SupplyCode *models.SupplyCode        `json:"supply_code,omitempty"`
	History    []models.SupplyHistory    `json:"history,omitempty"`

	// NUEVOS CAMPOS PARA TRAZABILIDAD
	RequestAssignment *models.SupplyRequestQRAssignment `json:"request_assignment,omitempty"`
	SupplyRequest     *models.SupplyRequest             `json:"supply_request,omitempty"`
	Traceability      *QRTraceability                   `json:"traceability,omitempty"`
}

// QRTraceability contiene información completa de trazabilidad
type QRTraceability struct {
	QRCode              string                             `json:"qr_code"`
	CurrentStatus       string                             `json:"current_status"`
	IsAssignedToRequest bool                               `json:"is_assigned_to_request"`
	RequestHistory      []models.SupplyRequestQRAssignment `json:"request_history"`
	SupplyHistory       []models.SupplyHistory             `json:"supply_history"`
	CreatedDate         time.Time                          `json:"created_date"`
	LastUpdated         time.Time                          `json:"last_updated"`
	TotalMovements      int                                `json:"total_movements"`
	CurrentLocation     *LocationInfo                      `json:"current_location,omitempty"`
}

// LocationInfo representa información de ubicación actual
type LocationInfo struct {
	Type      string    `json:"type"` // "pavilion", "store", "in_transit"
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MedicalSupplyWithDetails extiende MedicalSupply con información adicional
type MedicalSupplyWithDetails struct {
	models.MedicalSupply
	BatchInfo    *models.Batch         `json:"batch_info,omitempty"`
	SupplyCode   *models.SupplyCode    `json:"supply_code,omitempty"`
	IsConsumed   bool                  `json:"is_consumed"`
	LastMovement *models.SupplyHistory `json:"last_movement,omitempty"`
	DaysToExpire *int                  `json:"days_to_expire,omitempty"`
}

// ScanQRWithTraceability escanea un QR con información completa de trazabilidad
func (s *QRService) ScanQRWithTraceability(qrCode string) (*QRInfo, error) {
	result := &QRInfo{
		QRCode: qrCode,
	}

	// Determinar el tipo de QR y obtener información básica
	isValid, qrType, err := s.ValidateQRCode(qrCode)
	if err != nil || !isValid {
		return nil, fmt.Errorf("código QR no válido: %s", qrCode)
	}

	result.Type = qrType

	if qrType == "medical_supply" {
		// Obtener información del insumo médico
		var supply models.MedicalSupply
		if err := s.DB.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
			return nil, fmt.Errorf("insumo no encontrado: %v", err)
		}

		result.ID = supply.ID

		// Obtener información detallada del insumo
		supplyDetails, err := s.GetMedicalSupplyWithDetails(supply.ID)
		if err != nil {
			return nil, fmt.Errorf("error obteniendo detalles del insumo: %v", err)
		}
		result.SupplyInfo = supplyDetails

		// Obtener información de trazabilidad
		traceability, err := s.GetQRTraceability(qrCode)
		if err != nil {
			// No fallar si no hay trazabilidad, solo registrar
			fmt.Printf("Advertencia: Error obteniendo trazabilidad para %s: %v\n", qrCode, err)
		} else {
			result.Traceability = traceability
		}

		// Obtener asignación a solicitud activa si existe
		var assignment models.SupplyRequestQRAssignment
		if err := s.DB.Where("qr_code = ? AND status NOT IN (?)", qrCode,
			[]string{models.AssignmentStatusConsumed}).
			Preload("SupplyRequest").
			Preload("SupplyRequestItem").
			Order("assigned_date DESC").
			First(&assignment).Error; err == nil {
			result.RequestAssignment = &assignment
			result.SupplyRequest = &assignment.SupplyRequest
		}

	} else if qrType == "batch" {
		// Obtener información del lote
		var batch models.Batch
		if err := s.DB.Where("qr_code = ?", qrCode).First(&batch).Error; err != nil {
			return nil, fmt.Errorf("lote no encontrado: %v", err)
		}

		result.ID = batch.ID
		result.BatchInfo = &batch

		// Para lotes, obtener información de insumos individuales
		var supplies []models.MedicalSupply
		s.DB.Where("batch_id = ?", batch.ID).Find(&supplies)

		// Obtener estadísticas del lote
		var consumedCount int64
		s.DB.Model(&models.SupplyHistory{}).
			Joins("JOIN medical_supply ON medical_supply.id = supply_history.medical_supply_id").
			Where("medical_supply.batch_id = ? AND supply_history.status = ?", batch.ID, "consumido").
			Count(&consumedCount)
	}

	// Obtener historial general del QR
	history, err := s.GetSupplyHistory(qrCode)
	if err == nil {
		result.History = history
	}

	// Obtener código de insumo
	if result.SupplyInfo != nil && result.SupplyInfo.Code > 0 {
		var supplyCode models.SupplyCode
		if err := s.DB.First(&supplyCode, result.SupplyInfo.Code).Error; err == nil {
			result.SupplyCode = &supplyCode
		}
	}

	return result, nil
}

// GetQRTraceability obtiene información completa de trazabilidad de un QR
func (s *QRService) GetQRTraceability(qrCode string) (*QRTraceability, error) {
	traceability := &QRTraceability{
		QRCode:         qrCode,
		CurrentStatus:  "unknown",
		RequestHistory: []models.SupplyRequestQRAssignment{},
		SupplyHistory:  []models.SupplyHistory{},
	}

	// Obtener información del insumo médico
	var supply models.MedicalSupply
	if err := s.DB.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
		return nil, fmt.Errorf("insumo no encontrado: %v", err)
	}

	// supply.CreatedAt does not exist, use a valid field (e.g., supply.CreatedAt if present, or set to time.Now())
	traceability.CreatedDate = time.Now() // Replace with actual creation time if available

	// Obtener historial de solicitudes
	if err := s.DB.Where("qr_code = ?", qrCode).
		Preload("SupplyRequest").
		Preload("SupplyRequestItem").
		Order("assigned_date DESC").
		Find(&traceability.RequestHistory).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo historial de solicitudes: %v", err)
	}

	// Obtener historial de movimientos del insumo
	if err := s.DB.Where("medical_supply_id = ?", supply.ID).
		Order("date_time DESC").
		Find(&traceability.SupplyHistory).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo historial de movimientos: %v", err)
	}

	// Determinar estado actual y ubicación
	traceability.IsAssignedToRequest = len(traceability.RequestHistory) > 0
	traceability.TotalMovements = len(traceability.SupplyHistory) + len(traceability.RequestHistory)

	// Determinar estado actual basado en el historial más reciente
	if len(traceability.RequestHistory) > 0 {
		latestAssignment := traceability.RequestHistory[0]
		traceability.CurrentStatus = latestAssignment.Status
		traceability.LastUpdated = latestAssignment.UpdatedAt

		// Determinar ubicación actual
		if latestAssignment.Status == models.AssignmentStatusDelivered {
			// Obtener información del pabellón
			var pavilion models.Pavilion
			if err := s.DB.First(&pavilion, latestAssignment.SupplyRequest.PavilionID).Error; err == nil {
				traceability.CurrentLocation = &LocationInfo{
					Type:      "pavilion",
					ID:        pavilion.ID,
					Name:      pavilion.Name,
					UpdatedAt: *latestAssignment.DeliveredDate,
				}
			}
		}
	} else if len(traceability.SupplyHistory) > 0 {
		latestHistory := traceability.SupplyHistory[0]
		traceability.CurrentStatus = latestHistory.Status
		traceability.LastUpdated = latestHistory.DateTime

		// Determinar ubicación basada en el historial
		if latestHistory.DestinationType == "pavilion" {
			var pavilion models.Pavilion
			if err := s.DB.First(&pavilion, latestHistory.DestinationID).Error; err == nil {
				traceability.CurrentLocation = &LocationInfo{
					Type:      "pavilion",
					ID:        pavilion.ID,
					Name:      pavilion.Name,
					UpdatedAt: latestHistory.DateTime,
				}
			}
		} else if latestHistory.DestinationType == "store" {
			var store models.Store
			if err := s.DB.First(&store, latestHistory.DestinationID).Error; err == nil {
				traceability.CurrentLocation = &LocationInfo{
					Type:      "store",
					ID:        store.ID,
					Name:      store.Name,
					UpdatedAt: latestHistory.DateTime,
				}
			}
		}
	}

	return traceability, nil
}

// GetMedicalSupplyWithDetails obtiene información detallada de un insumo médico
func (s *QRService) GetMedicalSupplyWithDetails(supplyID int) (*MedicalSupplyWithDetails, error) {
	var supply models.MedicalSupply
	if err := s.DB.First(&supply, supplyID).Error; err != nil {
		return nil, fmt.Errorf("insumo no encontrado: %v", err)
	}

	result := &MedicalSupplyWithDetails{
		MedicalSupply: supply,
	}

	// Obtener información del lote
	var batch models.Batch
	if err := s.DB.First(&batch, supply.BatchID).Error; err == nil {
		result.BatchInfo = &batch

		// Calcular días hasta vencimiento
		daysToExpire := int(time.Until(batch.ExpirationDate).Hours() / 24)
		result.DaysToExpire = &daysToExpire
	}

	// Obtener información del código de insumo
	var supplyCode models.SupplyCode
	if err := s.DB.First(&supplyCode, supply.Code).Error; err == nil {
		result.SupplyCode = &supplyCode
	}

	// Verificar si está consumido
	var consumedHistory models.SupplyHistory
	if err := s.DB.Where("medical_supply_id = ? AND status = ?",
		supply.ID, "consumido").First(&consumedHistory).Error; err == nil {
		result.IsConsumed = true
	}

	// Obtener último movimiento
	var lastMovement models.SupplyHistory
	if err := s.DB.Where("medical_supply_id = ?", supply.ID).
		Order("date_time DESC").
		First(&lastMovement).Error; err == nil {
		result.LastMovement = &lastMovement
	}

	return result, nil
}

// GetQRRequestHistory obtiene el historial de solicitudes para un QR específico
func (s *QRService) GetQRRequestHistory(qrCode string) ([]models.SupplyRequestQRAssignment, error) {
	var assignments []models.SupplyRequestQRAssignment

	if err := s.DB.Where("qr_code = ?", qrCode).
		Preload("SupplyRequest").
		Preload("SupplyRequestItem").
		Preload("MedicalSupply").
		Order("assigned_date DESC").
		Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo historial de solicitudes: %v", err)
	}

	return assignments, nil
}

// IsQRAvailableForRequest verifica si un QR está disponible para asignar a una solicitud
func (s *QRService) IsQRAvailableForRequest(qrCode string) (bool, string, error) {
	// Verificar que el QR existe
	var supply models.MedicalSupply
	if err := s.DB.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
		return false, "QR no encontrado", err
	}

	// Verificar que no esté consumido
	var consumedHistory models.SupplyHistory
	if err := s.DB.Where("medical_supply_id = ? AND status = ?",
		supply.ID, "consumido").First(&consumedHistory).Error; err == nil {
		return false, "El insumo ya fue consumido", nil
	}

	// Verificar que no esté asignado a otra solicitud activa
	var activeAssignment models.SupplyRequestQRAssignment
	if err := s.DB.Where("qr_code = ? AND status NOT IN (?)", qrCode,
		[]string{models.AssignmentStatusConsumed, models.AssignmentStatusReturned}).
		First(&activeAssignment).Error; err == nil {
		return false, "El QR ya está asignado a otra solicitud activa", nil
	}

	// Verificar que el lote no haya expirado
	var batch models.Batch
	if err := s.DB.First(&batch, supply.BatchID).Error; err == nil {
		if batch.ExpirationDate.Before(time.Now()) {
			return false, "El lote ha expirado", nil
		}
	}

	return true, "Disponible", nil
}

// UpdateQRAssignmentStatus actualiza el estado de una asignación QR
func (s *QRService) UpdateQRAssignmentStatus(qrCode, newStatus, updatedBy, updatedByName string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var assignment models.SupplyRequestQRAssignment
		if err := tx.Where("qr_code = ?", qrCode).
			Order("assigned_date DESC").
			First(&assignment).Error; err != nil {
			return fmt.Errorf("asignación no encontrada: %v", err)
		}

		updates := map[string]interface{}{
			"status":     newStatus,
			"updated_at": time.Now(),
		}

		// Actualizar campos específicos según el nuevo estado
		if newStatus == models.AssignmentStatusDelivered {
			now := time.Now()
			updates["delivered_date"] = now
			updates["delivered_by"] = updatedBy
			updates["delivered_by_name"] = updatedByName
		}

		return tx.Model(&assignment).Updates(updates).Error
	})
}

// GetAvailableQRsForSupplyCode obtiene QRs disponibles para un código de insumo específico
func (s *QRService) GetAvailableQRsForSupplyCode(supplyCode int, limit int) ([]models.MedicalSupply, error) {
	var supplies []models.MedicalSupply

	// Subconsulta para obtener IDs de insumos consumidos
	consumedSubQuery := s.DB.Model(&models.SupplyHistory{}).
		Select("medical_supply_id").
		Where("status = ?", "consumido")

	// Subconsulta para obtener QRs asignados activamente
	assignedSubQuery := s.DB.Model(&models.SupplyRequestQRAssignment{}).
		Select("qr_code").
		Where("status NOT IN (?)", []string{models.AssignmentStatusConsumed, models.AssignmentStatusReturned})

	query := s.DB.Where("code = ?", supplyCode).
		Where("id NOT IN (?)", consumedSubQuery).
		Where("qr_code NOT IN (?)", assignedSubQuery).
		Preload("Batch").
		Order("created_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&supplies).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo QRs disponibles: %v", err)
	}

	// Filtrar por fecha de vencimiento
	var availableSupplies []models.MedicalSupply
	now := time.Now()

	for _, supply := range supplies {
		// Verificar que el lote no haya expirado (se asume que hay una relación Batch cargada)
		var batch models.Batch
		if err := s.DB.First(&batch, supply.BatchID).Error; err == nil {
			if batch.ExpirationDate.After(now) {
				availableSupplies = append(availableSupplies, supply)
			}
		}
	}

	return availableSupplies, nil
}

// MÉTODOS ORIGINALES ACTUALIZADOS

// ScanQRCode escanea un código QR (método original actualizado)
func (s *QRService) ScanQRCode(qrCode string) (map[string]interface{}, error) {
	// Usar el nuevo método con trazabilidad y convertir al formato original
	qrInfo, err := s.ScanQRWithTraceability(qrCode)
	if err != nil {
		return nil, err
	}

	// Convertir a formato compatible con el frontend existente
	result := map[string]interface{}{
		"qr_code": qrInfo.QRCode,
		"type":    qrInfo.Type,
		"id":      qrInfo.ID,
	}

	if qrInfo.SupplyInfo != nil {
		result["supply_info"] = qrInfo.SupplyInfo
		result["is_consumed"] = qrInfo.SupplyInfo.IsConsumed
		result["available_for_use"] = !qrInfo.SupplyInfo.IsConsumed
	}

	if qrInfo.BatchInfo != nil {
		result["batch_info"] = qrInfo.BatchInfo
	}

	if qrInfo.SupplyCode != nil {
		result["supply_code"] = qrInfo.SupplyCode
	}

	if qrInfo.History != nil {
		result["history"] = qrInfo.History
	}

	// Agregar información de trazabilidad si está disponible
	if qrInfo.Traceability != nil {
		result["traceability"] = qrInfo.Traceability
	}

	if qrInfo.RequestAssignment != nil {
		result["request_assignment"] = qrInfo.RequestAssignment
		result["assigned_to_request"] = true
	} else {
		result["assigned_to_request"] = false
	}

	return result, nil
}

// Mantener métodos originales existentes...
// GenerateBatchQR, GenerateSupplyQR, ConsumeSupplyByQR, etc.
// (El resto de métodos permanecen igual que en tu implementación original)

// generateSecureToken genera un token seguro
func generateSecureToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateUniqueQRCode genera un código QR único
func (s *QRService) generateUniqueQRCode(prefix string, id int) (string, error) {
	token, err := generateSecureToken()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%d_%s", prefix, id, token[:8]), nil
}

// GenerateBatchQR genera un código QR para un lote
func (s *QRService) GenerateBatchQR() (*QRGenerationResponse, error) {
	// Crear un lote temporal para generar el QR
	// (Este método necesitaría ser ajustado según tu lógica de negocio)
	qrCode, err := s.generateUniqueQRCode("BATCH", 1)
	if err != nil {
		return nil, err
	}

	// Generar imagen QR
	qrImage, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	return &QRGenerationResponse{
		QRCode:    qrCode,
		Type:      "batch",
		ImageData: base64.StdEncoding.EncodeToString(qrImage),
	}, nil
}

// GenerateSupplyQR genera un código QR para un insumo médico
func (s *QRService) GenerateSupplyQR() (*QRGenerationResponse, error) {
	// Crear un insumo temporal para generar el QR
	// (Este método necesitaría ser ajustado según tu lógica de negocio)
	qrCode, err := s.generateUniqueQRCode("SUPPLY", 1)
	if err != nil {
		return nil, err
	}

	// Generar imagen QR
	qrImage, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	return &QRGenerationResponse{
		QRCode:    qrCode,
		Type:      "medical_supply",
		ImageData: base64.StdEncoding.EncodeToString(qrImage),
	}, nil
}

// ValidateQRCode valida si un código QR existe y retorna su tipo
func (s *QRService) ValidateQRCode(qrCode string) (bool, string, error) {
	// Buscar en batch
	var batchCount int64
	s.DB.Model(&models.Batch{}).Where("qr_code = ?", qrCode).Count(&batchCount)
	if batchCount > 0 {
		return true, "batch", nil
	}

	// Buscar en medical_supply
	var supplyCount int64
	s.DB.Model(&models.MedicalSupply{}).Where("qr_code = ?", qrCode).Count(&supplyCount)
	if supplyCount > 0 {
		return true, "medical_supply", nil
	}

	return false, "", fmt.Errorf("código QR no válido: %s", qrCode)
}

// GetSupplyHistory obtiene el historial completo de un insumo por su código QR
func (s *QRService) GetSupplyHistory(qrCode string) ([]models.SupplyHistory, error) {
	var supply models.MedicalSupply
	if err := s.DB.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
		return nil, fmt.Errorf("insumo no encontrado con QR: %s", qrCode)
	}

	var history []models.SupplyHistory
	if err := s.DB.Where("medical_supply_id = ?", supply.ID).Order("date_time DESC").Find(&history).Error; err != nil {
		return nil, err
	}

	return history, nil
}

// ConsumeSupplyByQR procesa el consumo de un insumo (método original mantenido)
func (s *QRService) ConsumeSupplyByQR(request QRConsumptionRequest) (*QRConsumptionResponse, error) {
	var supply models.MedicalSupply
	var response QRConsumptionResponse

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Buscar el insumo por QR
		if err := tx.Where("qr_code = ?", request.QRCode).First(&supply).Error; err != nil {
			return fmt.Errorf("insumo no encontrado con QR %s: %v", request.QRCode, err)
		}

		// Verificar que el insumo no haya sido consumido previamente
		var existingHistory models.SupplyHistory
		err := tx.Where("medical_supply_id = ? AND status = ?", supply.ID, "consumido").First(&existingHistory).Error
		if err == nil {
			return fmt.Errorf("el insumo con QR %s ya ha sido consumido el %s", request.QRCode, existingHistory.DateTime.Format("2006-01-02 15:04:05"))
		}

		// Obtener información del lote antes de la actualización
		var batch models.Batch
		if err := tx.Where("id = ?", supply.BatchID).First(&batch).Error; err != nil {
			return fmt.Errorf("lote no encontrado: %v", err)
		}

		// Verificar que hay stock disponible
		if batch.Amount <= 0 {
			return fmt.Errorf("no hay stock disponible en el lote %d", batch.ID)
		}

		// Crear historial de consumo
		history := models.SupplyHistory{
			DateTime:        time.Now(),
			Status:          "consumido",
			DestinationType: request.DestinationType,
			DestinationID:   request.DestinationID,
			MedicalSupplyID: supply.ID,
			UserRUT:         request.UserRUT,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("error creando historial de consumo: %v", err)
		}

		// Actualizar cantidad del lote (restar 1)
		newAmount := batch.Amount - 1
		if err := tx.Model(&batch).Update("amount", newAmount).Error; err != nil {
			return fmt.Errorf("error actualizando cantidad del lote: %v", err)
		}

		// Si hay una asignación activa, actualizarla a consumido
		var assignment models.SupplyRequestQRAssignment
		if err := tx.Where("qr_code = ? AND status NOT IN (?)", request.QRCode,
			[]string{models.AssignmentStatusConsumed}).
			First(&assignment).Error; err == nil {
			tx.Model(&assignment).Updates(map[string]interface{}{
				"status":     models.AssignmentStatusConsumed,
				"updated_at": time.Now(),
			})
		}

		// Preparar respuesta
		response.Success = true
		response.Message = fmt.Sprintf("Insumo consumido exitosamente. Lote actualizado: %d -> %d unidades", batch.Amount, newAmount)
		response.ConsumedSupply = &supply
		response.RemainingAmount = newAmount

		// Información actualizada del lote
		updatedBatchInfo := map[string]interface{}{
			"batch_id":        batch.ID,
			"previous_amount": batch.Amount,
			"new_amount":      newAmount,
			"supplier":        batch.Supplier,
			"expiration_date": batch.ExpirationDate,
		}
		response.UpdatedBatch = updatedBatchInfo

		// Información del historial de consumo
		consumptionInfo := map[string]interface{}{
			"consumption_id":   history.ID,
			"consumed_at":      history.DateTime,
			"consumed_by":      request.UserRUT,
			"destination_type": request.DestinationType,
			"destination_id":   request.DestinationID,
		}
		response.ConsumptionHistory = consumptionInfo

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response, nil
}
