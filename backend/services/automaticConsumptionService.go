package services

import (
	"fmt"
	"log"
	"meditrack/mailer"
	"meditrack/models"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type AutomaticConsumptionService struct {
	DB                   *gorm.DB
	QRService            *QRService
	SupplyRequestService *SupplyRequestService
	CartService          *CartService
}

func NewAutomaticConsumptionService(db *gorm.DB, qrService *QRService, supplyRequestService *SupplyRequestService) *AutomaticConsumptionService {
	return &AutomaticConsumptionService{
		DB:                   db,
		QRService:            qrService,
		SupplyRequestService: supplyRequestService,
		CartService:          nil, // Se establecerá después
	}
}

// SetCartService establece el servicio de carritos
func (s *AutomaticConsumptionService) SetCartService(cartService *CartService) {
	s.CartService = cartService
}

// ProcessAutomaticConsumption verifica cirugías completadas y envía notificaciones para insumos pendientes
func (s *AutomaticConsumptionService) ProcessAutomaticConsumption() error {
	now := time.Now()
	log.Printf("🔄 Verificando cirugías completadas y notificaciones pendientes - %s", now.Format("2006-01-02 15:04:05"))

	// Buscar solicitudes con cirugías que ya deberían haber terminado
	// Usar una consulta SQL más eficiente que filtre directamente por fecha
	var requests []models.SupplyRequest

	// Obtener todas las solicitudes con cirugías y calcular en memoria cuáles ya terminaron
	// Nota: Esto podría optimizarse con una consulta SQL más compleja, pero por simplicidad
	// lo hacemos así para asegurar que tenemos la información completa de la cirugía
	err := s.DB.Where("surgery_id IS NOT NULL").
		Where("status IN ?", []string{
			models.RequestStatusApproved,
			models.RequestStatusInProcess,
			models.RequestStatusCompleted,
		}).
		Preload("Surgery").
		Find(&requests).Error
	if err != nil {
		return fmt.Errorf("error obteniendo solicitudes: %v", err)
	}

	log.Printf("📋 Encontradas %d solicitudes con cirugías asociadas", len(requests))

	var totalProcessed int
	var totalConsumed int
	var errors []string

	for _, request := range requests {
		if request.Surgery == nil {
			log.Printf("⚠️  Solicitud %d: sin cirugía asociada", request.ID)
			continue
		}

		// Calcular cuándo debería terminar la cirugía
		surgeryEndTime := request.SurgeryDatetime.Add(time.Duration(request.Surgery.Duration * float64(time.Hour)))

		log.Printf("🔍 Verificando solicitud %s: cirugía programada %s, duración %.2f horas, debería terminar %s, ahora es %s",
			request.RequestNumber,
			request.SurgeryDatetime.Format("2006-01-02 15:04:05"),
			request.Surgery.Duration,
			surgeryEndTime.Format("2006-01-02 15:04:05"),
			now.Format("2006-01-02 15:04:05"))

		// Solo procesar si la cirugía ya debería haber terminado
		if now.Before(surgeryEndTime) {
			log.Printf("⏳ Solicitud %s: cirugía aún no ha terminado (termina en %s)",
				request.RequestNumber, surgeryEndTime.Format("2006-01-02 15:04:05"))
			continue
		}

		// El estado ya fue filtrado en la consulta SQL, pero verificamos por seguridad

		// Buscar insumos recepcionados en el pabellón de esta solicitud
		// En lugar de buscar asignaciones con estado "delivered", buscamos directamente
		// insumos que estén recepcionados en el pabellón y asociados a esta solicitud
		var assignments []models.SupplyRequestQRAssignment
		err := s.DB.Where("supply_request_id = ?", request.ID).
			Where("status IN ?", []string{
				models.AssignmentStatusAssigned,
				models.AssignmentStatusDelivered,
			}).
			Preload("MedicalSupply").
			Find(&assignments).Error

		if err != nil {
			errors = append(errors, fmt.Sprintf("Error obteniendo asignaciones para solicitud %d: %v", request.ID, err))
			log.Printf("❌ Error obteniendo asignaciones para solicitud %d: %v", request.ID, err)
			continue
		}

		log.Printf("📦 Solicitud %s: encontradas %d asignaciones asociadas", request.RequestNumber, len(assignments))

		if len(assignments) == 0 {
			log.Printf("⚠️  Solicitud %s: no hay asignaciones para procesar", request.RequestNumber)
			continue
		}

		// Procesar cada asignación - SOLO NOTIFICAR, NO CONSUMIR
		notifiedCount := 0
		for _, assignment := range assignments {
			// Cargar el insumo médico
			var supply models.MedicalSupply
			// Verificar si el Preload cargó el insumo (si ID es 0, no se cargó)
			if assignment.MedicalSupply.ID == 0 {
				// Si no está precargado, cargarlo manualmente
				if err := s.DB.First(&supply, assignment.MedicalSupplyID).Error; err != nil {
					errors = append(errors, fmt.Sprintf("Error cargando insumo %d: %v", assignment.MedicalSupplyID, err))
					continue
				}
			} else {
				supply = assignment.MedicalSupply
			}

			// Verificar que el insumo no esté ya consumido o devuelto
			if supply.Status == models.StatusConsumed || supply.Status == models.StatusAvailable {
				continue
			}

			// Solo notificar insumos en estado "recepcionado"
			if supply.Status != models.StatusReceived {
				log.Printf("⚠️  Insumo %s (QR: %s): estado %s, esperado %s",
					supply.QRCode, supply.QRCode, supply.Status, models.StatusReceived)
				continue
			}

			if supply.LocationType != models.SupplyLocationPavilion {
				log.Printf("⚠️  Insumo %s (QR: %s): ubicación tipo %s, esperado %s",
					supply.QRCode, supply.QRCode, supply.LocationType, models.SupplyLocationPavilion)
				continue
			}

			// Verificar que el insumo esté en el pabellón correcto (el de la solicitud)
			if supply.LocationID != request.PavilionID {
				log.Printf("⚠️  Insumo %s (QR: %s): pabellón %d, esperado %d",
					supply.QRCode, supply.QRCode, supply.LocationID, request.PavilionID)
				continue
			}

			// Verificar si debemos enviar notificación según el intervalo configurado
			shouldNotify, isFirstNotification := s.shouldSendNotification(&assignment, surgeryEndTime)
			if !shouldNotify {
				log.Printf("⏭️  Insumo %s: notificación ya enviada recientemente", supply.QRCode)
				continue
			}

			if isFirstNotification {
				log.Printf("📧 Insumo %s (QR: %s) requiere atención - cirugía completada (primera notificación)", supply.QRCode, supply.QRCode)
			} else {
				log.Printf("🔔 Insumo %s (QR: %s) - enviando recordatorio (notificación #%d)", supply.QRCode, assignment.NotificationCount+1)
			}

			// Enviar notificación por correo en lugar de consumir automáticamente
			if err := s.sendUnconsumedSupplyAlert(request, supply, surgeryEndTime, assignment.NotificationCount+1); err != nil {
				errorMsg := fmt.Sprintf("Error enviando notificación para insumo %s: %v", supply.QRCode, err)
				errors = append(errors, errorMsg)
				log.Printf("❌ %s", errorMsg)
				continue
			}

			// Actualizar el registro de notificación
			now := time.Now()
			assignment.LastNotificationSent = &now
			assignment.NotificationCount++
			if err := s.DB.Save(&assignment).Error; err != nil {
				log.Printf("⚠️  Error actualizando registro de notificación: %v", err)
			}

			log.Printf("✅ Notificación enviada para insumo %s (QR: %s)", supply.QRCode, supply.QRCode)
			notifiedCount++
			totalConsumed++
		}

		if notifiedCount > 0 {
			totalProcessed++
			log.Printf("📨 Solicitud %s: %d notificación(es) enviada(s)",
				request.RequestNumber, notifiedCount)
		}
	}

	log.Printf("✅ Procesamiento completado: %d solicitud(es) procesada(s), %d notificación(es) enviada(s)",
		totalProcessed, totalConsumed)

	if len(errors) > 0 {
		log.Printf("⚠️  Se encontraron %d error(es) durante el procesamiento", len(errors))
		for _, errMsg := range errors {
			log.Printf("   - %s", errMsg)
		}
	}

	return nil
}

// shouldSendNotification determina si se debe enviar una notificación según el intervalo configurado
// Retorna (debeEnviar, esPrimeraNotificacion)
func (s *AutomaticConsumptionService) shouldSendNotification(assignment *models.SupplyRequestQRAssignment, surgeryEndTime time.Time) (bool, bool) {
	// Si nunca se ha enviado notificación, enviar inmediatamente después de que termine la cirugía
	if assignment.LastNotificationSent == nil {
		return true, true
	}

	// Obtener el intervalo de reenvío desde variable de entorno (en horas)
	// Por defecto: 2 horas
	resendIntervalHours := 2.0
	if envInterval := os.Getenv("NOTIFICATION_RESEND_INTERVAL_HOURS"); envInterval != "" {
		if parsed, err := strconv.ParseFloat(envInterval, 64); err == nil && parsed > 0 {
			resendIntervalHours = parsed
		}
	}

	// Calcular cuánto tiempo ha pasado desde la última notificación
	timeSinceLastNotification := time.Since(*assignment.LastNotificationSent)
	resendInterval := time.Duration(resendIntervalHours * float64(time.Hour))

	// Enviar si ha pasado el intervalo configurado
	shouldSend := timeSinceLastNotification >= resendInterval

	return shouldSend, false
}

// sendUnconsumedSupplyAlert envía una notificación por correo para un insumo no consumido
func (s *AutomaticConsumptionService) sendUnconsumedSupplyAlert(request models.SupplyRequest, supply models.MedicalSupply, surgeryEndTime time.Time, notificationNumber int) error {
	// Cargar el usuario solicitante
	var requester models.User
	if err := s.DB.First(&requester, "rut = ?", request.RequestedBy).Error; err != nil {
		return fmt.Errorf("error cargando usuario solicitante: %v", err)
	}

	// Cargar información del batch con el código de insumo
	var batch models.Batch
	if err := s.DB.First(&batch, supply.BatchID).Error; err != nil {
		log.Printf("⚠️  No se pudo cargar información del lote: %v", err)
	}

	// Cargar información del código de insumo
	var supplyCode models.SupplyCode
	if err := s.DB.First(&supplyCode, supply.Code).Error; err != nil {
		log.Printf("⚠️  No se pudo cargar información del código de insumo: %v", err)
	}

	// Calcular horas transcurridas desde el fin de la cirugía
	hoursElapsed := time.Since(surgeryEndTime).Hours()

	// Preparar datos para el template
	data := struct {
		SupplyID           int
		SupplyName         string
		SupplyCode         int
		QRCode             string
		BatchID            int
		ReceivedAt         string
		HoursElapsed       string
		Date               string
		RequestNumber      string
		SurgeryDate        string
		NotificationNumber int
		IsReminder         bool
	}{
		SupplyID:           supply.ID,
		SupplyName:         supplyCode.Name,
		SupplyCode:         supply.Code,
		QRCode:             supply.QRCode,
		BatchID:            batch.ID,
		ReceivedAt:         supply.UpdatedAt.Format("2006-01-02 15:04:05"),
		HoursElapsed:       fmt.Sprintf("%.2f", hoursElapsed),
		Date:               time.Now().Format("2006-01-02 15:04:05"),
		RequestNumber:      request.RequestNumber,
		SurgeryDate:        request.SurgeryDatetime.Format("2006-01-02 15:04:05"),
		NotificationNumber: notificationNumber,
		IsReminder:         notificationNumber > 1,
	}

	templatePath := filepath.Join("mailer", "templates", "unconsumed_supply_alert.html")
	subject := "🚨 Insumo Pendiente de Consumo/Devolución - " + request.RequestNumber
	if notificationNumber > 1 {
		subject = fmt.Sprintf("🔔 RECORDATORIO #%d - Insumo Pendiente - %s", notificationNumber, request.RequestNumber)
	}
	mailReq := mailer.NewRequest([]string{requester.Email}, subject)

	if err := mailReq.SendMailSkipTLS(templatePath, data); err != nil {
		return fmt.Errorf("error enviando correo: %v", err)
	}

	return nil
}

// StartAutomaticConsumptionChecker inicia un verificador que verifica cirugías completadas y envía notificaciones
func (s *AutomaticConsumptionService) StartAutomaticConsumptionChecker() {
	// Ejecutar cada 5 minutos
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Ejecutar inmediatamente al inicio
	go func() {
		log.Println("🔄 Ejecutando verificación inicial de cirugías completadas...")
		if err := s.ProcessAutomaticConsumption(); err != nil {
			log.Printf("❌ Error en verificación inicial: %v", err)
		}
	}()

	// Ejecutar periódicamente
	for range ticker.C {
		if err := s.ProcessAutomaticConsumption(); err != nil {
			log.Printf("❌ Error en verificación de cirugías: %v", err)
		}
	}
}
