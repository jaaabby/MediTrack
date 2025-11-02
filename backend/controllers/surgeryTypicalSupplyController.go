package controllers

import (
	"meditrack/models"
	"meditrack/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SurgeryTypicalSupplyController struct {
	typicalSupplyService *services.SurgeryTypicalSupplyService
}

func NewSurgeryTypicalSupplyController(typicalSupplyService *services.SurgeryTypicalSupplyService) *SurgeryTypicalSupplyController {
	return &SurgeryTypicalSupplyController{
		typicalSupplyService: typicalSupplyService,
	}
}

// CreateSurgeryTypicalSupply crea un nuevo insumo típico para una cirugía
func (c *SurgeryTypicalSupplyController) CreateSurgeryTypicalSupply(ctx *gin.Context) {
	var typicalSupply models.SurgeryTypicalSupply

	if err := ctx.ShouldBindJSON(&typicalSupply); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	if err := c.typicalSupplyService.CreateSurgeryTypicalSupply(&typicalSupply); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al crear insumo típico",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Insumo típico creado exitosamente",
		Data:    typicalSupply,
	})
}

// GetSurgeryTypicalSupplyByID obtiene un insumo típico por ID
func (c *SurgeryTypicalSupplyController) GetSurgeryTypicalSupplyByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID inválido",
			Error:   err.Error(),
		})
		return
	}

	typicalSupply, err := c.typicalSupplyService.GetSurgeryTypicalSupplyByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Insumo típico no encontrado",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Insumo típico encontrado",
		Data:    typicalSupply,
	})
}

// GetTypicalSuppliesBySurgeryID obtiene todos los insumos típicos de una cirugía
func (c *SurgeryTypicalSupplyController) GetTypicalSuppliesBySurgeryID(ctx *gin.Context) {
	surgeryID, err := strconv.Atoi(ctx.Param("surgery_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID de cirugía inválido",
			Error:   err.Error(),
		})
		return
	}

	typicalSupplies, err := c.typicalSupplyService.GetTypicalSuppliesBySurgeryID(surgeryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al obtener insumos típicos",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Insumos típicos obtenidos",
		Data: gin.H{
			"typical_supplies": typicalSupplies,
			"count":            len(typicalSupplies),
		},
	})
}

// GetSurgeriesBySupplyCode obtiene todas las cirugías que requieren un insumo específico
func (c *SurgeryTypicalSupplyController) GetSurgeriesBySupplyCode(ctx *gin.Context) {
	supplyCode, err := strconv.Atoi(ctx.Param("supply_code"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Código de insumo inválido",
			Error:   err.Error(),
		})
		return
	}

	typicalSupplies, err := c.typicalSupplyService.GetSurgeriesBySupplyCode(supplyCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al obtener cirugías",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Cirugías obtenidas",
		Data: gin.H{
			"surgeries": typicalSupplies,
			"count":  len(typicalSupplies),
		},
	})
}

// UpdateSurgeryTypicalSupply actualiza un insumo típico
func (c *SurgeryTypicalSupplyController) UpdateSurgeryTypicalSupply(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID inválido",
			Error:   err.Error(),
		})
		return
	}

	var typicalSupply models.SurgeryTypicalSupply
	if err := ctx.ShouldBindJSON(&typicalSupply); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	updatedTypicalSupply, err := c.typicalSupplyService.UpdateSurgeryTypicalSupply(id, &typicalSupply)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al actualizar insumo típico",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Insumo típico actualizado exitosamente",
		Data:    updatedTypicalSupply,
	})
}

// DeleteSurgeryTypicalSupply elimina un insumo típico
func (c *SurgeryTypicalSupplyController) DeleteSurgeryTypicalSupply(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID inválido",
			Error:   err.Error(),
		})
		return
	}

	if err := c.typicalSupplyService.DeleteSurgeryTypicalSupply(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al eliminar insumo típico",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Insumo típico eliminado exitosamente",
	})
}

// BulkCreateSurgeryTypicalSupplies crea múltiples insumos típicos para una cirugía
func (c *SurgeryTypicalSupplyController) BulkCreateSurgeryTypicalSupplies(ctx *gin.Context) {
	surgeryID, err := strconv.Atoi(ctx.Param("surgery_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID de cirugía inválido",
			Error:   err.Error(),
		})
		return
	}

	var typicalSupplies []models.SurgeryTypicalSupply
	if err := ctx.ShouldBindJSON(&typicalSupplies); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	if err := c.typicalSupplyService.BulkCreateSurgeryTypicalSupplies(surgeryID, typicalSupplies); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al crear insumos típicos",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Insumos típicos creados exitosamente",
		Data: gin.H{
			"created_count": len(typicalSupplies),
		},
	})
}

// GetTypicalSuppliesCount obtiene el conteo total de insumos típicos
func (c *SurgeryTypicalSupplyController) GetTypicalSuppliesCount(ctx *gin.Context) {
	count, err := c.typicalSupplyService.GetTypicalSuppliesCount()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al obtener conteo de insumos típicos",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Conteo de insumos típicos obtenido",
		Data: gin.H{
			"count": count,
		},
	})
}

