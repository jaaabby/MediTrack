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
    code_supplier INTEGER NOT NULL
);

CREATE TABLE medical_supply (
    id SERIAL PRIMARY KEY,
    code INTEGER NOT NULL REFERENCES supply_code(code),
    batch_id INTEGER NOT NULL REFERENCES batch(id),
    qr_code VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE "user" (
    rut VARCHAR(20) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL, -- Este campo almacenará el hash de la contraseña
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'pabellón', 'encargado de bodega')),
    medical_center_id INTEGER NOT NULL REFERENCES medical_center(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE supply_history (
    id SERIAL PRIMARY KEY,
    date_time TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    destination_type VARCHAR(50) NOT NULL,
    destination_id INTEGER NOT NULL,
    medical_supply_id INTEGER NOT NULL REFERENCES medical_supply(id) ON DELETE CASCADE,
    user_rut VARCHAR(20) NOT NULL REFERENCES "user"(rut)
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
-- COMENTARIOS Y DOCUMENTACIÓN
-- =======================

COMMENT ON TABLE supply_request IS 'Solicitudes de insumo con trazabilidad QR';
COMMENT ON TABLE supply_request_item IS 'Items individuales dentro de una solicitud de insumo';
COMMENT ON TABLE supply_request_qr_assignment IS 'Asignaciones específicas de códigos QR a items de solicitud';

COMMENT ON COLUMN supply_request.request_number IS 'Número único de solicitud generado automáticamente';
COMMENT ON COLUMN supply_request.status IS 'Estado: pending, approved, rejected, in_process, completed, cancelled';
COMMENT ON COLUMN supply_request.priority IS 'Prioridad: low, normal, high, critical';

COMMENT ON COLUMN supply_request_item.specifications IS 'Especificaciones técnicas del insumo (medidas, tipo pediátrico, etc.)';
COMMENT ON COLUMN supply_request_item.is_pediatric IS 'Indica si el insumo es para uso pediátrico';
COMMENT ON COLUMN supply_request_item.urgency_level IS 'Nivel de urgencia específico del item';

COMMENT ON COLUMN supply_request_qr_assignment.qr_code IS 'Código QR específico asignado al item';
COMMENT ON COLUMN supply_request_qr_assignment.status IS 'Estado de la asignación: assigned, delivered, consumed, returned, lost';

COMMENT ON VIEW v_supply_requests_detail IS 'Vista con información completa de solicitudes incluyendo totales';
COMMENT ON VIEW v_qr_traceability IS 'Vista para trazabilidad completa de códigos QR';

-- =======================
-- PERMISOS (ajustar según necesidades)
-- =======================

-- Otorgar permisos básicos (ajustar según roles en tu aplicación)
-- GRANT SELECT, INSERT, UPDATE ON supply_request TO app_user;
-- GRANT SELECT, INSERT, UPDATE ON supply_request_item TO app_user;
-- GRANT SELECT, INSERT, UPDATE ON supply_request_qr_assignment TO app_user;
-- GRANT SELECT ON v_supply_requests_detail TO app_user;
-- GRANT SELECT ON v_qr_traceability TO app_user;

-- =======================
-- VERIFICACIÓN DE INSTALACIÓN
-- =======================

-- Consulta para verificar que las tablas se crearon correctamente
SELECT 
    schemaname, 
    tablename, 
    tableowner
FROM pg_tables 
WHERE tablename IN ('supply_request', 'supply_request_item', 'supply_request_qr_assignment')
ORDER BY tablename;

-- Consulta para verificar las vistas
SELECT 
    schemaname,
    viewname,
    viewowner
FROM pg_views
WHERE viewname IN ('v_supply_requests_detail', 'v_qr_traceability')
ORDER BY viewname;