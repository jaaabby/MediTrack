package services

import (
	"context"

	"meditrack/models"
	"meditrack/repository"
)

// LocationService define los servicios para ubicaciones
type LocationService interface {
	CreateLocation(ctx context.Context, location *models.Location) error
	GetLocationByID(ctx context.Context, id string) (*models.Location, error)
	GetAllLocations(ctx context.Context) ([]*models.Location, error)
	UpdateLocation(ctx context.Context, location *models.Location) error
	DeleteLocation(ctx context.Context, id string) error
}

// locationService implementa LocationService
type locationService struct {
	locationRepo repository.LocationRepository
}

// NewLocationService crea una nueva instancia de LocationService
func NewLocationService(locationRepo repository.LocationRepository) LocationService {
	return &locationService{
		locationRepo: locationRepo,
	}
}

// CreateLocation implementa la creación de una ubicación
func (s *locationService) CreateLocation(ctx context.Context, location *models.Location) error {
	return s.locationRepo.Create(ctx, location)
}

// GetLocationByID implementa la obtención de una ubicación por ID
func (s *locationService) GetLocationByID(ctx context.Context, id string) (*models.Location, error) {
	return s.locationRepo.GetByID(ctx, id)
}

// GetAllLocations implementa la obtención de todas las ubicaciones
func (s *locationService) GetAllLocations(ctx context.Context) ([]*models.Location, error) {
	return s.locationRepo.GetAll(ctx)
}

// UpdateLocation implementa la actualización de una ubicación
func (s *locationService) UpdateLocation(ctx context.Context, location *models.Location) error {
	return s.locationRepo.Update(ctx, location)
}

// DeleteLocation implementa la eliminación de una ubicación
func (s *locationService) DeleteLocation(ctx context.Context, id string) error {
	return s.locationRepo.Delete(ctx, id)
}
