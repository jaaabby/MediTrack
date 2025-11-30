package controllers

import (
	"meditrack/pkg/response"
	"meditrack/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type InventoryController struct {
	inventoryService *services.InventoryService
}

func NewInventoryController(inventoryService *services.InventoryService) *InventoryController {
	return &InventoryController{
		inventoryService: inventoryService,
	}
}

// Helper function to send success response
func sendSuccess(ctx *gin.Context, status int, message string, data interface{}) {
	ctx.JSON(status, response.Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Helper function to send error response
func sendError(ctx *gin.Context, status int, message string, err string) {
	ctx.JSON(status, response.Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// parseIntQuery parsea un parámetro de query como entero opcional
func parseIntQuery(ctx *gin.Context, paramName string) (*int, error) {
	valueStr := ctx.Query(paramName)
	if valueStr == "" {
		return nil, nil
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// parseStringQuery parsea un parámetro de query como string opcional
func parseStringQuery(ctx *gin.Context, paramName string) *string {
	value := ctx.Query(paramName)
	if value == "" {
		return nil
	}
	return &value
}

// Helper function to parse date with multiple formats
func parseFlexibleDate(dateStr string) (time.Time, error) {
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

// GetStoreInventory obtiene el inventario de bodegas con filtros
func (c *InventoryController) GetStoreInventory(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	storeID, _ := parseIntQuery(ctx, "store_id")
	surgeryID, _ := parseIntQuery(ctx, "surgery_id")
	supplyCode, _ := parseIntQuery(ctx, "supply_code")
	supplier := parseStringQuery(ctx, "supplier")

	nearExpiration := ctx.Query("near_expiration") == "true"
	lowStock := ctx.Query("low_stock") == "true"

	inventory, total, err := c.inventoryService.GetStoreInventory(
		storeID,
		surgeryID,
		supplyCode,
		supplier,
		nearExpiration,
		lowStock,
		page,
		pageSize,
	)

	if err != nil {
		sendError(ctx, http.StatusInternalServerError, "Error al obtener inventario de bodega", err.Error())
		return
	}

	sendSuccess(ctx, http.StatusOK, "Inventario de bodega obtenido", gin.H{
		"inventory":   inventory,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (int(total) + pageSize - 1) / pageSize,
	})
}

// parseIDFromParam parsea el ID del parámetro de la URL
func (c *InventoryController) parseIDFromParam(ctx *gin.Context, paramName string) (int, error) {
	id := ctx.Param(paramName)
	intID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return intID, nil
}

// GetPavilionInventory obtiene el inventario de un pabellón
func (c *InventoryController) GetPavilionInventory(ctx *gin.Context) {
	pavilionID, err := c.parseIDFromParam(ctx, "pavilion_id")
	if err != nil {
		sendError(ctx, http.StatusBadRequest, "ID de pabellón inválido", err.Error())
		return
	}

	includeInTransit := ctx.Query("include_in_transit") == "true"
	supplier := parseStringQuery(ctx, "supplier")

	inventory, err := c.inventoryService.GetPavilionInventory(pavilionID, includeInTransit, supplier)
	if err != nil {
		sendError(ctx, http.StatusInternalServerError, "Error al obtener inventario de pabellón", err.Error())
		return
	}

	sendSuccess(ctx, http.StatusOK, "Inventario de pabellón obtenido", gin.H{
		"pavilion_id": pavilionID,
		"inventory":   inventory,
		"count":       len(inventory),
	})
}

// GetMovementHistory obtiene el historial de movimientos
func (c *InventoryController) GetMovementHistory(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	locationType := parseStringQuery(ctx, "location_type")
	locationID, _ := parseIntQuery(ctx, "location_id")

	var startDate *time.Time
	if sdStr := ctx.Query("start_date"); sdStr != "" {
		if sd, err := parseFlexibleDate(sdStr); err == nil {
			startDate = &sd
		}
	}

	var endDate *time.Time
	if edStr := ctx.Query("end_date"); edStr != "" {
		if ed, err := parseFlexibleDate(edStr); err == nil {
			endDate = &ed
		}
	}

	movementType := parseStringQuery(ctx, "movement_type")

	history, total, err := c.inventoryService.GetMovementHistory(
		locationType,
		locationID,
		startDate,
		endDate,
		movementType,
		page,
		pageSize,
	)

	if err != nil {
		sendError(ctx, http.StatusInternalServerError, "Error al obtener historial de movimientos", err.Error())
		return
	}

	sendSuccess(ctx, http.StatusOK, "Historial de movimientos obtenido", gin.H{
		"history":     history,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (int(total) + pageSize - 1) / pageSize,
	})
}

// SyncInventory sincroniza los contadores de inventario
func (c *InventoryController) SyncInventory(ctx *gin.Context) {
	result, err := c.inventoryService.SyncInventory()
	if err != nil {
		sendError(ctx, http.StatusInternalServerError, "Error al sincronizar inventario", err.Error())
		return
	}

	sendSuccess(ctx, http.StatusOK, "Inventario sincronizado exitosamente", result)
}

// GetInventorySummary obtiene un resumen general del inventario
func (c *InventoryController) GetInventorySummary(ctx *gin.Context) {
	medicalCenterID, _ := parseIntQuery(ctx, "medical_center_id")

	summary, err := c.inventoryService.GetInventorySummary(medicalCenterID)
	if err != nil {
		sendError(ctx, http.StatusInternalServerError, "Error al obtener resumen de inventario", err.Error())
		return
	}

	sendSuccess(ctx, http.StatusOK, "Resumen de inventario obtenido", summary)
}

// GetInventoryBySurgeryType obtiene inventario agrupado por tipo de cirugía
func (c *InventoryController) GetInventoryBySurgeryType(ctx *gin.Context) {
	storeID, _ := parseIntQuery(ctx, "store_id")

	inventory, err := c.inventoryService.GetInventoryBySurgeryType(storeID)
	if err != nil {
		sendError(ctx, http.StatusInternalServerError, "Error al obtener inventario por tipo de cirugía", err.Error())
		return
	}

	sendSuccess(ctx, http.StatusOK, "Inventario por tipo de cirugía obtenido", gin.H{
		"inventory": inventory,
		"count":     len(inventory),
	})
}

// GetTransferReport obtiene un reporte de transferencias
func (c *InventoryController) GetTransferReport(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")
	groupBy := ctx.DefaultQuery("group_by", "date")

	if startDateStr == "" || endDateStr == "" {
		sendError(ctx, http.StatusBadRequest, "start_date y end_date son requeridos", "")
		return
	}

	startDate, err := parseFlexibleDate(startDateStr)
	if err != nil {
		sendError(ctx, http.StatusBadRequest, "Formato de start_date inválido (use YYYY-MM-DD o ISO 8601)", err.Error())
		return
	}

	endDate, err := parseFlexibleDate(endDateStr)
	if err != nil {
		sendError(ctx, http.StatusBadRequest, "Formato de end_date inválido (use YYYY-MM-DD o ISO 8601)", err.Error())
		return
	}

	report, err := c.inventoryService.GetTransferReport(startDate, endDate, groupBy)
	if err != nil {
		sendError(ctx, http.StatusInternalServerError, "Error al obtener reporte de transferencias", err.Error())
		return
	}

	sendSuccess(ctx, http.StatusOK, "Reporte de transferencias obtenido", gin.H{
		"report":     report,
		"start_date": startDateStr,
		"end_date":   endDateStr,
		"group_by":   groupBy,
	})
}
