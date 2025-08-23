-- Migración: Agregar campos de código QR a las tablas
-- Fecha: 2025-08-19
-- Descripción: Agregar campos qr_code únicos a batch y medical_supply para trazabilidad

-- Agregar campo qr_code a la tabla batch
ALTER TABLE batch ADD COLUMN qr_code VARCHAR(255);

-- Agregar campo qr_code a la tabla medical_supply
ALTER TABLE medical_supply ADD COLUMN qr_code VARCHAR(255);

-- Poblar tablas base primero
INSERT INTO medical_center (id, name, address, phone, email) VALUES
(1, 'Centro Médico Central', 'Av. Principal 123', '123456789', 'central@meditrack.com'),
(2, 'Clínica Norte', 'Calle Norte 456', '987654321', 'norte@meditrack.com');

INSERT INTO pavilion (id, name, medical_center_id) VALUES
(1, 'Pabellón A', 1),
(2, 'Pabellón B', 2);

INSERT INTO store (id, name, type, medical_center_id) VALUES
(1, 'Bodega Principal', 'central', 1),
(2, 'Bodega Secundaria', 'secundaria', 2);

-- Poblar tabla batch (sin qr_code primero)
INSERT INTO batch (id, expiration_date, amount, supplier, store_id) VALUES
(1, '2026-12-31', 10, 'Proveedor Uno', 1),
(2, '2025-08-16', 5, 'Proveedor Dos', 2);

-- Poblar tabla supply_code
INSERT INTO supply_code (code, name, code_supplier) VALUES
(1001, 'Guantes', 5001),
(1002, 'Mascarillas', 5002);

-- Poblar tabla medical_supply
INSERT INTO medical_supply (id, code, batch_id) VALUES
-- Guantes del código 1001 (múltiples guantes individuales)
(1, 1001, 1),
(2, 1001, 1),
(3, 1001, 1),
(4, 1001, 1),
(5, 1001, 1),
(6, 1001, 1),
(7, 1001, 1),
(8, 1001, 1),
(9, 1001, 1),
(10, 1001, 1),
-- Mascarillas del código 1002
(11, 1002, 2),
(12, 1002, 2),
(13, 1002, 2),
(14, 1002, 2),
(15, 1002, 2);

INSERT INTO "user" (rut, name, email, password, role, medical_center_id) VALUES
('12345678-9', 'Juan Pérez', 'juan@meditrack.com', 'password123', 'admin', 1),
('98765432-1', 'Ana Gómez', 'ana@meditrack.com', 'password456', 'doctor', 2);

INSERT INTO supply_history (id, date_time, status, destination_type, destination_id, medical_supply_id, user_rut) VALUES
(1, '2025-08-16 10:00:00', 'entregado', 'pavilion', 1, 1, '12345678-9'),
(2, '2025-08-16 11:00:00', 'recibido', 'store', 2, 2, '98765432-1');

<<<<<<< HEAD
-- Generar códigos QR únicos para datos existentes
UPDATE batch SET qr_code = 'BATCH_' || EXTRACT(EPOCH FROM NOW())::bigint || '_' || id WHERE qr_code IS NULL;
UPDATE medical_supply SET qr_code = 'SUPPLY_' || EXTRACT(EPOCH FROM NOW())::bigint || '_' || id WHERE qr_code IS NULL;

-- Hacer los campos NOT NULL después de establecer valores
ALTER TABLE batch ALTER COLUMN qr_code SET NOT NULL;
ALTER TABLE medical_supply ALTER COLUMN qr_code SET NOT NULL;

-- Crear índices únicos DESPUÉS de que todos los valores estén establecidos
CREATE UNIQUE INDEX idx_batch_qr_code ON batch(qr_code);
CREATE UNIQUE INDEX idx_medical_supply_qr_code ON medical_supply(qr_code);
=======
-- Poblar tabla batch_history
INSERT INTO batch_history (date_time, change_details, previous_values, new_values, user_name, batch_id, user_rut, batch_number) VALUES
('2025-08-16 10:00:00', 'Lote creado', '{"expiration_date": "2026-12-31", "amount": 10, "supplier": "Proveedor Uno", "store_id": 1}', '{"expiration_date": "2026-12-31", "amount": 10, "supplier": "Proveedor Uno", "store_id": 1}', 'Juan Pérez', 1, '12345678-9', 1),
('2025-08-16 11:00:00', 'Lote actualizado', '{"expiration_date": "2025-08-16", "amount": 5, "supplier": "Proveedor Dos", "store_id": 2}', '{"expiration_date": "2025-08-16", "amount": 5, "supplier": "Proveedor Dos", "store_id": 2}', 'Ana Gómez', 2, '98765432-1', 2);
>>>>>>> origin/main
