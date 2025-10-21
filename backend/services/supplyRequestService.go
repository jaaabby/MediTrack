package services

import (
	"fmt"
	"meditrack/mailer"
	"meditrack/models"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SupplyRequestService struct {
	DB *gorm.DB
}

func NewSupplyRequestService(db *gorm.DB) *SupplyRequestService {
	return &SupplyRequestService{DB: db}
}

// FlexibleTime es un tipo personalizado que acepta múltiples formatos de fecha
type FlexibleTime struct {
	time.Time
}

// UnmarshalJSON implementa el unmarshaler personalizado para aceptar múltiples formatos
func (ft *FlexibleTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		ft.Time = time.Time{}
		return nil
	}

	// Intentar diferentes formatos
	formats := []string{
		time.RFC3339,          // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05", // 2006-01-02T15:04:05
		"2006-01-02T15:04",    // 2006-01-02T15:04 (datetime-local)
		"2006-01-02 15:04:05", // 2006-01-02 15:04:05
		"2006-01-02 15:04",    // 2006-01-02 15:04
	}

	var err error
	for _, format := range formats {
		ft.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("no se pudo parsear la fecha: %s", s)
}

// CreateSupplyRequestRequest representa la estructura para crear una solicitud
type CreateSupplyRequestRequest struct {
	PavilionID      int                              `json:"pavilion_id" binding:"required"`
	RequestedBy     string                           `json:"requested_by" binding:"required"`
	RequestedByName string                           `json:"requested_by_name" binding:"required"`
	SurgeryDatetime FlexibleTime                     `json:"surgery_datetime" binding:"required"`
	Notes           string                           `json:"notes"`
	Items           []CreateSupplyRequestItemRequest `json:"items" binding:"required,dive"`
}

type CreateSupplyRequestItemRequest struct {
	SupplyCode        int    `json:"supply_code" binding:"required"`
	SupplyName        string `json:"supply_name" binding:"required"`
	QuantityRequested int    `json:"quantity_requested" binding:"required,min=1"`
	Specifications    string `json:"specifications"`
	IsPediatric       bool   `json:"is_pediatric"`
	Size              string `json:"size"`
	Brand             string `json:"brand"`
	SpecialRequests   string `json:"special_requests"`
	UrgencyLevel      string `json:"urgency_level"`
}

// ApproveSupplyRequestRequest representa la estructura para aprobar una solicitud
type ApproveSupplyRequestRequest struct {
	ApprovedBy     string                            `json:"approved_by" binding:"required"`
	ApprovedByName string                            `json:"approved_by_name" binding:"required"`
	Notes          string                            `json:"notes"`
	ItemApprovals  []ApproveSupplyRequestItemRequest `json:"item_approvals" binding:"required,dive"`
}

type ApproveSupplyRequestItemRequest struct {
	ItemID           int `json:"item_id" binding:"required"`
	QuantityApproved int `json:"quantity_approved" binding:"min=0"`
}

// AssignQRToRequestRequest representa la estructura para asignar QRs
type AssignQRToRequestRequest struct {
	SupplyRequestID     int    `json:"supply_request_id" binding:"required"`
	SupplyRequestItemID int    `json:"supply_request_item_id" binding:"required"`
	QRCode              string `json:"qr_code" binding:"required"`
	AssignedBy          string `json:"assigned_by" binding:"required"`
	AssignedByName      string `json:"assigned_by_name" binding:"required"`
	Notes               string `json:"notes"`
}

// CreateSupplyRequest crea una nueva solicitud de insumo con sus items
func (s *SupplyRequestService) CreateSupplyRequest(request CreateSupplyRequestRequest) (*models.SupplyRequest, error) {
	var supplyRequest models.SupplyRequest
	var medicalCenterID int

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Obtener el medical_center_id desde el pavilion
		var pavilion models.Pavilion
		if err := tx.First(&pavilion, request.PavilionID).Error; err != nil {
			return fmt.Errorf("pabellón no encontrado: %v", err)
		}
		medicalCenterID = pavilion.MedicalCenterID

		// Crear la solicitud principal
		supplyRequest = models.SupplyRequest{
			RequestNumber:   models.GenerateRequestNumber(),
			PavilionID:      request.PavilionID,
			RequestedBy:     request.RequestedBy,
			RequestedByName: request.RequestedByName,
			RequestDate:     time.Now(),
			SurgeryDatetime: request.SurgeryDatetime.Time,
			Status:          models.RequestStatusPendingPavedad, // Estado inicial: pendiente de asignación por Pavedad
			Notes:           request.Notes,
			MedicalCenterID: medicalCenterID,
		}

		if err := tx.Create(&supplyRequest).Error; err != nil {
			return fmt.Errorf("error creando solicitud: %v", err)
		}

		// Crear los items de la solicitud
		for _, itemReq := range request.Items {
			// Validar que el código de insumo existe
			var supplyCode models.SupplyCode
			if err := tx.First(&supplyCode, itemReq.SupplyCode).Error; err != nil {
				return fmt.Errorf("código de insumo %d no encontrado: %v", itemReq.SupplyCode, err)
			}

			item := models.SupplyRequestItem{
				SupplyRequestID:   supplyRequest.ID,
				SupplyCode:        itemReq.SupplyCode,
				SupplyName:        itemReq.SupplyName,
				QuantityRequested: itemReq.QuantityRequested,
				IsPediatric:       itemReq.IsPediatric,
			}

			if err := tx.Create(&item).Error; err != nil {
				return fmt.Errorf("error creando item de solicitud: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Retornar la solicitud completa con items
	return s.GetSupplyRequestByID(supplyRequest.ID)
}

// GetSupplyRequestByID obtiene una solicitud por ID con todos sus items y relaciones
func (s *SupplyRequestService) GetSupplyRequestByID(id int) (*models.SupplyRequest, error) {
	var request models.SupplyRequest

	if err := s.DB.First(&request, id).Error; err != nil {
		return nil, fmt.Errorf("solicitud no encontrada: %v", err)
	}

	// Obtener items con información del código de insumo
	var items []models.SupplyRequestItem
	if err := s.DB.Where("supply_request_id = ?", id).
		Preload("SupplyCodeInfo").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo items: %v", err)
	}

	// Obtener información del pabellón
	var pavilion models.Pavilion
	s.DB.First(&pavilion, request.PavilionID)

	// Obtener asignaciones de QR
	var qrAssignments []models.SupplyRequestQRAssignment
	s.DB.Where("supply_request_id = ?", id).
		Preload("MedicalSupply").
		Find(&qrAssignments)

	// Agregar información adicional al request (usando un mapa temporal)
	// Nota: En una implementación real, podrías crear un DTO para esto

	return &request, nil
}

// SupplyRequestWithItems representa una solicitud con información adicional calculada
type SupplyRequestWithItems struct {
	models.SupplyRequest
	TotalItems int64 `json:"total_items"`
}

// GetAllSupplyRequests obtiene todas las solicitudes con paginación
func (s *SupplyRequestService) GetAllSupplyRequests(limit, offset int, status string) ([]SupplyRequestWithItems, int64, error) {
	var requests []models.SupplyRequest
	var total int64

	query := s.DB.Model(&models.SupplyRequest{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Contar total
	query.Count(&total)

	// Obtener registros con paginación
	if err := query.Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	// Convertir a SupplyRequestWithItems y calcular total_items
	var requestsWithItems []SupplyRequestWithItems
	for _, request := range requests {
		var itemCount int64
		s.DB.Model(&models.SupplyRequestItem{}).
			Where("supply_request_id = ?", request.ID).
			Count(&itemCount)

		requestsWithItems = append(requestsWithItems, SupplyRequestWithItems{
			SupplyRequest: request,
			TotalItems:    itemCount,
		})
	}

	return requestsWithItems, total, nil
}

// GetSupplyRequestsByPavilion obtiene solicitudes por pabellón
func (s *SupplyRequestService) GetSupplyRequestsByPavilion(pavilionID int, limit, offset int) ([]SupplyRequestWithItems, int64, error) {
	var requests []models.SupplyRequest
	var total int64

	query := s.DB.Model(&models.SupplyRequest{}).Where("pavilion_id = ?", pavilionID)

	// Contar total
	query.Count(&total)

	// Obtener registros
	if err := query.Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	// Convertir a SupplyRequestWithItems y calcular total_items
	var requestsWithItems []SupplyRequestWithItems
	for _, request := range requests {
		var itemCount int64
		s.DB.Model(&models.SupplyRequestItem{}).
			Where("supply_request_id = ?", request.ID).
			Count(&itemCount)

		requestsWithItems = append(requestsWithItems, SupplyRequestWithItems{
			SupplyRequest: request,
			TotalItems:    itemCount,
		})
	}

	return requestsWithItems, total, nil
}

// ApproveSupplyRequest aprueba una solicitud y establece cantidades aprobadas
func (s *SupplyRequestService) ApproveSupplyRequest(requestID int, approval ApproveSupplyRequestRequest) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que la solicitud existe y puede ser aprobada
		var request models.SupplyRequest
		if err := tx.First(&request, requestID).Error; err != nil {
			return fmt.Errorf("solicitud no encontrada: %v", err)
		}

		if !request.CanBeApproved() {
			return fmt.Errorf("la solicitud no puede ser aprobada en su estado actual: %s", request.Status)
		}

		// Actualizar el estado de la solicitud
		now := time.Now()
		if err := tx.Model(&request).Updates(map[string]interface{}{
			"status":           models.RequestStatusApproved,
			"approved_by":      approval.ApprovedBy,
			"approved_by_name": approval.ApprovedByName,
			"approval_date":    now,
			"updated_at":       now,
		}).Error; err != nil {
			return fmt.Errorf("error actualizando solicitud: %v", err)
		}

		// Actualizar cantidades aprobadas para cada item
		for _, itemApproval := range approval.ItemApprovals {
			if err := tx.Model(&models.SupplyRequestItem{}).
				Where("id = ? AND supply_request_id = ?", itemApproval.ItemID, requestID).
				Update("quantity_approved", itemApproval.QuantityApproved).Error; err != nil {
				return fmt.Errorf("error actualizando item %d: %v", itemApproval.ItemID, err)
			}
		}

		return nil
	})
}

// RejectSupplyRequest rechaza una solicitud
func (s *SupplyRequestService) RejectSupplyRequest(requestID int, rejectedBy, rejectedByName, reason string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var request models.SupplyRequest
		if err := tx.First(&request, requestID).Error; err != nil {
			return fmt.Errorf("solicitud no encontrada: %v", err)
		}

		if !request.CanBeApproved() {
			return fmt.Errorf("la solicitud no puede ser rechazada en su estado actual: %s", request.Status)
		}

		now := time.Now()
		return tx.Model(&request).Updates(map[string]interface{}{
			"status":           models.RequestStatusRejected,
			"approved_by":      rejectedBy,
			"approved_by_name": rejectedByName,
			"approval_date":    now,
			"notes":            request.Notes + "\n\nMOTIVO DEL RECHAZO: " + reason,
			"updated_at":       now,
		}).Error
	})
}

// AssignQRToRequest asigna un código QR específico a un item de solicitud
func (s *SupplyRequestService) AssignQRToRequest(assignment AssignQRToRequestRequest) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que la solicitud existe y está aprobada
		var request models.SupplyRequest
		if err := tx.First(&request, assignment.SupplyRequestID).Error; err != nil {
			return fmt.Errorf("solicitud no encontrada: %v", err)
		}

		if !request.CanBeProcessed() {
			return fmt.Errorf("la solicitud no puede ser procesada en su estado actual: %s", request.Status)
		}

		// Verificar que el item existe
		var item models.SupplyRequestItem
		if err := tx.Where("id = ? AND supply_request_id = ?",
			assignment.SupplyRequestItemID, assignment.SupplyRequestID).
			First(&item).Error; err != nil {
			return fmt.Errorf("item de solicitud no encontrado: %v", err)
		}

		// Buscar el insumo médico por código QR
		var medicalSupply models.MedicalSupply
		if err := tx.Where("qr_code = ?", assignment.QRCode).First(&medicalSupply).Error; err != nil {
			return fmt.Errorf("insumo médico con QR %s no encontrado: %v", assignment.QRCode, err)
		}

		// Crear la asignación
		qrAssignment := models.SupplyRequestQRAssignment{
			SupplyRequestID:     assignment.SupplyRequestID,
			SupplyRequestItemID: assignment.SupplyRequestItemID,
			MedicalSupplyID:     medicalSupply.ID,
			AssignedDate:        time.Now(),
			AssignedBy:          assignment.AssignedBy,
			AssignedByName:      assignment.AssignedByName,
			Status:              models.AssignmentStatusAssigned,
			Notes:               assignment.Notes,
		}

		if err := tx.Create(&qrAssignment).Error; err != nil {
			return fmt.Errorf("error creando asignación QR: %v", err)
		}

		// Actualizar el estado de la solicitud a "en proceso" si es la primera asignación
		if err := tx.Model(&request).Update("status", models.RequestStatusInProcess).Error; err != nil {
			return fmt.Errorf("error actualizando estado de solicitud: %v", err)
		}

		// Incrementar el contador de items entregados
		if err := tx.Model(&item).Update("quantity_delivered", gorm.Expr("quantity_delivered + 1")).Error; err != nil {
			return fmt.Errorf("error actualizando cantidad entregada: %v", err)
		}

		return nil
	})
}

// GetQRTraceability obtiene la trazabilidad completa de un código QR
func (s *SupplyRequestService) GetQRTraceability(qrCode string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Buscar asignaciones del QR
	var assignments []models.SupplyRequestQRAssignment
	if err := s.DB.Where("qr_code = ?", qrCode).
		Preload("SupplyRequest").
		Preload("SupplyRequestItem").
		Preload("MedicalSupply").
		Order("assigned_date DESC").
		Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo trazabilidad: %v", err)
	}

	result["qr_code"] = qrCode
	result["assignments"] = assignments
	result["total_assignments"] = len(assignments)

	if len(assignments) > 0 {
		latest := assignments[0]
		result["current_status"] = latest.Status
		result["current_request"] = latest.SupplyRequest
		result["current_item"] = latest.SupplyRequestItem
		result["last_updated"] = latest.UpdatedAt
	}

	// Obtener historial del insumo médico
	if len(assignments) > 0 {
		var history []models.SupplyHistory
		s.DB.Where("medical_supply_id = ?", assignments[0].MedicalSupplyID).
			Order("date_time DESC").
			Find(&history)
		result["supply_history"] = history
	}

	return result, nil
}

// MarkQRAsDelivered marca un QR como entregado al pabellón
func (s *SupplyRequestService) MarkQRAsDelivered(qrCode, deliveredBy, deliveredByName string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var assignment models.SupplyRequestQRAssignment
		if err := tx.Where("qr_code = ? AND status = ?",
			qrCode, models.AssignmentStatusAssigned).
			First(&assignment).Error; err != nil {
			return fmt.Errorf("asignación no encontrada o ya procesada: %v", err)
		}

		now := time.Now()
		return tx.Model(&assignment).Updates(map[string]interface{}{
			"status":            models.AssignmentStatusDelivered,
			"delivered_date":    now,
			"delivered_by":      deliveredBy,
			"delivered_by_name": deliveredByName,
			"updated_at":        now,
		}).Error
	})
}

// GetSupplyRequestItems obtiene los items de una solicitud con información completa
func (s *SupplyRequestService) GetSupplyRequestItems(requestID int) ([]models.SupplyRequestItem, error) {
	var items []models.SupplyRequestItem

	if err := s.DB.Where("supply_request_id = ?", requestID).
		Preload("SupplyCodeInfo").
		Order("created_at ASC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo items: %v", err)
	}

	return items, nil
}

// GetSupplyRequestQRAssignments obtiene las asignaciones QR de una solicitud
func (s *SupplyRequestService) GetSupplyRequestQRAssignments(requestID int) ([]models.SupplyRequestQRAssignment, error) {
	var assignments []models.SupplyRequestQRAssignment

	if err := s.DB.Where("supply_request_id = ?", requestID).
		Preload("SupplyRequestItem").
		Preload("MedicalSupply").
		Order("assigned_date DESC").
		Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo asignaciones QR: %v", err)
	}

	return assignments, nil
}

// CompleteSupplyRequest marca una solicitud como completada
func (s *SupplyRequestService) CompleteSupplyRequest(requestID int, completedBy, completedByName string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var request models.SupplyRequest
		if err := tx.First(&request, requestID).Error; err != nil {
			return fmt.Errorf("solicitud no encontrada: %v", err)
		}

		// Verificar que todos los items aprobados han sido entregados
		var items []models.SupplyRequestItem
		if err := tx.Where("supply_request_id = ?", requestID).Find(&items).Error; err != nil {
			return fmt.Errorf("error obteniendo items: %v", err)
		}

		for _, item := range items {
			approvedQuantity := 0
			if item.QuantityApproved != nil {
				approvedQuantity = *item.QuantityApproved
			}

			if item.QuantityDelivered < approvedQuantity {
				return fmt.Errorf("no todos los items han sido entregados")
			}
		}

		now := time.Now()
		return tx.Model(&request).Updates(map[string]interface{}{
			"status":         models.RequestStatusCompleted,
			"completed_date": now,
			"updated_at":     now,
		}).Error
	})
}

// DeleteSupplyRequest elimina una solicitud (solo si está en estado pending)
func (s *SupplyRequestService) DeleteSupplyRequest(requestID int) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var request models.SupplyRequest
		if err := tx.First(&request, requestID).Error; err != nil {
			return fmt.Errorf("solicitud no encontrada: %v", err)
		}

		if !request.IsEditable() {
			return fmt.Errorf("la solicitud no puede ser eliminada en su estado actual: %s", request.Status)
		}

		// Eliminar items primero
		if err := tx.Where("supply_request_id = ?", requestID).Delete(&models.SupplyRequestItem{}).Error; err != nil {
			return fmt.Errorf("error eliminando items: %v", err)
		}

		// Eliminar asignaciones QR si existen
		if err := tx.Where("supply_request_id = ?", requestID).Delete(&models.SupplyRequestQRAssignment{}).Error; err != nil {
			return fmt.Errorf("error eliminando asignaciones QR: %v", err)
		}

		// Eliminar la solicitud
		return tx.Delete(&request).Error
	})
}

// GetSupplyRequestStats obtiene estadísticas de las solicitudes
func (s *SupplyRequestService) GetSupplyRequestStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Contar por estado
	var statusCounts []struct {
		Status string
		Count  int64
	}

	s.DB.Model(&models.SupplyRequest{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&statusCounts)

	statusMap := make(map[string]int64)
	for _, sc := range statusCounts {
		statusMap[sc.Status] = sc.Count
	}

	stats["by_status"] = statusMap

	// Contar total de solicitudes
	var total int64
	s.DB.Model(&models.SupplyRequest{}).Count(&total)
	stats["total_requests"] = total

	// Solicitudes del día actual
	today := time.Now().Format("2006-01-02")
	var todayCount int64
	s.DB.Model(&models.SupplyRequest{}).
		Where("DATE(created_at) = ?", today).
		Count(&todayCount)
	stats["today_requests"] = todayCount

	// Items más solicitados
	var topItems []struct {
		SupplyCode int    `json:"supply_code"`
		SupplyName string `json:"supply_name"`
		TotalQty   int64  `json:"total_quantity"`
	}

	s.DB.Model(&models.SupplyRequestItem{}).
		Select("supply_code, supply_name, sum(quantity_requested) as total_qty").
		Group("supply_code, supply_name").
		Order("total_qty DESC").
		Limit(10).
		Scan(&topItems)

	stats["top_requested_items"] = topItems

	stats["generated_at"] = time.Now()

	return stats, nil
}

// AssignRequestToWarehouseManager asigna una solicitud a un encargado de bodega (usado por Pavedad)
type AssignRequestToWarehouseManagerRequest struct {
	AssignedTo            string `json:"assigned_to" binding:"required"`              // RUT del encargado de bodega
	AssignedToName        string `json:"assigned_to_name" binding:"required"`         // Nombre del encargado
	AssignedByPavedad     string `json:"assigned_by_pavedad" binding:"required"`      // RUT del usuario Pavedad
	AssignedByPavedadName string `json:"assigned_by_pavedad_name" binding:"required"` // Nombre del usuario Pavedad
	PavedadNotes          string `json:"pavedad_notes"`                               // Notas de Pavedad
}

func (s *SupplyRequestService) AssignRequestToWarehouseManager(requestID int, req AssignRequestToWarehouseManagerRequest) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Obtener la solicitud
		var supplyRequest models.SupplyRequest
		if err := tx.First(&supplyRequest, requestID).Error; err != nil {
			return fmt.Errorf("solicitud no encontrada: %v", err)
		}

		// Validar que la solicitud está en estado pendiente_pavedad
		if supplyRequest.Status != "pendiente_pavedad" {
			return fmt.Errorf("la solicitud debe estar en estado 'pendiente_pavedad' para ser asignada. Estado actual: %s", supplyRequest.Status)
		}

		// Validar que el usuario asignado existe y es encargado de bodega
		var assignedUser models.User
		if err := tx.Where("rut = ?", req.AssignedTo).First(&assignedUser).Error; err != nil {
			return fmt.Errorf("usuario a asignar no encontrado: %v", err)
		}
		if assignedUser.Role != "encargado de bodega" {
			return fmt.Errorf("el usuario debe tener el rol 'encargado de bodega'. Rol actual: %s", assignedUser.Role)
		}

		// Validar que el usuario que asigna es Pavedad
		var pavedadUser models.User
		if err := tx.Where("rut = ?", req.AssignedByPavedad).First(&pavedadUser).Error; err != nil {
			return fmt.Errorf("usuario Pavedad no encontrado: %v", err)
		}
		if pavedadUser.Role != "pavedad" {
			return fmt.Errorf("solo usuarios con rol 'pavedad' pueden asignar solicitudes. Rol actual: %s", pavedadUser.Role)
		}

		// Actualizar la solicitud
		now := time.Now()
		supplyRequest.AssignedTo = &req.AssignedTo
		supplyRequest.AssignedToName = &req.AssignedToName
		supplyRequest.AssignedDate = &now
		supplyRequest.AssignedByPavedad = &req.AssignedByPavedad
		supplyRequest.AssignedByPavedadName = &req.AssignedByPavedadName
		if req.PavedadNotes != "" {
			supplyRequest.PavedadNotes = &req.PavedadNotes
		}
		supplyRequest.Status = "asignado_bodega"

		if err := tx.Save(&supplyRequest).Error; err != nil {
			return fmt.Errorf("error asignando solicitud: %v", err)
		}

		// Enviar correo al solicitante notificando la asignación
		go func() {
			if err := s.sendRequestAssignedEmail(&supplyRequest); err != nil {
				fmt.Printf("Error enviando correo de asignación: %v\n", err)
			}
		}()

		return nil
	})
}

// GetPendingRequestsForPavedad obtiene todas las solicitudes para Pavedad (pendientes y asignadas)
func (s *SupplyRequestService) GetPendingRequestsForPavedad() ([]models.SupplyRequest, error) {
	var requests []models.SupplyRequest
	// Pavedad ve tanto las pendientes como las ya asignadas
	err := s.DB.Where("status IN ?", []string{"pendiente_pavedad", "asignado_bodega", "en_proceso", "aprobado", "rechazado", "completado", "parcialmente_aprobado", "pendiente_revision", "devuelto"}).
		Order("request_date DESC").
		Find(&requests).Error

	if err != nil {
		return nil, fmt.Errorf("error obteniendo solicitudes para Pavedad: %v", err)
	}

	return requests, nil
}

// GetAssignedRequestsForWarehouseManager obtiene las solicitudes asignadas a un encargado de bodega específico
func (s *SupplyRequestService) GetAssignedRequestsForWarehouseManager(warehouseManagerRut string) ([]models.SupplyRequest, error) {
	var requests []models.SupplyRequest
	// Incluir todos los estados relevantes para el encargado de bodega
	err := s.DB.Where("assigned_to = ? AND status IN ?", warehouseManagerRut, []string{"asignado_bodega", "en_proceso", "aprobado", "rechazado", "parcialmente_aprobado", "pendiente_revision", "devuelto", "completado"}).
		Order("assigned_date DESC").
		Find(&requests).Error

	if err != nil {
		return nil, fmt.Errorf("error obteniendo solicitudes asignadas: %v", err)
	}

	return requests, nil
}

// ReviewItemRequest representa la estructura para revisar un item individual
type ReviewItemRequest struct {
	ItemStatus     string  `json:"item_status" binding:"required,oneof=aceptado rechazado devuelto"`
	ReviewedBy     string  `json:"reviewed_by" binding:"required"`
	ReviewedByName string  `json:"reviewed_by_name" binding:"required"`
	Comment        *string `json:"comment"` // Observaciones o comentarios sobre el item
}

// ReviewSupplyRequestItem permite al encargado de bodega revisar un item individual
func (s *SupplyRequestService) ReviewSupplyRequestItem(itemID int, req ReviewItemRequest) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var item models.SupplyRequestItem
		if err := tx.First(&item, itemID).Error; err != nil {
			return fmt.Errorf("item no encontrado: %v", err)
		}

		// Actualizar estado del item
		now := time.Now()
		item.ItemStatus = req.ItemStatus
		item.ReviewedBy = &req.ReviewedBy
		item.ReviewedByName = &req.ReviewedByName
		item.ReviewedAt = &now

		// Asignar comentario
		if req.Comment != nil && *req.Comment != "" {
			item.ItemNotes = req.Comment
		}

		if err := tx.Save(&item).Error; err != nil {
			return fmt.Errorf("error actualizando item: %v", err)
		}

		// Verificar si todos los items han sido revisados
		var totalItems, reviewedItems int64
		tx.Model(&models.SupplyRequestItem{}).Where("supply_request_id = ?", item.SupplyRequestID).Count(&totalItems)
		tx.Model(&models.SupplyRequestItem{}).Where("supply_request_id = ? AND item_status != ?", item.SupplyRequestID, "pendiente").Count(&reviewedItems)

		// Si todos los items fueron revisados, actualizar el estado de la solicitud
		if totalItems == reviewedItems {
			var allAccepted, hasReturned int64
			tx.Model(&models.SupplyRequestItem{}).Where("supply_request_id = ? AND item_status = ?", item.SupplyRequestID, "aceptado").Count(&allAccepted)
			tx.Model(&models.SupplyRequestItem{}).Where("supply_request_id = ? AND item_status = ?", item.SupplyRequestID, "devuelto").Count(&hasReturned)

			var request models.SupplyRequest
			if err := tx.First(&request, item.SupplyRequestID).Error; err != nil {
				return fmt.Errorf("error obteniendo solicitud: %v", err)
			}

			// Determinar nuevo estado de la solicitud
			if hasReturned > 0 {
				request.Status = "devuelto" // Al menos un item devuelto

				// Agregar comentarios de items devueltos al campo notes de la solicitud
				var returnedItems []models.SupplyRequestItem
				tx.Where("supply_request_id = ? AND item_status = ?", item.SupplyRequestID, "devuelto").Find(&returnedItems)

				var returnComments string
				for _, ri := range returnedItems {
					if ri.ItemNotes != nil && *ri.ItemNotes != "" {
						if returnComments != "" {
							returnComments += "\n"
						}
						returnComments += fmt.Sprintf("- %s: %s", ri.SupplyName, *ri.ItemNotes)
					}
				}

				if returnComments != "" {
					reviewerName := "Encargado de Bodega"
					if req.ReviewedByName != "" {
						reviewerName = req.ReviewedByName
					}

					notesHeader := fmt.Sprintf("\n\n[Devolución por %s]:\n", reviewerName)
					if request.Notes == "" {
						request.Notes = notesHeader + returnComments
					} else {
						request.Notes = request.Notes + notesHeader + returnComments
					}
				}

				// Enviar correo de devolución
				go func() {
					if err := s.sendRequestReturnedEmail(&request); err != nil {
						fmt.Printf("Error enviando correo de devolución: %v\n", err)
					}
				}()

			} else if allAccepted == totalItems {
				request.Status = "aprobado" // Todos aceptados

				// Enviar correo de aprobación
				go func() {
					if err := s.sendRequestApprovedEmail(&request); err != nil {
						fmt.Printf("Error enviando correo de aprobación: %v\n", err)
					}
				}()

			} else {
				request.Status = "parcialmente_aprobado" // Algunos rechazados

				// Enviar correo de aprobación parcial (usar template de aprobado)
				go func() {
					if err := s.sendRequestApprovedEmail(&request); err != nil {
						fmt.Printf("Error enviando correo de aprobación parcial: %v\n", err)
					}
				}()
			}

			if err := tx.Save(&request).Error; err != nil {
				return fmt.Errorf("error actualizando estado de solicitud: %v", err)
			}
		}

		return nil
	})
}

// UpdatedItemRequest representa la actualización de un item devuelto
type UpdatedItemRequest struct {
	ItemID   int `json:"item_id"`
	Quantity int `json:"quantity"`
}

// ResubmitReturnedRequest permite al doctor reenviar una solicitud devuelta
// Solo resetea los items devueltos a pendiente, manteniendo los aceptados
func (s *SupplyRequestService) ResubmitReturnedRequest(requestID int, updatedItems []UpdatedItemRequest) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var request models.SupplyRequest
		if err := tx.First(&request, requestID).Error; err != nil {
			return fmt.Errorf("solicitud no encontrada: %v", err)
		}

		// Verificar que la solicitud esté en estado devuelto
		if request.Status != "devuelto" {
			return fmt.Errorf("solo se pueden reenviar solicitudes devueltas")
		}

		// Actualizar los items según los cambios del doctor
		for _, updatedItem := range updatedItems {
			var item models.SupplyRequestItem
			if err := tx.First(&item, updatedItem.ItemID).Error; err != nil {
				continue // Si no se encuentra el item, continuar
			}

			// Solo actualizar items que fueron devueltos
			if item.ItemStatus == "devuelto" {
				// Actualizar cantidad si se modificó
				if updatedItem.Quantity > 0 {
					item.QuantityRequested = updatedItem.Quantity
				}

				// Resetear el estado a pendiente para que sea revisado nuevamente
				item.ItemStatus = "pendiente"
				item.ItemNotes = nil
				item.ReviewedBy = nil
				item.ReviewedByName = nil
				item.ReviewedAt = nil

				if err := tx.Save(&item).Error; err != nil {
					return fmt.Errorf("error actualizando item: %v", err)
				}
			}
			// Los items aceptados no se tocan, se mantienen como están
		}

		// Cambiar el estado de la solicitud de vuelta a asignado_bodega
		request.Status = "asignado_bodega"

		if err := tx.Save(&request).Error; err != nil {
			return fmt.Errorf("error actualizando solicitud: %v", err)
		}

		return nil
	})
}

// Helper functions para envío de correos

func (s *SupplyRequestService) sendRequestAssignedEmail(request *models.SupplyRequest) error {
	// Obtener email del solicitante
	var requester models.User
	if err := s.DB.Where("rut = ?", request.RequestedBy).First(&requester).Error; err != nil {
		return fmt.Errorf("error obteniendo solicitante: %v", err)
	}

	if requester.Email == "" {
		return fmt.Errorf("el solicitante no tiene email registrado")
	}

	// Obtener nombre del pabellón
	var pavilion models.Pavilion
	s.DB.First(&pavilion, request.PavilionID)

	templatePath := filepath.Join("mailer", "templates", "request_assigned.html")

	data := map[string]interface{}{
		"RecipientName":  request.RequestedByName,
		"RequestNumber":  request.RequestNumber,
		"PavilionName":   pavilion.Name,
		"SurgeryDate":    request.SurgeryDatetime.Format("02/01/2006 15:04"),
		"AssignedByName": *request.AssignedByPavedadName,
		"Notes":          request.Notes,
	}

	mailReq := mailer.NewRequest([]string{requester.Email}, "Solicitud Asignada a Bodega - "+request.RequestNumber)
	return mailReq.SendMailSkipTLS(templatePath, data)
}

func (s *SupplyRequestService) sendRequestApprovedEmail(request *models.SupplyRequest) error {
	// Obtener email del solicitante
	var requester models.User
	if err := s.DB.Where("rut = ?", request.RequestedBy).First(&requester).Error; err != nil {
		return fmt.Errorf("error obteniendo solicitante: %v", err)
	}

	if requester.Email == "" {
		return fmt.Errorf("el solicitante no tiene email registrado")
	}

	// Obtener nombre del pabellón
	var pavilion models.Pavilion
	s.DB.First(&pavilion, request.PavilionID)

	// Obtener items aprobados
	var items []models.SupplyRequestItem
	s.DB.Where("supply_request_id = ? AND item_status = ?", request.ID, "aceptado").Find(&items)

	itemsData := []map[string]interface{}{}
	for _, item := range items {
		itemsData = append(itemsData, map[string]interface{}{
			"SupplyName": item.SupplyName,
			"Quantity":   item.QuantityRequested,
		})
	}

	templatePath := filepath.Join("mailer", "templates", "request_approved.html")

	data := map[string]interface{}{
		"RecipientName":  request.RequestedByName,
		"RequestNumber":  request.RequestNumber,
		"PavilionName":   pavilion.Name,
		"SurgeryDate":    request.SurgeryDatetime.Format("02/01/2006 15:04"),
		"ApprovedByName": *request.AssignedToName,
		"Items":          itemsData,
		"Notes":          request.Notes,
	}

	mailReq := mailer.NewRequest([]string{requester.Email}, "Solicitud Aprobada - "+request.RequestNumber)
	return mailReq.SendMailSkipTLS(templatePath, data)
}

func (s *SupplyRequestService) sendRequestReturnedEmail(request *models.SupplyRequest) error {
	// Obtener email del solicitante
	var requester models.User
	if err := s.DB.Where("rut = ?", request.RequestedBy).First(&requester).Error; err != nil {
		return fmt.Errorf("error obteniendo solicitante: %v", err)
	}

	if requester.Email == "" {
		return fmt.Errorf("el solicitante no tiene email registrado")
	}

	// Obtener nombre del pabellón
	var pavilion models.Pavilion
	s.DB.First(&pavilion, request.PavilionID)

	// Obtener items devueltos con sus comentarios
	var items []models.SupplyRequestItem
	s.DB.Where("supply_request_id = ? AND item_status = ?", request.ID, "devuelto").Find(&items)

	returnedItems := []map[string]interface{}{}
	for _, item := range items {
		comments := "Sin comentarios"
		if item.ItemNotes != nil {
			comments = *item.ItemNotes
		}
		returnedItems = append(returnedItems, map[string]interface{}{
			"SupplyName": item.SupplyName,
			"Comments":   comments,
		})
	}

	templatePath := filepath.Join("mailer", "templates", "request_returned.html")

	data := map[string]interface{}{
		"RecipientName":  request.RequestedByName,
		"RequestNumber":  request.RequestNumber,
		"PavilionName":   pavilion.Name,
		"SurgeryDate":    request.SurgeryDatetime.Format("02/01/2006 15:04"),
		"ReturnedByName": *request.AssignedToName,
		"ReturnedItems":  returnedItems,
		"Notes":          request.Notes,
	}

	mailReq := mailer.NewRequest([]string{requester.Email}, "Solicitud Devuelta para Revisión - "+request.RequestNumber)
	return mailReq.SendMailSkipTLS(templatePath, data)
}

func (s *SupplyRequestService) sendRequestRejectedEmail(request *models.SupplyRequest) error {
	// Obtener email del solicitante
	var requester models.User
	if err := s.DB.Where("rut = ?", request.RequestedBy).First(&requester).Error; err != nil {
		return fmt.Errorf("error obteniendo solicitante: %v", err)
	}

	if requester.Email == "" {
		return fmt.Errorf("el solicitante no tiene email registrado")
	}

	// Obtener nombre del pabellón
	var pavilion models.Pavilion
	s.DB.First(&pavilion, request.PavilionID)

	templatePath := filepath.Join("mailer", "templates", "request_rejected.html")

	data := map[string]interface{}{
		"RecipientName":  request.RequestedByName,
		"RequestNumber":  request.RequestNumber,
		"PavilionName":   pavilion.Name,
		"SurgeryDate":    request.SurgeryDatetime.Format("02/01/2006 15:04"),
		"RejectedByName": *request.AssignedToName,
		"Notes":          request.Notes,
	}

	mailReq := mailer.NewRequest([]string{requester.Email}, "Solicitud Rechazada - "+request.RequestNumber)
	return mailReq.SendMailSkipTLS(templatePath, data)
}
