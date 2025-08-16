package services

import (
	"context"

	"meditrack/models"
	"meditrack/repository"
)

// OperatingRoomService define los servicios para salas de operación
type OperatingRoomService interface {
	CreateOperatingRoom(ctx context.Context, room *models.OperatingRoom) error
	GetOperatingRoomByID(ctx context.Context, id string) (*models.OperatingRoom, error)
	GetAllOperatingRooms(ctx context.Context) ([]*models.OperatingRoom, error)
	UpdateOperatingRoom(ctx context.Context, room *models.OperatingRoom) error
	DeleteOperatingRoom(ctx context.Context, id string) error
}

// operatingRoomService implementa OperatingRoomService
type operatingRoomService struct {
	operatingRoomRepo repository.OperatingRoomRepository
}

// NewOperatingRoomService crea una nueva instancia de OperatingRoomService
func NewOperatingRoomService(operatingRoomRepo repository.OperatingRoomRepository) OperatingRoomService {
	return &operatingRoomService{
		operatingRoomRepo: operatingRoomRepo,
	}
}

// CreateOperatingRoom implementa la creación de una sala de operación
func (s *operatingRoomService) CreateOperatingRoom(ctx context.Context, room *models.OperatingRoom) error {
	return s.operatingRoomRepo.Create(ctx, room)
}

// GetOperatingRoomByID implementa la obtención de una sala de operación por ID
func (s *operatingRoomService) GetOperatingRoomByID(ctx context.Context, id string) (*models.OperatingRoom, error) {
	return s.operatingRoomRepo.GetByID(ctx, id)
}

// GetAllOperatingRooms implementa la obtención de todas las salas de operación
func (s *operatingRoomService) GetAllOperatingRooms(ctx context.Context) ([]*models.OperatingRoom, error) {
	return s.operatingRoomRepo.GetAll(ctx)
}

// UpdateOperatingRoom implementa la actualización de una sala de operación
func (s *operatingRoomService) UpdateOperatingRoom(ctx context.Context, room *models.OperatingRoom) error {
	return s.operatingRoomRepo.Update(ctx, room)
}

// DeleteOperatingRoom implementa la eliminación de una sala de operación
func (s *operatingRoomService) DeleteOperatingRoom(ctx context.Context, id string) error {
	return s.operatingRoomRepo.Delete(ctx, id)
}
