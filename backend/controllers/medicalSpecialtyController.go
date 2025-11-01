package controllers

import (
	"meditrack/models"
	"meditrack/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MedicalSpecialtyController struct {
	specialtyService *services.MedicalSpecialtyService
}

func NewMedicalSpecialtyController(specialtyService *services.MedicalSpecialtyService) *MedicalSpecialtyController {
	return &MedicalSpecialtyController{
		specialtyService: specialtyService,
	}
}

// CreateMedicalSpecialty crea una nueva especialidad médica
func (c *MedicalSpecialtyController) CreateMedicalSpecialty(ctx *gin.Context) {
	var specialty models.MedicalSpecialty

	if err := ctx.ShouldBindJSON(&specialty); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	if err := c.specialtyService.CreateMedicalSpecialty(&specialty); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al crear especialidad médica",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Especialidad médica creada exitosamente",
		Data:    specialty,
	})
}

// GetMedicalSpecialtyByID obtiene una especialidad médica por ID
func (c *MedicalSpecialtyController) GetMedicalSpecialtyByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID inválido",
			Error:   err.Error(),
		})
		return
	}

	specialty, err := c.specialtyService.GetMedicalSpecialtyByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Especialidad médica no encontrada",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Especialidad médica encontrada",
		Data:    specialty,
	})
}

// GetAllMedicalSpecialties obtiene todas las especialidades médicas
func (c *MedicalSpecialtyController) GetAllMedicalSpecialties(ctx *gin.Context) {
	specialties, err := c.specialtyService.GetAllMedicalSpecialties()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al obtener especialidades médicas",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Especialidades médicas obtenidas",
		Data: gin.H{
			"specialties": specialties,
			"count":       len(specialties),
		},
	})
}

// GetMedicalSpecialtiesPaginated obtiene especialidades médicas con paginación
func (c *MedicalSpecialtyController) GetMedicalSpecialtiesPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	var search *string
	if s := ctx.Query("search"); s != "" {
		search = &s
	}

	specialties, total, err := c.specialtyService.GetMedicalSpecialtiesPaginated(page, pageSize, search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al obtener especialidades médicas",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Especialidades médicas obtenidas",
		Data: gin.H{
			"specialties":   specialties,
			"total":          total,
			"page":           page,
			"page_size":      pageSize,
			"total_pages":    (int(total) + pageSize - 1) / pageSize,
		},
	})
}

// UpdateMedicalSpecialty actualiza una especialidad médica
func (c *MedicalSpecialtyController) UpdateMedicalSpecialty(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID inválido",
			Error:   err.Error(),
		})
		return
	}

	var specialty models.MedicalSpecialty
	if err := ctx.ShouldBindJSON(&specialty); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	updatedSpecialty, err := c.specialtyService.UpdateMedicalSpecialty(id, &specialty)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al actualizar especialidad médica",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Especialidad médica actualizada exitosamente",
		Data:    updatedSpecialty,
	})
}

// DeleteMedicalSpecialty elimina una especialidad médica
func (c *MedicalSpecialtyController) DeleteMedicalSpecialty(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID inválido",
			Error:   err.Error(),
		})
		return
	}

	if err := c.specialtyService.DeleteMedicalSpecialty(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al eliminar especialidad médica",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Especialidad médica eliminada exitosamente",
	})
}

// SearchMedicalSpecialties busca especialidades médicas por nombre
func (c *MedicalSpecialtyController) SearchMedicalSpecialties(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Parámetro 'name' es requerido",
		})
		return
	}

	specialties, err := c.specialtyService.SearchMedicalSpecialtiesByName(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al buscar especialidades médicas",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Búsqueda completada",
		Data: gin.H{
			"specialties": specialties,
			"count":       len(specialties),
		},
	})
}

