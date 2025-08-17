package services

import (
	"meditrack/models"
	"meditrack/repository"
)

type MedicalSupplyService interface {
	CreateMedicalSupply(supply *models.MedicalSupply) error
	GetMedicalSupplyByID(id int) (*models.MedicalSupply, error)
	GetAllMedicalSupplies() ([]models.MedicalSupply, error)
	UpdateMedicalSupply(supply *models.MedicalSupply) error
	DeleteMedicalSupply(id int) error
}

type medicalSupplyService struct {
	repo *repository.MedicalSupplyRepository
}

func NewMedicalSupplyService(repo *repository.MedicalSupplyRepository) MedicalSupplyService {
	return &medicalSupplyService{
		repo: repo,
	}
}

func (s *medicalSupplyService) CreateMedicalSupply(supply *models.MedicalSupply) error {
	return s.repo.Create(supply)
}

func (s *medicalSupplyService) GetMedicalSupplyByID(id int) (*models.MedicalSupply, error) {
	return s.repo.GetByID(id)
}

func (s *medicalSupplyService) GetAllMedicalSupplies() ([]models.MedicalSupply, error) {
	return s.repo.GetAll()
}

func (s *medicalSupplyService) UpdateMedicalSupply(supply *models.MedicalSupply) error {
	return s.repo.Update(supply)
}

func (s *medicalSupplyService) DeleteMedicalSupply(id int) error {
	return s.repo.Delete(id)
}
