package controllers

import (
	"net/http"

	"meditrack/models"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// DoctorSupplyRouteController maneja las peticiones HTTP relacionadas con la relación doctor-ruta de suministro
type DoctorSupplyRouteController struct {
	doctorSupplyRouteService services.DoctorSupplyRouteService
}

// NewDoctorSupplyRouteController crea una nueva instancia de DoctorSupplyRouteController
func NewDoctorSupplyRouteController(doctorSupplyRouteService services.DoctorSupplyRouteService) *DoctorSupplyRouteController {
	return &DoctorSupplyRouteController{
		doctorSupplyRouteService: doctorSupplyRouteService,
	}
}

// CreateDoctorSupplyRoute crea una nueva relación doctor-ruta de suministro
func (c *DoctorSupplyRouteController) CreateDoctorSupplyRoute(ctx *gin.Context) {
	var relation models.DoctorSupplyRoute
	if err := ctx.ShouldBindJSON(&relation); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de relación doctor-ruta inválidos: " + err.Error(),
		})
		return
	}

	if err := c.doctorSupplyRouteService.CreateDoctorSupplyRoute(ctx, &relation); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear relación doctor-ruta: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Relación doctor-ruta creada exitosamente",
		Data:    relation,
	})
}

// GetDoctorSupplyRouteByID obtiene una relación doctor-ruta por ID
func (c *DoctorSupplyRouteController) GetDoctorSupplyRouteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de relación doctor-ruta requerido",
		})
		return
	}

	relation, err := c.doctorSupplyRouteService.GetDoctorSupplyRouteByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Relación doctor-ruta no encontrada: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    relation,
	})
}

// GetDoctorSupplyRoutesByDoctorID obtiene relaciones por ID de doctor
func (c *DoctorSupplyRouteController) GetDoctorSupplyRoutesByDoctorID(ctx *gin.Context) {
	doctorID := ctx.Query("doctor_id")
	if doctorID == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de doctor requerido",
		})
		return
	}

	relations, err := c.doctorSupplyRouteService.GetDoctorSupplyRoutesByDoctorID(ctx, doctorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener relaciones doctor-ruta: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    relations,
	})
}

// GetDoctorSupplyRoutesBySupplyRouteID obtiene relaciones por ID de ruta de suministro
func (c *DoctorSupplyRouteController) GetDoctorSupplyRoutesBySupplyRouteID(ctx *gin.Context) {
	routeID := ctx.Query("supply_route_id")
	if routeID == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de ruta de suministro requerido",
		})
		return
	}

	relations, err := c.doctorSupplyRouteService.GetDoctorSupplyRoutesBySupplyRouteID(ctx, routeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener relaciones doctor-ruta: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    relations,
	})
}

// GetAllDoctorSupplyRoutes obtiene todas las relaciones doctor-ruta
func (c *DoctorSupplyRouteController) GetAllDoctorSupplyRoutes(ctx *gin.Context) {
	relations, err := c.doctorSupplyRouteService.GetAllDoctorSupplyRoutes(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener relaciones doctor-ruta: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    relations,
	})
}

// UpdateDoctorSupplyRoute actualiza una relación doctor-ruta existente
func (c *DoctorSupplyRouteController) UpdateDoctorSupplyRoute(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de relación doctor-ruta requerido",
		})
		return
	}

	var relation models.DoctorSupplyRoute
	if err := ctx.ShouldBindJSON(&relation); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos de relación doctor-ruta inválidos: " + err.Error(),
		})
		return
	}

	if err := c.doctorSupplyRouteService.UpdateDoctorSupplyRoute(ctx, &relation); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar relación doctor-ruta: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Relación doctor-ruta actualizada exitosamente",
		Data:    relation,
	})
}

// DeleteDoctorSupplyRoute elimina una relación doctor-ruta
func (c *DoctorSupplyRouteController) DeleteDoctorSupplyRoute(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "ID de relación doctor-ruta requerido",
		})
		return
	}

	if err := c.doctorSupplyRouteService.DeleteDoctorSupplyRoute(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al eliminar relación doctor-ruta: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Relación doctor-ruta eliminada exitosamente",
	})
}
