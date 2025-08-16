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
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de usuario inválidos: " + err.Error(),
		})
		return
	}

	if err := c.userService.CreateUser(ctx, &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Usuario creado exitosamente",
		Data:    user,
	})
}

// GetUserByID obtiene un usuario por ID
func (c *UserController) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de usuario requerido",
		})
		return
	}

	user, err := c.userService.GetUserByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Usuario no encontrado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

// GetAllUsers obtiene todos los usuarios
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.userService.GetAllUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener usuarios: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
	})
}

// UpdateUser actualiza un usuario existente
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de usuario requerido",
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

	if err := c.userService.UpdateUser(ctx, &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Usuario actualizado exitosamente",
		Data:    user,
	})
}

// DeleteUser elimina un usuario
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de usuario requerido",
		})
		return
	}

	if err := c.userService.DeleteUser(ctx, id); err != nil {
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
