package sms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Send envía un SMS al número destino con el mensaje dado usando la API de Twilio.
// Requiere las variables de entorno:
//
//	TWILIO_ACCOUNT_SID  - Account SID de tu cuenta Twilio
//	TWILIO_AUTH_TOKEN   - Auth Token de tu cuenta Twilio
//	TWILIO_FROM_NUMBER  - Número de origen en formato E.164 (ej: +15005550006)
//
// En entorno de desarrollo (ENV=development o ENV vacío), si TWILIO_ACCOUNT_SID
// no está configurado, imprime el mensaje en consola en lugar de enviarlo.
func Send(to, message string) error {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromNumber := os.Getenv("TWILIO_FROM_NUMBER")
	env := os.Getenv("ENV")

	// Modo desarrollo: mostrar el código en consola si Twilio no está configurado
	if accountSID == "" || authToken == "" || fromNumber == "" {
		if env == "production" {
			return fmt.Errorf("configuración de Twilio incompleta: revisa TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN y TWILIO_FROM_NUMBER")
		}
		fmt.Printf("\n📱 [DEV - SMS simulado] Para: %s | Mensaje: %s\n\n", to, message)
		return nil
	}

	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSID)

	formData := url.Values{}
	formData.Set("To", to)
	formData.Set("From", fromNumber)
	formData.Set("Body", message)

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("error creando request Twilio: %w", err)
	}

	req.SetBasicAuth(accountSID, authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando SMS via Twilio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errBody map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return fmt.Errorf("Twilio respondió con status %d: %v", resp.StatusCode, errBody)
	}

	return nil
}
