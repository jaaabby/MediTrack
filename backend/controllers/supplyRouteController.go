package controllers

import (
	"net/http"

	"meditrack/models"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SupplyRouteController maneja las peticiones HTTP relacionadas con rutas de suministro
type SupplyRouteController struct {
	supplyRouteService services.SupplyRouteService
}

// NewSupplyRouteController crea una nueva instancia de SupplyRouteController
func NewSupplyRouteController(supplyRouteService services.SupplyRouteService) *SupplyRouteController {
	return &SupplyRouteController{
		supplyRouteService: supplyRouteService,
	}
}

// CreateSupplyRoute crea una nueva ruta de suministro
func (c *SupplyRouteController) CreateSupplyRoute(ctx *gin.Context) {
	var route models.SupplyRoute
	if err := ctx.ShouldBindJSON(&route); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de ruta de suministro inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyRouteService.CreateSupplyRoute(ctx, &route); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear ruta de suministro: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Ruta de suministro creada exitosamente",
		Data:    route,
	})
}

// GetSupplyRouteByID obtiene una ruta de suministro por ID
func (c *SupplyRouteController) GetSupplyRouteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de ruta de suministro requerido",
		})
		return
	}

	route, err := c.supplyRouteService.GetSupplyRouteByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Ruta de suministro no encontrada: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    route,
	})
}

// GetSupplyRoutesByPatientID obtiene rutas de suministro por ID de paciente
func (c *SupplyRouteController) GetSupplyRoutesByPatientID(ctx *gin.Context) {
	patientID := ctx.Query("patient_id")
	if patientID == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de paciente requerido",
		})
		return
	}

	routes, err := c.supplyRouteService.GetSupplyRoutesByPatientID(ctx, patientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener rutas de suministro: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    routes,
	})
}

// GetSupplyRoutesByOperatingRoomID obtiene rutas de suministro por ID de sala de operación
func (c *SupplyRouteController) GetSupplyRoutesByOperatingRoomID(ctx *gin.Context) {
	roomID := ctx.Query("operating_room_id")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de sala de operación requerido",
		})
		return
	}

	routes, err := c.supplyRouteService.GetSupplyRoutesByOperatingRoomID(ctx, roomID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener rutas de suministro: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    routes,
	})
}

// GetAllSupplyRoutes obtiene todas las rutas de suministro
func (c *SupplyRouteController) GetAllSupplyRoutes(ctx *gin.Context) {
	routes, err := c.supplyRouteService.GetAllSupplyRoutes(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener rutas de suministro: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    routes,
	})
}

// UpdateSupplyRoute actualiza una ruta de suministro existente
func (c *SupplyRouteController) UpdateSupplyRoute(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de ruta de suministro requerido",
		})
		return
	}

	var route models.SupplyRoute
	if err := ctx.ShouldBindJSON(&route); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de ruta de suministro inválidos: " + err.Error(),
		})
		return
	}

	if err := c.supplyRouteService.UpdateSupplyRoute(ctx, &route); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar ruta de suministro: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Ruta de suministro actualizada exitosamente",
		Data:    route,
	})
}

// DeleteSupplyRoute elimina una ruta de suministro
func (c *SupplyRouteController) DeleteSupplyRoute(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de ruta de suministro requerido",
		})
		return
	}

	if err := c.supplyRouteService.DeleteSupplyRoute(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al eliminar ruta de suministro: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Ruta de suministro eliminada exitosamente",
	})
}
