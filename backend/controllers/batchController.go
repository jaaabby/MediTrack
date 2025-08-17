package controllers

import (
	"fmt"
	"meditrack/models"
	"meditrack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BatchController struct {
	batchService services.BatchService
}

func NewBatchController(batchService services.BatchService) *BatchController {
	return &BatchController{batchService: batchService}
}

func (c *BatchController) CreateBatch(ctx *gin.Context) {
	var batch models.Batch
	if err := ctx.ShouldBindJSON(&batch); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}
	if err := c.batchService.CreateBatch(&batch); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al crear batch: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, Response{Success: true, Message: "Batch creado", Data: batch})
}

func (c *BatchController) GetBatchByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	batch, err := c.batchService.GetBatchByID(intID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{Success: false, Error: "Batch no encontrado: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: batch})
}

func (c *BatchController) GetAllBatches(ctx *gin.Context) {
	batches, err := c.batchService.GetAllBatches()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener batches: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: batches})
}

func (c *BatchController) UpdateBatch(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	var batch models.Batch
	if err := ctx.ShouldBindJSON(&batch); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}
	batch.ID = intID
	if err := c.batchService.UpdateBatch(&batch); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al actualizar batch: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Batch actualizado", Data: batch})
}

func (c *BatchController) DeleteBatch(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	if err := c.batchService.DeleteBatch(intID); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al eliminar batch: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Batch eliminado"})
}
