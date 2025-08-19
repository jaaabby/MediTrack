-- Migración de rollback: Eliminar campos de código QR
-- Fecha: 2025-08-19
-- Descripción: Eliminar campos qr_code de batch y medical_supply

-- Eliminar índices únicos
DROP INDEX IF EXISTS idx_batch_qr_code;
DROP INDEX IF EXISTS idx_medical_supply_qr_code;

-- Eliminar columnas qr_code
ALTER TABLE batch DROP COLUMN IF EXISTS qr_code;
ALTER TABLE medical_supply DROP COLUMN IF EXISTS qr_code;