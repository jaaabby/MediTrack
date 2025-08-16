package services

import (
	"context"

	"meditrack/models"
	"meditrack/repository"
)

// SupplyRouteMedicalSupplyService define los servicios para la relación ruta de suministro-insumo médico
type SupplyRouteMedicalSupplyService interface {
	CreateSupplyRouteMedicalSupply(ctx context.Context, relation *models.SupplyRouteMedicalSupply) error
	GetSupplyRouteMedicalSupplyByID(ctx context.Context, id string) (*models.SupplyRouteMedicalSupply, error)
	GetSupplyRouteMedicalSuppliesBySupplyRouteID(ctx context.Context, routeID string) ([]*models.SupplyRouteMedicalSupply, error)
	GetSupplyRouteMedicalSuppliesByMedicalSupplyID(ctx context.Context, supplyID string) ([]*models.SupplyRouteMedicalSupply, error)
	GetAllSupplyRouteMedicalSupplies(ctx context.Context) ([]*models.SupplyRouteMedicalSupply, error)
	UpdateSupplyRouteMedicalSupply(ctx context.Context, relation *models.SupplyRouteMedicalSupply) error
	DeleteSupplyRouteMedicalSupply(ctx context.Context, id string) error
}

// supplyRouteMedicalSupplyService implementa SupplyRouteMedicalSupplyService
type supplyRouteMedicalSupplyService struct {
	supplyRouteMedicalSupplyRepo repository.SupplyRouteMedicalSupplyRepository
}

// NewSupplyRouteMedicalSupplyService crea una nueva instancia de SupplyRouteMedicalSupplyService
func NewSupplyRouteMedicalSupplyService(supplyRouteMedicalSupplyRepo repository.SupplyRouteMedicalSupplyRepository) SupplyRouteMedicalSupplyService {
	return &supplyRouteMedicalSupplyService{
		supplyRouteMedicalSupplyRepo: supplyRouteMedicalSupplyRepo,
	}
}

// CreateSupplyRouteMedicalSupply implementa la creación de una relación ruta-insumo
func (s *supplyRouteMedicalSupplyService) CreateSupplyRouteMedicalSupply(ctx context.Context, relation *models.SupplyRouteMedicalSupply) error {
	return s.supplyRouteMedicalSupplyRepo.Create(ctx, relation)
}

// GetSupplyRouteMedicalSupplyByID implementa la obtención de una relación ruta-insumo por ID
func (s *supplyRouteMedicalSupplyService) GetSupplyRouteMedicalSupplyByID(ctx context.Context, id string) (*models.SupplyRouteMedicalSupply, error) {
	return s.supplyRouteMedicalSupplyRepo.GetByID(ctx, id)
}

// GetSupplyRouteMedicalSuppliesBySupplyRouteID implementa la obtención de relaciones por ID de ruta de suministro
func (s *supplyRouteMedicalSupplyService) GetSupplyRouteMedicalSuppliesBySupplyRouteID(ctx context.Context, routeID string) ([]*models.SupplyRouteMedicalSupply, error) {
	return s.supplyRouteMedicalSupplyRepo.GetBySupplyRouteID(ctx, routeID)
}

// GetSupplyRouteMedicalSuppliesByMedicalSupplyID implementa la obtención de relaciones por ID de insumo médico
func (s *supplyRouteMedicalSupplyService) GetSupplyRouteMedicalSuppliesByMedicalSupplyID(ctx context.Context, supplyID string) ([]*models.SupplyRouteMedicalSupply, error) {
	return s.supplyRouteMedicalSupplyRepo.GetByMedicalSupplyID(ctx, supplyID)
}

// GetAllSupplyRouteMedicalSupplies implementa la obtención de todas las relaciones ruta-insumo
func (s *supplyRouteMedicalSupplyService) GetAllSupplyRouteMedicalSupplies(ctx context.Context) ([]*models.SupplyRouteMedicalSupply, error) {
	return s.supplyRouteMedicalSupplyRepo.GetAll(ctx)
}

// UpdateSupplyRouteMedicalSupply implementa la actualización de una relación ruta-insumo
func (s *supplyRouteMedicalSupplyService) UpdateSupplyRouteMedicalSupply(ctx context.Context, relation *models.SupplyRouteMedicalSupply) error {
	return s.supplyRouteMedicalSupplyRepo.Update(ctx, relation)
}

// DeleteSupplyRouteMedicalSupply implementa la eliminación de una relación ruta-insumo
func (s *supplyRouteMedicalSupplyService) DeleteSupplyRouteMedicalSupply(ctx context.Context, id string) error {
	return s.supplyRouteMedicalSupplyRepo.Delete(ctx, id)
}
