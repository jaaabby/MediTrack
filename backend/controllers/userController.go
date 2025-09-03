package controllers

import (
	"net/http"

	"meditrack/models"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// UserController maneja las peticiones HTTP relacionadas con usuarios
type UserController struct {
	userService services.UserService
}

// NewUserController crea una nueva instancia de UserController
func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser crea un nuevo usuario
func (c *UserController) CreateUser(ctx *gin.Context) {
	var userRequest struct {
		RUT             string `json:"rut" binding:"required"`
		Name            string `json:"name" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=6"`
		Role            string `json:"role" binding:"required"`
		MedicalCenterID int    `json:"medical_center_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de usuario inválidos: " + err.Error(),
		})
		return
	}

	// Validar rol
	tempUser := models.User{Role: userRequest.Role}
	if !tempUser.IsValidRole() {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Rol inválido. Roles permitidos: admin, pabellón, encargado de bodega",
		})
		return
	}

	// Crear modelo
	user := models.User{
		RUT:             userRequest.RUT,
		Name:            userRequest.Name,
		Email:           userRequest.Email,
		Password:        userRequest.Password,
		Role:            userRequest.Role,
		MedicalCenterID: userRequest.MedicalCenterID,
		IsActive:        true,
	}

	if err := c.userService.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Usuario creado exitosamente",
		Data:    user.ToResponse(),
	})
}

// GetUserByID obtiene un usuario por ID
func (c *UserController) GetUserByID(ctx *gin.Context) {
	rut := ctx.Param("id")
	if rut == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "RUT de usuario requerido",
		})
		return
	}

	user, err := c.userService.GetUserByID(rut)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Usuario no encontrado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user.ToResponse(),
	})
}

// GetAllUsers obtiene todos los usuarios
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.userService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener usuarios: " + err.Error(),
		})
		return
	}

	// Convertir a UserResponse
	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    userResponses,
	})
}

// GetUsersByRole obtiene usuarios por rol
func (c *UserController) GetUsersByRole(ctx *gin.Context) {
	role := ctx.Query("role")
	if role == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Rol requerido",
		})
		return
	}

	user := models.User{Role: role}
	if !user.IsValidRole() {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Rol inválido",
		})
		return
	}

	users, err := c.userService.GetUsersByRole(role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener usuarios: " + err.Error(),
		})
		return
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    userResponses,
	})
}

// UpdateUser actualiza un usuario existente
func (c *UserController) UpdateUser(ctx *gin.Context) {
	rut := ctx.Param("id")
	if rut == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "RUT de usuario requerido",
		})
		return
	}

	var userRequest struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		Role            string `json:"role"`
		MedicalCenterID int    `json:"medical_center_id"`
		IsActive        *bool  `json:"is_active"`
	}

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de usuario inválidos: " + err.Error(),
		})
		return
	}

	// Validar rol si se proporciona
	if userRequest.Role != "" {
		user := models.User{Role: userRequest.Role}
		if !user.IsValidRole() {
			ctx.JSON(http.StatusBadRequest, Response{
				Success: false,
				Error:   "Rol inválido. Roles permitidos: admin, pabellón, encargado de bodega",
			})
			return
		}
	}

	// Crear usuario para actualización
	user := models.User{
		RUT:             rut,
		Name:            userRequest.Name,
		Email:           userRequest.Email,
		Password:        userRequest.Password,
		Role:            userRequest.Role,
		MedicalCenterID: userRequest.MedicalCenterID,
	}

	if userRequest.IsActive != nil {
		user.IsActive = *userRequest.IsActive
	}

	updatedUser, err := c.userService.UpdateUser(rut, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Usuario actualizado exitosamente",
		Data:    updatedUser.ToResponse(),
	})
}

// DeleteUser elimina un usuario
func (c *UserController) DeleteUser(ctx *gin.Context) {
	rut := ctx.Param("id")
	if rut == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "RUT de usuario requerido",
		})
		return
	}

	if err := c.userService.DeleteUser(rut); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al eliminar usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Usuario eliminado exitosamente",
	})
}

// DeactivateUser desactiva un usuario
func (c *UserController) DeactivateUser(ctx *gin.Context) {
	rut := ctx.Param("id")
	if rut == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "RUT de usuario requerido",
		})
		return
	}

	if err := c.userService.DeactivateUser(rut); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al desactivar usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Usuario desactivado exitosamente",
	})
}

// ActivateUser activa un usuario
func (c *UserController) ActivateUser(ctx *gin.Context) {
	rut := ctx.Param("id")
	if rut == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "RUT de usuario requerido",
		})
		return
	}

	if err := c.userService.ActivateUser(rut); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al activar usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Usuario activado exitosamente",
	})
}

func (c *UserController) GetUserProfileByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Email requerido",
		})
		return
	}

	user, err := c.userService.GetUserProfileByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Usuario no encontrado: " + err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user.ToResponse(),
	})
}
