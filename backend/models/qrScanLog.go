package models

import (
	"database/sql/driver"
	"encoding/json"
	"net"
	"time"
)

// QRScanEvent representa un evento de escaneo de QR para trazabilidad completa
type QRScanEvent struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	QRCode        string    `json:"qr_code" gorm:"not null;index"`
	ScannedAt     time.Time `json:"scanned_at" gorm:"not null;default:now();index:idx_scanned_at,sort:desc"`
	ScannedByRUT  *string   `json:"scanned_by_rut,omitempty" gorm:"index"`
	ScannedByName *string   `json:"scanned_by_name,omitempty"`
	ScanSource    string    `json:"scan_source" gorm:"not null;default:'web'"`
	UserAgent     *string   `json:"user_agent,omitempty"`
	IPAddress     *net.IP   `json:"ip_address,omitempty" gorm:"type:inet"`

	// Información del dispositivo/sesión
	DeviceInfo  *DeviceInfo  `json:"device_info,omitempty" gorm:"type:jsonb"`
	BrowserInfo *BrowserInfo `json:"browser_info,omitempty" gorm:"type:jsonb"`

	// Ubicación/contexto del escaneo
	PavilionID        *int    `json:"pavilion_id,omitempty"`
	PavilionName      *string `json:"pavilion_name,omitempty"`
	MedicalCenterID   *int    `json:"medical_center_id,omitempty"`
	MedicalCenterName *string `json:"medical_center_name,omitempty"`

	// Información adicional del contexto
	ScanPurpose  *string `json:"scan_purpose,omitempty"` // 'lookup', 'consume', 'assign', 'verify', 'inventory_check'
	ScanResult   string  `json:"scan_result" gorm:"not null;default:'success'"`
	ErrorMessage *string `json:"error_message,omitempty"`

	// Datos del QR al momento del escaneo (snapshot)
	QRType        *string `json:"qr_type,omitempty"` // 'SUPPLY', 'BATCH'
	SupplyID      *int    `json:"supply_id,omitempty"`
	BatchID       *int    `json:"batch_id,omitempty"`
	SupplyCode    *int    `json:"supply_code,omitempty"`
	SupplyName    *string `json:"supply_name,omitempty"`
	BatchSupplier *string `json:"batch_supplier,omitempty"`
	CurrentStatus *string `json:"current_status,omitempty"`

	// Información de trazabilidad
	PreviousLocation *string `json:"previous_location,omitempty"`
	CurrentLocation  *string `json:"current_location,omitempty"`
	MovementType     *string `json:"movement_type,omitempty"` // 'scan_only', 'location_change', 'status_change'

	// Metadatos
	SessionID *string `json:"session_id,omitempty"`
	RequestID *string `json:"request_id,omitempty"`
	Notes     *string `json:"notes,omitempty"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relaciones opcionales
	User          *User          `json:"user,omitempty" gorm:"foreignKey:ScannedByRUT;references:RUT"`
	Pavilion      *Pavilion      `json:"pavilion,omitempty" gorm:"foreignKey:PavilionID"`
	MedicalCenter *MedicalCenter `json:"medical_center,omitempty" gorm:"foreignKey:MedicalCenterID"`
}

// DeviceInfo contiene información del dispositivo desde el cual se escaneó
type DeviceInfo struct {
	Platform     string `json:"platform,omitempty"`    // 'Windows', 'Mac', 'Linux', 'iOS', 'Android'
	DeviceType   string `json:"device_type,omitempty"` // 'desktop', 'mobile', 'tablet'
	ScreenSize   string `json:"screen_size,omitempty"` // '1920x1080'
	TouchEnabled bool   `json:"touch_enabled,omitempty"`
	Language     string `json:"language,omitempty"` // 'es-ES', 'en-US'
}

// BrowserInfo contiene información del navegador
type BrowserInfo struct {
	Name       string `json:"name,omitempty"`    // 'Chrome', 'Firefox', 'Safari'
	Version    string `json:"version,omitempty"` // '91.0.4472.124'
	Engine     string `json:"engine,omitempty"`  // 'Blink', 'Gecko', 'WebKit'
	Cookies    bool   `json:"cookies_enabled,omitempty"`
	JavaScript bool   `json:"javascript_enabled,omitempty"`
}

// TableName especifica el nombre de la tabla para GORM
func (QRScanEvent) TableName() string {
	return "qr_scan_event"
}

// Implementar interfaz driver.Valuer para DeviceInfo
func (d DeviceInfo) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Implementar interfaz sql.Scanner para DeviceInfo
func (d *DeviceInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, d)
}

// Implementar interfaz driver.Valuer para BrowserInfo
func (b BrowserInfo) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Implementar interfaz sql.Scanner para BrowserInfo
func (b *BrowserInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, b)
}

// QRCompleteTraceability representa la vista de trazabilidad completa
type QRCompleteTraceability struct {
	QRCode                string     `json:"qr_code"`
	ScannedAt             time.Time  `json:"scanned_at"`
	ScannedByRUT          *string    `json:"scanned_by_rut"`
	ScannedByName         *string    `json:"scanned_by_name"`
	ScanSource            string     `json:"scan_source"`
	PavilionName          *string    `json:"pavilion_name"`
	MedicalCenterName     *string    `json:"medical_center_name"`
	ScanPurpose           *string    `json:"scan_purpose"`
	ScanResult            string     `json:"scan_result"`
	QRType                *string    `json:"qr_type"`
	SupplyName            *string    `json:"supply_name"`
	BatchSupplier         *string    `json:"batch_supplier"`
	CurrentStatus         *string    `json:"current_status"`
	CurrentLocation       *string    `json:"current_location"`
	MovementType          *string    `json:"movement_type"`
	SessionID             *string    `json:"session_id"`
	Notes                 *string    `json:"notes"`
	UserFullName          *string    `json:"user_full_name"`
	UserEmail             *string    `json:"user_email"`
	PavilionFullName      *string    `json:"pavilion_full_name"`
	MedicalCenterFullName *string    `json:"medical_center_full_name"`
	ScanSequence          int        `json:"scan_sequence"`
	PreviousScanTime      *time.Time `json:"previous_scan_time"`
	MinutesSinceLastScan  *float64   `json:"minutes_since_last_scan"`
}

// TableName para la vista de trazabilidad completa
func (QRCompleteTraceability) TableName() string {
	return "v_qr_complete_traceability"
}

// QRScanStatistics representa las estadísticas de escaneo de un QR
type QRScanStatistics struct {
	QRCode            string    `json:"qr_code"`
	TotalScans        int       `json:"total_scans"`
	UniqueScanners    int       `json:"unique_scanners"`
	LocationsVisited  int       `json:"locations_visited"`
	FirstScan         time.Time `json:"first_scan"`
	LastScan          time.Time `json:"last_scan"`
	HoursInSystem     float64   `json:"hours_in_system"`
	SuccessfulScans   int       `json:"successful_scans"`
	ErrorScans        int       `json:"error_scans"`
	WebScans          int       `json:"web_scans"`
	MobileScans       int       `json:"mobile_scans"`
	APIScans          int       `json:"api_scans"`
	ConsumptionScans  int       `json:"consumption_scans"`
	LookupScans       int       `json:"lookup_scans"`
	VerificationScans int       `json:"verification_scans"`
}

// TableName para las estadísticas de escaneo
func (QRScanStatistics) TableName() string {
	return "v_qr_scan_statistics"
}

// Constantes para los valores de escaneo
const (
	// Fuentes de escaneo
	ScanSourceWeb     = "web"
	ScanSourceMobile  = "mobile"
	ScanSourceAPI     = "api"
	ScanSourceScanner = "scanner"

	// Resultados de escaneo
	ScanResultSuccess      = "success"
	ScanResultError        = "error"
	ScanResultNotFound     = "not_found"
	ScanResultUnauthorized = "unauthorized"

	// Propósitos de escaneo
	ScanPurposeLookup               = "lookup"
	ScanPurposeConsume              = "consume"
	ScanPurposeAssign               = "assign"
	ScanPurposeVerify               = "verify"
	ScanPurposeInventoryCheck       = "inventory_check"
	ScanPurposeTransferVerification = "transfer_verification"
	ScanPurposeTransferExecution    = "transfer_execution"
	ScanPurposeTraceabilityView     = "traceability_view"

	// Tipos de movimiento
	MovementTypeScanOnly       = "scan_only"
	MovementTypeLocationChange = "location_change"
	MovementTypeStatusChange   = "status_change"
)
