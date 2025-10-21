package controllers

import (
	"meditrack/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SupplyTransferController struct {
	transferService *services.SupplyTransferService
}

func NewSupplyTransferController(transferService *services.SupplyTransferService) *SupplyTransferController {
	return &SupplyTransferController{
		transferService: transferService,
	}
}

// Helper function to parse date with multiple formats
func parseFlexibleDateTransfer(dateStr string) (time.Time, error) {
	// Formatos soportados
	formats := []string{
		time.RFC3339,           // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05Z", // 2006-01-02T15:04:05Z
		"2006-01-02",           // 2006-01-02
		"2006-01-02 15:04:05",  // 2006-01-02 15:04:05
	}

	var lastErr error
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		} else {
			lastErr = err
		}
	}

	return time.Time{}, lastErr
}

// TransferToPavilion transfiere insumos de bodega a pabellón
func (c *SupplyTransferController) TransferToPavilion(ctx *gin.Context) {
	var request struct {
		QRCodes    []string `json:"qr_codes" binding:"required,min=1"`
		PavilionID int      `json:"pavilion_id" binding:"required"`
		Reason     string   `json:"reason" binding:"required"`
		Notes      string   `json:"notes"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	// Obtener información del usuario del contexto (viene del middleware de auth)
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "Usuario no autenticado",
		})
		return
	}

	userEmail, _ := ctx.Get("user_email")

	transfers, err := c.transferService.TransferToPavilion(
		request.QRCodes,
		request.PavilionID,
		userID.(string),    // UserID es el RUT
		userEmail.(string), // Usamos email en lugar de nombre
		request.Reason,
		request.Notes,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Error al transferir insumos",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Transferencia iniciada exitosamente",
		Data: gin.H{
			"transfers": transfers,
			"count":     len(transfers),
		},
	})
}

// ConfirmReception confirma la recepción de una transferencia
func (c *SupplyTransferController) ConfirmReception(ctx *gin.Context) {
	transferCode := ctx.Param("code")

	var request struct {
		Notes string `json:"notes"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	// Obtener información del usuario del contexto
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "Usuario no autenticado",
		})
		return
	}

	userEmail, _ := ctx.Get("user_email")

	transfer, err := c.transferService.ConfirmReception(
		transferCode,
		userID.(string),
		userEmail.(string),
		request.Notes,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Error al confirmar recepción",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Recepción confirmada exitosamente",
		Data:    transfer,
	})
}

// ReturnToStore devuelve insumos de pabellón a bodega
func (c *SupplyTransferController) ReturnToStore(ctx *gin.Context) {
	var request struct {
		QRCodes []string `json:"qr_codes" binding:"required,min=1"`
		StoreID int      `json:"store_id" binding:"required"`
		Reason  string   `json:"reason" binding:"required"`
		Notes   string   `json:"notes"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	// Obtener información del usuario del contexto
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "Usuario no autenticado",
		})
		return
	}

	userEmail, _ := ctx.Get("user_email")

	transfers, err := c.transferService.ReturnToStore(
		request.QRCodes,
		request.StoreID,
		userID.(string),
		userEmail.(string),
		request.Reason,
		request.Notes,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Error al devolver insumos",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Devolución procesada exitosamente",
		Data: gin.H{
			"transfers": transfers,
			"count":     len(transfers),
		},
	})
}

// GetTransferByCode obtiene una transferencia por su código
func (c *SupplyTransferController) GetTransferByCode(ctx *gin.Context) {
	transferCode := ctx.Param("code")

	transfer, err := c.transferService.GetTransferByCode(transferCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Transferencia no encontrada",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Transferencia encontrada",
		Data:    transfer,
	})
}

// GetTransfers obtiene listado de transferencias con filtros
func (c *SupplyTransferController) GetTransfers(ctx *gin.Context) {
	// Parsear parámetros de query
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	var status *string
	if s := ctx.Query("status"); s != "" {
		status = &s
	}

	var originType *string
	if ot := ctx.Query("origin_type"); ot != "" {
		originType = &ot
	}

	var originID *int
	if oidStr := ctx.Query("origin_id"); oidStr != "" {
		if oid, err := strconv.Atoi(oidStr); err == nil {
			originID = &oid
		}
	}

	var destinationType *string
	if dt := ctx.Query("destination_type"); dt != "" {
		destinationType = &dt
	}

	var destinationID *int
	if didStr := ctx.Query("destination_id"); didStr != "" {
		if did, err := strconv.Atoi(didStr); err == nil {
			destinationID = &did
		}
	}

	var startDate *time.Time
	if sdStr := ctx.Query("start_date"); sdStr != "" {
		if sd, err := parseFlexibleDateTransfer(sdStr); err == nil {
			startDate = &sd
		}
	}

	var endDate *time.Time
	if edStr := ctx.Query("end_date"); edStr != "" {
		if ed, err := parseFlexibleDateTransfer(edStr); err == nil {
			endDate = &ed
		}
	}

	transfers, total, err := c.transferService.GetTransfersByFilters(
		status,
		originType,
		originID,
		destinationType,
		destinationID,
		startDate,
		endDate,
		page,
		pageSize,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al obtener transferencias",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Transferencias obtenidas",
		Data: gin.H{
			"transfers":   transfers,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": (int(total) + pageSize - 1) / pageSize,
		},
	})
}

// CancelTransfer cancela una transferencia pendiente
func (c *SupplyTransferController) CancelTransfer(ctx *gin.Context) {
	transferCode := ctx.Param("code")

	var request struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	// Obtener información del usuario del contexto
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Message: "Usuario no autenticado",
		})
		return
	}

	transfer, err := c.transferService.CancelTransfer(
		transferCode,
		userID.(string),
		request.Reason,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Error al cancelar transferencia",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Transferencia cancelada exitosamente",
		Data:    transfer,
	})
}
