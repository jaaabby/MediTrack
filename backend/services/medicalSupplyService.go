package services

import (
	"context"

	"meditrack/models"
	"meditrack/repository"
)

// MedicalSupplyService define los servicios para insumos médicos
type MedicalSupplyService interface {
	CreateMedicalSupply(ctx context.Context, supply *models.MedicalSupply) error
	GetMedicalSupplyByID(ctx context.Context, id string) (*models.MedicalSupply, error)
	GetMedicalSupplyByQRCode(ctx context.Context, qrCode string) (*models.MedicalSupply, error)
	GetAllMedicalSupplies(ctx context.Context) ([]*models.MedicalSupply, error)
	UpdateMedicalSupply(ctx context.Context, supply *models.MedicalSupply) error
	DeleteMedicalSupply(ctx context.Context, id string) error
}

// medicalSupplyService implementa MedicalSupplyService
type medicalSupplyService struct {
	medicalSupplyRepo repository.MedicalSupplyRepositoryInterface
}

// NewMedicalSupplyService crea una nueva instancia de MedicalSupplyService
func NewMedicalSupplyService(medicalSupplyRepo repository.MedicalSupplyRepositoryInterface) MedicalSupplyService {
	return &medicalSupplyService{
		medicalSupplyRepo: medicalSupplyRepo,
	}
}

// CreateMedicalSupply implementa la creación de un insumo médico
func (s *medicalSupplyService) CreateMedicalSupply(ctx context.Context, supply *models.MedicalSupply) error {
	return s.medicalSupplyRepo.Create(ctx, supply)
}

// GetMedicalSupplyByID implementa la obtención de un insumo médico por ID
func (s *medicalSupplyService) GetMedicalSupplyByID(ctx context.Context, id string) (*models.MedicalSupply, error) {
	return s.medicalSupplyRepo.GetByID(ctx, id)
}

// GetMedicalSupplyByQRCode implementa la obtención de un insumo médico por código QR
func (s *medicalSupplyService) GetMedicalSupplyByQRCode(ctx context.Context, qrCode string) (*models.MedicalSupply, error) {
	return s.medicalSupplyRepo.GetByQRCode(ctx, qrCode)
}

// GetAllMedicalSupplies implementa la obtención de todos los insumos médicos
func (s *medicalSupplyService) GetAllMedicalSupplies(ctx context.Context) ([]*models.MedicalSupply, error) {
	return s.medicalSupplyRepo.GetAll(ctx)
}

// UpdateMedicalSupply implementa la actualización de un insumo médico
func (s *medicalSupplyService) UpdateMedicalSupply(ctx context.Context, supply *models.MedicalSupply) error {
	return s.medicalSupplyRepo.Update(ctx, supply)
}

// DeleteMedicalSupply implementa la eliminación de un insumo médico
func (s *medicalSupplyService) DeleteMedicalSupply(ctx context.Context, id string) error {
	return s.medicalSupplyRepo.Delete(ctx, id)
}
