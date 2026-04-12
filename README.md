# Meditrack - Sistema de Gestión de Inventario Médico

## Descripción
Meditrack es un sistema completo para la gestión y control de inventario médico, desarrollado con Go en el backend y Vue.js en el frontend.

## Funcionalidades Principales

### Gestión de Inventario
- Visualización completa del inventario de insumos médicos
- Búsqueda y filtrado por múltiples criterios
- Ordenamiento por diferentes campos
- Paginación de resultados
- Edición de lotes existentes

### Historial de Movimientos
- **Historial por Lote**: Botón de historial en cada fila del inventario que muestra el historial completo de un lote específico
- **Historial Global**: Botón "Historial Global" en el header que permite buscar y visualizar el historial de todos los lotes
- **Buscador Inteligente**: Búsqueda por número de lote, nombre del insumo, fecha, tipo de movimiento, cantidad o usuario
- **Filtros Avanzados**: Ordenamiento por fecha, tipo, cantidad y usuario
- **Paginación**: Navegación eficiente a través de grandes volúmenes de datos históricos

## Estructura del Proyecto

```
Meditrack/
├── backend/                 # Servidor Go
│   ├── controllers/         # Controladores de la API
│   ├── models/             # Modelos de datos
│   ├── services/           # Lógica de negocio
│   ├── routes/             # Definición de rutas
│   └── config/             # Configuración de base de datos
├── frontend/               # Aplicación Vue.js
│   ├── src/
│   │   ├── views/          # Vistas principales
│   │   ├── services/       # Servicios de API
│   │   └── router/         # Configuración de rutas
└── database/               # Migraciones y scripts SQL
```

## Instalación y Configuración

### Backend (Go)
1. Navegar al directorio `backend/`
2. Instalar dependencias: `go mod download`
3. Configurar variables de entorno (base de datos, JWT, etc.)
4. Ejecutar: `go run main.go`

### Frontend (Vue.js)
1. Navegar al directorio `frontend/`
2. Instalar dependencias: `npm install`
3. Configurar la URL del backend en las variables de entorno
4. Ejecutar: `npm run dev`

### Base de Datos
1. Configurar PostgreSQL
2. Ejecutar las migraciones en `database/migrations/`
3. Verificar la conexión desde el backend

## API Endpoints

### Historial de Lotes
- `GET /api/v1/batch-histories/batch/:batchId` - Obtener historial de un lote específico
- `GET /api/v1/batch-histories/search?q=:term` - Buscar historial por término
- `GET /api/v1/batch-histories` - Obtener historial global

### Inventario
- `GET /api/v1/medical-supplies/list` - Listar todo el inventario
- `PUT /api/v1/batches/:id` - Actualizar un lote

## Características del Historial

### Información Mostrada
- **Fecha y Hora**: Timestamp exacto de cada movimiento
- **Tipo de Movimiento**: Estado actual, modificaciones, etc.
- **Cantidad**: Cambios en el stock del lote
- **Usuario**: RUT del usuario que realizó la acción
- **Detalles**: Información adicional del movimiento

### Funcionalidades de Búsqueda
- Búsqueda por texto libre en todos los campos
- Filtrado en tiempo real
- Ordenamiento por múltiples criterios
- Paginación automática

## Tecnologías Utilizadas

- **Backend**: Go, Gin, GORM, PostgreSQL
- **Frontend**: Vue.js 3, Tailwind CSS, Axios
- **Base de Datos**: PostgreSQL con migraciones
- **Autenticación**: JWT
- **Notificaciones**: Sistema de alertas en tiempo real

## Contribución

1. Fork del repositorio
2. Crear rama para nueva funcionalidad
3. Implementar cambios
4. Ejecutar pruebas
5. Crear Pull Request

## Licencia

Este proyecto está bajo la Licencia MIT. 