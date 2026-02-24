package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"meditrack/mailer"
	"meditrack/models"
	"net"
	"os"
	"path/filepath"
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

	return qrInfo, nil
}

// logScanEvent registra automáticamente un evento de escaneo
func (s *QRService) logScanEvent(qrCode string, context *ScanContext, qrInfo *QRInfo, result string, errorMsg string) (*models.QRScanEvent, error) {
	// Verificar si debemos prevenir logging duplicado
	if context != nil && context.ScanPurpose != nil {
		// No registrar escaneos de verificación para transferencias si ya hay uno reciente
		if *context.ScanPurpose == models.ScanPurposeTransferVerification || *context.ScanPurpose == "transfer_check" {
			// Buscar si ya hay un escaneo similar en los últimos 5 minutos
			var recentScan models.QRScanEvent
			fiveMinutesAgo := time.Now().Add(-5 * time.Minute)

			err := s.DB.Where("qr_code = ? AND scan_purpose = ? AND scanned_at > ?",
				qrCode, *context.ScanPurpose, fiveMinutesAgo).
				First(&recentScan).Error

			if err == nil {
				// Ya existe un escaneo reciente, no crear otro
				return &recentScan, nil
			}
		}
	}

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
			u.name as user_name,
			COALESCE(
				sh.location,
				CASE
					WHEN sh.destination_type = 'pavilion' THEN 
						COALESCE('Pabellon: ' || p.name || COALESCE(' (' || mc.name || ')', ''), 'Ubicacion no especificada')
					WHEN sh.destination_type = 'store' THEN 
						COALESCE('Almacen: ' || st.name, 'Ubicacion no especificada')
					ELSE 'Ubicacion no especificada'
				END
			) as location`).
		Joins("LEFT JOIN pavilion p ON sh.destination_type = 'pavilion' AND sh.destination_id = p.id").
		Joins("LEFT JOIN store st ON sh.destination_type = 'store' AND sh.destination_id = st.id").
		Joins("LEFT JOIN medical_center mc ON p.medical_center_id = mc.id").
		Joins("LEFT JOIN \"user\" u ON sh.user_rut = u.rut").
		Where("sh.medical_supply_id = ?", supply.ID).
		Order("sh.date_time DESC").
		Find(&traceability.SupplyHistory).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo historial de movimientos: %v", err)
	}
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
			// No fallar si no hay trazabilidad
		} else {
			result.Traceability = traceability
		}

		// Obtener asignación ACTIVA basada en el carrito activo
		// Primero buscar si existe un item de carrito ACTIVO para este QR
		var cartItem models.SupplyCartItem
		var assignment models.SupplyRequestQRAssignment

		fmt.Printf("\n🔍 ===== DEBUG ESCANEO QR: %s =====\n", qrCode)

		// Debug: Verificar TODAS las asignaciones de este QR
		var allAssignments []models.SupplyRequestQRAssignment
		s.DB.Where("qr_code = ?", qrCode).Order("assigned_date DESC").Find(&allAssignments)
		fmt.Printf("📊 Total de asignaciones para este QR: %d\n", len(allAssignments))
		for i, a := range allAssignments {
			fmt.Printf("  [%d] ID=%d, Status=%s, RequestID=%d, AssignedDate=%s\n",
				i+1, a.ID, a.Status, a.SupplyRequestID, a.AssignedDate.Format("2006-01-02 15:04:05"))
		}

		// Debug: Verificar TODOS los items de carrito para estas asignaciones
		if len(allAssignments) > 0 {
			var assignmentIDs []int
			for _, a := range allAssignments {
				assignmentIDs = append(assignmentIDs, a.ID)
			}
			var allCartItems []models.SupplyCartItem
			s.DB.Where("supply_request_qr_assignment_id IN ?", assignmentIDs).
				Preload("SupplyCart").
				Order("added_at DESC").
				Find(&allCartItems)
			fmt.Printf("🛒 Total de items de carrito para estas asignaciones: %d\n", len(allCartItems))
			for i, ci := range allCartItems {
				fmt.Printf("  [%d] CartItemID=%d, AssignmentID=%d, IsActive=%v, CartNumber=%s, CartStatus=%s, AddedAt=%s\n",
					i+1, ci.ID, ci.SupplyRequestQRAssignmentID, ci.IsActive,
					ci.SupplyCart.CartNumber, ci.SupplyCart.Status, ci.AddedAt.Format("2006-01-02 15:04:05"))
			}
		}

		// Buscar item de carrito activo que contenga una asignación de este QR
		fmt.Printf("\n🔎 Buscando carrito ACTIVO específicamente...\n")
		if err := s.DB.Table("supply_cart_item sci").
			Select("sci.*").
			Joins("INNER JOIN supply_request_qr_assignment srqa ON sci.supply_request_qr_assignment_id = srqa.id").
			Joins("INNER JOIN supply_cart sc ON sci.supply_cart_id = sc.id").
			Where("srqa.qr_code = ? AND sci.is_active = ? AND sc.status = ?",
				qrCode, true, models.CartStatusActive).
			Preload("SupplyCart").
			Preload("SupplyRequestQRAssignment").
			Preload("SupplyRequestQRAssignment.SupplyRequest").
			Preload("SupplyRequestQRAssignment.SupplyRequestItem").
			Order("sci.added_at DESC").
			First(&cartItem).Error; err == nil {

			// Encontramos un carrito activo con este QR
			fmt.Printf("✅ CARRITO ACTIVO ENCONTRADO:\n")
			fmt.Printf("   - CartNumber: %s\n", cartItem.SupplyCart.CartNumber)
			fmt.Printf("   - CartStatus: %s\n", cartItem.SupplyCart.Status)
			fmt.Printf("   - AssignmentID: %d\n", cartItem.SupplyRequestQRAssignmentID)
			fmt.Printf("   - AssignmentStatus: %s\n", cartItem.SupplyRequestQRAssignment.Status)
			fmt.Printf("   - RequestID: %d\n", cartItem.SupplyRequestQRAssignment.SupplyRequestID)
			fmt.Printf("===== FIN DEBUG =====\n\n")

			assignment = cartItem.SupplyRequestQRAssignment
			assignment.Cart = &cartItem.SupplyCart
			result.RequestAssignment = &assignment
			result.SupplyRequest = &assignment.SupplyRequest

		} else {
			fmt.Printf("⚠️ NO se encontró carrito activo. Error: %v\n", err)

			// No hay carrito activo, buscar la última asignación no devuelta
			// IMPORTANTE: NO cargar el carrito aquí, solo la asignación
			if err := s.DB.Where("qr_code = ? AND status != ?", qrCode, models.AssignmentStatusReturned).
				Preload("SupplyRequest").
				Preload("SupplyRequestItem").
				Order("assigned_date DESC").
				First(&assignment).Error; err == nil {

				fmt.Printf("📋 ASIGNACIÓN SIN CARRITO ACTIVO:\n")
				fmt.Printf("   - AssignmentID: %d\n", assignment.ID)
				fmt.Printf("   - Status: %s\n", assignment.Status)
				fmt.Printf("   - RequestID: %d\n", assignment.SupplyRequestID)

				// NO cargar el carrito aquí - dejar Cart como nil
				result.RequestAssignment = &assignment
				result.SupplyRequest = &assignment.SupplyRequest

				// Verificar explícitamente si hay algún carrito (para debug)
				var debugCartItem models.SupplyCartItem
				if err := s.DB.Where("supply_request_qr_assignment_id = ?", assignment.ID).
					Preload("SupplyCart").
					Order("added_at DESC").
					First(&debugCartItem).Error; err == nil {
					fmt.Printf("   ⚠️ ADVERTENCIA: Existe carrito para esta asignación:\n")
					fmt.Printf("      - CartNumber: %s\n", debugCartItem.SupplyCart.CartNumber)
					fmt.Printf("      - CartStatus: %s\n", debugCartItem.SupplyCart.Status)
					fmt.Printf("      - IsActive: %v\n", debugCartItem.IsActive)
					fmt.Printf("      - (NO se cargará porque no está activo)\n")
				} else {
					fmt.Printf("   ℹ️ No hay carrito asociado a esta asignación\n")
				}
			} else {
				fmt.Printf("❌ No se encontró ninguna asignación activa\n")
			}
			fmt.Printf("===== FIN DEBUG =====\n\n")
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
			u.name as user_name,
			COALESCE(
				sh.location,
				CASE
					WHEN sh.destination_type = 'pavilion' THEN 
						COALESCE('Pabellon: ' || p.name || COALESCE(' (' || mc.name || ')', ''), 'Ubicacion no especificada')
					WHEN sh.destination_type = 'store' THEN 
						COALESCE('Almacen: ' || st.name, 'Ubicacion no especificada')
					ELSE 'Ubicacion no especificada'
				END
			) as location`).
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
			"Status":       qrInfo.SupplyInfo.Status, // ✅ CAMPO CRÍTICO AGREGADO
			"status":       qrInfo.SupplyInfo.Status, // ✅ También en minúsculas para compatibilidad
			"InTransit":    qrInfo.SupplyInfo.InTransit,
			"LocationType": qrInfo.SupplyInfo.LocationType,
			"LocationID":   qrInfo.SupplyInfo.LocationID,
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

	return result, nil
}

// TransferSupplyByQR transfiere un insumo individual por su código QR
func (s *QRService) TransferSupplyByQR(qrCode, userRUT, receiverRUT, destinationType string, destinationID int, notes string) (map[string]interface{}, error) {
	var supply models.MedicalSupply
	var result map[string]interface{}

	// Validar formato del RUT
	if len(userRUT) < 9 || !strings.Contains(userRUT, "-") {
		// Si el RUT no tiene el formato correcto, asumir que falta el guión
		if len(userRUT) == 9 {
			// Agregar el guión antes del último dígito
			userRUT = userRUT[:8] + "-" + userRUT[8:]
		} else {
			return nil, fmt.Errorf("formato de RUT inválido: %s", userRUT)
		}
	}

	// DEBUG: Verificar si ya hay una transferencia reciente para este QR (COMENTADO - permitir duplicados)
	/*
		var existingTransfer models.SupplyHistory
		transferInProgress := s.DB.Where(`medical_supply_id = (SELECT id FROM medical_supply WHERE qr_code = ?)
			AND status IN ('en_camino_a_pabellon', 'en_camino_a_bodega')
			AND date_time > ?`, qrCode, time.Now().Add(-5*time.Minute)).First(&existingTransfer).Error == nil

		if transferInProgress {
			fmt.Printf("DEBUG - PREVENCIÓN DUPLICADO: Ya existe transferencia reciente para QR=%s, TransferID=%d, Status=%s, Time=%s\n",
				qrCode, existingTransfer.ID, existingTransfer.Status, existingTransfer.DateTime.Format("2006-01-02 15:04:05"))
			return nil, fmt.Errorf("ya existe una transferencia en progreso para este insumo - ID: %d", existingTransfer.ID)
		}
	*/

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Buscar el insumo por QR
		if err := tx.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
			return fmt.Errorf("insumo no encontrado con QR %s: %v", qrCode, err)
		}

		// Verificar que el insumo no haya sido consumido previamente
		if supply.Status == models.StatusConsumed {
			return fmt.Errorf("el insumo con QR %s ya ha sido consumido y no puede ser transferido", qrCode)
		}

		// Verificar que el insumo esté disponible para transferencia (permitir desde estados válidos)
		validStates := []string{models.StatusAvailable, models.StatusReceived}
		isValidState := false
		for _, state := range validStates {
			if supply.Status == state {
				isValidState = true
				break
			}
		}

		if !isValidState {
			return fmt.Errorf("el insumo tiene estado '%s' y no está disponible para transferencia", supply.Status)
		}

		// Para transferencias a pabellón, el insumo queda en "pendiente_retiro"
		// hasta que se escanee para registrar el retiro físico
		var newStatus string
		if destinationType == "pavilion" {
			newStatus = models.StatusPendingPickup // Pendiente de retiro físico
		} else {
			// Para transferencias a bodega, mantener el flujo anterior
			newStatus = models.StatusEnRouteToStore
		}

		// Actualizar el estado del insumo
		if err := tx.Model(&supply).Update("status", newStatus).Error; err != nil {
			return fmt.Errorf("error actualizando estado del insumo: %v", err)
		}

		// Crear entrada en el historial de suministros para la transferencia
		historyStatus := newStatus
		historyNotes := fmt.Sprintf("Transferencia a %s - %s", destinationType, notes)
		if destinationType == "pavilion" {
			historyNotes = fmt.Sprintf("Preparado para transferencia a pabellón. Debe ser escaneado al retirar de bodega. - %s", notes)
		}

		history := models.SupplyHistory{
			DateTime:        time.Now(),
			Status:          historyStatus,
			DestinationType: destinationType,
			DestinationID:   destinationID,
			MedicalSupplyID: supply.ID,
			UserRUT:         userRUT,
			Notes:           historyNotes,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("error creando historial de transferencia: %v", err)
		}

		// Crear registro en supply_transfer para tracking en la vista de transferencias
		transferCode := fmt.Sprintf("TRF-%d-%s", supply.ID, time.Now().Format("20060102150405"))

		// Determinar origen basado en la ubicación actual del insumo
		originType := supply.LocationType
		originID := supply.LocationID

		// Si no tiene ubicación válida, asumir bodega principal
		if originID == 0 {
			originType = models.TransferLocationStore
			originID = 1
		}

		// Obtener información del usuario
		var userName string
		var user models.User
		if err := tx.Where("rut = ?", userRUT).First(&user).Error; err == nil {
			userName = user.Name
		} else {
			userName = "Usuario Desconocido"
		}

		// Determinar el estado de la transferencia
		transferStatus := models.TransferStatusInTransit
		if destinationType == "pavilion" {
			transferStatus = models.TransferStatusPending // Pendiente hasta que se retire físicamente
		}

		supplyTransfer := models.SupplyTransfer{
			TransferCode:    transferCode,
			QRCode:          qrCode,
			MedicalSupplyID: supply.ID,
			OriginType:      originType,
			OriginID:        originID,
			DestinationType: destinationType,
			DestinationID:   destinationID,
			SentBy:          userRUT,
			SentByName:      userName,
			Status:          transferStatus,
			TransferReason:  "Transferencia realizada desde escáner QR",
			SendDate:        time.Now(),
			Notes:           notes,
		}

		if err := tx.Create(&supplyTransfer).Error; err != nil {
			return fmt.Errorf("error creando registro de transferencia: %v", err)
		}

		// Actualizar contadores de inventario si se transfiere desde bodega
		if originType == models.TransferLocationStore || supply.LocationType == models.SupplyLocationStore {
			var batch models.Batch
			if err := tx.First(&batch, supply.BatchID).Error; err != nil {
				return fmt.Errorf("error obteniendo batch: %v", err)
			}

			var storeSummary models.StoreInventorySummary
			storeID := originID
			if storeID == 0 {
				storeID = supply.LocationID
			}
			if err := tx.Where("store_id = ? AND batch_id = ?", storeID, supply.BatchID).
				First(&storeSummary).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Crear resumen si no existe
					storeSummary = models.StoreInventorySummary{
						StoreID:             storeID,
						BatchID:             supply.BatchID,
						SupplyCode:          supply.Code,
						SurgeryID:           batch.SurgeryID,
						OriginalAmount:      batch.Amount,
						CurrentInStore:      batch.Amount - 1,
						TotalTransferredOut: 1,
					}
					now := time.Now()
					storeSummary.LastTransferOutDate = &now
					if err := tx.Create(&storeSummary).Error; err != nil {
						return fmt.Errorf("error creando resumen de bodega: %v", err)
					}
				} else {
					return fmt.Errorf("error obteniendo resumen de bodega: %v", err)
				}
			} else {
				// Actualizar resumen existente
				storeSummary.CurrentInStore--
				storeSummary.TotalTransferredOut++
				now := time.Now()
				storeSummary.LastTransferOutDate = &now
				if err := tx.Save(&storeSummary).Error; err != nil {
					return fmt.Errorf("error actualizando resumen de bodega: %v", err)
				}
			}
		}

		// NO crear QRScanEvent aquí - SupplyHistory es suficiente para las transferencias
		// Los QRScanEvent son solo para escaneos, no para cambios de estado

		// Crear resultado de la transferencia
		result = map[string]interface{}{
			"success":            true,
			"supply_id":          supply.ID,
			"qr_code":            qrCode,
			"old_status":         models.StatusAvailable,
			"new_status":         newStatus,
			"destination_type":   destinationType,
			"destination_id":     destinationID,
			"user_rut":           userRUT,
			"receiver_rut":       receiverRUT,
			"transfer_timestamp": time.Now(),
			"status_change": map[string]string{
				"from": models.StatusAvailable,
				"to":   newStatus,
			},
			"requires_pickup": destinationType == "pavilion", // Indica si requiere escaneo de retiro
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// PickupSupplyFromStore registra el retiro físico de un insumo de bodega
// Paso 1: Cuando alguien viene a retirar físicamente el insumo
func (s *QRService) PickupSupplyFromStore(qrCode, userRUT string, notes string) (map[string]interface{}, error) {
	var supply models.MedicalSupply
	var transfer models.SupplyTransfer
	var userName string

	// Obtener información del usuario
	var user models.User
	if err := s.DB.Where("rut = ?", userRUT).First(&user).Error; err == nil {
		userName = user.Name
	} else {
		userName = "Usuario Desconocido"
	}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Buscar el insumo por QR
		if err := tx.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
			return fmt.Errorf("insumo no encontrado con QR %s: %v", qrCode, err)
		}

		// 2. Verificar que el insumo esté en estado "pendiente_retiro"
		if supply.Status != models.StatusPendingPickup {
			return fmt.Errorf("el insumo debe estar en estado 'pendiente_retiro' para ser retirado, estado actual: %s", supply.Status)
		}

		// 3. Buscar la transferencia pendiente asociada a este QR
		if err := tx.Where("qr_code = ? AND status = ?", qrCode, models.TransferStatusPending).
			First(&transfer).Error; err != nil {
			return fmt.Errorf("transferencia pendiente no encontrada para QR %s: %v", qrCode, err)
		}

		// 3.5. Validar quién puede retirar el insumo según la configuración de la solicitud
		var assignment models.SupplyRequestQRAssignment
		if err := tx.Where("qr_code = ?", qrCode).
			Preload("SupplyRequest").
			Order("assigned_date DESC").
			First(&assignment).Error; err == nil {
			// Si hay una asignación, verificar la configuración de retiro
			request := assignment.SupplyRequest
			if !request.AllowAnyoneToPickup {
				// Solo una persona específica puede retirar
				if request.AuthorizedPickupRUT == nil || *request.AuthorizedPickupRUT == "" {
					return fmt.Errorf("la solicitud está configurada para retiro restringido, pero no se ha especificado una persona autorizada")
				}
				if *request.AuthorizedPickupRUT != userRUT {
					authorizedName := "Usuario Desconocido"
					if request.AuthorizedPickupName != nil {
						authorizedName = *request.AuthorizedPickupName
					}
					return fmt.Errorf("solo %s (RUT: %s) está autorizado para retirar este insumo. Su RUT es %s", authorizedName, *request.AuthorizedPickupRUT, userRUT)
				}
			}
			// Si AllowAnyoneToPickup es true, cualquiera puede retirar (no se valida)
		}
		// Si no hay asignación a solicitud, se permite el retiro (comportamiento anterior)

		// 4. Registrar quién retiró físicamente
		now := time.Now()
		transfer.PickedUpBy = &userRUT
		transfer.PickedUpByName = &userName
		transfer.PickedUpDate = &now
		transfer.Status = models.TransferStatusInTransit
		if notes != "" {
			if transfer.Notes != "" {
				transfer.Notes = transfer.Notes + "\nRetiro: " + notes
			} else {
				transfer.Notes = "Retiro: " + notes
			}
		}

		if err := tx.Save(&transfer).Error; err != nil {
			return fmt.Errorf("error al actualizar transferencia: %v", err)
		}

		// 5. Actualizar el estado del insumo a "en_camino_a_pabellon" y cambiar ubicación
		// Ahora SÍ cambia la ubicación porque físicamente está siendo retirado
		supply.Status = models.StatusEnRouteToPavilion
		supply.LocationType = models.SupplyLocationPavilion
		supply.LocationID = transfer.DestinationID
		supply.InTransit = true
		supply.UpdatedAt = now

		if err := tx.Save(&supply).Error; err != nil {
			return fmt.Errorf("error al actualizar estado del insumo: %v", err)
		}

		// 6. Registrar en historial
		// IMPORTANTE: Usar OriginType y OriginID de la transferencia para registrar correctamente
		// la bodega de origen (puede ser bodega secundaria, no necesariamente la principal)
		originType := transfer.OriginType
		originID := transfer.OriginID
		history := models.SupplyHistory{
			DateTime:        now,
			Status:          models.StatusEnRouteToPavilion,
			DestinationType: transfer.DestinationType,
			DestinationID:   transfer.DestinationID,
			MedicalSupplyID: supply.ID,
			UserRUT:         userRUT,
			Notes:           fmt.Sprintf("Retirado físicamente de bodega por %s. En camino a pabellón.", userName),
			OriginType:      &originType,
			OriginID:        &originID,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("error al crear historial de retiro: %v", err)
		}

		// 7. Retirar automáticamente todos los insumos del mismo carrito que estén pendientes de retiro
		// Buscar el carrito asociado a este insumo
		var cartAssignment models.SupplyRequestQRAssignment
		if err := tx.Where("medical_supply_id = ?", supply.ID).First(&cartAssignment).Error; err == nil {
			// Buscar el item del carrito asociado a esta asignación
			var cartItem models.SupplyCartItem
			if err := tx.Where("supply_request_qr_assignment_id = ? AND is_active = ?", cartAssignment.ID, true).
				First(&cartItem).Error; err == nil {
				// Obtener todos los items activos del carrito
				var cartItems []models.SupplyCartItem
				if err := tx.Where("supply_cart_id = ? AND is_active = ?", cartItem.SupplyCartID, true).
					Preload("SupplyRequestQRAssignment").
					Preload("SupplyRequestQRAssignment.MedicalSupply").
					Find(&cartItems).Error; err == nil {

					// Retirar todos los insumos del carrito que estén pendientes de retiro
					retiradosCount := 0
					for _, item := range cartItems {
						// Saltar el insumo que ya fue retirado
						if item.SupplyRequestQRAssignment.MedicalSupplyID == supply.ID {
							continue
						}

						otherSupply := item.SupplyRequestQRAssignment.MedicalSupply

						// Solo retirar insumos que estén pendientes de retiro y pertenezcan al mismo pabellón
						if otherSupply.Status == models.StatusPendingPickup &&
							!otherSupply.InTransit &&
							otherSupply.LocationType == models.SupplyLocationStore {

							// Buscar la transferencia pendiente asociada
							var otherTransfer models.SupplyTransfer
							if err := tx.Where("qr_code = ? AND status = ?", otherSupply.QRCode, models.TransferStatusPending).
								First(&otherTransfer).Error; err == nil {

								// Verificar que el destino sea el mismo pabellón
								if otherTransfer.DestinationType == transfer.DestinationType &&
									otherTransfer.DestinationID == transfer.DestinationID {

									// Actualizar la transferencia con información de retiro
									otherTransfer.Status = models.TransferStatusInTransit
									otherTransfer.PickedUpBy = &userRUT
									otherTransfer.PickedUpByName = &userName
									otherTransfer.PickedUpDate = &now
									if notes != "" {
										if otherTransfer.Notes != "" {
											otherTransfer.Notes = otherTransfer.Notes + "\n" + notes + " (retirado automáticamente con el carrito)"
										} else {
											otherTransfer.Notes = notes + " (retirado automáticamente con el carrito)"
										}
									} else {
										otherTransfer.Notes = "Retirado automáticamente con otros insumos del carrito"
									}

									if err := tx.Save(&otherTransfer).Error; err != nil {
										continue
									}

									// Actualizar el estado del insumo
									otherSupply.Status = models.StatusEnRouteToPavilion
									otherSupply.LocationType = models.SupplyLocationPavilion
									otherSupply.LocationID = transfer.DestinationID
									otherSupply.InTransit = true
									otherSupply.UpdatedAt = now

									if err := tx.Save(&otherSupply).Error; err != nil {
										continue
									}

									// Crear historial para el retiro automático
									// IMPORTANTE: Usar OriginType y OriginID de la transferencia para registrar correctamente
									// la bodega de origen (puede ser bodega secundaria, no necesariamente la principal)
									otherOriginType := otherTransfer.OriginType
									otherOriginID := otherTransfer.OriginID
									otherHistory := models.SupplyHistory{
										DateTime:        now,
										Status:          models.StatusEnRouteToPavilion,
										DestinationType: transfer.DestinationType,
										DestinationID:   transfer.DestinationID,
										MedicalSupplyID: otherSupply.ID,
										UserRUT:         userRUT,
										Notes:           fmt.Sprintf("Retirado automáticamente con el carrito. En camino a pabellón."),
										OriginType:      &otherOriginType,
										OriginID:        &otherOriginID,
									}

									if err := tx.Create(&otherHistory).Error; err != nil {
										continue
									}

									retiradosCount++
								}
							}
						}
					}

					// Log de cuántos insumos adicionales fueron retirados
					if retiradosCount > 0 {
						fmt.Printf("✅ Retirados automáticamente %d insumos adicionales del carrito\n", retiradosCount)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":           true,
		"qr_code":           qrCode,
		"picked_up_by":      userRUT,
		"picked_up_by_name": userName,
		"new_status":        models.StatusEnRouteToPavilion,
		"pickup_timestamp":  time.Now(),
		"message":           "Retiro registrado exitosamente. El insumo está en camino al pabellón.",
	}, nil
}

// ReceiveSupplyByQR recepciona un insumo que está en estado "en_camino_a_pabellon"
// Paso 2: Cuando llega al pabellón y se confirma la recepción
// El parámetro willBeConsumed indica si el insumo será usado (true) o devuelto inmediatamente (false)
func (s *QRService) ReceiveSupplyByQR(qrCode, userRUT, destinationType string, destinationID int, notes string, willBeConsumed bool) (map[string]interface{}, error) {
	var supply models.MedicalSupply
	var transfer models.SupplyTransfer
	var result map[string]interface{}
	var userName string

	// Obtener información del usuario
	var user models.User
	if err := s.DB.Where("rut = ?", userRUT).First(&user).Error; err == nil {
		userName = user.Name
	} else {
		userName = "Usuario Desconocido"
	}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Buscar el insumo por QR
		if err := tx.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
			return fmt.Errorf("insumo no encontrado con QR %s: %v", qrCode, err)
		}

		// 2. Verificar que el insumo esté en estado "en_camino_a_pabellon"
		if supply.Status != models.StatusEnRouteToPavilion {
			return fmt.Errorf("el insumo debe estar en estado 'en_camino_a_pabellon' para ser recepcionado, estado actual: %s", supply.Status)
		}

		// 3. Buscar la transferencia activa asociada a este QR
		if err := tx.Where("qr_code = ? AND status = ?", qrCode, models.TransferStatusInTransit).
			First(&transfer).Error; err != nil {
			return fmt.Errorf("transferencia en tránsito no encontrada para QR %s: %v", qrCode, err)
		}

		// 3.5. Validar que el insumo haya sido retirado físicamente
		if transfer.PickedUpBy == nil {
			return fmt.Errorf("el insumo no ha sido retirado físicamente de bodega aún. Debe escanearlo primero para registrar el retiro.")
		}

		// 3.6. Validar que solo usuarios con rol "pabellón" pueden recepcionar
		if transfer.DestinationType == models.TransferLocationPavilion {
			// Obtener el usuario que está recepcionando
			var receivingUser models.User
			if err := tx.Where("rut = ?", userRUT).First(&receivingUser).Error; err != nil {
				return fmt.Errorf("usuario no encontrado: %v", err)
			}

			// Solo usuarios con rol "pabellón" pueden recepcionar insumos
			if receivingUser.Role != "pabellón" {
				return fmt.Errorf("solo usuarios con rol 'pabellón' pueden recepcionar insumos")
			}

			// Verificar que el usuario tenga un pabellón asignado
			if receivingUser.PavilionID == nil {
				return fmt.Errorf("el usuario de pabellón no tiene un pabellón asignado")
			}

			// Verificar que el PavilionID del usuario coincida con el destino de la transferencia
			if *receivingUser.PavilionID != transfer.DestinationID {
				// Obtener nombre del pabellón esperado
				var expectedPavilion models.Pavilion
				pavilionName := fmt.Sprintf("Pabellón %d", transfer.DestinationID)
				if err := tx.First(&expectedPavilion, transfer.DestinationID).Error; err == nil {
					pavilionName = expectedPavilion.Name
				}
				return fmt.Errorf("este insumo debe ser recepcionado por el usuario del %s. Su usuario pertenece al pabellón %d", pavilionName, *receivingUser.PavilionID)
			}
		}

		// 4. Actualizar el estado de la transferencia
		now := time.Now()
		transfer.Status = models.TransferStatusReceived
		transfer.ReceiveDate = &now
		transfer.ReceivedBy = &userRUT
		transfer.ReceivedByName = &userName
		if notes != "" {
			if transfer.Notes != "" {
				transfer.Notes = transfer.Notes + "\n" + notes
			} else {
				transfer.Notes = notes
			}
		}

		if err := tx.Save(&transfer).Error; err != nil {
			return fmt.Errorf("error al actualizar transferencia: %v", err)
		}

		// 5. Actualizar el estado del insumo según si será consumido o no
		if willBeConsumed {
			// Si será consumido, queda en estado "recepcionado" para uso posterior
			supply.Status = models.StatusReceived
		} else {
			// Si no será consumido, se marca como disponible para devolución
			supply.Status = models.StatusAvailable
		}
		supply.InTransit = false
		supply.LocationType = transfer.DestinationType
		supply.LocationID = transfer.DestinationID
		supply.UpdatedAt = now

		if err := tx.Save(&supply).Error; err != nil {
			return fmt.Errorf("error al actualizar estado del insumo: %v", err)
		}

		// 5.5. Actualizar la asignación QR según la acción
		if transfer.DestinationType == models.TransferLocationPavilion {
			var assignment models.SupplyRequestQRAssignment
			if err := tx.Where("qr_code = ? AND status = ?", qrCode, models.AssignmentStatusAssigned).
				First(&assignment).Error; err == nil {

				if willBeConsumed {
					// Si será consumido, marcar como entregado
					assignment.Status = models.AssignmentStatusDelivered
				} else {
					// Si no será consumido, marcar como devuelto inmediatamente
					assignment.Status = models.AssignmentStatusReturned
				}

				assignment.DeliveredDate = &now
				assignment.DeliveredBy = &userRUT
				assignment.DeliveredByName = &userName
				assignment.UpdatedAt = now

				if err := tx.Save(&assignment).Error; err != nil {
					// No fallar si no se puede actualizar la asignación, solo loguear
					fmt.Printf("⚠️  Advertencia: No se pudo actualizar asignación QR: %v\n", err)
				}
			}
		}

		// 6. Actualizar inventario según el tipo de destino
		if transfer.DestinationType == models.TransferLocationPavilion {
			// Incrementar stock en pabellón
			var pavilionSummary models.PavilionInventorySummary
			if err := tx.Where("pavilion_id = ? AND batch_id = ?", transfer.DestinationID, supply.BatchID).
				First(&pavilionSummary).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Crear resumen si no existe
					pavilionSummary = models.PavilionInventorySummary{
						PavilionID:       transfer.DestinationID,
						BatchID:          supply.BatchID,
						SupplyCode:       supply.Code,
						TotalReceived:    1,
						CurrentAvailable: 1,
						LastReceivedDate: &now,
					}
					if err := tx.Create(&pavilionSummary).Error; err != nil {
						return fmt.Errorf("error al crear resumen de pabellón: %v", err)
					}
				} else {
					return fmt.Errorf("error al obtener resumen de pabellón: %v", err)
				}
			} else {
				// Actualizar resumen existente
				pavilionSummary.TotalReceived++
				pavilionSummary.CurrentAvailable++
				pavilionSummary.LastReceivedDate = &now

				if err := tx.Save(&pavilionSummary).Error; err != nil {
					return fmt.Errorf("error al actualizar resumen de pabellón: %v", err)
				}
			}
		} else if transfer.DestinationType == models.TransferLocationStore {
			// Incrementar stock en bodega
			var storeSummary models.StoreInventorySummary
			if err := tx.Where("store_id = ? AND batch_id = ?", transfer.DestinationID, supply.BatchID).
				First(&storeSummary).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Obtener información del batch para el original amount
					var batch models.Batch
					if err := tx.First(&batch, supply.BatchID).Error; err != nil {
						return fmt.Errorf("error al obtener información del batch: %v", err)
					}

					// Crear resumen si no existe
					storeSummary = models.StoreInventorySummary{
						StoreID:          transfer.DestinationID,
						BatchID:          supply.BatchID,
						SupplyCode:       supply.Code,
						OriginalAmount:   batch.Amount,
						CurrentInStore:   1,
						TotalReturnedIn:  1,
						LastReturnInDate: &now,
					}
					if err := tx.Create(&storeSummary).Error; err != nil {
						return fmt.Errorf("error al crear resumen de bodega: %v", err)
					}
				} else {
					return fmt.Errorf("error al obtener resumen de bodega: %v", err)
				}
			} else {
				// Actualizar resumen existente - incrementar devoluciones
				storeSummary.CurrentInStore++
				storeSummary.TotalReturnedIn++
				storeSummary.LastReturnInDate = &now

				if err := tx.Save(&storeSummary).Error; err != nil {
					return fmt.Errorf("error al actualizar resumen de bodega: %v", err)
				}
			}
		}

		// 7. Crear entrada en el historial de suministros
		originType := transfer.OriginType
		originID := transfer.OriginID
		history := models.SupplyHistory{
			DateTime:        now,
			Status:          models.StatusReceived,
			DestinationType: transfer.DestinationType,
			DestinationID:   transfer.DestinationID,
			MedicalSupplyID: supply.ID,
			UserRUT:         userRUT,
			Notes:           fmt.Sprintf("Recepción confirmada desde escáner QR - %s", notes),
			OriginType:      &originType,
			OriginID:        &originID,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("error al crear historial de recepción: %v", err)
		}

		// 8. Recepcionar automáticamente todos los insumos del mismo carrito que estén en tránsito
		// Buscar el carrito asociado a este insumo
		var assignment models.SupplyRequestQRAssignment
		if err := tx.Where("medical_supply_id = ?", supply.ID).First(&assignment).Error; err == nil {
			// Buscar el item del carrito asociado a esta asignación
			var cartItem models.SupplyCartItem
			if err := tx.Where("supply_request_qr_assignment_id = ? AND is_active = ?", assignment.ID, true).
				First(&cartItem).Error; err == nil {
				// Obtener todos los items activos del carrito
				var cartItems []models.SupplyCartItem
				if err := tx.Where("supply_cart_id = ? AND is_active = ?", cartItem.SupplyCartID, true).
					Preload("SupplyRequestQRAssignment").
					Preload("SupplyRequestQRAssignment.MedicalSupply").
					Find(&cartItems).Error; err == nil {

					// Recepcionar todos los insumos del carrito que estén en tránsito al pabellón
					recepcionadosCount := 0
					for _, item := range cartItems {
						// Saltar el insumo que ya fue recepcionado
						if item.SupplyRequestQRAssignment.MedicalSupplyID == supply.ID {
							continue
						}

						otherSupply := item.SupplyRequestQRAssignment.MedicalSupply

						// Solo recepcionar insumos que estén en tránsito al pabellón
						if otherSupply.Status == models.StatusEnRouteToPavilion &&
							otherSupply.InTransit &&
							otherSupply.LocationType == models.SupplyLocationPavilion &&
							otherSupply.LocationID == transfer.DestinationID {

							// Buscar la transferencia asociada
							var otherTransfer models.SupplyTransfer
							if err := tx.Where("qr_code = ? AND status = ?", otherSupply.QRCode, models.TransferStatusInTransit).
								First(&otherTransfer).Error; err == nil {

								// Actualizar la transferencia
								otherTransfer.Status = models.TransferStatusReceived
								otherTransfer.ReceiveDate = &now
								otherTransfer.ReceivedBy = &userRUT
								otherTransfer.ReceivedByName = &userName
								if notes != "" {
									if otherTransfer.Notes != "" {
										otherTransfer.Notes = otherTransfer.Notes + "\n" + notes + " (recepcionado automáticamente con el carrito)"
									} else {
										otherTransfer.Notes = notes + " (recepcionado automáticamente con el carrito)"
									}
								} else {
									otherTransfer.Notes = "Recepcionado automáticamente con otros insumos del carrito"
								}

								if err := tx.Save(&otherTransfer).Error; err != nil {
									continue
								}

								// Actualizar el estado del insumo
								otherSupply.Status = models.StatusReceived
								otherSupply.InTransit = false
								otherSupply.UpdatedAt = now

								if err := tx.Save(&otherSupply).Error; err != nil {
									continue
								}

								// Actualizar inventario del pabellón
								var otherPavilionSummary models.PavilionInventorySummary
								if err := tx.Where("pavilion_id = ? AND batch_id = ?", transfer.DestinationID, otherSupply.BatchID).
									First(&otherPavilionSummary).Error; err != nil {
									if errors.Is(err, gorm.ErrRecordNotFound) {
										// Crear resumen si no existe
										otherPavilionSummary = models.PavilionInventorySummary{
											PavilionID:       transfer.DestinationID,
											BatchID:          otherSupply.BatchID,
											SupplyCode:       otherSupply.Code,
											TotalReceived:    1,
											CurrentAvailable: 1,
											LastReceivedDate: &now,
										}
										if err := tx.Create(&otherPavilionSummary).Error; err != nil {
											continue
										}
									} else {
										continue
									}
								} else {
									// Actualizar resumen existente
									otherPavilionSummary.TotalReceived++
									otherPavilionSummary.CurrentAvailable++
									otherPavilionSummary.LastReceivedDate = &now

									if err := tx.Save(&otherPavilionSummary).Error; err != nil {
										continue
									}
								}

								// Registrar en historial
								otherHistory := models.SupplyHistory{
									MedicalSupplyID: otherSupply.ID,
									DateTime:        now,
									Status:          models.StatusReceived,
									DestinationType: models.DestinationTypePavilion,
									DestinationID:   transfer.DestinationID,
									UserRUT:         userRUT,
									Notes:           fmt.Sprintf("Recepcionado automáticamente junto con el insumo %s del mismo carrito", qrCode),
									OriginType:      &otherTransfer.OriginType,
									OriginID:        &otherTransfer.OriginID,
								}

								if err := tx.Create(&otherHistory).Error; err != nil {
									continue
								}

								recepcionadosCount++
							}
						}
					}

				}
			}
		}

		// 9. Programar alerta para insumo recepcionado (después de la transacción)
		go func() {
			// Esperar 1 minuto para pruebas (cambiar a 12 * time.Hour en producción)
			time.Sleep(1 * time.Minute)

			// Verificar si el insumo sigue en estado "Recepcionado"
			var updatedSupply models.MedicalSupply
			if err := s.DB.First(&updatedSupply, supply.ID).Error; err != nil {
				return
			}

			// Si sigue en estado "Recepcionado", enviar alerta
			if updatedSupply.Status == models.StatusReceived {
				s.sendUnconsumedSupplyAlert(&updatedSupply)
			}
		}()

		// Preparar respuesta
		result = map[string]interface{}{
			"success":            true,
			"transfer_code":      transfer.TransferCode,
			"supply_id":          supply.ID,
			"qr_code":            supply.QRCode,
			"previous_status":    models.StatusEnRouteToPavilion,
			"new_status":         models.StatusReceived,
			"received_at":        now,
			"received_by":        userRUT,
			"received_by_name":   userName,
			"destination_type":   transfer.DestinationType,
			"destination_id":     transfer.DestinationID,
			"notes":              notes,
			"inventory_updated":  true,
			"transfer_completed": true,
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

	// Variables para guardar información del batch fuera de la transacción
	var batchID int
	var previousAmount int

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
		if err := tx.Preload("SupplierConfig").Where("id = ?", supply.BatchID).First(&batch).Error; err != nil {
			return fmt.Errorf("lote no encontrado: %v", err)
		}

		// Guardar información del batch para uso fuera de la transacción
		batchID = batch.ID
		previousAmount = batch.Amount

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

		// Actualizar contadores de inventario según la ubicación del insumo
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
						SurgeryID:            batch.SurgeryID,    // Obtener del batch
						OriginalAmount:       int(realCount) + 1, // Cantidad antes del consumo (real + 1 consumido)
						CurrentInStore:       int(realCount),     // Stock actual en bodega (sin el que se consumió)
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

	// Verificar si se debe enviar alerta de stock bajo después de consumir el insumo
	// Obtener el batch actualizado usando el BatchID guardado
	var updatedBatch models.Batch
	if err := s.DB.First(&updatedBatch, batchID).Error; err == nil {
		// Obtener el stock crítico del tipo de insumo
		var supplyCode models.SupplyCode
		if err := s.DB.Joins("JOIN medical_supply ON medical_supply.code = supply_code.code").
			Where("medical_supply.batch_id = ?", batchID).
			First(&supplyCode).Error; err == nil {

			// Verificar si el stock está en nivel crítico
			// Enviar alerta cuando el stock llega al nivel crítico o por debajo
			newAmount := updatedBatch.Amount

			// Solo enviar alerta si:
			// 1. El nuevo stock está en o por debajo del crítico
			// 2. El stock anterior era mayor al crítico (para evitar alertas repetidas cuando ya está en crítico)
			if newAmount > 0 && newAmount <= supplyCode.CriticalStock {
				// Si el stock anterior era mayor al crítico, es la primera vez que llega al crítico
				if previousAmount > supplyCode.CriticalStock {
					// Enviar alerta en una goroutine para no bloquear la respuesta
					go func() {
						_ = s.sendLowStockAlertForBatch(updatedBatch, supplyCode)
					}()
				}
			}
		}
	}

	return &response, nil
}

// sendLowStockAlertForBatch envía correo de alerta de stock bajo para un lote
func (s *QRService) sendLowStockAlertForBatch(batch models.Batch, supplyCode models.SupplyCode) error {
	// Verificar que las variables de entorno del mailer estén configuradas
	emailDir := os.Getenv("EMAIL_DIR")
	emailPass := os.Getenv("EMAIL_PASS")
	emailServer := os.Getenv("EMAIL_SERVER")
	emailPort := os.Getenv("EMAIL_PORT")

	if emailDir == "" || emailPass == "" || emailServer == "" || emailPort == "" {
		return fmt.Errorf("variables de entorno del mailer no configuradas")
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

// getSupplyWithBatchInfo obtiene información completa de un insumo con datos del lote
func (s *QRService) getSupplyWithBatchInfo(qrCode string) (map[string]interface{}, error) {
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

// sendUnconsumedSupplyAlert envía una alerta por correo para un insumo no consumido
func (s *QRService) sendUnconsumedSupplyAlert(supply *models.MedicalSupply) {
	// Obtener información adicional del insumo
	supplyInfo, err := s.getSupplyWithBatchInfo(supply.QRCode)
	if err != nil {
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
		return
	}

	templatePath := filepath.Join(wd, "mailer", "templates", "unconsumed_supply_alert.html")
	_ = request.SendMailSkipTLS(templatePath, emailData)
}
