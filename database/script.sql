-- Script de Poblado Completo para MediTrack
-- Fecha: 2025-08-19
-- Descripción: Poblado completo de todas las tablas del sistema

-- ============================================
-- POBLADO DE TABLAS BASE
-- ============================================

-- Poblar centros médicos
INSERT INTO medical_center (id, name, address, phone, email) VALUES
(1, 'Centro Médico Principal', 'Av. Principal 123', '+56 2 2345 6789', 'info@centromedico.cl'),
(2, 'Clínica Norte', 'Calle Norte 456', '987654321', 'norte@meditrack.com')
ON CONFLICT (id) DO NOTHING;

-- Poblar pabellones
INSERT INTO pavilion (id, name, medical_center_id) VALUES
(1, 'Pabellón A', 1),
(2, 'Pabellón B', 2)
ON CONFLICT (id) DO NOTHING;

-- Poblar bodegas
INSERT INTO store (id, name, type, medical_center_id) VALUES
(1, 'Bodega Principal', 'central', 1),
(2, 'Bodega Secundaria', 'secundaria', 2)
ON CONFLICT (id) DO NOTHING;

-- Poblar códigos de insumos
INSERT INTO supply_code (code, name, code_supplier) VALUES
(1001, 'Guantes', 5001),
(1002, 'Mascarillas', 5002),
(1003, 'Jeringas', 5003),
(1004, 'Agujas', 5004),
(1005, 'Gasas', 5005)
ON CONFLICT (code) DO NOTHING;

-- Poblar lotes
INSERT INTO batch (id, expiration_date, amount, supplier, store_id, qr_code) VALUES
(1, '2026-12-31', 10, 'Proveedor Uno', 1, 'BATCH_1_1'),
(2, '2025-12-31', 5, 'Proveedor Dos', 2, 'BATCH_2_1'),
(3, '2026-06-30', 15, 'Proveedor Tres', 1, 'BATCH_3_1'),
(4, '2025-10-15', 8, 'Proveedor Cuatro', 2, 'BATCH_4_1')
ON CONFLICT (id) DO NOTHING;

-- Poblar insumos médicos
INSERT INTO medical_supply (code, batch_id, qr_code) VALUES
-- Guantes del lote 1
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
-- Mascarillas del lote 2
(1002, 2, 'SUPPLY_11_1'),
(1002, 2, 'SUPPLY_12_1'),
(1002, 2, 'SUPPLY_13_1'),
(1002, 2, 'SUPPLY_14_1'),
(1002, 2, 'SUPPLY_15_1'),
-- Jeringas del lote 3
(1003, 3, 'SUPPLY_16_1'),
(1003, 3, 'SUPPLY_17_1'),
(1003, 3, 'SUPPLY_18_1'),
(1003, 3, 'SUPPLY_19_1'),
(1003, 3, 'SUPPLY_20_1'),
-- Agujas del lote 4
(1004, 4, 'SUPPLY_21_1'),
(1004, 4, 'SUPPLY_22_1'),
(1004, 4, 'SUPPLY_23_1'),
(1004, 4, 'SUPPLY_24_1'),
(1004, 4, 'SUPPLY_25_1')
ON CONFLICT (qr_code) DO NOTHING;

-- ============================================
-- POBLADO DE USUARIOS
-- ============================================

-- Script para insertar usuarios del sistema
-- Contraseña: admin123 (hasheada con bcrypt)
-- Nota: Las contraseñas se almacenan como hash, no como texto plano

-- Usuario administrador del sistema
INSERT INTO "user" (
    rut, 
    name, 
    email, 
    password, 
    role, 
    medical_center_id, 
    is_active, 
    created_at, 
    updated_at
) VALUES (
    '12345678-9',
    'Administrador del Sistema',
    'admin@meditrack.com',
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', -- admin123 hasheada con bcrypt
    'admin',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

-- Usuario de pabellón de ejemplo
INSERT INTO "user" (
    rut, 
    name, 
    email, 
    password, 
    role, 
    medical_center_id, 
    is_active, 
    created_at, 
    updated_at
) VALUES (
    '87654321-0',
    'Usuario Pabellón',
    'pabellon@meditrack.com',
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', -- admin123 hasheada con bcrypt
    'pabellón',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

-- Usuario encargado de bodega de ejemplo
INSERT INTO "user" (
    rut, 
    name, 
    email, 
    password, 
    role, 
    medical_center_id, 
    is_active, 
    created_at, 
    updated_at
) VALUES (
    '11111111-1',
    'Encargado Bodega',
    'bodega@meditrack.com',
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', -- admin123 hasheada con bcrypt
    'encargado de bodega',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

-- ============================================
-- POBLADO DE HISTORIALES
-- ============================================

-- Poblar historial de insumos
INSERT INTO supply_history (date_time, status, destination_type, destination_id, medical_supply_id, user_rut) VALUES
('2025-08-16 10:00:00', 'entregado', 'pavilion', 1, 1, '12345678-9'),
('2025-08-16 11:00:00', 'recibido', 'store', 2, 2, '87654321-0'),
('2025-08-16 12:00:00', 'entregado', 'pavilion', 2, 3, '11111111-1'),
('2025-08-16 13:00:00', 'consumido', 'pavilion', 1, 4, '87654321-0'),
('2025-08-16 14:00:00', 'recibido', 'store', 1, 5, '12345678-9')
ON CONFLICT DO NOTHING;

-- Poblar historial de lotes
INSERT INTO batch_history (date_time, change_details, previous_values, new_values, user_name, batch_id, user_rut, batch_number) VALUES
('2025-08-16 10:00:00', 'Lote creado', NULL, '{"expiration_date": "2026-12-31", "amount": 10, "supplier": "Proveedor Uno", "store_id": 1}', 'Administrador del Sistema', 1, '12345678-9', 1),
('2025-08-16 11:00:00', 'Lote creado', NULL, '{"expiration_date": "2025-12-31", "amount": 5, "supplier": "Proveedor Dos", "store_id": 2}', 'Usuario Pabellón', 2, '87654321-0', 2),
('2025-08-16 12:00:00', 'Lote creado', NULL, '{"expiration_date": "2026-06-30", "amount": 15, "supplier": "Proveedor Tres", "store_id": 1}', 'Encargado Bodega', 3, '11111111-1', 3),
('2025-08-16 13:00:00', 'Lote creado', NULL, '{"expiration_date": "2025-10-15", "amount": 8, "supplier": "Proveedor Cuatro", "store_id": 2}', 'Administrador del Sistema', 4, '12345678-9', 4),
('2025-08-16 14:00:00', 'Cantidad actualizada', '{"amount": 10}', '{"amount": 8}', 'Encargado Bodega', 1, '11111111-1', 1)
ON CONFLICT DO NOTHING;