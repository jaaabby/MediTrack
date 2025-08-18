package repository

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type MedicalSupplyRepository struct {
	DB *gorm.DB
}

func NewMedicalSupplyRepository(db *gorm.DB) *MedicalSupplyRepository {
	return &MedicalSupplyRepository{DB: db}
}

// CRUD básico
func (r *MedicalSupplyRepository) Create(supply *models.MedicalSupply) error {
	return r.DB.Create(supply).Error
}

func (r *MedicalSupplyRepository) GetByID(id int) (*models.MedicalSupply, error) {
	var supply models.MedicalSupply
	if err := r.DB.First(&supply, id).Error; err != nil {
		return nil, err
	}
	return &supply, nil
}

func (r *MedicalSupplyRepository) GetAll() ([]models.MedicalSupply, error) {
	var supplies []models.MedicalSupply
	if err := r.DB.Find(&supplies).Error; err != nil {
		return nil, err
	}
	return supplies, nil
}

func (r *MedicalSupplyRepository) Update(supply *models.MedicalSupply) error {
	return r.DB.Save(supply).Error
}

func (r *MedicalSupplyRepository) Delete(id int) error {
	return r.DB.Delete(&models.MedicalSupply{}, id).Error
}
