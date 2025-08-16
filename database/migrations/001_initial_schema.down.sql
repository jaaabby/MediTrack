-- Migración de rollback: Eliminar esquema de base de datos para MediTrack
-- Fecha: 2024-01-01
-- Descripción: Eliminar todas las tablas creadas en la migración inicial

-- Eliminar triggers primero
DROP TRIGGER IF EXISTS update_supply_movements_updated_at ON supply_movements;
DROP TRIGGER IF EXISTS update_supply_route_medical_supply_updated_at ON supply_route_medical_supply;
DROP TRIGGER IF EXISTS update_doctor_supply_route_updated_at ON doctor_supply_route;
DROP TRIGGER IF EXISTS update_supply_routes_updated_at ON supply_routes;
DROP TRIGGER IF EXISTS update_medical_supplies_updated_at ON medical_supplies;
DROP TRIGGER IF EXISTS update_locations_updated_at ON locations;
DROP TRIGGER IF EXISTS update_doctors_updated_at ON doctors;
DROP TRIGGER IF EXISTS update_operating_rooms_updated_at ON operating_rooms;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Eliminar función de trigger
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Eliminar tablas en orden de dependencias (hijos primero)
DROP TABLE IF EXISTS supply_movements;
DROP TABLE IF EXISTS supply_route_medical_supply;
DROP TABLE IF EXISTS doctor_supply_route;
DROP TABLE IF EXISTS supply_routes;
DROP TABLE IF EXISTS medical_supplies;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS doctors;
DROP TABLE IF EXISTS operating_rooms;
DROP TABLE IF EXISTS users;

-- Eliminar extensión de UUIDs
DROP EXTENSION IF EXISTS "uuid-ossp"; 