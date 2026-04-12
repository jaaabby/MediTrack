package models

import (
	"encoding/json"

	"github.com/go-webauthn/webauthn/webauthn"
)

// PasskeyCredential almacena una credencial WebAuthn/Passkey de un usuario
type PasskeyCredential struct {
	ID             uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	UserRUT        string `json:"user_rut" gorm:"not null;index;column:user_rut"`
	CredentialID   []byte `json:"-" gorm:"not null;uniqueIndex;column:credential_id"`
	CredentialData []byte `json:"-" gorm:"not null;column:credential_data"`
	Name           string `json:"name" gorm:"column:name;default:''"`
	CreatedAt      int64  `json:"created_at" gorm:"autoCreateTime;column:created_at"`
	UpdatedAt      int64  `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
}

func (p PasskeyCredential) TableName() string { return "passkey_credential" }

// ToWebAuthnCredential deserializa el credential almacenado
func (p *PasskeyCredential) ToWebAuthnCredential() (*webauthn.Credential, error) {
	var cred webauthn.Credential
	if err := json.Unmarshal(p.CredentialData, &cred); err != nil {
		return nil, err
	}
	return &cred, nil
}

// PasskeyResponse es la representación pública de una passkey (sin datos sensibles)
type PasskeyResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}

func (p PasskeyCredential) ToResponse() PasskeyResponse {
	return PasskeyResponse{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
	}
}

// WebAuthnUser adapta el modelo User para implementar webauthn.User
type WebAuthnUser struct {
	User        User
	Credentials []webauthn.Credential
}

func (u *WebAuthnUser) WebAuthnID() []byte          { return []byte(u.User.RUT) }
func (u *WebAuthnUser) WebAuthnName() string         { return u.User.Email }
func (u *WebAuthnUser) WebAuthnDisplayName() string  { return u.User.Name }
func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential { return u.Credentials }
