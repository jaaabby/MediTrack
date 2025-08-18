package repository

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type SupplyHistoryRepository struct {
	DB *gorm.DB
}

func NewSupplyHistoryRepository(db *gorm.DB) *SupplyHistoryRepository {
	return &SupplyHistoryRepository{DB: db}
}

// Create a new supply history
func (r *SupplyHistoryRepository) Create(history *models.SupplyHistory) error {
	return r.DB.Create(history).Error
}

// Get supply history by ID
func (r *SupplyHistoryRepository) GetByID(id int) (*models.SupplyHistory, error) {
	var history models.SupplyHistory
	if err := r.DB.First(&history, id).Error; err != nil {
		return nil, err
	}
	return &history, nil
}

// Get all supply histories
func (r *SupplyHistoryRepository) GetAll() ([]models.SupplyHistory, error) {
	var histories []models.SupplyHistory
	if err := r.DB.Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

// Update supply history
func (r *SupplyHistoryRepository) Update(history *models.SupplyHistory) error {
	return r.DB.Save(history).Error
}

// Delete supply history
func (r *SupplyHistoryRepository) Delete(id int) error {
	return r.DB.Delete(&models.SupplyHistory{}, id).Error
}
