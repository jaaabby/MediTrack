package repository

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type StoreRepository struct {
	DB *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{DB: db}
}

// Create a new store
func (r *StoreRepository) Create(store *models.Store) error {
	return r.DB.Create(store).Error
}

// Get store by ID
func (r *StoreRepository) GetByID(id int) (*models.Store, error) {
	var store models.Store
	if err := r.DB.First(&store, id).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

// Get all stores
func (r *StoreRepository) GetAll() ([]models.Store, error) {
	var stores []models.Store
	if err := r.DB.Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

// Update store
func (r *StoreRepository) Update(store *models.Store) error {
	return r.DB.Save(store).Error
}

// Delete store
func (r *StoreRepository) Delete(id int) error {
	return r.DB.Delete(&models.Store{}, id).Error
}
