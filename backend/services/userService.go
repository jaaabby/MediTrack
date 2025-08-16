package services

import (
	"context"

	"meditrack/models"
	"meditrack/repository"
)

// UserService define los servicios para usuarios
type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
}

// userService implementa UserService
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService crea una nueva instancia de UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser implementa la creación de un usuario
func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Create(ctx, user)
}

// GetUserByID implementa la obtención de un usuario por ID
func (s *userService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetUserByEmail implementa la obtención de un usuario por email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// GetAllUsers implementa la obtención de todos los usuarios
func (s *userService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.GetAll(ctx)
}

// UpdateUser implementa la actualización de un usuario
func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

// DeleteUser implementa la eliminación de un usuario
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
