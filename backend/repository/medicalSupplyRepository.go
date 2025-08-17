package repository

import (
	"context"
	"meditrack/models"

	"gorm.io/gorm"
)

// MedicalSupplyRepositoryInterface define los métodos que debe implementar el repositorio de insumos médicos
type MedicalSupplyRepositoryInterface interface {
	Create(ctx context.Context, supply *models.MedicalSupply) error
	GetByID(ctx context.Context, id string) (*models.MedicalSupply, error)
	GetByQRCode(ctx context.Context, qrCode string) (*models.MedicalSupply, error)
	GetAll(ctx context.Context) ([]*models.MedicalSupply, error)
	Update(ctx context.Context, supply *models.MedicalSupply) error
	Delete(ctx context.Context, id string) error
}

type MedicalSupplyRepository struct {
	DB *gorm.DB
}

func NewMedicalSupplyRepository(db *gorm.DB) *MedicalSupplyRepository {
	return &MedicalSupplyRepository{DB: db}
}

// Create a new medical supply
func (r *MedicalSupplyRepository) Create(ctx context.Context, supply *models.MedicalSupply) error {
	return r.DB.WithContext(ctx).Create(supply).Error
}

// Get medical supply by ID
func (r *MedicalSupplyRepository) GetByID(ctx context.Context, id string) (*models.MedicalSupply, error) {
	var supply models.MedicalSupply
	if err := r.DB.WithContext(ctx).First(&supply, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &supply, nil
}

// Get medical supply by QRCode (code)
func (r *MedicalSupplyRepository) GetByQRCode(ctx context.Context, qrCode string) (*models.MedicalSupply, error) {
	var supply models.MedicalSupply
	if err := r.DB.WithContext(ctx).First(&supply, "code = ?", qrCode).Error; err != nil {
		return nil, err
	}
	return &supply, nil
}

// Get all medical supplies
func (r *MedicalSupplyRepository) GetAll(ctx context.Context) ([]*models.MedicalSupply, error) {
	var supplies []*models.MedicalSupply
	if err := r.DB.WithContext(ctx).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

// Update medical supply
func (r *MedicalSupplyRepository) Update(ctx context.Context, supply *models.MedicalSupply) error {
	return r.DB.WithContext(ctx).Save(supply).Error
}

// Delete medical supply
func (r *MedicalSupplyRepository) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&models.MedicalSupply{}, "id = ?", id).Error
}
