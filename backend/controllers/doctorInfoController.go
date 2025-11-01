package controllers

import (
	"meditrack/models"
	"meditrack/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DoctorInfoController struct {
	doctorInfoService *services.DoctorInfoService
}

func NewDoctorInfoController(doctorInfoService *services.DoctorInfoService) *DoctorInfoController {
	return &DoctorInfoController{
		doctorInfoService: doctorInfoService,
	}
}

// CreateDoctorInfo crea información extendida de un doctor
func (c *DoctorInfoController) CreateDoctorInfo(ctx *gin.Context) {
	var doctorInfo models.DoctorInfo

	if err := ctx.ShouldBindJSON(&doctorInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	if err := c.doctorInfoService.CreateDoctorInfo(&doctorInfo); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al crear información del doctor",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Información del doctor creada exitosamente",
		Data:    doctorInfo,
	})
}

// GetDoctorInfoByRUT obtiene información de un doctor por RUT
func (c *DoctorInfoController) GetDoctorInfoByRUT(ctx *gin.Context) {
	rut := ctx.Param("rut")
	if rut == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "RUT inválido",
		})
		return
	}

	doctorInfo, err := c.doctorInfoService.GetDoctorInfoByRUT(rut)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Información del doctor no encontrada",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Información del doctor encontrada",
		Data:    doctorInfo,
	})
}

// GetAllDoctorInfo obtiene información de todos los doctores
func (c *DoctorInfoController) GetAllDoctorInfo(ctx *gin.Context) {
	doctorsInfo, err := c.doctorInfoService.GetAllDoctorInfo()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al obtener información de doctores",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Información de doctores obtenida",
		Data: gin.H{
			"doctors": doctorsInfo,
			"count":   len(doctorsInfo),
		},
	})
}

// GetDoctorsBySpecialtyID obtiene todos los doctores de una especialidad
func (c *DoctorInfoController) GetDoctorsBySpecialtyID(ctx *gin.Context) {
	specialtyID, err := strconv.Atoi(ctx.Param("specialty_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "ID de especialidad inválido",
			Error:   err.Error(),
		})
		return
	}

	doctorsInfo, err := c.doctorInfoService.GetDoctorsBySpecialtyID(specialtyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al obtener doctores de la especialidad",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Doctores obtenidos",
		Data: gin.H{
			"doctors": doctorsInfo,
			"count":   len(doctorsInfo),
		},
	})
}

// GetDoctorsPaginated obtiene doctores con paginación
func (c *DoctorInfoController) GetDoctorsPaginated(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	var search *string
	if s := ctx.Query("search"); s != "" {
		search = &s
	}

	var specialtyID *int
	if sid := ctx.Query("specialty_id"); sid != "" {
		id, err := strconv.Atoi(sid)
		if err == nil {
			specialtyID = &id
		}
	}

	doctorsInfo, total, err := c.doctorInfoService.GetDoctorsPaginated(page, pageSize, search, specialtyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al obtener doctores",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Doctores obtenidos",
		Data: gin.H{
			"doctors":      doctorsInfo,
			"total":        total,
			"page":         page,
			"page_size":    pageSize,
			"total_pages":  (int(total) + pageSize - 1) / pageSize,
		},
	})
}

// UpdateDoctorInfo actualiza información de un doctor
func (c *DoctorInfoController) UpdateDoctorInfo(ctx *gin.Context) {
	rut := ctx.Param("rut")
	if rut == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "RUT inválido",
		})
		return
	}

	var doctorInfo models.DoctorInfo
	if err := ctx.ShouldBindJSON(&doctorInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	updatedDoctorInfo, err := c.doctorInfoService.UpdateDoctorInfo(rut, &doctorInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al actualizar información del doctor",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Información del doctor actualizada exitosamente",
		Data:    updatedDoctorInfo,
	})
}

// DeleteDoctorInfo elimina información de un doctor
func (c *DoctorInfoController) DeleteDoctorInfo(ctx *gin.Context) {
	rut := ctx.Param("rut")
	if rut == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "RUT inválido",
		})
		return
	}

	if err := c.doctorInfoService.DeleteDoctorInfo(rut); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al eliminar información del doctor",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Información del doctor eliminada exitosamente",
	})
}

// GetDoctorsBySpecialtyCode obtiene todos los doctores de una especialidad por código
func (c *DoctorInfoController) GetDoctorsBySpecialtyCode(ctx *gin.Context) {
	specialtyCode := ctx.Param("code")
	if specialtyCode == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Código de especialidad inválido",
		})
		return
	}

	doctorsInfo, err := c.doctorInfoService.GetDoctorsBySpecialtyCode(specialtyCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al obtener doctores de la especialidad",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Doctores obtenidos",
		Data: gin.H{
			"doctors": doctorsInfo,
			"count":   len(doctorsInfo),
		},
	})
}

