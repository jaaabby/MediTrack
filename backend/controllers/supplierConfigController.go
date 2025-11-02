package controllers

import (
	"meditrack/models"
	"meditrack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SupplierConfigController struct {
	supplierConfigService services.SupplierConfigService
}

func NewSupplierConfigController(supplierConfigService services.SupplierConfigService) *SupplierConfigController {
	return &SupplierConfigController{
		supplierConfigService: supplierConfigService,
	}
}

// CreateSupplierConfig crea una nueva configuración de proveedor
func (c *SupplierConfigController) CreateSupplierConfig(ctx *gin.Context) {
	var configRequest struct {
		SupplierName       string `json:"supplier_name" binding:"required"`
		ExpirationAlertDays int    `json:"expiration_alert_days" binding:"required,min=1"`
		Notes              string `json:"notes,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&configRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	config := models.SupplierConfig{
		SupplierName:       configRequest.SupplierName,
		ExpirationAlertDays: configRequest.ExpirationAlertDays,
		Notes:              configRequest.Notes,
	}

	if err := c.supplierConfigService.CreateSupplierConfig(&config); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al crear configuración de proveedor: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Configuración de proveedor creada exitosamente",
		Data:    config,
	})
}

// GetSupplierConfig obtiene la configuración de un proveedor
func (c *SupplierConfigController) GetSupplierConfig(ctx *gin.Context) {
	supplierName := ctx.Param("name")
	if supplierName == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Nombre de proveedor requerido",
		})
		return
	}

	config, err := c.supplierConfigService.GetSupplierConfig(supplierName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Configuración de proveedor no encontrada: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    config,
	})
}

// GetAllSupplierConfigs obtiene todas las configuraciones de proveedores
func (c *SupplierConfigController) GetAllSupplierConfigs(ctx *gin.Context) {
	configs, err := c.supplierConfigService.GetAllSupplierConfigs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al obtener configuraciones de proveedores: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data:    configs,
	})
}

// UpdateSupplierConfig actualiza la configuración de un proveedor
func (c *SupplierConfigController) UpdateSupplierConfig(ctx *gin.Context) {
	supplierName := ctx.Param("name")
	if supplierName == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Nombre de proveedor requerido",
		})
		return
	}

	var configRequest struct {
		ExpirationAlertDays int    `json:"expiration_alert_days" binding:"required,min=1"`
		Notes              string `json:"notes,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&configRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	config := models.SupplierConfig{
		ExpirationAlertDays: configRequest.ExpirationAlertDays,
		Notes:              configRequest.Notes,
	}

	updatedConfig, err := c.supplierConfigService.UpdateSupplierConfig(supplierName, &config)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al actualizar configuración de proveedor: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Configuración de proveedor actualizada exitosamente",
		Data:    updatedConfig,
	})
}

// DeleteSupplierConfig elimina la configuración de un proveedor
func (c *SupplierConfigController) DeleteSupplierConfig(ctx *gin.Context) {
	supplierName := ctx.Param("name")
	if supplierName == "" {
		ctx.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Nombre de proveedor requerido",
		})
		return
	}

	if err := c.supplierConfigService.DeleteSupplierConfig(supplierName); err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Error al eliminar configuración de proveedor: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Configuración de proveedor eliminada exitosamente",
	})
}

