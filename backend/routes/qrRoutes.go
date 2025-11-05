package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupQRRoutes configura las rutas de códigos QR con funcionalidades mejoradas y trazabilidad completa
func SetupQRRoutes(router *gin.RouterGroup, qrService services.QRService, medicalSupplyService services.MedicalSupplyService, cartService *services.CartService) {
	qrController := controllers.NewQRController(qrService)
	// Set cart service if available
	if cartService != nil {
		qrController.SetCartService(cartService)
	}
	qrController.SetMedicalSupplyService(medicalSupplyService)

	qr := router.Group("/qr")
	{
		// === RUTAS BÁSICAS DE QR CON TRAZABILIDAD AUTOMÁTICA ===

		// Escanear un código QR y obtener toda su información (REGISTRA AUTOMÁTICAMENTE EL ESCANEO)
		qr.GET("/scan/:qrcode", qrController.ScanQR)

		// Validar si un código QR es válido
		qr.GET("/validate/:qrcode", qrController.ValidateQR)

		// Obtener historial de un insumo por código QR (sin registrar escaneo)
		qr.GET("/history/:qrcode", qrController.GetSupplyHistory)

		// === NUEVAS RUTAS DE TRAZABILIDAD COMPLETA ===

		// Obtener trazabilidad completa incluyendo todos los escaneos
		qr.GET("/traceability/:qrcode", qrController.GetCompleteTraceability)

		// Obtener historial de escaneos específicamente
		qr.GET("/scan-history/:qrcode", qrController.GetScanHistory)

		// Obtener estadísticas de escaneo de un QR
		qr.GET("/scan-stats/:qrcode", qrController.GetScanStatistics)

		// Registrar manualmente un evento de escaneo (para integraciones externas)
		qr.POST("/register-scan", qrController.RegisterManualScanEvent)

		// === GENERACIÓN DE QR CODES ===

		// Generar códigos QR con imagen
		qr.POST("/generate/batch", qrController.GenerateBatchQR)
		qr.POST("/generate/supply", qrController.GenerateSupplyQR)

		// === IMÁGENES QR ===

		// Servir imágenes QR (para mostrar en la UI)
		qr.GET("/image/:qrcode", qrController.GetQRImage)

		// Descargar imágenes QR con diferentes resoluciones
		// Query param: resolution=normal|high (default: normal)
		qr.GET("/download/:qrcode", qrController.DownloadQRImage)

		// === FUNCIONALIDADES DE CONSUMO Y TRANSFERENCIA ===

		// Consumir un insumo por QR (actualiza automáticamente las cantidades del lote)
		qr.POST("/consume", qrController.ConsumeSupply)

		// Transferir un insumo individual por QR
		qr.POST("/transfer", qrController.TransferSupply)

		// Recepcionar un insumo que está en camino al pabellón
		qr.POST("/receive", qrController.ReceiveSupply)

		// === FUNCIONALIDADES DE RETORNO A BODEGA ===

		// Regresar un insumo a bodega manualmente
		qr.POST("/return-to-store", qrController.ReturnSupplyToStore)

		// Confirmar llegada de insumo a bodega
		qr.POST("/confirm-arrival-to-store", qrController.ConfirmArrivalToStore)

		// Obtener lista de insumos que deben regresar a bodega (8 horas laborales sin consumir)
		qr.GET("/supplies-for-return", qrController.GetSuppliesForReturn)

		// Ejecutar manualmente el proceso automático de retornos
		qr.POST("/process-automatic-returns", qrController.ProcessAutomaticReturns)

		// Consumir un insumo individual específico (REGISTRA AUTOMÁTICAMENTE EL ESCANEO DE CONSUMO)
		qr.POST("/consume/individual", qrController.ConsumeIndividualSupply)

		// Consumir múltiples insumos en lote
		qr.POST("/consume/bulk", qrController.BulkConsumeSupplies)

		// Verificar disponibilidad de un insumo para consumo
		qr.GET("/verify/:qrcode", qrController.VerifySupplyAvailability)

		// === INFORMACIÓN DETALLADA DE INSUMOS INDIVIDUALES ===

		// Obtener información detallada de un insumo con datos del lote
		qr.GET("/details/:qrcode", qrController.GetSupplyDetails)

		// === ADMINISTRACIÓN Y ESTADÍSTICAS ===

		// Sincronizar cantidades de lotes con productos individuales
		qr.POST("/sync/batch-amounts", qrController.SyncBatchAmounts)

		// Obtener estadísticas generales de uso de QR codes
		qr.GET("/stats", qrController.GetQRStats)

		// === RUTAS ADICIONALES PARA MEJOR ORGANIZACIÓN ===

		// Grupo de rutas específicas para insumos individuales
		individual := qr.Group("/individual")
		{
			// Escanear insumo individual con contexto específico
			individual.GET("/scan/:qrcode", qrController.ScanQR)

			// Obtener información completa de un insumo individual
			individual.GET("/details/:qrcode", qrController.GetSupplyDetails)

			// Consumir insumo individual
			individual.POST("/consume", qrController.ConsumeIndividualSupply)

			// Verificar disponibilidad para consumo
			individual.GET("/verify/:qrcode", qrController.VerifySupplyAvailability)

			// Obtener trazabilidad específica de insumo individual
			individual.GET("/traceability/:qrcode", qrController.GetCompleteTraceability)

			// Obtener información de insumo individual (ruta alternativa)
			individual.GET("/:qrcode", qrController.ScanQR)

			// Verificar disponibilidad de insumo individual (ruta alternativa)
			individual.GET("/:qrcode/availability", qrController.VerifySupplyAvailability)

			// Obtener historial de insumo individual (ruta alternativa)
			individual.GET("/:qrcode/history", qrController.GetSupplyHistory)
		}

		// Grupo de rutas específicas para lotes
		batch := qr.Group("/batch")
		{
			// Escanear lote (solo información, no para consumo directo)
			batch.GET("/scan/:qrcode", qrController.ScanQR)

			// Obtener información del lote
			batch.GET("/details/:qrcode", qrController.GetSupplyDetails)

			// Generar códigos QR para insumos individuales de un lote
			batch.POST("/generate-supplies", qrController.GenerateSuppliesFromBatch)

			// Sincronizar cantidades del lote
			batch.POST("/sync-amounts/:batch_id", qrController.SyncBatchAmounts)

			// Obtener insumos individuales de un lote
			batch.GET("/:batch_id/supplies", qrController.GetIndividualSuppliesByBatch)

			// Generar QR de lote
			batch.POST("/generate", qrController.GenerateBatchQR)
		}

		// === RUTAS PARA REPORTES Y ANALÍTICA ===
		analytics := qr.Group("/analytics")
		{
			// Obtener resumen de actividad de escaneos por período
			analytics.GET("/scan-summary", qrController.GetScanAnalytics)

			// Obtener top QRs más escaneados
			analytics.GET("/top-scanned", qrController.GetTopScannedQRs)

			// Obtener actividad por usuario
			analytics.GET("/user-activity", qrController.GetUserScanActivity)

			// Obtener actividad por pabellón
			analytics.GET("/pavilion-activity", qrController.GetPavilionScanActivity)

			// Obtener patrones de movimiento
			analytics.GET("/movement-patterns", qrController.GetMovementPatterns)
		}

		// === RUTAS PARA ADMINISTRACIÓN ===
		admin := qr.Group("/admin")
		{
			// Limpiar eventos de escaneo antiguos
			admin.POST("/cleanup-old-events", qrController.CleanupOldScanEvents)

			// Obtener estadísticas del sistema de QR
			admin.GET("/system-stats", qrController.GetSystemQRStats)

			// Exportar datos de trazabilidad
			admin.GET("/export/traceability", qrController.ExportTraceabilityData)

			// Verificar integridad de datos de trazabilidad
			admin.GET("/verify/data-integrity", qrController.VerifyDataIntegrity)

			// Estadísticas generales (ruta alternativa)
			admin.GET("/stats", qrController.GetQRStats)

			// Sincronizar cantidades (ruta alternativa)
			admin.POST("/sync", qrController.SyncBatchAmounts)
		}

		// === WEBHOOKS Y INTEGRACIONES ===
		webhooks := qr.Group("/webhooks")
		{
			// Webhook para sistemas externos cuando se escanea un QR
			webhooks.POST("/scan-event", qrController.HandleScanWebhook)

			// Webhook para notificaciones de consumo
			webhooks.POST("/consume-event", qrController.HandleConsumeWebhook)
		}

		// Grupo de rutas para imágenes (organización mejorada)
		images := qr.Group("/images")
		{
			// Servir imagen QR
			images.GET("/:qrcode", qrController.GetQRImage)

			// Descargar imagen QR
			images.GET("/:qrcode/download", qrController.DownloadQRImage)
		}

		// === RUTAS ESPECIALES PARA COMPATIBILIDAD ===

		// Ruta legacy para compatibilidad con frontend existente
		qr.GET("/scan-legacy/:qrcode", func(c *gin.Context) {
			// Redirigir a la nueva ruta con parámetros de compatibilidad
			qrCode := c.Param("qrcode")
			c.Redirect(302, "/api/qr/scan/"+qrCode+"?scan_purpose=lookup&scan_source=web")
		})

		// === RUTAS DE TRAZABILIDAD EXTENDIDA ===

		// Grupo específico para funcionalidades de trazabilidad
		trace := qr.Group("/trace")
		{
			// Obtener trazabilidad completa
			trace.GET("/:qrcode", qrController.GetCompleteTraceability)

			// Obtener solo historial de escaneos
			trace.GET("/:qrcode/scans", qrController.GetScanHistory)

			// Obtener solo estadísticas
			trace.GET("/:qrcode/stats", qrController.GetScanStatistics)

			// Registrar evento manual
			trace.POST("/register", qrController.RegisterManualScanEvent)
		}

		// === RUTAS DE VERIFICACIÓN Y VALIDACIÓN ===

		// Grupo para verificaciones
		verify := qr.Group("/verify")
		{

			// Verificar disponibilidad
			verify.GET("/:qrcode/availability", qrController.VerifySupplyAvailability)
		}

		// === DOCUMENTACIÓN DE LA API ===

		// Endpoint para obtener documentación de los endpoints de QR
		qr.GET("/docs", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "QR API Documentation with Complete Traceability",
				"version": "2.0",
				"features": []string{
					"Automatic scan event logging",
					"Complete traceability with scan history",
					"Real-time analytics and statistics",
					"Movement tracking and location monitoring",
					"User activity monitoring",
					"Location-based scanning",
					"Multi-device support",
					"Webhook integrations",
					"Data integrity verification",
					"Comprehensive reporting",
				},
				"endpoints": gin.H{
					"basic": []string{
						"GET /scan/:qrcode - Scan QR with automatic logging",
						"GET /validate/:qrcode - Validate QR code",
						"GET /history/:qrcode - Get supply history",
						"GET /details/:qrcode - Get detailed supply info",
					},
					"traceability": []string{
						"GET /traceability/:qrcode - Complete traceability",
						"GET /scan-history/:qrcode - Scan event history",
						"GET /scan-stats/:qrcode - Scan statistics",
						"POST /register-scan - Manual scan registration",
					},
					"consumption": []string{
						"POST /consume - Consume supply (legacy)",
						"POST /consume/individual - Consume individual supply",
						"POST /consume/bulk - Bulk consume supplies",
						"GET /verify/:qrcode - Verify availability",
					},
					"analytics": []string{
						"GET /analytics/scan-summary - Scan activity summary",
						"GET /analytics/top-scanned - Most scanned QRs",
						"GET /analytics/user-activity - User scan activity",
						"GET /analytics/pavilion-activity - Location activity",
						"GET /analytics/movement-patterns - Movement patterns",
					},
					"administration": []string{
						"POST /admin/cleanup-old-events - Clean old events",
						"GET /admin/system-stats - System statistics",
						"GET /admin/export/traceability - Export data",
						"GET /admin/verify/data-integrity - Verify integrity",
					},
					"organized_routes": gin.H{
						"individual": "Routes specific to individual supplies",
						"batch":      "Routes specific to batches",
						"analytics":  "Analytics and reporting routes",
						"admin":      "Administrative routes",
						"webhooks":   "Webhook integration routes",
						"images":     "QR image generation and serving",
						"trace":      "Dedicated traceability routes",
						"verify":     "Verification and validation routes",
					},
				},
				"automatic_logging": gin.H{
					"description": "Every QR scan is automatically logged with complete context",
					"logged_data": []string{
						"User information (RUT, name)",
						"Timestamp (exact scan time)",
						"Location (pavilion, medical center)",
						"Device information (browser, OS, IP)",
						"Scan purpose (lookup, consume, verify, etc.)",
						"QR information snapshot",
						"Scan result (success/error)",
						"Session and request tracking",
					},
				},
				"usage_examples": gin.H{
					"scan_with_context": "GET /qr/scan/SUPPLY_123_abc?user_rut=12345678-9&pavilion_id=1&scan_purpose=consume",
					"get_full_trace":    "GET /qr/traceability/SUPPLY_123_abc",
					"scan_statistics":   "GET /qr/scan-stats/SUPPLY_123_abc",
					"manual_register":   "POST /qr/register-scan {qr_code, user_rut, scan_purpose, ...}",
				},
			})
		})

		// === RUTAS DE AYUDA Y SOPORTE ===

		// Endpoint de salud del sistema QR
		qr.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "healthy",
				"service": "QR Traceability System",
				"version": "2.0",
				"features_enabled": []string{
					"automatic_scan_logging",
					"complete_traceability",
					"real_time_analytics",
					"webhook_integrations",
				},
			})
		})

		// Endpoint para obtener información de configuración
		qr.GET("/config", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"scan_sources":         []string{"web", "mobile", "api", "scanner"},
				"scan_purposes":        []string{"lookup", "consume", "assign", "verify", "inventory_check"},
				"movement_types":       []string{"scan_only", "location_change", "status_change"},
				"supported_qr_types":   []string{"SUPPLY", "BATCH"},
				"logging_enabled":      true,
				"traceability_enabled": true,
			})
		})
	}
}
