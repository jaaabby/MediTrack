package services

import (
	"errors"
	"fmt"
	"log"
	"meditrack/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SupplyTransferService struct {
	DB *gorm.DB
}

func NewSupplyTransferService(db *gorm.DB) *SupplyTransferService {
	return &SupplyTransferService{DB: db}
}

// TransferToPavilion transfiere insumos de bodega a pabellón con descuentos automáticos
func (s *SupplyTransferService) TransferToPavilion(
	qrCodes []string,
	pavilionID int,
	userRUT string,
	userName string,
	reason string,
	notes string,
) ([]models.SupplyTransfer, error) {
	var transfers []models.SupplyTransfer

	// Iniciar transacción
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		for _, qrCode := range qrCodes {
			// 1. Validar que el insumo existe y está en bodega
			var supply models.MedicalSupply
			if err := tx.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
				return fmt.Errorf("insumo con QR %s no encontrado: %v", qrCode, err)
			}

			if supply.LocationType != models.SupplyLocationStore {
				return fmt.Errorf("insumo %s no está en bodega", qrCode)
			}

			if supply.InTransit {
				return fmt.Errorf("insumo %s ya está en tránsito", qrCode)
			}

			if supply.Status == models.StatusConsumed {
				return fmt.Errorf("insumo %s ya fue consumido", qrCode)
			}

			// 2. Obtener información del batch
			var batch models.Batch
			if err := tx.First(&batch, supply.BatchID).Error; err != nil {
				return fmt.Errorf("lote no encontrado: %v", err)
			}

			// 3. Crear registro de transferencia
			transferCode := fmt.Sprintf("TRANS-%d-%s", time.Now().Unix(), qrCode[len(qrCode)-5:])
			transfer := models.SupplyTransfer{
				TransferCode:    transferCode,
				QRCode:          qrCode,
				MedicalSupplyID: supply.ID,
				OriginType:      models.TransferLocationStore,
				OriginID:        supply.LocationID,
				DestinationType: models.TransferLocationPavilion,
				DestinationID:   pavilionID,
				SentBy:          userRUT,
				SentByName:      userName,
				Status:          models.TransferStatusInTransit,
				TransferReason:  reason,
				SendDate:        time.Now(),
				Notes:           notes,
			}

			if err := tx.Create(&transfer).Error; err != nil {
				return fmt.Errorf("error al crear transferencia: %v", err)
			}

			// 4. Actualizar ubicación del insumo y marcarlo en tránsito
			now := time.Now()
			supply.LocationType = models.SupplyLocationPavilion
			supply.LocationID = pavilionID
			supply.InTransit = true
			supply.TransferDate = &now
			supply.TransferredBy = &userRUT
			supply.Status = models.StatusEnRouteToPavilion

			if err := tx.Save(&supply).Error; err != nil {
				return fmt.Errorf("error al actualizar insumo: %v", err)
			}

			// 5. Descontar del stock de bodega
			var storeSummary models.StoreInventorySummary
			if err := tx.Where("store_id = ? AND batch_id = ?", supply.LocationID, batch.ID).First(&storeSummary).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Si no existe el resumen, calcular el stock real en bodega
					var realCount int64
					tx.Model(&models.MedicalSupply{}).
						Where("batch_id = ? AND location_type = ? AND location_id = ? AND status != ?",
							batch.ID, models.SupplyLocationStore, supply.LocationID, models.StatusConsumed).
						Count(&realCount)

					// Crear resumen si no existe
					storeSummary = models.StoreInventorySummary{
						StoreID:             supply.LocationID,
						BatchID:             batch.ID,
						SupplyCode:          supply.Code,
						SurgeryID:           batch.SurgeryID,
						OriginalAmount:      int(realCount) + 1, // +1 porque vamos a transferir uno
						CurrentInStore:      int(realCount),     // Stock actual sin contar el que se transfiere
						TotalTransferredOut: 1,
					}
					now := time.Now()
					storeSummary.LastTransferOutDate = &now
					if err := tx.Create(&storeSummary).Error; err != nil {
						return fmt.Errorf("error al crear resumen de bodega: %v", err)
					}
				} else {
					return fmt.Errorf("error al obtener resumen de bodega: %v", err)
				}
			} else {
				// Actualizar resumen existente
				if storeSummary.CurrentInStore > 0 {
					storeSummary.CurrentInStore--
				}
				storeSummary.TotalTransferredOut++
				now := time.Now()
				storeSummary.LastTransferOutDate = &now

				if err := tx.Save(&storeSummary).Error; err != nil {
					return fmt.Errorf("error al actualizar resumen de bodega: %v", err)
				}
			}

			// 6. Actualizar cantidad del lote
			batch.Amount--
			if err := tx.Save(&batch).Error; err != nil {
				return fmt.Errorf("error al actualizar lote: %v", err)
			}

			// 7. Registrar en historial
			history := models.SupplyHistory{
				DateTime:        time.Now(),
				Status:          models.StatusEnRouteToPavilion,
				DestinationType: models.DestinationTypePavilion,
				DestinationID:   pavilionID,
				MedicalSupplyID: supply.ID,
				UserRUT:         userRUT,
				Notes:           fmt.Sprintf("Transferencia iniciada: %s", reason),
				OriginType:      func() *string { s := models.DestinationTypeStore; return &s }(),
				OriginID:        &supply.LocationID,
				TransferNotes:   notes,
			}

			if err := tx.Create(&history).Error; err != nil {
				return fmt.Errorf("error al crear historial: %v", err)
			}

			transfers = append(transfers, transfer)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return transfers, nil
}

// ConfirmReception confirma la recepción de insumos en el pabellón
func (s *SupplyTransferService) ConfirmReception(
	transferCode string,
	userRUT string,
	userName string,
	notes string,
) (*models.SupplyTransfer, error) {
	var transfer models.SupplyTransfer
	var supply models.MedicalSupply

	// Iniciar transacción
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Validar que la transferencia existe y está en tránsito
		if err := tx.Where("transfer_code = ?", transferCode).First(&transfer).Error; err != nil {
			return fmt.Errorf("transferencia no encontrada: %v", err)
		}

		if transfer.Status != models.TransferStatusInTransit {
			return fmt.Errorf("la transferencia no está en tránsito (estado: %s)", transfer.Status)
		}

		// 1.5. Validar que la persona que recibe sea la misma que retiró físicamente
		// La persona que sacó de bodega debe confirmar que el insumo llegó
		if transfer.PickedUpBy == nil {
			return fmt.Errorf("el insumo no ha sido retirado físicamente de bodega aún. Debe escanearlo primero para registrar el retiro.")
		}
		if *transfer.PickedUpBy != userRUT {
			// Obtener el nombre del usuario que retiró el insumo
			var pickedUpUser models.User
			pickedUpUserName := "Usuario Desconocido"
			if err := tx.Where("rut = ?", *transfer.PickedUpBy).First(&pickedUpUser).Error; err == nil {
				pickedUpUserName = pickedUpUser.Name
			}
			
			// Obtener el nombre del usuario actual si no está disponible
			currentUserName := userName
			if currentUserName == "" {
				var currentUser models.User
				if err := tx.Where("rut = ?", userRUT).First(&currentUser).Error; err == nil {
					currentUserName = currentUser.Name
				} else {
					currentUserName = "Usuario Desconocido"
				}
			}
			
			return fmt.Errorf("solo la persona que retiró el insumo de bodega (%s) puede confirmar su recepción. Usted es %s", pickedUpUserName, currentUserName)
		}

		// 2. Obtener el insumo
		if err := tx.First(&supply, transfer.MedicalSupplyID).Error; err != nil {
			return fmt.Errorf("insumo no encontrado: %v", err)
		}

		// 3. Actualizar transferencia
		now := time.Now()
		transfer.Status = models.TransferStatusReceived
		transfer.ReceiveDate = &now
		transfer.ReceivedBy = &userRUT
		transfer.ReceivedByName = &userName
		if notes != "" {
			transfer.Notes = transfer.Notes + "\n" + notes
		}

		if err := tx.Save(&transfer).Error; err != nil {
			return fmt.Errorf("error al actualizar transferencia: %v", err)
		}

		// 4. Quitar marca de tránsito del insumo
		supply.InTransit = false
		supply.Status = models.StatusReceived

		if err := tx.Save(&supply).Error; err != nil {
			return fmt.Errorf("error al actualizar insumo: %v", err)
		}

		// 5. Incrementar contador de stock disponible en pabellón
		var pavilionSummary models.PavilionInventorySummary
		if err := tx.Where("pavilion_id = ? AND batch_id = ?", transfer.DestinationID, supply.BatchID).
			First(&pavilionSummary).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Crear resumen si no existe
				pavilionSummary = models.PavilionInventorySummary{
					PavilionID:       transfer.DestinationID,
					BatchID:          supply.BatchID,
					SupplyCode:       supply.Code,
					TotalReceived:    1,
					CurrentAvailable: 1,
					LastReceivedDate: &now,
				}
				if err := tx.Create(&pavilionSummary).Error; err != nil {
					return fmt.Errorf("error al crear resumen de pabellón: %v", err)
				}
			} else {
				return fmt.Errorf("error al obtener resumen de pabellón: %v", err)
			}
		} else {
			// Actualizar resumen existente
			pavilionSummary.TotalReceived++
			pavilionSummary.CurrentAvailable++
			pavilionSummary.LastReceivedDate = &now

			if err := tx.Save(&pavilionSummary).Error; err != nil {
				return fmt.Errorf("error al actualizar resumen de pabellón: %v", err)
			}
		}

		// 6. Registrar confirmación en historial
		history := models.SupplyHistory{
			DateTime:         now,
			Status:           models.StatusReceived,
			DestinationType:  models.DestinationTypePavilion,
			DestinationID:    transfer.DestinationID,
			MedicalSupplyID:  supply.ID,
			UserRUT:          userRUT,
			Notes:            "Recepción confirmada en pabellón",
			OriginType:       &transfer.OriginType,
			OriginID:         &transfer.OriginID,
			ConfirmedBy:      &userRUT,
			ConfirmationDate: &now,
			TransferNotes:    notes,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("error al crear historial de confirmación: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &transfer, nil
}

// ReturnToStore devuelve insumos de pabellón a bodega
func (s *SupplyTransferService) ReturnToStore(
	qrCodes []string,
	storeID int,
	userRUT string,
	userName string,
	reason string,
	notes string,
) ([]models.SupplyTransfer, error) {
	var transfers []models.SupplyTransfer

	// Iniciar transacción
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		for _, qrCode := range qrCodes {
			// 1. Validar que el insumo existe y está en pabellón
			var supply models.MedicalSupply
			if err := tx.Where("qr_code = ?", qrCode).First(&supply).Error; err != nil {
				return fmt.Errorf("insumo con QR %s no encontrado: %v", qrCode, err)
			}

			if supply.LocationType != models.SupplyLocationPavilion {
				return fmt.Errorf("insumo %s no está en pabellón", qrCode)
			}

			// Permitir devolver insumos consumidos automáticamente
			// Verificar si fue consumido automáticamente revisando el historial
			if supply.Status == models.StatusConsumed {
				// Verificar si fue consumido automáticamente
				var lastConsumptionHistory models.SupplyHistory
				if err := tx.Where("medical_supply_id = ? AND status = ?", supply.ID, models.StatusConsumed).
					Order("date_time DESC").
					First(&lastConsumptionHistory).Error; err == nil {
					// Si las notas contienen el prefijo de consumo automático, permitir devolución
					if strings.Contains(lastConsumptionHistory.Notes, "[CONSUMO_AUTOMATICO]") {
						// Permitir devolución de insumo consumido automáticamente
						log.Printf("🔄 Permitiendo devolución de insumo %s consumido automáticamente", qrCode)
					} else {
						return fmt.Errorf("insumo %s ya fue consumido manualmente y no puede ser devuelto", qrCode)
					}
				} else {
					return fmt.Errorf("insumo %s ya fue consumido", qrCode)
				}
			}

			// 2. Obtener información del batch
			var batch models.Batch
			if err := tx.First(&batch, supply.BatchID).Error; err != nil {
				return fmt.Errorf("lote no encontrado: %v", err)
			}

			pavilionID := supply.LocationID

			// 3. Crear registro de transferencia de devolución
			transferCode := fmt.Sprintf("RETURN-%d-%s", time.Now().Unix(), qrCode[len(qrCode)-5:])
			transfer := models.SupplyTransfer{
				TransferCode:    transferCode,
				QRCode:          qrCode,
				MedicalSupplyID: supply.ID,
				OriginType:      models.TransferLocationPavilion,
				OriginID:        pavilionID,
				DestinationType: models.TransferLocationStore,
				DestinationID:   storeID,
				SentBy:          userRUT,
				SentByName:      userName,
				ReceivedBy:      &userRUT,
				ReceivedByName:  &userName,
				Status:          models.TransferStatusReceived, // Devolución se considera recibida inmediatamente
				TransferReason:  reason,
				SendDate:        time.Now(),
				Notes:           notes,
			}
			now := time.Now()
			transfer.ReceiveDate = &now

			if err := tx.Create(&transfer).Error; err != nil {
				return fmt.Errorf("error al crear transferencia de devolución: %v", err)
			}

			// 4. Decrementar stock del pabellón (solo si el insumo no estaba consumido)
			// Si estaba consumido automáticamente, no estaba en el inventario disponible
			if supply.Status != models.StatusConsumed {
				var pavilionSummary models.PavilionInventorySummary
				if err := tx.Where("pavilion_id = ? AND batch_id = ?", pavilionID, batch.ID).
					First(&pavilionSummary).Error; err != nil {
					return fmt.Errorf("resumen de pabellón no encontrado: %v", err)
				}

				pavilionSummary.CurrentAvailable--
				pavilionSummary.TotalReturned++
				pavilionSummary.LastReturnedDate = &now

				if err := tx.Save(&pavilionSummary).Error; err != nil {
					return fmt.Errorf("error al actualizar resumen de pabellón: %v", err)
				}
			} else {
				// Si estaba consumido automáticamente, solo actualizar contadores de devolución
				var pavilionSummary models.PavilionInventorySummary
				if err := tx.Where("pavilion_id = ? AND batch_id = ?", pavilionID, batch.ID).
					First(&pavilionSummary).Error; err == nil {
					pavilionSummary.TotalReturned++
					pavilionSummary.LastReturnedDate = &now
					tx.Save(&pavilionSummary)
				}
			}

			// 5. Incrementar stock de bodega
			var storeSummary models.StoreInventorySummary
			if err := tx.Where("batch_id = ?", batch.ID).First(&storeSummary).Error; err != nil {
				return fmt.Errorf("resumen de bodega no encontrado: %v", err)
			}

			storeSummary.CurrentInStore++
			storeSummary.TotalReturnedIn++
			storeSummary.LastReturnInDate = &now

			if err := tx.Save(&storeSummary).Error; err != nil {
				return fmt.Errorf("error al actualizar resumen de bodega: %v", err)
			}

			// 6. Incrementar cantidad del lote
			batch.Amount++
			if err := tx.Save(&batch).Error; err != nil {
				return fmt.Errorf("error al actualizar lote: %v", err)
			}

			// 7. Actualizar ubicación del insumo
			supply.LocationType = models.SupplyLocationStore
			supply.LocationID = storeID
			supply.InTransit = false
			supply.Status = models.StatusAvailable
			supply.TransferDate = &now
			supply.TransferredBy = &userRUT

			if err := tx.Save(&supply).Error; err != nil {
				return fmt.Errorf("error al actualizar insumo: %v", err)
			}

			// 8. Registrar en historial
			history := models.SupplyHistory{
				DateTime:         now,
				Status:           models.StatusAvailable,
				DestinationType:  models.DestinationTypeStore,
				DestinationID:    storeID,
				MedicalSupplyID:  supply.ID,
				UserRUT:          userRUT,
				Notes:            fmt.Sprintf("Devolución de pabellón: %s", reason),
				OriginType:       func() *string { s := models.DestinationTypePavilion; return &s }(),
				OriginID:         &pavilionID,
				ConfirmedBy:      &userRUT,
				ConfirmationDate: &now,
				TransferNotes:    notes,
			}

			if err := tx.Create(&history).Error; err != nil {
				return fmt.Errorf("error al crear historial: %v", err)
			}

			transfers = append(transfers, transfer)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return transfers, nil
}

// GetTransferByCode obtiene una transferencia por su código
func (s *SupplyTransferService) GetTransferByCode(transferCode string) (*models.SupplyTransferWithDetails, error) {
	var transfer models.SupplyTransferWithDetails

	err := s.DB.Table("supply_transfer st").
		Select(`st.*, 
			ms.qr_code as qr_code,
			sc.code as supply_code,
			sc.name as supply_name,
			b.id as batch_number,
			CASE 
				WHEN st.origin_type = 'store' THEN s.name
				WHEN st.origin_type = 'pavilion' THEN p.name
			END as origin_name,
			CASE 
				WHEN st.destination_type = 'store' THEN s2.name
				WHEN st.destination_type = 'pavilion' THEN p2.name
			END as destination_name,
			mc.name as medical_center_name`).
		Joins("LEFT JOIN medical_supply ms ON st.medical_supply_id = ms.id").
		Joins("LEFT JOIN supply_code sc ON ms.code = sc.code").
		Joins("LEFT JOIN batch b ON ms.batch_id = b.id").
		Joins("LEFT JOIN store s ON st.origin_type = 'store' AND st.origin_id = s.id").
		Joins("LEFT JOIN pavilion p ON st.origin_type = 'pavilion' AND st.origin_id = p.id").
		Joins("LEFT JOIN store s2 ON st.destination_type = 'store' AND st.destination_id = s2.id").
		Joins("LEFT JOIN pavilion p2 ON st.destination_type = 'pavilion' AND st.destination_id = p2.id").
		Joins("LEFT JOIN medical_center mc ON (s.medical_center_id = mc.id OR p.medical_center_id = mc.id)").
		Where("st.transfer_code = ?", transferCode).
		First(&transfer).Error

	if err != nil {
		return nil, err
	}

	return &transfer, nil
}

// GetTransfersByFilters obtiene transferencias con filtros
func (s *SupplyTransferService) GetTransfersByFilters(
	status *string,
	originType *string,
	originID *int,
	destinationType *string,
	destinationID *int,
	startDate *time.Time,
	endDate *time.Time,
	page int,
	pageSize int,
) ([]models.SupplyTransferWithDetails, int64, error) {
	var transfers []models.SupplyTransferWithDetails
	var total int64

	query := s.DB.Table("supply_transfer st").
		Select(`st.*, 
			sc.code as supply_code,
			sc.name as supply_name,
			b.id as batch_number,
			CASE 
				WHEN st.origin_type = 'store' THEN s.name
				WHEN st.origin_type = 'pavilion' THEN p.name
			END as origin_name,
			CASE 
				WHEN st.destination_type = 'store' THEN s2.name
				WHEN st.destination_type = 'pavilion' THEN p2.name
			END as destination_name`).
		Joins("LEFT JOIN medical_supply ms ON st.medical_supply_id = ms.id").
		Joins("LEFT JOIN supply_code sc ON ms.code = sc.code").
		Joins("LEFT JOIN batch b ON ms.batch_id = b.id").
		Joins("LEFT JOIN store s ON st.origin_type = 'store' AND st.origin_id = s.id").
		Joins("LEFT JOIN pavilion p ON st.origin_type = 'pavilion' AND st.origin_id = p.id").
		Joins("LEFT JOIN store s2 ON st.destination_type = 'store' AND st.destination_id = s2.id").
		Joins("LEFT JOIN pavilion p2 ON st.destination_type = 'pavilion' AND st.destination_id = p2.id")

	// Aplicar filtros
	if status != nil {
		query = query.Where("st.status = ?", *status)
	}
	if originType != nil {
		query = query.Where("st.origin_type = ?", *originType)
	}
	if originID != nil {
		query = query.Where("st.origin_id = ?", *originID)
	}
	if destinationType != nil {
		query = query.Where("st.destination_type = ?", *destinationType)
	}
	if destinationID != nil {
		query = query.Where("st.destination_id = ?", *destinationID)
	}
	if startDate != nil {
		query = query.Where("st.send_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("st.send_date <= ?", *endDate)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	offset := (page - 1) * pageSize
	if err := query.Order("st.send_date DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&transfers).Error; err != nil {
		return nil, 0, err
	}

	return transfers, total, nil
}

// CancelTransfer cancela una transferencia pendiente
func (s *SupplyTransferService) CancelTransfer(
	transferCode string,
	userRUT string,
	reason string,
) (*models.SupplyTransfer, error) {
	var transfer models.SupplyTransfer

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Obtener transferencia
		if err := tx.Where("transfer_code = ?", transferCode).First(&transfer).Error; err != nil {
			return fmt.Errorf("transferencia no encontrada: %v", err)
		}

		if !transfer.CanBeCancelled() {
			return fmt.Errorf("la transferencia no puede ser cancelada (estado: %s)", transfer.Status)
		}

		// 2. Obtener insumo
		var supply models.MedicalSupply
		if err := tx.First(&supply, transfer.MedicalSupplyID).Error; err != nil {
			return fmt.Errorf("insumo no encontrado: %v", err)
		}

		// 3. Revertir cambios en el insumo
		supply.LocationType = models.SupplyLocationStore
		supply.LocationID = transfer.OriginID
		supply.InTransit = false
		supply.Status = models.StatusAvailable

		if err := tx.Save(&supply).Error; err != nil {
			return fmt.Errorf("error al revertir insumo: %v", err)
		}

		// 4. Revertir cambios en resumen de bodega y lote
		var batch models.Batch
		if err := tx.First(&batch, supply.BatchID).Error; err != nil {
			return fmt.Errorf("lote no encontrado: %v", err)
		}

		var storeSummary models.StoreInventorySummary
		if err := tx.Where("batch_id = ?", batch.ID).First(&storeSummary).Error; err == nil {
			storeSummary.CurrentInStore++
			storeSummary.TotalTransferredOut--
			if err := tx.Save(&storeSummary).Error; err != nil {
				return fmt.Errorf("error al revertir resumen de bodega: %v", err)
			}
		}

		batch.Amount++
		if err := tx.Save(&batch).Error; err != nil {
			return fmt.Errorf("error al revertir lote: %v", err)
		}

		// 5. Actualizar transferencia
		transfer.Status = models.TransferStatusCancelled
		transfer.RejectionReason = &reason

		if err := tx.Save(&transfer).Error; err != nil {
			return fmt.Errorf("error al cancelar transferencia: %v", err)
		}

		// 6. Registrar en historial
		history := models.SupplyHistory{
			DateTime:        time.Now(),
			Status:          models.StatusAvailable,
			DestinationType: models.DestinationTypeStore,
			DestinationID:   transfer.OriginID,
			MedicalSupplyID: supply.ID,
			UserRUT:         userRUT,
			Notes:           fmt.Sprintf("Transferencia cancelada: %s", reason),
			TransferNotes:   reason,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("error al crear historial: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &transfer, nil
}
