package controllers

import (
	"meditrack/models"
	"meditrack/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PavilionController struct {
	pavilionService services.PavilionService
}

func NewPavilionController(pavilionService services.PavilionService) *PavilionController {
	return &PavilionController{pavilionService: pavilionService}
}

func (c *PavilionController) CreatePavilion(ctx *gin.Context) {
	var pavilion models.Pavilion

	if err := ctx.ShouldBindJSON(&pavilion); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}

	if pavilion.Name == "" {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "El nombre es requerido"})
		return
	}

	if err := c.pavilionService.CreatePavilion(&pavilion); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al crear pavilion: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, Response{Success: true, Message: "Pavilion creado", Data: pavilion})
}

func (c *PavilionController) GetPavilionByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	pavilion, err := c.pavilionService.GetPavilionByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{Success: false, Error: "Pavilion no encontrado: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: pavilion})
}

func (c *PavilionController) GetAllPavilions(ctx *gin.Context) {
	pavilions, err := c.pavilionService.GetAllPavilions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al obtener pavilions: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Data: pavilions})
}

func (c *PavilionController) UpdatePavilion(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	var pavilion models.Pavilion
	if err := ctx.ShouldBindJSON(&pavilion); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "Datos inválidos: " + err.Error()})
		return
	}
	updatedPavilion, err := c.pavilionService.UpdatePavilion(id, &pavilion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al actualizar pavilion: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Pavilion actualizado", Data: updatedPavilion})
}

func (c *PavilionController) DeletePavilion(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Success: false, Error: "ID inválido: " + err.Error()})
		return
	}
	if err := c.pavilionService.DeletePavilion(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{Success: false, Error: "Error al eliminar pavilion: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Response{Success: true, Message: "Pavilion eliminado"})
}
