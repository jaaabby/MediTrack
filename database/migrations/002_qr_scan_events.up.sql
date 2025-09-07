-- Migración: Crear tabla para eventos de escaneo de QR con trazabilidad completa
-- Fecha: 2025-01-21
-- Descripción: Registrar automáticamente cada escaneo de QR para trazabilidad completa

-- =======================
-- TABLA DE EVENTOS DE ESCANEO QR
-- =======================

CREATE TABLE IF NOT EXISTS qr_scan_event (
    id SERIAL PRIMARY KEY,
    qr_code VARCHAR(255) NOT NULL,
    scanned_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    scanned_by_rut VARCHAR(20) REFERENCES "user"(rut),
    scanned_by_name VARCHAR(255),
    scan_source VARCHAR(50) NOT NULL DEFAULT 'web', -- web, mobile, api
    user_agent TEXT,
    ip_address INET,
    
    -- Información del dispositivo/sesión
    device_info JSONB,
    browser_info JSONB,
    
    -- Ubicación/contexto del escaneo
    pavilion_id INTEGER REFERENCES pavilion(id),
    pavilion_name VARCHAR(255),
    medical_center_id INTEGER REFERENCES medical_center(id),
    medical_center_name VARCHAR(255),
    
    -- Información adicional del contexto
    scan_purpose VARCHAR(100), -- 'lookup', 'consume', 'assign', 'verify', 'inventory_check'
    scan_result VARCHAR(50) NOT NULL DEFAULT 'success', -- 'success', 'error', 'not_found'
    error_message TEXT,
    
    -- Datos del QR al momento del escaneo (snapshot)
    qr_type VARCHAR(20), -- 'SUPPLY', 'BATCH'
    supply_id INTEGER,
    batch_id INTEGER,
    supply_code INTEGER,
    supply_name VARCHAR(255),
    batch_supplier VARCHAR(255),
    current_status VARCHAR(50),
    
    -- Información de trazabilidad
    previous_location VARCHAR(255),
    current_location VARCHAR(255),
    movement_type VARCHAR(50), -- 'scan_only', 'location_change', 'status_change'
    
    -- Metadatos
    session_id VARCHAR(255),
    request_id VARCHAR(255), -- Para rastrear requests específicos
    notes TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_qr_scan_event_scan_source CHECK (scan_source IN ('web', 'mobile', 'api', 'scanner')),
    CONSTRAINT chk_qr_scan_event_scan_result CHECK (scan_result IN ('success', 'error', 'not_found', 'unauthorized')),
    CONSTRAINT chk_qr_scan_event_qr_type CHECK (qr_type IN ('SUPPLY', 'BATCH') OR qr_type IS NULL)
);

-- Índices para optimizar consultas de trazabilidad
CREATE INDEX idx_qr_scan_event_qr_code ON qr_scan_event(qr_code);
CREATE INDEX idx_qr_scan_event_scanned_at ON qr_scan_event(scanned_at DESC);
CREATE INDEX idx_qr_scan_event_scanned_by ON qr_scan_event(scanned_by_rut);
CREATE INDEX idx_qr_scan_event_pavilion ON qr_scan_event(pavilion_id);
CREATE INDEX idx_qr_scan_event_medical_center ON qr_scan_event(medical_center_id);
CREATE INDEX idx_qr_scan_event_scan_result ON qr_scan_event(scan_result);
CREATE INDEX idx_qr_scan_event_supply_id ON qr_scan_event(supply_id);
CREATE INDEX idx_qr_scan_event_batch_id ON qr_scan_event(batch_id);

-- Índice compuesto para consultas de trazabilidad por QR y fecha
CREATE INDEX idx_qr_scan_event_qr_time ON qr_scan_event(qr_code, scanned_at DESC);

-- Índice para consultas por sesión
CREATE INDEX idx_qr_scan_event_session ON qr_scan_event(session_id);

-- =======================
-- VISTA DE TRAZABILIDAD COMPLETA
-- =======================

CREATE OR REPLACE VIEW v_qr_complete_traceability AS
SELECT 
    qse.qr_code,
    qse.scanned_at,
    qse.scanned_by_rut,
    qse.scanned_by_name,
    qse.scan_source,
    qse.pavilion_name,
    qse.medical_center_name,
    qse.scan_purpose,
    qse.scan_result,
    qse.qr_type,
    qse.supply_name,
    qse.batch_supplier,
    qse.current_status,
    qse.current_location,
    qse.movement_type,
    qse.session_id,
    qse.notes,
    
    -- Información del usuario si existe
    u.name AS user_full_name,
    u.email AS user_email,
    
    -- Información del pabellón si existe
    p.name AS pavilion_full_name,
    
    -- Información del centro médico si existe
    mc.name AS medical_center_full_name,
    
    -- Numeración de escaneos para este QR
    ROW_NUMBER() OVER (PARTITION BY qse.qr_code ORDER BY qse.scanned_at) AS scan_sequence,
    
    -- Tiempo entre escaneos
    LAG(qse.scanned_at) OVER (PARTITION BY qse.qr_code ORDER BY qse.scanned_at) AS previous_scan_time,
    
    -- Diferencia de tiempo con el escaneo anterior
    EXTRACT(EPOCH FROM (qse.scanned_at - LAG(qse.scanned_at) OVER (PARTITION BY qse.qr_code ORDER BY qse.scanned_at))) / 60 AS minutes_since_last_scan

FROM qr_scan_event qse
LEFT JOIN "user" u ON qse.scanned_by_rut = u.rut
LEFT JOIN pavilion p ON qse.pavilion_id = p.id
LEFT JOIN medical_center mc ON qse.medical_center_id = mc.id
ORDER BY qse.qr_code, qse.scanned_at DESC;

-- =======================
-- VISTA DE ESTADÍSTICAS DE ESCANEO
-- =======================

CREATE OR REPLACE VIEW v_qr_scan_statistics AS
SELECT 
    qr_code,
    COUNT(*) AS total_scans,
    COUNT(DISTINCT scanned_by_rut) AS unique_scanners,
    COUNT(DISTINCT pavilion_id) AS locations_visited,
    MIN(scanned_at) AS first_scan,
    MAX(scanned_at) AS last_scan,
    EXTRACT(EPOCH FROM (MAX(scanned_at) - MIN(scanned_at))) / 3600 AS hours_in_system,
    
    -- Estadísticas por resultado
    COUNT(CASE WHEN scan_result = 'success' THEN 1 END) AS successful_scans,
    COUNT(CASE WHEN scan_result = 'error' THEN 1 END) AS error_scans,
    
    -- Estadísticas por fuente
    COUNT(CASE WHEN scan_source = 'web' THEN 1 END) AS web_scans,
    COUNT(CASE WHEN scan_source = 'mobile' THEN 1 END) AS mobile_scans,
    COUNT(CASE WHEN scan_source = 'api' THEN 1 END) AS api_scans,
    
    -- Estadísticas por propósito
    COUNT(CASE WHEN scan_purpose = 'consume' THEN 1 END) AS consumption_scans,
    COUNT(CASE WHEN scan_purpose = 'lookup' THEN 1 END) AS lookup_scans,
    COUNT(CASE WHEN scan_purpose = 'verify' THEN 1 END) AS verification_scans

FROM qr_scan_event 
GROUP BY qr_code;

-- =======================
-- FUNCIÓN PARA LIMPIAR EVENTOS ANTIGUOS
-- =======================

CREATE OR REPLACE FUNCTION cleanup_old_scan_events(days_to_keep INTEGER DEFAULT 90) 
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM qr_scan_event 
    WHERE scanned_at < NOW() - INTERVAL '1 day' * days_to_keep
    AND scan_result = 'success'; -- Solo eliminar escaneos exitosos, mantener errores para debugging
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- =======================
-- TRIGGER PARA ACTUALIZAR updated_at
-- =======================

CREATE OR REPLACE FUNCTION update_qr_scan_event_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_qr_scan_event_updated_at
    BEFORE UPDATE ON qr_scan_event
    FOR EACH ROW
    EXECUTE FUNCTION update_qr_scan_event_updated_at();