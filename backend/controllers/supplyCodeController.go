package controllers

import (
	"fmt"
	"meditrack/models"
	"meditrack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SupplyCodeController struct {
	supplyCodeService services.SupplyCodeService
}

func NewSupplyCodeController(supplyCodeService services.SupplyCodeService) *SupplyCodeController {
	return &SupplyCodeController{supplyCodeService: supplyCodeService}
}

func (c *SupplyCodeController) CreateSupplyCode(ctx *gin.Context) {
	var supplyCodeRequest struct {
		Code         int    `json:"code" binding:"required"`
		Name         string `json:"name" binding:"required"`
		CodeSupplier int    `json:"code_supplier" binding:"required"`
		CriticalStock int    `json:"critical_stock" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&supplyCodeRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
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
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al crear supply code: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, Response{Success: true, Message: "Supply code creado", Data: supplyCode})
}

func (c *SupplyCodeController) GetSupplyCodeByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	supplyCode, err := c.supplyCodeService.GetSupplyCodeByID(intID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{Success: false, Error: "Supply code no encontrado: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: supplyCode})
}

func (c *SupplyCodeController) GetAllSupplyCodes(ctx *gin.Context) {
	supplyCodes, err := c.supplyCodeService.GetAllSupplyCodes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener supply codes: " + err.Error()})
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: supplyCodes})
}

func (c *SupplyCodeController) UpdateSupplyCode(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	supplyCode, err := c.supplyCodeService.GetSupplyCodeByID(intID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{Success: false, Error: "Supply code no encontrado: " + err.Error()})
		return
	}
	
	var updateRequest struct {
		Name         string `json:"name" binding:"required"`
		CodeSupplier int    `json:"code_supplier" binding:"required"`
		CriticalStock int    `json:"critical_stock" binding:"required,min=1"`
	}
	
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}
	
	supplyCode.Name = updateRequest.Name
	supplyCode.CodeSupplier = updateRequest.CodeSupplier
	supplyCode.CriticalStock = updateRequest.CriticalStock
	if _, err := c.supplyCodeService.UpdateSupplyCode(intID, supplyCode); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al actualizar supply code: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Supply code actualizado", Data: supplyCode})
}

func (c *SupplyCodeController) DeleteSupplyCode(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	if err := c.supplyCodeService.DeleteSupplyCode(intID); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al eliminar supply code: " + err.Error()})
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Supply code eliminado"})

}
