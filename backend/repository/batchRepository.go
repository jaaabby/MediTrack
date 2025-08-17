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

// CRUD básico
func (r *BatchRepository) Create(batch *models.Batch) error {
	return r.DB.Create(batch).Error
}

func (r *BatchRepository) GetByID(id int) (*models.Batch, error) {
	var batch models.Batch
	if err := r.DB.First(&batch, id).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *BatchRepository) GetAll() ([]models.Batch, error) {
	var batches []models.Batch
	if err := r.DB.Find(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

func (r *BatchRepository) Update(batch *models.Batch) error {
	return r.DB.Save(batch).Error
}

func (r *BatchRepository) Delete(id int) error {
	return r.DB.Delete(&models.Batch{}, id).Error
}
