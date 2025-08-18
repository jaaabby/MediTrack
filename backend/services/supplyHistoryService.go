package services

import (
	"meditrack/models"
	"meditrack/repository"
)

type SupplyHistoryService interface {
	CreateSupplyHistory(history *models.SupplyHistory) error
	GetSupplyHistoryByID(id int) (*models.SupplyHistory, error)
	GetAllSupplyHistories() ([]models.SupplyHistory, error)
	UpdateSupplyHistory(history *models.SupplyHistory) error
	DeleteSupplyHistory(id int) error
}

type supplyHistoryService struct {
	repo *repository.SupplyHistoryRepository
}

func NewSupplyHistoryService(repo *repository.SupplyHistoryRepository) SupplyHistoryService {
	return &supplyHistoryService{repo: repo}
}

func (s *supplyHistoryService) CreateSupplyHistory(history *models.SupplyHistory) error {
	return s.repo.Create(history)
}

func (s *supplyHistoryService) GetSupplyHistoryByID(id int) (*models.SupplyHistory, error) {
	return s.repo.GetByID(id)
}

func (s *supplyHistoryService) GetAllSupplyHistories() ([]models.SupplyHistory, error) {
	return s.repo.GetAll()
}

func (s *supplyHistoryService) UpdateSupplyHistory(history *models.SupplyHistory) error {
	return s.repo.Update(history)
}

func (s *supplyHistoryService) DeleteSupplyHistory(id int) error {
	return s.repo.Delete(id)
}
