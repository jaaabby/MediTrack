package services

import (
	"meditrack/models"
	"meditrack/repository"
)

type BatchService interface {
	CreateBatch(batch *models.Batch) error
	GetBatchByID(id int) (*models.Batch, error)
	GetAllBatches() ([]models.Batch, error)
	UpdateBatch(batch *models.Batch) error
	DeleteBatch(id int) error
}

type batchService struct {
	repo *repository.BatchRepository
}

func NewBatchService(repo *repository.BatchRepository) BatchService {
	return &batchService{repo: repo}
}

func (s *batchService) CreateBatch(batch *models.Batch) error {
	return s.repo.Create(batch)
}

func (s *batchService) GetBatchByID(id int) (*models.Batch, error) {
	return s.repo.GetByID(id)
}

func (s *batchService) GetAllBatches() ([]models.Batch, error) {
	return s.repo.GetAll()
}

func (s *batchService) UpdateBatch(batch *models.Batch) error {
	return s.repo.Update(batch)
}

func (s *batchService) DeleteBatch(id int) error {
	return s.repo.Delete(id)
}
