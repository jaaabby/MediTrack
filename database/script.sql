-- Poblar tabla medical_center
INSERT INTO medical_center (id, name, address, phone, email) VALUES
(1, 'Centro Médico Central', 'Av. Principal 123', '123456789', 'central@meditrack.com'),
(2, 'Clínica Norte', 'Calle Norte 456', '987654321', 'norte@meditrack.com');

-- Poblar tabla pavilion
INSERT INTO pavilion (id, name, medical_center_id) VALUES
(1, 'Pabellón A', 1),
(2, 'Pabellón B', 2);

-- Poblar tabla store
INSERT INTO store (id, name, type, medical_center_id) VALUES
(1, 'Bodega Principal', 'central', 1),
(2, 'Bodega Secundaria', 'secundaria', 2);

-- Poblar tabla batch
INSERT INTO batch (id, expiration_date, amount, supplier, store_id) VALUES
(1, '2026-12-31', 10, 'Proveedor Uno', 1),
(2, '2025-08-16', 5, 'Proveedor Dos', 2);

-- Poblar tabla supply_code
INSERT INTO supply_code (code, name, code_supplier, batch_id) VALUES
(1001, 'Guantes', 5001, 1),
(1002, 'Mascarillas', 5002, 2);

-- Poblar tabla medical_supply
INSERT INTO medical_supply (id, code) VALUES
-- Guantes del código 1001 (múltiples guantes individuales)
(1, 1001),
(2, 1001),
(3, 1001),
(4, 1001),
(5, 1001),
(6, 1001),
(7, 1001),
(8, 1001),
(9, 1001),
(10, 1001),
-- Mascarillas del código 1002
(11, 1002),
(12, 1002),
(13, 1002),
(14, 1002),
(15, 1002);

-- Poblar tabla "user"
INSERT INTO "user" (rut, name, email, password, role, medical_center_id) VALUES
('12345678-9', 'Juan Pérez', 'juan@meditrack.com', 'password123', 'admin', 1),
('98765432-1', 'Ana Gómez', 'ana@meditrack.com', 'password456', 'doctor', 2);

-- Poblar tabla supply_history
INSERT INTO supply_history (id, date_time, status, destination_type, destination_id, medical_supply_id, user_rut) VALUES
(1, '2025-08-16 10:00:00', 'entregado', 'pavilion', 1, 1, '12345678-9'),
(2, '2025-08-16 11:00:00', 'recibido', 'store', 2, 2, '98765432-1');
