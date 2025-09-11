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
INSERT INTO medical_supply (code, batch_id, qr_code, status) VALUES
-- Guantes del lote 1
(1001, 1, 'SUPPLY_1_1', 'disponible'),
(1001, 1, 'SUPPLY_2_1', 'disponible'),
(1001, 1, 'SUPPLY_3_1', 'disponible'),
(1001, 1, 'SUPPLY_4_1', 'en_camino_a_pabellon'),
(1001, 1, 'SUPPLY_5_1', 'recepcionado'),
(1001, 1, 'SUPPLY_6_1', 'consumido'),
(1001, 1, 'SUPPLY_7_1', 'disponible'),
(1001, 1, 'SUPPLY_8_1', 'disponible'),
(1001, 1, 'SUPPLY_9_1', 'disponible'),
(1001, 1, 'SUPPLY_10_1', 'disponible'),
-- Mascarillas del lote 2
(1002, 2, 'SUPPLY_11_1', 'disponible'),
(1002, 2, 'SUPPLY_12_1', 'disponible'),
(1002, 2, 'SUPPLY_13_1', 'en_camino_a_pabellon'),
(1002, 2, 'SUPPLY_14_1', 'recepcionado'),
(1002, 2, 'SUPPLY_15_1', 'consumido'),
-- Jeringas del lote 3
(1003, 3, 'SUPPLY_16_1', 'disponible'),
(1003, 3, 'SUPPLY_17_1', 'disponible'),
(1003, 3, 'SUPPLY_18_1', 'disponible'),
(1003, 3, 'SUPPLY_19_1', 'disponible'),
(1003, 3, 'SUPPLY_20_1', 'disponible'),
-- Agujas del lote 4
(1004, 4, 'SUPPLY_21_1', 'disponible'),
(1004, 4, 'SUPPLY_22_1', 'disponible'),
(1004, 4, 'SUPPLY_23_1', 'disponible'),
(1004, 4, 'SUPPLY_24_1', 'disponible'),
(1004, 4, 'SUPPLY_25_1', 'disponible')
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

-- Usuario enfermera de ejemplo
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
    '22222222-2',
    'María González',
    'enfermera@meditrack.com',
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', -- admin123 hasheada con bcrypt
    'enfermera',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

-- Usuario doctor de ejemplo
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
    '33333333-3',
    'Dr. Carlos Pérez',
    'doctor@meditrack.com',
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', -- admin123 hasheada con bcrypt
    'doctor',
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
('2025-08-16 10:00:00', 'recepcionado', 'pavilion', 1, 1, '12345678-9'),
('2025-08-16 11:00:00', 'disponible', 'store', 2, 2, '87654321-0'),
('2025-08-16 12:00:00', 'recepcionado', 'pavilion', 2, 3, '11111111-1'),
('2025-08-16 13:00:00', 'consumido', 'pavilion', 1, 4, '87654321-0'),
('2025-08-16 14:00:00', 'disponible', 'store', 1, 5, '12345678-9'),
-- Historial para insumos con estados específicos
('2025-08-16 15:00:00', 'en_camino_a_pabellon', 'pavilion', 1, 4, '11111111-1'),
('2025-08-16 16:00:00', 'recepcionado', 'pavilion', 1, 5, '12345678-9'),
('2025-08-16 17:00:00', 'consumido', 'pavilion', 1, 6, '87654321-0'),
('2025-08-16 18:00:00', 'en_camino_a_pabellon', 'pavilion', 2, 13, '11111111-1'),
('2025-08-16 19:00:00', 'recepcionado', 'pavilion', 2, 14, '12345678-9'),
('2025-08-16 20:00:00', 'consumido', 'pavilion', 2, 15, '87654321-0')
ON CONFLICT DO NOTHING;

-- Poblar historial de lotes
INSERT INTO batch_history (date_time, change_details, previous_values, new_values, user_name, batch_id, user_rut, batch_number) VALUES
('2025-08-16 10:00:00', 'Lote creado', NULL, '{"expiration_date": "2026-12-31", "amount": 10, "supplier": "Proveedor Uno", "store_id": 1}', 'Administrador del Sistema', 1, '12345678-9', 1),
('2025-08-16 11:00:00', 'Lote creado', NULL, '{"expiration_date": "2025-12-31", "amount": 5, "supplier": "Proveedor Dos", "store_id": 2}', 'Usuario Pabellón', 2, '87654321-0', 2),
('2025-08-16 12:00:00', 'Lote creado', NULL, '{"expiration_date": "2026-06-30", "amount": 15, "supplier": "Proveedor Tres", "store_id": 1}', 'Encargado Bodega', 3, '11111111-1', 3),
('2025-08-16 13:00:00', 'Lote creado', NULL, '{"expiration_date": "2025-10-15", "amount": 8, "supplier": "Proveedor Cuatro", "store_id": 2}', 'Administrador del Sistema', 4, '12345678-9', 4),
('2025-08-16 14:00:00', 'Cantidad actualizada', '{"amount": 10}', '{"amount": 8}', 'Encargado Bodega', 1, '11111111-1', 1)
ON CONFLICT DO NOTHING;
-- Insertar solicitud de ejemplo
INSERT INTO supply_request (
    request_number, pavilion_id, requested_by, requested_by_name,
    request_date, status, priority, notes, medical_center_id
) VALUES (
    'SOL-20250120140000', 1, '12345678-9', 'Juan Pérez',
    NOW() - INTERVAL '1 hour', 'pending', 'normal',
    'Solicitud de prueba para implementación de trazabilidad QR', 1
);
-- Insertar items de ejemplo

INSERT INTO supply_request_item (
    supply_request_id, supply_code, supply_name, quantity_requested,
    specifications, is_pediatric, size, urgency_level
) VALUES 
(1, 1001, 'Guantes', 50, 'Talla M, látex libre', FALSE, 'M', 'normal'),
(1, 1002, 'Mascarillas', 100, 'N95, uso pediátrico', TRUE, 'Pediatric', 'high');

-- Script para arreglar el problema de secuencia en la tabla batch
-- Este script resetea la secuencia de auto-incremento al valor correcto

-- Verificar el ID máximo actual en la tabla batch
SELECT 'ID máximo actual en batch:' as info, COALESCE(MAX(id), 0) as max_id FROM batch;

-- Resetear la secuencia al valor correcto
-- La secuencia debe ser mayor que el ID máximo existente
SELECT setval('batch_id_seq', COALESCE((SELECT MAX(id) FROM batch), 0) + 1, false);

-- Verificar el estado de la secuencia después del reset
SELECT 'Secuencia después del reset:' as info, last_value, is_called FROM batch_id_seq;

-- También verificar que no hay problemas con QR codes duplicados
SELECT 'QR codes duplicados:' as info, qr_code, COUNT(*) as count 
FROM batch 
WHERE qr_code IS NOT NULL 
GROUP BY qr_code 
HAVING COUNT(*) > 1;