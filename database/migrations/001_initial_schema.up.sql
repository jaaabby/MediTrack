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

CREATE TABLE medical_specialty (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    code VARCHAR(50) UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_medical_specialty_name ON medical_specialty(name);
CREATE INDEX idx_medical_specialty_code ON medical_specialty(code);
CREATE INDEX idx_medical_specialty_active ON medical_specialty(is_active);

CREATE TABLE surgery (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    duration FLOAT NOT NULL,
    specialty_id INTEGER REFERENCES medical_specialty(id)
);

CREATE INDEX idx_surgery_specialty ON surgery(specialty_id);

-- =======================
-- TABLA DE PROVEEDORES
-- =======================
CREATE TABLE IF NOT EXISTS supplier_config (
    id SERIAL PRIMARY KEY,
    supplier_name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    notes TEXT
);

COMMENT ON TABLE supplier_config IS 'Tabla de datos de proveedores';
COMMENT ON COLUMN supplier_config.supplier_name IS 'Nombre del proveedor (referenciado por batch.supplier_id → supplier_config.id)';

CREATE OR REPLACE FUNCTION update_supplier_config_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_supplier_config_updated_at
    BEFORE UPDATE ON supplier_config
    FOR EACH ROW
    EXECUTE FUNCTION update_supplier_config_updated_at();

CREATE TABLE supply_code (
    code INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code_supplier INTEGER NOT NULL,
    critical_stock INTEGER NOT NULL
);

CREATE TABLE batch (
    id SERIAL PRIMARY KEY,
    expiration_date DATE NOT NULL,
    amount INTEGER NOT NULL,
    supplier_id INTEGER NOT NULL REFERENCES supplier_config(id),
    store_id INTEGER NOT NULL REFERENCES store(id),
    supply_code INTEGER NOT NULL REFERENCES supply_code(code),
    qr_code VARCHAR(255) UNIQUE,
    surgery_id INTEGER REFERENCES surgery(id),
    location_type VARCHAR(50) NOT NULL DEFAULT 'store' CHECK (location_type IN ('store', 'pavilion')),
    location_id INTEGER NOT NULL DEFAULT 0,
    expiration_alert_days INTEGER NOT NULL DEFAULT 90 CHECK (expiration_alert_days > 0)
);

-- Trigger de seguridad: si location_id no se proporcionó, inicializarlo con store_id
CREATE OR REPLACE FUNCTION set_batch_initial_location()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.location_id = 0 OR NEW.location_id IS NULL THEN
        NEW.location_type := 'store';
        NEW.location_id := NEW.store_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_set_batch_initial_location
    BEFORE INSERT ON batch
    FOR EACH ROW
    EXECUTE FUNCTION set_batch_initial_location();

CREATE TABLE "user" (
    rut VARCHAR(20) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'pabellón', 'encargado de bodega', 'enfermera', 'doctor', 'pavedad')),
    medical_center_id INTEGER NOT NULL REFERENCES medical_center(id),
    pavilion_id INTEGER REFERENCES pavilion(id),
    specialty_id INTEGER REFERENCES medical_specialty(id),
    is_active BOOLEAN DEFAULT TRUE,
    must_change_password BOOLEAN DEFAULT FALSE,
    reset_password_token VARCHAR(255),
    reset_password_expires_at BIGINT,
    totp_secret VARCHAR(64),
    totp_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE INDEX idx_user_email ON "user"(email);
CREATE INDEX idx_user_specialty ON "user"(specialty_id);
CREATE INDEX idx_user_pavilion ON "user"(pavilion_id);
CREATE INDEX idx_user_reset_token ON "user"(reset_password_token);

CREATE TABLE medical_supply (
    id SERIAL PRIMARY KEY,
    batch_id INTEGER NOT NULL REFERENCES batch(id),
    qr_code VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'disponible',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    location_type VARCHAR(50) NOT NULL DEFAULT 'store' CHECK (location_type IN ('store', 'pavilion')),
    location_id INTEGER NOT NULL,
    in_transit BOOLEAN DEFAULT FALSE,
    transfer_date TIMESTAMP,
    transferred_by VARCHAR(20) REFERENCES "user"(rut),
    CONSTRAINT chk_medical_supply_status CHECK (status IN ('disponible', 'pendiente_retiro', 'en_camino_a_pabellon', 'recepcionado', 'consumido', 'en_camino_a_bodega'))
);

CREATE INDEX idx_medical_supply_status ON medical_supply(status);

CREATE TABLE supply_history (
    id SERIAL PRIMARY KEY,
    date_time TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    destination_type VARCHAR(50) NOT NULL,
    destination_id INTEGER NOT NULL,
    medical_supply_id INTEGER NOT NULL REFERENCES medical_supply(id) ON DELETE CASCADE,
    user_rut VARCHAR(20) NOT NULL REFERENCES "user"(rut),
    notes TEXT,
    location TEXT,
    origin_type VARCHAR(50),
    origin_id INTEGER,
    confirmed_by VARCHAR(20) REFERENCES "user"(rut),
    confirmation_date TIMESTAMP,
    transfer_notes TEXT
);

CREATE INDEX idx_supply_history_location ON supply_history(location);
COMMENT ON COLUMN supply_history.location IS 'Ubicación legible del evento (ej: "Pabellon: Pabellón Central 01 (Centro Médico Principal)")';

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

CREATE INDEX idx_batch_history_batch_number ON batch_history(batch_number);

CREATE UNIQUE INDEX idx_batch_qr_code ON batch(qr_code);
CREATE UNIQUE INDEX idx_medical_supply_qr_code ON medical_supply(qr_code);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = EXTRACT(EPOCH FROM NOW());
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_updated_at
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE FUNCTION trg_set_batch_number() RETURNS trigger AS $$
BEGIN
    IF NEW.batch_number IS NULL THEN
        NEW.batch_number := COALESCE(
            NEW.batch_id, 
            (NEW.previous_values->>'id')::int, 
            (NEW.new_values->>'id')::int,
            0
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_batch_number 
    BEFORE INSERT OR UPDATE ON batch_history 
    FOR EACH ROW 
    EXECUTE FUNCTION trg_set_batch_number();

CREATE OR REPLACE FUNCTION log_batch_delete() RETURNS trigger AS $$
DECLARE 
    who text := current_setting('app.current_user', true);
    default_user_rut text := '12345678-9';
    default_user_name text := 'Juan Pérez';
    batch_num int;
BEGIN
    IF who IS NULL OR who = '' THEN
        who := default_user_rut;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM "user" WHERE rut = who) THEN
        who := default_user_rut;
    END IF;
    
    batch_num := OLD.id;
    IF batch_num IS NULL THEN
        batch_num := 0;
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
        NULL,
        default_user_name,
        who,
        batch_num
    );
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_log_batch_delete 
    AFTER DELETE ON batch 
    FOR EACH ROW 
    EXECUTE FUNCTION log_batch_delete();

CREATE TABLE IF NOT EXISTS supply_request (
    id SERIAL PRIMARY KEY,
    request_number VARCHAR(50) NOT NULL UNIQUE,
    pavilion_id INTEGER NOT NULL REFERENCES pavilion(id),
    requested_by VARCHAR(20) NOT NULL REFERENCES "user"(rut),
    requested_by_name VARCHAR(255) NOT NULL,
    request_date TIMESTAMP WITH TIME ZONE NOT NULL,
    surgery_datetime TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pendiente_pavedad',
    notes TEXT,
    -- Campos de Pavedad
    assigned_to VARCHAR(20) REFERENCES "user"(rut),
    assigned_to_name VARCHAR(255),
    assigned_date TIMESTAMP WITH TIME ZONE,
    assigned_by_pavedad VARCHAR(20) REFERENCES "user"(rut),
    assigned_by_pavedad_name VARCHAR(255),
    pavedad_notes TEXT,
    -- Campos de aprobación/rechazo
    approved_by VARCHAR(20) REFERENCES "user"(rut),
    approved_by_name VARCHAR(255),
    approval_date TIMESTAMP WITH TIME ZONE,
    completed_date TIMESTAMP WITH TIME ZONE,
    medical_center_id INTEGER NOT NULL REFERENCES medical_center(id),
    -- Campos de médico responsable
    surgeon_id VARCHAR(20) REFERENCES "user"(rut),
    surgeon_name VARCHAR(255),
    surgery_id INTEGER REFERENCES surgery(id),
    specialty_id INTEGER REFERENCES medical_specialty(id),
    -- Campos para control de retiro
    allow_anyone_to_pickup BOOLEAN DEFAULT TRUE,
    authorized_pickup_rut VARCHAR(20) REFERENCES "user"(rut),
    authorized_pickup_name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Estados posibles:
    -- 'pendiente_pavedad': Doctor crea solicitud
    -- 'asignado_bodega': Pavedad asigna a encargado de bodega
    -- 'en_proceso': Encargado está procesando
    -- 'aprobado', 'rechazado': Decisión del encargado
    -- 'completado', 'cancelado': Estados finales
    -- 'parcialmente_aprobado', 'pendiente_revision': Estados intermedios
    -- 'devuelto': Encargado devuelve items al solicitante para que los modifique
    -- 'devuelto_al_encargado': Doctor reenvía solicitud modificada al encargado
    CONSTRAINT chk_supply_request_status CHECK (status IN ('pendiente_pavedad', 'asignado_bodega', 'en_proceso', 'aprobado', 'rechazado', 'completado', 'cancelado', 'parcialmente_aprobado', 'pendiente_revision', 'devuelto', 'devuelto_al_encargado')),
    CONSTRAINT chk_surgery_datetime_future CHECK (surgery_datetime >= request_date)
);

CREATE INDEX idx_supply_request_status ON supply_request(status);
CREATE INDEX idx_supply_request_pavilion ON supply_request(pavilion_id);
CREATE INDEX idx_supply_request_requested_by ON supply_request(requested_by);
CREATE INDEX idx_supply_request_date ON supply_request(request_date);
CREATE INDEX idx_supply_request_surgery_datetime ON supply_request(surgery_datetime);
CREATE INDEX idx_supply_request_number ON supply_request(request_number);
CREATE INDEX idx_supply_request_assigned_to ON supply_request(assigned_to);
CREATE INDEX idx_supply_request_assigned_by_pavedad ON supply_request(assigned_by_pavedad);
CREATE INDEX idx_supply_request_surgeon ON supply_request(surgeon_id);
CREATE INDEX idx_supply_request_surgery ON supply_request(surgery_id);
CREATE INDEX idx_supply_request_specialty ON supply_request(specialty_id);

CREATE TABLE IF NOT EXISTS supply_request_item (
    id SERIAL PRIMARY KEY,
    supply_request_id INTEGER NOT NULL REFERENCES supply_request(id) ON DELETE CASCADE,
    supply_code INTEGER NOT NULL REFERENCES supply_code(code),
    supply_name VARCHAR(255) NOT NULL,
    quantity_requested INTEGER NOT NULL CHECK (quantity_requested > 0),
    quantity_approved INTEGER CHECK (quantity_approved >= 0),
    quantity_delivered INTEGER NOT NULL DEFAULT 0 CHECK (quantity_delivered >= 0),
    is_pediatric BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Campos para revisión individual por encargado de bodega
    item_status VARCHAR(50) DEFAULT 'pendiente',
    item_notes TEXT,
    reviewed_by VARCHAR(20) REFERENCES "user"(rut),
    reviewed_by_name VARCHAR(255),
    reviewed_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_quantities_logical CHECK (
        quantity_approved IS NULL OR quantity_approved <= quantity_requested
    ),
    CONSTRAINT chk_delivered_not_exceed_approved CHECK (
        quantity_approved IS NULL OR quantity_delivered <= quantity_approved
    )
);

CREATE INDEX idx_supply_request_item_request ON supply_request_item(supply_request_id);
CREATE INDEX idx_supply_request_item_supply_code ON supply_request_item(supply_code);
CREATE INDEX idx_supply_request_item_pediatric ON supply_request_item(is_pediatric);
CREATE INDEX idx_supply_request_item_status ON supply_request_item(item_status);
CREATE INDEX idx_supply_request_item_request_status ON supply_request_item(supply_request_id, item_status);

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
    last_notification_sent TIMESTAMP WITH TIME ZONE,
    notification_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT chk_qr_assignment_status CHECK (status IN ('assigned', 'delivered', 'consumed', 'returned', 'lost')),
    CONSTRAINT chk_delivered_fields_consistency CHECK (
        (delivered_date IS NULL AND delivered_by IS NULL AND delivered_by_name IS NULL) OR
        (delivered_date IS NOT NULL AND delivered_by IS NOT NULL AND delivered_by_name IS NOT NULL)
    )
);

CREATE INDEX idx_qr_assignment_request ON supply_request_qr_assignment(supply_request_id);
CREATE INDEX idx_qr_assignment_item ON supply_request_qr_assignment(supply_request_item_id);
CREATE INDEX idx_qr_assignment_qr_code ON supply_request_qr_assignment(qr_code);
CREATE INDEX idx_qr_assignment_medical_supply ON supply_request_qr_assignment(medical_supply_id);
CREATE INDEX idx_qr_assignment_status ON supply_request_qr_assignment(status);
CREATE INDEX idx_qr_assignment_assigned_date ON supply_request_qr_assignment(assigned_date);
CREATE INDEX idx_qr_assignment_assigned_by ON supply_request_qr_assignment(assigned_by);
CREATE INDEX idx_assignment_notification ON supply_request_qr_assignment(last_notification_sent);

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

CREATE OR REPLACE FUNCTION generate_request_number()
RETURNS VARCHAR(50) AS $$
BEGIN
    RETURN 'SOL-' || TO_CHAR(NOW(), 'YYYYMMDDHH24MISS') || '-' || LPAD(nextval('supply_request_id_seq')::TEXT, 3, '0');
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION is_qr_available_for_assignment(p_qr_code VARCHAR)
RETURNS BOOLEAN AS $$
DECLARE
    supply_exists BOOLEAN := FALSE;
    is_consumed BOOLEAN := FALSE;
    has_active_assignment BOOLEAN := FALSE;
BEGIN
    SELECT EXISTS(SELECT 1 FROM medical_supply WHERE qr_code = p_qr_code) INTO supply_exists;
    
    IF NOT supply_exists THEN
        RETURN FALSE;
    END IF;
    
    SELECT EXISTS(
        SELECT 1 FROM supply_history sh
        JOIN medical_supply ms ON ms.id = sh.medical_supply_id
        WHERE ms.qr_code = p_qr_code AND sh.status = 'consumido'
    ) INTO is_consumed;
    
    IF is_consumed THEN
        RETURN FALSE;
    END IF;
    
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
    scan_source VARCHAR(50) NOT NULL DEFAULT 'web',
    user_agent TEXT,
    ip_address INET,
    
    device_info JSONB,
    browser_info JSONB,
    
    pavilion_id INTEGER REFERENCES pavilion(id),
    pavilion_name VARCHAR(255),
    medical_center_id INTEGER REFERENCES medical_center(id),
    medical_center_name VARCHAR(255),
    
    scan_purpose VARCHAR(100),
    scan_result VARCHAR(50) NOT NULL DEFAULT 'success',
    error_message TEXT,
    
    qr_type VARCHAR(20),
    supply_id INTEGER,
    batch_id INTEGER,
    supply_code INTEGER,
    supply_name VARCHAR(255),
    batch_supplier VARCHAR(255),
    current_status VARCHAR(50),
    
    previous_location VARCHAR(255),
    current_location VARCHAR(255),
    movement_type VARCHAR(50),
    
    session_id VARCHAR(255),
    request_id VARCHAR(255),
    notes TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT chk_qr_scan_event_scan_source CHECK (scan_source IN ('web', 'mobile', 'api', 'scanner', 'pickup_view', 'reception_view', 'transfer_view', 'consumption_view')),
    CONSTRAINT chk_qr_scan_event_scan_result CHECK (scan_result IN ('success', 'error', 'not_found', 'unauthorized')),
    CONSTRAINT chk_qr_scan_event_qr_type CHECK (qr_type IN ('SUPPLY', 'BATCH') OR qr_type IS NULL)
);

CREATE INDEX idx_qr_scan_event_qr_code ON qr_scan_event(qr_code);
CREATE INDEX idx_qr_scan_event_scanned_at ON qr_scan_event(scanned_at DESC);
CREATE INDEX idx_qr_scan_event_scanned_by ON qr_scan_event(scanned_by_rut);
CREATE INDEX idx_qr_scan_event_pavilion ON qr_scan_event(pavilion_id);
CREATE INDEX idx_qr_scan_event_medical_center ON qr_scan_event(medical_center_id);
CREATE INDEX idx_qr_scan_event_scan_result ON qr_scan_event(scan_result);
CREATE INDEX idx_qr_scan_event_supply_id ON qr_scan_event(supply_id);
CREATE INDEX idx_qr_scan_event_batch_id ON qr_scan_event(batch_id);
CREATE INDEX idx_qr_scan_event_qr_time ON qr_scan_event(qr_code, scanned_at DESC);
CREATE INDEX idx_qr_scan_event_session ON qr_scan_event(session_id);

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
    
    u.name AS user_full_name,
    u.email AS user_email,
    
    p.name AS pavilion_full_name,
    
    mc.name AS medical_center_full_name,
    
    ROW_NUMBER() OVER (PARTITION BY qse.qr_code ORDER BY qse.scanned_at) AS scan_sequence,
    
    LAG(qse.scanned_at) OVER (PARTITION BY qse.qr_code ORDER BY qse.scanned_at) AS previous_scan_time,
    
    EXTRACT(EPOCH FROM (qse.scanned_at - LAG(qse.scanned_at) OVER (PARTITION BY qse.qr_code ORDER BY qse.scanned_at))) / 60 AS minutes_since_last_scan

FROM qr_scan_event qse
LEFT JOIN "user" u ON qse.scanned_by_rut = u.rut
LEFT JOIN pavilion p ON qse.pavilion_id = p.id
LEFT JOIN medical_center mc ON qse.medical_center_id = mc.id
ORDER BY qse.qr_code, qse.scanned_at DESC;

CREATE OR REPLACE VIEW v_qr_scan_statistics AS
SELECT 
    qr_code,
    COUNT(*) AS total_scans,
    COUNT(DISTINCT scanned_by_rut) AS unique_scanners,
    COUNT(DISTINCT pavilion_id) AS locations_visited,
    MIN(scanned_at) AS first_scan,
    MAX(scanned_at) AS last_scan,
    EXTRACT(EPOCH FROM (MAX(scanned_at) - MIN(scanned_at))) / 3600 AS hours_in_system,
    
    COUNT(CASE WHEN scan_result = 'success' THEN 1 END) AS successful_scans,
    COUNT(CASE WHEN scan_result = 'error' THEN 1 END) AS error_scans,
    
    COUNT(CASE WHEN scan_source = 'web' THEN 1 END) AS web_scans,
    COUNT(CASE WHEN scan_source = 'mobile' THEN 1 END) AS mobile_scans,
    COUNT(CASE WHEN scan_source = 'api' THEN 1 END) AS api_scans,
    
    COUNT(CASE WHEN scan_purpose = 'consume' THEN 1 END) AS consumption_scans,
    COUNT(CASE WHEN scan_purpose = 'lookup' THEN 1 END) AS lookup_scans,
    COUNT(CASE WHEN scan_purpose = 'verify' THEN 1 END) AS verification_scans

FROM qr_scan_event 
GROUP BY qr_code;

CREATE OR REPLACE FUNCTION cleanup_old_scan_events(days_to_keep INTEGER DEFAULT 90) 
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM qr_scan_event 
    WHERE scanned_at < NOW() - INTERVAL '1 day' * days_to_keep
    AND scan_result = 'success';
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

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

CREATE TABLE IF NOT EXISTS supply_transfer (
    id SERIAL PRIMARY KEY,
    transfer_code VARCHAR(100) NOT NULL UNIQUE,
    qr_code VARCHAR(255) NOT NULL,
    medical_supply_id INTEGER NOT NULL REFERENCES medical_supply(id),
    origin_type VARCHAR(50) NOT NULL CHECK (origin_type IN ('store', 'pavilion')),
    origin_id INTEGER NOT NULL,
    destination_type VARCHAR(50) NOT NULL CHECK (destination_type IN ('store', 'pavilion')),
    destination_id INTEGER NOT NULL,
    sent_by VARCHAR(20) NOT NULL REFERENCES "user"(rut),
    sent_by_name VARCHAR(255) NOT NULL,
    picked_up_by VARCHAR(20) REFERENCES "user"(rut),
    picked_up_by_name VARCHAR(255),
    picked_up_date TIMESTAMP WITH TIME ZONE,
    received_by VARCHAR(20) REFERENCES "user"(rut),
    received_by_name VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'pendiente' CHECK (status IN ('pendiente', 'en_transito', 'recibido', 'rechazado', 'cancelado')),
    transfer_reason TEXT,
    send_date TIMESTAMP WITH TIME ZONE NOT NULL,
    receive_date TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    rejection_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_supply_transfer_qr_code ON supply_transfer(qr_code);
CREATE INDEX idx_supply_transfer_medical_supply ON supply_transfer(medical_supply_id);
CREATE INDEX idx_supply_transfer_picked_up_by ON supply_transfer(picked_up_by);
CREATE INDEX idx_supply_transfer_status ON supply_transfer(status);
CREATE INDEX idx_supply_transfer_origin ON supply_transfer(origin_type, origin_id);
CREATE INDEX idx_supply_transfer_destination ON supply_transfer(destination_type, destination_id);
CREATE INDEX idx_supply_transfer_sent_by ON supply_transfer(sent_by);
CREATE INDEX idx_supply_transfer_send_date ON supply_transfer(send_date);

CREATE TABLE IF NOT EXISTS pavilion_inventory_summary (
    id SERIAL PRIMARY KEY,
    pavilion_id INTEGER NOT NULL REFERENCES pavilion(id),
    batch_id INTEGER NOT NULL REFERENCES batch(id),
    supply_code INTEGER NOT NULL REFERENCES supply_code(code),
    total_received INTEGER NOT NULL DEFAULT 0,
    current_available INTEGER NOT NULL DEFAULT 0,
    total_consumed INTEGER NOT NULL DEFAULT 0,
    total_returned INTEGER NOT NULL DEFAULT 0,
    last_received_date TIMESTAMP WITH TIME ZONE,
    last_consumed_date TIMESTAMP WITH TIME ZONE,
    last_returned_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT uq_pavilion_batch UNIQUE (pavilion_id, batch_id)
);

CREATE INDEX idx_pavilion_inventory_pavilion ON pavilion_inventory_summary(pavilion_id);
CREATE INDEX idx_pavilion_inventory_batch ON pavilion_inventory_summary(batch_id);
CREATE INDEX idx_pavilion_inventory_supply_code ON pavilion_inventory_summary(supply_code);
CREATE INDEX idx_pavilion_inventory_available ON pavilion_inventory_summary(current_available);

CREATE TABLE IF NOT EXISTS store_inventory_summary (
    id SERIAL PRIMARY KEY,
    store_id INTEGER NOT NULL REFERENCES store(id),
    batch_id INTEGER NOT NULL REFERENCES batch(id) UNIQUE,
    supply_code INTEGER NOT NULL REFERENCES supply_code(code),
    surgery_id INTEGER REFERENCES surgery(id),
    original_amount INTEGER NOT NULL,
    current_in_store INTEGER NOT NULL,
    total_transferred_out INTEGER NOT NULL DEFAULT 0,
    total_returned_in INTEGER NOT NULL DEFAULT 0,
    total_consumed_in_store INTEGER NOT NULL DEFAULT 0,
    last_transfer_out_date TIMESTAMP WITH TIME ZONE,
    last_return_in_date TIMESTAMP WITH TIME ZONE,
    last_consumed_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_store_inventory_store ON store_inventory_summary(store_id);
CREATE INDEX idx_store_inventory_batch ON store_inventory_summary(batch_id);
CREATE INDEX idx_store_inventory_supply_code ON store_inventory_summary(supply_code);
CREATE INDEX idx_store_inventory_surgery ON store_inventory_summary(surgery_id);
CREATE INDEX idx_store_inventory_current ON store_inventory_summary(current_in_store);

CREATE OR REPLACE FUNCTION update_supply_transfer_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_supply_transfer_updated_at
    BEFORE UPDATE ON supply_transfer
    FOR EACH ROW
    EXECUTE FUNCTION update_supply_transfer_updated_at();

CREATE OR REPLACE FUNCTION update_pavilion_inventory_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_pavilion_inventory_updated_at
    BEFORE UPDATE ON pavilion_inventory_summary
    FOR EACH ROW
    EXECUTE FUNCTION update_pavilion_inventory_updated_at();

CREATE OR REPLACE FUNCTION update_store_inventory_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_store_inventory_updated_at
    BEFORE UPDATE ON store_inventory_summary
    FOR EACH ROW
    EXECUTE FUNCTION update_store_inventory_updated_at();

-- =======================
-- TABLA DE INSUMOS TÍPICOS POR CIRUGÍA
-- =======================
CREATE TABLE IF NOT EXISTS surgery_typical_supply (
    id SERIAL PRIMARY KEY,
    surgery_id INTEGER NOT NULL REFERENCES surgery(id) ON DELETE CASCADE,
    supply_code INTEGER NOT NULL REFERENCES supply_code(code) ON DELETE CASCADE,
    typical_quantity INTEGER DEFAULT 1,
    is_required BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT uq_surgery_supply UNIQUE (surgery_id, supply_code)
);

CREATE INDEX idx_surgery_typical_supply_surgery ON surgery_typical_supply(surgery_id);
CREATE INDEX idx_surgery_typical_supply_code ON surgery_typical_supply(supply_code);
CREATE INDEX idx_surgery_typical_supply_required ON surgery_typical_supply(is_required);

-- =======================
-- TRIGGERS PARA CONFIGURACIÓN MÉDICA
-- =======================
CREATE OR REPLACE FUNCTION update_medical_specialty_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_medical_specialty_updated_at
    BEFORE UPDATE ON medical_specialty
    FOR EACH ROW
    EXECUTE FUNCTION update_medical_specialty_updated_at();

CREATE OR REPLACE FUNCTION update_surgery_typical_supply_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_surgery_typical_supply_updated_at
    BEFORE UPDATE ON surgery_typical_supply
    FOR EACH ROW
    EXECUTE FUNCTION update_surgery_typical_supply_updated_at();

-- =======================
-- SISTEMA DE CARRITOS DE INSUMOS
-- =======================

-- Tabla principal de carritos
CREATE TABLE IF NOT EXISTS supply_cart (
    id SERIAL PRIMARY KEY,
    supply_request_id INTEGER NOT NULL REFERENCES supply_request(id) ON DELETE CASCADE,
    cart_number VARCHAR(50) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by VARCHAR(20) NOT NULL REFERENCES "user"(rut),
    created_by_name VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    closed_at TIMESTAMP WITH TIME ZONE,
    closed_by VARCHAR(20) REFERENCES "user"(rut),
    closed_by_name VARCHAR(255),
    notes TEXT,
    
    CONSTRAINT chk_supply_cart_status CHECK (status IN ('active', 'closed', 'cancelled')),
    CONSTRAINT uq_supply_cart_request UNIQUE (supply_request_id)
);

CREATE INDEX idx_supply_cart_request ON supply_cart(supply_request_id);
CREATE INDEX idx_supply_cart_status ON supply_cart(status);
CREATE INDEX idx_supply_cart_number ON supply_cart(cart_number);
CREATE INDEX idx_supply_cart_created_at ON supply_cart(created_at DESC);

COMMENT ON TABLE supply_cart IS 'Carritos de insumos generados automáticamente al aprobar solicitudes';
COMMENT ON COLUMN supply_cart.cart_number IS 'Número único del carrito generado automáticamente';
COMMENT ON COLUMN supply_cart.status IS 'Estado del carrito: active, closed, cancelled';

-- Tabla de items del carrito
CREATE TABLE IF NOT EXISTS supply_cart_item (
    id SERIAL PRIMARY KEY,
    supply_cart_id INTEGER NOT NULL REFERENCES supply_cart(id) ON DELETE CASCADE,
    supply_request_qr_assignment_id INTEGER NOT NULL REFERENCES supply_request_qr_assignment(id) ON DELETE CASCADE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    added_by VARCHAR(20) NOT NULL REFERENCES "user"(rut),
    added_by_name VARCHAR(255) NOT NULL,
    removed_at TIMESTAMP WITH TIME ZONE,
    removed_by VARCHAR(20) REFERENCES "user"(rut),
    removed_by_name VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,
    
    CONSTRAINT uq_cart_qr_assignment UNIQUE (supply_cart_id, supply_request_qr_assignment_id)
);

CREATE INDEX idx_supply_cart_item_cart ON supply_cart_item(supply_cart_id);
CREATE INDEX idx_supply_cart_item_assignment ON supply_cart_item(supply_request_qr_assignment_id);
CREATE INDEX idx_supply_cart_item_active ON supply_cart_item(supply_cart_id, is_active);
CREATE INDEX idx_supply_cart_item_added_at ON supply_cart_item(added_at DESC);

COMMENT ON TABLE supply_cart_item IS 'Items individuales del carrito vinculados a asignaciones QR';
COMMENT ON COLUMN supply_cart_item.is_active IS 'Indica si el item está activo en el carrito';

-- Función para generar número de carrito
CREATE OR REPLACE FUNCTION generate_cart_number()
RETURNS VARCHAR(50) AS $$
BEGIN
    RETURN 'CART-' || TO_CHAR(NOW(), 'YYYYMMDDHH24MISS') || '-' || LPAD(nextval('supply_cart_id_seq')::TEXT, 3, '0');
END;
$$ LANGUAGE plpgsql;

-- Trigger para actualizar updated_at en supply_cart
CREATE OR REPLACE FUNCTION update_supply_cart_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_supply_cart_updated_at
    BEFORE UPDATE ON supply_cart
    FOR EACH ROW
    EXECUTE FUNCTION update_supply_cart_updated_at();

-- Vista para obtener carritos con detalles
CREATE OR REPLACE VIEW v_supply_cart_details AS
SELECT 
    sc.id AS cart_id,
    sc.cart_number,
    sc.supply_request_id,
    sr.request_number,
    sc.status AS cart_status,
    sc.created_at AS cart_created_at,
    sc.created_by,
    sc.created_by_name,
    sc.updated_at AS cart_updated_at,
    sc.closed_at,
    sc.closed_by,
    sc.closed_by_name,
    sc.notes AS cart_notes,
    COUNT(sci.id) FILTER (WHERE sci.is_active = TRUE) AS active_items_count,
    COUNT(sci.id) AS total_items_count,
    sr.status AS request_status,
    sr.requested_by_name,
    sr.surgery_datetime,
    sr.pavilion_id
FROM supply_cart sc
JOIN supply_request sr ON sr.id = sc.supply_request_id
LEFT JOIN supply_cart_item sci ON sci.supply_cart_id = sc.id
GROUP BY 
    sc.id, sc.cart_number, sc.supply_request_id, sr.request_number,
    sc.status, sc.created_at, sc.created_by, sc.created_by_name,
    sc.updated_at, sc.closed_at, sc.closed_by, sc.closed_by_name,
    sc.notes, sr.status, sr.requested_by_name, sr.surgery_datetime, sr.pavilion_id;

COMMENT ON VIEW v_supply_cart_details IS 'Vista con detalles completos de carritos incluyendo contadores de items';

-- =======================
-- ACTUALIZACIÓN DE STATUS DE INSUMOS MÉDICOS DESDE HISTORIAL
-- =======================
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