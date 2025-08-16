package services

import (
	"context"

	"meditrack/models"
	"meditrack/repository"
)

// SupplyMovementService define los servicios para movimientos de insumos
type SupplyMovementService interface {
	CreateSupplyMovement(ctx context.Context, movement *models.SupplyMovement) error
	GetSupplyMovementByID(ctx context.Context, id string) (*models.SupplyMovement, error)
	GetSupplyMovementsByMedicalSupplyID(ctx context.Context, supplyID string) ([]*models.SupplyMovement, error)
	GetSupplyMovementsByLocationID(ctx context.Context, locationID string) ([]*models.SupplyMovement, error)
	GetSupplyMovementsByStatus(ctx context.Context, status string) ([]*models.SupplyMovement, error)
	GetAllSupplyMovements(ctx context.Context) ([]*models.SupplyMovement, error)
	UpdateSupplyMovement(ctx context.Context, movement *models.SupplyMovement) error
	DeleteSupplyMovement(ctx context.Context, id string) error
}

// supplyMovementService implementa SupplyMovementService
type supplyMovementService struct {
	supplyMovementRepo repository.SupplyMovementRepository
}

// NewSupplyMovementService crea una nueva instancia de SupplyMovementService
func NewSupplyMovementService(supplyMovementRepo repository.SupplyMovementRepository) SupplyMovementService {
	return &supplyMovementService{
		supplyMovementRepo: supplyMovementRepo,
	}
}

// CreateSupplyMovement implementa la creación de un movimiento de insumo
func (s *supplyMovementService) CreateSupplyMovement(ctx context.Context, movement *models.SupplyMovement) error {
	return s.supplyMovementRepo.Create(ctx, movement)
}

// GetSupplyMovementByID implementa la obtención de un movimiento por ID
func (s *supplyMovementService) GetSupplyMovementByID(ctx context.Context, id string) (*models.SupplyMovement, error) {
	return s.supplyMovementRepo.GetByID(ctx, id)
}

// GetSupplyMovementsByMedicalSupplyID implementa la obtención de movimientos por ID de insumo
func (s *supplyMovementService) GetSupplyMovementsByMedicalSupplyID(ctx context.Context, supplyID string) ([]*models.SupplyMovement, error) {
	return s.supplyMovementRepo.GetByMedicalSupplyID(ctx, supplyID)
}

// GetSupplyMovementsByLocationID implementa la obtención de movimientos por ID de ubicación
func (s *supplyMovementService) GetSupplyMovementsByLocationID(ctx context.Context, locationID string) ([]*models.SupplyMovement, error) {
	return s.supplyMovementRepo.GetByLocationID(ctx, locationID)
}

// GetSupplyMovementsByStatus implementa la obtención de movimientos por estado
func (s *supplyMovementService) GetSupplyMovementsByStatus(ctx context.Context, status string) ([]*models.SupplyMovement, error) {
	return s.supplyMovementRepo.GetByStatus(ctx, status)
}

// GetAllSupplyMovements implementa la obtención de todos los movimientos
func (s *supplyMovementService) GetAllSupplyMovements(ctx context.Context) ([]*models.SupplyMovement, error) {
	return s.supplyMovementRepo.GetAll(ctx)
}

// UpdateSupplyMovement implementa la actualización de un movimiento
func (s *supplyMovementService) UpdateSupplyMovement(ctx context.Context, movement *models.SupplyMovement) error {
	return s.supplyMovementRepo.Update(ctx, movement)
}

// DeleteSupplyMovement implementa la eliminación de un movimiento
func (s *supplyMovementService) DeleteSupplyMovement(ctx context.Context, id string) error {
	return s.supplyMovementRepo.Delete(ctx, id)
}
