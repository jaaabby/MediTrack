package services

import (
	"github.com/skip2/go-qrcode"
)

// GenerateQRCodeImage genera una imagen PNG de un QR dado el código y tamaño
func GenerateQRCodeImage(qrCode string, size int) ([]byte, error) {
	return qrcode.Encode(qrCode, qrcode.Medium, size)
}
