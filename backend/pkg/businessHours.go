package pkg

import "time"

// BusinessHoursConfig configuración de horas laborales
type BusinessHoursConfig struct {
	StartHour int // Hora de inicio (0-23)
	EndHour   int // Hora de fin (0-23)
	AllDays   bool // Si es true, incluye fines de semana
}

// DefaultBusinessHoursConfig configuración por defecto (8:00 - 17:00, todos los días)
func DefaultBusinessHoursConfig() BusinessHoursConfig {
	return BusinessHoursConfig{
		StartHour: 8,
		EndHour:   17,
		AllDays:   true, // Hospital funciona todos los días
	}
}

// CalculateBusinessHours calcula las horas laborales entre dos fechas
// Considera el horario configurado y si incluye fines de semana
func CalculateBusinessHours(start, end time.Time, config BusinessHoursConfig) float64 {
	if end.Before(start) || end.Equal(start) {
		return 0
	}

	var totalHours float64
	current := start

	// Ajustar inicio: si está antes del horario de inicio, empezar desde el horario de inicio
	if current.Hour() < config.StartHour {
		current = time.Date(current.Year(), current.Month(), current.Day(), config.StartHour, 0, 0, 0, current.Location())
	}

	// Ajustar fin: si está después del horario de fin, terminar en el horario de fin del día
	if end.Hour() > config.EndHour || (end.Hour() == config.EndHour && (end.Minute() > 0 || end.Second() > 0)) {
		end = time.Date(end.Year(), end.Month(), end.Day(), config.EndHour, 0, 0, 0, end.Location())
	}

	// Si ambos están en el mismo día
	if current.Year() == end.Year() && current.Month() == end.Month() && current.Day() == end.Day() {
		// Verificar si es fin de semana
		if !config.AllDays {
			weekday := current.Weekday()
			if weekday == time.Saturday || weekday == time.Sunday {
				return 0
			}
		}

		// Calcular horas del mismo día
		if end.After(current) {
			return end.Sub(current).Hours()
		}
		return 0
	}

	// Iterar día por día
	for {
		currentDay := time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
		
		// Verificar si es fin de semana
		if !config.AllDays {
			weekday := currentDay.Weekday()
			if weekday == time.Saturday || weekday == time.Sunday {
				// Saltar al siguiente día a las 00:00
				current = currentDay.AddDate(0, 0, 1)
				current = time.Date(current.Year(), current.Month(), current.Day(), config.StartHour, 0, 0, 0, current.Location())
				continue
			}
		}

		// Verificar si estamos en el último día
		isLastDay := currentDay.Year() == end.Year() && currentDay.Month() == end.Month() && currentDay.Day() == end.Day()

		if isLastDay {
			// Último día: calcular desde current hasta end
			dayStart := time.Date(current.Year(), current.Month(), current.Day(), config.StartHour, 0, 0, 0, current.Location())
			if current.Before(dayStart) {
				current = dayStart
			}
			if end.After(current) {
				totalHours += end.Sub(current).Hours()
			}
			break
		}

		// Día completo o parcial
		dayStart := time.Date(current.Year(), current.Month(), current.Day(), config.StartHour, 0, 0, 0, current.Location())
		dayEnd := time.Date(current.Year(), current.Month(), current.Day(), config.EndHour, 0, 0, 0, current.Location())

		if current.Before(dayStart) {
			// Empezamos antes del horario de inicio: día completo
			totalHours += float64(config.EndHour - config.StartHour)
		} else if current.Before(dayEnd) || current.Equal(dayEnd) {
			// Empezamos durante el horario laboral: horas restantes del día
			totalHours += dayEnd.Sub(current).Hours()
		}
		// Si current >= dayEnd, no se cuenta nada para este día

		// Avanzar al siguiente día
		nextDay := currentDay.AddDate(0, 0, 1)
		current = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), config.StartHour, 0, 0, 0, nextDay.Location())
	}

	return totalHours
}

// HasPassedBusinessHours verifica si han pasado las horas laborales especificadas desde una fecha
func HasPassedBusinessHours(since time.Time, hoursRequired float64, config BusinessHoursConfig) bool {
	now := time.Now()
	hoursElapsed := CalculateBusinessHours(since, now, config)
	return hoursElapsed >= hoursRequired
}

