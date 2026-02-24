package services

import (
	"fmt"
	"meditrack/models"
	"time"

	"gorm.io/gorm"
)

type InventoryService struct {
	DB *gorm.DB
}

func NewInventoryService(db *gorm.DB) *InventoryService {
	return &InventoryService{DB: db}
}

// GetStoreInventory obtiene el inventario de una bodega con filtros
func (s *InventoryService) GetStoreInventory(
	storeID *int,
	surgeryID *int,
	supplyCode *int,
	supplier *string,
	nearExpiration bool,
	lowStock bool,
	page int,
	pageSize int,
) ([]models.StoreInventorySummaryWithDetails, int64, error) {
	var inventory []models.StoreInventorySummaryWithDetails
	var total int64

	query := s.DB.Table("store_inventory_summary sis").
		Select(`sis.id,
			sis.store_id,
			sis.batch_id,
			sis.supply_code,
			sis.surgery_id,
			sis.original_amount,
			sis.current_in_store,
			sis.total_transferred_out,
			sis.total_returned_in,
			sis.total_consumed_in_store,
			sis.last_transfer_out_date,
			sis.last_return_in_date,
			sis.last_consumed_date,
			sis.created_at,
			sis.updated_at,
			s.name as store_name,
			sc.name as supply_name,
			surg.name as surgery_name,
			b.supplier as batch_supplier,
			COALESCE(b.expiration_date::text, '') as expiration_date,
			mc.id as medical_center_id,
			mc.name as medical_center_name,
			COALESCE((SELECT SUM(pis.total_consumed)
				FROM pavilion_inventory_summary pis
				WHERE pis.batch_id = sis.batch_id), 0) AS total_consumed_from_pavilions`).
		Joins("LEFT JOIN store s ON sis.store_id = s.id").
		Joins("LEFT JOIN supply_code sc ON sis.supply_code = sc.code").
		Joins("LEFT JOIN surgery surg ON sis.surgery_id = surg.id").
		Joins("LEFT JOIN batch b ON sis.batch_id = b.id").
		Joins("LEFT JOIN medical_center mc ON s.medical_center_id = mc.id")

	// Aplicar filtros
	if storeID != nil {
		query = query.Where("sis.store_id = ?", *storeID)
	}
	if surgeryID != nil {
		query = query.Where("sis.surgery_id = ?", *surgeryID)
	}
	if supplyCode != nil {
		query = query.Where("sis.supply_code = ?", *supplyCode)
	}
	if supplier != nil {
		query = query.Where("b.supplier LIKE ?", "%"+*supplier+"%")
	}
	if nearExpiration {
		// Productos que vencen en los próximos 90 días
		expirationDate := time.Now().AddDate(0, 0, 90)
		query = query.Where("b.expiration_date <= ?", expirationDate)
	}
	if lowStock {
		// Stock bajo (menos del 20% del original)
		query = query.Where("sis.current_in_store < (sis.original_amount * 0.2)")
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	offset := (page - 1) * pageSize
	if err := query.Order("sis.current_in_store ASC, b.expiration_date ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&inventory).Error; err != nil {
		return nil, 0, err
	}

	return inventory, total, nil
}

// GetPavilionInventory obtiene el inventario de un pabellón
func (s *InventoryService) GetPavilionInventory(
	pavilionID int,
	includeInTransit bool,
	supplier *string,
) ([]models.PavilionInventorySummaryWithDetails, error) {
	var inventory []models.PavilionInventorySummaryWithDetails

	query := s.DB.Table("pavilion_inventory_summary pis").
		Select(`pis.*,
			p.name as pavilion_name,
			sc.name as supply_name,
			b.supplier as batch_supplier,
			b.expiration_date as expiration_date,
			mc.id as medical_center_id,
			mc.name as medical_center_name`).
		Joins("LEFT JOIN pavilion p ON pis.pavilion_id = p.id").
		Joins("LEFT JOIN supply_code sc ON pis.supply_code = sc.code").
		Joins("LEFT JOIN batch b ON pis.batch_id = b.id").
		Joins("LEFT JOIN medical_center mc ON p.medical_center_id = mc.id").
		Where("pis.pavilion_id = ?", pavilionID)

	// Aplicar filtro por proveedor (case-insensitive) si se especifica
	if supplier != nil {
		query = query.Where("LOWER(b.supplier) LIKE LOWER(?)", "%"+*supplier+"%")
	}

	if err := query.Find(&inventory).Error; err != nil {
		return nil, err
	}

	// Si se solicita incluir productos en tránsito:
	// Busca en medical_supply los insumos con status 'en_camino_a_pabellon'
	// destinados a este pabellón, agrupados por lote+código.
	if includeInTransit {
		type inTransitRow struct {
			PavilionID        int     `gorm:"column:pavilion_id"`
			BatchID           int     `gorm:"column:batch_id"`
			SupplyCode        int     `gorm:"column:supply_code"`
			CurrentAvailable  int     `gorm:"column:current_available"`
			PavilionName      string  `gorm:"column:pavilion_name"`
			SupplyName        string  `gorm:"column:supply_name"`
			BatchSupplier     string  `gorm:"column:batch_supplier"`
			ExpirationDate    string  `gorm:"column:expiration_date"`
			MedicalCenterID   int     `gorm:"column:medical_center_id"`
			MedicalCenterName *string `gorm:"column:medical_center_name"`
		}

		var inTransitRows []inTransitRow
		err := s.DB.Table("medical_supply ms").
			Select(`ms.location_id AS pavilion_id,
				ms.batch_id,
				ms.code AS supply_code,
				COUNT(*) AS current_available,
				p.name AS pavilion_name,
				sc.name AS supply_name,
				b.supplier AS batch_supplier,
				TO_CHAR(b.expiration_date, 'YYYY-MM-DD') AS expiration_date,
				COALESCE(mc.id, 0) AS medical_center_id,
				mc.name AS medical_center_name`).
			Joins("LEFT JOIN supply_code sc ON ms.code = sc.code").
			Joins("LEFT JOIN batch b ON ms.batch_id = b.id").
			Joins("LEFT JOIN pavilion p ON ms.location_id = p.id").
			Joins("LEFT JOIN medical_center mc ON p.medical_center_id = mc.id").
			Where("ms.status = ? AND ms.location_type = ? AND ms.location_id = ?",
				models.StatusEnRouteToPavilion,
				models.SupplyLocationPavilion,
				pavilionID).
			Group("ms.batch_id, ms.code, ms.location_id, p.name, sc.name, b.supplier, b.expiration_date, mc.id, mc.name").
			Scan(&inTransitRows).Error

		if err == nil {
			for _, row := range inTransitRows {
				item := models.PavilionInventorySummaryWithDetails{
					PavilionInventorySummary: models.PavilionInventorySummary{
						PavilionID:       row.PavilionID,
						BatchID:          row.BatchID,
						SupplyCode:       row.SupplyCode,
						CurrentAvailable: row.CurrentAvailable,
					},
					PavilionName:      row.PavilionName,
					SupplyName:        row.SupplyName,
					BatchSupplier:     row.BatchSupplier,
					ExpirationDate:    row.ExpirationDate,
					MedicalCenterID:   row.MedicalCenterID,
					MedicalCenterName: row.MedicalCenterName,
					InTransit:         true,
				}
				inventory = append(inventory, item)
			}
		}
	}

	return inventory, nil
}

// GetMovementHistory obtiene el historial de movimientos con filtros
func (s *InventoryService) GetMovementHistory(
	locationType *string,
	locationID *int,
	startDate *time.Time,
	endDate *time.Time,
	movementType *string, // 'entrada', 'salida', 'consumo', 'devolucion'
	page int,
	pageSize int,
) ([]models.SupplyHistoryWithDestination, int64, error) {
	var history []models.SupplyHistoryWithDestination
	var total int64

	query := s.DB.Table("supply_history sh").
		Select(`sh.*,
			CASE 
				WHEN sh.destination_type = 'store' THEN s.name
				WHEN sh.destination_type = 'pavilion' THEN p.name
			END as destination_name,
			CASE 
				WHEN sh.origin_type = 'store' THEN so.name
				WHEN sh.origin_type = 'pavilion' THEN po.name
			END as origin_name,
			mc.name as medical_center_name,
			u.name as user_name`).
		Joins("LEFT JOIN store s ON sh.destination_type = 'store' AND sh.destination_id = s.id").
		Joins("LEFT JOIN pavilion p ON sh.destination_type = 'pavilion' AND sh.destination_id = p.id").
		Joins("LEFT JOIN store so ON sh.origin_type = 'store' AND sh.origin_id = so.id").
		Joins("LEFT JOIN pavilion po ON sh.origin_type = 'pavilion' AND sh.origin_id = po.id").
		Joins("LEFT JOIN medical_center mc ON (s.medical_center_id = mc.id OR p.medical_center_id = mc.id OR so.medical_center_id = mc.id OR po.medical_center_id = mc.id)").
		Joins("LEFT JOIN \"user\" u ON sh.user_rut = u.rut")

	// Aplicar filtros
	if locationType != nil && locationID != nil {
		query = query.Where("(sh.destination_type = ? AND sh.destination_id = ?) OR (sh.origin_type = ? AND sh.origin_id = ?)",
			*locationType, *locationID, *locationType, *locationID)
	}
	if startDate != nil {
		query = query.Where("sh.date_time >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("sh.date_time <= ?", *endDate)
	}
	if movementType != nil {
		switch *movementType {
		case "entrada":
			query = query.Where("sh.destination_type = ? OR sh.status = ?", locationType, models.StatusReceived)
		case "salida":
			query = query.Where("sh.origin_type = ? OR sh.status = ?", locationType, models.StatusEnRouteToPavilion)
		case "consumo":
			query = query.Where("sh.status = ?", models.StatusConsumed)
		case "devolucion":
			query = query.Where("sh.origin_type = 'pavilion' AND sh.destination_type = 'store'")
		}
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	offset := (page - 1) * pageSize
	if err := query.Order("sh.date_time DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&history).Error; err != nil {
		return nil, 0, err
	}

	return history, total, nil
}

// SyncInventory sincroniza los contadores de inventario
func (s *InventoryService) SyncInventory() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var discrepancies []map[string]interface{}

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Sincronizar inventario de bodega
		var storeSummaries []models.StoreInventorySummary
		if err := tx.Find(&storeSummaries).Error; err != nil {
			return fmt.Errorf("error al obtener resúmenes de bodega: %v", err)
		}

		for _, summary := range storeSummaries {
			// Contar productos reales en la bodega
			var realCount int64
			tx.Model(&models.MedicalSupply{}).
				Where("batch_id = ? AND location_type = ? AND status != ?",
					summary.BatchID, models.SupplyLocationStore, models.StatusConsumed).
				Count(&realCount)

			if int(realCount) != summary.CurrentInStore {
				discrepancies = append(discrepancies, map[string]interface{}{
					"type":     "store",
					"batch_id": summary.BatchID,
					"store_id": summary.StoreID,
					"expected": summary.CurrentInStore,
					"real":     realCount,
				})

				// Corregir discrepancia
				summary.CurrentInStore = int(realCount)
				if err := tx.Save(&summary).Error; err != nil {
					return fmt.Errorf("error al actualizar resumen de bodega: %v", err)
				}
			}
		}

		// 2. Sincronizar inventario de pabellones
		var pavilionSummaries []models.PavilionInventorySummary
		if err := tx.Find(&pavilionSummaries).Error; err != nil {
			return fmt.Errorf("error al obtener resúmenes de pabellón: %v", err)
		}

		for _, summary := range pavilionSummaries {
			// Contar productos reales en el pabellón
			var realCount int64
			tx.Model(&models.MedicalSupply{}).
				Where("batch_id = ? AND location_type = ? AND location_id = ? AND status NOT IN ?",
					summary.BatchID, models.SupplyLocationPavilion, summary.PavilionID,
					[]string{models.StatusConsumed}).
				Count(&realCount)

			if int(realCount) != summary.CurrentAvailable {
				discrepancies = append(discrepancies, map[string]interface{}{
					"type":        "pavilion",
					"batch_id":    summary.BatchID,
					"pavilion_id": summary.PavilionID,
					"expected":    summary.CurrentAvailable,
					"real":        realCount,
				})

				// Corregir discrepancia
				summary.CurrentAvailable = int(realCount)
				if err := tx.Save(&summary).Error; err != nil {
					return fmt.Errorf("error al actualizar resumen de pabellón: %v", err)
				}
			}
		}

		// 3. Sincronizar cantidades de lotes
		var batches []models.Batch
		if err := tx.Find(&batches).Error; err != nil {
			return fmt.Errorf("error al obtener lotes: %v", err)
		}

		for _, batch := range batches {
			// Contar productos reales no consumidos del lote
			var realCount int64
			tx.Model(&models.MedicalSupply{}).
				Where("batch_id = ? AND status != ?", batch.ID, models.StatusConsumed).
				Count(&realCount)

			if int(realCount) != batch.Amount {
				discrepancies = append(discrepancies, map[string]interface{}{
					"type":     "batch",
					"batch_id": batch.ID,
					"expected": batch.Amount,
					"real":     realCount,
				})

				// Corregir discrepancia
				batch.Amount = int(realCount)
				if err := tx.Save(&batch).Error; err != nil {
					return fmt.Errorf("error al actualizar lote: %v", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	result["discrepancies_found"] = len(discrepancies)
	result["discrepancies"] = discrepancies
	result["sync_date"] = time.Now()
	result["status"] = "completed"

	return result, nil
}

// GetInventorySummary obtiene un resumen general del inventario
func (s *InventoryService) GetInventorySummary(medicalCenterID *int) (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	query := s.DB

	// Filtrar por centro médico si se especifica
	if medicalCenterID != nil {
		query = query.Where("medical_center_id = ?", *medicalCenterID)
	}

	// Total en bodegas
	var totalInStores int64
	query.Model(&models.StoreInventorySummary{}).
		Select("COALESCE(SUM(current_in_store), 0)").
		Scan(&totalInStores)
	summary["total_in_stores"] = totalInStores

	// Total en pabellones
	var totalInPavilions int64
	query.Model(&models.PavilionInventorySummary{}).
		Select("COALESCE(SUM(current_available), 0)").
		Scan(&totalInPavilions)
	summary["total_in_pavilions"] = totalInPavilions

	// Total transferido
	var totalTransferred int64
	query.Model(&models.StoreInventorySummary{}).
		Select("COALESCE(SUM(total_transferred_out), 0)").
		Scan(&totalTransferred)
	summary["total_transferred"] = totalTransferred

	// Total consumido
	var totalConsumed int64
	query.Model(&models.PavilionInventorySummary{}).
		Select("COALESCE(SUM(total_consumed), 0)").
		Scan(&totalConsumed)
	summary["total_consumed"] = totalConsumed

	// Stock bajo en bodegas
	var lowStockStores int64
	query.Model(&models.StoreInventorySummary{}).
		Where("current_in_store < (original_amount * 0.2)").
		Count(&lowStockStores)
	summary["low_stock_stores"] = lowStockStores

	// Productos próximos a vencer (90 días)
	var nearExpiration int64
	expirationDate := time.Now().AddDate(0, 0, 90)
	s.DB.Table("store_inventory_summary sis").
		Joins("LEFT JOIN batch b ON sis.batch_id = b.id").
		Where("b.expiration_date <= ?", expirationDate).
		Count(&nearExpiration)
	summary["near_expiration"] = nearExpiration

	// Transferencias pendientes
	var pendingTransfers int64
	query.Model(&models.SupplyTransfer{}).
		Where("status IN ?", []string{models.TransferStatusPending, models.TransferStatusInTransit}).
		Count(&pendingTransfers)
	summary["pending_transfers"] = pendingTransfers

	return summary, nil
}

// GetInventoryBySurgeryType obtiene inventario agrupado por tipo de cirugía
func (s *InventoryService) GetInventoryBySurgeryType(storeID *int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := s.DB.Table("store_inventory_summary sis").
		Select(`surg.id as surgery_id,
			surg.name as surgery_name,
			SUM(sis.current_in_store) as total_in_store,
			SUM(sis.total_transferred_out) as total_transferred,
			COUNT(DISTINCT sis.batch_id) as batch_count`).
		Joins("LEFT JOIN surgery surg ON sis.surgery_id = surg.id").
		Group("surg.id, surg.name")

	if storeID != nil {
		query = query.Where("sis.store_id = ?", *storeID)
	}

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// GetTransferReport obtiene un reporte de transferencias
func (s *InventoryService) GetTransferReport(
	startDate time.Time,
	endDate time.Time,
	groupBy string, // 'origin', 'destination', 'status', 'date'
) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := s.DB.Table("supply_transfer st")

	switch groupBy {
	case "origin":
		query = query.Select(`st.origin_type,
			st.origin_id,
			COUNT(*) as transfer_count,
			COUNT(CASE WHEN st.status = 'recibido' THEN 1 END) as completed_count`).
			Group("st.origin_type, st.origin_id")
	case "destination":
		query = query.Select(`st.destination_type,
			st.destination_id,
			COUNT(*) as transfer_count,
			COUNT(CASE WHEN st.status = 'recibido' THEN 1 END) as completed_count`).
			Group("st.destination_type, st.destination_id")
	case "status":
		query = query.Select(`st.status,
			COUNT(*) as transfer_count`).
			Group("st.status")
	case "date":
		query = query.Select(`DATE(st.send_date) as transfer_date,
			COUNT(*) as transfer_count,
			COUNT(CASE WHEN st.status = 'recibido' THEN 1 END) as completed_count`).
			Group("DATE(st.send_date)")
	default:
		return nil, fmt.Errorf("tipo de agrupación no válido: %s", groupBy)
	}

	// Usar >= y <= para incluir todo el día de endDate
	// Agregar 24 horas a endDate para incluir todo el día
	endDateInclusive := endDate.AddDate(0, 0, 1).Add(-time.Second)
	query = query.Where("st.send_date >= ? AND st.send_date <= ?", startDate, endDateInclusive)

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
