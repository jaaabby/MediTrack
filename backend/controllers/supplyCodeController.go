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
	var supplyCode models.SupplyCode
	if err := ctx.ShouldBindJSON(&supplyCode); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
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
	}
	var newSupplyCode models.SupplyCode
	if err := ctx.ShouldBindJSON(&newSupplyCode); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}
	supplyCode.Name = newSupplyCode.Name
	supplyCode.CodeSupplier = newSupplyCode.CodeSupplier
	supplyCode.BatchID = newSupplyCode.BatchID
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
