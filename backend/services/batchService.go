package services

import (
	"fmt"
	"meditrack/mailer"
	"meditrack/models"
	"time"

	"gorm.io/gorm"
)

type BatchService struct {
	DB                  *gorm.DB
	BatchHistoryService *BatchHistoryService
}

func NewBatchService(db *gorm.DB) *BatchService {
	return &BatchService{
		DB:                  db,
		BatchHistoryService: NewBatchHistoryService(db),
	}
}

func (s *BatchService) CreateBatch(batch *models.Batch) error {
	if err := s.DB.Create(batch).Error; err != nil {
		return err
	}

	// Registrar en el historial (RUT hardcodeado por ahora)
	userRUT := "12345678-9" // RUT hardcodeado temporalmente
	if err := s.BatchHistoryService.RegisterBatchCreation(batch.ID, userRUT); err != nil {
		// Solo log del error, no fallar la creación
		fmt.Printf("Error al registrar en historial: %v\n", err)
	}

	return nil
}

func (s *BatchService) GetBatchByID(id int) (*models.Batch, error) {
	var batch models.Batch
	if err := s.DB.First(&batch, id).Error; err != nil {
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

	// Guardar valores anteriores para el historial
	previousBatch := batch

	batch.ExpirationDate = newBatch.ExpirationDate
	batch.Amount = newBatch.Amount
	batch.Supplier = newBatch.Supplier
	batch.StoreID = newBatch.StoreID

	if err := s.DB.Save(&batch).Error; err != nil {
		return nil, err
	}

	// Registrar en el historial (RUT hardcodeado por ahora)
	userRUT := "12345678-9" // RUT hardcodeado temporalmente
	if err := s.BatchHistoryService.RegisterBatchUpdate(batch.ID, userRUT, &previousBatch, &batch); err != nil {
		// Solo log del error, no fallar la actualización
		fmt.Printf("Error al registrar en historial: %v\n", err)
	}

	// Verificar si se debe enviar correo por stock bajo
	if batch.Amount < 10 {
		if err := s.sendLowStockAlert(batch); err != nil {
			// Solo log del error, no fallar la actualización
			fmt.Printf("Error al enviar alerta de stock bajo: %v\n", err)
		}
	}

	// Verificar si se debe enviar correo por vencimiento próximo (30 días)
	expirationThreshold := time.Now().AddDate(0, 0, 30)
	if batch.ExpirationDate.Before(expirationThreshold) && batch.ExpirationDate.After(time.Now()) {
		if err := s.sendExpirationAlert(batch); err != nil {
			// Solo log del error, no fallar la actualización
			fmt.Printf("Error al enviar alerta de vencimiento: %v\n", err)
		}
	}

	return &batch, nil
}

func (s *BatchService) DeleteBatch(id int) error {
	// Usar transacción para asegurar consistencia
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Obtener el lote antes de eliminarlo para el historial
		var batch models.Batch
		if err := tx.First(&batch, id).Error; err != nil {
			return err
		}

		// NOTA: El historial se registra automáticamente por el trigger log_batch_delete
		// No es necesario llamar manualmente a RegisterBatchDeletion

		// 1. Eliminar registros de medical_supply (supply_history se elimina automáticamente por CASCADE)
		if err := tx.Where("batch_id = ?", id).Delete(&models.MedicalSupply{}).Error; err != nil {
			return fmt.Errorf("error al eliminar insumos médicos del lote: %v", err)
		}

		// 2. Eliminar el lote (el trigger log_batch_delete se ejecutará automáticamente)
		if err := tx.Delete(&models.Batch{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// sendLowStockAlert envía correo de alerta de stock bajo
func (s *BatchService) sendLowStockAlert(batch models.Batch) error {
	// Obtener información completa del insumo desde supply_code a través de medical_supply
	var supplyCode models.SupplyCode
	if err := s.DB.Joins("JOIN medical_supply ON medical_supply.code = supply_code.code").
		Where("medical_supply.batch_id = ?", batch.ID).
		First(&supplyCode).Error; err != nil {
		return fmt.Errorf("error al obtener información del código de insumo: %v", err)
	}

	// Crear datos para la plantilla
	data := map[string]interface{}{
		"BatchID":      batch.ID,
		"Code":         supplyCode.Code,
		"Name":         supplyCode.Name,
		"CurrentStock": batch.Amount,
		"Date":         time.Now().Format("02/01/2006"),
	}

	// Crear solicitud de correo
	req := mailer.NewRequest([]string{"vergara.javiera12@gmail.com"}, "Alerta: Stock Bajo en Lote")

	// Enviar correo usando la plantilla de stock bajo
	templatePath := "mailer/templates/low_stock.html"
	return req.SendMailSkipTLS(templatePath, data)
}

// sendExpirationAlert envía correo de alerta de vencimiento próximo
func (s *BatchService) sendExpirationAlert(batch models.Batch) error {
	// Obtener información completa del insumo desde supply_code a través de medical_supply
	var supplyCode models.SupplyCode
	if err := s.DB.Joins("JOIN medical_supply ON medical_supply.code = supply_code.code").
		Where("medical_supply.batch_id = ?", batch.ID).
		First(&supplyCode).Error; err != nil {
		return fmt.Errorf("error al obtener información del código de vencimiento: %v", err)
	}

	// Crear datos para la plantilla
	data := map[string]interface{}{
		"BatchID":        batch.ID,
		"Code":           supplyCode.Code,
		"Name":           supplyCode.Name,
		"ExpirationDate": batch.ExpirationDate.Format("02/01/2006"),
		"DaysUntilExp":   int(time.Until(batch.ExpirationDate).Hours() / 24),
		"Date":           time.Now().Format("02/01/2006"),
	}

	// Crear solicitud de correo
	req := mailer.NewRequest([]string{"vergara.javiera12@gmail.com"}, "Alerta: Lote Próximo a Vencer")

	// Enviar correo usando la plantilla de vencimiento
	templatePath := "mailer/templates/expiration_warning.html"
	return req.SendMailSkipTLS(templatePath, data)
}
