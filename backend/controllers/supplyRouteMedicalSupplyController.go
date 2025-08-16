package controllers

import (
	"net/http"

	"meditrack/models"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SupplyRouteMedicalSupplyController maneja las peticiones HTTP relacionadas con la relación ruta de suministro-insumo médico
type SupplyRouteMedicalSupplyController struct {
	supplyRouteMedicalSupplyService services.SupplyRouteMedicalSupplyService
}

// NewSupplyRouteMedicalSupplyController crea una nueva instancia de SupplyRouteMedicalSupplyController
func NewSupplyRouteMedicalSupplyController(supplyRouteMedicalSupplyService services.SupplyRouteMedicalSupplyService) *SupplyRouteMedicalSupplyController {
	return &SupplyRouteMedicalSupplyController{
		supplyRouteMedicalSupplyService: supplyRouteMedicalSupplyService,
	}
}

// CreateSupplyRouteMedicalSupply crea una nueva relación ruta de suministro-insumo médico
func (c *SupplyRouteMedicalSupplyController) CreateSupplyRouteMedicalSupply(ctx *gin.Context) {
	var relation models.SupplyRouteMedicalSupply
	if err := ctx.ShouldBindJSON(&relation); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de relación ruta-insumo inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyRouteMedicalSupplyService.CreateSupplyRouteMedicalSupply(ctx, &relation); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear relación ruta-insumo: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Relación ruta-insumo creada exitosamente",
		Data:    relation,
	})
}

// GetSupplyRouteMedicalSupplyByID obtiene una relación ruta-insumo por ID
func (c *SupplyRouteMedicalSupplyController) GetSupplyRouteMedicalSupplyByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de relación ruta-insumo requerido",
		})
		return
	}

	relation, err := c.supplyRouteMedicalSupplyService.GetSupplyRouteMedicalSupplyByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Relación ruta-insumo no encontrada: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    relation,
	})
}

// GetSupplyRouteMedicalSuppliesBySupplyRouteID obtiene relaciones por ID de ruta de suministro
func (c *SupplyRouteMedicalSupplyController) GetSupplyRouteMedicalSuppliesBySupplyRouteID(ctx *gin.Context) {
	routeID := ctx.Query("supply_route_id")
	if routeID == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de ruta de suministro requerido",
		})
		return
	}

	relations, err := c.supplyRouteMedicalSupplyService.GetSupplyRouteMedicalSuppliesBySupplyRouteID(ctx, routeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener relaciones ruta-insumo: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    relations,
	})
}

// GetSupplyRouteMedicalSuppliesByMedicalSupplyID obtiene relaciones por ID de insumo médico
func (c *SupplyRouteMedicalSupplyController) GetSupplyRouteMedicalSuppliesByMedicalSupplyID(ctx *gin.Context) {
	supplyID := ctx.Query("medical_supply_id")
	if supplyID == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de insumo médico requerido",
		})
		return
	}

	relations, err := c.supplyRouteMedicalSupplyService.GetSupplyRouteMedicalSuppliesByMedicalSupplyID(ctx, supplyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener relaciones ruta-insumo: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    relations,
	})
}

// GetAllSupplyRouteMedicalSupplies obtiene todas las relaciones ruta-insumo
func (c *SupplyRouteMedicalSupplyController) GetAllSupplyRouteMedicalSupplies(ctx *gin.Context) {
	relations, err := c.supplyRouteMedicalSupplyService.GetAllSupplyRouteMedicalSupplies(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener relaciones ruta-insumo: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    relations,
	})
}

// UpdateSupplyRouteMedicalSupply actualiza una relación ruta-insumo existente
func (c *SupplyRouteMedicalSupplyController) UpdateSupplyRouteMedicalSupply(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de relación ruta-insumo requerido",
		})
		return
	}

	var relation models.SupplyRouteMedicalSupply
	if err := ctx.ShouldBindJSON(&relation); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de relación ruta-insumo inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyRouteMedicalSupplyService.UpdateSupplyRouteMedicalSupply(ctx, &relation); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar relación ruta-insumo: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Relación ruta-insumo actualizada exitosamente",
		Data:    relation,
	})
}

// DeleteSupplyRouteMedicalSupply elimina una relación ruta-insumo
func (c *SupplyRouteMedicalSupplyController) DeleteSupplyRouteMedicalSupply(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de relación ruta-insumo requerido",
		})
		return
	}

	if err := c.supplyRouteMedicalSupplyService.DeleteSupplyRouteMedicalSupply(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al eliminar relación ruta-insumo: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Relación ruta-insumo eliminada exitosamente",
	})
}
