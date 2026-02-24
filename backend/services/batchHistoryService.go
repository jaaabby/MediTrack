package services

import (
	"encoding/json"
	"fmt"
	"meditrack/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type BatchHistoryService struct {
	DB *gorm.DB
}

func NewBatchHistoryService(db *gorm.DB) *BatchHistoryService {
	return &BatchHistoryService{DB: db}
}

// ========================
// MÉTODOS CRUD BÁSICOS
// ========================

func (s *BatchHistoryService) CreateBatchHistory(history *models.BatchHistory) error {
	return s.DB.Create(history).Error
}

func (s *BatchHistoryService) GetBatchHistoryByID(id int) (*models.BatchHistory, error) {
	var history models.BatchHistory
	if err := s.DB.First(&history, id).Error; err != nil {
		return nil, err
	}
	return &history, nil
}

func (s *BatchHistoryService) GetAllBatchHistories() ([]models.BatchHistory, error) {
	var histories []models.BatchHistory
	if err := s.DB.Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

func (s *BatchHistoryService) DeleteBatchHistory(id int) error {
	return s.DB.Delete(&models.BatchHistory{}, id).Error
}

func (s *BatchHistoryService) UpdateBatchHistory(id int, history *models.BatchHistory) error {
	var existing models.BatchHistory
	if err := s.DB.First(&existing, id).Error; err != nil {
		return err
	}
	return s.DB.Model(&existing).Omit("id", "created_at", "updated_at").Updates(history).Error
}

// ========================
// HELPER PRIVADO
// ========================

// formatHistoryEntry formatea una entrada de historial con valores JSON parseados
func (s *BatchHistoryService) formatHistoryEntry(history models.BatchHistory) map[string]interface{} {
	entry := map[string]interface{}{
		"id":             history.ID,
		"batch_id":       history.BatchNumber,
		"change_details": history.ChangeDetails,
		"date_time":      history.DateTime,
		"user_rut":       history.UserRUT,
		"user_name":      history.UserName,
	}

	if history.PreviousValues != "" && history.PreviousValues != "null" {
		var prevValues map[string]interface{}
		if err := json.Unmarshal([]byte(history.PreviousValues), &prevValues); err == nil {
			entry["previous_values"] = prevValues
		}
	}

	if history.NewValues != "" && history.NewValues != "{}" {
		var newValues map[string]interface{}
		if err := json.Unmarshal([]byte(history.NewValues), &newValues); err == nil {
			entry["new_values"] = newValues
		}
	}

	return entry
}

// ========================
// REGISTRO DE EVENTOS
// ========================

// getUserRUTOrDefault obtiene el RUT del usuario o retorna un valor por defecto
func (s *BatchHistoryService) getUserRUTOrDefault(userRUT string) string {
	if userRUT == "" {
		return "12345678-9"
	}
	return userRUT
}

func (s *BatchHistoryService) RegisterBatchCreation(batchID int, userRUT string) error {
	userRUT = s.getUserRUTOrDefault(userRUT)

	var batch models.Batch
	if err := s.DB.Preload("SupplierConfig").First(&batch, batchID).Error; err != nil {
		return fmt.Errorf("error obteniendo lote: %v", err)
	}

	var user models.User
	if err := s.DB.Where("rut = ?", userRUT).First(&user).Error; err != nil {
		return fmt.Errorf("error obteniendo información del usuario: %v", err)
	}

	newValues, _ := json.Marshal(map[string]interface{}{
		"expiration_date": batch.ExpirationDate,
		"amount":          batch.Amount,
		"supplier":        batch.Supplier,
		"store_id":        batch.StoreID,
	})

	history := &models.BatchHistory{
		DateTime:       time.Now(),
		BatchID:        &batchID,
		UserRUT:        userRUT,
		UserName:       user.Name,
		ChangeDetails:  "Lote creado",
		PreviousValues: "null",
		NewValues:      string(newValues),
		BatchNumber:    batchID,
	}

	return s.CreateBatchHistory(history)
}

func (s *BatchHistoryService) RegisterBatchUpdate(batchID int, userRUT string, previousBatch, newBatch *models.Batch) error {
	userRUT = s.getUserRUTOrDefault(userRUT)

	var user models.User
	if err := s.DB.Where("rut = ?", userRUT).First(&user).Error; err != nil {
		return fmt.Errorf("error obteniendo información del usuario: %v", err)
	}

	var changes []string
	previousValues := make(map[string]interface{})
	newValues := make(map[string]interface{})

	if !previousBatch.ExpirationDate.Equal(newBatch.ExpirationDate) {
		changes = append(changes, "fecha de expiración")
		previousValues["expiration_date"] = previousBatch.ExpirationDate
		newValues["expiration_date"] = newBatch.ExpirationDate
	}
	if previousBatch.Amount != newBatch.Amount {
		changes = append(changes, "cantidad")
		previousValues["amount"] = previousBatch.Amount
		newValues["amount"] = newBatch.Amount
	}
	if previousBatch.SupplierID != newBatch.SupplierID {
		changes = append(changes, "proveedor")
		previousValues["supplier_id"] = previousBatch.SupplierID
		newValues["supplier_id"] = newBatch.SupplierID
	}
	if previousBatch.StoreID != newBatch.StoreID {
		changes = append(changes, "almacén")
		previousValues["store_id"] = previousBatch.StoreID
		newValues["store_id"] = newBatch.StoreID
	}

	previousValuesJSON, _ := json.Marshal(previousValues)
	newValuesJSON, _ := json.Marshal(newValues)

	history := &models.BatchHistory{
		DateTime:       time.Now(),
		BatchID:        &batchID,
		UserRUT:        userRUT,
		UserName:       user.Name,
		ChangeDetails:  "Lote actualizado: " + strings.Join(changes, ", "),
		PreviousValues: string(previousValuesJSON),
		NewValues:      string(newValuesJSON),
		BatchNumber:    batchID,
	}

	return s.CreateBatchHistory(history)
}

func (s *BatchHistoryService) RegisterBatchDeletion(batchID int, userRUT string, deletedBatch *models.Batch) error {
	userRUT = s.getUserRUTOrDefault(userRUT)

	var user models.User
	if err := s.DB.Where("rut = ?", userRUT).First(&user).Error; err != nil {
		return fmt.Errorf("error obteniendo información del usuario: %v", err)
	}

	previousValues, _ := json.Marshal(map[string]interface{}{
		"expiration_date": deletedBatch.ExpirationDate,
		"amount":          deletedBatch.Amount,
		"supplier":        deletedBatch.Supplier,
		"store_id":        deletedBatch.StoreID,
	})

	history := &models.BatchHistory{
		DateTime:       time.Now(),
		BatchID:        &batchID,
		UserRUT:        userRUT,
		UserName:       user.Name,
		ChangeDetails:  "Lote eliminado",
		PreviousValues: string(previousValues),
		NewValues:      "{}",
		BatchNumber:    batchID,
	}

	return s.CreateBatchHistory(history)
}

// ========================
// CONSULTAS CON FORMATO
// ========================

func (s *BatchHistoryService) GetAllBatchHistoriesWithDetails() ([]map[string]interface{}, error) {
	var histories []models.BatchHistory
	if err := s.DB.Order("date_time DESC").Find(&histories).Error; err != nil {
		return nil, err
	}

	detailedHistories := make([]map[string]interface{}, 0, len(histories))
	for _, history := range histories {
		detailedHistories = append(detailedHistories, s.formatHistoryEntry(history))
	}

	return detailedHistories, nil
}

func (s *BatchHistoryService) GetAllBatchHistoriesWithDetailsPaginated(page, pageSize int) ([]map[string]interface{}, int64, error) {
	var total int64
	if err := s.DB.Model(&models.BatchHistory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var histories []models.BatchHistory
	if err := s.DB.Order("date_time DESC").Offset(offset).Limit(pageSize).Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	detailedHistories := make([]map[string]interface{}, 0, len(histories))
	for _, history := range histories {
		detailedHistories = append(detailedHistories, s.formatHistoryEntry(history))
	}

	return detailedHistories, total, nil
}

func (s *BatchHistoryService) SearchBatchHistoryByBatchNumber(batchNumber int) ([]map[string]interface{}, error) {
	var histories []models.BatchHistory
	if err := s.DB.Where("batch_number = ?", batchNumber).Order("date_time DESC").Find(&histories).Error; err != nil {
		return nil, err
	}

	detailedHistories := make([]map[string]interface{}, 0, len(histories))
	for _, history := range histories {
		detailedHistories = append(detailedHistories, s.formatHistoryEntry(history))
	}

	return detailedHistories, nil
}

// ELIMINADO: GetBatchHistoryByBatchNumber - duplicado de SearchBatchHistoryByBatchNumber
