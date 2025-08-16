package controllers

import (
	"net/http"

	"meditrack/models"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// DoctorController maneja las peticiones HTTP relacionadas con doctores
type DoctorController struct {
	doctorService services.DoctorService
}

// NewDoctorController crea una nueva instancia de DoctorController
func NewDoctorController(doctorService services.DoctorService) *DoctorController {
	return &DoctorController{
		doctorService: doctorService,
	}
}

// CreateDoctor crea un nuevo doctor
func (c *DoctorController) CreateDoctor(ctx *gin.Context) {
	var doctor models.Doctor
	if err := ctx.ShouldBindJSON(&doctor); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de doctor inválidos: " + err.Error(),
		})
		return
	}

	if err := c.doctorService.CreateDoctor(ctx, &doctor); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear doctor: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Doctor creado exitosamente",
		Data:    doctor,
	})
}

// GetDoctorByID obtiene un doctor por ID
func (c *DoctorController) GetDoctorByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de doctor requerido",
		})
		return
	}

	doctor, err := c.doctorService.GetDoctorByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Doctor no encontrado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    doctor,
	})
}

// GetAllDoctors obtiene todos los doctores
func (c *DoctorController) GetAllDoctors(ctx *gin.Context) {
	doctors, err := c.doctorService.GetAllDoctors(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener doctores: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    doctors,
	})
}

// GetDoctorsBySpecialty obtiene doctores por especialidad
func (c *DoctorController) GetDoctorsBySpecialty(ctx *gin.Context) {
	specialty := ctx.Query("specialty")
	if specialty == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Especialidad requerida",
		})
		return
	}

	doctors, err := c.doctorService.GetDoctorsBySpecialty(ctx, specialty)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener doctores por especialidad: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    doctors,
	})
}

// UpdateDoctor actualiza un doctor existente
func (c *DoctorController) UpdateDoctor(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de doctor requerido",
		})
		return
	}

	var doctor models.Doctor
	if err := ctx.ShouldBindJSON(&doctor); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de doctor inválidos: " + err.Error(),
		})
		return
	}

	if err := c.doctorService.UpdateDoctor(ctx, &doctor); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar doctor: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Doctor actualizado exitosamente",
		Data:    doctor,
	})
}

// DeleteDoctor elimina un doctor
func (c *DoctorController) DeleteDoctor(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de doctor requerido",
		})
		return
	}

	if err := c.doctorService.DeleteDoctor(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al eliminar doctor: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Doctor eliminado exitosamente",
	})
}
