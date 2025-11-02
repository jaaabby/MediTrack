-- Script de rollback completo: Eliminar esquema completo de MediTrack
-- Fecha: 2025-01-20
-- Descripción: Eliminar todas las tablas, vistas, funciones, triggers e índices del sistema

-- =======================
-- ELIMINAR VISTAS
-- =======================
DROP VIEW IF EXISTS v_qr_scan_statistics CASCADE;
DROP VIEW IF EXISTS v_qr_complete_traceability CASCADE;
DROP VIEW IF EXISTS v_qr_traceability CASCADE;
DROP VIEW IF EXISTS v_supply_requests_detail CASCADE;

-- =======================
-- ELIMINAR TRIGGERS
-- =======================
DROP TRIGGER IF EXISTS trg_update_store_inventory_updated_at ON store_inventory_summary;
DROP TRIGGER IF EXISTS trg_update_pavilion_inventory_updated_at ON pavilion_inventory_summary;
DROP TRIGGER IF EXISTS trg_update_supply_transfer_updated_at ON supply_transfer;
DROP TRIGGER IF EXISTS trg_update_medical_supply_updated_at ON medical_supply;
DROP TRIGGER IF EXISTS trg_log_medical_supply_status_change ON medical_supply;
DROP TRIGGER IF EXISTS trg_update_qr_scan_event_updated_at ON qr_scan_event;
DROP TRIGGER IF EXISTS tr_qr_assignment_updated_at ON supply_request_qr_assignment;
DROP TRIGGER IF EXISTS tr_supply_request_item_updated_at ON supply_request_item;
DROP TRIGGER IF EXISTS tr_supply_request_updated_at ON supply_request;
DROP TRIGGER IF EXISTS trg_log_batch_delete ON batch;
DROP TRIGGER IF EXISTS set_batch_number ON batch_history;
DROP TRIGGER IF EXISTS update_user_updated_at ON "user";
DROP TRIGGER IF EXISTS trg_update_doctor_info_updated_at ON doctor_info;
DROP TRIGGER IF EXISTS trg_update_surgery_typical_supply_updated_at ON surgery_typical_supply;
DROP TRIGGER IF EXISTS trg_update_medical_specialty_updated_at ON medical_specialty;
DROP TRIGGER IF EXISTS trg_update_supplier_config_updated_at ON supplier_config;

-- =======================
-- ELIMINAR FUNCIONES
-- =======================
DROP FUNCTION IF EXISTS update_store_inventory_updated_at();
DROP FUNCTION IF EXISTS update_pavilion_inventory_updated_at();
DROP FUNCTION IF EXISTS update_supply_transfer_updated_at();
DROP FUNCTION IF EXISTS update_medical_supply_updated_at();
DROP FUNCTION IF EXISTS log_medical_supply_status_change();
DROP FUNCTION IF EXISTS cleanup_old_scan_events(INTEGER);
DROP FUNCTION IF EXISTS update_qr_scan_event_updated_at();
DROP FUNCTION IF EXISTS is_qr_available_for_assignment(VARCHAR);
DROP FUNCTION IF EXISTS generate_request_number();
DROP FUNCTION IF EXISTS update_qr_assignment_updated_at();
DROP FUNCTION IF EXISTS update_supply_request_item_updated_at();
DROP FUNCTION IF EXISTS update_supply_request_updated_at();
DROP FUNCTION IF EXISTS log_batch_delete();
DROP FUNCTION IF EXISTS trg_set_batch_number();
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS update_doctor_info_updated_at();
DROP FUNCTION IF EXISTS update_surgery_typical_supply_updated_at();
DROP FUNCTION IF EXISTS update_medical_specialty_updated_at();
DROP FUNCTION IF EXISTS update_supplier_config_updated_at();

-- =======================
-- ELIMINAR ÍNDICES
-- =======================

-- Índices de supplier_config
DROP INDEX IF EXISTS idx_supplier_config_alert_days;

-- Índices de store_inventory_summary
DROP INDEX IF EXISTS idx_store_inventory_current;
DROP INDEX IF EXISTS idx_store_inventory_surgery;
DROP INDEX IF EXISTS idx_store_inventory_supply_code;
DROP INDEX IF EXISTS idx_store_inventory_batch;
DROP INDEX IF EXISTS idx_store_inventory_store;

-- Índices de pavilion_inventory_summary
DROP INDEX IF EXISTS idx_pavilion_inventory_available;
DROP INDEX IF EXISTS idx_pavilion_inventory_supply_code;
DROP INDEX IF EXISTS idx_pavilion_inventory_batch;
DROP INDEX IF EXISTS idx_pavilion_inventory_pavilion;

-- Índices de supply_transfer
DROP INDEX IF EXISTS idx_supply_transfer_send_date;
DROP INDEX IF EXISTS idx_supply_transfer_sent_by;
DROP INDEX IF EXISTS idx_supply_transfer_destination;
DROP INDEX IF EXISTS idx_supply_transfer_origin;
DROP INDEX IF EXISTS idx_supply_transfer_status;
DROP INDEX IF EXISTS idx_supply_transfer_medical_supply;
DROP INDEX IF EXISTS idx_supply_transfer_qr_code;

-- Índices de medical_supply
DROP INDEX IF EXISTS idx_medical_supply_status;

-- Índices de qr_scan_event
DROP INDEX IF EXISTS idx_qr_scan_event_session;
DROP INDEX IF EXISTS idx_qr_scan_event_qr_time;
DROP INDEX IF EXISTS idx_qr_scan_event_batch_id;
DROP INDEX IF EXISTS idx_qr_scan_event_supply_id;
DROP INDEX IF EXISTS idx_qr_scan_event_scan_result;
DROP INDEX IF EXISTS idx_qr_scan_event_medical_center;
DROP INDEX IF EXISTS idx_qr_scan_event_pavilion;
DROP INDEX IF EXISTS idx_qr_scan_event_scanned_by;
DROP INDEX IF EXISTS idx_qr_scan_event_scanned_at;
DROP INDEX IF EXISTS idx_qr_scan_event_qr_code;

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
DROP INDEX IF EXISTS idx_supply_request_item_supply_code;
DROP INDEX IF EXISTS idx_supply_request_item_request;

-- Índices de supply_request
DROP INDEX IF EXISTS idx_supply_request_specialty;
DROP INDEX IF EXISTS idx_supply_request_surgery;
DROP INDEX IF EXISTS idx_supply_request_surgeon;
DROP INDEX IF EXISTS idx_supply_request_number;
DROP INDEX IF EXISTS idx_supply_request_date;
DROP INDEX IF EXISTS idx_supply_request_requested_by;
DROP INDEX IF EXISTS idx_supply_request_pavilion;
DROP INDEX IF EXISTS idx_supply_request_status;

-- Índices de surgery
DROP INDEX IF EXISTS idx_surgery_specialty;

-- Índices de doctor_info
DROP INDEX IF EXISTS idx_doctor_info_license;
DROP INDEX IF EXISTS idx_doctor_info_available;
DROP INDEX IF EXISTS idx_doctor_info_specialty;

-- Índices de surgery_typical_supply
DROP INDEX IF EXISTS idx_surgery_typical_supply_required;
DROP INDEX IF EXISTS idx_surgery_typical_supply_code;
DROP INDEX IF EXISTS idx_surgery_typical_supply_surgery;

-- Índices de medical_specialty
DROP INDEX IF EXISTS idx_medical_specialty_active;
DROP INDEX IF EXISTS idx_medical_specialty_code;
DROP INDEX IF EXISTS idx_medical_specialty_name;

-- Índices de batch_history
DROP INDEX IF EXISTS idx_batch_history_batch_number;

-- Índices de user
DROP INDEX IF EXISTS idx_user_email;

-- Índices únicos de QR
DROP INDEX IF EXISTS idx_medical_supply_qr_code;
DROP INDEX IF EXISTS idx_batch_qr_code;

-- =======================
-- ELIMINAR TABLAS
-- =======================

-- Tablas de resumen de inventario
DROP TABLE IF EXISTS store_inventory_summary CASCADE;
DROP TABLE IF EXISTS pavilion_inventory_summary CASCADE;

-- Tabla de transferencias
DROP TABLE IF EXISTS supply_transfer CASCADE;

-- Tablas de solicitudes y asignaciones QR
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

-- Tablas de configuración médica
DROP TABLE IF EXISTS doctor_info CASCADE;
DROP TABLE IF EXISTS surgery_typical_supply CASCADE;
DROP TABLE IF EXISTS medical_specialty CASCADE;

-- Tabla de configuración de proveedores
DROP TABLE IF EXISTS supplier_config CASCADE;

-- Tablas de estructura organizacional
DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS store CASCADE;
DROP TABLE IF EXISTS pavilion CASCADE;
DROP TABLE IF EXISTS surgery CASCADE;
DROP TABLE IF EXISTS medical_center CASCADE;

-- =======================
-- ELIMINAR SECUENCIAS (si existen independientemente)
-- =======================
DROP SEQUENCE IF EXISTS surgery_id_seq CASCADE;
DROP SEQUENCE IF EXISTS medical_center_id_seq CASCADE;
DROP SEQUENCE IF EXISTS pavilion_id_seq CASCADE;
DROP SEQUENCE IF EXISTS store_id_seq CASCADE;
DROP SEQUENCE IF EXISTS batch_id_seq CASCADE;
DROP SEQUENCE IF EXISTS medical_supply_id_seq CASCADE;
DROP SEQUENCE IF EXISTS supply_history_id_seq CASCADE;
DROP SEQUENCE IF EXISTS batch_history_id_seq CASCADE;
DROP SEQUENCE IF EXISTS supply_request_id_seq CASCADE;
DROP SEQUENCE IF EXISTS supply_request_item_id_seq CASCADE;
DROP SEQUENCE IF EXISTS supply_request_qr_assignment_id_seq CASCADE;
DROP SEQUENCE IF EXISTS qr_scan_event_id_seq CASCADE;
DROP SEQUENCE IF EXISTS supply_transfer_id_seq CASCADE;
DROP SEQUENCE IF EXISTS pavilion_inventory_summary_id_seq CASCADE;
DROP SEQUENCE IF EXISTS store_inventory_summary_id_seq CASCADE;

-- =======================
-- MENSAJE DE CONFIRMACIÓN
-- =======================
DO $$
BEGIN
    RAISE NOTICE 'Base de datos MediTrack eliminada completamente.';
    RAISE NOTICE 'Todas las tablas, vistas, funciones, triggers, índices y secuencias han sido eliminadas.';
END $$;