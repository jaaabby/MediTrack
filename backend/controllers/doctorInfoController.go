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

// CreateDoctor crea un nuevo usuario doctor
func (c *DoctorInfoController) CreateDoctor(ctx *gin.Context) {
	var doctor models.User

	if err := ctx.ShouldBindJSON(&doctor); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	// Validar que el rol sea doctor
	if doctor.Role != models.RoleDoctor {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "El rol debe ser 'doctor'",
		})
		return
	}

	if err := c.doctorInfoService.CreateDoctor(&doctor); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al crear doctor",
			Error:   err.Error(),
		})
		return
	}

	// Recargar con relaciones
	createdDoctor, err := c.doctorInfoService.GetDoctorByRUT(doctor.RUT)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Doctor creado pero error al cargar datos completos",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Doctor creado exitosamente",
		Data:    createdDoctor,
	})
}

// GetDoctorByRUT obtiene un doctor por RUT
func (c *DoctorInfoController) GetDoctorByRUT(ctx *gin.Context) {
	rut := ctx.Param("rut")
	if rut == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "RUT inválido",
		})
		return
	}

	doctor, err := c.doctorInfoService.GetDoctorByRUT(rut)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Message: "Doctor no encontrado",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Doctor encontrado",
		Data:    doctor,
	})
}

// GetAllDoctors obtiene todos los doctores
func (c *DoctorInfoController) GetAllDoctors(ctx *gin.Context) {
	doctors, err := c.doctorInfoService.GetAllDoctors()
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
			"doctors": doctors,
			"count":   len(doctors),
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

	doctors, err := c.doctorInfoService.GetDoctorsBySpecialtyID(specialtyID)
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
			"doctors": doctors,
			"count":   len(doctors),
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

	doctors, total, err := c.doctorInfoService.GetDoctorsPaginated(page, pageSize, search, specialtyID)
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
			"doctors":     doctors,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": (int(total) + pageSize - 1) / pageSize,
		},
	})
}

// UpdateDoctor actualiza información de un doctor
func (c *DoctorInfoController) UpdateDoctor(ctx *gin.Context) {
	rut := ctx.Param("rut")
	if rut == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "RUT inválido",
		})
		return
	}

	var doctor models.User
	if err := ctx.ShouldBindJSON(&doctor); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Datos inválidos",
			Error:   err.Error(),
		})
		return
	}

	updatedDoctor, err := c.doctorInfoService.UpdateDoctor(rut, &doctor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al actualizar doctor",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Doctor actualizado exitosamente",
		Data:    updatedDoctor,
	})
}

// DeleteDoctor elimina un doctor de la base de datos
func (c *DoctorInfoController) DeleteDoctor(ctx *gin.Context) {
	rut := ctx.Param("rut")
	if rut == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "RUT inválido",
		})
		return
	}

	if err := c.doctorInfoService.DeleteDoctor(rut); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error al eliminar doctor",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Doctor eliminado exitosamente",
	})
}

// Métodos de compatibilidad con la API anterior (deprecated)
// Estos métodos redirigen a los nuevos métodos

// CreateDoctorInfo crea información extendida de un doctor (compatibilidad)
func (c *DoctorInfoController) CreateDoctorInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: "Este endpoint está deprecado, use POST /doctors en su lugar",
	})
}

// GetDoctorInfoByRUT obtiene información de un doctor por RUT (compatibilidad)
func (c *DoctorInfoController) GetDoctorInfoByRUT(ctx *gin.Context) {
	// Redirigir al nuevo método
	c.GetDoctorByRUT(ctx)
}

// GetAllDoctorInfo obtiene información de todos los doctores (compatibilidad)
func (c *DoctorInfoController) GetAllDoctorInfo(ctx *gin.Context) {
	// Redirigir al nuevo método
	c.GetAllDoctors(ctx)
}

// GetDoctorsBySpecialtyID ya está implementado arriba

// GetDoctorsPaginated ya está implementado arriba

// UpdateDoctorInfo actualiza información de un doctor (compatibilidad)
func (c *DoctorInfoController) UpdateDoctorInfo(ctx *gin.Context) {
	// Redirigir al nuevo método
	c.UpdateDoctor(ctx)
}

// DeleteDoctorInfo elimina información de un doctor (compatibilidad)
func (c *DoctorInfoController) DeleteDoctorInfo(ctx *gin.Context) {
	// Redirigir al nuevo método
	c.DeleteDoctor(ctx)
}
