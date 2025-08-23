package controllers

import (
	"fmt"
	"meditrack/models"
	"meditrack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BatchHistoryController struct {
	batchHistoryService services.BatchHistoryService
}

func NewBatchHistoryController(batchHistoryService services.BatchHistoryService) *BatchHistoryController {
	return &BatchHistoryController{batchHistoryService: batchHistoryService}
}

func (c *BatchHistoryController) CreateBatchHistory(ctx *gin.Context) {
	var history models.BatchHistory
	if err := ctx.ShouldBindJSON(&history); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
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
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
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
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
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

// GetBatchHistoryByBatchNumber obtiene el historial de un lote específico por su número
func (c *BatchHistoryController) GetBatchHistoryByBatchNumber(ctx *gin.Context) {
	batchNumber := ctx.Param("batchNumber")
	var intBatchNumber int
	if _, err := fmt.Sscanf(batchNumber, "%d", &intBatchNumber); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Número de lote inválido: " + err.Error()})
		return
	}

	histories, err := c.batchHistoryService.GetBatchHistoryByBatchNumber(intBatchNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener historial del lote: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Response{Success: true, Data: histories})
}

// SearchBatchHistoryByBatchNumber busca en el historial por número de lote con detalles formateados
func (c *BatchHistoryController) SearchBatchHistoryByBatchNumber(ctx *gin.Context) {
	batchNumber := ctx.Param("batchNumber")
	var intBatchNumber int
	if _, err := fmt.Sscanf(batchNumber, "%d", &intBatchNumber); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Número de lote inválido: " + err.Error()})
		return
	}

	histories, err := c.batchHistoryService.SearchBatchHistoryByBatchNumber(intBatchNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al buscar historial del lote: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Response{Success: true, Data: histories})
}
