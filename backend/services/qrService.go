package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"meditrack/models"
	"time"

	"gorm.io/gorm"
)

type QRService struct {
	DB *gorm.DB
}

func NewQRService(db *gorm.DB) *QRService {
	return &QRService{DB: db}
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

// GenerateBatchQRCode genera un código QR para un lote
func (s *QRService) GenerateBatchQRCode() (string, error) {
	return s.GenerateUniqueQRCode("BATCH")
}

// GenerateMedicalSupplyQRCode genera un código QR para un insumo médico individual
func (s *QRService) GenerateMedicalSupplyQRCode() (string, error) {
	return s.GenerateUniqueQRCode("SUPPLY")
}

// QRInfo representa la información completa de un código QR escaneado
type QRInfo struct {
	Type         string      `json:"type"`         // "batch" o "medical_supply"
	ID           int         `json:"id"`
	QRCode       string      `json:"qr_code"`
	BatchInfo    *models.Batch `json:"batch_info,omitempty"`
	SupplyInfo   *MedicalSupplyWithDetails `json:"supply_info,omitempty"`
	SupplyCode   *models.SupplyCode `json:"supply_code,omitempty"`
	History      []models.SupplyHistory `json:"history,omitempty"`
}

type MedicalSupplyWithDetails struct {
	models.MedicalSupply
	SupplyCodeName string `json:"supply_code_name"`
	BatchID        int    `json:"batch_id"`
	Supplier       string `json:"supplier"`
	ExpirationDate time.Time `json:"expiration_date"`
	StoreName      string `json:"store_name"`
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
		
		// Obtener supply codes relacionados al batch
		var supplyCodes []models.SupplyCode
		s.DB.Where("batch_id = ?", batch.ID).Find(&supplyCodes)
		
		return info, nil
	}
	
	// Buscar en medical_supply
	var supply models.MedicalSupply
	if err := s.DB.Where("qr_code = ?", qrCode).First(&supply).Error; err == nil {
		info.Type = "medical_supply"
		info.ID = supply.ID
		
		// Obtener información detallada del supply
		supplyDetails, err := s.getMedicalSupplyDetails(supply.ID)
		if err != nil {
			return nil, err
		}
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
			sc.name as supply_code_name,
			sc.batch_id,
			b.supplier,
			b.expiration_date,
			st.name as store_name
		FROM medical_supply ms
		JOIN supply_code sc ON ms.code = sc.code
		JOIN batch b ON sc.batch_id = b.id
		JOIN store st ON b.store_id = st.id
		WHERE ms.id = ?
	`
	
	row := s.DB.Raw(query, supplyID).Row()
	err := row.Scan(
		&result.ID,
		&result.Code,
		&result.QRCode,
		&result.SupplyCodeName,
		&result.BatchID,
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