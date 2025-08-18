package controllers

import (
	"fmt"
	"meditrack/models"
	"meditrack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MedicalSupplyController maneja las peticiones HTTP relacionadas con insumos médicos
type MedicalSupplyController struct {
	medicalSupplyService services.MedicalSupplyService
}

// NewMedicalSupplyController crea una nueva instancia de MedicalSupplyController
func NewMedicalSupplyController(medicalSupplyService services.MedicalSupplyService) *MedicalSupplyController {
	return &MedicalSupplyController{
		medicalSupplyService: medicalSupplyService,
	}
}

// CreateMedicalSupply crea un nuevo insumo médico
func (c *MedicalSupplyController) CreateMedicalSupply(ctx *gin.Context) {
	var supply models.MedicalSupply
	if err := ctx.ShouldBindJSON(&supply); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de insumo médico inválidos: " + err.Error(),
		})
		return
	}

	if err := c.medicalSupplyService.CreateMedicalSupply(&supply); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear insumo médico: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Insumo médico creado exitosamente",
		Data:    supply,
	})
}

// GetMedicalSupplyByID obtiene un insumo médico por ID
func (c *MedicalSupplyController) GetMedicalSupplyByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de insumo médico requerido",
		})
		return
	}

	// Convertir id a int
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID inválido: " + err.Error(),
		})
		return
	}

	supply, err := c.medicalSupplyService.GetMedicalSupplyByID(intID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Insumo médico no encontrado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    supply,
	})
}

// ...eliminar método GetMedicalSupplyByQRCode, no existe en el service CRUD...

// GetAllMedicalSupplies obtiene todos los insumos médicos
func (c *MedicalSupplyController) GetAllMedicalSupplies(ctx *gin.Context) {
	supplies, err := c.medicalSupplyService.GetAllMedicalSupplies()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener insumos médicos: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    supplies,
	})
}

// UpdateMedicalSupply actualiza un insumo médico existente
func (c *MedicalSupplyController) UpdateMedicalSupply(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de insumo médico requerido",
		})
		return
	}

	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID inválido: " + err.Error(),
		})
		return
	}

	var supply models.MedicalSupply
	if err := ctx.ShouldBindJSON(&supply); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de insumo médico inválidos: " + err.Error(),
		})
		return
	}

	supply.ID = intID

	if err := c.medicalSupplyService.UpdateMedicalSupply(&supply); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar insumo médico: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Insumo médico actualizado exitosamente",
		Data:    supply,
	})
}

// DeleteMedicalSupply elimina un insumo médico
func (c *MedicalSupplyController) DeleteMedicalSupply(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de insumo médico requerido",
		})
		return
	}

	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID inválido: " + err.Error(),
		})
		return
	}

	if err := c.medicalSupplyService.DeleteMedicalSupply(intID); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al eliminar insumo médico: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Insumo médico eliminado exitosamente",
	})
}
