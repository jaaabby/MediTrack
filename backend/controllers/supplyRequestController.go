package controllers

import (
	"fmt"
	"meditrack/services"
	"net/http"
	"strconv"

	"meditrack/pkg/response"

	"github.com/gin-gonic/gin"
)

type SupplyRequestController struct {
	supplyRequestService *services.SupplyRequestService
	qrService            *services.QRService
}

func NewSupplyRequestController(supplyRequestService *services.SupplyRequestService, qrService *services.QRService) *SupplyRequestController {
	return &SupplyRequestController{
		supplyRequestService: supplyRequestService,
		qrService:            qrService,
	}
}

// ...removed duplicate Response struct, use controllers.Response from common.go...

// CreateSupplyRequest crea una nueva solicitud de insumo
func (c *SupplyRequestController) CreateSupplyRequest(ctx *gin.Context) {
	var request services.CreateSupplyRequestRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	// Validar que hay al menos un item
	if len(request.Items) == 0 {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Debe incluir al menos un item en la solicitud",
		})
		return
	}

	supplyRequest, err := c.supplyRequestService.CreateSupplyRequest(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al crear solicitud: " + err.Error(),
		})
		return
	}

	// Obtener items de la solicitud creada
	items, _ := c.supplyRequestService.GetSupplyRequestItems(supplyRequest.ID)

	ctx.JSON(http.StatusCreated, response.Response{
		Success: true,
		Message: "Solicitud creada exitosamente",
		Data: map[string]interface{}{
			"request": supplyRequest,
			"items":   items,
		},
	})
}

// GetSupplyRequestByID obtiene una solicitud por ID
func (c *SupplyRequestController) GetSupplyRequestByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "ID de solicitud requerido",
		})
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "ID inválido: debe ser un número entero",
		})
		return
	}

	request, err := c.supplyRequestService.GetSupplyRequestByID(intID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   "Solicitud no encontrada: " + err.Error(),
		})
		return
	}

	// Obtener items y asignaciones QR
	items, _ := c.supplyRequestService.GetSupplyRequestItems(intID)
	assignments, _ := c.supplyRequestService.GetSupplyRequestQRAssignments(intID)

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Data: map[string]interface{}{
			"request":     request,
			"items":       items,
			"assignments": assignments,
		},
	})
}

// GetAllSupplyRequests obtiene todas las solicitudes con paginación y filtros
func (c *SupplyRequestController) GetAllSupplyRequests(ctx *gin.Context) {
	// Parámetros de paginación
	limit := 20
	offset := 0

	if l := ctx.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o := ctx.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	status := ctx.Query("status")

	requests, total, err := c.supplyRequestService.GetAllSupplyRequests(limit, offset, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al obtener solicitudes: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Data: map[string]interface{}{
			"requests": requests,
			"total":    total,
			"limit":    limit,
			"offset":   offset,
		},
	})
}

// GetSupplyRequestsByPavilion obtiene solicitudes por pabellón
func (c *SupplyRequestController) GetSupplyRequestsByPavilion(ctx *gin.Context) {
	pavilionIDStr := ctx.Param("pavilion_id")
	if pavilionIDStr == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "ID de pabellón requerido",
		})
		return
	}

	pavilionID, err := strconv.Atoi(pavilionIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de pabellón inválido",
		})
		return
	}

	// Parámetros de paginación
	limit := 20
	offset := 0

	if l := ctx.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o := ctx.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	requests, total, err := c.supplyRequestService.GetSupplyRequestsByPavilion(pavilionID, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener solicitudes del pabellón: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"requests":    requests,
			"total":       total,
			"limit":       limit,
			"offset":      offset,
			"pavilion_id": pavilionID,
		},
	})
}

// ApproveSupplyRequest aprueba una solicitud
func (c *SupplyRequestController) ApproveSupplyRequest(ctx *gin.Context) {
	// Verificar permisos - solo encargado de bodega puede aprobar
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Usuario no autenticado",
		})
		return
	}

	userClaims, ok := user.(map[string]interface{})
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error procesando información del usuario",
		})
		return
	}

	userRole, exists := userClaims["role"].(string)
	if !exists || userRole != "encargado de bodega" {
		ctx.JSON(http.StatusForbidden, Response{
			Success: false,
			Error:   "No tiene permisos para aprobar solicitudes. Solo encargados de bodega pueden realizar esta acción",
		})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de solicitud requerido",
		})
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID inválido",
		})
		return
	}

	var request services.ApproveSupplyRequestRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyRequestService.ApproveSupplyRequest(intID, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al aprobar solicitud: " + err.Error(),
		})
		return
	}

	// Obtener la solicitud actualizada
	updatedRequest, _ := c.supplyRequestService.GetSupplyRequestByID(intID)
	items, _ := c.supplyRequestService.GetSupplyRequestItems(intID)

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Solicitud aprobada exitosamente",
		Data: map[string]interface{}{
			"request": updatedRequest,
			"items":   items,
		},
	})
}

// RejectSupplyRequest rechaza una solicitud
func (c *SupplyRequestController) RejectSupplyRequest(ctx *gin.Context) {
	// Verificar permisos - solo encargado de bodega puede rechazar
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Usuario no autenticado",
		})
		return
	}

	userClaims, ok := user.(map[string]interface{})
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error procesando información del usuario",
		})
		return
	}

	userRole, exists := userClaims["role"].(string)
	if !exists || userRole != "encargado de bodega" {
		ctx.JSON(http.StatusForbidden, Response{
			Success: false,
			Error:   "No tiene permisos para rechazar solicitudes. Solo encargados de bodega pueden realizar esta acción",
		})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de solicitud requerido",
		})
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID inválido",
		})
		return
	}

	var requestData struct {
		RejectedBy     string `json:"rejected_by" binding:"required"`
		RejectedByName string `json:"rejected_by_name" binding:"required"`
		Reason         string `json:"reason" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyRequestService.RejectSupplyRequest(intID, requestData.RejectedBy, requestData.RejectedByName, requestData.Reason); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al rechazar solicitud: " + err.Error(),
		})
		return
	}

	// Obtener la solicitud actualizada
	updatedRequest, _ := c.supplyRequestService.GetSupplyRequestByID(intID)

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Solicitud rechazada",
		Data: map[string]interface{}{
			"request": updatedRequest,
		},
	})
}

// AssignQRToRequest asigna un código QR a un item de solicitud
func (c *SupplyRequestController) AssignQRToRequest(ctx *gin.Context) {
	var request services.AssignQRToRequestRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyRequestService.AssignQRToRequest(request); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al asignar código QR: " + err.Error(),
		})
		return
	}

	// Obtener información actualizada de la asignación
	traceability, _ := c.supplyRequestService.GetQRTraceability(request.QRCode)

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Código QR asignado exitosamente",
		Data: map[string]interface{}{
			"assignment":   request,
			"traceability": traceability,
		},
	})
}

// GetQRTraceability obtiene la trazabilidad completa de un código QR
func (c *SupplyRequestController) GetQRTraceability(ctx *gin.Context) {
	qrCode := ctx.Param("qr_code")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	traceability, err := c.supplyRequestService.GetQRTraceability(qrCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener trazabilidad: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    traceability,
	})
}

// MarkQRAsDelivered marca un QR como entregado
func (c *SupplyRequestController) MarkQRAsDelivered(ctx *gin.Context) {
	qrCode := ctx.Param("qr_code")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	var requestData struct {
		DeliveredBy     string `json:"delivered_by" binding:"required"`
		DeliveredByName string `json:"delivered_by_name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyRequestService.MarkQRAsDelivered(qrCode, requestData.DeliveredBy, requestData.DeliveredByName); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al marcar QR como entregado: " + err.Error(),
		})
		return
	}

	// Obtener trazabilidad actualizada
	traceability, _ := c.supplyRequestService.GetQRTraceability(qrCode)

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "QR marcado como entregado exitosamente",
		Data:    traceability,
	})
}

// CompleteSupplyRequest marca una solicitud como completada
func (c *SupplyRequestController) CompleteSupplyRequest(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de solicitud requerido",
		})
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID inválido",
		})
		return
	}

	var requestData struct {
		CompletedBy     string `json:"completed_by" binding:"required"`
		CompletedByName string `json:"completed_by_name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyRequestService.CompleteSupplyRequest(intID, requestData.CompletedBy, requestData.CompletedByName); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al completar solicitud: " + err.Error(),
		})
		return
	}

	// Obtener la solicitud actualizada
	updatedRequest, _ := c.supplyRequestService.GetSupplyRequestByID(intID)

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Solicitud completada exitosamente",
		Data: map[string]interface{}{
			"request": updatedRequest,
		},
	})
}

// DeleteSupplyRequest elimina una solicitud
func (c *SupplyRequestController) DeleteSupplyRequest(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de solicitud requerido",
		})
		return
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID inválido",
		})
		return
	}

	if err := c.supplyRequestService.DeleteSupplyRequest(intID); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al eliminar solicitud: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Solicitud eliminada exitosamente",
	})
}

// GetSupplyRequestStats obtiene estadísticas de solicitudes
func (c *SupplyRequestController) GetSupplyRequestStats(ctx *gin.Context) {
	stats, err := c.supplyRequestService.GetSupplyRequestStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener estadísticas: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    stats,
	})
}

// BulkAssignQRs asigna múltiples códigos QR de una sola vez
func (c *SupplyRequestController) BulkAssignQRs(ctx *gin.Context) {
	var requestData struct {
		Assignments []services.AssignQRToRequestRequest `json:"assignments" binding:"required,dive"`
	}

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	var successCount int
	var errors []string

	for _, assignment := range requestData.Assignments {
		if err := c.supplyRequestService.AssignQRToRequest(assignment); err != nil {
			errors = append(errors, fmt.Sprintf("QR %s: %s", assignment.QRCode, err.Error()))
		} else {
			successCount++
		}
	}

	response := map[string]interface{}{
		"total_processed":    len(requestData.Assignments),
		"successful_assigns": successCount,
		"errors":             errors,
	}

	if len(errors) > 0 {
		ctx.JSON(http.StatusPartialContent, Response{
			Success: successCount > 0,
			Message: fmt.Sprintf("%d de %d asignaciones completadas exitosamente", successCount, len(requestData.Assignments)),
			Data:    response,
		})
	} else {
		ctx.JSON(http.StatusOK, Response{
			Success: true,
			Message: "Todas las asignaciones completadas exitosamente",
			Data:    response,
		})
	}
}

// AssignRequestToWarehouseManager asigna una solicitud a un encargado de bodega (usado por Pavedad)
func (c *SupplyRequestController) AssignRequestToWarehouseManager(ctx *gin.Context) {
	requestID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de solicitud inválido",
		})
		return
	}

	var request services.AssignRequestToWarehouseManagerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyRequestService.AssignRequestToWarehouseManager(requestID, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Solicitud asignada exitosamente a encargado de bodega",
	})
}

// GetPendingRequestsForPavedad obtiene todas las solicitudes pendientes de asignación por Pavedad
func (c *SupplyRequestController) GetPendingRequestsForPavedad(ctx *gin.Context) {
	requests, err := c.supplyRequestService.GetPendingRequestsForPavedad()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    requests,
	})
}

// GetAssignedRequestsForWarehouseManager obtiene las solicitudes asignadas a un encargado de bodega
func (c *SupplyRequestController) GetAssignedRequestsForWarehouseManager(ctx *gin.Context) {
	warehouseManagerRut := ctx.Param("rut")

	requests, err := c.supplyRequestService.GetAssignedRequestsForWarehouseManager(warehouseManagerRut)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    requests,
	})
}

// ReviewSupplyRequestItem permite revisar un item individual (aceptar/rechazar/devolver)
func (c *SupplyRequestController) ReviewSupplyRequestItem(ctx *gin.Context) {
	itemIDStr := ctx.Param("itemId")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de item inválido",
		})
		return
	}

	var req services.ReviewItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyRequestService.ReviewSupplyRequestItem(itemID, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Item revisado exitosamente",
	})
}

// GetSupplyRequestItemsController obtiene todos los items de una solicitud
func (c *SupplyRequestController) GetSupplyRequestItemsController(ctx *gin.Context) {
	requestIDStr := ctx.Param("id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de solicitud inválido",
		})
		return
	}

	items, err := c.supplyRequestService.GetSupplyRequestItems(requestID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    items,
	})
}

// ResubmitReturnedRequest permite al doctor reenviar una solicitud devuelta
func (c *SupplyRequestController) ResubmitReturnedRequest(ctx *gin.Context) {
	requestIDStr := ctx.Param("id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de solicitud inválido",
		})
		return
	}

	var req struct {
		UpdatedItems []services.UpdatedItemRequest `json:"updated_items"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyRequestService.ResubmitReturnedRequest(requestID, req.UpdatedItems); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Solicitud reenviada exitosamente",
	})
}
