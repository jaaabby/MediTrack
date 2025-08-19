-- Migración: Agregar campos de código QR a las tablas
-- Fecha: 2025-08-19
-- Descripción: Agregar campos qr_code únicos a batch y medical_supply para trazabilidad

-- Agregar campo qr_code a la tabla batch
ALTER TABLE batch ADD COLUMN qr_code VARCHAR(255);

-- Agregar campo qr_code a la tabla medical_supply
ALTER TABLE medical_supply ADD COLUMN qr_code VARCHAR(255);

-- Generar códigos QR únicos para datos existentes ANTES de crear índices
UPDATE batch SET qr_code = 'BATCH_' || EXTRACT(EPOCH FROM NOW())::bigint || '_' || id WHERE qr_code IS NULL;
UPDATE medical_supply SET qr_code = 'SUPPLY_' || EXTRACT(EPOCH FROM NOW())::bigint || '_' || id WHERE qr_code IS NULL;

-- Hacer los campos NOT NULL después de establecer valores
ALTER TABLE batch ALTER COLUMN qr_code SET NOT NULL;
ALTER TABLE medical_supply ALTER COLUMN qr_code SET NOT NULL;

-- Crear índices únicos DESPUÉS de que todos los valores estén establecidos
CREATE UNIQUE INDEX idx_batch_qr_code ON batch(qr_code);
CREATE UNIQUE INDEX idx_medical_supply_qr_code ON medical_supply(qr_code);