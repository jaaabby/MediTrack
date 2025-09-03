package controllers

import (
	"fmt"
	"net/http"
	"time"

	"meditrack/config"
	"meditrack/models"
	"meditrack/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthController maneja las peticiones HTTP relacionadas con autenticación
type AuthController struct {
	userService services.UserService
	secretKey   string
}

// NewAuthController crea una nueva instancia de AuthController
func NewAuthController(userService services.UserService, secretKey string) *AuthController {
	return &AuthController{
		userService: userService,
		secretKey:   secretKey,
	}
}

// LoginRequest representa la solicitud de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest representa la solicitud de registro
type RegisterRequest struct {
	RUT             string `json:"rut" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	Role            string `json:"role" binding:"required,oneof=admin pabellón encargado de bodega"`
	MedicalCenterID int    `json:"medical_center_id" binding:"required"`
}

// LoginResponse representa la respuesta del login
type LoginResponse struct {
	Token     string              `json:"token"`
	User      models.UserResponse `json:"user"`
	ExpiresIn int64               `json:"expires_in"`
	TokenType string              `json:"token_type"`
}

// Login autentica un usuario y retorna un token JWT
func (c *AuthController) Login(ctx *gin.Context) {
	var loginReq LoginRequest
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de login inválidos: " + err.Error(),
		})
		return
	}

	// Buscar usuario por email
	user, err := c.userService.GetUserByEmail(loginReq.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Credenciales inválidas",
		})
		return
	}

	// Verificar si el usuario está activo
	if !user.IsActive {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Usuario inactivo",
		})
		return
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		fmt.Println("user.Password", user.Password)
		fmt.Println("loginReq.Password", loginReq.Password)
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Credenciales inválidas",
		})
		return
	}

	// Generar token JWT (24 horas de duración)
	duration := 24 * time.Hour
	token, err := config.GenerateToken(user.RUT, user.Email, user.Role, c.secretKey, duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al generar token: " + err.Error(),
		})
		return
	}

	// Crear respuesta
	loginResp := LoginResponse{
		Token:     token,
		User:      user.ToResponse(),
		ExpiresIn: int64(duration.Seconds()),
		TokenType: "Bearer",
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Login exitoso",
		Data:    loginResp,
	})
}

// GetProfile obtiene el perfil del usuario autenticado
func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Usuario no autenticado",
		})
		return
	}

	user, err := c.userService.GetUserByID(userID.(string))
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

// ChangePassword cambia la contraseña del usuario autenticado
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var changePassReq struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := ctx.ShouldBindJSON(&changePassReq); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	userID, _ := ctx.Get("user_id")
	user, err := c.userService.GetUserByID(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Usuario no encontrado",
		})
		return
	}

	// Verificar contraseña actual
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePassReq.CurrentPassword)); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Contraseña actual incorrecta",
		})
		return
	}

	// Hashear nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePassReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al procesar nueva contraseña",
		})
		return
	}

	// Actualizar contraseña
	user.Password = string(hashedPassword)
	if _, err := c.userService.UpdateUser(user.RUT, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar contraseña: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Contraseña cambiada exitosamente",
	})
}

// Register registra un nuevo usuario en el sistema
func (c *AuthController) Register(ctx *gin.Context) {
	var registerReq RegisterRequest
	if err := ctx.ShouldBindJSON(&registerReq); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de registro inválidos: " + err.Error(),
		})
		return
	}

	// Verificar si el usuario ya existe por RUT
	existingUserByRUT, _ := c.userService.GetUserByID(registerReq.RUT)
	if existingUserByRUT != nil {
		ctx.JSON(http.StatusConflict, Response{
			Success: false,
			Error:   "Ya existe un usuario con este RUT",
		})
		return
	}

	// Verificar si el usuario ya existe por email
	existingUserByEmail, _ := c.userService.GetUserByEmail(registerReq.Email)
	if existingUserByEmail != nil {
		ctx.JSON(http.StatusConflict, Response{
			Success: false,
			Error:   "Ya existe un usuario con este email",
		})
		return
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al procesar contraseña: " + err.Error(),
		})
		return
	}

	// Crear nuevo usuario
	newUser := &models.User{
		RUT:             registerReq.RUT,
		Name:            registerReq.Name,
		Email:           registerReq.Email,
		Password:        string(hashedPassword),
		Role:            registerReq.Role,
		MedicalCenterID: registerReq.MedicalCenterID,
		IsActive:        true,
	}

	// Guardar usuario en la base de datos
	if err := c.userService.CreateUser(newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Usuario registrado exitosamente",
		Data:    newUser.ToResponse(),
	})
}
