package controllers

import (
	"net/http"

	"meditrack/models"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// OperatingRoomController maneja las peticiones HTTP relacionadas con salas de operación
type OperatingRoomController struct {
	operatingRoomService services.OperatingRoomService
}

// NewOperatingRoomController crea una nueva instancia de OperatingRoomController
func NewOperatingRoomController(operatingRoomService services.OperatingRoomService) *OperatingRoomController {
	return &OperatingRoomController{
		operatingRoomService: operatingRoomService,
	}
}

// CreateOperatingRoom crea una nueva sala de operación
func (c *OperatingRoomController) CreateOperatingRoom(ctx *gin.Context) {
	var room models.OperatingRoom
	if err := ctx.ShouldBindJSON(&room); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de sala de operación inválidos: " + err.Error(),
		})
		return
	}

	if err := c.operatingRoomService.CreateOperatingRoom(ctx, &room); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear sala de operación: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Sala de operación creada exitosamente",
		Data:    room,
	})
}

// GetOperatingRoomByID obtiene una sala de operación por ID
func (c *OperatingRoomController) GetOperatingRoomByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de sala de operación requerido",
		})
		return
	}

	room, err := c.operatingRoomService.GetOperatingRoomByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Sala de operación no encontrada: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    room,
	})
}

// GetAllOperatingRooms obtiene todas las salas de operación
func (c *OperatingRoomController) GetAllOperatingRooms(ctx *gin.Context) {
	rooms, err := c.operatingRoomService.GetAllOperatingRooms(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener salas de operación: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    rooms,
	})
}

// UpdateOperatingRoom actualiza una sala de operación existente
func (c *OperatingRoomController) UpdateOperatingRoom(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de sala de operación requerido",
		})
		return
	}

	var room models.OperatingRoom
	if err := ctx.ShouldBindJSON(&room); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de sala de operación inválidos: " + err.Error(),
		})
		return
	}

	if err := c.operatingRoomService.UpdateOperatingRoom(ctx, &room); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar sala de operación: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Sala de operación actualizada exitosamente",
		Data:    room,
	})
}

// DeleteOperatingRoom elimina una sala de operación
func (c *OperatingRoomController) DeleteOperatingRoom(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de sala de operación requerido",
		})
		return
	}

	if err := c.operatingRoomService.DeleteOperatingRoom(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al eliminar sala de operación: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Sala de operación eliminada exitosamente",
	})
}
