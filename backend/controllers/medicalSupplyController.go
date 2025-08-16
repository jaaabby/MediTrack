package controllers

import (
	"net/http"

	"meditrack/models"
	"meditrack/services"

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

	if err := c.medicalSupplyService.CreateMedicalSupply(ctx, &supply); err != nil {
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

	supply, err := c.medicalSupplyService.GetMedicalSupplyByID(ctx, id)
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

// GetMedicalSupplyByQRCode obtiene un insumo médico por código QR
func (c *MedicalSupplyController) GetMedicalSupplyByQRCode(ctx *gin.Context) {
	qrCode := ctx.Query("qr_code")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	supply, err := c.medicalSupplyService.GetMedicalSupplyByQRCode(ctx, qrCode)
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

// GetAllMedicalSupplies obtiene todos los insumos médicos
func (c *MedicalSupplyController) GetAllMedicalSupplies(ctx *gin.Context) {
	supplies, err := c.medicalSupplyService.GetAllMedicalSupplies(ctx)
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

	var supply models.MedicalSupply
	if err := ctx.ShouldBindJSON(&supply); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de insumo médico inválidos: " + err.Error(),
		})
		return
	}

	if err := c.medicalSupplyService.UpdateMedicalSupply(ctx, &supply); err != nil {
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

	if err := c.medicalSupplyService.DeleteMedicalSupply(ctx, id); err != nil {
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
