package repository

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type MedicalCenterRepository struct {
	DB *gorm.DB
}

func NewMedicalCenterRepository(db *gorm.DB) *MedicalCenterRepository {
	return &MedicalCenterRepository{DB: db}
}

// CRUD básico
func (r *MedicalCenterRepository) Create(center *models.MedicalCenter) error {
	return r.DB.Create(center).Error
}

func (r *MedicalCenterRepository) GetByID(id int) (*models.MedicalCenter, error) {
	var center models.MedicalCenter
	if err := r.DB.First(&center, id).Error; err != nil {
		return nil, err
	}
	return &center, nil
}

func (r *MedicalCenterRepository) GetAll() ([]models.MedicalCenter, error) {
	var centers []models.MedicalCenter
	if err := r.DB.Find(&centers).Error; err != nil {
		return nil, err
	}
	return centers, nil
}

func (r *MedicalCenterRepository) Update(center *models.MedicalCenter) error {
	return r.DB.Save(center).Error
}

func (r *MedicalCenterRepository) Delete(id int) error {
	return r.DB.Delete(&models.MedicalCenter{}, id).Error
}
