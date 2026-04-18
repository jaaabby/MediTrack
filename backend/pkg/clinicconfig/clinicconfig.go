package clinicconfig

import (
	"os"

	"gorm.io/gorm"
)

// GetAlertEmail devuelve el correo de alertas con la siguiente prioridad:
//  1. medical_center.alert_email (configurado explícitamente para alertas)
//  2. medical_center.email (correo general del centro médico)
//  3. Variable de entorno ALERT_EMAIL
func GetAlertEmail(db *gorm.DB) string {
	type row struct {
		AlertEmail string
		Email      string
	}
	var r row
	db.Table("medical_center").
		Select("alert_email, email").
		Order("id").
		Limit(1).
		Scan(&r)
	if r.AlertEmail != "" {
		return r.AlertEmail
	}
	if r.Email != "" {
		return r.Email
	}
	return os.Getenv("ALERT_EMAIL")
}
