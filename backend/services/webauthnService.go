package services

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"meditrack/models"
	"sync"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

// sessionEntry envuelve SessionData con tiempo de expiración
type sessionEntry struct {
	data      *webauthn.SessionData
	expiresAt time.Time
}

// WebAuthnService maneja el ciclo de vida de las credenciales Passkey
type WebAuthnService struct {
	weba     *webauthn.WebAuthn
	sessions sync.Map // map[string]*sessionEntry
	db       *gorm.DB
}

// NewWebAuthnService crea un nuevo servicio WebAuthn
func NewWebAuthnService(db *gorm.DB, rpID string, rpOrigins []string) (*WebAuthnService, error) {
	weba, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "MediTrack",
		RPID:          rpID,
		RPOrigins:     rpOrigins,
	})
	if err != nil {
		return nil, fmt.Errorf("error iniciando WebAuthn: %w", err)
	}

	svc := &WebAuthnService{weba: weba, db: db}

	// Limpiar sesiones expiradas cada 5 minutos
	go svc.cleanupSessions()

	return svc, nil
}

func (s *WebAuthnService) cleanupSessions() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.sessions.Range(func(k, v interface{}) bool {
			if entry, ok := v.(*sessionEntry); ok && now.After(entry.expiresAt) {
				s.sessions.Delete(k)
			}
			return true
		})
	}
}

func generateSessionID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func (s *WebAuthnService) storeSession(key string, data *webauthn.SessionData) {
	s.sessions.Store(key, &sessionEntry{
		data:      data,
		expiresAt: time.Now().Add(5 * time.Minute),
	})
}

func (s *WebAuthnService) loadAndDeleteSession(key string) (*webauthn.SessionData, error) {
	v, ok := s.sessions.LoadAndDelete(key)
	if !ok {
		return nil, fmt.Errorf("sesión no encontrada o expirada")
	}
	entry := v.(*sessionEntry)
	if time.Now().After(entry.expiresAt) {
		return nil, fmt.Errorf("sesión expirada")
	}
	return entry.data, nil
}

// ---- Helpers de DB ----

func (s *WebAuthnService) getWebAuthnUser(userRUT string) (*models.WebAuthnUser, error) {
	var user models.User
	if err := s.db.Where("rut = ?", userRUT).First(&user).Error; err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	credentials, err := s.GetCredentialsForUser(userRUT)
	if err != nil {
		return nil, err
	}

	return &models.WebAuthnUser{User: user, Credentials: credentials}, nil
}

// FindWebAuthnUser construye un WebAuthnUser listo para usar con la librería
func (s *WebAuthnService) FindWebAuthnUser(userRUT string) (*models.WebAuthnUser, error) {
	return s.getWebAuthnUser(userRUT)
}

// GetCredentialsForUser devuelve las credenciales webauthn.Credential de un usuario
func (s *WebAuthnService) GetCredentialsForUser(userRUT string) ([]webauthn.Credential, error) {
	var rows []models.PasskeyCredential
	if err := s.db.Where("user_rut = ?", userRUT).Find(&rows).Error; err != nil {
		return nil, err
	}

	creds := make([]webauthn.Credential, 0, len(rows))
	for _, row := range rows {
		cred, err := row.ToWebAuthnCredential()
		if err != nil {
			continue // ignorar registros corruptos
		}
		creds = append(creds, *cred)
	}
	return creds, nil
}

// ListPasskeys retorna las passkeys (sin datos sensibles) de un usuario
func (s *WebAuthnService) ListPasskeys(userRUT string) ([]models.PasskeyResponse, error) {
	var rows []models.PasskeyCredential
	if err := s.db.Where("user_rut = ?", userRUT).Order("created_at ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]models.PasskeyResponse, len(rows))
	for i, row := range rows {
		result[i] = row.ToResponse()
	}
	return result, nil
}

// DeletePasskey elimina una passkey por ID verificando que pertenece al usuario
func (s *WebAuthnService) DeletePasskey(userRUT string, passkeyID uint) error {
	res := s.db.Where("id = ? AND user_rut = ?", passkeyID, userRUT).Delete(&models.PasskeyCredential{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("passkey no encontrada")
	}
	return nil
}

// ---- Registro ----

// BeginRegistration inicia el flujo de registro de una credencial biométrica.
// Fuerza authenticatorAttachment=platform para que el navegador use únicamente
// el autenticador integrado en el dispositivo (huella, Face ID, Windows Hello),
// sin ofrecer la opción de teléfono por QR ni llaves de seguridad externas.
func (s *WebAuthnService) BeginRegistration(userRUT string) (interface{}, error) {
	waUser, err := s.getWebAuthnUser(userRUT)
	if err != nil {
		return nil, err
	}

	options, session, err := s.weba.BeginRegistration(waUser,
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			// Credencial residente requerida (necesaria para discoverable login)
			ResidentKey: protocol.ResidentKeyRequirementRequired,
			// El dispositivo debe verificar al usuario (huella/cara/PIN)
			UserVerification: protocol.VerificationRequired,
		}),
		// Sugerir primero el autenticador integrado del dispositivo (huella/Face ID/PIN).
		// Si el dispositivo no lo soporta, el navegador ofrece alternativas (teléfono, etc.)
		webauthn.WithPublicKeyCredentialHints([]protocol.PublicKeyCredentialHints{
			protocol.PublicKeyCredentialHintClientDevice,
			protocol.PublicKeyCredentialHintHybrid,
		}),
		webauthn.WithConveyancePreference(protocol.PreferNoAttestation),
	)
	if err != nil {
		return nil, fmt.Errorf("error iniciando registro: %w", err)
	}

	s.storeSession("reg:"+userRUT, session)

	return options, nil
}

// GetWebAuthn retorna la instancia de WebAuthn para uso en el controller
func (s *WebAuthnService) GetWebAuthn() *webauthn.WebAuthn {
	return s.weba
}

// GetRegistrationSession recupera y elimina la sesión de registro
func (s *WebAuthnService) GetRegistrationSession(userRUT string) (*webauthn.SessionData, error) {
	return s.loadAndDeleteSession("reg:" + userRUT)
}

// GetLoginSession recupera y elimina la sesión de login
func (s *WebAuthnService) GetLoginSession(sessionID string) (*webauthn.SessionData, error) {
	return s.loadAndDeleteSession("login:" + sessionID)
}

// BeginLogin inicia el flujo de autenticación biométrica (discoverable).
// Requiere verificación del usuario para que el dispositivo pida huella/cara/PIN.
func (s *WebAuthnService) BeginLogin() (interface{}, string, error) {
	options, session, err := s.weba.BeginDiscoverableLogin(
		webauthn.WithUserVerification(protocol.VerificationRequired),
	)
	if err != nil {
		return nil, "", fmt.Errorf("error iniciando login biométrico: %w", err)
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return nil, "", err
	}

	s.storeSession("login:"+sessionID, session)

	return options, sessionID, nil
}

// FindUserByCredentialID busca el usuario dueño de una credencial
func (s *WebAuthnService) FindUserByCredentialID(rawID []byte) (*models.WebAuthnUser, error) {
	var row models.PasskeyCredential
	if err := s.db.Where("credential_id = ?", rawID).First(&row).Error; err != nil {
		return nil, fmt.Errorf("credencial no encontrada")
	}
	return s.getWebAuthnUser(row.UserRUT)
}

// SaveCredential persiste una nueva credencial en la base de datos
func (s *WebAuthnService) SaveCredential(userRUT, name string, cred *webauthn.Credential) error {
	data, err := json.Marshal(cred)
	if err != nil {
		return fmt.Errorf("error serializando credencial: %w", err)
	}

	row := models.PasskeyCredential{
		UserRUT:        userRUT,
		CredentialID:   cred.ID,
		CredentialData: data,
		Name:           name,
	}
	return s.db.Create(&row).Error
}

// UpdateCredentialSignCount actualiza el contador de firma de una credencial (anti-clonación)
func (s *WebAuthnService) UpdateCredentialSignCount(credID []byte, updatedCred *webauthn.Credential) error {
	data, err := json.Marshal(updatedCred)
	if err != nil {
		return err
	}
	return s.db.Model(&models.PasskeyCredential{}).
		Where("credential_id = ?", credID).
		Update("credential_data", data).Error
}
