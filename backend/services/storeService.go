package services

import (
	"meditrack/models"
	"meditrack/repository"
)

type StoreService interface {
	CreateStore(store *models.Store) error
	GetStoreByID(id int) (*models.Store, error)
	GetAllStores() ([]models.Store, error)
	UpdateStore(store *models.Store) error
	DeleteStore(id int) error
}

type storeService struct {
	repo *repository.StoreRepository
}

func NewStoreService(repo *repository.StoreRepository) StoreService {
	return &storeService{repo: repo}
}

func (s *storeService) CreateStore(store *models.Store) error {
	return s.repo.Create(store)
}

func (s *storeService) GetStoreByID(id int) (*models.Store, error) {
	return s.repo.GetByID(id)
}

func (s *storeService) GetAllStores() ([]models.Store, error) {
	return s.repo.GetAll()
}

func (s *storeService) UpdateStore(store *models.Store) error {
	return s.repo.Update(store)
}

func (s *storeService) DeleteStore(id int) error {
	return s.repo.Delete(id)
}
