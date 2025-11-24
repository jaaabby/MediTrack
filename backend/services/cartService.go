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

const (
	DefaultReturnReason = "Sin especificar"
)

type CartService struct {
	DB *gorm.DB
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{DB: db}
}

// ========================
// CRUD BÁSICO
// ========================

func (s *CartService) CreateCartForRequest(supplyRequestID int, createdByRUT, createdByName string) (*models.SupplyCart, error) {
	// Verificar si ya existe un carrito para esta solicitud
	var existingCart models.SupplyCart
	if err := s.DB.Where("supply_request_id = ?", supplyRequestID).First(&existingCart).Error; err == nil {
		// Si ya existe, cargarlo con sus relaciones y retornarlo
		if err := s.DB.Preload("Items.SupplyRequestQRAssignment.MedicalSupply").
			Preload("Items.SupplyRequestQRAssignment.SupplyRequestItem.SupplyCodeInfo").
			Preload("SupplyRequest").
			First(&existingCart, existingCart.ID).Error; err != nil {
			return nil, fmt.Errorf("error cargando carrito existente: %w", err)
		}
		return &existingCart, nil
	}

	var request models.SupplyRequest
	if err := s.DB.First(&request, supplyRequestID).Error; err != nil {
		return nil, fmt.Errorf("solicitud no encontrada: %w", err)
	}

	cart := &models.SupplyCart{
		SupplyRequestID: supplyRequestID,
		CartNumber:      models.GenerateCartNumber(),
		Status:          models.CartStatusActive,
		CreatedBy:       createdByRUT,
		CreatedByName:   createdByName,
	}

	// Usar la misma transacción si s.DB ya es una transacción, o crear una nueva si no
	// Verificar si s.DB es una transacción verificando si tiene un método específico
	// Por ahora, usaremos s.DB directamente ya que si se pasa una transacción, debería funcionar
	return cart, s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(cart).Error; err != nil {
			return fmt.Errorf("error creando carrito: %w", err)
		}

		var assignments []models.SupplyRequestQRAssignment
		if err := tx.Where("supply_request_id = ? AND status = ?", supplyRequestID, models.AssignmentStatusAssigned).
			Find(&assignments).Error; err != nil {
			return fmt.Errorf("error obteniendo asignaciones: %w", err)
		}

		for _, assignment := range assignments {
			cartItem := &models.SupplyCartItem{
				SupplyCartID:                cart.ID,
				SupplyRequestQRAssignmentID: assignment.ID,
				AddedBy:                     createdByRUT,
				AddedByName:                 createdByName,
				IsActive:                    true,
			}
			if err := tx.Create(cartItem).Error; err != nil {
				return fmt.Errorf("error agregando item: %w", err)
			}
		}

		return tx.Preload("Items.SupplyRequestQRAssignment.MedicalSupply").
			Preload("Items.SupplyRequestQRAssignment.SupplyRequestItem.SupplyCodeInfo").
			Preload("SupplyRequest").
			First(cart, cart.ID).Error
	})
}

func (s *CartService) GetCartByRequestID(supplyRequestID int) (*models.SupplyCart, error) {
	var cart models.SupplyCart
	err := s.DB.Where("supply_request_id = ?", supplyRequestID).
		Preload("Items", "is_active = ?", true).
		Preload("Items.SupplyRequestQRAssignment.MedicalSupply").
		Preload("Items.SupplyRequestQRAssignment.SupplyRequestItem.SupplyCodeInfo").
		Preload("SupplyRequest").
		First(&cart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("carrito no encontrado para la solicitud")
		}
		return nil, fmt.Errorf("error obteniendo carrito: %w", err)
	}
	return &cart, nil
}

func (s *CartService) GetCartByID(cartID int) (*models.SupplyCart, error) {
	var cart models.SupplyCart
	err := s.DB.Preload("Items", "is_active = ?", true).
		Preload("Items.SupplyRequestQRAssignment.MedicalSupply").
		Preload("Items.SupplyRequestQRAssignment.SupplyRequestItem.SupplyCodeInfo").
		Preload("SupplyRequest").
		First(&cart, cartID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("carrito no encontrado")
		}
		return nil, fmt.Errorf("error obteniendo carrito: %w", err)
	}
	return &cart, nil
}

func (s *CartService) GetCartByQRCode(qrCode string) (*models.SupplyCart, error) {
	var assignment models.SupplyRequestQRAssignment
	if err := s.DB.Where("qr_code = ?", qrCode).First(&assignment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no se encontró asignación para el código QR")
		}
		return nil, fmt.Errorf("error buscando asignación QR: %w", err)
	}
	return s.GetCartByRequestID(assignment.SupplyRequestID)
}

func (s *CartService) GetCartDetails(cartID int) (*models.SupplyCartDetailView, error) {
	var details models.SupplyCartDetailView
	if err := s.DB.Where("cart_id = ?", cartID).First(&details).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("detalles del carrito no encontrados")
		}
		return nil, fmt.Errorf("error obteniendo detalles: %w", err)
	}
	return &details, nil
}

func (s *CartService) GetAllCarts(page, pageSize int, status string) ([]models.SupplyCart, int64, error) {
	var carts []models.SupplyCart
	var total int64

	query := s.DB.Model(&models.SupplyCart{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error contando carritos: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("SupplyRequest").
		Preload("Items", "is_active = ?", true).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&carts).Error; err != nil {
		return nil, 0, fmt.Errorf("error obteniendo carritos: %w", err)
	}

	return carts, total, nil
}

// ========================
// GESTIÓN DE ITEMS
// ========================

func (s *CartService) AddItemToCart(cartID, assignmentID int, addedByRUT, addedByName string) (*models.SupplyCartItem, error) {
	cart, err := s.GetCartByID(cartID)
	if err != nil {
		return nil, err
	}

	if !cart.CanAddItems() {
		return nil, fmt.Errorf("el carrito no está activo")
	}

	var assignment models.SupplyRequestQRAssignment
	if err := s.DB.First(&assignment, assignmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("asignación QR no encontrada")
		}
		return nil, fmt.Errorf("error buscando asignación: %w", err)
	}

	var existingItem models.SupplyCartItem
	if err := s.DB.Where("supply_cart_id = ? AND supply_request_qr_assignment_id = ?", cartID, assignmentID).
		First(&existingItem).Error; err == nil {
		if existingItem.IsActive {
			return &existingItem, nil
		}
		existingItem.IsActive = true
		existingItem.RemovedAt = nil
		existingItem.RemovedBy = nil
		existingItem.RemovedByName = nil
		if err := s.DB.Save(&existingItem).Error; err != nil {
			return nil, fmt.Errorf("error reactivando item: %w", err)
		}
		return &existingItem, nil
	}

	cartItem := &models.SupplyCartItem{
		SupplyCartID:                cartID,
		SupplyRequestQRAssignmentID: assignmentID,
		AddedBy:                     addedByRUT,
		AddedByName:                 addedByName,
		IsActive:                    true,
	}

	if err := s.DB.Create(cartItem).Error; err != nil {
		return nil, fmt.Errorf("error agregando item: %w", err)
	}

	if err := s.DB.Preload("SupplyRequestQRAssignment.MedicalSupply").
		Preload("SupplyRequestQRAssignment.SupplyRequestItem.SupplyCodeInfo").
		First(cartItem, cartItem.ID).Error; err != nil {
		return nil, fmt.Errorf("error cargando item: %w", err)
	}

	return cartItem, nil
}

func (s *CartService) RemoveItemFromCart(cartID, itemID int, removedByRUT, removedByName string) error {
	var cartItem models.SupplyCartItem
	if err := s.DB.Where("id = ? AND supply_cart_id = ?", itemID, cartID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("item no encontrado en el carrito")
		}
		return fmt.Errorf("error buscando item: %w", err)
	}

	if !cartItem.IsActive {
		return fmt.Errorf("el item ya está inactivo")
	}

	now := time.Now()
	cartItem.IsActive = false
	cartItem.RemovedAt = &now
	cartItem.RemovedBy = &removedByRUT
	cartItem.RemovedByName = &removedByName

	return s.DB.Save(&cartItem).Error
}

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

	return s.DB.Save(cart).Error
}

// ========================
// OPERACIONES DE ITEMS
// ========================

func (s *CartService) MarkItemAsUsed(cartID, itemID int, userRUT, userName string) error {
	return s.processCartItem(cartID, itemID, userRUT, userName, "use", "")
}

func (s *CartService) MarkItemForReturn(cartID, itemID int, userRUT, userName, reason string) error {
	if reason == "" {
		reason = DefaultReturnReason
	}
	return s.processCartItem(cartID, itemID, userRUT, userName, "return", reason)
}

// processCartItem es el método centralizado para procesar items (uso o devolución)
func (s *CartService) processCartItem(cartID, itemID int, userRUT, userName, action, reason string) error {
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		cart, cartItem, supply, batch, err := s.validateCartItem(tx, cartID, itemID)
		if err != nil {
			return err
		}

		now := time.Now()

		if action == "use" {
			return s.processItemUse(tx, cart, cartItem, supply, batch, userRUT, userName, now)
		} else if action == "return" {
			return s.processItemReturn(tx, cart, cartItem, supply, batch, userRUT, userName, reason, now)
		}

		return fmt.Errorf("acción inválida: %s", action)
	})

	if err != nil {
		return err
	}

	return s.checkAndAutoCloseCart(cartID, userRUT, userName)
}

// validateCartItem valida que el cart item exista y esté en estado correcto
func (s *CartService) validateCartItem(tx *gorm.DB, cartID, itemID int) (*models.SupplyCart, *models.SupplyCartItem, *models.MedicalSupply, *models.Batch, error) {
	var cart models.SupplyCart
	if err := tx.First(&cart, cartID).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("carrito no encontrado: %w", err)
	}

	if !cart.CanAddItems() {
		return nil, nil, nil, nil, fmt.Errorf("el carrito no está activo")
	}

	var cartItem models.SupplyCartItem
	if err := tx.Where("id = ? AND supply_cart_id = ?", itemID, cartID).
		Preload("SupplyRequestQRAssignment").
		First(&cartItem).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("item no encontrado: %w", err)
	}

	if !cartItem.IsActive {
		return nil, nil, nil, nil, fmt.Errorf("el item no está activo")
	}

	var supply models.MedicalSupply
	if err := tx.First(&supply, cartItem.SupplyRequestQRAssignment.MedicalSupplyID).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error obteniendo insumo: %w", err)
	}

	if supply.Status != models.StatusReceived || supply.LocationType != models.SupplyLocationPavilion {
		return nil, nil, nil, nil, fmt.Errorf("el insumo debe estar recibido en el pabellón. Estado: %s, Ubicación: %s", supply.Status, supply.LocationType)
	}

	var batch models.Batch
	if err := tx.First(&batch, supply.BatchID).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error obteniendo lote: %w", err)
	}

	return &cart, &cartItem, &supply, &batch, nil
}

// processItemUse procesa el uso/consumo de un item
func (s *CartService) processItemUse(tx *gorm.DB, cart *models.SupplyCart, cartItem *models.SupplyCartItem, supply *models.MedicalSupply, batch *models.Batch, userRUT, userName string, now time.Time) error {
	// Actualizar asignación QR
	if err := tx.Model(&models.SupplyRequestQRAssignment{}).
		Where("id = ?", cartItem.SupplyRequestQRAssignmentID).
		Updates(map[string]interface{}{
			"status":            models.AssignmentStatusConsumed,
			"delivered_date":    now,
			"delivered_by":      userRUT,
			"delivered_by_name": userName,
			"updated_at":        now,
		}).Error; err != nil {
		return fmt.Errorf("error actualizando asignación: %w", err)
	}

	// Actualizar insumo
	if err := tx.Model(supply).Update("status", models.StatusConsumed).Error; err != nil {
		return fmt.Errorf("error actualizando insumo: %w", err)
	}

	// Actualizar lote
	if err := tx.Model(batch).Update("amount", batch.Amount-1).Error; err != nil {
		return fmt.Errorf("error actualizando lote: %w", err)
	}

	// Actualizar inventarios
	if err := s.updateInventoryOnConsumption(tx, supply, batch, now); err != nil {
		return err
	}

	// Obtener pavilion ID
	pavilionID := 0
	var supplyRequest models.SupplyRequest
	if err := tx.First(&supplyRequest, cartItem.SupplyRequestQRAssignment.SupplyRequestID).Error; err == nil {
		pavilionID = supplyRequest.PavilionID
	}

	// Registrar historial
	history := models.SupplyHistory{
		MedicalSupplyID: supply.ID,
		DateTime:        now,
		Status:          "consumido",
		DestinationType: "pabellon",
		DestinationID:   pavilionID,
		UserRUT:         userRUT,
		Notes:           fmt.Sprintf("Insumo utilizado desde carrito %s", cart.CartNumber),
	}
	if err := tx.Create(&history).Error; err != nil {
		return fmt.Errorf("error registrando historial: %w", err)
	}

	// Actualizar notas del item
	notes := fmt.Sprintf("Marcado como utilizado el %s por %s", now.Format("02/01/2006 15:04"), userName)
	if cartItem.Notes != "" {
		cartItem.Notes += "\n" + notes
	} else {
		cartItem.Notes = notes
	}
	return tx.Save(cartItem).Error
}

// processItemReturn procesa la devolución de un item
func (s *CartService) processItemReturn(tx *gorm.DB, cart *models.SupplyCart, cartItem *models.SupplyCartItem, supply *models.MedicalSupply, batch *models.Batch, userRUT, userName, reason string, now time.Time) error {
	// Verificar si el insumo fue consumido automáticamente
	wasAutoConsumed := false
	if supply.Status == models.StatusConsumed {
		var lastConsumptionHistory models.SupplyHistory
		if err := tx.Where("medical_supply_id = ? AND status = ?", supply.ID, models.StatusConsumed).
			Order("date_time DESC").
			First(&lastConsumptionHistory).Error; err == nil {
			if strings.Contains(lastConsumptionHistory.Notes, "[CONSUMO_AUTOMATICO]") {
				wasAutoConsumed = true
				log.Printf("🔄 Devolviendo insumo %s consumido automáticamente desde carrito", supply.QRCode)
			}
		}
	}

	// Actualizar asignación
	assignmentStatus := models.AssignmentStatusReturned
	if wasAutoConsumed {
		// Si fue consumido automáticamente, cambiar el estado a devuelto
		assignmentStatus = models.AssignmentStatusReturned
	}
	
	if err := tx.Model(&models.SupplyRequestQRAssignment{}).
		Where("id = ?", cartItem.SupplyRequestQRAssignmentID).
		Updates(map[string]interface{}{
			"status":     assignmentStatus,
			"notes":      reason,
			"updated_at": now,
		}).Error; err != nil {
		return fmt.Errorf("error actualizando asignación: %w", err)
	}

	// Guardar ubicación anterior
	oldLocationType := supply.LocationType
	oldLocationID := supply.LocationID
	storeID := batch.StoreID

	// Crear transferencia
	transferCode := fmt.Sprintf("RETURN-CART-%d-%s", time.Now().Unix(), supply.QRCode[len(supply.QRCode)-5:])
	transfer := models.SupplyTransfer{
		TransferCode:    transferCode,
		QRCode:          supply.QRCode,
		MedicalSupplyID: supply.ID,
		OriginType:      models.TransferLocationPavilion,
		OriginID:        oldLocationID,
		DestinationType: models.TransferLocationStore,
		DestinationID:   storeID,
		SentBy:          userRUT,
		SentByName:      userName,
		Status:          models.TransferStatusInTransit,
		TransferReason:  fmt.Sprintf("Devolución desde carrito %s", cart.CartNumber),
		SendDate:        now,
		Notes:           reason,
	}
	if err := tx.Create(&transfer).Error; err != nil {
		return fmt.Errorf("error creando transferencia: %w", err)
	}

	// Actualizar insumo
	if err := tx.Model(supply).Updates(map[string]interface{}{
		"status":        models.StatusEnRouteToStore,
		"location_type": models.SupplyLocationStore,
		"location_id":   storeID,
		"in_transit":    true,
	}).Error; err != nil {
		return fmt.Errorf("error actualizando insumo: %w", err)
	}

	// Actualizar inventario de pabellón (solo si no estaba consumido automáticamente)
	if !wasAutoConsumed {
		if err := s.updatePavilionInventoryOnReturn(tx, oldLocationID, supply.BatchID, now); err != nil {
			return err
		}
	} else {
		// Si estaba consumido automáticamente, solo actualizar contadores de devolución
		var pavilionSummary models.PavilionInventorySummary
		if err := tx.Where("pavilion_id = ? AND batch_id = ?", oldLocationID, supply.BatchID).
			First(&pavilionSummary).Error; err == nil {
			pavilionSummary.TotalReturned++
			pavilionSummary.LastReturnedDate = &now
			tx.Save(&pavilionSummary)
		}
	}

	// Registrar historial
	originType := oldLocationType
	originID := oldLocationID
	history := models.SupplyHistory{
		MedicalSupplyID: supply.ID,
		DateTime:        now,
		Status:          models.StatusEnRouteToStore,
		DestinationType: models.DestinationTypeStore,
		DestinationID:   storeID,
		UserRUT:         userRUT,
		Notes:           fmt.Sprintf("Devuelto desde carrito %s (en tránsito). Motivo: %s", cart.CartNumber, reason),
		OriginType:      &originType,
		OriginID:        &originID,
	}
	if err := tx.Create(&history).Error; err != nil {
		return fmt.Errorf("error registrando historial: %w", err)
	}

	// Actualizar notas del item
	notes := fmt.Sprintf("Marcado para devolución el %s por %s. Motivo: %s", now.Format("02/01/2006 15:04"), userName, reason)
	if cartItem.Notes != "" {
		cartItem.Notes += "\n" + notes
	} else {
		cartItem.Notes = notes
	}
	return tx.Save(cartItem).Error
}

// updateInventoryOnConsumption actualiza inventarios al consumir
func (s *CartService) updateInventoryOnConsumption(tx *gorm.DB, supply *models.MedicalSupply, batch *models.Batch, now time.Time) error {
	if supply.LocationType == models.SupplyLocationPavilion {
		var pavilionSummary models.PavilionInventorySummary
		if err := tx.Where("pavilion_id = ? AND batch_id = ?", supply.LocationID, batch.ID).
			First(&pavilionSummary).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				pavilionSummary = models.PavilionInventorySummary{
					PavilionID:       supply.LocationID,
					BatchID:          batch.ID,
					SupplyCode:       supply.Code,
					TotalReceived:    1,
					CurrentAvailable: 0,
					TotalConsumed:    1,
					LastConsumedDate: &now,
				}
				return tx.Create(&pavilionSummary).Error
			}
			return fmt.Errorf("error obteniendo resumen de pabellón: %w", err)
		}

		pavilionSummary.CurrentAvailable--
		pavilionSummary.TotalConsumed++
		pavilionSummary.LastConsumedDate = &now
		return tx.Save(&pavilionSummary).Error
	}
	return nil
}

// updatePavilionInventoryOnReturn actualiza inventario de pabellón al devolver
func (s *CartService) updatePavilionInventoryOnReturn(tx *gorm.DB, pavilionID, batchID int, now time.Time) error {
	if pavilionID > 0 {
		var pavilionSummary models.PavilionInventorySummary
		if err := tx.Where("pavilion_id = ? AND batch_id = ?", pavilionID, batchID).
			First(&pavilionSummary).Error; err == nil {
			pavilionSummary.CurrentAvailable--
			pavilionSummary.TotalReturned++
			pavilionSummary.LastReturnedDate = &now
			return tx.Save(&pavilionSummary).Error
		}
	}
	return nil
}

// ========================
// OPERACIONES MÚLTIPLES
// ========================

type BatchOperationItem struct {
	ItemID int    `json:"item_id"`
	Action string `json:"action"`
	Reason string `json:"reason"`
}

type BatchOperationResult struct {
	SuccessCount int      `json:"success_count"`
	ErrorCount   int      `json:"error_count"`
	Errors       []string `json:"errors,omitempty"`
	Processed    []int    `json:"processed"`
}

func (s *CartService) BatchOperationItems(cartID int, items []BatchOperationItem, userRUT, userName string) (*BatchOperationResult, error) {
	result := &BatchOperationResult{
		Processed: []int{},
		Errors:    []string{},
	}

	var cart models.SupplyCart
	if err := s.DB.First(&cart, cartID).Error; err != nil {
		return nil, fmt.Errorf("carrito no encontrado: %w", err)
	}

	if !cart.CanAddItems() {
		return nil, fmt.Errorf("el carrito no está activo")
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no hay items para procesar")
	}

	for _, itemOp := range items {
		err := s.processCartItem(cartID, itemOp.ItemID, userRUT, userName, itemOp.Action, itemOp.Reason)
		if err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Item %d: %v", itemOp.ItemID, err))
		} else {
			result.SuccessCount++
			result.Processed = append(result.Processed, itemOp.ItemID)
		}
	}

	return result, nil
}

// ========================
// HELPERS
// ========================

func (s *CartService) checkAndAutoCloseCart(cartID int, closedByRUT, closedByName string) error {
	var cart models.SupplyCart
	if err := s.DB.Where("id = ? AND status = ?", cartID, models.CartStatusActive).
		Preload("Items", "is_active = ?", true).
		Preload("Items.SupplyRequestQRAssignment").
		Preload("Items.SupplyRequestQRAssignment.MedicalSupply").
		First(&cart).Error; err != nil {
		return nil
	}

	if len(cart.Items) == 0 {
		return nil
	}

	allProcessed := true
	for _, item := range cart.Items {
		if !item.IsActive {
			continue
		}
		status := item.SupplyRequestQRAssignment.Status
		supply := item.SupplyRequestQRAssignment.MedicalSupply

		if status == models.AssignmentStatusConsumed {
			continue
		} else if status == models.AssignmentStatusReturned {
			if supply.InTransit && supply.Status == models.StatusEnRouteToStore {
				allProcessed = false
				break
			}
			continue
		} else {
			allProcessed = false
			break
		}
	}

	if allProcessed {
		now := time.Now()
		cart.Status = models.CartStatusClosed
		cart.ClosedAt = &now
		cart.ClosedBy = &closedByRUT
		cart.ClosedByName = &closedByName
		cart.Notes = "Cerrado automáticamente: todos los items procesados"
		return s.DB.Save(&cart).Error
	}

	return nil
}

func (s *CartService) CheckAndAutoCloseCartForSupply(supplyID int, closedByRUT, closedByName string) error {
	var assignment models.SupplyRequestQRAssignment
	if err := s.DB.Where("medical_supply_id = ?", supplyID).First(&assignment).Error; err != nil {
		return nil
	}

	var cartItem models.SupplyCartItem
	if err := s.DB.Where("supply_request_qr_assignment_id = ? AND is_active = ?", assignment.ID, true).
		First(&cartItem).Error; err != nil {
		return nil
	}

	return s.checkAndAutoCloseCart(cartItem.SupplyCartID, closedByRUT, closedByName)
}

// TransferCartToPavilion - MANTENIDO SIN CAMBIOS (no hay duplicación aquí)
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
						StoreID:             originalStoreID,
						BatchID:             batch.ID,
						SupplyCode:          supply.Code,
						SurgeryID:           batch.SurgeryID,
						OriginalAmount:      int(realCount) + 1, // +1 porque vamos a transferir uno
						CurrentInStore:      int(realCount),     // Stock actual sin contar el que se transfiere
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
				Status:          models.TransferStatusPending, // Pendiente hasta que se retire físicamente
				TransferReason:  fmt.Sprintf("Transferencia desde carrito %s", cart.CartNumber),
				SendDate:        now,
				Notes:           fmt.Sprintf("Transferencia automática desde carrito de solicitud %s. Pendiente de retiro físico.", cart.SupplyRequest.RequestNumber),
			}

			if err := tx.Create(&transfer).Error; err != nil {
				return fmt.Errorf("error al crear transferencia para insumo %s: %w", supply.QRCode, err)
			}

			// Actualizar ubicación del insumo - queda pendiente de retiro físico
			// El estado será "pendiente_retiro" hasta que alguien lo escanee para retirarlo
			// IMPORTANTE: NO cambiar location_type a pavilion hasta que se retire físicamente
			// El insumo sigue físicamente en bodega, solo está preparado para retiro
			// supply.LocationType = models.SupplyLocationStore // Mantener en bodega
			// supply.LocationID = originalStoreID // Mantener ID de bodega original
			supply.InTransit = false // No está en tránsito aún, está pendiente de retiro
			supply.TransferDate = &now
			supply.TransferredBy = &userRUT
			supply.Status = models.StatusPendingPickup // Pendiente de retiro físico

			if err := tx.Save(&supply).Error; err != nil {
				return fmt.Errorf("error al actualizar insumo %s: %w", supply.QRCode, err)
			}

			// Registrar en historial
			// IMPORTANTE: El destino es pavilion pero la ubicación actual sigue siendo store
			// porque físicamente el insumo aún está en bodega
			// Usar originalStoreID para registrar correctamente la bodega de origen (puede ser bodega secundaria)
			originType := models.DestinationTypeStore
			history := models.SupplyHistory{
				MedicalSupplyID: supply.ID,
				DateTime:        now,
				Status:          models.StatusPendingPickup,
				DestinationType: models.DestinationTypePavilion, // Destino final
				DestinationID:   pavilionID,
				UserRUT:         userRUT,
				Notes:           fmt.Sprintf("Preparado para retiro desde carrito %s. Debe ser escaneado al retirar de bodega. El insumo físicamente sigue en bodega.", cart.CartNumber),
				OriginType:      &originType,
				OriginID:        &originalStoreID,
			}
			if err := tx.Create(&history).Error; err != nil {
				return fmt.Errorf("error registrando historial: %w", err)
			}
		}

		return nil
	})
}
