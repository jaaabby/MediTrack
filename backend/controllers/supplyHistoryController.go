package controllers

import (
	"fmt"
	"meditrack/models"
	"meditrack/services"
	"net/http"

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
	if err := c.supplyHistoryService.CreateSupplyHistory(&history); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al crear supply history: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, Response{Success: true, Message: "Supply history creado", Data: history})
}

func (c *SupplyHistoryController) GetSupplyHistoryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	history, err := c.supplyHistoryService.GetSupplyHistoryByID(intID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{Success: false, Error: "Supply history no encontrado: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: history})
}

func (c *SupplyHistoryController) GetAllSupplyHistories(ctx *gin.Context) {
	histories, err := c.supplyHistoryService.GetAllSupplyHistories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener supply histories: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: histories})
}

func (c *SupplyHistoryController) DeleteSupplyHistory(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	if err := c.supplyHistoryService.DeleteSupplyHistory(intID); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al eliminar supply history: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Supply history eliminado"})
}
