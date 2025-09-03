package services

import (
	"fmt"
	"meditrack/models"
	"time"

	"gorm.io/gorm"
)

type SupplyRequestService struct {
	DB *gorm.DB
}

func NewSupplyRequestService(db *gorm.DB) *SupplyRequestService {
	return &SupplyRequestService{DB: db}
}

// CreateSupplyRequestRequest representa la estructura para crear una solicitud
type CreateSupplyRequestRequest struct {
	PavilionID      int                              `json:"pavilion_id" binding:"required"`
	RequestedBy     string                           `json:"requested_by" binding:"required"`
	RequestedByName string                           `json:"requested_by_name" binding:"required"`
	Priority        string                           `json:"priority"`
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
			Status:          models.RequestStatusPending,
			Priority:        request.Priority,
			Notes:           request.Notes,
			MedicalCenterID: medicalCenterID,
		}

		// Validar y establecer prioridad por defecto
		if supplyRequest.Priority == "" {
			supplyRequest.Priority = models.RequestPriorityNormal
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
				Specifications:    itemReq.Specifications,
				IsPediatric:       itemReq.IsPediatric,
				SpecialRequests:   itemReq.SpecialRequests,
				UrgencyLevel:      itemReq.UrgencyLevel,
			}

			// Establecer valores opcionales
			if itemReq.Size != "" {
				item.Size = &itemReq.Size
			}
			if itemReq.Brand != "" {
				item.Brand = &itemReq.Brand
			}
			if item.UrgencyLevel == "" {
				item.UrgencyLevel = models.UrgencyLevelNormal
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

// GetAllSupplyRequests obtiene todas las solicitudes con paginación
func (s *SupplyRequestService) GetAllSupplyRequests(limit, offset int, status string) ([]models.SupplyRequest, int64, error) {
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

	return requests, total, nil
}

// GetSupplyRequestsByPavilion obtiene solicitudes por pabellón
func (s *SupplyRequestService) GetSupplyRequestsByPavilion(pavilionID int, limit, offset int) ([]models.SupplyRequest, int64, error) {
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

	return requests, total, nil
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

		// Verificar que el código QR existe y está disponible
		var medicalSupply models.MedicalSupply
		if err := tx.Where("qr_code = ?", assignment.QRCode).First(&medicalSupply).Error; err != nil {
			return fmt.Errorf("código QR no encontrado: %v", err)
		}

		// Verificar que el QR no esté ya asignado a otra solicitud
		var existingAssignment models.SupplyRequestQRAssignment
		if err := tx.Where("qr_code = ? AND status != ?",
			assignment.QRCode, models.AssignmentStatusConsumed).
			First(&existingAssignment).Error; err == nil {
			return fmt.Errorf("el código QR ya está asignado a otra solicitud")
		}

		// Verificar que el insumo corresponda al código solicitado
		var batch models.Batch
		if err := tx.First(&batch, medicalSupply.BatchID).Error; err != nil {
			return fmt.Errorf("lote no encontrado: %v", err)
		}

		// Verificar compatibilidad del código de insumo
		if medicalSupply.Code != item.SupplyCode {
			return fmt.Errorf("el código QR no corresponde al insumo solicitado")
		}

		// Crear la asignación
		qrAssignment := models.SupplyRequestQRAssignment{
			SupplyRequestID:     assignment.SupplyRequestID,
			SupplyRequestItemID: assignment.SupplyRequestItemID,
			QRCode:              assignment.QRCode,
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
