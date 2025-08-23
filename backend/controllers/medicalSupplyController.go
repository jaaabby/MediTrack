package controllers

import (
	"meditrack/models"
	"meditrack/services"
	"net/http"
	"strconv"

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

// CreateMedicalSupply crea un nuevo insumo médico individual
func (c *MedicalSupplyController) CreateMedicalSupply(ctx *gin.Context) {
	var supplyRequest struct {
		Code    int `json:"code" binding:"required"`
		BatchID int `json:"batch_id" binding:"required"`
		// No incluimos QRCode ni ID - se generarán automáticamente
	}

	if err := ctx.ShouldBindJSON(&supplyRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de insumo médico inválidos: " + err.Error(),
		})
		return
	}

	// Crear el modelo sin ID - se auto-generará
	supply := models.MedicalSupply{
		Code:    supplyRequest.Code,
		BatchID: supplyRequest.BatchID,
		// ID se auto-generará
		// QRCode se generará en el servicio
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

// CreateMultipleSupplies crea múltiples insumos individuales para un lote
func (c *MedicalSupplyController) CreateMultipleSupplies(ctx *gin.Context) {
	var request struct {
		BatchID  int `json:"batch_id" binding:"required"`
		Code     int `json:"code" binding:"required"`
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	supplies, err := c.medicalSupplyService.CreateMultipleIndividualSupplies(
		request.BatchID,
		request.Code,
		request.Quantity,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear insumos individuales: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Insumos individuales creados exitosamente",
		Data: map[string]interface{}{
			"created_supplies": supplies,
			"total_created":    len(supplies),
			"batch_id":         request.BatchID,
		},
	})
}

// ConsumeSupply consume un insumo por su código QR
func (c *MedicalSupplyController) ConsumeSupply(ctx *gin.Context) {
	var request struct {
		QRCode          string `json:"qr_code" binding:"required"`
		UserRUT         string `json:"user_rut" binding:"required"`
		DestinationType string `json:"destination_type" binding:"required"`
		DestinationID   int    `json:"destination_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de consumo inválidos: " + err.Error(),
		})
		return
	}

	supply, err := c.medicalSupplyService.ConsumeSupplyByQR(
		request.QRCode,
		request.UserRUT,
		request.DestinationType,
		request.DestinationID,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al consumir insumo: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Insumo consumido exitosamente",
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
	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID inválido: debe ser un número entero",
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

// GetMedicalSupplyByQR obtiene un insumo médico por su código QR
func (c *MedicalSupplyController) GetMedicalSupplyByQR(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	supply, err := c.medicalSupplyService.GetMedicalSupplyByQR(qrCode)
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

// GetSupplyWithBatchInfo obtiene información completa de un insumo con datos del lote
func (c *MedicalSupplyController) GetSupplyWithBatchInfo(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	info, err := c.medicalSupplyService.GetSupplyWithBatchInfo(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Información del insumo obtenida exitosamente",
		Data:    info,
	})
}

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

// GetInventoryList obtiene el inventario con información actualizada de productos consumidos
func (c *MedicalSupplyController) GetInventoryList(ctx *gin.Context) {
	inventory, err := c.medicalSupplyService.GetInventoryList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener inventario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Inventario obtenido exitosamente",
		Data: map[string]interface{}{
			"inventory_items": inventory,
			"total_items":     len(inventory),
		},
	})
}

// GetIndividualSuppliesByCode obtiene todos los insumos individuales de un código específico
func (c *MedicalSupplyController) GetIndividualSuppliesByCode(ctx *gin.Context) {
	codeStr := ctx.Param("code")
	if codeStr == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código de insumo requerido",
		})
		return
	}

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código inválido: debe ser un número entero",
		})
		return
	}

	supplies, err := c.medicalSupplyService.GetIndividualSuppliesByCode(code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener insumos individuales: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Insumos individuales obtenidos exitosamente",
		Data: map[string]interface{}{
			"supplies":    supplies,
			"total_count": len(supplies),
			"supply_code": code,
		},
	})
}

// GetAvailableSuppliesByBatch obtiene todos los insumos disponibles de un lote
func (c *MedicalSupplyController) GetAvailableSuppliesByBatch(ctx *gin.Context) {
	batchIDStr := ctx.Param("batch_id")
	if batchIDStr == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de lote requerido",
		})
		return
	}

	batchID, err := strconv.Atoi(batchIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de lote inválido: debe ser un número entero",
		})
		return
	}

	supplies, err := c.medicalSupplyService.GetAvailableSuppliesByBatch(batchID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener insumos disponibles: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Insumos disponibles obtenidos exitosamente",
		Data: map[string]interface{}{
			"available_supplies": supplies,
			"total_available":    len(supplies),
			"batch_id":           batchID,
		},
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

	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID inválido: debe ser un número entero",
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

	updatedSupply, err := c.medicalSupplyService.UpdateMedicalSupply(intID, &supply)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar insumo médico: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Insumo médico actualizado exitosamente",
		Data:    updatedSupply,
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

	intID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID inválido: debe ser un número entero",
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

// SyncBatchAmounts sincroniza las cantidades de los lotes con los productos individuales disponibles
func (c *MedicalSupplyController) SyncBatchAmounts(ctx *gin.Context) {
	err := c.medicalSupplyService.SyncBatchAmounts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error sincronizando cantidades de lotes: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Cantidades de lotes sincronizadas exitosamente",
	})
}
