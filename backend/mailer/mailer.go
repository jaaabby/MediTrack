package mailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"mime"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/mail.v2"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

var (
	email, pass, server, port string
)

func Setup() {
	email = os.Getenv("EMAIL_DIR")
	pass = os.Getenv("EMAIL_PASS")
	server = os.Getenv("EMAIL_SERVER")
	port = os.Getenv("EMAIL_PORT")
}

func NewRequest(to []string, subject string) *Request {
	return &Request{
		from:    email,
		to:      to,
		subject: subject,
	}
}

func (r *Request) parseTemplate(templatePath string, data interface{}) error {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	if err := t.Execute(buffer, data); err != nil {
		return err
	}

	r.body = buffer.String()
	return nil
}

func (r *Request) sendMail() error {
	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%s", server, port)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", email, pass, server), email, r.to, []byte(body)); err != nil {
		return err
	}
	return nil
}

func (r *Request) Send(templateName string, items interface{}) error {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		return err
	}
	if err := r.sendMail(); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *Request) SendMailSkipTLS(templateName string, items interface{}) error {

	err := r.parseTemplate(templateName, items)
	if err != nil {
		fmt.Printf("Error parseando template: %s\n", err.Error())
		return err
	}

	m := mail.NewMessage()

	m.SetHeader("From", email)
	// Enviar a todos los destinatarios (incluyendo el correo de prueba)
	m.SetHeader("To", r.to...)
	m.SetHeader("Subject", r.subject)
	m.SetBody("text/html", r.body)

	var puerto int
	puerto, err = strconv.Atoi(port)
	if err != nil {
		fmt.Printf("No se pudo parsear el puerto: %s\n", err.Error())
		return err
	}
	
	// Verificar que las variables de entorno estén configuradas
	if email == "" || pass == "" || server == "" || port == "" {
		fmt.Printf("ERROR: Variables de entorno de email no configuradas. EMAIL_DIR=%s, EMAIL_SERVER=%s, EMAIL_PORT=%s\n", 
			email, server, port)
		return fmt.Errorf("variables de entorno de email no configuradas")
	}
	
	fmt.Printf("Intentando enviar correo a: %v desde: %s usando servidor: %s:%s\n", r.to, email, server, port)
	
	d := mail.NewDialer(server, puerto, email, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // Habilitar TLS con skip verify

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("ERROR al enviar correo: %s\n", err.Error())
		return err
	}

	fmt.Printf("Correo enviado exitosamente a: %v\n", r.to)
	return nil
}

// SendMailWithAttachment envía un correo con un archivo adjunto
func (r *Request) SendMailWithAttachment(templateName string, data interface{}, attachmentName string, attachmentData []byte) error {
	err := r.parseTemplate(templateName, data)
	if err != nil {
		return err
	}

	ext := filepath.Ext(attachmentName)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	m := mail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", r.to...) // Enviar a múltiples destinatarios si es necesario
	m.SetHeader("Subject", r.subject)
	m.SetBody("text/html", r.body)

	m.AttachReader(
		attachmentName,
		bytes.NewReader(attachmentData), // Convertir buffer en un lector
		mail.SetHeader(map[string][]string{
			"Content-Type": {mimeType},
		}),
	)

	puerto, err := strconv.Atoi(port)
	if err != nil {
		fmt.Printf("No se pudo parsear el puerto: %s\n", err.Error())
		return err
	}
	d := mail.NewDialer(server, puerto, email, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // Opcional, si el servidor no tiene SSL válido

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("No se pudo enviar el email con adjunto: %s\n", err.Error())
		return err
	}

	fmt.Println("Correo enviado correctamente con adjunto:", attachmentName)
	return nil
}
