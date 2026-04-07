package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

// Parámetros scrypt — OWASP Interactive Login mínimo: N=32768, r=8, p=1 (~32MB RAM/hash)
const (
	scryptN      = 32768
	scryptR      = 8
	scryptP      = 1
	scryptKeyLen = 32
	scryptSaltLen = 32
)

// HashPassword genera un hash scrypt con salt aleatorio.
// Formato resultante: $scrypt$N=32768,r=8,p=1$<salt_b64>$<hash_b64>
func HashPassword(password string) (string, error) {
	salt := make([]byte, scryptSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("error generando salt: %w", err)
	}

	hash, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, scryptKeyLen)
	if err != nil {
		return "", fmt.Errorf("error calculando scrypt: %w", err)
	}

	encoded := fmt.Sprintf(
		"$scrypt$N=%d,r=%d,p=%d$%s$%s",
		scryptN, scryptR, scryptP,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	return encoded, nil
}

// ComparePassword verifica una contraseña contra su hash almacenado.
// Soporta tanto hashes scrypt (nuevos) como bcrypt (migración de usuarios existentes).
func ComparePassword(storedHash, password string) error {
	if strings.HasPrefix(storedHash, "$scrypt$") {
		return compareScrypt(storedHash, password)
	}
	// Fallback bcrypt para usuarios existentes (hashes que empiezan con $2a$ o $2b$)
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
}

// IsScryptHash indica si un hash almacenado ya usa scrypt.
func IsScryptHash(h string) bool {
	return strings.HasPrefix(h, "$scrypt$")
}

// compareScrypt verifica una contraseña contra un hash en formato $scrypt$...
func compareScrypt(storedHash, password string) error {
	// Formato: $scrypt$N=<n>,r=<r>,p=<p>$<salt_b64>$<hash_b64>
	parts := strings.Split(storedHash, "$")
	// parts[0]="" parts[1]="scrypt" parts[2]="N=...,r=...,p=..." parts[3]=salt parts[4]=hash
	if len(parts) != 5 || parts[1] != "scrypt" {
		return errors.New("formato de hash inválido")
	}

	n, r, p, err := parseScryptParams(parts[2])
	if err != nil {
		return err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return errors.New("salt inválido")
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return errors.New("hash inválido")
	}

	computed, err := scrypt.Key([]byte(password), salt, n, r, p, len(expectedHash))
	if err != nil {
		return fmt.Errorf("error calculando scrypt: %w", err)
	}

	if subtle.ConstantTimeCompare(computed, expectedHash) != 1 {
		return errors.New("contraseña incorrecta")
	}
	return nil
}

// parseScryptParams parsea "N=32768,r=8,p=1" y retorna los valores enteros.
func parseScryptParams(params string) (n, r, p int, err error) {
	for _, part := range strings.Split(params, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return 0, 0, 0, errors.New("parámetros scrypt inválidos")
		}
		val, e := strconv.Atoi(kv[1])
		if e != nil {
			return 0, 0, 0, fmt.Errorf("valor inválido en parámetros scrypt: %w", e)
		}
		switch kv[0] {
		case "N":
			n = val
		case "r":
			r = val
		case "p":
			p = val
		}
	}
	if n == 0 || r == 0 || p == 0 {
		return 0, 0, 0, errors.New("parámetros scrypt incompletos")
	}
	return n, r, p, nil
}
