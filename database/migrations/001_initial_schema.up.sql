-- Migración inicial: Crear esquema de base de datos para MediTrack
-- Fecha: 2025-08-16
-- Descripción: Tablas principales según modelos Go

CREATE TABLE medical_center (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL
);

CREATE TABLE pavilion (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    medical_center_id INTEGER NOT NULL REFERENCES medical_center(id)
);

CREATE TABLE store (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    medical_center_id INTEGER NOT NULL REFERENCES medical_center(id)
);

CREATE TABLE batch (
    id SERIAL PRIMARY KEY,
    expiration_date DATE NOT NULL,
    amount INTEGER NOT NULL,
    supplier VARCHAR(255) NOT NULL,
    store_id INTEGER NOT NULL REFERENCES store(id)
);

CREATE TABLE medical_supply (
    id SERIAL PRIMARY KEY,
    code INTEGER NOT NULL,
    code_supplier INTEGER,
    name VARCHAR(255) NOT NULL,
    batch_id INTEGER NOT NULL REFERENCES batch(id)
);

CREATE TABLE "user" (
    rut VARCHAR(20) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    medical_center_id INTEGER NOT NULL REFERENCES medical_center(id)
);

CREATE TABLE supply_history (
    id SERIAL PRIMARY KEY,
    date_time TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    destination_type VARCHAR(50) NOT NULL,
    destination_id INTEGER NOT NULL,
    medical_supply_id INTEGER NOT NULL REFERENCES medical_supply(id),
    user_rut VARCHAR(20) NOT NULL REFERENCES "user"(rut)
);