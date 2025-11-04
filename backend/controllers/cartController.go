package controllers

import (
	"meditrack/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CartController maneja las peticiones HTTP relacionadas con carritos
type CartController struct {
	cartService *services.CartService
}

// NewCartController crea una nueva instancia del controlador de carritos
func NewCartController(cartService *services.CartService) *CartController {
	return &CartController{
		cartService: cartService,
	}
}

// GetCartByRequestID obtiene el carrito asociado a una solicitud
// GET /api/carts/request/:requestId
func (c *CartController) GetCartByRequestID(ctx *gin.Context) {
	requestID, err := strconv.Atoi(ctx.Param("requestId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de solicitud inválido: " + err.Error(),
		})
		return
	}

	cart, err := c.cartService.GetCartByRequestID(requestID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Carrito no encontrado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Carrito obtenido exitosamente",
		Data:    cart,
	})
}

// GetCartByID obtiene un carrito por su ID
// GET /api/carts/:id
func (c *CartController) GetCartByID(ctx *gin.Context) {
	cartID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de carrito inválido: " + err.Error(),
		})
		return
	}

	cart, err := c.cartService.GetCartByID(cartID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Carrito no encontrado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
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
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	cart, err := c.cartService.GetCartByQRCode(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Carrito no encontrado para este QR: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Carrito obtenido exitosamente",
		Data:    cart,
	})
}

// AddItemToCartRequest estructura para agregar items al carrito
type AddItemToCartRequest struct {
	AssignmentID int `json:"assignment_id" binding:"required"`
}

// AddItemToCart agrega un item al carrito
// POST /api/carts/:id/items
func (c *CartController) AddItemToCart(ctx *gin.Context) {
	cartID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de carrito inválido: " + err.Error(),
		})
		return
	}

	var req AddItemToCartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	// Obtener información del usuario desde el contexto
	userRUT, _ := ctx.Get("userRUT")
	userName, _ := ctx.Get("userName")

	cartItem, err := c.cartService.AddItemToCart(cartID, req.AssignmentID, userRUT.(string), userName.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Error al agregar item al carrito: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Item agregado al carrito exitosamente",
		Data:    cartItem,
	})
}

// RemoveItemFromCart marca un item como inactivo en el carrito
// DELETE /api/carts/:id/items/:itemId
func (c *CartController) RemoveItemFromCart(ctx *gin.Context) {
	cartID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de carrito inválido: " + err.Error(),
		})
		return
	}

	itemID, err := strconv.Atoi(ctx.Param("itemId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de item inválido: " + err.Error(),
		})
		return
	}

	// Obtener información del usuario desde el contexto
	userRUT, _ := ctx.Get("userRUT")
	userName, _ := ctx.Get("userName")

	if err := c.cartService.RemoveItemFromCart(cartID, itemID, userRUT.(string), userName.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Error al remover item del carrito: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Item removido del carrito exitosamente",
	})
}

// CloseCart cierra un carrito
// POST /api/carts/:id/close
func (c *CartController) CloseCart(ctx *gin.Context) {
	cartID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de carrito inválido: " + err.Error(),
		})
		return
	}

	// Obtener información del usuario desde el contexto
	userRUT, _ := ctx.Get("userRUT")
	userName, _ := ctx.Get("userName")

	if err := c.cartService.CloseCart(cartID, userRUT.(string), userName.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Error al cerrar carrito: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Carrito cerrado exitosamente",
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
		ctx.JSON(http.StatusInternalServerError, Response{
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

// GetCartDetails obtiene los detalles completos de un carrito
// GET /api/carts/:id/details
func (c *CartController) GetCartDetails(ctx *gin.Context) {
	cartID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de carrito inválido: " + err.Error(),
		})
		return
	}

	details, err := c.cartService.GetCartDetails(cartID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Detalles del carrito no encontrados: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Detalles obtenidos exitosamente",
		Data:    details,
	})
}

// CreateCartForRequest crea un carrito para una solicitud (usado internamente o manualmente)
// POST /api/carts/request/:requestId
func (c *CartController) CreateCartForRequest(ctx *gin.Context) {
	requestID, err := strconv.Atoi(ctx.Param("requestId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de solicitud inválido: " + err.Error(),
		})
		return
	}

	// Obtener información del usuario desde el contexto
	userRUT, _ := ctx.Get("userRUT")
	userName, _ := ctx.Get("userName")

	cart, err := c.cartService.CreateCartForRequest(requestID, userRUT.(string), userName.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Error al crear carrito: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Carrito creado exitosamente",
		Data:    cart,
	})
}

// MarkItemAsUsed marca un item del carrito como utilizado (consumido)
// POST /api/carts/:id/items/:itemId/use
func (c *CartController) MarkItemAsUsed(ctx *gin.Context) {
	cartID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de carrito inválido: " + err.Error(),
		})
		return
	}

	itemID, err := strconv.Atoi(ctx.Param("itemId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de item inválido: " + err.Error(),
		})
		return
	}

	// Obtener información del usuario desde el contexto
	userRUT, _ := ctx.Get("userRUT")
	userName, _ := ctx.Get("userName")

	if err := c.cartService.MarkItemAsUsed(cartID, itemID, userRUT.(string), userName.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Error al marcar item como utilizado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Item marcado como utilizado exitosamente",
	})
}

// MarkItemForReturn marca un item del carrito para devolución
// POST /api/carts/:id/items/:itemId/return
func (c *CartController) MarkItemForReturn(ctx *gin.Context) {
	cartID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de carrito inválido: " + err.Error(),
		})
		return
	}

	itemID, err := strconv.Atoi(ctx.Param("itemId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de item inválido: " + err.Error(),
		})
		return
	}

	// Obtener información del usuario desde el contexto
	userRUT, _ := ctx.Get("userRUT")
	userName, _ := ctx.Get("userName")

	// Obtener el motivo de devolución del body
	var req struct {
		Reason string `json:"reason"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	if req.Reason == "" {
		req.Reason = "Sin especificar"
	}

	if err := c.cartService.MarkItemForReturn(cartID, itemID, userRUT.(string), userName.(string), req.Reason); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Error al marcar item para devolución: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Item marcado para devolución exitosamente",
	})
}
