package controllers

import (
	"fmt"
	"meditrack/models"
	"meditrack/services"
	"net/http"
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
	var historyRequest struct {
		DateTime        time.Time `json:"date_time"`
		Status          string    `json:"status" binding:"required"`
		DestinationType string    `json:"destination_type" binding:"required"`
		DestinationID   int       `json:"destination_id" binding:"required"`
		MedicalSupplyID int       `json:"medical_supply_id" binding:"required"`
		UserRUT         string    `json:"user_rut" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&historyRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}

	// Crear modelo sin ID
	history := models.SupplyHistory{
		DateTime:        historyRequest.DateTime,
		Status:          historyRequest.Status,
		DestinationType: historyRequest.DestinationType,
		DestinationID:   historyRequest.DestinationID,
		MedicalSupplyID: historyRequest.MedicalSupplyID,
		UserRUT:         historyRequest.UserRUT,
		// ID se auto-generará
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

// GetAllSupplyHistory - obtiene todos los supply histories
func (c *SupplyHistoryController) GetAllSupplyHistory(ctx *gin.Context) {
	histories, err := c.supplyHistoryService.GetAllSupplyHistories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener supply histories: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: histories})
}

// UpdateSupplyHistory - actualiza un supply history por ID
func (c *SupplyHistoryController) UpdateSupplyHistory(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	var history models.SupplyHistory
	if err := ctx.ShouldBindJSON(&history); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}
	if err := c.supplyHistoryService.UpdateSupplyHistory(intID, &history); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al actualizar supply history: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Supply history actualizado", Data: history})
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
