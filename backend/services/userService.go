package services

import (
	"meditrack/models"
	"meditrack/repository"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(rut string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(rut string) error
}

type userService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(user *models.User) error {
	return s.repo.Create(user)
}

func (s *userService) GetUserByID(rut string) (*models.User, error) {
	return s.repo.GetByID(rut)
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *userService) UpdateUser(user *models.User) error {
	return s.repo.Update(user)
}

func (s *userService) DeleteUser(rut string) error {
	return s.repo.Delete(rut)
}
