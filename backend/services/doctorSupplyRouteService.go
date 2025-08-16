package services

import (
	"context"

	"meditrack/models"
	"meditrack/repository"
)

// DoctorSupplyRouteService define los servicios para la relación doctor-ruta de suministro
type DoctorSupplyRouteService interface {
	CreateDoctorSupplyRoute(ctx context.Context, relation *models.DoctorSupplyRoute) error
	GetDoctorSupplyRouteByID(ctx context.Context, id string) (*models.DoctorSupplyRoute, error)
	GetDoctorSupplyRoutesByDoctorID(ctx context.Context, doctorID string) ([]*models.DoctorSupplyRoute, error)
	GetDoctorSupplyRoutesBySupplyRouteID(ctx context.Context, routeID string) ([]*models.DoctorSupplyRoute, error)
	GetAllDoctorSupplyRoutes(ctx context.Context) ([]*models.DoctorSupplyRoute, error)
	UpdateDoctorSupplyRoute(ctx context.Context, relation *models.DoctorSupplyRoute) error
	DeleteDoctorSupplyRoute(ctx context.Context, id string) error
}

// doctorSupplyRouteService implementa DoctorSupplyRouteService
type doctorSupplyRouteService struct {
	doctorSupplyRouteRepo repository.DoctorSupplyRouteRepository
}

// NewDoctorSupplyRouteService crea una nueva instancia de DoctorSupplyRouteService
func NewDoctorSupplyRouteService(doctorSupplyRouteRepo repository.DoctorSupplyRouteRepository) DoctorSupplyRouteService {
	return &doctorSupplyRouteService{
		doctorSupplyRouteRepo: doctorSupplyRouteRepo,
	}
}

// CreateDoctorSupplyRoute implementa la creación de una relación doctor-ruta
func (s *doctorSupplyRouteService) CreateDoctorSupplyRoute(ctx context.Context, relation *models.DoctorSupplyRoute) error {
	return s.doctorSupplyRouteRepo.Create(ctx, relation)
}

// GetDoctorSupplyRouteByID implementa la obtención de una relación doctor-ruta por ID
func (s *doctorSupplyRouteService) GetDoctorSupplyRouteByID(ctx context.Context, id string) (*models.DoctorSupplyRoute, error) {
	return s.doctorSupplyRouteRepo.GetByID(ctx, id)
}

// GetDoctorSupplyRoutesByDoctorID implementa la obtención de relaciones por ID de doctor
func (s *doctorSupplyRouteService) GetDoctorSupplyRoutesByDoctorID(ctx context.Context, doctorID string) ([]*models.DoctorSupplyRoute, error) {
	return s.doctorSupplyRouteRepo.GetByDoctorID(ctx, doctorID)
}

// GetDoctorSupplyRoutesBySupplyRouteID implementa la obtención de relaciones por ID de ruta de suministro
func (s *doctorSupplyRouteService) GetDoctorSupplyRoutesBySupplyRouteID(ctx context.Context, routeID string) ([]*models.DoctorSupplyRoute, error) {
	return s.doctorSupplyRouteRepo.GetBySupplyRouteID(ctx, routeID)
}

// GetAllDoctorSupplyRoutes implementa la obtención de todas las relaciones doctor-ruta
func (s *doctorSupplyRouteService) GetAllDoctorSupplyRoutes(ctx context.Context) ([]*models.DoctorSupplyRoute, error) {
	return s.doctorSupplyRouteRepo.GetAll(ctx)
}

// UpdateDoctorSupplyRoute implementa la actualización de una relación doctor-ruta
func (s *doctorSupplyRouteService) UpdateDoctorSupplyRoute(ctx context.Context, relation *models.DoctorSupplyRoute) error {
	return s.doctorSupplyRouteRepo.Update(ctx, relation)
}

// DeleteDoctorSupplyRoute implementa la eliminación de una relación doctor-ruta
func (s *doctorSupplyRouteService) DeleteDoctorSupplyRoute(ctx context.Context, id string) error {
	return s.doctorSupplyRouteRepo.Delete(ctx, id)
}
