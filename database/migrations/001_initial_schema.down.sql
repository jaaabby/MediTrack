-- Migración de rollback: Eliminar esquema completo de MediTrack
-- Fecha: 2025-08-19
-- Descripción: Eliminar todas las tablas, estructuras y campos QR

-- =======================
-- ELIMINAR VISTAS PRIMERO
-- =======================
DROP VIEW IF EXISTS v_qr_traceability CASCADE;
DROP VIEW IF EXISTS v_supply_requests_detail CASCADE;
DROP VIEW IF EXISTS v_qr_complete_traceability CASCADE;
DROP VIEW IF EXISTS v_qr_scan_statistics CASCADE;

-- =======================
-- ELIMINAR TRIGGERS
-- =======================
DROP TRIGGER IF EXISTS tr_qr_assignment_updated_at ON supply_request_qr_assignment;
DROP TRIGGER IF EXISTS tr_supply_request_item_updated_at ON supply_request_item;
DROP TRIGGER IF EXISTS tr_supply_request_updated_at ON supply_request;
DROP TRIGGER IF EXISTS trg_log_batch_delete ON batch;
DROP TRIGGER IF EXISTS set_batch_number ON batch_history;
DROP TRIGGER IF EXISTS update_user_updated_at ON "user";
DROP TRIGGER IF EXISTS trg_update_qr_scan_event_updated_at ON qr_scan_event;

-- =======================
-- ELIMINAR FUNCIONES
-- =======================
DROP FUNCTION IF EXISTS is_qr_available_for_assignment(VARCHAR);
DROP FUNCTION IF EXISTS generate_request_number();
DROP FUNCTION IF EXISTS update_qr_assignment_updated_at();
DROP FUNCTION IF EXISTS update_supply_request_item_updated_at();
DROP FUNCTION IF EXISTS update_supply_request_updated_at();
DROP FUNCTION IF EXISTS log_batch_delete();
DROP FUNCTION IF EXISTS trg_set_batch_number();
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS cleanup_old_scan_events(INTEGER);
DROP FUNCTION IF EXISTS update_qr_scan_event_updated_at();

-- =======================
-- ELIMINAR ÍNDICES
-- =======================
-- Índices de supply_request_qr_assignment
DROP INDEX IF EXISTS idx_qr_assignment_assigned_by;
DROP INDEX IF EXISTS idx_qr_assignment_assigned_date;
DROP INDEX IF EXISTS idx_qr_assignment_status;
DROP INDEX IF EXISTS idx_qr_assignment_medical_supply;
DROP INDEX IF EXISTS idx_qr_assignment_qr_code;
DROP INDEX IF EXISTS idx_qr_assignment_item;
DROP INDEX IF EXISTS idx_qr_assignment_request;

-- Índices de supply_request_item
DROP INDEX IF EXISTS idx_supply_request_item_pediatric;
DROP INDEX IF EXISTS idx_supply_request_item_urgency;
DROP INDEX IF EXISTS idx_supply_request_item_supply_code;
DROP INDEX IF EXISTS idx_supply_request_item_request;

-- Índices de supply_request
DROP INDEX IF EXISTS idx_supply_request_number;
DROP INDEX IF EXISTS idx_supply_request_date;
DROP INDEX IF EXISTS idx_supply_request_requested_by;
DROP INDEX IF EXISTS idx_supply_request_pavilion;
DROP INDEX IF EXISTS idx_supply_request_status;

-- Índices de batch_history
DROP INDEX IF EXISTS idx_batch_history_batch_number;

-- Índices de user
DROP INDEX IF EXISTS idx_user_email;

-- Índices únicos de QR
DROP INDEX IF EXISTS idx_batch_qr_code;
DROP INDEX IF EXISTS idx_medical_supply_qr_code;

-- Índices de qr_scan_event
DROP INDEX IF EXISTS idx_qr_scan_event_qr_code;
DROP INDEX IF EXISTS idx_qr_scan_event_scanned_at;
DROP INDEX IF EXISTS idx_qr_scan_event_scanned_by;
DROP INDEX IF EXISTS idx_qr_scan_event_pavilion;
DROP INDEX IF EXISTS idx_qr_scan_event_medical_center;
DROP INDEX IF EXISTS idx_qr_scan_event_scan_result;
DROP INDEX IF EXISTS idx_qr_scan_event_supply_id;
DROP INDEX IF EXISTS idx_qr_scan_event_batch_id;
DROP INDEX IF EXISTS idx_qr_scan_event_qr_time;
DROP INDEX IF EXISTS idx_qr_scan_event_session;

-- =======================
-- ELIMINAR TABLAS EN ORDEN CORRECTO (respetando dependencias)
-- =======================

-- Tablas de solicitudes (más dependientes primero)
DROP TABLE IF EXISTS supply_request_qr_assignment CASCADE;
DROP TABLE IF EXISTS supply_request_item CASCADE;
DROP TABLE IF EXISTS supply_request CASCADE;

-- Tabla de eventos de escaneo QR
DROP TABLE IF EXISTS qr_scan_event CASCADE;

-- Tablas de historial
DROP TABLE IF EXISTS supply_history CASCADE;
DROP TABLE IF EXISTS batch_history CASCADE;

-- Tablas principales
DROP TABLE IF EXISTS medical_supply CASCADE;
DROP TABLE IF EXISTS batch CASCADE;
DROP TABLE IF EXISTS supply_code CASCADE;

-- Tablas de estructura organizacional
DROP TABLE IF EXISTS store CASCADE;
DROP TABLE IF EXISTS pavilion CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS medical_center CASCADE;