package services

import (
	"fmt"
	"log"
	"meditrack/models"
	"time"

	"gorm.io/gorm"
)

type AutomaticConsumptionService struct {
	DB                *gorm.DB
	QRService         *QRService
	SupplyRequestService *SupplyRequestService
	CartService       *CartService
}

func NewAutomaticConsumptionService(db *gorm.DB, qrService *QRService, supplyRequestService *SupplyRequestService) *AutomaticConsumptionService {
	return &AutomaticConsumptionService{
		DB:                db,
		QRService:         qrService,
		SupplyRequestService: supplyRequestService,
		CartService:       nil, // Se establecerá después
	}
}

// SetCartService establece el servicio de carritos
func (s *AutomaticConsumptionService) SetCartService(cartService *CartService) {
	s.CartService = cartService
}

// ProcessAutomaticConsumption procesa el consumo automático de insumos para cirugías completadas
func (s *AutomaticConsumptionService) ProcessAutomaticConsumption() error {
	now := time.Now()
	log.Printf("🔄 Iniciando procesamiento de consumo automático de insumos - %s", now.Format("2006-01-02 15:04:05"))

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

		// Procesar cada asignación
		consumedCount := 0
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

			// Solo consumir insumos en estado "recepcionado"
			// Los que están en "disponible" ya fueron devueltos al momento de la recepción
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

			log.Printf("✅ Insumo %s (QR: %s) cumple condiciones para consumo automático", supply.QRCode, supply.QRCode)

			// Consumir el insumo automáticamente
			// Nota: Agregamos un prefijo especial para identificar consumo automático
			notes := fmt.Sprintf("[CONSUMO_AUTOMATICO] Cirugía completada el %s (solicitud %s). Este insumo puede ser devuelto a bodega si no fue utilizado.", 
				surgeryEndTime.Format("2006-01-02 15:04:05"), request.RequestNumber)
			
			consumptionRequest := QRConsumptionRequest{
				QRCode:          supply.QRCode,
				UserRUT:         "12345678-9", // Administrador del Sistema
				DestinationType: "pavilion", // Usar string literal en lugar de constante
				DestinationID:   supply.LocationID,
				Notes:           notes,
			}

			log.Printf("🔄 Intentando consumir insumo %s (QR: %s) automáticamente", supply.QRCode, supply.QRCode)
			_, err := s.QRService.ConsumeSupplyByQR(consumptionRequest)
			if err != nil {
				errorMsg := fmt.Sprintf("Error consumiendo insumo %s (QR: %s): %v", 
					supply.QRCode, supply.QRCode, err)
				errors = append(errors, errorMsg)
				log.Printf("❌ %s", errorMsg)
				continue
			}

			log.Printf("✅ Insumo %s (QR: %s) consumido exitosamente", supply.QRCode, supply.QRCode)

			// Verificar y cerrar el carrito automáticamente si todos los insumos fueron consumidos
			if s.CartService != nil {
				if err := s.CartService.CheckAndAutoCloseCartForSupply(supply.ID, "12345678-9", "Sistema Automático"); err != nil {
					log.Printf("⚠️  Error verificando cierre automático de carrito para insumo %s: %v", supply.QRCode, err)
				} else {
					log.Printf("✅ Verificado cierre automático de carrito para insumo %s", supply.QRCode)
				}
			}

			consumedCount++
			totalConsumed++
		}

		if consumedCount > 0 {
			totalProcessed++
			log.Printf("✅ Solicitud %s: %d insumo(s) consumido(s) automáticamente", 
				request.RequestNumber, consumedCount)
		}
	}

	log.Printf("✅ Procesamiento completado: %d solicitud(es) procesada(s), %d insumo(s) consumido(s)", 
		totalProcessed, totalConsumed)

	if len(errors) > 0 {
		log.Printf("⚠️  Se encontraron %d error(es) durante el procesamiento", len(errors))
		for _, errMsg := range errors {
			log.Printf("   - %s", errMsg)
		}
	}

	return nil
}

// StartAutomaticConsumptionChecker inicia un verificador que ejecuta el consumo automático periódicamente
func (s *AutomaticConsumptionService) StartAutomaticConsumptionChecker() {
	// Ejecutar cada 5 minutos
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Ejecutar inmediatamente al inicio
	go func() {
		log.Println("🔄 Ejecutando consumo automático inicial...")
		if err := s.ProcessAutomaticConsumption(); err != nil {
			log.Printf("❌ Error en consumo automático inicial: %v", err)
		}
	}()

	// Ejecutar periódicamente
	for range ticker.C {
		if err := s.ProcessAutomaticConsumption(); err != nil {
			log.Printf("❌ Error en consumo automático: %v", err)
		}
	}
}

