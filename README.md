# MediTrack - Sistema de Trazabilidad Médica

Sistema completo de trazabilidad de insumos médicos con arquitectura limpia, backend en Go y frontend en Vue 3.

## 🏗️ Arquitectura del Proyecto

```
MediTrack/
├── backend/          # API REST en Go con Clean Architecture
├── frontend/         # Aplicación Vue 3 + TailwindCSS
├── database/         # Migraciones y scripts de base de datos
├── docker-compose.yml # Orquestación de servicios
└── README.md         # Este archivo
```

## 🚀 Tecnologías Utilizadas

### Backend
- **Go 1.21+** con Clean Architecture
- **Gin** como framework web
- **PostgreSQL** como base de datos
- **golang-migrate** para migraciones
- **JWT** para autenticación (preparado para implementación)

### Frontend
- **Vue 3** con Composition API
- **Vite** como bundler
- **TailwindCSS** para estilos
- **Vue Router** para navegación
- **Pinia** para manejo de estado

## 📋 Prerrequisitos

- Docker y Docker Compose
- Go 1.21 o superior
- Node.js 18 o superior
- npm o yarn

## 🛠️ Instalación y Configuración

### 1. Clonar y configurar el proyecto

```bash
# Clonar el repositorio
git clone <tu-repositorio>
cd MediTrack

# Copiar archivos de configuración
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

### 2. Configurar variables de entorno

Editar `backend/.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=meditrack_user
DB_PASSWORD=meditrack_password
DB_NAME=meditrack_db
DB_SSL_MODE=disable
JWT_SECRET=tu_jwt_secret_aqui
PORT=8080
```

Editar `frontend/.env`:
```env
VITE_API_BASE_URL=http://localhost:8080/api
```

### 3. Levantar servicios con Docker

```bash
# Levantar todos los servicios
docker-compose up -d

# Ver logs
docker-compose logs -f
```

### 4. Ejecutar migraciones de base de datos

```bash
# Desde el directorio backend
cd backend
go run cmd/migrate/main.go
```

### 5. Instalar dependencias del frontend

```bash
cd frontend
npm install
```

### 6. Ejecutar el proyecto

```bash
# Terminal 1: Backend
cd backend
go run cmd/server/main.go

# Terminal 2: Frontend
cd frontend
npm run dev
```

## 🌐 Acceso a la Aplicación

- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080
- **PostgreSQL:** localhost:5432

## 📊 Estructura de la Base de Datos

El sistema incluye las siguientes entidades principales:

- **supply_routes**: Rutas de suministro médico
- **operating_rooms**: Salas de operación
- **doctors**: Médicos y especialistas
- **medical_supplies**: Insumos médicos
- **supply_movements**: Movimientos de insumos
- **locations**: Ubicaciones físicas
- **users**: Usuarios del sistema

## 🔧 Desarrollo

### Backend

```bash
cd backend

# Ejecutar tests
go test ./...

# Ejecutar con hot reload (requiere air)
air

# Generar documentación
swag init -g cmd/server/main.go
```

### Frontend

```bash
cd frontend

# Ejecutar tests
npm run test

# Build para producción
npm run build

# Preview build
npm run preview
```

## 📁 Estructura de Carpetas

### Backend (Clean Architecture)
```
backend/
├── cmd/           # Puntos de entrada de la aplicación
├── internal/      # Lógica de negocio interna
│   ├── domain/    # Entidades y reglas de negocio
│   ├── usecase/   # Casos de uso
│   ├── repository/# Interfaces de repositorio
│   └── delivery/  # Controladores HTTP
├── pkg/           # Paquetes reutilizables
├── config/        # Configuración
└── migrations/    # Migraciones de base de datos
```

### Frontend
```
frontend/
├── src/
│   ├── components/    # Componentes Vue
│   ├── views/         # Páginas/Vistas
│   ├── stores/        # Stores de Pinia
│   ├── router/        # Configuración de rutas
│   └── assets/        # Recursos estáticos
├── public/            # Archivos públicos
└── dist/              # Build de producción
```

## 🚧 Próximas Funcionalidades

- [ ] Autenticación JWT completa
- [ ] Sistema de trazabilidad con QR codes
- [ ] Integración con Smart Contracts
- [ ] Alertas automatizadas
- [ ] Dashboard de estadísticas avanzado
- [ ] Sistema de auditoría
- [ ] API para dispositivos móviles

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 🆘 Soporte

Para soporte técnico o preguntas, crear un issue en el repositorio o contactar al equipo de desarrollo.

---

**MediTrack** - Trazabilidad médica del futuro 🏥✨ 