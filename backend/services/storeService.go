package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type StoreService struct {
	DB *gorm.DB
}

func NewStoreService(db *gorm.DB) *StoreService {
	return &StoreService{DB: db}
}

func (s *StoreService) CreateStore(store *models.Store) error {
	return s.DB.Create(store).Error
}

func (s *StoreService) GetStoreByID(id int) (*models.Store, error) {
	var store models.Store
	if err := s.DB.First(&store, id).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (s *StoreService) GetAllStores() ([]models.Store, error) {
	var stores []models.Store
	if err := s.DB.Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

func (s *StoreService) UpdateStore(id int, newStore *models.Store) (*models.Store, error) {
	var store models.Store
	if err := s.DB.First(&store, id).Error; err != nil {
		return nil, err
	}
	store.Name = newStore.Name
	store.Type = newStore.Type
	store.MedicalCenterID = newStore.MedicalCenterID

	if err := s.DB.Save(&store).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (s *StoreService) DeleteStore(id int) error {
	if err := s.DB.Delete(&models.Store{}, id).Error; err != nil {
		return err
	}
	return nil
}
