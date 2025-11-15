package controllers

import (
	"fmt"
	"meditrack/models"
	"meditrack/pkg/response"
	"meditrack/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SupplyCodeController struct {
	supplyCodeService services.SupplyCodeService
}

func NewSupplyCodeController(supplyCodeService services.SupplyCodeService) *SupplyCodeController {
	return &SupplyCodeController{supplyCodeService: supplyCodeService}
}

// parseIDFromParam parsea el ID del parámetro de la URL
func (c *SupplyCodeController) parseIDFromParam(ctx *gin.Context) (int, error) {
	id := ctx.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("ID inválido: debe ser un número entero")
	}
	return intID, nil
}

func (c *SupplyCodeController) CreateSupplyCode(ctx *gin.Context) {
	var supplyCodeRequest struct {
		Code         int    `json:"code" binding:"required"`
		Name         string `json:"name" binding:"required"`
		CodeSupplier int    `json:"code_supplier" binding:"required"`
		CriticalStock int    `json:"critical_stock" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&supplyCodeRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}

	// Crear modelo - Code es la primary key
	supplyCode := models.SupplyCode{
		Code:         supplyCodeRequest.Code,
		Name:         supplyCodeRequest.Name,
		CodeSupplier: supplyCodeRequest.CodeSupplier,
		CriticalStock: supplyCodeRequest.CriticalStock,
	}

	if err := c.supplyCodeService.CreateSupplyCode(&supplyCode); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: "Error al crear supply code: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, response.Response{Success: true, Message: "Supply code creado", Data: supplyCode})
}

func (c *SupplyCodeController) GetSupplyCodeByID(ctx *gin.Context) {
	intID, err := c.parseIDFromParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Success: false, Error: err.Error()})
		return
	}
	supplyCode, err := c.supplyCodeService.GetSupplyCodeByID(intID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{Success: false, Error: "Supply code no encontrado: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{Success: true, Data: supplyCode})
}

func (c *SupplyCodeController) GetAllSupplyCodes(ctx *gin.Context) {
	supplyCodes, err := c.supplyCodeService.GetAllSupplyCodes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: "Error al obtener supply codes: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{Success: true, Data: supplyCodes})
}

func (c *SupplyCodeController) UpdateSupplyCode(ctx *gin.Context) {
	intID, err := c.parseIDFromParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Success: false, Error: err.Error()})
		return
	}
	
	// Verificar que el supply code existe
	_, err = c.supplyCodeService.GetSupplyCodeByID(intID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{Success: false, Error: "Supply code no encontrado: " + err.Error()})
		return
	}
	
	var updateRequest struct {
		Name         string `json:"name" binding:"required"`
		CodeSupplier int    `json:"code_supplier" binding:"required"`
		CriticalStock int    `json:"critical_stock" binding:"required,min=1"`
	}
	
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}
	
	// Crear modelo con los datos actualizados
	updatedSupplyCode := models.SupplyCode{
		Code:         intID,
		Name:         updateRequest.Name,
		CodeSupplier: updateRequest.CodeSupplier,
		CriticalStock: updateRequest.CriticalStock,
	}
	
	result, err := c.supplyCodeService.UpdateSupplyCode(intID, &updatedSupplyCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: "Error al actualizar supply code: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{Success: true, Message: "Supply code actualizado", Data: result})
}

func (c *SupplyCodeController) DeleteSupplyCode(ctx *gin.Context) {
	intID, err := c.parseIDFromParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Success: false, Error: err.Error()})
		return
	}
	if err := c.supplyCodeService.DeleteSupplyCode(intID); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: "Error al eliminar supply code: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{Success: true, Message: "Supply code eliminado"})
}
