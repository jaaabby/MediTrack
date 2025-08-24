package controllers

import (
	"fmt"
	"meditrack/models"
	"meditrack/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BatchHistoryController struct {
	batchHistoryService services.BatchHistoryService
}

func NewBatchHistoryController(batchHistoryService services.BatchHistoryService) *BatchHistoryController {
	return &BatchHistoryController{batchHistoryService: batchHistoryService}
}

func (c *BatchHistoryController) CreateBatchHistory(ctx *gin.Context) {
	var historyRequest struct {
		DateTime       time.Time `json:"date_time"`
		ChangeDetails  string    `json:"change_details" binding:"required"`
		PreviousValues string    `json:"previous_values"`
		NewValues      string    `json:"new_values"`
		UserName       string    `json:"user_name" binding:"required"`
		BatchID        *int      `json:"batch_id"`
		UserRUT        string    `json:"user_rut" binding:"required"`
		BatchNumber    int       `json:"batch_number" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&historyRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}

	// Crear modelo sin ID
	history := models.BatchHistory{
		DateTime:       historyRequest.DateTime,
		ChangeDetails:  historyRequest.ChangeDetails,
		PreviousValues: historyRequest.PreviousValues,
		NewValues:      historyRequest.NewValues,
		UserName:       historyRequest.UserName,
		BatchID:        historyRequest.BatchID,
		UserRUT:        historyRequest.UserRUT,
		BatchNumber:    historyRequest.BatchNumber,
		// ID se auto-generará
	}

	if err := c.batchHistoryService.CreateBatchHistory(&history); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al crear batch history: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, Response{Success: true, Message: "Batch history creado", Data: history})
}

func (c *BatchHistoryController) GetBatchHistoryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID invÃ¡lido: " + err.Error()})
		return
	}
	history, err := c.batchHistoryService.GetBatchHistoryByID(intID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{Success: false, Error: "Batch history no encontrado: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: history})
}

func (c *BatchHistoryController) GetAllBatchHistories(ctx *gin.Context) {
	histories, err := c.batchHistoryService.GetAllBatchHistories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener batch histories: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: histories})
}

func (c *BatchHistoryController) DeleteBatchHistory(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID invÃ¡lido: " + err.Error()})
		return
	}
	if err := c.batchHistoryService.DeleteBatchHistory(intID); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al eliminar batch history: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Batch history eliminado"})
}

func (c *BatchHistoryController) GetAllBatchHistoriesWithDetails(ctx *gin.Context) {
	histories, err := c.batchHistoryService.GetAllBatchHistoriesWithDetails()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener batch histories: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: histories})
}

// GetAllBatchHistoriesWithDetailsPaginated - MÃ©todo faltante agregado
func (c *BatchHistoryController) GetAllBatchHistoriesWithDetailsPaginated(ctx *gin.Context) {
	// Obtener parÃ¡metros de paginaciÃ³n
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	histories, total, err := c.batchHistoryService.GetAllBatchHistoriesWithDetailsPaginated(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener batch histories paginados: " + err.Error()})
		return
	}

	// Calcular informaciÃ³n de paginaciÃ³n
	totalPages := (int(total) + pageSize - 1) / pageSize

	response := map[string]interface{}{
		"data": histories,
		"pagination": map[string]interface{}{
			"current_page": page,
			"page_size":    pageSize,
			"total_items":  total,
			"total_pages":  totalPages,
			"has_next":     page < totalPages,
			"has_prev":     page > 1,
		},
	}

	ctx.JSON(http.StatusOK, Response{Success: true, Data: response})
}

// GetBatchHistoryByBatchNumber obtiene el historial de un lote especÃ­fico por su nÃºmero
func (c *BatchHistoryController) GetBatchHistoryByBatchNumber(ctx *gin.Context) {
	batchNumber := ctx.Param("batchNumber")
	var intBatchNumber int
	if _, err := fmt.Sscanf(batchNumber, "%d", &intBatchNumber); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "NÃºmero de lote invÃ¡lido: " + err.Error()})
		return
	}

	histories, err := c.batchHistoryService.GetBatchHistoryByBatchNumber(intBatchNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener historial del lote: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Response{Success: true, Data: histories})
}

// SearchBatchHistoryByBatchNumber busca en el historial por nÃºmero de lote con detalles formateados
// GetAllBatchHistory - obtiene todos los batch histories
func (c *BatchHistoryController) GetAllBatchHistory(ctx *gin.Context) {
	histories, err := c.batchHistoryService.GetAllBatchHistories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener batch histories: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: histories})
}

// UpdateBatchHistory - actualiza un batch history por ID
func (c *BatchHistoryController) UpdateBatchHistory(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	var history models.BatchHistory
	if err := ctx.ShouldBindJSON(&history); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}
	if err := c.batchHistoryService.UpdateBatchHistory(intID, &history); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al actualizar batch history: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Batch history actualizado", Data: history})
}
func (c *BatchHistoryController) SearchBatchHistoryByBatchNumber(ctx *gin.Context) {
	batchNumber := ctx.Param("batchNumber")
	var intBatchNumber int
	if _, err := fmt.Sscanf(batchNumber, "%d", &intBatchNumber); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "NÃºmero de lote invÃ¡lido: " + err.Error()})
		return
	}

	histories, err := c.batchHistoryService.SearchBatchHistoryByBatchNumber(intBatchNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al buscar historial del lote: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Response{Success: true, Data: histories})
}
