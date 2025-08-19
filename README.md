# MediTrack - Sistema de Gestión de Inventario Médico

MediTrack es una aplicación web completa para la gestión y control de inventario de insumos médicos en centros de salud.

## 🏗️ Arquitectura

- **Backend**: API REST en Go con Gin y GORM
- **Frontend**: Aplicación Vue.js 3 con Tailwind CSS
- **Base de datos**: PostgreSQL
- **Contenedores**: Docker y Docker Compose

## 🚀 Instalación y Ejecución

### Prerrequisitos

- Docker y Docker Compose
- Go 1.21+ (para desarrollo del backend)
- Node.js 18+ (para desarrollo del frontend)

### 1. Clonar el repositorio

```bash
git clone <url-del-repositorio>
cd MediTrack
```

### 2. Configurar la base de datos

```bash
# Crear y poblar la base de datos
cd database
docker-compose up -d postgres

# Ejecutar las migraciones
psql -h localhost -U postgres -d meditrack -f migrations/001_initial_schema.up.sql

# Insertar datos de prueba (opcional)
psql -h localhost -U postgres -d meditrack -f script.sql
```

### 3. Configurar variables de entorno

Crear archivo `.env` en el directorio raíz:

```env
# Backend
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=meditrack
DB_SSL_MODE=disable

# Frontend
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### 4. Ejecutar con Docker Compose

```bash
# Ejecutar toda la aplicación
docker-compose up -d

# Ver logs
docker-compose logs -f
```

### 5. Desarrollo local

#### Backend
```bash
cd backend
go mod download
go run main.go
```

#### Frontend
```bash
cd frontend
npm install
npm run dev
```

## 📊 Estructura de la Base de Datos

### Tablas principales:

- **medical_center**: Centros médicos
- **store**: Bodegas de almacenamiento
- **batch**: Lotes de insumos
- **medical_supply**: Insumos médicos
- **user**: Usuarios del sistema
- **supply_history**: Historial de movimientos

### Consulta de inventario:

La vista de inventario obtiene datos de múltiples tablas mediante JOINs:

```sql
SELECT 
    medical_supply.id,
    medical_supply.code,
    medical_supply.name,
    batch.expiration_date,
    batch.amount,
    batch.supplier,
    store.name as store_name,
    medical_center.name as medical_center_name
FROM medical_supply
LEFT JOIN batch ON medical_supply.batch_id = batch.id
LEFT JOIN store ON batch.store_id = store.id
LEFT JOIN medical_center ON store.medical_center_id = medical_center.id
```

## 🔧 API Endpoints

### Inventario
- `GET /api/v1/medical-supplies/inventory` - Obtener inventario completo
- `GET /api/v1/medical-supplies` - Obtener todos los insumos
- `POST /api/v1/medical-supplies` - Crear insumo
- `PUT /api/v1/medical-supplies/:id` - Actualizar insumo
- `DELETE /api/v1/medical-supplies/:id` - Eliminar insumo

## 🎯 Funcionalidades del Frontend

- **Vista de inventario** con filtros y búsqueda
- **Paginación** de resultados
- **Ordenamiento** por múltiples criterios
- **Filtros** por nombre, lote, centro médico
- **Indicadores visuales** para fechas de vencimiento
- **Acciones** para ver, editar y eliminar insumos

## 🐛 Solución de Problemas

### Error de conexión a la base de datos
- Verificar que PostgreSQL esté ejecutándose
- Comprobar credenciales en el archivo `.env`
- Verificar que la base de datos `meditrack` exista

### Error en el frontend
- Verificar que el backend esté ejecutándose en el puerto 8080
- Comprobar la variable `VITE_API_BASE_URL`
- Revisar la consola del navegador para errores

## 📝 Contribución

1. Fork el proyecto
2. Crear una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles. 