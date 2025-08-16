package services

import (
	"context"

	"meditrack/models"
	"meditrack/repository"
)

// SupplyRouteService define los servicios para rutas de suministro
type SupplyRouteService interface {
	CreateSupplyRoute(ctx context.Context, route *models.SupplyRoute) error
	GetSupplyRouteByID(ctx context.Context, id string) (*models.SupplyRoute, error)
	GetSupplyRoutesByPatientID(ctx context.Context, patientID string) ([]*models.SupplyRoute, error)
	GetSupplyRoutesByOperatingRoomID(ctx context.Context, roomID string) ([]*models.SupplyRoute, error)
	GetAllSupplyRoutes(ctx context.Context) ([]*models.SupplyRoute, error)
	UpdateSupplyRoute(ctx context.Context, route *models.SupplyRoute) error
	DeleteSupplyRoute(ctx context.Context, id string) error
}

// supplyRouteService implementa SupplyRouteService
type supplyRouteService struct {
	supplyRouteRepo repository.SupplyRouteRepository
}

// NewSupplyRouteService crea una nueva instancia de SupplyRouteService
func NewSupplyRouteService(supplyRouteRepo repository.SupplyRouteRepository) SupplyRouteService {
	return &supplyRouteService{
		supplyRouteRepo: supplyRouteRepo,
	}
}

// CreateSupplyRoute implementa la creación de una ruta de suministro
func (s *supplyRouteService) CreateSupplyRoute(ctx context.Context, route *models.SupplyRoute) error {
	return s.supplyRouteRepo.Create(ctx, route)
}

// GetSupplyRouteByID implementa la obtención de una ruta de suministro por ID
func (s *supplyRouteService) GetSupplyRouteByID(ctx context.Context, id string) (*models.SupplyRoute, error) {
	return s.supplyRouteRepo.GetByID(ctx, id)
}

// GetSupplyRoutesByPatientID implementa la obtención de rutas por ID de paciente
func (s *supplyRouteService) GetSupplyRoutesByPatientID(ctx context.Context, patientID string) ([]*models.SupplyRoute, error) {
	return s.supplyRouteRepo.GetByPatientID(ctx, patientID)
}

// GetSupplyRoutesByOperatingRoomID implementa la obtención de rutas por ID de sala de operación
func (s *supplyRouteService) GetSupplyRoutesByOperatingRoomID(ctx context.Context, roomID string) ([]*models.SupplyRoute, error) {
	return s.supplyRouteRepo.GetByOperatingRoomID(ctx, roomID)
}

// GetAllSupplyRoutes implementa la obtención de todas las rutas de suministro
func (s *supplyRouteService) GetAllSupplyRoutes(ctx context.Context) ([]*models.SupplyRoute, error) {
	return s.supplyRouteRepo.GetAll(ctx)
}

// UpdateSupplyRoute implementa la actualización de una ruta de suministro
func (s *supplyRouteService) UpdateSupplyRoute(ctx context.Context, route *models.SupplyRoute) error {
	return s.supplyRouteRepo.Update(ctx, route)
}

// DeleteSupplyRoute implementa la eliminación de una ruta de suministro
func (s *supplyRouteService) DeleteSupplyRoute(ctx context.Context, id string) error {
	return s.supplyRouteRepo.Delete(ctx, id)
}
