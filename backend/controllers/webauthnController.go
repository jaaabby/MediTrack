package controllers

import (
	"fmt"
	"net/http"
	"time"

	"meditrack/config"
	"meditrack/models"
	"meditrack/pkg/response"
	"meditrack/services"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
)

const jwtDuration = 24 * time.Hour

// WebAuthnController maneja los endpoints de registro y login con Passkey
type WebAuthnController struct {
	webauthnService *services.WebAuthnService
	userService     services.UserService
	secretKey       string
}

// NewWebAuthnController crea una nueva instancia
func NewWebAuthnController(
	webauthnService *services.WebAuthnService,
	userService services.UserService,
	secretKey string,
) *WebAuthnController {
	return &WebAuthnController{
		webauthnService: webauthnService,
		userService:     userService,
		secretKey:       secretKey,
	}
}

// ──────────────────────────────────────────────────────────────────────────
// REGISTRO  (rutas protegidas – usuario ya autenticado con JWT)
// ──────────────────────────────────────────────────────────────────────────

// BeginRegistration inicia el registro de una nueva passkey
// POST /auth/passkey/register/begin
func (c *WebAuthnController) BeginRegistration(ctx *gin.Context) {
	userRUT, err := getContextUserRUT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{Success: false, Error: err.Error()})
		return
	}

	options, err := c.webauthnService.BeginRegistration(userRUT)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Inicia el registro de passkey en tu dispositivo",
		Data:    options,
	})
}

// FinishRegistration verifica la respuesta del autenticador y guarda la credencial
// POST /auth/passkey/register/finish?name=MiDispositivo
// Body: PublicKeyCredential JSON tal como lo devuelve el navegador
func (c *WebAuthnController) FinishRegistration(ctx *gin.Context) {
	userRUT, err := getContextUserRUT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{Success: false, Error: err.Error()})
		return
	}

	name := ctx.Query("name")
	if name == "" {
		name = "Mi Passkey"
	}

	session, err := c.webauthnService.GetRegistrationSession(userRUT)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Sesión de registro expirada, inicia el proceso nuevamente",
		})
		return
	}

	waUser, err := c.webauthnService.FindWebAuthnUser(userRUT)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: err.Error()})
		return
	}

	credential, err := c.webauthnService.GetWebAuthn().FinishRegistration(waUser, *session, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Credencial inválida: " + err.Error(),
		})
		return
	}

	if err := c.webauthnService.SaveCredential(userRUT, name, credential); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error guardando passkey: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Passkey registrada exitosamente",
	})
}

// ──────────────────────────────────────────────────────────────────────────
// LOGIN  (rutas públicas)
// ──────────────────────────────────────────────────────────────────────────

// BeginLogin inicia la autenticación passkey (discoverable – sin email)
// POST /auth/passkey/login/begin
func (c *WebAuthnController) BeginLogin(ctx *gin.Context) {
	options, sessionID, err := c.webauthnService.BeginLogin()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Autentícate con tu passkey",
		Data: map[string]interface{}{
			"options":    options,
			"session_id": sessionID,
		},
	})
}

// FinishLogin verifica la respuesta del autenticador y devuelve un JWT
// POST /auth/passkey/login/finish?session_id=xxx
// Body: PublicKeyCredential JSON tal como lo devuelve el navegador
func (c *WebAuthnController) FinishLogin(ctx *gin.Context) {
	sessionID := ctx.Query("session_id")
	if sessionID == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Falta el parámetro session_id",
		})
		return
	}

	session, err := c.webauthnService.GetLoginSession(sessionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Sesión de login expirada, inicia el proceso nuevamente",
		})
		return
	}

	var resolvedUser *models.WebAuthnUser

	handler := func(rawID, userHandle []byte) (webauthn.User, error) {
		// Primero buscar por credentialID
		u, err := c.webauthnService.FindUserByCredentialID(rawID)
		if err != nil {
			// Fallback: user handle = RUT
			u, err = c.webauthnService.FindWebAuthnUser(string(userHandle))
			if err != nil {
				return nil, fmt.Errorf("usuario no encontrado")
			}
		}
		resolvedUser = u
		return u, nil
	}

	credential, err := c.webauthnService.GetWebAuthn().FinishDiscoverableLogin(handler, *session, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   "Autenticación con passkey fallida: " + err.Error(),
		})
		return
	}

	if resolvedUser == nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{Success: false, Error: "Usuario no identificado"})
		return
	}

	// Actualizar sign count (protección anti-clonación)
	_ = c.webauthnService.UpdateCredentialSignCount(credential.ID, credential)

	user := resolvedUser.User
	if !user.IsActive {
		ctx.JSON(http.StatusUnauthorized, response.Response{Success: false, Error: "Usuario inactivo"})
		return
	}

	token, tokenErr := config.GenerateToken(
		user.RUT, user.Email, user.Role, user.TokenVersion, c.secretKey, jwtDuration,
	)
	if tokenErr != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error generando token",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Autenticación con passkey exitosa",
		Data: LoginResponse{
			Token:     token,
			User:      user.ToResponse(),
			ExpiresIn: int64(jwtDuration.Seconds()),
			TokenType: "Bearer",
		},
	})
}

// ──────────────────────────────────────────────────────────────────────────
// GESTIÓN  (rutas protegidas)
// ──────────────────────────────────────────────────────────────────────────

// ListPasskeys devuelve las passkeys registradas del usuario autenticado
// GET /auth/passkey/credentials
func (c *WebAuthnController) ListPasskeys(ctx *gin.Context) {
	userRUT, err := getContextUserRUT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{Success: false, Error: err.Error()})
		return
	}

	passkeys, err := c.webauthnService.ListPasskeys(userRUT)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{Success: false, Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{Success: true, Data: passkeys})
}

// DeletePasskey elimina una passkey del usuario autenticado
// DELETE /auth/passkey/credentials/:id
func (c *WebAuthnController) DeletePasskey(ctx *gin.Context) {
	userRUT, err := getContextUserRUT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{Success: false, Error: err.Error()})
		return
	}

	var passkeyID uint
	if _, scanErr := fmt.Sscanf(ctx.Param("id"), "%d", &passkeyID); scanErr != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Success: false, Error: "ID inválido"})
		return
	}

	if err := c.webauthnService.DeletePasskey(userRUT, passkeyID); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{Success: false, Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Passkey eliminada exitosamente",
	})
}

// ──────────────────────────────────────────────────────────────────────────
// helpers
// ──────────────────────────────────────────────────────────────────────────

func getContextUserRUT(ctx *gin.Context) (string, error) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return "", fmt.Errorf("usuario no autenticado")
	}
	return userID.(string), nil
}
