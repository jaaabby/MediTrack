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

		// Actualizar el insumo médico a estado "consumido"
		if err := tx.Model(&models.MedicalSupply{}).
			Where("id = ?", cartItem.SupplyRequestQRAssignment.MedicalSupplyID).
			Update("status", models.StatusConsumed).Error; err != nil {
			return fmt.Errorf("error actualizando estado del insumo: %w", err)
		}

		// Registrar en historial
		history := models.SupplyHistory{
			MedicalSupplyID: cartItem.SupplyRequestQRAssignment.MedicalSupplyID,
			DateTime:        now,
			Status:          "consumido",
			DestinationType: "pabellon",
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
}

// MarkItemForReturn marca un item del carrito para devolución
func (s *CartService) MarkItemForReturn(cartID, itemID int, userRUT, userName, reason string) error {
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

		// Devolver el insumo a bodega (disponible)
		if err := tx.Model(&models.MedicalSupply{}).
			Where("id = ?", supply.ID).
			Updates(map[string]interface{}{
				"status":        models.StatusAvailable,
				"location_type": models.SupplyLocationStore,
			}).Error; err != nil {
			return fmt.Errorf("error actualizando estado del insumo: %w", err)
		}

		// Registrar en historial
		history := models.SupplyHistory{
			MedicalSupplyID: supply.ID,
			DateTime:        now,
			Status:          "disponible",
			DestinationType: models.DestinationTypeStore,
			DestinationID:   supply.LocationID,
			UserRUT:         userRUT,
			Notes:           fmt.Sprintf("Devuelto desde carrito %s. Motivo: %s", cart.CartNumber, reason),
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
}
