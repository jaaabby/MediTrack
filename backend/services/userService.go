package services

import (
	"meditrack/models"
	"strings"

	"gorm.io/gorm"
)

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

	if err := s.DB.Model(&user).Omit("rut", "created_at", "updated_at").Updates(newUser).Error; err != nil {
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

