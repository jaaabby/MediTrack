-- Migración inicial: Crear esquema de base de datos para MediTrack
-- Fecha: 2025-08-16
-- Descripción: Tablas principales según modelos Go

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
    store_id INTEGER NOT NULL REFERENCES store(id)
);

CREATE TABLE supply_code (
    code INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code_supplier INTEGER NOT NULL
);

CREATE TABLE medical_supply (
    id SERIAL PRIMARY KEY,
    code INTEGER NOT NULL REFERENCES supply_code(code),
    batch_id INTEGER NOT NULL REFERENCES batch(id)
);

CREATE TABLE "user" (
    rut VARCHAR(20) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    medical_center_id INTEGER NOT NULL REFERENCES medical_center(id)
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