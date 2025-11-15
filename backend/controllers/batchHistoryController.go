package controllers

import (
	"meditrack/models"
	"meditrack/pkg/response"
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

// ========================
// CRUD BÁSICO
// ========================

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
		ctx.JSON(http.StatusBadRequest, response.Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}

	history := models.BatchHistory{
		DateTime:       historyRequest.DateTime,
		ChangeDetails:  historyRequest.ChangeDetails,
		PreviousValues: historyRequest.PreviousValues,
		NewValues:      historyRequest.NewValues,
		UserName:       historyRequest.UserName,
		BatchID:        historyRequest.BatchID,
		UserRUT:        historyRequest.UserRUT,
		BatchNumber:    historyRequest.BatchNumber,
	}

	if err := c.batchHistoryService.CreateBatchHistory(&history); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: "Error al crear batch history: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, response.Response{Success: true, Message: "Batch history creado", Data: history})
}

func (c *BatchHistoryController) GetBatchHistoryByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}

	history, err := c.batchHistoryService.GetBatchHistoryByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{Success: false, Error: "Batch history no encontrado: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{Success: true, Data: history})
}

func (c *BatchHistoryController) GetAllBatchHistories(ctx *gin.Context) {
	histories, err := c.batchHistoryService.GetAllBatchHistories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: "Error al obtener batch histories: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{Success: true, Data: histories})
}

func (c *BatchHistoryController) UpdateBatchHistory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}

	var history models.BatchHistory
	if err := ctx.ShouldBindJSON(&history); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}

	if err := c.batchHistoryService.UpdateBatchHistory(id, &history); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: "Error al actualizar batch history: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{Success: true, Message: "Batch history actualizado", Data: history})
}

func (c *BatchHistoryController) DeleteBatchHistory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}

	if err := c.batchHistoryService.DeleteBatchHistory(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: "Error al eliminar batch history: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{Success: true, Message: "Batch history eliminado"})
}

// ========================
// CONSULTAS ESPECIALIZADAS
// ========================

func (c *BatchHistoryController) GetAllBatchHistoriesWithDetails(ctx *gin.Context) {
	histories, err := c.batchHistoryService.GetAllBatchHistoriesWithDetails()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: "Error al obtener batch histories: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{Success: true, Data: histories})
}

func (c *BatchHistoryController) GetAllBatchHistoriesWithDetailsPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	histories, total, err := c.batchHistoryService.GetAllBatchHistoriesWithDetailsPaginated(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: "Error al obtener batch histories paginados: " + err.Error()})
		return
	}

	totalPages := (int(total) + pageSize - 1) / pageSize

	paginationData := map[string]interface{}{
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

	ctx.JSON(http.StatusOK, response.Response{Success: true, Data: paginationData})
}

func (c *BatchHistoryController) SearchBatchHistoryByBatchNumber(ctx *gin.Context) {
	batchNumber, err := strconv.Atoi(ctx.Param("batchNumber"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Success: false, Error: "Número de lote inválido: " + err.Error()})
		return
	}

	histories, err := c.batchHistoryService.SearchBatchHistoryByBatchNumber(batchNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: "Error al buscar historial del lote: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{Success: true, Data: histories})
}

// ELIMINADO: GetAllBatchHistory - duplicado de GetAllBatchHistories
// ELIMINADO: GetBatchHistoryByBatchNumber - usa SearchBatchHistoryByBatchNumber
