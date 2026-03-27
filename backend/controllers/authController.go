package controllers

import (
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"math/big"
	"meditrack/config"
	"meditrack/mailer"
	"meditrack/models"
	"meditrack/pkg/response"
	"meditrack/services"
	"meditrack/sms"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

// getUserIDFromContext obtiene el user_id del contexto de Gin
func (c *AuthController) getUserIDFromContext(ctx *gin.Context) (string, error) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return "", fmt.Errorf("usuario no autenticado")
	}
	return userID.(string), nil
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
	Role            string `json:"role" binding:"required,oneof=admin pabellón 'encargado de bodega' enfermera doctor"`
	MedicalCenterID int    `json:"medical_center_id" binding:"required"`
}

// LoginResponse representa la respuesta del login completo (post-OTP)
type LoginResponse struct {
	Token     string              `json:"token"`
	User      models.UserResponse `json:"user"`
	ExpiresIn int64               `json:"expires_in"`
	TokenType string              `json:"token_type"`
}

// LoginPendingResponse representa la respuesta cuando se requiere verificación OTP
type LoginPendingResponse struct {
	OtpSessionID string `json:"otp_session_id"`
	PhoneMasked  string `json:"phone_masked"`
}

// generateOTPCode genera un código numérico de 6 dígitos criptográficamente seguro
func generateOTPCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// maskPhone enmascara el número de teléfono mostrando solo los últimos 4 dígitos
func maskPhone(phone string) string {
	if len(phone) <= 4 {
		return "****"
	}
	return "****" + phone[len(phone)-4:]
}

// Login verifica credenciales y, si el usuario tiene teléfono registrado,
// envía un OTP por SMS retornando un session_id para el segundo paso.
// Si el usuario no tiene teléfono, hace login directo (compatibilidad).
func (c *AuthController) Login(ctx *gin.Context) {
	var loginReq LoginRequest
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos de login inválidos: " + err.Error(),
		})
		return
	}

	// Mensaje genérico para no revelar si el correo existe o no (evita user enumeration)
	const errCredenciales = "Credenciales inválidas. Verifica tu correo y contraseña."

	// Buscar usuario por email
	user, err := c.userService.GetUserByEmail(loginReq.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   errCredenciales,
		})
		return
	}

	// Verificar si el usuario está activo
	if !user.IsActive {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   errCredenciales,
		})
		return
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   errCredenciales,
		})
		return
	}

	// Si el usuario tiene teléfono, iniciar flujo OTP
	if user.Phone != nil && *user.Phone != "" {
		code, err := generateOTPCode()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.Response{
				Success: false,
				Error:   "Error al generar código de verificación",
			})
			return
		}

		// Guardar sesión OTP en la base de datos (expira en 10 minutos)
		expiresAt := time.Now().Add(10 * time.Minute).Unix()
		otpSession := &models.OtpSession{
			UserRUT:   user.RUT,
			Code:      code,
			ExpiresAt: expiresAt,
		}
		if err := c.userService.CreateOtpSession(otpSession); err != nil {
			ctx.JSON(http.StatusInternalServerError, response.Response{
				Success: false,
				Error:   "Error al crear sesión OTP",
			})
			return
		}

		// Enviar SMS
		message := fmt.Sprintf("Tu código de verificación MediTrack es: %s. Expira en 10 minutos.", code)
		if err := sms.Send(*user.Phone, message); err != nil {
			fmt.Printf("Error enviando SMS OTP: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, response.Response{
				Success: false,
				Error:   "Error al enviar código de verificación por SMS",
			})
			return
		}

		ctx.JSON(http.StatusOK, response.Response{
			Success: true,
			Message: "Código de verificación enviado por SMS",
			Data: LoginPendingResponse{
				OtpSessionID: otpSession.ID,
				PhoneMasked:  maskPhone(*user.Phone),
			},
		})
		return
	}

	// Sin teléfono: login directo (compatibilidad con cuentas sin 2FA)
	duration := 24 * time.Hour
	token, err := config.GenerateToken(user.RUT, user.Email, user.Role, c.secretKey, duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al generar token: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Login exitoso",
		Data: LoginResponse{
			Token:     token,
			User:      user.ToResponse(),
			ExpiresIn: int64(duration.Seconds()),
			TokenType: "Bearer",
		},
	})
}

// otpCodeRegex valida que el código OTP sea exactamente 6 dígitos numéricos
var otpCodeRegex = regexp.MustCompile(`^\d{6}$`)

// VerifyOTPRequest representa la solicitud de verificación del código OTP
type VerifyOTPRequest struct {
	OtpSessionID string `json:"otp_session_id" binding:"required"`
	Code         string `json:"code" binding:"required"`
}

// VerifyOTP verifica el código OTP recibido por SMS y, si es válido, devuelve el JWT.
// Protecciones implementadas:
//   - Máximo 5 intentos fallidos antes de invalidar la sesión
//   - Comparación segura contra timing attacks (crypto/subtle)
//   - Validación de formato (exactamente 6 dígitos)
//   - Mensajes de error genéricos para evitar enumeración de información
func (c *AuthController) VerifyOTP(ctx *gin.Context) {
	var req VerifyOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos inválidos",
		})
		return
	}

	// Validar formato del código: exactamente 6 dígitos
	if !otpCodeRegex.MatchString(req.Code) {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "El código debe ser de 6 dígitos numéricos",
		})
		return
	}

	otpSession, err := c.userService.GetOtpSession(req.OtpSessionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Código de verificación inválido o expirado",
		})
		return
	}

	// Verificar que no esté usado o bloqueado por intentos
	if otpSession.Used {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Código de verificación inválido o expirado",
		})
		return
	}

	// Verificar expiración
	if time.Now().Unix() > otpSession.ExpiresAt {
		_ = c.userService.InvalidateOtpSession(req.OtpSessionID)
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "El código ha expirado. Inicia sesión nuevamente para recibir un nuevo código",
		})
		return
	}

	// Comparación segura contra timing attacks usando crypto/subtle
	codeMatch := subtle.ConstantTimeCompare([]byte(otpSession.Code), []byte(req.Code)) == 1

	if !codeMatch {
		// Incrementar contador de intentos fallidos
		attempts, err := c.userService.IncrementOtpAttempts(req.OtpSessionID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.Response{
				Success: false,
				Error:   "Error al procesar verificación",
			})
			return
		}

		remaining := models.MaxOTPAttempts - attempts
		if remaining <= 0 {
			// Bloquear la sesión definitivamente
			_ = c.userService.InvalidateOtpSession(req.OtpSessionID)
			ctx.JSON(http.StatusBadRequest, response.Response{
				Success: false,
				Error:   "Demasiados intentos fallidos. Inicia sesión nuevamente para recibir un nuevo código",
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   fmt.Sprintf("Código incorrecto. Te quedan %d intentos", remaining),
		})
		return
	}

	// Marcar sesión como usada (un solo uso)
	if err := c.userService.MarkOtpSessionUsed(req.OtpSessionID); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al procesar verificación",
		})
		return
	}

	// Obtener usuario y generar JWT
	user, err := c.userService.GetUserByRut(otpSession.UserRUT)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al procesar verificación",
		})
		return
	}

	duration := 24 * time.Hour
	token, err := config.GenerateToken(user.RUT, user.Email, user.Role, c.secretKey, duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al procesar verificación",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Verificación exitosa",
		Data: LoginResponse{
			Token:     token,
			User:      user.ToResponse(),
			ExpiresIn: int64(duration.Seconds()),
			TokenType: "Bearer",
		},
	})
}

// GetProfile obtiene el perfil del usuario autenticado
func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, err := c.userService.GetUserByRut(userID)
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

// ChangePassword cambia la contraseña del usuario autenticado
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var changePassReq struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := ctx.ShouldBindJSON(&changePassReq); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, err := c.userService.GetUserByRut(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   "Usuario no encontrado",
		})
		return
	}

	// Verificar contraseña actual
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePassReq.CurrentPassword)); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Contraseña actual incorrecta",
		})
		return
	}

	// Hashear nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePassReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al procesar nueva contraseña",
		})
		return
	}

	// Actualizar contraseña y desactivar flag de cambio obligatorio
	user.Password = string(hashedPassword)
	user.MustChangePassword = false
	if _, err := c.userService.UpdateUser(user.RUT, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al actualizar contraseña: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Contraseña cambiada exitosamente",
	})
}

// FirstTimePasswordChange cambia la contraseña en el primer inicio de sesión
// No requiere contraseña actual si el usuario tiene MustChangePassword=true
func (c *AuthController) FirstTimePasswordChange(ctx *gin.Context) {
	var changePassReq struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := ctx.ShouldBindJSON(&changePassReq); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, err := c.userService.GetUserByRut(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   "Usuario no encontrado",
		})
		return
	}

	// Si MustChangePassword es true, no verificar contraseña actual
	// Si es false, requerir contraseña actual
	if !user.MustChangePassword {
		if changePassReq.CurrentPassword == "" {
			ctx.JSON(http.StatusBadRequest, response.Response{
				Success: false,
				Error:   "Se requiere la contraseña actual",
			})
			return
		}

		// Verificar contraseña actual
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePassReq.CurrentPassword)); err != nil {
			ctx.JSON(http.StatusBadRequest, response.Response{
				Success: false,
				Error:   "Contraseña actual incorrecta",
			})
			return
		}
	} else {
		// Si está en primer inicio, verificar contraseña temporal
		if changePassReq.CurrentPassword == "" {
			ctx.JSON(http.StatusBadRequest, response.Response{
				Success: false,
				Error:   "Se requiere la contraseña temporal",
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePassReq.CurrentPassword)); err != nil {
			ctx.JSON(http.StatusBadRequest, response.Response{
				Success: false,
				Error:   "Contraseña temporal incorrecta",
			})
			return
		}
	}

	// Hashear nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePassReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al procesar nueva contraseña",
		})
		return
	}

	// Actualizar contraseña y desactivar flag de cambio obligatorio
	user.Password = string(hashedPassword)
	user.MustChangePassword = false
	if _, err := c.userService.UpdateUser(user.RUT, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al actualizar contraseña: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Contraseña actualizada exitosamente",
	})
}

// Register registra un nuevo usuario en el sistema
func (c *AuthController) Register(ctx *gin.Context) {
	var registerReq RegisterRequest
	if err := ctx.ShouldBindJSON(&registerReq); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos de registro inválidos: " + err.Error(),
		})
		return
	}

	// Verificar si el usuario ya existe por RUT
	existingUserByRUT, _ := c.userService.GetUserByRut(registerReq.RUT)
	if existingUserByRUT != nil {
		ctx.JSON(http.StatusConflict, response.Response{
			Success: false,
			Error:   "Ya existe un usuario con este RUT",
		})
		return
	}

	// Verificar si el usuario ya existe por email
	existingUserByEmail, _ := c.userService.GetUserByEmail(registerReq.Email)
	if existingUserByEmail != nil {
		ctx.JSON(http.StatusConflict, response.Response{
			Success: false,
			Error:   "Ya existe un usuario con este email",
		})
		return
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
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
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al crear usuario: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response.Response{
		Success: true,
		Message: "Usuario registrado exitosamente",
		Data:    newUser.ToResponse(),
	})
}

// RequestPasswordResetRequest representa la solicitud para recuperar contraseña
type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// RequestPasswordReset genera un token de recuperación y envía email
func (c *AuthController) RequestPasswordReset(ctx *gin.Context) {
	var req RequestPasswordResetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Email inválido: " + err.Error(),
		})
		return
	}

	// Obtener usuario primero para verificar su estado
	user, err := c.userService.GetUserByEmail(req.Email)
	if err != nil {
		// Verificar si el error es porque el usuario no existe
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, response.Response{
				Success: false,
				Error:   "El correo electrónico ingresado no está registrado en el sistema",
			})
			return
		}
		// Otro tipo de error
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al procesar la solicitud",
		})
		return
	}

	// Verificar si el usuario debe cambiar su contraseña (primer inicio)
	if user.MustChangePassword {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "No puedes recuperar tu contraseña porque aún no has ingresado al sistema por primera vez. Usa la contraseña temporal enviada a tu correo.",
		})
		return
	}

	// Generar token de reset
	token, err := c.userService.GenerateResetToken(req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al procesar la solicitud",
		})
		return
	}

	// Construir URL de reset
	// 1. Intentar desde variable de entorno
	frontendURL := os.Getenv("FRONTEND_URL")

	// 2. Si no está configurada, detectar desde headers de la petición
	if frontendURL == "" {
		// Intentar obtener desde Origin header (más confiable)
		origin := ctx.GetHeader("Origin")

		// Si no hay Origin, intentar desde Referer
		if origin == "" {
			referer := ctx.GetHeader("Referer")
			if referer != "" {
				// Extraer solo protocolo://host:puerto del Referer (sin path)
				// Ejemplo: http://localhost:3000/some/path -> http://localhost:3000
				idx := len("http://")
				if len(referer) > 8 && referer[:8] == "https://" {
					idx = len("https://")
				}

				// Buscar el siguiente / después del protocolo
				pathStart := idx
				for i := idx; i < len(referer); i++ {
					if referer[i] == '/' {
						pathStart = i
						break
					}
				}
				origin = referer[:pathStart]
			}
		}

		if origin != "" {
			frontendURL = origin
		} else {
			// Último recurso: detectar entorno
			env := os.Getenv("ENV")
			if env == "production" {
				frontendURL = "https://localhost:3000"
			} else {
				// Desarrollo: usar HTTPS si el frontend está configurado con HTTPS (mkcert)
				// Por defecto, el frontend usa HTTPS en desarrollo
				frontendURL = "https://localhost:3000"
			}
		}
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	// Preparar datos para el template
	templateData := struct {
		Name     string
		ResetURL string
		Email    string
	}{
		Name:     user.Name,
		ResetURL: resetURL,
		Email:    user.Email,
	}

	// Enviar email
	templatePath := filepath.Join("mailer", "templates", "password_reset.html")
	mailReq := mailer.NewRequest([]string{user.Email}, "Recuperación de Contraseña - MediTrack")

	if err := mailReq.SendMailSkipTLS(templatePath, templateData); err != nil {
		fmt.Printf("Error enviando email: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al enviar el correo de recuperación",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Correo enviado exitosamente. Por favor, revisa tu bandeja de entrada",
	})
}

// ResetPasswordRequest representa la solicitud para cambiar la contraseña
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ResetPassword cambia la contraseña usando un token válido
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	// Resetear contraseña
	if err := c.userService.ResetPassword(req.Token, req.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Token inválido o expirado",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Contraseña actualizada exitosamente",
	})
}

// ValidateResetTokenRequest representa la solicitud para validar un token
type ValidateResetTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// ValidateResetToken verifica si un token es válido
func (c *AuthController) ValidateResetToken(ctx *gin.Context) {
	var req ValidateResetTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Token requerido",
		})
		return
	}

	// Validar token
	user, err := c.userService.ValidateResetToken(req.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Token inválido o expirado",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Token válido",
		Data: map[string]interface{}{
			"email": user.Email,
			"name":  user.Name,
		},
	})
}
