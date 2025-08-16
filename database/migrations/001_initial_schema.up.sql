-- Migración inicial: Crear esquema de base de datos para MediTrack
-- Fecha: 2024-01-01
-- Descripción: Crear todas las tablas principales del sistema

-- Habilitar extensión para UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tabla de usuarios
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de salas de operación
CREATE TABLE operating_rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de doctores
CREATE TABLE doctors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    specialty VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de ubicaciones
CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de insumos médicos
CREATE TABLE medical_supplies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(500) NOT NULL,
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de rutas de suministro
CREATE TABLE supply_routes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    patient_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    operating_room_id UUID NOT NULL REFERENCES operating_rooms(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de relación doctor-ruta de suministro
CREATE TABLE doctor_supply_route (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    doctor_id UUID NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    supply_route_id UUID NOT NULL REFERENCES supply_routes(id) ON DELETE CASCADE,
    role VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de relación ruta de suministro-insumo médico
CREATE TABLE supply_route_medical_supply (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    supply_route_id UUID NOT NULL REFERENCES supply_routes(id) ON DELETE CASCADE,
    medical_supply_id UUID NOT NULL REFERENCES medical_supplies(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de movimientos de insumos
CREATE TABLE supply_movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    medical_supply_id UUID NOT NULL REFERENCES medical_supplies(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    status VARCHAR(100) NOT NULL,
    scanned_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scanned_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Crear índices para mejorar rendimiento
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_medical_supplies_qr_code ON medical_supplies(qr_code);
CREATE INDEX idx_supply_routes_patient_id ON supply_routes(patient_id);
CREATE INDEX idx_supply_routes_operating_room_id ON supply_routes(operating_room_id);
CREATE INDEX idx_doctor_supply_route_doctor_id ON doctor_supply_route(doctor_id);
CREATE INDEX idx_doctor_supply_route_supply_route_id ON doctor_supply_route(supply_route_id);
CREATE INDEX idx_supply_route_medical_supply_supply_route_id ON supply_route_medical_supply(supply_route_id);
CREATE INDEX idx_supply_route_medical_supply_medical_supply_id ON supply_route_medical_supply(medical_supply_id);
CREATE INDEX idx_supply_movements_medical_supply_id ON supply_movements(medical_supply_id);
CREATE INDEX idx_supply_movements_location_id ON supply_movements(location_id);
CREATE INDEX idx_supply_movements_status ON supply_movements(status);
CREATE INDEX idx_supply_movements_scanned_at ON supply_movements(scanned_at);

-- Crear trigger para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Aplicar trigger a todas las tablas
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_operating_rooms_updated_at BEFORE UPDATE ON operating_rooms FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_doctors_updated_at BEFORE UPDATE ON doctors FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_locations_updated_at BEFORE UPDATE ON locations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_medical_supplies_updated_at BEFORE UPDATE ON medical_supplies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_supply_routes_updated_at BEFORE UPDATE ON supply_routes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_doctor_supply_route_updated_at BEFORE UPDATE ON doctor_supply_route FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_supply_route_medical_supply_updated_at BEFORE UPDATE ON supply_route_medical_supply FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_supply_movements_updated_at BEFORE UPDATE ON supply_movements FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); 