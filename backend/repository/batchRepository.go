package repository

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type BatchRepository struct {
	DB *gorm.DB
}

func NewBatchRepository(db *gorm.DB) *BatchRepository {
	return &BatchRepository{DB: db}
}

// Create a new batch
func (r *BatchRepository) Create(batch *models.Batch) error {
	return r.DB.Create(batch).Error
}

// Get batch by ID
func (r *BatchRepository) GetByID(id int) (*models.Batch, error) {
	var batch models.Batch
	if err := r.DB.First(&batch, id).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

// Get all batches
func (r *BatchRepository) GetAll() ([]models.Batch, error) {
	var batches []models.Batch
	if err := r.DB.Find(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

// Update batch
func (r *BatchRepository) Update(batch *models.Batch) error {
	return r.DB.Save(batch).Error
}

// Delete batch
func (r *BatchRepository) Delete(id int) error {
	return r.DB.Delete(&models.Batch{}, id).Error
}
