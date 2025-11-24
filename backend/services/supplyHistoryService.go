package services

import (
	"database/sql"
	"meditrack/models"
	"time"

	"gorm.io/gorm"
)

type SupplyHistoryService struct {
	DB *gorm.DB
}

func NewSupplyHistoryService(db *gorm.DB) *SupplyHistoryService {
	return &SupplyHistoryService{DB: db}
}

// Funcionalidades básicas de la versión anterior (MANTENIDAS)
func (s *SupplyHistoryService) CreateSupplyHistory(history *models.SupplyHistory) error {
	return s.DB.Create(history).Error
}

func (s *SupplyHistoryService) GetSupplyHistoryByID(id int) (*models.SupplyHistory, error) {
	var history models.SupplyHistory
	if err := s.DB.First(&history, id).Error; err != nil {
		return nil, err
	}
	return &history, nil
}

func (s *SupplyHistoryService) DeleteSupplyHistory(id int) error {
	return s.DB.Delete(&models.SupplyHistory{}, id).Error
}

func (s *SupplyHistoryService) GetAllSupplyHistories() ([]models.SupplyHistory, error) {
	var histories []models.SupplyHistory
	if err := s.DB.Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

// GetAllSupplyHistoriesWithDetails obtiene todos los historiales con información del insumo
func (s *SupplyHistoryService) GetAllSupplyHistoriesWithDetails() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT 
			sh.id,
			sh.date_time,
			sh.status,
			sh.destination_type,
			sh.destination_id,
			sh.medical_supply_id,
			sh.user_rut,
			sh.notes,
			sh.origin_type,
			sh.origin_id,
			sh.confirmed_by,
			sh.confirmation_date,
			sh.transfer_notes,
			sc.name as supply_name,
			sc.code as supply_code,
			ms.qr_code
		FROM supply_history sh
		LEFT JOIN medical_supply ms ON sh.medical_supply_id = ms.id
		LEFT JOIN supply_code sc ON ms.code = sc.code
		ORDER BY sh.date_time DESC
	`

	rows, err := s.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, destinationID, medicalSupplyID int
		var destinationType, status, userRut, notes string
		var dateTime time.Time
		var confirmationDate sql.NullTime
		var originType, originID, confirmedBy, transferNotes, supplyName, qrCode *string
		var supplyCode *int

		err := rows.Scan(
			&id, &dateTime, &status, &destinationType, &destinationID,
			&medicalSupplyID, &userRut, &notes, &originType, &originID,
			&confirmedBy, &confirmationDate, &transferNotes,
			&supplyName, &supplyCode, &qrCode,
		)
		if err != nil {
			return nil, err
		}

		// Formatear fecha en formato ISO 8601 con zona horaria para evitar problemas de interpretación
		// Esto asegura que la fecha se interprete correctamente en el frontend
		dateTimeISO := dateTime.Format(time.RFC3339)
		
		// También formatear confirmation_date si existe
		var confirmationDateISO *string
		if confirmationDate.Valid {
			iso := confirmationDate.Time.Format(time.RFC3339)
			confirmationDateISO = &iso
		}

		result := map[string]interface{}{
			"id":                id,
			"date_time":         dateTimeISO, // Formato ISO 8601 con zona horaria
			"status":            status,
			"destination_type":  destinationType,
			"destination_id":    destinationID,
			"medical_supply_id": medicalSupplyID,
			"user_rut":          userRut,
			"notes":             notes,
			"origin_type":       originType,
			"origin_id":         originID,
			"confirmed_by":      confirmedBy,
			"confirmation_date": confirmationDateISO,
			"transfer_notes":    transferNotes,
			"supply_name":       supplyName,
			"supply_code":       supplyCode,
			"qr_code":           qrCode,
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *SupplyHistoryService) UpdateSupplyHistory(id int, history *models.SupplyHistory) error {
	var existing models.SupplyHistory
	if err := s.DB.First(&existing, id).Error; err != nil {
		return err
	}
	// Actualizar campos omitiendo ID y timestamps
	return s.DB.Model(&existing).Omit("id", "created_at", "updated_at").Updates(history).Error
}
