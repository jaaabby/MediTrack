package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"meditrack/models"
	"time"

	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type QRService struct {
	DB *gorm.DB
}

func NewQRService(db *gorm.DB) *QRService {
	return &QRService{DB: db}
}

// QRGenerationResponse representa la respuesta al generar un QR
type QRGenerationResponse struct {
	QRCode    string `json:"qr_code"`
	Type      string `json:"type"`
	ImageData string `json:"image_data,omitempty"` // Base64 encoded image
	ImageURL  string `json:"image_url,omitempty"`  // URL para descargar la imagen
}

// QRConsumptionRequest representa la solicitud de consumo de un QR
type QRConsumptionRequest struct {
	QRCode          string `json:"qr_code" binding:"required"`
	UserRUT         string `json:"user_rut" binding:"required"`
	DestinationType string `json:"destination_type" binding:"required"` // "pavilion" o "store"
	DestinationID   int    `json:"destination_id" binding:"required"`
	Notes           string `json:"notes,omitempty"`
}

// QRConsumptionResponse representa la respuesta del consumo de un QR
type QRConsumptionResponse struct {
	Success            bool                   `json:"success"`
	Message            string                 `json:"message"`
	ConsumedSupply     *models.MedicalSupply  `json:"consumed_supply,omitempty"`
	UpdatedBatch       map[string]interface{} `json:"updated_batch,omitempty"`
	RemainingAmount    int                    `json:"remaining_amount"`
	ConsumptionHistory map[string]interface{} `json:"consumption_history,omitempty"`
}

// GenerateUniqueQRCode genera un código QR único
func (s *QRService) GenerateUniqueQRCode(prefix string) (string, error) {
	for attempts := 0; attempts < 10; attempts++ {
		// Generar bytes aleatorios
		bytes := make([]byte, 8)
		if _, err := rand.Read(bytes); err != nil {
			return "", fmt.Errorf("error al generar bytes aleatorios: %w", err)
		}

		// Crear código con prefijo, timestamp y bytes aleatorios
		timestamp := time.Now().Unix()
		qrCode := fmt.Sprintf("%s_%d_%s", prefix, timestamp, hex.EncodeToString(bytes))

		// Verificar que no existe en batch
		var batchCount int64
		s.DB.Model(&models.Batch{}).Where("qr_code = ?", qrCode).Count(&batchCount)

		// Verificar que no existe en medical_supply
		var supplyCount int64
		s.DB.Model(&models.MedicalSupply{}).Where("qr_code = ?", qrCode).Count(&supplyCount)

		if batchCount == 0 && supplyCount == 0 {
			return qrCode, nil
		}
	}

	return "", fmt.Errorf("no se pudo generar un código QR único después de 10 intentos")
}

// GenerateQRImage genera una imagen QR y la retorna en base64
func (s *QRService) GenerateQRImage(qrCode string) (string, error) {
	// Generar QR con tamaño 256x256 y nivel de corrección Medium
	qrBytes, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("error al generar imagen QR: %w", err)
	}

	// Convertir a base64
	base64String := base64.StdEncoding.EncodeToString(qrBytes)
	return base64String, nil
}

// GenerateBatchQRCode genera un código QR para un lote con imagen
func (s *QRService) GenerateBatchQRCode() (*QRGenerationResponse, error) {
	qrCode, err := s.GenerateUniqueQRCode("BATCH")
	if err != nil {
		return nil, err
	}

	// Generar imagen QR
	imageData, err := s.GenerateQRImage(qrCode)
	if err != nil {
		return nil, err
	}

	response := &QRGenerationResponse{
		QRCode:    qrCode,
		Type:      "batch",
		ImageData: imageData,
		ImageURL:  fmt.Sprintf("/api/v1/qr/image/%s", qrCode),
	}

	return response, nil
}

// GenerateMedicalSupplyQRCode genera un código QR para un insumo médico individual con imagen
func (s *QRService) GenerateMedicalSupplyQRCode() (*QRGenerationResponse, error) {
	qrCode, err := s.GenerateUniqueQRCode("SUPPLY")
	if err != nil {
		return nil, err
	}

	// Generar imagen QR
	imageData, err := s.GenerateQRImage(qrCode)
	if err != nil {
		return nil, err
	}

	response := &QRGenerationResponse{
		QRCode:    qrCode,
		Type:      "medical_supply",
		ImageData: imageData,
		ImageURL:  fmt.Sprintf("/api/v1/qr/image/%s", qrCode),
	}

	return response, nil
}

// ConsumeSupplyByQR procesa el consumo de un insumo y actualiza automáticamente las cantidades del lote
func (s *QRService) ConsumeSupplyByQR(request QRConsumptionRequest) (*QRConsumptionResponse, error) {
	var supply models.MedicalSupply
	var response QRConsumptionResponse

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Buscar el insumo por QR
		if err := tx.Where("qr_code = ?", request.QRCode).First(&supply).Error; err != nil {
			return fmt.Errorf("insumo no encontrado con QR %s: %v", request.QRCode, err)
		}

		// Verificar que el insumo no haya sido consumido previamente
		var existingHistory models.SupplyHistory
		err := tx.Where("medical_supply_id = ? AND status = ?", supply.ID, "consumido").First(&existingHistory).Error
		if err == nil {
			return fmt.Errorf("el insumo con QR %s ya ha sido consumido el %s", request.QRCode, existingHistory.DateTime.Format("2006-01-02 15:04:05"))
		}

		// Obtener información del lote antes de la actualización
		var batch models.Batch
		if err := tx.Where("id = ?", supply.BatchID).First(&batch).Error; err != nil {
			return fmt.Errorf("lote no encontrado: %v", err)
		}

		// Verificar que hay stock disponible
		if batch.Amount <= 0 {
			return fmt.Errorf("no hay stock disponible en el lote %d", batch.ID)
		}

		// Crear historial de consumo
		history := models.SupplyHistory{
			DateTime:        time.Now(),
			Status:          "consumido",
			DestinationType: request.DestinationType,
			DestinationID:   request.DestinationID,
			MedicalSupplyID: supply.ID,
			UserRUT:         request.UserRUT,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("error creando historial de consumo: %v", err)
		}

		// Actualizar cantidad del lote (restar 1)
		newAmount := batch.Amount - 1
		if err := tx.Model(&batch).Update("amount", newAmount).Error; err != nil {
			return fmt.Errorf("error actualizando cantidad del lote: %v", err)
		}

		// Preparar respuesta
		response.Success = true
		response.Message = fmt.Sprintf("Insumo consumido exitosamente. Lote actualizado: %d -> %d unidades", batch.Amount, newAmount)
		response.ConsumedSupply = &supply
		response.RemainingAmount = newAmount

		// Obtener información actualizada del lote
		updatedBatchInfo := map[string]interface{}{
			"batch_id":        batch.ID,
			"previous_amount": batch.Amount,
			"new_amount":      newAmount,
			"supplier":        batch.Supplier,
			"expiration_date": batch.ExpirationDate,
		}
		response.UpdatedBatch = updatedBatchInfo

		// Información del historial de consumo
		consumptionInfo := map[string]interface{}{
			"consumption_id":   history.ID,
			"consumed_at":      history.DateTime,
			"consumed_by":      request.UserRUT,
			"destination_type": request.DestinationType,
			"destination_id":   request.DestinationID,
			"notes":            request.Notes,
		}
		response.ConsumptionHistory = consumptionInfo

		// Si el stock llegó a 0, podríamos agregar una alerta
		if newAmount == 0 {
			response.Message += " ¡ATENCIÓN: El lote se ha agotado completamente!"
		} else if newAmount <= 5 { // Alerta de stock bajo
			response.Message += fmt.Sprintf(" ¡ADVERTENCIA: Stock bajo, quedan %d unidades!", newAmount)
		}

		return nil
	})

	if err != nil {
		return &QRConsumptionResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &response, nil
}

// GetQRImage devuelve la imagen QR como bytes
func (s *QRService) GetQRImage(qrCode string) ([]byte, error) {
	// Verificar que el código QR existe
	_, _, err := s.ValidateQRCode(qrCode)
	if err != nil {
		return nil, fmt.Errorf("código QR no válido: %w", err)
	}

	// Generar imagen QR
	qrBytes, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("error al generar imagen QR: %w", err)
	}

	return qrBytes, nil
}

// GetQRImageHighRes devuelve la imagen QR en alta resolución
func (s *QRService) GetQRImageHighRes(qrCode string) ([]byte, error) {
	// Verificar que el código QR existe
	_, _, err := s.ValidateQRCode(qrCode)
	if err != nil {
		return nil, fmt.Errorf("código QR no válido: %w", err)
	}

	// Generar imagen QR en alta resolución (512x512)
	qrBytes, err := qrcode.Encode(qrCode, qrcode.High, 512)
	if err != nil {
		return nil, fmt.Errorf("error al generar imagen QR en alta resolución: %w", err)
	}

	return qrBytes, nil
}

// QRInfo representa la información completa de un código QR escaneado
type QRInfo struct {
	Type        string                    `json:"type"` // "batch" o "medical_supply"
	ID          int                       `json:"id"`
	QRCode      string                    `json:"qr_code"`
	BatchInfo   *models.Batch             `json:"batch_info,omitempty"`
	SupplyInfo  *MedicalSupplyWithDetails `json:"supply_info,omitempty"`
	SupplyCode  *models.SupplyCode        `json:"supply_code,omitempty"`
	History     []models.SupplyHistory    `json:"history,omitempty"`
	IsConsumed  bool                      `json:"is_consumed"`
	CanConsume  bool                      `json:"can_consume"`
	BatchStatus map[string]interface{}    `json:"batch_status,omitempty"`
}

type MedicalSupplyWithDetails struct {
	models.MedicalSupply
	SupplyCodeName string    `json:"supply_code_name"`
	BatchID        int       `json:"batch_id"`
	Supplier       string    `json:"supplier"`
	ExpirationDate time.Time `json:"expiration_date"`
	StoreName      string    `json:"store_name"`
	IsConsumed     bool      `json:"is_consumed"`
}

// GetQRInfo obtiene toda la información relacionada con un código QR
func (s *QRService) GetQRInfo(qrCode string) (*QRInfo, error) {
	info := &QRInfo{
		QRCode: qrCode,
	}

	// Buscar en batch primero
	var batch models.Batch
	if err := s.DB.Where("qr_code = ?", qrCode).First(&batch).Error; err == nil {
		info.Type = "batch"
		info.ID = batch.ID
		info.BatchInfo = &batch
		info.CanConsume = false // Los lotes no se consumen directamente

		// Obtener estadísticas del lote
		var totalSupplies, consumedSupplies int64
		s.DB.Model(&models.MedicalSupply{}).Where("batch_id = ?", batch.ID).Count(&totalSupplies)

		subquery := s.DB.Model(&models.SupplyHistory{}).
			Select("medical_supply_id").
			Where("status = ?", "consumido")
		s.DB.Model(&models.MedicalSupply{}).
			Where("batch_id = ? AND id IN (?)", batch.ID, subquery).
			Count(&consumedSupplies)

		info.BatchStatus = map[string]interface{}{
			"total_individual_supplies": totalSupplies,
			"consumed_supplies":         consumedSupplies,
			"available_supplies":        totalSupplies - consumedSupplies,
			"current_batch_amount":      batch.Amount,
			"amounts_synchronized":      batch.Amount == int(totalSupplies-consumedSupplies),
		}

		return info, nil
	}

	// Buscar en medical_supply
	var supply models.MedicalSupply
	if err := s.DB.Where("qr_code = ?", qrCode).First(&supply).Error; err == nil {
		info.Type = "medical_supply"
		info.ID = supply.ID

		// Verificar si ya fue consumido
		var consumptionHistory models.SupplyHistory
		isConsumed := s.DB.Where("medical_supply_id = ? AND status = ?", supply.ID, "consumido").First(&consumptionHistory).Error == nil
		info.IsConsumed = isConsumed
		info.CanConsume = !isConsumed

		// Obtener información detallada del supply
		supplyDetails, err := s.getMedicalSupplyDetails(supply.ID)
		if err != nil {
			return nil, err
		}
		supplyDetails.IsConsumed = isConsumed
		info.SupplyInfo = supplyDetails

		// Obtener supply code
		var supplyCode models.SupplyCode
		if err := s.DB.Where("code = ?", supply.Code).First(&supplyCode).Error; err == nil {
			info.SupplyCode = &supplyCode
		}

		// Obtener historial
		var history []models.SupplyHistory
		s.DB.Where("medical_supply_id = ?", supply.ID).Order("date_time DESC").Find(&history)
		info.History = history

		// Obtener información del lote asociado
		var batch models.Batch
		if err := s.DB.Where("id = ?", supply.BatchID).First(&batch).Error; err == nil {
			info.BatchStatus = map[string]interface{}{
				"batch_id":            batch.ID,
				"current_amount":      batch.Amount,
				"expiration_date":     batch.ExpirationDate,
				"supplier":            batch.Supplier,
				"has_available_stock": batch.Amount > 0,
			}
		}

		return info, nil
	}

	return nil, fmt.Errorf("código QR no encontrado: %s", qrCode)
}

// getMedicalSupplyDetails obtiene información detallada de un insumo médico
func (s *QRService) getMedicalSupplyDetails(supplyID int) (*MedicalSupplyWithDetails, error) {
	var result MedicalSupplyWithDetails

	query := `
		SELECT 
			ms.id,
			ms.code,
			ms.qr_code,
			ms.batch_id,
			sc.name as supply_code_name,
			b.supplier,
			b.expiration_date,
			st.name as store_name
		FROM medical_supply ms
		JOIN supply_code sc ON ms.code = sc.code
		JOIN batch b ON ms.batch_id = b.id
		JOIN store st ON b.store_id = st.id
		WHERE ms.id = ?
	`

	row := s.DB.Raw(query, supplyID).Row()
	err := row.Scan(
		&result.ID,
		&result.Code,
		&result.QRCode,
		&result.BatchID,
		&result.SupplyCodeName,
		&result.Supplier,
		&result.ExpirationDate,
		&result.StoreName,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ValidateQRCode valida si un código QR existe y está activo
func (s *QRService) ValidateQRCode(qrCode string) (bool, string, error) {
	// Buscar en batch
	var batchCount int64
	s.DB.Model(&models.Batch{}).Where("qr_code = ?", qrCode).Count(&batchCount)
	if batchCount > 0 {
		return true, "batch", nil
	}

	// Buscar en medical_supply
	var supplyCount int64
	s.DB.Model(&models.MedicalSupply{}).Where("qr_code = ?", qrCode).Count(&supplyCount)
	if supplyCount > 0 {
		return true, "medical_supply", nil
	}

	return false, "", fmt.Errorf("código QR no válido: %s", qrCode)
}

// GetSupplyHistory obtiene el historial completo de un insumo por su código QR
func (s *QRService) GetSupplyHistory(qrCode string) ([]models.SupplyHistory, error) {
	var supply models.MedicalSupply
	if err := s.DB.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
		return nil, fmt.Errorf("insumo no encontrado con QR: %s", qrCode)
	}

	var history []models.SupplyHistory
	if err := s.DB.Where("medical_supply_id = ?", supply.ID).Order("date_time DESC").Find(&history).Error; err != nil {
		return nil, err
	}

	return history, nil
}

// GetQRStats obtiene estadísticas generales de uso de QR codes
func (s *QRService) GetQRStats() (map[string]interface{}, error) {
	var totalBatches, totalSupplies, consumedSupplies int64

	// Contar lotes totales
	s.DB.Model(&models.Batch{}).Count(&totalBatches)

	// Contar insumos totales
	s.DB.Model(&models.MedicalSupply{}).Count(&totalSupplies)

	// Contar insumos consumidos
	s.DB.Model(&models.SupplyHistory{}).Where("status = ?", "consumido").Count(&consumedSupplies)

	return map[string]interface{}{
		"total_batches":      totalBatches,
		"total_supplies":     totalSupplies,
		"consumed_supplies":  consumedSupplies,
		"available_supplies": totalSupplies - consumedSupplies,
		"consumption_rate":   float64(consumedSupplies) / float64(totalSupplies) * 100,
		"generated_at":       time.Now(),
	}, nil
}
