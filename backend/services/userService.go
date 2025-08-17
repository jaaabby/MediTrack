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
	GetUsersBySpecialty(ctx context.Context, specialty string) ([]*models.User, error)
}

// userService implementa UserService
type userService struct {
	userRepo repository.UserRepositoryInterface
}

// NewUserService crea una nueva instancia de UserService
func NewUserService(userRepo repository.UserRepositoryInterface) UserService {
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

// GetUsersBySpecialty implementa la obtención de usuarios por especialidad y rol doctor
func (s *userService) GetUsersBySpecialty(ctx context.Context, specialty string) ([]*models.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var doctors []*models.User
	for _, u := range users {
		if u.Role == "doctor" && u.Name == specialty {
			doctors = append(doctors, u)
		}
	}
	return doctors, nil
}
