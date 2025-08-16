package controllers

import (
	"net/http"

	"meditrack/models"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// LocationController maneja las peticiones HTTP relacionadas con ubicaciones
type LocationController struct {
	locationService services.LocationService
}

// NewLocationController crea una nueva instancia de LocationController
func NewLocationController(locationService services.LocationService) *LocationController {
	return &LocationController{
		locationService: locationService,
	}
}

// CreateLocation crea una nueva ubicación
func (c *LocationController) CreateLocation(ctx *gin.Context) {
	var location models.Location
	if err := ctx.ShouldBindJSON(&location); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de ubicación inválidos: " + err.Error(),
		})
		return
	}

	if err := c.locationService.CreateLocation(ctx, &location); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear ubicación: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Ubicación creada exitosamente",
		Data:    location,
	})
}

// GetLocationByID obtiene una ubicación por ID
func (c *LocationController) GetLocationByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de ubicación requerido",
		})
		return
	}

	location, err := c.locationService.GetLocationByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Ubicación no encontrada: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    location,
	})
}

// GetAllLocations obtiene todas las ubicaciones
func (c *LocationController) GetAllLocations(ctx *gin.Context) {
	locations, err := c.locationService.GetAllLocations(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener ubicaciones: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    locations,
	})
}

// UpdateLocation actualiza una ubicación existente
func (c *LocationController) UpdateLocation(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de ubicación requerido",
		})
		return
	}

	var location models.Location
	if err := ctx.ShouldBindJSON(&location); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de ubicación inválidos: " + err.Error(),
		})
		return
	}

	if err := c.locationService.UpdateLocation(ctx, &location); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar ubicación: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Ubicación actualizada exitosamente",
		Data:    location,
	})
}

// DeleteLocation elimina una ubicación
func (c *LocationController) DeleteLocation(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de ubicación requerido",
		})
		return
	}

	if err := c.locationService.DeleteLocation(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al eliminar ubicación: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Ubicación eliminada exitosamente",
	})
}
