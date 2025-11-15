package controllers

import (
	"fmt"
	"meditrack/pkg/response"
	"meditrack/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	cartService *services.CartService
}

func NewCartController(cartService *services.CartService) *CartController {
	return &CartController{cartService: cartService}
}

// parseIDFromParam parsea el ID del parámetro de la URL
func (c *CartController) parseIDFromParam(ctx *gin.Context, paramName string) (int, error) {
	id := ctx.Param(paramName)
	intID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("ID inválido: debe ser un número entero")
	}
	return intID, nil
}

// getUserInfo extrae userRUT y userName del contexto
func getUserInfo(ctx *gin.Context) (string, string, error) {
	userRUT, exists := ctx.Get("userRUT")
	if !exists {
		userID, idExists := ctx.Get("user_id")
		if !idExists {
			return "", "", fmt.Errorf("usuario no autenticado")
		}
		userRUT = userID
	}

	userName, nameExists := ctx.Get("userName")
	if !nameExists {
		email, emailExists := ctx.Get("user_email")
		if emailExists {
			userName = email
		} else {
			userName = "Usuario Sistema"
		}
	}

	return userRUT.(string), userName.(string), nil
}

// ========================
// OBTENER CARRITOS
// ========================

// GetCartByRequestID obtiene el carrito asociado a una solicitud
// GET /api/carts/request/:requestId
func (c *CartController) GetCartByRequestID(ctx *gin.Context) {
	requestID, err := c.parseIDFromParam(ctx, "requestId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	cart, err := c.cartService.GetCartByRequestID(requestID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   "Carrito no encontrado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Carrito obtenido exitosamente",
		Data:    cart,
	})
}

// GetCartByID obtiene un carrito por su ID
// GET /api/carts/:id
func (c *CartController) GetCartByID(ctx *gin.Context) {
	cartID, err := c.parseIDFromParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	cart, err := c.cartService.GetCartByID(cartID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   "Carrito no encontrado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Carrito obtenido exitosamente",
		Data:    cart,
	})
}

// GetCartByQRCode obtiene el carrito asociado a un código QR
// GET /api/carts/qr/:qrCode
func (c *CartController) GetCartByQRCode(ctx *gin.Context) {
	qrCode := ctx.Param("qrCode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	cart, err := c.cartService.GetCartByQRCode(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   "Carrito no encontrado para este QR: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Carrito obtenido exitosamente",
		Data:    cart,
	})
}

// GetCartDetails obtiene los detalles completos de un carrito
// GET /api/carts/:id/details
func (c *CartController) GetCartDetails(ctx *gin.Context) {
	cartID, err := c.parseIDFromParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	details, err := c.cartService.GetCartDetails(cartID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   "Detalles del carrito no encontrados: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Detalles obtenidos exitosamente",
		Data:    details,
	})
}

// GetAllCarts obtiene todos los carritos con paginación
// GET /api/carts
func (c *CartController) GetAllCarts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	status := ctx.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	carts, total, err := c.cartService.GetAllCarts(page, pageSize, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al obtener carritos: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Carritos obtenidos exitosamente",
		"data":       carts,
		"page":       page,
		"pageSize":   pageSize,
		"total":      total,
		"totalPages": (int(total) + pageSize - 1) / pageSize,
	})
}

// ========================
// CREAR Y GESTIONAR CARRITOS
// ========================

// CreateCartForRequest crea un carrito para una solicitud (usado internamente o manualmente)
// POST /api/carts/request/:requestId
func (c *CartController) CreateCartForRequest(ctx *gin.Context) {
	requestID, err := c.parseIDFromParam(ctx, "requestId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userRUT, userName, err := getUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	cart, err := c.cartService.CreateCartForRequest(requestID, userRUT, userName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Error al crear carrito: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response.Response{
		Success: true,
		Message: "Carrito creado exitosamente",
		Data:    cart,
	})
}

// AddItemToCart agrega un item al carrito
// POST /api/carts/:id/items
func (c *CartController) AddItemToCart(ctx *gin.Context) {
	cartID, err := c.parseIDFromParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var req struct {
		AssignmentID int `json:"assignment_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	userRUT, userName, err := getUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	cartItem, err := c.cartService.AddItemToCart(cartID, req.AssignmentID, userRUT, userName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Error al agregar item al carrito: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response.Response{
		Success: true,
		Message: "Item agregado al carrito exitosamente",
		Data:    cartItem,
	})
}

// RemoveItemFromCart marca un item como inactivo en el carrito
// DELETE /api/carts/:id/items/:itemId
func (c *CartController) RemoveItemFromCart(ctx *gin.Context) {
	cartID, err := c.parseIDFromParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	itemID, err := c.parseIDFromParam(ctx, "itemId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userRUT, userName, err := getUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := c.cartService.RemoveItemFromCart(cartID, itemID, userRUT, userName); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Error al remover item del carrito: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Item removido del carrito exitosamente",
	})
}

// CloseCart cierra un carrito
// POST /api/carts/:id/close
func (c *CartController) CloseCart(ctx *gin.Context) {
	cartID, err := c.parseIDFromParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userRUT, userName, err := getUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := c.cartService.CloseCart(cartID, userRUT, userName); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al cerrar carrito: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Carrito cerrado exitosamente",
	})
}

// ========================
// OPERACIONES DE ITEMS
// ========================

// MarkItemAsUsed marca un item del carrito como utilizado (consumido)
// POST /api/carts/:id/items/:itemId/use
func (c *CartController) MarkItemAsUsed(ctx *gin.Context) {
	cartID, err := c.parseIDFromParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	itemID, err := c.parseIDFromParam(ctx, "itemId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userRUT, userName, err := getUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := c.cartService.MarkItemAsUsed(cartID, itemID, userRUT, userName); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Error al marcar item como utilizado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Item marcado como utilizado exitosamente",
	})
}

// MarkItemForReturn marca un item del carrito para devolución
// POST /api/carts/:id/items/:itemId/return
func (c *CartController) MarkItemForReturn(ctx *gin.Context) {
	cartID, err := c.parseIDFromParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	itemID, err := c.parseIDFromParam(ctx, "itemId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userRUT, userName, err := getUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	if err := c.cartService.MarkItemForReturn(cartID, itemID, userRUT, userName, req.Reason); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Error al marcar item para devolución: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Item marcado para devolución exitosamente",
	})
}

// BatchOperationItems procesa múltiples items del carrito en una sola operación
// POST /api/carts/:id/items/batch-operation
func (c *CartController) BatchOperationItems(ctx *gin.Context) {
	cartID, err := c.parseIDFromParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var req struct {
		Items []services.BatchOperationItem `json:"items" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	if len(req.Items) == 0 {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Debe proporcionar al menos un item para procesar",
		})
		return
	}

	userRUT, userName, err := getUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	result, err := c.cartService.BatchOperationItems(cartID, req.Items, userRUT, userName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Error al procesar items: " + err.Error(),
		})
		return
	}

	// Determinar el mensaje según los resultados
	var message string
	if result.ErrorCount == 0 {
		message = fmt.Sprintf("Todos los items procesados exitosamente (%d items)", result.SuccessCount)
	} else if result.SuccessCount > 0 {
		message = fmt.Sprintf("Procesados %d items exitosamente, %d fallaron", result.SuccessCount, result.ErrorCount)
	} else {
		message = fmt.Sprintf("Error al procesar items: todos fallaron (%d items)", result.ErrorCount)
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: result.ErrorCount == 0,
		Message: message,
		Data:    result,
	})
}

// ========================
// TRANSFERENCIAS
// ========================

// TransferCartToPavilion transfiere todos los items del carrito al pabellón
// POST /api/carts/:id/transfer-to-pavilion
func (c *CartController) TransferCartToPavilion(ctx *gin.Context) {
	cartID, err := c.parseIDFromParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	userRUT, userName, err := getUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := c.cartService.TransferCartToPavilion(cartID, userRUT, userName); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al transferir carrito al pabellón: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Carrito transferido al pabellón exitosamente. El pabellón debe confirmar recepción.",
	})
}
