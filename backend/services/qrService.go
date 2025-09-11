package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"meditrack/models"
	"net"
	"strings"
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

// ScanContext contiene información contextual para el escaneo
type ScanContext struct {
	UserRUT         *string             `json:"user_rut,omitempty"`
	UserName        *string             `json:"user_name,omitempty"`
	PavilionID      *int                `json:"pavilion_id,omitempty"`
	MedicalCenterID *int                `json:"medical_center_id,omitempty"`
	ScanSource      string              `json:"scan_source"` // web, mobile, api
	ScanPurpose     *string             `json:"scan_purpose,omitempty"`
	UserAgent       *string             `json:"user_agent,omitempty"`
	IPAddress       *net.IP             `json:"ip_address,omitempty"`
	DeviceInfo      *models.DeviceInfo  `json:"device_info,omitempty"`
	BrowserInfo     *models.BrowserInfo `json:"browser_info,omitempty"`
	SessionID       *string             `json:"session_id,omitempty"`
	RequestID       *string             `json:"request_id,omitempty"`
	Notes           *string             `json:"notes,omitempty"`
}

// QRInfo representa la información completa de un código QR escaneado
type QRInfo struct {
	Type              string                            `json:"type"` // "batch" o "medical_supply"
	ID                int                               `json:"id"`
	QRCode            string                            `json:"qr_code"`
	BatchInfo         *models.Batch                     `json:"batch_info,omitempty"`
	SupplyInfo        *MedicalSupplyWithDetails         `json:"supply_info,omitempty"`
	SupplyCode        *models.SupplyCode                `json:"supply_code,omitempty"`
	History           []models.SupplyHistory            `json:"history,omitempty"`
	RequestAssignment *models.SupplyRequestQRAssignment `json:"request_assignment,omitempty"`
	SupplyRequest     *models.SupplyRequest             `json:"supply_request,omitempty"`
	Traceability      *QRTraceability                   `json:"traceability,omitempty"`
	ScanEvents        []models.QRScanEvent              `json:"scan_events,omitempty"`
	ScanStatistics    *models.QRScanStatistics          `json:"scan_statistics,omitempty"`
}

// QRTraceability contiene información completa de trazabilidad
type QRTraceability struct {
	QRCode              string                                `json:"qr_code"`
	CurrentStatus       string                                `json:"current_status"`
	IsAssignedToRequest bool                                  `json:"is_assigned_to_request"`
	RequestHistory      []models.SupplyRequestQRAssignment    `json:"request_history"`
	SupplyHistory       []models.SupplyHistoryWithDestination `json:"supply_history"`
	ScanHistory         []models.QRCompleteTraceability       `json:"scan_history"`
	CreatedDate         time.Time                             `json:"created_date"`
	LastUpdated         time.Time                             `json:"last_updated"`
	TotalMovements      int                                   `json:"total_movements"`
	CurrentLocation     *LocationInfo                         `json:"current_location,omitempty"`
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

// =============================================
// MÉTODO PRINCIPAL CON REGISTRO AUTOMÁTICO
// =============================================

// ScanQRWithAutoLogging escanea un QR y registra automáticamente el evento
func (s *QRService) ScanQRWithAutoLogging(qrCode string, context *ScanContext) (*QRInfo, error) {
	// Primero obtener la información del QR
	qrInfo, err := s.ScanQRWithTraceability(qrCode)
	if err != nil {
		// Registrar el escaneo fallido
		s.logScanEvent(qrCode, context, nil, models.ScanResultError, err.Error())
		return nil, err
	}

	// Registrar el escaneo exitoso
	scanEvent, logErr := s.logScanEvent(qrCode, context, qrInfo, models.ScanResultSuccess, "")
	if logErr != nil {
		// Log el error pero no fallar el escaneo
		fmt.Printf("Error logging scan event: %v\n", logErr)
	}

	// Agregar el evento de escaneo a la respuesta
	if scanEvent != nil {
		qrInfo.ScanEvents = append(qrInfo.ScanEvents, *scanEvent)
	}

	// Obtener estadísticas actualizadas
	stats, _ := s.GetQRScanStatistics(qrCode)
	qrInfo.ScanStatistics = stats

	// Log para debug del frontend
	fmt.Printf("DEBUG - QR Response from ScanQRWithAutoLogging for %s:\n", qrInfo.QRCode)
	fmt.Printf("  - SupplyInfo exists: %v\n", qrInfo.SupplyInfo != nil)
	if qrInfo.SupplyInfo != nil {
		fmt.Printf("  - BatchInfo exists: %v\n", qrInfo.SupplyInfo.BatchInfo != nil)
		if qrInfo.SupplyInfo.BatchInfo != nil {
			fmt.Printf("  - Batch ID: %d, Amount: %d\n", qrInfo.SupplyInfo.BatchInfo.ID, qrInfo.SupplyInfo.BatchInfo.Amount)
		}
	}

	return qrInfo, nil
}

// logScanEvent registra automáticamente un evento de escaneo
func (s *QRService) logScanEvent(qrCode string, context *ScanContext, qrInfo *QRInfo, result string, errorMsg string) (*models.QRScanEvent, error) {
	event := &models.QRScanEvent{
		QRCode:     qrCode,
		ScannedAt:  time.Now(),
		ScanSource: models.ScanSourceWeb, // default
		ScanResult: result,
	}

	// Aplicar contexto si está disponible
	if context != nil {
		event.ScannedByRUT = context.UserRUT
		event.ScannedByName = context.UserName
		event.PavilionID = context.PavilionID
		event.MedicalCenterID = context.MedicalCenterID
		event.ScanPurpose = context.ScanPurpose
		event.UserAgent = context.UserAgent
		event.IPAddress = context.IPAddress
		event.DeviceInfo = context.DeviceInfo
		event.BrowserInfo = context.BrowserInfo
		event.SessionID = context.SessionID
		event.RequestID = context.RequestID
		event.Notes = context.Notes

		if context.ScanSource != "" {
			event.ScanSource = context.ScanSource
		}
	}

	// Si hay error, registrarlo
	if errorMsg != "" {
		event.ErrorMessage = &errorMsg
	}

	// Si el escaneo fue exitoso, agregar información del QR
	if qrInfo != nil && result == models.ScanResultSuccess {
		// Mapear el tipo interno al formato de la base de datos
		dbQRType := mapQRTypeToDatabase(qrInfo.Type)
		event.QRType = &dbQRType

		if qrInfo.SupplyInfo != nil {
			event.SupplyID = &qrInfo.SupplyInfo.ID
			event.SupplyCode = &qrInfo.SupplyInfo.Code
			if qrInfo.SupplyCode != nil {
				event.SupplyName = &qrInfo.SupplyCode.Name
			}
			event.BatchID = &qrInfo.SupplyInfo.BatchID
		}

		if qrInfo.BatchInfo != nil {
			if event.BatchID == nil {
				event.BatchID = &qrInfo.BatchInfo.ID
			}
			event.BatchSupplier = &qrInfo.BatchInfo.Supplier
		}

		// Determinar status actual
		currentStatus := "available"
		if qrInfo.SupplyInfo != nil {
			// Verificar si está consumido usando el campo status
			if qrInfo.SupplyInfo.Status == models.StatusConsumed {
				currentStatus = "consumed"
			}
		}
		event.CurrentStatus = &currentStatus

		// Determinar ubicación actual
		if context != nil && context.PavilionID != nil {
			var pavilion models.Pavilion
			if err := s.DB.First(&pavilion, *context.PavilionID).Error; err == nil {
				location := pavilion.Name
				event.CurrentLocation = &location
				event.PavilionName = &pavilion.Name
			}
		}

		if context != nil && context.MedicalCenterID != nil {
			var medicalCenter models.MedicalCenter
			if err := s.DB.First(&medicalCenter, *context.MedicalCenterID).Error; err == nil {
				event.MedicalCenterName = &medicalCenter.Name
			}
		}

		// Determinar tipo de movimiento
		movementType := models.MovementTypeScanOnly
		if context != nil && context.ScanPurpose != nil {
			switch *context.ScanPurpose {
			case models.ScanPurposeConsume:
				movementType = models.MovementTypeStatusChange
			case models.ScanPurposeAssign:
				movementType = models.MovementTypeLocationChange
			}
		}
		event.MovementType = &movementType
	}

	// Guardar el evento en la base de datos
	if err := s.DB.Create(event).Error; err != nil {
		return nil, fmt.Errorf("error guardando evento de escaneo: %v", err)
	}

	return event, nil
}

// =============================================
// MÉTODOS DE TRAZABILIDAD AVANZADA
// =============================================

// GetCompleteTraceability obtiene la trazabilidad completa de un QR incluyendo escaneos
func (s *QRService) GetCompleteTraceability(qrCode string) (*QRTraceability, error) {
	traceability := &QRTraceability{
		QRCode:         qrCode,
		CurrentStatus:  "disponible", // Estado por defecto más descriptivo
		RequestHistory: []models.SupplyRequestQRAssignment{},
		SupplyHistory:  []models.SupplyHistoryWithDestination{},
		ScanHistory:    []models.QRCompleteTraceability{},
	}

	// Obtener información del insumo médico
	var supply models.MedicalSupply
	if err := s.DB.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
		return nil, fmt.Errorf("insumo no encontrado: %v", err)
	}

	traceability.CreatedDate = time.Now() // Replace with actual creation time if available

	// Obtener historial de solicitudes
	if err := s.DB.Where("qr_code = ?", qrCode).
		Preload("SupplyRequest").
		Preload("SupplyRequestItem").
		Order("assigned_date DESC").
		Find(&traceability.RequestHistory).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo historial de solicitudes: %v", err)
	}

	// Obtener historial de movimientos del insumo con información de destino
	if err := s.DB.Table("supply_history sh").
		Select(`sh.*, 
			CASE 
				WHEN sh.destination_type = 'pavilion' THEN p.name
				WHEN sh.destination_type = 'store' THEN st.name
				ELSE NULL
			END as destination_name,
			CASE
				WHEN sh.destination_type = 'pavilion' THEN mc.name
				ELSE NULL
			END as medical_center_name,
			u.name as user_name`).
		Joins("LEFT JOIN pavilion p ON sh.destination_type = 'pavilion' AND sh.destination_id = p.id").
		Joins("LEFT JOIN store st ON sh.destination_type = 'store' AND sh.destination_id = st.id").
		Joins("LEFT JOIN medical_center mc ON p.medical_center_id = mc.id").
		Joins("LEFT JOIN \"user\" u ON sh.user_rut = u.rut").
		Where("sh.medical_supply_id = ?", supply.ID).
		Order("sh.date_time DESC").
		Find(&traceability.SupplyHistory).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo historial de movimientos: %v", err)
	}

	// Obtener historial de escaneos
	if err := s.DB.Where("qr_code = ?", qrCode).
		Order("scanned_at DESC").
		Find(&traceability.ScanHistory).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo historial de escaneos: %v", err)
	}

	// Determinar estado actual y ubicación
	traceability.IsAssignedToRequest = len(traceability.RequestHistory) > 0
	traceability.TotalMovements = len(traceability.SupplyHistory) + len(traceability.RequestHistory) + len(traceability.ScanHistory)

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

// GetQRScanStatistics obtiene estadísticas de escaneo para un QR específico
func (s *QRService) GetQRScanStatistics(qrCode string) (*models.QRScanStatistics, error) {
	var stats models.QRScanStatistics

	err := s.DB.Where("qr_code = ?", qrCode).First(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("error obteniendo estadísticas de escaneo: %v", err)
	}

	return &stats, nil
}

// GetScanEventHistory obtiene el historial completo de eventos de escaneo
func (s *QRService) GetScanEventHistory(qrCode string, limit int) ([]models.QRCompleteTraceability, error) {
	var history []models.QRCompleteTraceability

	query := s.DB.Where("qr_code = ?", qrCode).Order("scanned_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&history).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo historial de escaneos: %v", err)
	}

	return history, nil
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
		CurrentStatus:  "disponible", // Estado por defecto más descriptivo
		RequestHistory: []models.SupplyRequestQRAssignment{},
		SupplyHistory:  []models.SupplyHistoryWithDestination{},
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

	// Obtener historial de movimientos del insumo con información de destino
	if err := s.DB.Table("supply_history sh").
		Select(`sh.*, 
			CASE 
				WHEN sh.destination_type = 'pavilion' THEN p.name
				WHEN sh.destination_type = 'store' THEN st.name
				ELSE NULL
			END as destination_name,
			CASE
				WHEN sh.destination_type = 'pavilion' THEN mc.name
				ELSE NULL
			END as medical_center_name,
			u.name as user_name`).
		Joins("LEFT JOIN pavilion p ON sh.destination_type = 'pavilion' AND sh.destination_id = p.id").
		Joins("LEFT JOIN store st ON sh.destination_type = 'store' AND sh.destination_id = st.id").
		Joins("LEFT JOIN medical_center mc ON p.medical_center_id = mc.id").
		Joins("LEFT JOIN \"user\" u ON sh.user_rut = u.rut").
		Where("sh.medical_supply_id = ?", supply.ID).
		Order("sh.date_time DESC").
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

	// Verificar si está consumido usando el campo status
	result.IsConsumed = supply.IsConsumed()

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
	if supply.Status == models.StatusConsumed {
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

	// Subconsulta para obtener QRs asignados activamente
	assignedSubQuery := s.DB.Model(&models.SupplyRequestQRAssignment{}).
		Select("qr_code").
		Where("status NOT IN (?)", []string{models.AssignmentStatusConsumed, models.AssignmentStatusReturned})

	query := s.DB.Where("code = ?", supplyCode).
		Where("status != ?", models.StatusConsumed).
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

// =============================================
// MÉTODOS ORIGINALES ACTUALIZADOS
// =============================================

// ScanQRCode escanea un código QR (método original actualizado para compatibilidad)
func (s *QRService) ScanQRCode(qrCode string) (map[string]interface{}, error) {
	// Crear contexto básico para el escaneo
	context := &ScanContext{
		ScanSource:  models.ScanSourceWeb,
		ScanPurpose: stringPtr(models.ScanPurposeLookup),
	}

	// Usar el nuevo método con logging automático
	qrInfo, err := s.ScanQRWithAutoLogging(qrCode, context)
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
		// Crear una copia del SupplyInfo con la información del batch anidada
		supplyInfoMap := map[string]interface{}{
			"ID":           qrInfo.SupplyInfo.ID,
			"Code":         qrInfo.SupplyInfo.Code,
			"BatchID":      qrInfo.SupplyInfo.BatchID,
			"QRCode":       qrInfo.SupplyInfo.QRCode,
			"IsConsumed":   qrInfo.SupplyInfo.IsConsumed,
			"LastMovement": qrInfo.SupplyInfo.LastMovement,
			"DaysToExpire": qrInfo.SupplyInfo.DaysToExpire,
		}

		// Agregar nombre del insumo directamente
		if qrInfo.SupplyInfo.SupplyCode != nil {
			supplyInfoMap["name"] = qrInfo.SupplyInfo.SupplyCode.Name
			supplyInfoMap["supply_code_name"] = qrInfo.SupplyInfo.SupplyCode.Name
		}

		// Agregar información del batch dentro de supply_info
		if qrInfo.SupplyInfo.BatchInfo != nil {
			supplyInfoMap["batch"] = qrInfo.SupplyInfo.BatchInfo
		}

		// Agregar información del código de insumo
		if qrInfo.SupplyInfo.SupplyCode != nil {
			supplyInfoMap["SupplyCode"] = qrInfo.SupplyInfo.SupplyCode
		}

		result["supply_info"] = supplyInfoMap
		result["is_consumed"] = qrInfo.SupplyInfo.IsConsumed
		result["available_for_use"] = !qrInfo.SupplyInfo.IsConsumed
		result["can_consume"] = !qrInfo.SupplyInfo.IsConsumed
	}

	if qrInfo.BatchInfo != nil {
		result["batch_info"] = qrInfo.BatchInfo
		result["batch_status"] = "active"
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

	// Agregar información de escaneos
	result["scan_events"] = qrInfo.ScanEvents
	result["scan_statistics"] = qrInfo.ScanStatistics

	// Log para debug del frontend
	fmt.Printf("DEBUG - QR Response for %s:\n", qrInfo.QRCode)
	fmt.Printf("  - SupplyInfo exists: %v\n", qrInfo.SupplyInfo != nil)
	if qrInfo.SupplyInfo != nil {
		fmt.Printf("  - BatchInfo exists: %v\n", qrInfo.SupplyInfo.BatchInfo != nil)
		if qrInfo.SupplyInfo.BatchInfo != nil {
			fmt.Printf("  - Batch ID: %d, Amount: %d\n", qrInfo.SupplyInfo.BatchInfo.ID, qrInfo.SupplyInfo.BatchInfo.Amount)
		}
	}

	return result, nil
}

// TransferSupplyByQR transfiere un insumo individual por su código QR
func (s *QRService) TransferSupplyByQR(qrCode, userRUT, receiverRUT, destinationType string, destinationID int, notes string) (map[string]interface{}, error) {
	var supply models.MedicalSupply
	var result map[string]interface{}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Buscar el insumo por QR
		if err := tx.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
			return fmt.Errorf("insumo no encontrado con QR %s: %v", qrCode, err)
		}

		// Verificar que el insumo no haya sido consumido previamente
		if supply.Status == models.StatusConsumed {
			return fmt.Errorf("el insumo con QR %s ya ha sido consumido y no puede ser transferido", qrCode)
		}

		// Verificar que el insumo esté disponible para transferencia
		if supply.Status != models.StatusAvailable {
			if supply.Status == models.StatusReceived {
				return fmt.Errorf("el insumo tiene estado 'recepcionado' y solo puede ser consumido, no transferido")
			}
			return fmt.Errorf("el insumo tiene estado '%s' y no está disponible para transferencia", supply.Status)
		}

		// Determinar el estado de tránsito según el tipo de destino
		var transitStatus string
		switch destinationType {
		case "pavilion":
			transitStatus = models.StatusEnRouteToPavilion
		case "store", "warehouse":
			transitStatus = models.StatusEnRouteToStore
		default:
			transitStatus = models.StatusEnRouteToPavilion // Default
		}

		// Actualizar el estado del insumo a "en tránsito"
		if err := tx.Model(&supply).Update("status", transitStatus).Error; err != nil {
			return fmt.Errorf("error actualizando estado del insumo: %v", err)
		}

		// Crear entrada en el historial de suministros para la transferencia
		history := models.SupplyHistory{
			DateTime:        time.Now(),
			Status:          transitStatus,
			DestinationType: destinationType,
			DestinationID:   destinationID,
			MedicalSupplyID: supply.ID,
			UserRUT:         userRUT,
			Notes:           notes,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("error creando historial de transferencia: %v", err)
		}

		// Crear resultado de la transferencia
		result = map[string]interface{}{
			"success":            true,
			"supply_id":          supply.ID,
			"qr_code":            qrCode,
			"old_status":         models.StatusAvailable,
			"new_status":         transitStatus,
			"destination_type":   destinationType,
			"destination_id":     destinationID,
			"user_rut":           userRUT,
			"receiver_rut":       receiverRUT,
			"transfer_timestamp": time.Now(),
			"status_change": map[string]string{
				"from": models.StatusAvailable,
				"to":   transitStatus,
			},
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// ReceiveSupplyByQR recepciona un insumo que está en estado "en_camino_a_pabellon"
func (s *QRService) ReceiveSupplyByQR(qrCode, userRUT, destinationType string, destinationID int, notes string) (map[string]interface{}, error) {
	var supply models.MedicalSupply
	var result map[string]interface{}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Buscar el insumo por QR
		if err := tx.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
			return fmt.Errorf("insumo no encontrado con QR %s: %v", qrCode, err)
		}

		// Verificar que el insumo esté en estado "en_camino_a_pabellon"
		if supply.Status != "en_camino_a_pabellon" {
			return fmt.Errorf("el insumo debe estar en estado 'en_camino_a_pabellon' para ser recepcionado, estado actual: %s", supply.Status)
		}

		// Crear entrada en el historial
		history := models.SupplyHistory{
			DateTime:        time.Now(),
			Status:          "recepcionado",
			DestinationType: destinationType,
			DestinationID:   destinationID,
			MedicalSupplyID: supply.ID,
			UserRUT:         userRUT,
			Notes:           notes,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("error al crear historial de recepción: %v", err)
		}

		// Actualizar el estado del insumo
		supply.Status = "recepcionado"
		supply.UpdatedAt = time.Now()

		if err := tx.Save(&supply).Error; err != nil {
			return fmt.Errorf("error al actualizar estado del insumo: %v", err)
		}

		// Preparar respuesta
		result = map[string]interface{}{
			"supply_id":        supply.ID,
			"qr_code":          supply.QRCode,
			"previous_status":  "en_camino_a_pabellon",
			"new_status":       "recepcionado",
			"received_at":      time.Now(),
			"user_rut":         userRUT,
			"destination_type": destinationType,
			"destination_id":   destinationID,
			"notes":            notes,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
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
		if supply.Status == models.StatusConsumed {
			return fmt.Errorf("el insumo con QR %s ya ha sido consumido", request.QRCode)
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
			Status:          models.StatusConsumed,
			DestinationType: request.DestinationType,
			DestinationID:   request.DestinationID,
			MedicalSupplyID: supply.ID,
			UserRUT:         request.UserRUT,
			Notes:           request.Notes,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("error creando historial de consumo: %v", err)
		}

		// Actualizar el estado del insumo a consumido
		if err := tx.Model(&supply).Update("status", models.StatusConsumed).Error; err != nil {
			return fmt.Errorf("error actualizando estado del insumo: %v", err)
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

// mapQRTypeToDatabase mapea tipos internos a tipos de base de datos
func mapQRTypeToDatabase(internalType string) string {
	switch internalType {
	case "medical_supply":
		return "SUPPLY"
	case "batch":
		return "BATCH"
	default:
		return internalType
	}
}

// mapQRTypeFromDatabase mapea tipos de base de datos a tipos internos
func mapQRTypeFromDatabase(dbType string) string {
	switch dbType {
	case "SUPPLY":
		return "medical_supply"
	case "BATCH":
		return "batch"
	default:
		return strings.ToLower(dbType)
	}
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

// generateSecureToken genera un token seguro
func generateSecureToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateUniqueQRCode genera un código QR único garantizando unicidad en la base de datos
func (s *QRService) generateUniqueQRCode(prefix string, id int) (string, error) {
	maxAttempts := 10

	for attempt := 0; attempt < maxAttempts; attempt++ {
		token, err := generateSecureToken()
		if err != nil {
			return "", err
		}

		qrCode := fmt.Sprintf("%s_%d_%s", prefix, id, token[:8])

		// Verificar unicidad en batch
		var batchCount int64
		s.DB.Model(&models.Batch{}).Where("qr_code = ?", qrCode).Count(&batchCount)

		// Verificar unicidad en medical_supply
		var supplyCount int64
		s.DB.Model(&models.MedicalSupply{}).Where("qr_code = ?", qrCode).Count(&supplyCount)

		// Si no existe en ninguna tabla, el QR es único
		if batchCount == 0 && supplyCount == 0 {
			return qrCode, nil
		}

		// Si llegamos aquí, el QR ya existe, intentar de nuevo
	}

	return "", fmt.Errorf("no se pudo generar un código QR único después de %d intentos", maxAttempts)
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

// stringPtr helper para crear punteros a string
func stringPtr(s string) *string {
	return &s
}
