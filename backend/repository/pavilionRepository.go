package repository

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type PavilionRepository struct {
	DB *gorm.DB
}

func NewPavilionRepository(db *gorm.DB) *PavilionRepository {
	return &PavilionRepository{DB: db}
}

// Create a new pavilion
func (r *PavilionRepository) Create(pavilion *models.Pavilion) error {
	return r.DB.Create(pavilion).Error
}

// Get pavilion by ID
func (r *PavilionRepository) GetByID(id int) (*models.Pavilion, error) {
	var pavilion models.Pavilion
	if err := r.DB.First(&pavilion, id).Error; err != nil {
		return nil, err
	}
	return &pavilion, nil
}

// Get all pavilions
func (r *PavilionRepository) GetAll() ([]models.Pavilion, error) {
	var pavilions []models.Pavilion
	if err := r.DB.Find(&pavilions).Error; err != nil {
		return nil, err
	}
	return pavilions, nil
}

// Update pavilion
func (r *PavilionRepository) Update(pavilion *models.Pavilion) error {
	return r.DB.Save(pavilion).Error
}

// Delete pavilion
func (r *PavilionRepository) Delete(id int) error {
	return r.DB.Delete(&models.Pavilion{}, id).Error
}
