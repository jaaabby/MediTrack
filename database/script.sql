-- Migración: Agregar campos de código QR a las tablas
-- Fecha: 2025-08-19
-- Descripción: Agregar campos qr_code únicos a batch y medical_supply para trazabilidad

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

INSERT INTO batch (id, expiration_date, amount, supplier, store_id, qr_code) VALUES
(1, '2026-12-31', 10, 'Proveedor Uno', 1, 'BATCH_1_1'),
(2, '2025-08-16', 5, 'Proveedor Dos', 2, 'BATCH_2_1');

-- Corregir la secuencia de la tabla batch para evitar errores de clave duplicada
-- (Solo para PostgreSQL)
SELECT setval(pg_get_serial_sequence('batch', 'id'), (SELECT MAX(id) FROM batch));

INSERT INTO supply_code (code, name, code_supplier) VALUES
(1001, 'Guantes', 5001),
(1002, 'Mascarillas', 5002);

INSERT INTO medical_supply (code, batch_id, qr_code) VALUES
(1001, 1, 'SUPPLY_1_1'),
(1001, 1, 'SUPPLY_2_1'),
(1001, 1, 'SUPPLY_3_1'),
(1001, 1, 'SUPPLY_4_1'),
(1001, 1, 'SUPPLY_5_1'),
(1001, 1, 'SUPPLY_6_1'),
(1001, 1, 'SUPPLY_7_1'),
(1001, 1, 'SUPPLY_8_1'),
(1001, 1, 'SUPPLY_9_1'),
(1001, 1, 'SUPPLY_10_1'),
(1002, 2, 'SUPPLY_11_1'),
(1002, 2, 'SUPPLY_12_1'),
(1002, 2, 'SUPPLY_13_1'),
(1002, 2, 'SUPPLY_14_1'),
(1002, 2, 'SUPPLY_15_1');

-- Corregir la secuencia de la tabla medical_supply para evitar errores de clave duplicada
-- (Solo para PostgreSQL)
SELECT setval(pg_get_serial_sequence('medical_supply', 'id'), (SELECT MAX(id) FROM medical_supply));

INSERT INTO "user" (rut, name, email, password, role, medical_center_id) VALUES
('12345678-9', 'Juan Pérez', 'juan@meditrack.com', 'password123', 'admin', 1),
('98765432-1', 'Ana Gómez', 'ana@meditrack.com', 'password456', 'doctor', 2);

INSERT INTO supply_history (date_time, status, destination_type, destination_id, medical_supply_id, user_rut) VALUES
('2025-08-16 10:00:00', 'entregado', 'pavilion', 1, 1, '12345678-9'),
('2025-08-16 11:00:00', 'recibido', 'store', 2, 2, '98765432-1');

-- Poblar tabla batch_history
INSERT INTO batch_history (date_time, change_details, previous_values, new_values, user_name, batch_id, user_rut, batch_number) VALUES
('2025-08-16 10:00:00', 'Lote creado', NULL, '{"expiration_date": "2026-12-31", "amount": 10, "supplier": "Proveedor Uno", "store_id": 1}', 'Juan Pérez', 1, '12345678-9', 1);