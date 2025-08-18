package services

import (
	"meditrack/models"
	"meditrack/repository"
)

type MedicalCenterService interface {
	CreateMedicalCenter(center *models.MedicalCenter) error
	GetMedicalCenterByID(id int) (*models.MedicalCenter, error)
	GetAllMedicalCenters() ([]models.MedicalCenter, error)
	UpdateMedicalCenter(center *models.MedicalCenter) error
	DeleteMedicalCenter(id int) error
}

type medicalCenterService struct {
	repo *repository.MedicalCenterRepository
}

func NewMedicalCenterService(repo *repository.MedicalCenterRepository) MedicalCenterService {
	return &medicalCenterService{repo: repo}
}

func (s *medicalCenterService) CreateMedicalCenter(center *models.MedicalCenter) error {
	return s.repo.Create(center)
}

func (s *medicalCenterService) GetMedicalCenterByID(id int) (*models.MedicalCenter, error) {
	return s.repo.GetByID(id)
}

func (s *medicalCenterService) GetAllMedicalCenters() ([]models.MedicalCenter, error) {
	return s.repo.GetAll()
}

func (s *medicalCenterService) UpdateMedicalCenter(center *models.MedicalCenter) error {
	return s.repo.Update(center)
}

func (s *medicalCenterService) DeleteMedicalCenter(id int) error {
	return s.repo.Delete(id)
}
