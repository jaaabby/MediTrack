-- Migración inicial: Crear esquema de base de datos para MediTrack
-- Fecha: 2025-08-19
-- Descripción: Tablas principales según modelos Go con campos QR incluidos

CREATE TABLE medical_center (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL
);

CREATE TABLE pavilion (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    medical_center_id INTEGER NOT NULL REFERENCES medical_center(id)
);

CREATE TABLE store (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    medical_center_id INTEGER NOT NULL REFERENCES medical_center(id)
);

CREATE TABLE batch (
    id SERIAL PRIMARY KEY,
    expiration_date DATE NOT NULL,
    amount INTEGER NOT NULL,
    supplier VARCHAR(255) NOT NULL,
    store_id INTEGER NOT NULL REFERENCES store(id),
    qr_code VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE supply_code (
    code INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code_supplier INTEGER NOT NULL,
    critical_stock INTEGER NOT NULL
);

CREATE TABLE medical_supply (
    id SERIAL PRIMARY KEY,
    code INTEGER NOT NULL REFERENCES supply_code(code),
    batch_id INTEGER NOT NULL REFERENCES batch(id),
    qr_code VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE surgery (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    duration FLOAT NOT NULL
);

CREATE TABLE "user" (
    rut VARCHAR(20) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL, -- Este campo almacenará el hash de la contraseña
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'pabellón', 'encargado de bodega', 'enfermera', 'doctor')),
    medical_center_id INTEGER NOT NULL REFERENCES medical_center(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE INDEX idx_user_email ON "user"(email);

CREATE TABLE supply_history (
    id SERIAL PRIMARY KEY,
    date_time TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    destination_type VARCHAR(50) NOT NULL,
    destination_id INTEGER NOT NULL,
    medical_supply_id INTEGER NOT NULL REFERENCES medical_supply(id) ON DELETE CASCADE,
    user_rut VARCHAR(20) NOT NULL REFERENCES "user"(rut),
    notes TEXT
);

CREATE TABLE batch_history (
    id SERIAL PRIMARY KEY,
    date_time TIMESTAMP NOT NULL,
    change_details VARCHAR(255) NOT NULL,
    previous_values JSONB NULL,
    new_values JSONB NULL,
    user_name VARCHAR(255) NOT NULL,
    batch_id INTEGER REFERENCES batch(id) ON DELETE SET NULL,
    user_rut VARCHAR(20) NOT NULL REFERENCES "user"(rut),
    batch_number INTEGER NOT NULL
);

-- Crear índice para batch_number para optimizar búsquedas
CREATE INDEX idx_batch_history_batch_number ON batch_history(batch_number);

-- Crear índices únicos para códigos QR
CREATE UNIQUE INDEX idx_batch_qr_code ON batch(qr_code);
CREATE UNIQUE INDEX idx_medical_supply_qr_code ON medical_supply(qr_code);



-- Función para actualizar automáticamente updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = EXTRACT(EPOCH FROM NOW());
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para mantener updated_at actualizado en la tabla user
CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Función para establecer automáticamente batch_number en batch_history
CREATE OR REPLACE FUNCTION trg_set_batch_number() RETURNS trigger AS $$
BEGIN
    -- Solo establecer batch_number si no está ya establecido
    IF NEW.batch_number IS NULL THEN
        NEW.batch_number := COALESCE(
            NEW.batch_id, 
            (NEW.previous_values->>'id')::int, 
            (NEW.new_values->>'id')::int,
            0 -- Valor por defecto si todo lo demás es NULL
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para mantener batch_number actualizado
CREATE TRIGGER set_batch_number 
    BEFORE INSERT OR UPDATE ON batch_history 
    FOR EACH ROW 
    EXECUTE FUNCTION trg_set_batch_number();

-- Función para registrar eliminación de lotes después de que se eliminen
CREATE OR REPLACE FUNCTION log_batch_delete() RETURNS trigger AS $$
DECLARE 
    who text := current_setting('app.current_user', true);
    default_user_rut text := '12345678-9';
    default_user_name text := 'Juan Pérez';
    batch_num int;
BEGIN
    -- Verificar que el usuario existe, si no, usar el por defecto
    IF who IS NULL OR who = '' THEN
        who := default_user_rut;
    END IF;
    
    -- Verificar que el usuario existe en la tabla user
    IF NOT EXISTS (SELECT 1 FROM "user" WHERE rut = who) THEN
        who := default_user_rut;
    END IF;
    
    -- Asegurar que batch_number no sea NULL
    batch_num := OLD.id;
    IF batch_num IS NULL THEN
        batch_num := 0; -- Valor por defecto si OLD.id es NULL
    END IF;
    
    INSERT INTO batch_history(
        date_time, 
        change_details, 
        previous_values, 
        new_values, 
        batch_id, 
        user_name,
        user_rut,
        batch_number
    ) VALUES (
        now(), 
        'Lote eliminado', 
        to_jsonb(OLD), 
        NULL, 
        NULL, -- batch_id será NULL ya que el lote fue eliminado
        default_user_name,
        who,
        batch_num -- batch_number con valor garantizado no NULL
    );
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Trigger para auditar eliminación de lotes automáticamente
-- Este trigger se ejecuta DESPUÉS de eliminar un lote y registra
-- automáticamente la eliminación en batch_history con HARD_DELETE
CREATE TRIGGER trg_log_batch_delete 
    AFTER DELETE ON batch 
    FOR EACH ROW 
    EXECUTE FUNCTION log_batch_delete();

-- Migración: Crear tablas para sistema de solicitudes de insumo con trazabilidad QR
-- Fecha: 2025-01-20
-- Descripción: Agregar funcionalidad de trazabilidad de QR con solicitudes de insumo

-- =======================
-- TABLA PRINCIPAL DE SOLICITUDES
-- =======================

CREATE TABLE IF NOT EXISTS supply_request (
    id SERIAL PRIMARY KEY,
    request_number VARCHAR(50) NOT NULL UNIQUE,
    pavilion_id INTEGER NOT NULL REFERENCES pavilion(id),
    requested_by VARCHAR(20) NOT NULL REFERENCES "user"(rut),
    requested_by_name VARCHAR(255) NOT NULL,
    request_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority VARCHAR(20) NOT NULL DEFAULT 'normal',
    notes TEXT,
    approved_by VARCHAR(20) REFERENCES "user"(rut),
    approved_by_name VARCHAR(255),
    approval_date TIMESTAMP WITH TIME ZONE,
    completed_date TIMESTAMP WITH TIME ZONE,
    medical_center_id INTEGER NOT NULL REFERENCES medical_center(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_supply_request_status CHECK (status IN ('pending', 'approved', 'rejected', 'in_process', 'completed', 'cancelled')),
    CONSTRAINT chk_supply_request_priority CHECK (priority IN ('low', 'normal', 'high', 'critical'))
);

-- Índices para la tabla supply_request
CREATE INDEX idx_supply_request_status ON supply_request(status);
CREATE INDEX idx_supply_request_pavilion ON supply_request(pavilion_id);
CREATE INDEX idx_supply_request_requested_by ON supply_request(requested_by);
CREATE INDEX idx_supply_request_date ON supply_request(request_date);
CREATE INDEX idx_supply_request_number ON supply_request(request_number);

-- =======================
-- TABLA DE ITEMS DE SOLICITUD
-- =======================

CREATE TABLE IF NOT EXISTS supply_request_item (
    id SERIAL PRIMARY KEY,
    supply_request_id INTEGER NOT NULL REFERENCES supply_request(id) ON DELETE CASCADE,
    supply_code INTEGER NOT NULL REFERENCES supply_code(code),
    supply_name VARCHAR(255) NOT NULL,
    quantity_requested INTEGER NOT NULL CHECK (quantity_requested > 0),
    quantity_approved INTEGER CHECK (quantity_approved >= 0),
    quantity_delivered INTEGER NOT NULL DEFAULT 0 CHECK (quantity_delivered >= 0),
    specifications TEXT,
    is_pediatric BOOLEAN NOT NULL DEFAULT FALSE,
    size VARCHAR(50),
    brand VARCHAR(100),
    special_requests TEXT,
    urgency_level VARCHAR(20) NOT NULL DEFAULT 'normal',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_supply_request_item_urgency CHECK (urgency_level IN ('low', 'normal', 'high', 'critical')),
    CONSTRAINT chk_quantities_logical CHECK (
        quantity_approved IS NULL OR quantity_approved <= quantity_requested
    ),
    CONSTRAINT chk_delivered_not_exceed_approved CHECK (
        quantity_approved IS NULL OR quantity_delivered <= quantity_approved
    )
);

-- Índices para la tabla supply_request_item
CREATE INDEX idx_supply_request_item_request ON supply_request_item(supply_request_id);
CREATE INDEX idx_supply_request_item_supply_code ON supply_request_item(supply_code);
CREATE INDEX idx_supply_request_item_urgency ON supply_request_item(urgency_level);
CREATE INDEX idx_supply_request_item_pediatric ON supply_request_item(is_pediatric);

-- =======================
-- TABLA DE ASIGNACIONES QR
-- =======================

CREATE TABLE IF NOT EXISTS supply_request_qr_assignment (
    id SERIAL PRIMARY KEY,
    supply_request_id INTEGER NOT NULL REFERENCES supply_request(id) ON DELETE CASCADE,
    supply_request_item_id INTEGER NOT NULL REFERENCES supply_request_item(id) ON DELETE CASCADE,
    qr_code VARCHAR(100) NOT NULL,
    medical_supply_id INTEGER NOT NULL REFERENCES medical_supply(id),
    assigned_date TIMESTAMP WITH TIME ZONE NOT NULL,
    assigned_by VARCHAR(20) NOT NULL REFERENCES "user"(rut),
    assigned_by_name VARCHAR(255) NOT NULL,
    delivered_date TIMESTAMP WITH TIME ZONE,
    delivered_by VARCHAR(20) REFERENCES "user"(rut),
    delivered_by_name VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'assigned',
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_qr_assignment_status CHECK (status IN ('assigned', 'delivered', 'consumed', 'returned', 'lost')),
    CONSTRAINT chk_delivered_fields_consistency CHECK (
        (delivered_date IS NULL AND delivered_by IS NULL AND delivered_by_name IS NULL) OR
        (delivered_date IS NOT NULL AND delivered_by IS NOT NULL AND delivered_by_name IS NOT NULL)
    ),
    -- Un QR solo puede estar asignado una vez activamente
    CONSTRAINT uq_active_qr_assignment UNIQUE (qr_code, status) DEFERRABLE INITIALLY DEFERRED
);

-- Índices para la tabla supply_request_qr_assignment
CREATE INDEX idx_qr_assignment_request ON supply_request_qr_assignment(supply_request_id);
CREATE INDEX idx_qr_assignment_item ON supply_request_qr_assignment(supply_request_item_id);
CREATE INDEX idx_qr_assignment_qr_code ON supply_request_qr_assignment(qr_code);
CREATE INDEX idx_qr_assignment_medical_supply ON supply_request_qr_assignment(medical_supply_id);
CREATE INDEX idx_qr_assignment_status ON supply_request_qr_assignment(status);
CREATE INDEX idx_qr_assignment_assigned_date ON supply_request_qr_assignment(assigned_date);
CREATE INDEX idx_qr_assignment_assigned_by ON supply_request_qr_assignment(assigned_by);

-- =======================
-- TRIGGERS PARA ACTUALIZACIÓN AUTOMÁTICA
-- =======================

-- Trigger para actualizar updated_at automáticamente en supply_request
CREATE OR REPLACE FUNCTION update_supply_request_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_supply_request_updated_at
    BEFORE UPDATE ON supply_request
    FOR EACH ROW
    EXECUTE FUNCTION update_supply_request_updated_at();

-- Trigger para actualizar updated_at automáticamente en supply_request_item
CREATE OR REPLACE FUNCTION update_supply_request_item_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_supply_request_item_updated_at
    BEFORE UPDATE ON supply_request_item
    FOR EACH ROW
    EXECUTE FUNCTION update_supply_request_item_updated_at();

-- Trigger para actualizar updated_at automáticamente en supply_request_qr_assignment
CREATE OR REPLACE FUNCTION update_qr_assignment_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_qr_assignment_updated_at
    BEFORE UPDATE ON supply_request_qr_assignment
    FOR EACH ROW
    EXECUTE FUNCTION update_qr_assignment_updated_at();

-- =======================
-- FUNCIONES DE UTILIDAD
-- =======================

-- Función para generar número de solicitud único
CREATE OR REPLACE FUNCTION generate_request_number()
RETURNS VARCHAR(50) AS $$
BEGIN
    RETURN 'SOL-' || TO_CHAR(NOW(), 'YYYYMMDDHH24MISS') || '-' || LPAD(nextval('supply_request_id_seq')::TEXT, 3, '0');
END;
$$ LANGUAGE plpgsql;

-- Función para validar que un QR esté disponible para asignación
CREATE OR REPLACE FUNCTION is_qr_available_for_assignment(p_qr_code VARCHAR)
RETURNS BOOLEAN AS $$
DECLARE
    supply_exists BOOLEAN := FALSE;
    is_consumed BOOLEAN := FALSE;
    has_active_assignment BOOLEAN := FALSE;
BEGIN
    -- Verificar que el QR existe en medical_supply
    SELECT EXISTS(SELECT 1 FROM medical_supply WHERE qr_code = p_qr_code) INTO supply_exists;
    
    IF NOT supply_exists THEN
        RETURN FALSE;
    END IF;
    
    -- Verificar que no esté consumido
    SELECT EXISTS(
        SELECT 1 FROM supply_history sh
        JOIN medical_supply ms ON ms.id = sh.medical_supply_id
        WHERE ms.qr_code = p_qr_code AND sh.status = 'consumido'
    ) INTO is_consumed;
    
    IF is_consumed THEN
        RETURN FALSE;
    END IF;
    
    -- Verificar que no tenga asignación activa
    SELECT EXISTS(
        SELECT 1 FROM supply_request_qr_assignment
        WHERE qr_code = p_qr_code 
        AND status NOT IN ('consumed', 'returned')
    ) INTO has_active_assignment;
    
    IF has_active_assignment THEN
        RETURN FALSE;
    END IF;
    
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- =======================
-- VISTAS ÚTILES PARA REPORTING
-- =======================

-- Vista con información completa de solicitudes
CREATE OR REPLACE VIEW v_supply_requests_detail AS
SELECT 
    sr.id,
    sr.request_number,
    sr.status,
    sr.priority,
    sr.request_date,
    sr.requested_by,
    sr.requested_by_name,
    sr.approved_by,
    sr.approved_by_name,
    sr.approval_date,
    sr.completed_date,
    sr.notes,
    p.name AS pavilion_name,
    mc.name AS medical_center_name,
    COUNT(sri.id) AS total_items,
    SUM(sri.quantity_requested) AS total_quantity_requested,
    SUM(COALESCE(sri.quantity_approved, 0)) AS total_quantity_approved,
    SUM(sri.quantity_delivered) AS total_quantity_delivered,
    COUNT(srqa.id) AS total_qr_assignments,
    sr.created_at,
    sr.updated_at
FROM supply_request sr
JOIN pavilion p ON sr.pavilion_id = p.id
JOIN medical_center mc ON sr.medical_center_id = mc.id
LEFT JOIN supply_request_item sri ON sr.id = sri.supply_request_id
LEFT JOIN supply_request_qr_assignment srqa ON sr.id = srqa.supply_request_id
GROUP BY sr.id, p.name, mc.name;

-- Vista de trazabilidad de QR
CREATE OR REPLACE VIEW v_qr_traceability AS
SELECT 
    ms.qr_code,
    ms.id AS medical_supply_id,
    ms.code AS supply_code,
    sc.name AS supply_name,
    b.id AS batch_id,
    b.supplier,
    b.expiration_date,
    b.amount AS batch_amount,
    
    -- Información de solicitud si existe
    srqa.id AS assignment_id,
    srqa.status AS assignment_status,
    srqa.assigned_date,
    srqa.delivered_date,
    sr.id AS request_id,
    sr.request_number,
    sr.status AS request_status,
    p.name AS pavilion_name,
    
    -- Último movimiento del historial
    sh.status AS last_movement_status,
    sh.date_time AS last_movement_date,
    sh.destination_type AS last_destination_type,
    sh.destination_id AS last_destination_id,
    
    -- Indicadores
    CASE 
        WHEN sh.status = 'consumido' THEN 'consumed'
        WHEN srqa.status IS NOT NULL THEN srqa.status
        WHEN sh.status IS NOT NULL THEN sh.status
        ELSE 'available'
    END AS current_status,
    
    b.expiration_date < NOW() AS is_expired,
    b.expiration_date < (NOW() + INTERVAL '30 days') AS expires_soon
    
FROM medical_supply ms
JOIN supply_code sc ON ms.code = sc.code
JOIN batch b ON ms.batch_id = b.id
LEFT JOIN supply_request_qr_assignment srqa ON ms.qr_code = srqa.qr_code 
    AND srqa.status NOT IN ('consumed', 'returned')
LEFT JOIN supply_request sr ON srqa.supply_request_id = sr.id
LEFT JOIN pavilion p ON sr.pavilion_id = p.id
LEFT JOIN LATERAL (
    SELECT status, date_time, destination_type, destination_id
    FROM supply_history sh2 
    WHERE sh2.medical_supply_id = ms.id 
    ORDER BY date_time DESC 
    LIMIT 1
) sh ON TRUE;

-- =======================
-- MIGRACIÓN: EVENTOS DE ESCANEO QR CON TRAZABILIDAD COMPLETA
-- Fecha: 2025-01-21
-- Descripción: Registrar automáticamente cada escaneo de QR para trazabilidad completa
-- =======================

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

-- =======================
-- MIGRACIÓN: ESTADO DE INSUMOS MÉDICOS
-- Fecha: 2025-01-21
-- Descripción: Agregar campo status a medical_supply para almacenar el estado actual del insumo
-- =======================

-- Agregar columna status a medical_supply
ALTER TABLE medical_supply 
ADD COLUMN status VARCHAR(50) NOT NULL DEFAULT 'disponible';

-- Agregar constraint para validar los valores de status
ALTER TABLE medical_supply 
ADD CONSTRAINT chk_medical_supply_status 
CHECK (status IN ('disponible', 'en_camino_a_pabellon', 'recepcionado', 'consumido', 'en_camino_a_bodega'));

-- Crear índice para optimizar consultas por status
CREATE INDEX idx_medical_supply_status ON medical_supply(status);

-- Agregar columna updated_at para tracking de cambios
ALTER TABLE medical_supply 
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Crear trigger para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_medical_supply_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_medical_supply_updated_at
    BEFORE UPDATE ON medical_supply
    FOR EACH ROW
    EXECUTE FUNCTION update_medical_supply_updated_at();

-- NOTA: La función log_medical_supply_status_change y su trigger han sido eliminados
-- para evitar duplicación de registros en supply_history.
-- El historial se maneja exclusivamente a través de la lógica de la aplicación.

-- Migrar estados existentes de supply_history a medical_supply
UPDATE medical_supply 
SET status = (
    SELECT sh.status 
    FROM supply_history sh 
    WHERE sh.medical_supply_id = medical_supply.id 
    ORDER BY sh.date_time DESC 
    LIMIT 1
)
WHERE EXISTS (
    SELECT 1 
    FROM supply_history sh 
    WHERE sh.medical_supply_id = medical_supply.id
);