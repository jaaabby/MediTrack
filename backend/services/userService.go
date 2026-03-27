package services

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"meditrack/models"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// generateRandomPassword genera una contraseña aleatoria de 12 caracteres
// con al menos una mayúscula, una minúscula, un número y un carácter especial
func generateRandomPassword() (string, error) {
	const (
		length           = 12
		lowercaseLetters = "abcdefghijkmnopqrstuvwxyz"
		uppercaseLetters = "ABCDEFGHJKLMNPQRSTUVWXYZ"
		digits           = "23456789"
		specialChars     = "@#$%&*+-="
	)

	allChars := lowercaseLetters + uppercaseLetters + digits + specialChars
	password := make([]byte, length)

	// Asegurar al menos un carácter de cada tipo
	charSets := []string{lowercaseLetters, uppercaseLetters, digits, specialChars}
	for i := 0; i < 4; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSets[i]))))
		if err != nil {
			return "", err
		}
		password[i] = charSets[i][n.Int64()]
	}

	// Llenar el resto con caracteres aleatorios
	for i := 4; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}
		password[i] = allChars[n.Int64()]
	}

	// Mezclar la contraseña
	for i := len(password) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password), nil
}

// removeAccents elimina tildes y acentos de un string usando normalización Unicode
func removeAccents(s string) string {
	// Primero normalizar a NFD (descomponer caracteres con acentos)
	t := strings.ToLower(s)

	// Reemplazar caracteres acentuados manualmente
	replacements := map[rune]rune{
		'á': 'a', 'é': 'e', 'í': 'i', 'ó': 'o', 'ú': 'u',
		'à': 'a', 'è': 'e', 'ì': 'i', 'ò': 'o', 'ù': 'u',
		'ä': 'a', 'ë': 'e', 'ï': 'i', 'ö': 'o', 'ü': 'u',
		'â': 'a', 'ê': 'e', 'î': 'i', 'ô': 'o', 'û': 'u',
		'ñ': 'n', 'ç': 'c',
	}

	var result strings.Builder
	for _, r := range t {
		if replacement, ok := replacements[r]; ok {
			result.WriteRune(replacement)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.DB.Create(user).Error
}

// CreateUserWithTemporaryPassword crea un usuario con una contraseña temporal generada automáticamente
func (s *UserService) CreateUserWithTemporaryPassword(user *models.User) (string, error) {
	// Generar contraseña temporal
	tempPassword, err := generateRandomPassword()
	if err != nil {
		return "", err
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Asignar contraseña hasheada y marcar que debe cambiarla
	user.Password = string(hashedPassword)
	user.MustChangePassword = true
	user.IsActive = true

	// Crear usuario en la base de datos
	if err := s.DB.Create(user).Error; err != nil {
		return "", err
	}

	// Retornar la contraseña en texto plano para enviarla por correo
	return tempPassword, nil
}

func (s *UserService) GetUserByRut(rut string) (*models.User, error) {
	var user models.User
	if err := s.DB.Preload("MedicalCenter").Preload("Specialty").First(&user, "rut = ?", rut).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) DeleteUser(rut string) error {
	return s.DB.Delete(&models.User{}, "rut = ?", rut).Error
}

func (s *UserService) UpdateUser(rut string, newUser *models.User) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, "rut = ?", rut).Error; err != nil {
		return nil, err
	}

	// Usar Updates para campos generales, pero Update individual para is_active
	// porque Updates ignora valores zero (false)
	updates := map[string]interface{}{
		"name":                 newUser.Name,
		"email":                newUser.Email,
		"role":                 newUser.Role,
		"medical_center_id":    newUser.MedicalCenterID,
		"is_active":            newUser.IsActive,           // Explícitamente actualizar is_active
		"must_change_password": newUser.MustChangePassword, // Explícitamente actualizar must_change_password
	}

	// Actualizar campos opcionales
	if newUser.PavilionID != nil {
		updates["pavilion_id"] = newUser.PavilionID
	}
	if newUser.SpecialtyID != nil {
		updates["specialty_id"] = newUser.SpecialtyID
	}

	// Solo actualizar password si no está vacío
	if newUser.Password != "" {
		updates["password"] = newUser.Password
	}

	if err := s.DB.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Recargar el usuario actualizado
	if err := s.DB.First(&user, "rut = ?", rut).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := s.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.DB.Preload("MedicalCenter").Preload("Specialty").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUsersByRole(role string) ([]models.User, error) {
	var users []models.User
	if err := s.DB.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) DeactivateUser(rut string) error {
	if err := s.DB.Model(&models.User{}).Where("rut = ?", rut).Update("is_active", false).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserService) ActivateUser(rut string) error {
	if err := s.DB.Model(&models.User{}).Where("rut = ?", rut).Update("is_active", true).Error; err != nil {
		return err
	}
	return nil
}

// SearchUsers busca usuarios por nombre, RUT o email (insensible a tildes y mayúsculas)
func (s *UserService) SearchUsers(searchTerm string) ([]models.User, error) {
	var users []models.User

	// Normalizar el término de búsqueda (quitar tildes y convertir a minúsculas)
	normalizedSearch := removeAccents(strings.TrimSpace(searchTerm))
	if normalizedSearch == "" {
		return []models.User{}, nil
	}

	// Obtener todos los usuarios activos (limitamos a 200 para no sobrecargar)
	if err := s.DB.Where("is_active = ?", true).
		Order("name ASC").
		Limit(200).
		Find(&users).Error; err != nil {
		return nil, err
	}

	// Filtrar en memoria comparando versiones normalizadas
	var filteredUsers []models.User
	for _, user := range users {
		normalizedName := removeAccents(user.Name)
		normalizedRut := removeAccents(user.RUT)
		normalizedEmail := removeAccents(user.Email)

		if strings.Contains(normalizedName, normalizedSearch) ||
			strings.Contains(normalizedRut, normalizedSearch) ||
			strings.Contains(normalizedEmail, normalizedSearch) {
			filteredUsers = append(filteredUsers, user)
		}

		// Limitar a 50 resultados
		if len(filteredUsers) >= 50 {
			break
		}
	}

	return filteredUsers, nil
}

// GenerateResetToken genera un token único para recuperación de contraseña
func (s *UserService) GenerateResetToken(email string) (string, error) {
	// Buscar usuario por email
	var user models.User
	if err := s.DB.Where("email = ? AND is_active = ?", email, true).First(&user).Error; err != nil {
		return "", err
	}

	// Generar token aleatorio (32 bytes = 64 caracteres hex)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)

	// Establecer tiempo de expiración (1 hora desde ahora)
	expiresAt := time.Now().Add(1 * time.Hour).Unix()

	// Guardar token y expiración en la base de datos
	if err := s.DB.Model(&user).Updates(map[string]interface{}{
		"reset_password_token":      token,
		"reset_password_expires_at": expiresAt,
	}).Error; err != nil {
		return "", err
	}

	return token, nil
}

// ValidateResetToken valida si un token de reset es válido y no ha expirado
func (s *UserService) ValidateResetToken(token string) (*models.User, error) {
	var user models.User

	// Buscar usuario por token
	if err := s.DB.Where("reset_password_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}

	// Verificar si el token ha expirado
	if user.ResetPasswordExpiresAt == nil || *user.ResetPasswordExpiresAt < time.Now().Unix() {
		return nil, gorm.ErrRecordNotFound // Token expirado
	}

	return &user, nil
}

// ResetPassword cambia la contraseña del usuario usando un token de reset
func (s *UserService) ResetPassword(token, newPassword string) error {
	// Validar token
	user, err := s.ValidateResetToken(token)
	if err != nil {
		return err
	}

	// Hashear la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Actualizar contraseña y limpiar token
	if err := s.DB.Model(&user).Updates(map[string]interface{}{
		"password":                  string(hashedPassword),
		"reset_password_token":      nil,
		"reset_password_expires_at": nil,
	}).Error; err != nil {
		return err
	}

	return nil
}

// ClearResetToken limpia el token de reset de un usuario
func (s *UserService) ClearResetToken(rut string) error {
	return s.DB.Model(&models.User{}).Where("rut = ?", rut).Updates(map[string]interface{}{
		"reset_password_token":      nil,
		"reset_password_expires_at": nil,
	}).Error
}

// CreateOtpSession persiste una nueva sesión OTP en la base de datos
func (s *UserService) CreateOtpSession(session *models.OtpSession) error {
	return s.DB.Create(session).Error
}

// GetOtpSession obtiene una sesión OTP por su ID
func (s *UserService) GetOtpSession(id string) (*models.OtpSession, error) {
	var session models.OtpSession
	if err := s.DB.First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// MarkOtpSessionUsed marca una sesión OTP como utilizada
func (s *UserService) MarkOtpSessionUsed(id string) error {
	return s.DB.Model(&models.OtpSession{}).Where("id = ?", id).Update("used", true).Error
}

// IncrementOtpAttempts incrementa el contador de intentos fallidos y retorna el nuevo total
func (s *UserService) IncrementOtpAttempts(id string) (int, error) {
	if err := s.DB.Model(&models.OtpSession{}).Where("id = ?", id).
		UpdateColumn("attempts", gorm.Expr("attempts + 1")).Error; err != nil {
		return 0, err
	}
	var session models.OtpSession
	if err := s.DB.Select("attempts").First(&session, "id = ?", id).Error; err != nil {
		return 0, err
	}
	return session.Attempts, nil
}

// InvalidateOtpSession marca una sesión OTP como usada (para bloquearla)
func (s *UserService) InvalidateOtpSession(id string) error {
	return s.DB.Model(&models.OtpSession{}).Where("id = ?", id).Update("used", true).Error
}

// DeleteExpiredOtpSessions elimina las sesiones OTP expiradas o ya usadas
func (s *UserService) DeleteExpiredOtpSessions() error {
	return s.DB.Where("expires_at < ? OR used = ?", time.Now().Unix(), true).
		Delete(&models.OtpSession{}).Error
}
