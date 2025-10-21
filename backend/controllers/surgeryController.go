package controllers

import (
	"meditrack/models"
	"meditrack/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SurgeryController struct {
	surgeryService *services.SurgeryService
}

func NewSurgeryController(surgeryService *services.SurgeryService) *SurgeryController {
	return &SurgeryController{
		surgeryService: surgeryService,
	}
}

// CreateSurgery crea un nuevo tipo de cirugía
func (c *SurgeryController) CreateSurgery(ctx *gin.Context) {
	var surgery models.Surgery

	if err := ctx.ShouldBindJSON(&surgery); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	if err := c.surgeryService.CreateSurgery(&surgery); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al crear tipo de cirugía",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Tipo de cirugía creado exitosamente",
		Data:    surgery,
	})
}

// GetSurgeryByID obtiene un tipo de cirugía por ID
func (c *SurgeryController) GetSurgeryByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID inválido",
			Error:   err.Error(),
		})
		return
	}

	surgery, err := c.surgeryService.GetSurgeryByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Tipo de cirugía no encontrado",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Tipo de cirugía encontrado",
		Data:    surgery,
	})
}

// GetAllSurgeries obtiene todos los tipos de cirugía
func (c *SurgeryController) GetAllSurgeries(ctx *gin.Context) {
	surgeries, err := c.surgeryService.GetAllSurgeries()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al obtener tipos de cirugía",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Tipos de cirugía obtenidos",
		Data: gin.H{
			"surgeries": surgeries,
			"count":     len(surgeries),
		},
	})
}

// GetSurgeriesPaginated obtiene tipos de cirugía con paginación
func (c *SurgeryController) GetSurgeriesPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	var search *string
	if s := ctx.Query("search"); s != "" {
		search = &s
	}

	surgeries, total, err := c.surgeryService.GetSurgeriesPaginated(page, pageSize, search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al obtener tipos de cirugía",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Tipos de cirugía obtenidos",
		Data: gin.H{
			"surgeries":   surgeries,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": (int(total) + pageSize - 1) / pageSize,
		},
	})
}

// UpdateSurgery actualiza un tipo de cirugía
func (c *SurgeryController) UpdateSurgery(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID inválido",
			Error:   err.Error(),
		})
		return
	}

	var surgery models.Surgery
	if err := ctx.ShouldBindJSON(&surgery); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	updatedSurgery, err := c.surgeryService.UpdateSurgery(id, &surgery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al actualizar tipo de cirugía",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Tipo de cirugía actualizado exitosamente",
		Data:    updatedSurgery,
	})
}

// DeleteSurgery elimina un tipo de cirugía
func (c *SurgeryController) DeleteSurgery(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID inválido",
			Error:   err.Error(),
		})
		return
	}

	if err := c.surgeryService.DeleteSurgery(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al eliminar tipo de cirugía",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Tipo de cirugía eliminado exitosamente",
	})
}

// SearchSurgeries busca tipos de cirugía por nombre
func (c *SurgeryController) SearchSurgeries(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Parámetro 'name' es requerido",
		})
		return
	}

	surgeries, err := c.surgeryService.SearchSurgeriesByName(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al buscar tipos de cirugía",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Búsqueda completada",
		Data: gin.H{
			"surgeries": surgeries,
			"count":     len(surgeries),
		},
	})
}
