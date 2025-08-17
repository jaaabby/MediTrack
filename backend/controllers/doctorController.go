package controllers

import (
	"net/http"

	"meditrack/models"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// DoctorController maneja las peticiones HTTP relacionadas con doctores
type DoctorController struct {
	doctorService services.UserService
}

// NewDoctorController crea una nueva instancia de DoctorController
func NewDoctorController(doctorService services.UserService) *DoctorController {
	return &DoctorController{
		doctorService: doctorService,
	}
}

// CreateDoctor crea un nuevo doctor
func (c *DoctorController) CreateDoctor(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de usuario inválidos: " + err.Error(),
		})
		return
	}
	user.Role = "doctor"
	if err := c.doctorService.CreateUser(ctx, &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear doctor: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Doctor creado exitosamente",
		Data:    user,
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

	user, err := c.doctorService.GetUserByID(ctx, id)
	if err != nil || user.Role != "doctor" {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Doctor no encontrado: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

// GetAllDoctors obtiene todos los doctores
func (c *DoctorController) GetAllDoctors(ctx *gin.Context) {
	users, err := c.doctorService.GetAllUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener doctores: " + err.Error(),
		})
		return
	}
	var doctors []*models.User
	for _, u := range users {
		if u.Role == "doctor" {
			doctors = append(doctors, u)
		}
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

	users, err := c.doctorService.GetUsersBySpecialty(ctx, specialty)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener doctores por especialidad: " + err.Error(),
		})
		return
	}
	var doctors []*models.User
	for _, u := range users {
		if u.Role == "doctor" {
			doctors = append(doctors, u)
		}
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

	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de usuario inválidos: " + err.Error(),
		})
		return
	}
	user.Role = "doctor"
	if err := c.doctorService.UpdateUser(ctx, &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar doctor: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Doctor actualizado exitosamente",
		Data:    user,
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

	if err := c.doctorService.DeleteUser(ctx, id); err != nil {
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
