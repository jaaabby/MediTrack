-- Migración de rollback: Eliminar esquema de base de datos para MediTrack
-- Fecha: 2025-08-16
-- Descripción: Eliminar todas las tablas creadas en la migración inicial

-- Eliminar tablas en orden de dependencias (hijos primero)
DROP TABLE IF EXISTS supply_history;
DROP TABLE IF EXISTS medical_supply;
DROP TABLE IF EXISTS batch;
DROP TABLE IF EXISTS store;
DROP TABLE IF EXISTS pavilion;
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS medical_center;