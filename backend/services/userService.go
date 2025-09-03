package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.DB.Create(user).Error
}

func (s *UserService) GetUserByID(rut string) (*models.User, error) {
	var user models.User
	if err := s.DB.Preload("MedicalCenter").First(&user, "rut = ?", rut).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) DeleteUser(rut string) error {
	if err := s.DB.Delete(&models.User{}, "rut = ?", rut).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateUser(rut string, newUser *models.User) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, "rut = ?", rut).Error; err != nil {
		return nil, err
	}
	user.Name = newUser.Name
	user.Email = newUser.Email
	user.Password = newUser.Password
	user.Role = newUser.Role
	user.MedicalCenterID = newUser.MedicalCenterID

	if err := s.DB.Save(&user).Error; err != nil {
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
	if err := s.DB.Preload("MedicalCenter").First(&user, "email = ?", email).Error; err != nil {
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

func (s *UserService) GetUserProfileByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.DB.Preload("MedicalCenter").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
