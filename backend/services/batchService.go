package services

import (
	"meditrack/models"

	"gorm.io/gorm"
)

type BatchService struct {
	DB        *gorm.DB
	QRService *QRService
}

func NewBatchService(db *gorm.DB, qrService *QRService) *BatchService {
	return &BatchService{
		DB:        db,
		QRService: qrService,
	}
}

func (s *BatchService) CreateBatch(batch *models.Batch) error {
	// Generar código QR único para el batch
	qrCode, err := s.QRService.GenerateBatchQRCode()
	if err != nil {
		return err
	}
	batch.QRCode = qrCode

	return s.DB.Create(batch).Error
}

func (s *BatchService) GetBatchByID(id int) (*models.Batch, error) {
	var batch models.Batch
	if err := s.DB.First(&batch, id).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

func (s *BatchService) GetBatchByQRCode(qrCode string) (*models.Batch, error) {
	var batch models.Batch
	if err := s.DB.Where("qr_code = ?", qrCode).First(&batch).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

func (s *BatchService) GetAllBatches() ([]models.Batch, error) {
	var batches []models.Batch
	if err := s.DB.Find(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

func (s *BatchService) UpdateBatch(id int, newBatch *models.Batch) (*models.Batch, error) {
	var batch models.Batch
	if err := s.DB.First(&batch, id).Error; err != nil {
		return nil, err
	}

	// Actualizar campos pero mantener el QR code original
	batch.ExpirationDate = newBatch.ExpirationDate
	batch.Amount = newBatch.Amount
	batch.Supplier = newBatch.Supplier
	batch.StoreID = newBatch.StoreID
	// No actualizamos batch.QRCode para mantener la trazabilidad

	if err := s.DB.Save(&batch).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

func (s *BatchService) DeleteBatch(id int) error {
	if err := s.DB.Delete(&models.Batch{}, id).Error; err != nil {
		return err
	}
	return nil
}
