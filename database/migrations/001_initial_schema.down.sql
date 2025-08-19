-- Migración de rollback: Eliminar esquema completo de MediTrack
-- Fecha: 2025-08-16
-- Descripción: Eliminar todas las tablas y estructuras

-- Eliminar índices únicos de QR primero
DROP INDEX IF EXISTS idx_batch_qr_code;
DROP INDEX IF EXISTS idx_medical_supply_qr_code;

-- Eliminar columnas qr_code si existen
ALTER TABLE IF EXISTS batch DROP COLUMN IF EXISTS qr_code;
ALTER TABLE IF EXISTS medical_supply DROP COLUMN IF EXISTS qr_code;

-- Eliminar tablas en orden de dependencias (hijos primero)
DROP TABLE IF EXISTS supply_history;
DROP TABLE IF EXISTS medical_supply;
DROP TABLE IF EXISTS supply_code;
DROP TABLE IF EXISTS batch;
DROP TABLE IF EXISTS store;
DROP TABLE IF EXISTS pavilion;
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS medical_center;