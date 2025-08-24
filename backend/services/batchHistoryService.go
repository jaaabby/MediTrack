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
	if err := s.DB.Delete(&models.BatchHistory{}, id).Error; err != nil {
		return err
	}
	return nil
}

// UpdateBatchHistory actualiza un BatchHistory por ID
func (s *BatchHistoryService) UpdateBatchHistory(id int, history *models.BatchHistory) error {
	var existing models.BatchHistory
	if err := s.DB.First(&existing, id).Error; err != nil {
		return err
	}
	// Actualiza los campos necesarios
	if err := s.DB.Model(&existing).Updates(history).Error; err != nil {
		return err
	}
	return nil
}

// RegisterBatchCreation registra la creación de un nuevo lote
func (s *BatchHistoryService) RegisterBatchCreation(batchID int, userRUT string) error {
	// RUT hardcodeado para el usuario
	if userRUT == "" {
		userRUT = "12345678-9"
	}

	// Obtener información del lote y código de insumo
	var batch models.Batch
	var supplyCode models.SupplyCode

	if err := s.DB.First(&batch, batchID).Error; err != nil {
		return fmt.Errorf("error obteniendo lote: %v", err)
	}

	// Obtener información del código de insumo desde medical_supply
	var medicalSupply models.MedicalSupply
	if err := s.DB.Where("batch_id = ?", batchID).First(&medicalSupply).Error; err != nil {
		return fmt.Errorf("error obteniendo información del insumo médico: %v", err)
	}

	if err := s.DB.Where("code = ?", medicalSupply.Code).First(&supplyCode).Error; err != nil {
		return fmt.Errorf("error obteniendo código de insumo: %v", err)
	}

	// Obtener información del usuario
	var user models.User
	if err := s.DB.Where("rut = ?", userRUT).First(&user).Error; err != nil {
		return fmt.Errorf("error obteniendo información del usuario: %v", err)
	}

	// Crear valores nuevos en formato JSON
	newValues, _ := json.Marshal(map[string]interface{}{
		"expiration_date": batch.ExpirationDate,
		"amount":          batch.Amount,
		"supplier":        batch.Supplier,
		"store_id":        batch.StoreID,
	})

	history := &models.BatchHistory{
		DateTime:      time.Now(),
		BatchID:       &batchID,
		UserRUT:       userRUT,
		UserName:      user.Name, // Agregar el nombre del usuario
		ChangeDetails: "Lote creado",
		NewValues:     string(newValues),
		BatchNumber:   batchID, // Agregar batch_number
	}

	return s.CreateBatchHistory(history)
}

// RegisterBatchUpdate registra la actualización de un lote existente
func (s *BatchHistoryService) RegisterBatchUpdate(batchID int, userRUT string, previousBatch, newBatch *models.Batch) error {
	// RUT hardcodeado para el usuario
	if userRUT == "" {
		userRUT = "12345678-9"
	}
	// Obtener información del código de insumo desde medical_supply
	var medicalSupply models.MedicalSupply
	if err := s.DB.Where("batch_id = ?", batchID).First(&medicalSupply).Error; err != nil {
		return fmt.Errorf("error obteniendo información del insumo médico: %v", err)
	}

	// Obtener información del código de insumo
	var supplyCode models.SupplyCode
	if err := s.DB.Where("code = ?", medicalSupply.Code).First(&supplyCode).Error; err != nil {
		return fmt.Errorf("error obteniendo código de insumo: %v", err)
	}

	// Obtener información del usuario
	var user models.User
	if err := s.DB.Where("rut = ?", userRUT).First(&user).Error; err != nil {
		return fmt.Errorf("error obteniendo información del usuario: %v", err)
	}

	// Determinar qué campos cambiaron y crear valores solo para los campos modificados
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
	if previousBatch.Supplier != newBatch.Supplier {
		changes = append(changes, "proveedor")
		previousValues["supplier"] = previousBatch.Supplier
		newValues["supplier"] = newBatch.Supplier
	}
	if previousBatch.StoreID != newBatch.StoreID {
		changes = append(changes, "almacén")
		previousValues["store_id"] = newBatch.StoreID
	}

	// Convertir a JSON solo los campos que cambiaron
	previousValuesJSON, _ := json.Marshal(previousValues)
	newValuesJSON, _ := json.Marshal(newValues)

	changeDetails := "Lote actualizado: " + strings.Join(changes, ", ")

	history := &models.BatchHistory{
		DateTime:       time.Now(),
		BatchID:        &batchID,
		UserRUT:        userRUT,
		UserName:       user.Name,
		ChangeDetails:  changeDetails,
		PreviousValues: string(previousValuesJSON),
		NewValues:      string(newValuesJSON),
		BatchNumber:    batchID, // Agregar batch_number
	}

	return s.CreateBatchHistory(history)
}

// RegisterBatchDeletion registra la eliminación de un lote
func (s *BatchHistoryService) RegisterBatchDeletion(batchID int, userRUT string, deletedBatch *models.Batch) error {
	// RUT hardcodeado para el usuario
	if userRUT == "" {
		userRUT = "12345678-9"
	}

	// Obtener información del código de insumo desde medical_supply
	var medicalSupply models.MedicalSupply
	if err := s.DB.Where("batch_id = ?", batchID).First(&medicalSupply).Error; err != nil {
		return fmt.Errorf("error obteniendo información del insumo médico: %v", err)
	}

	// Obtener información del código de insumo
	var supplyCode models.SupplyCode
	if err := s.DB.Where("code = ?", medicalSupply.Code).First(&supplyCode).Error; err != nil {
		return fmt.Errorf("error obteniendo código de insumo: %v", err)
	}

	// Obtener información del usuario
	var user models.User
	if err := s.DB.Where("rut = ?", userRUT).First(&user).Error; err != nil {
		return fmt.Errorf("error obteniendo información del usuario: %v", err)
	}

	// Crear valores anteriores en formato JSON
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
		UserName:       user.Name, // Agregar el nombre del usuario
		ChangeDetails:  "Lote eliminado",
		PreviousValues: string(previousValues),
		NewValues:      "{}",    // JSON vacío válido para valores nuevos
		BatchNumber:    batchID, // Agregar batch_number
	}

	return s.CreateBatchHistory(history)
}

// GetAllBatchHistoriesWithDetails obtiene todos los registros del historial con información completa y formateada
func (s *BatchHistoryService) GetAllBatchHistoriesWithDetails() ([]map[string]interface{}, error) {
	var histories []models.BatchHistory

	// Obtener todos los historiales ordenados por fecha descendente (más recientes primero)
	if err := s.DB.Order("date_time DESC").Find(&histories).Error; err != nil {
		return nil, err
	}

	// Crear historial detallado formateado
	var detailedHistories []map[string]interface{}

	for _, history := range histories {
		entry := map[string]interface{}{
			"id":             history.ID,
			"batch_id":       history.BatchNumber,   // Usar batch_number en lugar de batch_id
			"change_details": history.ChangeDetails, // Detalles del cambio
			"date_time":      history.DateTime,      // Timestamp de cuando se realizó
			"user_rut":       history.UserRUT,       // RUT del usuario
			"user_name":      history.UserName,      // Nombre del usuario
		}

		// Agregar valores anteriores si existen
		if history.PreviousValues != "" {
			var prevValues map[string]interface{}
			if err := json.Unmarshal([]byte(history.PreviousValues), &prevValues); err == nil {
				entry["previous_values"] = prevValues
			}
		}

		// Agregar valores nuevos si existen
		if history.NewValues != "" {
			var newValues map[string]interface{}
			if err := json.Unmarshal([]byte(history.NewValues), &newValues); err == nil {
				entry["new_values"] = newValues
			}
		}

		detailedHistories = append(detailedHistories, entry)
	}

	return detailedHistories, nil
}

// GetAllBatchHistoriesWithDetailsPaginated obtiene todos los registros del historial con paginación
func (s *BatchHistoryService) GetAllBatchHistoriesWithDetailsPaginated(page, pageSize int) ([]map[string]interface{}, int64, error) {
	var total int64
	var histories []models.BatchHistory

	// Contar total de registros
	if err := s.DB.Model(&models.BatchHistory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calcular offset
	offset := (page - 1) * pageSize

	// Obtener historiales paginados ordenados por fecha descendente
	if err := s.DB.Order("date_time DESC").Offset(offset).Limit(pageSize).Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	// Crear historial detallado formateado
	var detailedHistories []map[string]interface{}

	for _, history := range histories {
		entry := map[string]interface{}{
			"id":             history.ID,
			"batch_id":       history.BatchNumber,   // Usar batch_number en lugar de batch_id
			"change_details": history.ChangeDetails, // Detalles del cambio
			"date_time":      history.DateTime,      // Timestamp de cuando se realizó
			"user_rut":       history.UserRUT,       // RUT del usuario
			"user_name":      history.UserName,      // Nombre del usuario
		}

		// Agregar valores anteriores si existen
		if history.PreviousValues != "" {
			var prevValues map[string]interface{}
			if err := json.Unmarshal([]byte(history.PreviousValues), &prevValues); err == nil {
				entry["previous_values"] = prevValues
			}
		}

		// Agregar valores nuevos si existen
		if history.NewValues != "" {
			var newValues map[string]interface{}
			if err := json.Unmarshal([]byte(history.NewValues), &newValues); err == nil {
				entry["new_values"] = newValues
			}
		}

		detailedHistories = append(detailedHistories, entry)
	}

	return detailedHistories, total, nil
}

// GetBatchHistoryByBatchNumber obtiene el historial de un lote específico por su número
func (s *BatchHistoryService) GetBatchHistoryByBatchNumber(batchNumber int) ([]models.BatchHistory, error) {
	var histories []models.BatchHistory
	if err := s.DB.Where("batch_number = ?", batchNumber).Order("date_time DESC").Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

// SearchBatchHistoryByBatchNumber busca en el historial por número de lote (incluye lotes eliminados)
func (s *BatchHistoryService) SearchBatchHistoryByBatchNumber(batchNumber int) ([]map[string]interface{}, error) {
	var histories []models.BatchHistory

	// Buscar por batch_number (que incluye lotes eliminados)
	if err := s.DB.Where("batch_number = ?", batchNumber).Order("date_time DESC").Find(&histories).Error; err != nil {
		return nil, err
	}

	// Formatear resultados
	var detailedHistories []map[string]interface{}
	for _, history := range histories {
		entry := map[string]interface{}{
			"id":             history.ID,
			"batch_id":       history.BatchNumber, // Usar batch_number
			"change_details": history.ChangeDetails,
			"date_time":      history.DateTime,
			"user_rut":       history.UserRUT,
			"user_name":      history.UserName,
		}

		// Agregar valores anteriores si existen
		if history.PreviousValues != "" {
			var prevValues map[string]interface{}
			if err := json.Unmarshal([]byte(history.PreviousValues), &prevValues); err == nil {
				entry["previous_values"] = prevValues
			}
		}

		// Agregar valores nuevos si existen
		if history.NewValues != "" {
			var newValues map[string]interface{}
			if err := json.Unmarshal([]byte(history.NewValues), &newValues); err == nil {
				entry["new_values"] = newValues
			}
		}

		detailedHistories = append(detailedHistories, entry)
	}

	return detailedHistories, nil
}
