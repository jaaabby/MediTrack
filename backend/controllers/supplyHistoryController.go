package controllers

import (
	"meditrack/models"
	"meditrack/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SupplyHistoryController struct {
	supplyHistoryService services.SupplyHistoryService
}

func NewSupplyHistoryController(supplyHistoryService services.SupplyHistoryService) *SupplyHistoryController {
	return &SupplyHistoryController{supplyHistoryService: supplyHistoryService}
}

func (c *SupplyHistoryController) CreateSupplyHistory(ctx *gin.Context) {
	var history models.SupplyHistory

	if err := ctx.ShouldBindJSON(&history); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}

	// Si no se proporciona fecha, usar la fecha actual
	if history.DateTime.IsZero() {
		history.DateTime = time.Now()
	}

	if history.Status == "" {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "El estado es requerido"})
		return
	}

	if err := c.supplyHistoryService.CreateSupplyHistory(&history); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al crear supply history: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, Response{Success: true, Message: "Supply history creado", Data: history})
}

func (c *SupplyHistoryController) GetSupplyHistoryByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	history, err := c.supplyHistoryService.GetSupplyHistoryByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{Success: false, Error: "Supply history no encontrado: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: history})
}

func (c *SupplyHistoryController) DeleteSupplyHistory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	if err := c.supplyHistoryService.DeleteSupplyHistory(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al eliminar supply history: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Supply history eliminado"})
}

func (c *SupplyHistoryController) GetAllSupplyHistory(ctx *gin.Context) {
	histories, err := c.supplyHistoryService.GetAllSupplyHistories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener supply histories: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: histories})
}

func (c *SupplyHistoryController) UpdateSupplyHistory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	var history models.SupplyHistory
	if err := ctx.ShouldBindJSON(&history); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}
	if err := c.supplyHistoryService.UpdateSupplyHistory(id, &history); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al actualizar supply history: " + err.Error()})
		return
	}
	updatedHistory, err := c.supplyHistoryService.GetSupplyHistoryByID(id)
	if err != nil {
		ctx.JSON(http.StatusOK, Response{Success: true, Message: "Supply history actualizado"})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Supply history actualizado", Data: updatedHistory})
}

// GetAllSupplyHistoriesWithDetails - obtiene todos los supply histories con información detallada del insumo
func (c *SupplyHistoryController) GetAllSupplyHistoriesWithDetails(ctx *gin.Context) {
	histories, err := c.supplyHistoryService.GetAllSupplyHistoriesWithDetails()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener supply histories: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: histories})
}
