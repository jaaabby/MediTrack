package controllers

import (
	"net/http"

	"meditrack/models"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SupplyMovementController maneja las peticiones HTTP relacionadas con movimientos de insumos
type SupplyMovementController struct {
	supplyMovementService services.SupplyMovementService
}

// NewSupplyMovementController crea una nueva instancia de SupplyMovementController
func NewSupplyMovementController(supplyMovementService services.SupplyMovementService) *SupplyMovementController {
	return &SupplyMovementController{
		supplyMovementService: supplyMovementService,
	}
}

// CreateSupplyMovement crea un nuevo movimiento de insumo
func (c *SupplyMovementController) CreateSupplyMovement(ctx *gin.Context) {
	var movement models.SupplyMovement
	if err := ctx.ShouldBindJSON(&movement); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de movimiento inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyMovementService.CreateSupplyMovement(ctx, &movement); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear movimiento: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Movimiento creado exitosamente",
		Data:    movement,
	})
}

// GetSupplyMovementByID obtiene un movimiento por ID
func (c *SupplyMovementController) GetSupplyMovementByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de movimiento requerido",
		})
		return
	}

	movement, err := c.supplyMovementService.GetSupplyMovementByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Movimiento no encontrado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    movement,
	})
}

// GetAllSupplyMovements obtiene todos los movimientos
func (c *SupplyMovementController) GetAllSupplyMovements(ctx *gin.Context) {
	movements, err := c.supplyMovementService.GetAllSupplyMovements(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener movimientos: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    movements,
	})
}

// GetSupplyMovementsByStatus obtiene movimientos por estado
func (c *SupplyMovementController) GetSupplyMovementsByStatus(ctx *gin.Context) {
	status := ctx.Query("status")
	if status == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Estado requerido",
		})
		return
	}

	movements, err := c.supplyMovementService.GetSupplyMovementsByStatus(ctx, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener movimientos: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    movements,
	})
}

// UpdateSupplyMovement actualiza un movimiento existente
func (c *SupplyMovementController) UpdateSupplyMovement(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de movimiento requerido",
		})
		return
	}

	var movement models.SupplyMovement
	if err := ctx.ShouldBindJSON(&movement); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de movimiento inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyMovementService.UpdateSupplyMovement(ctx, &movement); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar movimiento: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Movimiento actualizado exitosamente",
		Data:    movement,
	})
}

// DeleteSupplyMovement elimina un movimiento
func (c *SupplyMovementController) DeleteSupplyMovement(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de movimiento requerido",
		})
		return
	}

	if err := c.supplyMovementService.DeleteSupplyMovement(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al eliminar movimiento: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Movimiento eliminado exitosamente",
	})
}
