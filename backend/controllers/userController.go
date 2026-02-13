package controllers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"meditrack/mailer"
	"meditrack/models"
	"meditrack/pkg/response"
	"meditrack/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
		Role            string `json:"role" binding:"required"`
		MedicalCenterID int    `json:"medical_center_id" binding:"required"`
		PavilionID      *int   `json:"pavilion_id"`
		SpecialtyID     *int   `json:"specialty_id"`
	}

	// Log del body recibido
	bodyBytes, _ := ctx.GetRawData()
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	log.Printf("📝 CreateUser - Body recibido: %s", string(bodyBytes))

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		log.Printf("❌ Error en binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos de usuario inválidos: " + err.Error(),
		})
		return
	}

	log.Printf("✅ Datos parseados correctamente - RUT: %s, Role: %s, PavilionID: %v, SpecialtyID: %v",
		userRequest.RUT, userRequest.Role, userRequest.PavilionID, userRequest.SpecialtyID)

	// Validar rol
	tempUser := models.User{Role: userRequest.Role}
	if !tempUser.IsValidRole() {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Rol inválido. Roles permitidos: admin, pabellón, encargado de bodega, enfermera, doctor, pavedad",
		})
		return
	}

	// Crear modelo
	user := models.User{
		RUT:             userRequest.RUT,
		Name:            userRequest.Name,
		Email:           userRequest.Email,
		Role:            userRequest.Role,
		MedicalCenterID: userRequest.MedicalCenterID,
	}

	// Asignar campos opcionales
	if userRequest.PavilionID != nil {
		user.PavilionID = userRequest.PavilionID
	}
	if userRequest.SpecialtyID != nil {
		user.SpecialtyID = userRequest.SpecialtyID
	}

	// Crear usuario con contraseña temporal
	tempPassword, err := c.userService.CreateUserWithTemporaryPassword(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al crear usuario: " + err.Error(),
		})
		return
	}

	// Preparar datos para el template de email
	templateData := struct {
		Name     string
		Email    string
		RUT      string
		Role     string
		Password string
	}{
		Name:     user.Name,
		Email:    user.Email,
		RUT:      user.RUT,
		Role:     user.Role,
		Password: tempPassword,
	}

	// Enviar email con contraseña temporal
	templatePath := filepath.Join("mailer", "templates", "temporary_password.html")
	mailReq := mailer.NewRequest([]string{user.Email}, "¡Bienvenido a MediTrack! - Tus credenciales de acceso")

	if err := mailReq.SendMailSkipTLS(templatePath, templateData); err != nil {
		log.Printf("⚠️ Error enviando email a %s: %v", user.Email, err)
		// Continuamos aunque falle el email - el usuario fue creado
		ctx.JSON(http.StatusCreated, response.Response{
			Success: true,
			Message: "Usuario creado exitosamente, pero hubo un error al enviar el correo. Por favor, contacte al administrador para obtener su contraseña temporal.",
			Data:    user.ToResponse(),
		})
		return
	}

	log.Printf("✅ Usuario creado y correo enviado exitosamente a: %s", user.Email)

	ctx.JSON(http.StatusCreated, response.Response{
		Success: true,
		Message: "Usuario creado exitosamente. Se ha enviado un correo con la contraseña temporal.",
		Data:    user.ToResponse(),
	})
}

// GetUserByID obtiene un usuario por ID
func (c *UserController) GetUserByID(ctx *gin.Context) {
	rut := ctx.Param("id")

	user, err := c.userService.GetUserByRut(rut)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   "Usuario no encontrado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Data:    user.ToResponse(),
	})
}

// GetAllUsers obtiene todos los usuarios
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := c.userService.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
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

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Data:    userResponses,
	})
}

// SearchUsers busca usuarios por nombre, RUT o email (accesible para admin y encargado de bodega)
func (c *UserController) SearchUsers(ctx *gin.Context) {
	searchTerm := ctx.Query("q")
	if searchTerm == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Término de búsqueda requerido",
		})
		return
	}

	users, err := c.userService.SearchUsers(searchTerm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al buscar usuarios: " + err.Error(),
		})
		return
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Data:    userResponses,
	})
}

// GetUsersByRole obtiene usuarios por rol
func (c *UserController) GetUsersByRole(ctx *gin.Context) {
	role := ctx.Query("role")
	if role == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Rol requerido",
		})
		return
	}

	user := models.User{Role: role}
	if !user.IsValidRole() {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Rol inválido",
		})
		return
	}

	users, err := c.userService.GetUsersByRole(role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al obtener usuarios: " + err.Error(),
		})
		return
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Data:    userResponses,
	})
}

// UpdateUser actualiza un usuario existente
func (c *UserController) UpdateUser(ctx *gin.Context) {
	rut := ctx.Param("id")

	var userRequest struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		Role            string `json:"role"`
		MedicalCenterID int    `json:"medical_center_id"`
		PavilionID      *int   `json:"pavilion_id"`
		SpecialtyID     *int   `json:"specialty_id"`
		IsActive        *bool  `json:"is_active"`
	}

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos de usuario inválidos: " + err.Error(),
		})
		return
	}

	// Validar rol si se proporciona
	if userRequest.Role != "" {
		user := models.User{Role: userRequest.Role}
		if !user.IsValidRole() {
			ctx.JSON(http.StatusBadRequest, response.Response{
				Success: false,
				Error:   "Rol inválido. Roles permitidos: admin, pabellón, encargado de bodega",
			})
			return
		}
	}

	// Hashear la contraseña si se proporciona
	passwordToSave := userRequest.Password
	if userRequest.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.Response{
				Success: false,
				Error:   "Error al procesar contraseña: " + err.Error(),
			})
			return
		}
		passwordToSave = string(hashedPassword)
	}

	// Crear usuario para actualización
	user := models.User{
		RUT:             rut,
		Name:            userRequest.Name,
		Email:           userRequest.Email,
		Password:        passwordToSave,
		Role:            userRequest.Role,
		MedicalCenterID: userRequest.MedicalCenterID,
	}

	// Asignar campos opcionales
	if userRequest.PavilionID != nil {
		user.PavilionID = userRequest.PavilionID
	}
	if userRequest.SpecialtyID != nil {
		user.SpecialtyID = userRequest.SpecialtyID
	}

	if userRequest.IsActive != nil {
		user.IsActive = *userRequest.IsActive
	}

	updatedUser, err := c.userService.UpdateUser(rut, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al actualizar usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Usuario actualizado exitosamente",
		Data:    updatedUser.ToResponse(),
	})
}

// DeleteUser elimina un usuario
func (c *UserController) DeleteUser(ctx *gin.Context) {
	rut := ctx.Param("id")

	if err := c.userService.DeleteUser(rut); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al eliminar usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Usuario eliminado exitosamente",
	})
}

// DeactivateUser desactiva un usuario
func (c *UserController) DeactivateUser(ctx *gin.Context) {
	rut := ctx.Param("id")

	if err := c.userService.DeactivateUser(rut); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al desactivar usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Usuario desactivado exitosamente",
	})
}

// ActivateUser activa un usuario
func (c *UserController) ActivateUser(ctx *gin.Context) {
	rut := ctx.Param("id")

	if err := c.userService.ActivateUser(rut); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al activar usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Usuario activado exitosamente",
	})
}

func (c *UserController) GetUserProfileByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Email requerido",
		})
		return
	}

	user, err := c.userService.GetUserByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   "Usuario no encontrado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Data:    user.ToResponse(),
	})
}
