package controllers

import (
	"fmt"
	"meditrack/models"
	"meditrack/pkg/response"
	"meditrack/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// BatchController maneja las peticiones HTTP relacionadas con lotes
type BatchController struct {
	batchService services.BatchService
}

// NewBatchController crea una nueva instancia de BatchController
func NewBatchController(batchService services.BatchService) *BatchController {
	return &BatchController{
		batchService: batchService,
	}
}

// CreateBatch crea un nuevo lote
func (c *BatchController) CreateBatch(ctx *gin.Context) {
	// Crear struct temporal sin ID para evitar conflictos
	var batchRequest struct {
		ExpirationDate time.Time `json:"expiration_date"`
		Amount         int       `json:"amount"`
		Supplier       string    `json:"supplier"`
		StoreID        int       `json:"store_id"`
	}

	if err := ctx.ShouldBindJSON(&batchRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos de lote inválidos: " + err.Error(),
		})
		return
	}

	// Crear el modelo Batch sin ID (se auto-generará)
	batch := models.Batch{
		ExpirationDate: batchRequest.ExpirationDate,
		Amount:         batchRequest.Amount,
		Supplier:       batchRequest.Supplier,
		StoreID:        batchRequest.StoreID,
		// QRCode se generará en el servicio
		// ID se auto-generará
	}

	if err := c.batchService.CreateBatch(&batch); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al crear lote: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response.Response{
		Success: true,
		Message: "Lote creado exitosamente",
		Data:    batch,
	})
}

// CreateBatchWithIndividualSupplies crea un lote completo con sus insumos individuales
func (c *BatchController) CreateBatchWithIndividualSupplies(ctx *gin.Context) {
	var request struct {
		Batch struct {
			ExpirationDate     string `json:"expiration_date" binding:"required"`
			Amount             int    `json:"amount" binding:"required,min=1"`
			Supplier           string `json:"supplier" binding:"required"`
			StoreID            int    `json:"store_id" binding:"required"`
			ExpirationAlertDays *int   `json:"expiration_alert_days,omitempty"` // Opcional: días de alerta para el proveedor
		} `json:"batch" binding:"required"`
		SupplyCode struct {
			Code         int    `json:"code" binding:"required"`
			Name         string `json:"name" binding:"required"`
			CodeSupplier int    `json:"code_supplier" binding:"required"`
			CriticalStock int    `json:"critical_stock" binding:"required,min=1"`
		} `json:"supply_code" binding:"required"`
		IndividualCount int `json:"individual_count" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	// Convertir fecha
	expirationDate, err := parseDate(request.Batch.ExpirationDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Fecha de expiración inválida: " + err.Error(),
		})
		return
	}

	// Crear modelos SIN ID - se auto-generarán
	batch := &models.Batch{
		ExpirationDate: expirationDate,
		Amount:         request.Batch.Amount,
		Supplier:       request.Batch.Supplier,
		StoreID:        request.Batch.StoreID,
		// ID y QRCode se generarán automáticamente
	}

	supplyCode := &models.SupplyCode{
		Code:         request.SupplyCode.Code,
		Name:         request.SupplyCode.Name,
		CodeSupplier: request.SupplyCode.CodeSupplier,
		CriticalStock: request.SupplyCode.CriticalStock,
	}

	// Obtener días de alerta del request (si se proporciona, usar ese valor, sino usar 90 por defecto)
	expirationAlertDays := 90 // Valor por defecto
	if request.Batch.ExpirationAlertDays != nil && *request.Batch.ExpirationAlertDays > 0 {
		expirationAlertDays = *request.Batch.ExpirationAlertDays
	}

	// Crear lote con insumos individuales
	createdBatch, individualSupplies, err := c.batchService.CreateBatchWithIndividualSupplies(
		batch,
		supplyCode,
		request.IndividualCount,
		expirationAlertDays, // Pasar días de alerta para configuración del proveedor
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al crear lote con insumos individuales: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response.Response{
		Success: true,
		Message: "Lote con insumos individuales creado exitosamente",
		Data: map[string]interface{}{
			"batch":               createdBatch,
			"supply_code":         supplyCode,
			"individual_supplies": individualSupplies,
			"total_created":       len(individualSupplies),
		},
	})
}

// parseIDFromParam parsea el ID del parámetro de la URL
func (c *BatchController) parseIDFromParam(ctx *gin.Context) (int, error) {
	id := ctx.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("ID inválido: debe ser un número entero")
	}
	return intID, nil
}

// GetBatchByID obtiene un lote por ID
func (c *BatchController) GetBatchByID(ctx *gin.Context) {
	intID, err := c.parseIDFromParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	batch, err := c.batchService.GetBatchByID(intID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   "Lote no encontrado: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Data:    batch,
	})
}

// GetBatchWithSupplyInfo obtiene un lote con información completa de sus insumos
func (c *BatchController) GetBatchWithSupplyInfo(ctx *gin.Context) {
	intID, err := c.parseIDFromParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	batchInfo, err := c.batchService.GetBatchWithSupplyInfo(intID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Información del lote obtenida exitosamente",
		Data:    batchInfo,
	})
}

// GetBatchByQR obtiene información completa de un lote por su código QR
func (c *BatchController) GetBatchByQR(ctx *gin.Context) {
	qrCode := ctx.Param("qrcode")
	if qrCode == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Código QR requerido",
		})
		return
	}

	batchInfo, err := c.batchService.GetBatchQRInfo(qrCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Información del lote por QR obtenida exitosamente",
		Data:    batchInfo,
	})
}

// GetAllBatches obtiene todos los lotes con filtros opcionales
func (c *BatchController) GetAllBatches(ctx *gin.Context) {
	// Obtener parámetros de query opcionales
	var surgeryID *int
	if surgeryIDStr := ctx.Query("surgery_id"); surgeryIDStr != "" {
		if id, err := strconv.Atoi(surgeryIDStr); err == nil {
			surgeryID = &id
		}
	}

	var storeID *int
	if storeIDStr := ctx.Query("store_id"); storeIDStr != "" {
		if id, err := strconv.Atoi(storeIDStr); err == nil {
			storeID = &id
		}
	}

	supplier := ctx.Query("supplier")

	// Si no hay filtros, obtener todos
	if surgeryID == nil && storeID == nil && supplier == "" {
		batches, err := c.batchService.GetAllBatches()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.Response{
				Success: false,
				Error:   "Error al obtener lotes: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, response.Response{
			Success: true,
			Data:    batches,
		})
		return
	}

	// Obtener lotes con filtros
	batches, err := c.batchService.GetBatchesWithFilters(surgeryID, storeID, supplier)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al obtener lotes: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Data:    batches,
	})
}

// GetBatchesNeedingSync obtiene lotes que necesitan sincronización
func (c *BatchController) GetBatchesNeedingSync(ctx *gin.Context) {
	batches, err := c.batchService.GetBatchesNeedingSync()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al obtener lotes que necesitan sincronización: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Lotes que necesitan sincronización obtenidos exitosamente",
		Data: map[string]interface{}{
			"batches_needing_sync": batches,
			"total_count":          len(batches),
		},
	})
}

// UpdateBatch actualiza un lote existente
func (c *BatchController) UpdateBatch(ctx *gin.Context) {
	intID, err := c.parseIDFromParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var batchRequest struct {
		ExpirationDate time.Time `json:"expiration_date"`
		Amount         int       `json:"amount"`
		Supplier       string    `json:"supplier"`
		StoreID        int       `json:"store_id"`
		// No incluimos QRCode - debe mantenerse el original
	}

	if err := ctx.ShouldBindJSON(&batchRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos de lote inválidos: " + err.Error(),
		})
		return
	}

	// Crear modelo con el ID correcto pero sin QRCode
	batch := models.Batch{
		ID:             intID,
		ExpirationDate: batchRequest.ExpirationDate,
		Amount:         batchRequest.Amount,
		Supplier:       batchRequest.Supplier,
		StoreID:        batchRequest.StoreID,
		// QRCode no se actualiza para mantener trazabilidad
	}

	updatedBatch, err := c.batchService.UpdateBatch(intID, &batch)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al actualizar lote: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Lote actualizado exitosamente",
		Data:    updatedBatch,
	})
}

// UpdateBatchAmount actualiza solo la cantidad de un lote
func (c *BatchController) UpdateBatchAmount(ctx *gin.Context) {
	intID, err := c.parseIDFromParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var request struct {
		Amount int `json:"amount" binding:"required,min=0"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Datos inválidos: " + err.Error(),
		})
		return
	}

	err = c.batchService.UpdateBatchAmount(intID, request.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al actualizar cantidad del lote: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Cantidad del lote actualizada exitosamente",
		Data: map[string]interface{}{
			"batch_id":   intID,
			"new_amount": request.Amount,
		},
	})
}

// DeleteBatch elimina un lote
func (c *BatchController) DeleteBatch(ctx *gin.Context) {
	intID, err := c.parseIDFromParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := c.batchService.DeleteBatch(intID); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al eliminar lote: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Lote eliminado exitosamente",
	})
}

// SyncAllBatchAmounts sincroniza todas las cantidades de lotes
func (c *BatchController) SyncAllBatchAmounts(ctx *gin.Context) {
	err := c.batchService.SyncAllBatchAmounts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error sincronizando cantidades de lotes: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Todas las cantidades de lotes sincronizadas exitosamente",
	})
}

// CheckLowStockAlert verifica y envía alertas de stock bajo para un lote específico
func (c *BatchController) CheckLowStockAlert(ctx *gin.Context) {
	intID, err := c.parseIDFromParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Obtener threshold del query parameter, por defecto 5
	thresholdStr := ctx.DefaultQuery("threshold", "5")
	threshold, err := strconv.Atoi(thresholdStr)
	if err != nil {
		threshold = 5
	}

	err = c.batchService.CheckLowStockAlert(intID, threshold)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al verificar alerta de stock bajo: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Verificación de alerta de stock bajo completada",
		Data: map[string]interface{}{
			"batch_id":  intID,
			"threshold": threshold,
		},
	})
}

// CheckExpirationAlert verifica y envía alertas de vencimiento próximo para un lote específico
func (c *BatchController) CheckExpirationAlert(ctx *gin.Context) {
	intID, err := c.parseIDFromParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Obtener días de threshold del query parameter, por defecto 30
	daysStr := ctx.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	err = c.batchService.CheckExpirationAlert(intID, days)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al verificar alerta de vencimiento: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Verificación de alerta de vencimiento completada",
		Data: map[string]interface{}{
			"batch_id":       intID,
			"days_threshold": days,
		},
	})
}

// parseDate convierte una fecha en string a time.Time
func parseDate(dateStr string) (time.Time, error) {
	// Intentar diferentes formatos de fecha
	formats := []string{
		"2006-01-02",
		"02/01/2006",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("formato de fecha no reconocido: %s", dateStr)
}
