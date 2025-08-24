-- Migración de rollback: Eliminar esquema completo de MediTrack
-- Fecha: 2025-08-19
-- Descripción: Eliminar todas las tablas, estructuras y campos QR

-- Eliminar índices únicos de QR primero
DROP INDEX IF EXISTS idx_batch_qr_code;
DROP INDEX IF EXISTS idx_medical_supply_qr_code;

-- Eliminar columnas qr_code si existen
ALTER TABLE IF EXISTS batch DROP COLUMN IF EXISTS qr_code;
ALTER TABLE IF EXISTS medical_supply DROP COLUMN IF EXISTS qr_code;

-- Eliminar triggers y funciones primero
DROP TRIGGER IF EXISTS trg_log_batch_delete ON batch;
DROP TRIGGER IF EXISTS set_batch_number ON batch_history;
DROP FUNCTION IF EXISTS log_batch_delete();
DROP FUNCTION IF EXISTS trg_set_batch_number();

-- Eliminar todas las tablas usando CASCADE para manejar dependencias automáticamente
DROP TABLE IF EXISTS medical_center CASCADE;
DROP TABLE IF EXISTS supply_code CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS pavilion CASCADE;
DROP TABLE IF EXISTS store CASCADE;
DROP TABLE IF EXISTS batch CASCADE;
DROP TABLE IF EXISTS medical_supply CASCADE;
DROP TABLE IF EXISTS supply_history CASCADE;
DROP TABLE IF EXISTS batch_history CASCADE;