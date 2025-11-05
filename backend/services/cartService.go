package services

import (
	"errors"
	"fmt"
	"meditrack/models"
	"time"

	"gorm.io/gorm"
)

// CartService maneja la lógica de negocio de carritos
type CartService struct {
	DB *gorm.DB
}

// NewCartService crea una nueva instancia del servicio de carritos
func NewCartService(db *gorm.DB) *CartService {
	return &CartService{
		DB: db,
	}
}

// CreateCartForRequest crea un carrito automáticamente para una solicitud aprobada
func (s *CartService) CreateCartForRequest(supplyRequestID int, createdByRUT, createdByName string) (*models.SupplyCart, error) {
	// Verificar si ya existe un carrito para esta solicitud
	var existingCart models.SupplyCart
	if err := s.DB.Where("supply_request_id = ?", supplyRequestID).First(&existingCart).Error; err == nil {
		return &existingCart, nil
	}

	// Verificar que la solicitud existe
	var request models.SupplyRequest
	if err := s.DB.First(&request, supplyRequestID).Error; err != nil {
		return nil, fmt.Errorf("solicitud no encontrada: %w", err)
	}

	// Crear el carrito
	cart := &models.SupplyCart{
		SupplyRequestID: supplyRequestID,
		CartNumber:      models.GenerateCartNumber(),
		Status:          models.CartStatusActive,
		CreatedBy:       createdByRUT,
		CreatedByName:   createdByName,
	}

	if err := s.DB.Create(cart).Error; err != nil {
		return nil, fmt.Errorf("error al crear carrito: %w", err)
	}

	// Obtener todas las asignaciones QR de la solicitud y agregarlas al carrito
	var assignments []models.SupplyRequestQRAssignment
	if err := s.DB.Where("supply_request_id = ? AND status = ?", supplyRequestID, models.AssignmentStatusAssigned).
		Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("error al obtener asignaciones: %w", err)
	}

	// Agregar cada asignación al carrito
	for _, assignment := range assignments {
		cartItem := &models.SupplyCartItem{
			SupplyCartID:                cart.ID,
			SupplyRequestQRAssignmentID: assignment.ID,
			AddedBy:                     createdByRUT,
			AddedByName:                 createdByName,
			IsActive:                    true,
		}
		if err := s.DB.Create(cartItem).Error; err != nil {
			return nil, fmt.Errorf("error al agregar item al carrito: %w", err)
		}
	}

	// Cargar las relaciones
	if err := s.DB.Preload("Items.SupplyRequestQRAssignment.MedicalSupply").
		Preload("Items.SupplyRequestQRAssignment.SupplyRequestItem.SupplyCodeInfo").
		Preload("SupplyRequest").
		First(cart, cart.ID).Error; err != nil {
		return nil, fmt.Errorf("error al cargar carrito con relaciones: %w", err)
	}

	return cart, nil
}

// GetCartByRequestID obtiene el carrito asociado a una solicitud
func (s *CartService) GetCartByRequestID(supplyRequestID int) (*models.SupplyCart, error) {
	var cart models.SupplyCart
	if err := s.DB.Where("supply_request_id = ?", supplyRequestID).
		Preload("Items", "is_active = ?", true).
		Preload("Items.SupplyRequestQRAssignment.MedicalSupply").
		Preload("Items.SupplyRequestQRAssignment.SupplyRequestItem.SupplyCodeInfo").
		Preload("SupplyRequest").
		First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("carrito no encontrado para la solicitud")
		}
		return nil, fmt.Errorf("error al obtener carrito: %w", err)
	}
	return &cart, nil
}

// GetCartByID obtiene un carrito por su ID
func (s *CartService) GetCartByID(cartID int) (*models.SupplyCart, error) {
	var cart models.SupplyCart
	if err := s.DB.Preload("Items", "is_active = ?", true).
		Preload("Items.SupplyRequestQRAssignment.MedicalSupply").
		Preload("Items.SupplyRequestQRAssignment.SupplyRequestItem.SupplyCodeInfo").
		Preload("SupplyRequest").
		First(&cart, cartID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("carrito no encontrado")
		}
		return nil, fmt.Errorf("error al obtener carrito: %w", err)
	}
	return &cart, nil
}

// GetCartByQRCode obtiene el carrito asociado a un código QR
func (s *CartService) GetCartByQRCode(qrCode string) (*models.SupplyCart, error) {
	var assignment models.SupplyRequestQRAssignment
	if err := s.DB.Where("qr_code = ?", qrCode).First(&assignment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no se encontró asignación para el código QR")
		}
		return nil, fmt.Errorf("error al buscar asignación QR: %w", err)
	}

	return s.GetCartByRequestID(assignment.SupplyRequestID)
}

// AddItemToCart agrega un item (QR assignment) al carrito
func (s *CartService) AddItemToCart(cartID, assignmentID int, addedByRUT, addedByName string) (*models.SupplyCartItem, error) {
	// Verificar que el carrito existe y está activo
	cart, err := s.GetCartByID(cartID)
	if err != nil {
		return nil, err
	}

	if !cart.CanAddItems() {
		return nil, fmt.Errorf("el carrito no está activo, no se pueden agregar items")
	}

	// Verificar que la asignación existe
	var assignment models.SupplyRequestQRAssignment
	if err := s.DB.First(&assignment, assignmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("asignación QR no encontrada")
		}
		return nil, fmt.Errorf("error al buscar asignación: %w", err)
	}

	// Verificar si el item ya existe en el carrito
	var existingItem models.SupplyCartItem
	if err := s.DB.Where("supply_cart_id = ? AND supply_request_qr_assignment_id = ?", cartID, assignmentID).
		First(&existingItem).Error; err == nil {
		// El item ya existe
		if existingItem.IsActive {
			return &existingItem, nil
		}
		// Reactivar item
		existingItem.IsActive = true
		existingItem.RemovedAt = nil
		existingItem.RemovedBy = nil
		existingItem.RemovedByName = nil
		if err := s.DB.Save(&existingItem).Error; err != nil {
			return nil, fmt.Errorf("error al reactivar item: %w", err)
		}
		return &existingItem, nil
	}

	// Crear nuevo item
	cartItem := &models.SupplyCartItem{
		SupplyCartID:                cartID,
		SupplyRequestQRAssignmentID: assignmentID,
		AddedBy:                     addedByRUT,
		AddedByName:                 addedByName,
		IsActive:                    true,
	}

	if err := s.DB.Create(cartItem).Error; err != nil {
		return nil, fmt.Errorf("error al agregar item al carrito: %w", err)
	}

	// Cargar relaciones
	if err := s.DB.Preload("SupplyRequestQRAssignment.MedicalSupply").
		Preload("SupplyRequestQRAssignment.SupplyRequestItem.SupplyCodeInfo").
		First(cartItem, cartItem.ID).Error; err != nil {
		return nil, fmt.Errorf("error al cargar item con relaciones: %w", err)
	}

	return cartItem, nil
}

// RemoveItemFromCart marca un item como inactivo en el carrito
func (s *CartService) RemoveItemFromCart(cartID, itemID int, removedByRUT, removedByName string) error {
	var cartItem models.SupplyCartItem
	if err := s.DB.Where("id = ? AND supply_cart_id = ?", itemID, cartID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("item no encontrado en el carrito")
		}
		return fmt.Errorf("error al buscar item: %w", err)
	}

	if !cartItem.IsActive {
		return fmt.Errorf("el item ya está inactivo")
	}

	now := time.Now()
	cartItem.IsActive = false
	cartItem.RemovedAt = &now
	cartItem.RemovedBy = &removedByRUT
	cartItem.RemovedByName = &removedByName

	if err := s.DB.Save(&cartItem).Error; err != nil {
		return fmt.Errorf("error al remover item del carrito: %w", err)
	}

	return nil
}

// CloseCart cierra un carrito
func (s *CartService) CloseCart(cartID int, closedByRUT, closedByName string) error {
	cart, err := s.GetCartByID(cartID)
	if err != nil {
		return err
	}

	if cart.Status == models.CartStatusClosed {
		return fmt.Errorf("el carrito ya está cerrado")
	}

	now := time.Now()
	cart.Status = models.CartStatusClosed
	cart.ClosedAt = &now
	cart.ClosedBy = &closedByRUT
	cart.ClosedByName = &closedByName

	if err := s.DB.Save(cart).Error; err != nil {
		return fmt.Errorf("error al cerrar carrito: %w", err)
	}

	return nil
}

// GetAllCarts obtiene todos los carritos con paginación
func (s *CartService) GetAllCarts(page, pageSize int, status string) ([]models.SupplyCart, int64, error) {
	var carts []models.SupplyCart
	var total int64

	query := s.DB.Model(&models.SupplyCart{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error al contar carritos: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("SupplyRequest").
		Preload("Items", "is_active = ?", true).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&carts).Error; err != nil {
		return nil, 0, fmt.Errorf("error al obtener carritos: %w", err)
	}

	return carts, total, nil
}

// GetCartDetails obtiene los detalles completos del carrito usando la vista
func (s *CartService) GetCartDetails(cartID int) (*models.SupplyCartDetailView, error) {
	var details models.SupplyCartDetailView
	if err := s.DB.Where("cart_id = ?", cartID).First(&details).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("detalles del carrito no encontrados")
		}
		return nil, fmt.Errorf("error al obtener detalles: %w", err)
	}
	return &details, nil
}

// MarkItemAsUsed marca un item del carrito como utilizado (consumido)
func (s *CartService) MarkItemAsUsed(cartID, itemID int, userRUT, userName string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el carrito existe y está activo
		var cart models.SupplyCart
		if err := tx.First(&cart, cartID).Error; err != nil {
			return fmt.Errorf("carrito no encontrado: %w", err)
		}

		if !cart.CanAddItems() {
			return fmt.Errorf("el carrito no está activo")
		}

		// Verificar que el item existe y pertenece al carrito
		var cartItem models.SupplyCartItem
		if err := tx.Where("id = ? AND supply_cart_id = ?", itemID, cartID).
			Preload("SupplyRequestQRAssignment").
			First(&cartItem).Error; err != nil {
			return fmt.Errorf("item del carrito no encontrado: %w", err)
		}

		if !cartItem.IsActive {
			return fmt.Errorf("el item no está activo")
		}

		// Actualizar el estado de la asignación QR a "consumed"
		now := time.Now()
		if err := tx.Model(&models.SupplyRequestQRAssignment{}).
			Where("id = ?", cartItem.SupplyRequestQRAssignmentID).
			Updates(map[string]interface{}{
				"status":            models.AssignmentStatusConsumed,
				"delivered_date":    now,
				"delivered_by":      userRUT,
				"delivered_by_name": userName,
				"updated_at":        now,
			}).Error; err != nil {
			return fmt.Errorf("error actualizando asignación QR: %w", err)
		}

		// Obtener información completa del insumo
		var supply models.MedicalSupply
		if err := tx.First(&supply, cartItem.SupplyRequestQRAssignment.MedicalSupplyID).Error; err != nil {
			return fmt.Errorf("error obteniendo insumo: %w", err)
		}

		// Verificar que el insumo esté recibido en el pabellón antes de permitir uso
		if supply.Status != models.StatusReceived || supply.LocationType != models.SupplyLocationPavilion {
			return fmt.Errorf("el insumo debe estar recibido en el pabellón antes de ser utilizado. Estado actual: %s, Ubicación: %s", supply.Status, supply.LocationType)
		}

		// Obtener información del lote
		var batch models.Batch
		if err := tx.First(&batch, supply.BatchID).Error; err != nil {
			return fmt.Errorf("error obteniendo lote: %w", err)
		}

		// Obtener información del pabellón desde la solicitud
		pavilionID := 0
		var supplyRequest models.SupplyRequest
		if err := tx.First(&supplyRequest, cartItem.SupplyRequestQRAssignment.SupplyRequestID).Error; err == nil {
			pavilionID = supplyRequest.PavilionID
		}

		// Actualizar el insumo médico a estado "consumido"
		if err := tx.Model(&supply).
			Updates(map[string]interface{}{
				"status": models.StatusConsumed,
			}).Error; err != nil {
			return fmt.Errorf("error actualizando estado del insumo: %w", err)
		}

		// Actualizar cantidad del lote (restar 1)
		newAmount := batch.Amount - 1
		if err := tx.Model(&batch).Update("amount", newAmount).Error; err != nil {
			return fmt.Errorf("error actualizando cantidad del lote: %w", err)
		}

		// Actualizar inventarios según la ubicación del insumo
		if supply.LocationType == models.SupplyLocationStore {
			// Insumo consumido desde bodega - actualizar store_inventory_summary
			var storeSummary models.StoreInventorySummary
			if err := tx.Where("store_id = ? AND batch_id = ?", supply.LocationID, supply.BatchID).
				First(&storeSummary).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Si no existe el resumen, calcular el stock real en bodega
					var realCount int64
					tx.Model(&models.MedicalSupply{}).
						Where("batch_id = ? AND location_type = ? AND location_id = ? AND status != ?",
							supply.BatchID, models.SupplyLocationStore, supply.LocationID, models.StatusConsumed).
						Count(&realCount)

					// Crear resumen con valores calculados
					storeSummary = models.StoreInventorySummary{
						StoreID:              supply.LocationID,
						BatchID:              supply.BatchID,
						SupplyCode:           supply.Code,
						SurgeryID:            batch.SurgeryID,
						OriginalAmount:       int(realCount) + 1,
						CurrentInStore:       int(realCount),
						TotalConsumedInStore: 1,
						LastConsumedDate:     &now,
					}
					if err := tx.Create(&storeSummary).Error; err != nil {
						return fmt.Errorf("error creando resumen de bodega: %w", err)
					}
				} else {
					return fmt.Errorf("error obteniendo resumen de bodega: %w", err)
				}
			} else {
				// Actualizar resumen existente
				storeSummary.CurrentInStore--
				storeSummary.TotalConsumedInStore++
				storeSummary.LastConsumedDate = &now
				if err := tx.Save(&storeSummary).Error; err != nil {
					return fmt.Errorf("error actualizando resumen de bodega: %w", err)
				}
			}
		} else if supply.LocationType == models.SupplyLocationPavilion {
			// Insumo consumido desde pabellón - actualizar pavilion_inventory_summary
			var pavilionSummary models.PavilionInventorySummary
			if err := tx.Where("pavilion_id = ? AND batch_id = ?", supply.LocationID, supply.BatchID).
				First(&pavilionSummary).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Si no existe el resumen, crearlo
					pavilionSummary = models.PavilionInventorySummary{
						PavilionID:       supply.LocationID,
						BatchID:          supply.BatchID,
						SupplyCode:       supply.Code,
						TotalReceived:    1,
						CurrentAvailable: 0,
						TotalConsumed:    1,
						LastConsumedDate: &now,
					}
					if err := tx.Create(&pavilionSummary).Error; err != nil {
						return fmt.Errorf("error creando resumen de pabellón: %w", err)
					}
				} else {
					return fmt.Errorf("error obteniendo resumen de pabellón: %w", err)
				}
			} else {
				// Actualizar resumen existente
				pavilionSummary.CurrentAvailable--
				pavilionSummary.TotalConsumed++
				pavilionSummary.LastConsumedDate = &now
				if err := tx.Save(&pavilionSummary).Error; err != nil {
					return fmt.Errorf("error actualizando resumen de pabellón: %w", err)
				}
			}
		}

		// Registrar en historial
		history := models.SupplyHistory{
			MedicalSupplyID: supply.ID,
			DateTime:        now,
			Status:          "consumido",
			DestinationType: "pabellon",
			DestinationID:   pavilionID,
			UserRUT:         userRUT,
			Notes:           fmt.Sprintf("Insumo utilizado y marcado desde carrito %s", cart.CartNumber),
		}
		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("error registrando historial: %w", err)
		}

		// Agregar nota al item del carrito
		notes := fmt.Sprintf("Marcado como utilizado el %s por %s", now.Format("02/01/2006 15:04"), userName)
		if cartItem.Notes != "" {
			cartItem.Notes += "\n" + notes
		} else {
			cartItem.Notes = notes
		}
		if err := tx.Save(&cartItem).Error; err != nil {
			return fmt.Errorf("error actualizando notas del item: %w", err)
		}

		return nil
	})

	// Verificar si todos los items activos están procesados y cerrar automáticamente el carrito
	if err := s.checkAndAutoCloseCart(cartID, userRUT, userName); err != nil {
		// Log del error pero no fallar la operación
		fmt.Printf("Advertencia: Error al verificar cierre automático del carrito %d: %v\n", cartID, err)
	}

	return nil
}

// MarkItemForReturn marca un item del carrito para devolución
func (s *CartService) MarkItemForReturn(cartID, itemID int, userRUT, userName, reason string) error {
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el carrito existe y está activo
		var cart models.SupplyCart
		if err := tx.First(&cart, cartID).Error; err != nil {
			return fmt.Errorf("carrito no encontrado: %w", err)
		}

		if !cart.CanAddItems() {
			return fmt.Errorf("el carrito no está activo")
		}

		// Verificar que el item existe y pertenece al carrito
		var cartItem models.SupplyCartItem
		if err := tx.Where("id = ? AND supply_cart_id = ?", itemID, cartID).
			Preload("SupplyRequestQRAssignment").
			First(&cartItem).Error; err != nil {
			return fmt.Errorf("item del carrito no encontrado: %w", err)
		}

		if !cartItem.IsActive {
			return fmt.Errorf("el item no está activo")
		}

		// Actualizar el estado de la asignación QR a "returned"
		now := time.Now()
		if err := tx.Model(&models.SupplyRequestQRAssignment{}).
			Where("id = ?", cartItem.SupplyRequestQRAssignmentID).
			Updates(map[string]interface{}{
				"status":     models.AssignmentStatusReturned,
				"notes":      reason,
				"updated_at": now,
			}).Error; err != nil {
			return fmt.Errorf("error actualizando asignación QR: %w", err)
		}

		// Obtener información del insumo y su ubicación
		var supply models.MedicalSupply
		if err := tx.First(&supply, cartItem.SupplyRequestQRAssignment.MedicalSupplyID).Error; err != nil {
			return fmt.Errorf("error obteniendo insumo: %w", err)
		}

		// Verificar que el insumo esté recibido en el pabellón antes de permitir devolución
		if supply.Status != models.StatusReceived || supply.LocationType != models.SupplyLocationPavilion {
			return fmt.Errorf("el insumo debe estar recibido en el pabellón antes de ser devuelto. Estado actual: %s, Ubicación: %s", supply.Status, supply.LocationType)
		}

		// Obtener información del lote para determinar la bodega de destino
		var batch models.Batch
		if err := tx.First(&batch, supply.BatchID).Error; err != nil {
			return fmt.Errorf("error obteniendo lote: %w", err)
		}

		// Guardar ubicación anterior para el historial
		oldLocationType := supply.LocationType
		oldLocationID := supply.LocationID
		storeID := batch.StoreID

		// Crear transferencia de devolución (en tránsito de vuelta a bodega)
		transferCode := fmt.Sprintf("RETURN-CART-%d-%s", time.Now().Unix(), supply.QRCode[len(supply.QRCode)-5:])
		transfer := models.SupplyTransfer{
			TransferCode:    transferCode,
			QRCode:          supply.QRCode,
			MedicalSupplyID: supply.ID,
			OriginType:      models.TransferLocationPavilion,
			OriginID:        supply.LocationID,
			DestinationType: models.TransferLocationStore,
			DestinationID:   storeID,
			SentBy:          userRUT,
			SentByName:      userName,
			Status:          models.TransferStatusInTransit, // En tránsito, el bodeguero debe confirmar recepción
			TransferReason:  fmt.Sprintf("Devolución desde carrito %s", cart.CartNumber),
			SendDate:        now,
			Notes:           reason,
		}

		if err := tx.Create(&transfer).Error; err != nil {
			return fmt.Errorf("error al crear transferencia de devolución: %w", err)
		}

		// Marcar el insumo como en tránsito de vuelta a bodega
		if err := tx.Model(&supply).
			Updates(map[string]interface{}{
				"status":        models.StatusEnRouteToStore,
				"location_type": models.SupplyLocationStore,
				"location_id":   storeID,
				"in_transit":    true,
			}).Error; err != nil {
			return fmt.Errorf("error actualizando estado del insumo: %w", err)
		}

		// Actualizar inventarios
		// Decrementar el resumen del pabellón (ya que el insumo está en tránsito de vuelta)
		if oldLocationType == models.SupplyLocationPavilion && oldLocationID > 0 {
			var pavilionSummary models.PavilionInventorySummary
			if err := tx.Where("pavilion_id = ? AND batch_id = ?", oldLocationID, supply.BatchID).
				First(&pavilionSummary).Error; err == nil {
				// Solo actualizar si existe el resumen
				pavilionSummary.CurrentAvailable--
				pavilionSummary.TotalReturned++
				pavilionSummary.LastReturnedDate = &now
				if err := tx.Save(&pavilionSummary).Error; err != nil {
					return fmt.Errorf("error actualizando resumen de pabellón: %w", err)
				}
			}
		}

		// NO incrementar stock de bodega todavía - el insumo está en tránsito
		// El stock se incrementará cuando el bodeguero confirme la recepción

		// NO incrementar cantidad del lote todavía - se incrementará cuando se confirme la recepción

		// Registrar en historial como "en tránsito a bodega"
		originType := oldLocationType
		originID := oldLocationID
		history := models.SupplyHistory{
			MedicalSupplyID: supply.ID,
			DateTime:        now,
			Status:          models.StatusEnRouteToStore,
			DestinationType: models.DestinationTypeStore,
			DestinationID:   storeID,
			UserRUT:         userRUT,
			Notes:           fmt.Sprintf("Devuelto desde carrito %s (en tránsito). Motivo: %s. El bodeguero debe confirmar recepción.", cart.CartNumber, reason),
			OriginType:      &originType,
			OriginID:        &originID,
		}
		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("error registrando historial: %w", err)
		}

		// Agregar nota al item del carrito
		notes := fmt.Sprintf("Marcado para devolución el %s por %s. Motivo: %s", now.Format("02/01/2006 15:04"), userName, reason)
		if cartItem.Notes != "" {
			cartItem.Notes += "\n" + notes
		} else {
			cartItem.Notes = notes
		}
		if err := tx.Save(&cartItem).Error; err != nil {
			return fmt.Errorf("error actualizando notas del item: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Verificar si todos los items activos están procesados y cerrar automáticamente el carrito
	if err := s.checkAndAutoCloseCart(cartID, userRUT, userName); err != nil {
		// Log del error pero no fallar la operación
		fmt.Printf("Advertencia: Error al verificar cierre automático del carrito %d: %v\n", cartID, err)
	}

	return nil
}

// BatchOperationItem representa un item para operación múltiple
type BatchOperationItem struct {
	ItemID int    `json:"item_id"` // ID del item del carrito
	Action string `json:"action"`  // "use" o "return"
	Reason string `json:"reason"`  // Motivo (opcional, para devolución)
}

// BatchOperationResult representa el resultado de una operación múltiple
type BatchOperationResult struct {
	SuccessCount int      `json:"success_count"`
	ErrorCount   int      `json:"error_count"`
	Errors       []string `json:"errors,omitempty"`
	Processed    []int    `json:"processed"` // IDs de items procesados exitosamente
}

// BatchOperationItems procesa múltiples items del carrito en una sola operación
// Permite marcar algunos como usados y otros como devueltos en un solo paso
func (s *CartService) BatchOperationItems(cartID int, items []BatchOperationItem, userRUT, userName string) (*BatchOperationResult, error) {
	result := &BatchOperationResult{
		Processed: []int{},
		Errors:    []string{},
	}

	// Verificar que el carrito existe y está activo
	var cart models.SupplyCart
	if err := s.DB.First(&cart, cartID).Error; err != nil {
		return nil, fmt.Errorf("carrito no encontrado: %w", err)
	}

	if !cart.CanAddItems() {
		return nil, fmt.Errorf("el carrito no está activo")
	}

	// Validar que haya items para procesar
	if len(items) == 0 {
		return nil, fmt.Errorf("no hay items para procesar")
	}

	// Procesar cada item en una transacción
	for _, itemOp := range items {
		err := s.DB.Transaction(func(tx *gorm.DB) error {
			// Verificar que el item existe y pertenece al carrito
			var cartItem models.SupplyCartItem
			if err := tx.Where("id = ? AND supply_cart_id = ?", itemOp.ItemID, cartID).
				Preload("SupplyRequestQRAssignment").
				Preload("SupplyRequestQRAssignment.SupplyRequestItem").
				First(&cartItem).Error; err != nil {
				return fmt.Errorf("item %d no encontrado en el carrito", itemOp.ItemID)
			}

			if !cartItem.IsActive {
				return fmt.Errorf("el item %d no está activo", itemOp.ItemID)
			}

			// Verificar si ya fue procesado
			if cartItem.SupplyRequestQRAssignment.Status == models.AssignmentStatusConsumed ||
				cartItem.SupplyRequestQRAssignment.Status == models.AssignmentStatusReturned {
				return fmt.Errorf("el item %d ya fue procesado", itemOp.ItemID)
			}

			now := time.Now()

			switch itemOp.Action {
			case "use":
				// Marcar como usado (consumido)
				if err := tx.Model(&models.SupplyRequestQRAssignment{}).
					Where("id = ?", cartItem.SupplyRequestQRAssignmentID).
					Updates(map[string]interface{}{
						"status":            models.AssignmentStatusConsumed,
						"delivered_date":    now,
						"delivered_by":      userRUT,
						"delivered_by_name": userName,
						"updated_at":        now,
					}).Error; err != nil {
					return fmt.Errorf("error actualizando asignación QR: %w", err)
				}

				// Obtener información completa del insumo
				var supply models.MedicalSupply
				if err := tx.First(&supply, cartItem.SupplyRequestQRAssignment.MedicalSupplyID).Error; err != nil {
					return fmt.Errorf("error obteniendo insumo: %w", err)
				}

				// Verificar que el insumo esté recibido en el pabellón antes de permitir uso
				if supply.Status != models.StatusReceived || supply.LocationType != models.SupplyLocationPavilion {
					return fmt.Errorf("el insumo debe estar recibido en el pabellón antes de ser utilizado. Estado actual: %s, Ubicación: %s", supply.Status, supply.LocationType)
				}

				// Obtener información del lote
				var batch models.Batch
				if err := tx.First(&batch, supply.BatchID).Error; err != nil {
					return fmt.Errorf("error obteniendo lote: %w", err)
				}

				// Obtener información del pabellón desde la solicitud
				pavilionID := 0
				var supplyRequest models.SupplyRequest
				if err := tx.First(&supplyRequest, cartItem.SupplyRequestQRAssignment.SupplyRequestID).Error; err == nil {
					pavilionID = supplyRequest.PavilionID
				}

				// Actualizar el insumo médico a estado "consumido"
				if err := tx.Model(&supply).
					Updates(map[string]interface{}{
						"status": models.StatusConsumed,
					}).Error; err != nil {
					return fmt.Errorf("error actualizando estado del insumo: %w", err)
				}

				// Actualizar cantidad del lote (restar 1)
				newAmount := batch.Amount - 1
				if err := tx.Model(&batch).Update("amount", newAmount).Error; err != nil {
					return fmt.Errorf("error actualizando cantidad del lote: %w", err)
				}

				// Actualizar inventarios según la ubicación del insumo
				if supply.LocationType == models.SupplyLocationStore {
					// Insumo consumido desde bodega - actualizar store_inventory_summary
					var storeSummary models.StoreInventorySummary
					if err := tx.Where("store_id = ? AND batch_id = ?", supply.LocationID, supply.BatchID).
						First(&storeSummary).Error; err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							// Si no existe el resumen, calcular el stock real en bodega
							var realCount int64
							tx.Model(&models.MedicalSupply{}).
								Where("batch_id = ? AND location_type = ? AND location_id = ? AND status != ?",
									supply.BatchID, models.SupplyLocationStore, supply.LocationID, models.StatusConsumed).
								Count(&realCount)

							// Crear resumen con valores calculados
							storeSummary = models.StoreInventorySummary{
								StoreID:              supply.LocationID,
								BatchID:              supply.BatchID,
								SupplyCode:           supply.Code,
								SurgeryID:            batch.SurgeryID,
								OriginalAmount:       int(realCount) + 1,
								CurrentInStore:       int(realCount),
								TotalConsumedInStore: 1,
								LastConsumedDate:     &now,
							}
							if err := tx.Create(&storeSummary).Error; err != nil {
								return fmt.Errorf("error creando resumen de bodega: %w", err)
							}
						} else {
							return fmt.Errorf("error obteniendo resumen de bodega: %w", err)
						}
					} else {
						// Actualizar resumen existente
						storeSummary.CurrentInStore--
						storeSummary.TotalConsumedInStore++
						storeSummary.LastConsumedDate = &now
						if err := tx.Save(&storeSummary).Error; err != nil {
							return fmt.Errorf("error actualizando resumen de bodega: %w", err)
						}
					}
				} else if supply.LocationType == models.SupplyLocationPavilion {
					// Insumo consumido desde pabellón - actualizar pavilion_inventory_summary
					var pavilionSummary models.PavilionInventorySummary
					if err := tx.Where("pavilion_id = ? AND batch_id = ?", supply.LocationID, supply.BatchID).
						First(&pavilionSummary).Error; err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							// Si no existe el resumen, crearlo
							pavilionSummary = models.PavilionInventorySummary{
								PavilionID:       supply.LocationID,
								BatchID:          supply.BatchID,
								SupplyCode:       supply.Code,
								TotalReceived:    1,
								CurrentAvailable: 0,
								TotalConsumed:    1,
								LastConsumedDate: &now,
							}
							if err := tx.Create(&pavilionSummary).Error; err != nil {
								return fmt.Errorf("error creando resumen de pabellón: %w", err)
							}
						} else {
							return fmt.Errorf("error obteniendo resumen de pabellón: %w", err)
						}
					} else {
						// Actualizar resumen existente
						pavilionSummary.CurrentAvailable--
						pavilionSummary.TotalConsumed++
						pavilionSummary.LastConsumedDate = &now
						if err := tx.Save(&pavilionSummary).Error; err != nil {
							return fmt.Errorf("error actualizando resumen de pabellón: %w", err)
						}
					}
				}

				// Registrar en historial
				history := models.SupplyHistory{
					MedicalSupplyID: supply.ID,
					DateTime:        now,
					Status:          "consumido",
					DestinationType: "pabellon",
					DestinationID:   pavilionID,
					UserRUT:         userRUT,
					Notes:           fmt.Sprintf("Insumo utilizado y marcado desde carrito %s (operación múltiple)", cart.CartNumber),
				}
				if err := tx.Create(&history).Error; err != nil {
					return fmt.Errorf("error registrando historial: %w", err)
				}

				// Agregar nota al item del carrito
				notes := fmt.Sprintf("Marcado como utilizado el %s por %s (operación múltiple)", now.Format("02/01/2006 15:04"), userName)
				if cartItem.Notes != "" {
					cartItem.Notes += "\n" + notes
				} else {
					cartItem.Notes = notes
				}

			case "return":
				// Marcar para devolución
				reason := itemOp.Reason
				if reason == "" {
					reason = "Sin especificar"
				}

				if err := tx.Model(&models.SupplyRequestQRAssignment{}).
					Where("id = ?", cartItem.SupplyRequestQRAssignmentID).
					Updates(map[string]interface{}{
						"status":     models.AssignmentStatusReturned,
						"notes":      reason,
						"updated_at": now,
					}).Error; err != nil {
					return fmt.Errorf("error actualizando asignación QR: %w", err)
				}

				// Obtener información del insumo y su ubicación
				var supply models.MedicalSupply
				if err := tx.First(&supply, cartItem.SupplyRequestQRAssignment.MedicalSupplyID).Error; err != nil {
					return fmt.Errorf("error obteniendo insumo: %w", err)
				}

				// Verificar que el insumo esté recibido en el pabellón antes de permitir devolución
				if supply.Status != models.StatusReceived || supply.LocationType != models.SupplyLocationPavilion {
					return fmt.Errorf("el insumo debe estar recibido en el pabellón antes de ser devuelto. Estado actual: %s, Ubicación: %s", supply.Status, supply.LocationType)
				}

				// Obtener información del lote para determinar la bodega de destino
				var batch models.Batch
				if err := tx.First(&batch, supply.BatchID).Error; err != nil {
					return fmt.Errorf("error obteniendo lote: %w", err)
				}

				// Guardar ubicación anterior para el historial
				oldLocationType := supply.LocationType
				oldLocationID := supply.LocationID
				storeID := batch.StoreID

				// Crear transferencia de devolución (en tránsito de vuelta a bodega)
				transferCode := fmt.Sprintf("RETURN-CART-%d-%s", time.Now().Unix(), supply.QRCode[len(supply.QRCode)-5:])
				transfer := models.SupplyTransfer{
					TransferCode:    transferCode,
					QRCode:          supply.QRCode,
					MedicalSupplyID: supply.ID,
					OriginType:      models.TransferLocationPavilion,
					OriginID:        supply.LocationID,
					DestinationType: models.TransferLocationStore,
					DestinationID:   storeID,
					SentBy:          userRUT,
					SentByName:      userName,
					Status:          models.TransferStatusInTransit, // En tránsito, el bodeguero debe confirmar recepción
					TransferReason:  fmt.Sprintf("Devolución desde carrito %s", cart.CartNumber),
					SendDate:        now,
					Notes:           reason,
				}

				if err := tx.Create(&transfer).Error; err != nil {
					return fmt.Errorf("error al crear transferencia de devolución: %w", err)
				}

				// Marcar el insumo como en tránsito de vuelta a bodega
				if err := tx.Model(&supply).
					Updates(map[string]interface{}{
						"status":        models.StatusEnRouteToStore,
						"location_type": models.SupplyLocationStore,
						"location_id":   storeID,
						"in_transit":    true,
					}).Error; err != nil {
					return fmt.Errorf("error actualizando estado del insumo: %w", err)
				}

				// Actualizar inventarios
				// Decrementar el resumen del pabellón (ya que el insumo está en tránsito de vuelta)
				if oldLocationType == models.SupplyLocationPavilion && oldLocationID > 0 {
					var pavilionSummary models.PavilionInventorySummary
					if err := tx.Where("pavilion_id = ? AND batch_id = ?", oldLocationID, supply.BatchID).
						First(&pavilionSummary).Error; err == nil {
						// Solo actualizar si existe el resumen
						pavilionSummary.CurrentAvailable--
						pavilionSummary.TotalReturned++
						pavilionSummary.LastReturnedDate = &now
						if err := tx.Save(&pavilionSummary).Error; err != nil {
							return fmt.Errorf("error actualizando resumen de pabellón: %w", err)
						}
					}
				}

				// NO incrementar stock de bodega todavía - el insumo está en tránsito
				// El stock se incrementará cuando el bodeguero confirme la recepción

				// NO incrementar cantidad del lote todavía - se incrementará cuando se confirme la recepción

				// Registrar en historial como "en tránsito a bodega"
				originType := oldLocationType
				originID := oldLocationID
				history := models.SupplyHistory{
					MedicalSupplyID: supply.ID,
					DateTime:        now,
					Status:          models.StatusEnRouteToStore,
					DestinationType: models.DestinationTypeStore,
					DestinationID:   storeID,
					UserRUT:         userRUT,
					Notes:           fmt.Sprintf("Devuelto desde carrito %s (operación múltiple, en tránsito). Motivo: %s. El bodeguero debe confirmar recepción.", cart.CartNumber, reason),
					OriginType:      &originType,
					OriginID:        &originID,
				}
				if err := tx.Create(&history).Error; err != nil {
					return fmt.Errorf("error registrando historial: %w", err)
				}

				// Agregar nota al item del carrito
				notes := fmt.Sprintf("Marcado para devolución el %s por %s (operación múltiple). Motivo: %s", now.Format("02/01/2006 15:04"), userName, reason)
				if cartItem.Notes != "" {
					cartItem.Notes += "\n" + notes
				} else {
					cartItem.Notes = notes
				}

			default:
				return fmt.Errorf("acción inválida: %s (debe ser 'use' o 'return')", itemOp.Action)
			}

			// Guardar notas del item
			if err := tx.Save(&cartItem).Error; err != nil {
				return fmt.Errorf("error actualizando notas del item: %w", err)
			}

			return nil
		})

		if err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Item %d: %v", itemOp.ItemID, err))
		} else {
			result.SuccessCount++
			result.Processed = append(result.Processed, itemOp.ItemID)
		}
	}

	// Verificar si todos los items activos están procesados y cerrar automáticamente el carrito
	if err := s.checkAndAutoCloseCart(cartID, userRUT, userName); err != nil {
		// Log del error pero no fallar la operación
		fmt.Printf("Advertencia: Error al verificar cierre automático del carrito %d: %v\n", cartID, err)
	}

	return result, nil
}

// checkAndAutoCloseCart verifica si todos los items activos del carrito están procesados
// y cierra automáticamente el carrito si es así
func (s *CartService) checkAndAutoCloseCart(cartID int, closedByRUT, closedByName string) error {
	// Obtener el carrito con todos sus items activos
	var cart models.SupplyCart
	if err := s.DB.Where("id = ? AND status = ?", cartID, models.CartStatusActive).
		Preload("Items", "is_active = ?", true).
		Preload("Items.SupplyRequestQRAssignment").
		Preload("Items.SupplyRequestQRAssignment.MedicalSupply").
		First(&cart).Error; err != nil {
		return fmt.Errorf("carrito no encontrado o no activo: %w", err)
	}

	// Si no hay items activos, no hay nada que verificar
	if len(cart.Items) == 0 {
		return nil
	}

	// Verificar si todos los items activos están procesados (consumidos o devueltos)
	// Y si los items devueltos ya llegaron a bodega (no están en tránsito)
	allProcessed := true
	for _, item := range cart.Items {
		if !item.IsActive {
			continue
		}
		status := item.SupplyRequestQRAssignment.Status
		supply := item.SupplyRequestQRAssignment.MedicalSupply

		if status == models.AssignmentStatusConsumed {
			// Item consumido - está procesado
			continue
		} else if status == models.AssignmentStatusReturned {
			// Item devuelto - verificar si ya llegó a bodega (no está en tránsito)
			if supply.InTransit && supply.Status == models.StatusEnRouteToStore {
				// Item devuelto pero todavía en tránsito - NO cerrar el carrito todavía
				allProcessed = false
				fmt.Printf("Carrito %d no se cierra: item %d está devuelto pero en tránsito (QR: %s)\n",
					cartID, item.ID, supply.QRCode)
				break
			}
			// Item devuelto y ya recibido en bodega - está procesado
			continue
		} else {
			// Item no procesado
			allProcessed = false
			break
		}
	}

	// Si todos los items están procesados Y no hay devoluciones en tránsito, cerrar automáticamente el carrito
	if allProcessed {
		now := time.Now()
		cart.Status = models.CartStatusClosed
		cart.ClosedAt = &now
		cart.ClosedBy = &closedByRUT
		cart.ClosedByName = &closedByName
		cart.Notes = "Cerrado automáticamente: todos los items han sido procesados"

		if err := s.DB.Save(&cart).Error; err != nil {
			return fmt.Errorf("error al cerrar automáticamente el carrito: %w", err)
		}

		fmt.Printf("Carrito %d cerrado automáticamente: todos los items han sido procesados\n", cartID)
	}

	return nil
}

// CheckAndAutoCloseCartForSupply verifica y cierra automáticamente el carrito asociado a un insumo
// si todos los items del carrito están procesados (consumidos o devueltos y recibidos en bodega)
func (s *CartService) CheckAndAutoCloseCartForSupply(supplyID int, closedByRUT, closedByName string) error {
	// Buscar el assignment asociado al insumo
	var assignment models.SupplyRequestQRAssignment
	if err := s.DB.Where("medical_supply_id = ?", supplyID).First(&assignment).Error; err != nil {
		// No hay assignment - no hay carrito asociado
		return nil
	}

	// Buscar el item del carrito asociado a este assignment
	var cartItem models.SupplyCartItem
	if err := s.DB.Where("supply_request_qr_assignment_id = ? AND is_active = ?", assignment.ID, true).
		First(&cartItem).Error; err != nil {
		// No hay item activo del carrito - no hay carrito activo
		return nil
	}

	// Verificar y cerrar el carrito si corresponde
	return s.checkAndAutoCloseCart(cartItem.SupplyCartID, closedByRUT, closedByName)
}

// TransferCartToPavilion transfiere todos los items del carrito al pabellón
// Este método debe ser llamado por el bodeguero después de aprobar la solicitud
func (s *CartService) TransferCartToPavilion(cartID int, userRUT, userName string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el carrito existe y está activo
		var cart models.SupplyCart
		if err := tx.Preload("SupplyRequest").First(&cart, cartID).Error; err != nil {
			return fmt.Errorf("carrito no encontrado: %w", err)
		}

		if !cart.CanAddItems() {
			return fmt.Errorf("el carrito no está activo")
		}

		// Obtener todos los items activos del carrito
		var cartItems []models.SupplyCartItem
		if err := tx.Where("supply_cart_id = ? AND is_active = ?", cartID, true).
			Preload("SupplyRequestQRAssignment").
			Preload("SupplyRequestQRAssignment.MedicalSupply").
			Find(&cartItems).Error; err != nil {
			return fmt.Errorf("error al obtener items del carrito: %w", err)
		}

		if len(cartItems) == 0 {
			return fmt.Errorf("el carrito no tiene items para transferir")
		}

		// Obtener el pabellón de la solicitud
		pavilionID := cart.SupplyRequest.PavilionID
		if pavilionID == 0 {
			return fmt.Errorf("la solicitud no tiene pabellón asignado")
		}

		// Transferir cada insumo al pabellón
		for _, cartItem := range cartItems {
			supply := cartItem.SupplyRequestQRAssignment.MedicalSupply
			
			// Verificar que el insumo esté en bodega
			if supply.LocationType != models.SupplyLocationStore {
				return fmt.Errorf("el insumo %s no está en bodega (ubicación: %s)", supply.QRCode, supply.LocationType)
			}

			if supply.InTransit {
				return fmt.Errorf("el insumo %s ya está en tránsito", supply.QRCode)
			}

			if supply.Status == models.StatusConsumed {
				return fmt.Errorf("el insumo %s ya fue consumido", supply.QRCode)
			}

			// Guardar el store_id original antes de cambiar la ubicación
			originalStoreID := supply.LocationID
			
			// Si LocationID es 0, obtener el store_id del batch
			if originalStoreID == 0 {
				var batch models.Batch
				if err := tx.First(&batch, supply.BatchID).Error; err != nil {
					return fmt.Errorf("lote no encontrado para insumo %s: %w", supply.QRCode, err)
				}
				originalStoreID = batch.StoreID
				// Actualizar el supply con el LocationID correcto si estaba en 0
				if err := tx.Model(&supply).Updates(map[string]interface{}{
					"location_id":   originalStoreID,
					"location_type": models.SupplyLocationStore,
				}).Error; err != nil {
					return fmt.Errorf("error actualizando LocationID del insumo: %w", err)
				}
				supply.LocationID = originalStoreID
				supply.LocationType = models.SupplyLocationStore
			}

			// Obtener información del batch
			var batch models.Batch
			if err := tx.First(&batch, supply.BatchID).Error; err != nil {
				return fmt.Errorf("lote no encontrado para insumo %s: %w", supply.QRCode, err)
			}

			// Descontar del stock de bodega ANTES de actualizar la ubicación del insumo
			// El store_inventory_summary ya debería existir (se crea cuando se crea el lote)
			now := time.Now()
			var storeSummary models.StoreInventorySummary
			if err := tx.Where("batch_id = ?", batch.ID).First(&storeSummary).Error; err != nil {
				// Si no existe el resumen (no debería pasar), crearlo basado en el stock real actual
				if errors.Is(err, gorm.ErrRecordNotFound) {
					var realCount int64
					tx.Model(&models.MedicalSupply{}).
						Where("batch_id = ? AND location_type = ? AND location_id = ? AND status != ?",
							batch.ID, models.SupplyLocationStore, originalStoreID, models.StatusConsumed).
						Count(&realCount)

					storeSummary = models.StoreInventorySummary{
						StoreID:            originalStoreID,
						BatchID:            batch.ID,
						SupplyCode:         supply.Code,
						SurgeryID:          batch.SurgeryID,
						OriginalAmount:     int(realCount) + 1, // +1 porque vamos a transferir uno
						CurrentInStore:     int(realCount),     // Stock actual sin contar el que se transfiere
						TotalTransferredOut: 1,
						LastTransferOutDate: &now,
					}
					if err := tx.Create(&storeSummary).Error; err != nil {
						return fmt.Errorf("error creando resumen de bodega: %w", err)
					}
				} else {
					return fmt.Errorf("error obteniendo resumen de bodega: %w", err)
				}
			} else {
				// Actualizar resumen existente
				if storeSummary.CurrentInStore > 0 {
					storeSummary.CurrentInStore--
				}
				storeSummary.TotalTransferredOut++
				storeSummary.LastTransferOutDate = &now
				if err := tx.Save(&storeSummary).Error; err != nil {
					return fmt.Errorf("error actualizando resumen de bodega: %w", err)
				}
			}

			// Crear registro de transferencia
			transferCode := fmt.Sprintf("TRANS-CART-%d-%s", time.Now().Unix(), supply.QRCode[len(supply.QRCode)-5:])
			transfer := models.SupplyTransfer{
				TransferCode:    transferCode,
				QRCode:          supply.QRCode,
				MedicalSupplyID: supply.ID,
				OriginType:      models.TransferLocationStore,
				OriginID:        originalStoreID,
				DestinationType: models.TransferLocationPavilion,
				DestinationID:   pavilionID,
				SentBy:          userRUT,
				SentByName:      userName,
				Status:          models.TransferStatusInTransit,
				TransferReason:  fmt.Sprintf("Transferencia desde carrito %s", cart.CartNumber),
				SendDate:        now,
				Notes:           fmt.Sprintf("Transferencia automática desde carrito de solicitud %s", cart.SupplyRequest.RequestNumber),
			}

			if err := tx.Create(&transfer).Error; err != nil {
				return fmt.Errorf("error al crear transferencia para insumo %s: %w", supply.QRCode, err)
			}

			// Actualizar ubicación del insumo y marcarlo en tránsito
			supply.LocationType = models.SupplyLocationPavilion
			supply.LocationID = pavilionID
			supply.InTransit = true
			supply.TransferDate = &now
			supply.TransferredBy = &userRUT
			supply.Status = models.StatusEnRouteToPavilion

			if err := tx.Save(&supply).Error; err != nil {
				return fmt.Errorf("error al actualizar insumo %s: %w", supply.QRCode, err)
			}

			// Registrar en historial
			history := models.SupplyHistory{
				MedicalSupplyID: supply.ID,
				DateTime:        now,
				Status:          models.StatusEnRouteToPavilion,
				DestinationType: models.DestinationTypePavilion,
				DestinationID:   pavilionID,
				UserRUT:         userRUT,
				Notes:           fmt.Sprintf("Transferido desde carrito %s al pabellón. El pabellón debe confirmar recepción.", cart.CartNumber),
			}
			if err := tx.Create(&history).Error; err != nil {
				return fmt.Errorf("error registrando historial: %w", err)
			}
		}

		return nil
	})
}