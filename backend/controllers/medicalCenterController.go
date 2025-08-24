package controllers

import (
	"fmt"
	"meditrack/models"
	"meditrack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MedicalCenterController struct {
	medicalCenterService services.MedicalCenterService
}

func NewMedicalCenterController(medicalCenterService services.MedicalCenterService) *MedicalCenterController {
	return &MedicalCenterController{medicalCenterService: medicalCenterService}
}

func (c *MedicalCenterController) CreateMedicalCenter(ctx *gin.Context) {
	var centerRequest struct {
		Name    string `json:"name" binding:"required"`
		Address string `json:"address"`
		Phone   string `json:"phone"`
		Email   string `json:"email"`
	}

	if err := ctx.ShouldBindJSON(&centerRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}

	// Crear modelo sin ID
	center := models.MedicalCenter{
		Name:    centerRequest.Name,
		Address: centerRequest.Address,
		Phone:   centerRequest.Phone,
		Email:   centerRequest.Email,
		// ID se auto-generará
	}

	if err := c.medicalCenterService.CreateMedicalCenter(&center); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al crear centro médico: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, Response{Success: true, Message: "Centro médico creado", Data: center})
}

func (c *MedicalCenterController) GetMedicalCenterByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	center, err := c.medicalCenterService.GetMedicalCenterByID(intID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{Success: false, Error: "Centro médico no encontrado: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: center})
}

func (c *MedicalCenterController) GetAllMedicalCenters(ctx *gin.Context) {
	centers, err := c.medicalCenterService.GetAllMedicalCenters()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener centros médicos: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: centers})
}

func (c *MedicalCenterController) UpdateMedicalCenter(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	var center models.MedicalCenter
	if err := ctx.ShouldBindJSON(&center); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}
	center.ID = intID
	if _, err := c.medicalCenterService.UpdateMedicalCenter(intID, &center); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al actualizar centro médico: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Centro médico actualizado", Data: center})
}

func (c *MedicalCenterController) DeleteMedicalCenter(ctx *gin.Context) {
	id := ctx.Param("id")
	var intID int
	if _, err := fmt.Sscanf(id, "%d", &intID); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	if err := c.medicalCenterService.DeleteMedicalCenter(intID); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al eliminar centro médico: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Centro médico eliminado"})
}
