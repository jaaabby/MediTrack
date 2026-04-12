# Meditrack Backend

Este es el backend de la aplicación Meditrack, organizado siguiendo una arquitectura limpia y modular por entidad.

## Estructura del Proyecto

### Controllers (Controladores)
Los controladores manejan las peticiones HTTP y la lógica de presentación:

- `user_controller.go` - Controlador para usuarios
- `medical_supply_controller.go` - Controlador para insumos médicos  
- `supply_movement_controller.go` - Controlador para movimientos de insumos
- `operating_room_controller.go` - Controlador para salas de operación
- `doctor_controller.go` - Controlador para doctores
- `location_controller.go` - Controlador para ubicaciones
- `supply_route_controller.go` - Controlador para rutas de suministro
- `doctor_supply_route_controller.go` - Controlador para relación doctor-ruta
- `supply_route_medical_supply_controller.go` - Controlador para relación ruta-insumo
- `common.go` - Estructuras comunes de respuesta

### Services (Servicios)
Los servicios contienen la lógica de negocio:

- `user_service.go` - Servicios para usuarios
- `medical_supply_service.go` - Servicios para insumos médicos
- `supply_movement_service.go` - Servicios para movimientos de insumos
- `operating_room_service.go` - Servicios para salas de operación
- `doctor_service.go` - Servicios para doctores
- `location_service.go` - Servicios para ubicaciones
- `supply_route_service.go` - Servicios para rutas de suministro
- `doctor_supply_route_service.go` - Servicios para relación doctor-ruta
- `supply_route_medical_supply_service.go` - Servicios para relación ruta-insumo

### Repository (Repositorios)
Los repositorios manejan el acceso a datos:

- `interfaces.go` - Interfaces de todos los repositorios
- `user_repository.go` - Repositorio para usuarios
- `medical_supply_repository.go` - Repositorio para insumos médicos
- `supply_movement_repository.go` - Repositorio para movimientos de insumos
- `operating_room_repository.go` - Repositorio para salas de operación
- `doctor_repository.go` - Repositorio para doctores
- `location_repository.go` - Repositorio para ubicaciones
- `supply_route_repository.go` - Repositorio para rutas de suministro
- `doctor_supply_route_repository.go` - Repositorio para relación doctor-ruta
- `supply_route_medical_supply_repository.go` - Repositorio para relación ruta-insumo

### Models (Modelos)
Los modelos definen las estructuras de datos:

- `userModel.go` - Modelo de usuario
- `medicalSupplyModel.go` - Modelo de insumo médico
- `supplyMovementModel.go` - Modelo de movimiento de insumo
- `operatingRoomModel.go` - Modelo de sala de operación
- `doctorModel.go` - Modelo de doctor
- `locationModel.go` - Modelo de ubicación
- `supplyRouteModel.go` - Modelo de ruta de suministro
- `doctorSupplyRouteModel.go` - Modelo de relación doctor-ruta
- `supplyRouteMedicalSupplyModel.go` - Modelo de relación ruta-insumo

### Routes (Rutas)
Las rutas organizan los endpoints de la API:

- `routes.go` - Archivo principal de configuración de rutas
- `user_routes.go` - Rutas para usuarios
- `medical_supply_routes.go` - Rutas para insumos médicos
- `supply_movement_routes.go` - Rutas para movimientos de insumos
- `operating_room_routes.go` - Rutas para salas de operación
- `doctor_routes.go` - Rutas para doctores
- `location_routes.go` - Rutas para ubicaciones
- `supply_route_routes.go` - Rutas para rutas de suministro
- `doctor_supply_route_routes.go` - Rutas para relación doctor-ruta
- `supply_route_medical_supply_routes.go` - Rutas para relación ruta-insumo
- `additional_routes.go` - Rutas adicionales (trazabilidad, estadísticas, etc.)

## Arquitectura

La aplicación sigue el patrón de arquitectura en capas:

1. **Controllers** - Manejan las peticiones HTTP
2. **Services** - Contienen la lógica de negocio
3. **Repository** - Manejan el acceso a datos
4. **Models** - Definen las estructuras de datos

## Endpoints Disponibles

### Usuarios
- `POST /api/v1/users/` - Crear usuario
- `GET /api/v1/users/` - Obtener todos los usuarios
- `GET /api/v1/users/:id` - Obtener usuario por ID
- `PUT /api/v1/users/:id` - Actualizar usuario
- `DELETE /api/v1/users/:id` - Eliminar usuario

### Insumos Médicos
- `POST /api/v1/medical-supplies/` - Crear insumo médico
- `GET /api/v1/medical-supplies/` - Obtener todos los insumos
- `GET /api/v1/medical-supplies/:id` - Obtener insumo por ID
- `GET /api/v1/medical-supplies/qr` - Obtener insumo por código QR
- `PUT /api/v1/medical-supplies/:id` - Actualizar insumo
- `DELETE /api/v1/medical-supplies/:id` - Eliminar insumo

### Movimientos de Insumos
- `POST /api/v1/supply-movements/` - Crear movimiento
- `GET /api/v1/supply-movements/` - Obtener todos los movimientos
- `GET /api/v1/supply-movements/:id` - Obtener movimiento por ID
- `GET /api/v1/supply-movements/status` - Obtener movimientos por estado
- `PUT /api/v1/supply-movements/:id` - Actualizar movimiento
- `DELETE /api/v1/supply-movements/:id` - Eliminar movimiento

### Salas de Operación
- `POST /api/v1/operating-rooms/` - Crear sala de operación
- `GET /api/v1/operating-rooms/` - Obtener todas las salas
- `GET /api/v1/operating-rooms/:id` - Obtener sala por ID
- `PUT /api/v1/operating-rooms/:id` - Actualizar sala
- `DELETE /api/v1/operating-rooms/:id` - Eliminar sala

### Doctores
- `POST /api/v1/doctors/` - Crear doctor
- `GET /api/v1/doctors/` - Obtener todos los doctores
- `GET /api/v1/doctors/:id` - Obtener doctor por ID
- `GET /api/v1/doctors/specialty` - Obtener doctores por especialidad
- `PUT /api/v1/doctors/:id` - Actualizar doctor
- `DELETE /api/v1/doctors/:id` - Eliminar doctor

### Ubicaciones
- `POST /api/v1/locations/` - Crear ubicación
- `GET /api/v1/locations/` - Obtener todas las ubicaciones
- `GET /api/v1/locations/:id` - Obtener ubicación por ID
- `PUT /api/v1/locations/:id` - Actualizar ubicación
- `DELETE /api/v1/locations/:id` - Eliminar ubicación

### Rutas de Suministro
- `POST /api/v1/supply-routes/` - Crear ruta de suministro
- `GET /api/v1/supply-routes/` - Obtener todas las rutas
- `GET /api/v1/supply-routes/:id` - Obtener ruta por ID
- `GET /api/v1/supply-routes/patient` - Obtener rutas por paciente
- `GET /api/v1/supply-routes/operating-room` - Obtener rutas por sala
- `PUT /api/v1/supply-routes/:id` - Actualizar ruta
- `DELETE /api/v1/supply-routes/:id` - Eliminar ruta

### Relaciones Doctor-Ruta
- `POST /api/v1/doctor-supply-routes/` - Crear relación
- `GET /api/v1/doctor-supply-routes/` - Obtener todas las relaciones
- `GET /api/v1/doctor-supply-routes/:id` - Obtener relación por ID
- `GET /api/v1/doctor-supply-routes/doctor` - Obtener por doctor
- `GET /api/v1/doctor-supply-routes/supply-route` - Obtener por ruta
- `PUT /api/v1/doctor-supply-routes/:id` - Actualizar relación
- `DELETE /api/v1/doctor-supply-routes/:id` - Eliminar relación

### Relaciones Ruta-Insumo
- `POST /api/v1/supply-route-medical-supplies/` - Crear relación
- `GET /api/v1/supply-route-medical-supplies/` - Obtener todas las relaciones
- `GET /api/v1/supply-route-medical-supplies/:id` - Obtener relación por ID
- `GET /api/v1/supply-route-medical-supplies/supply-route` - Obtener por ruta
- `GET /api/v1/supply-route-medical-supplies/medical-supply` - Obtener por insumo
- `PUT /api/v1/supply-route-medical-supplies/:id` - Actualizar relación
- `DELETE /api/v1/supply-route-medical-supplies/:id` - Eliminar relación

## Cómo Ejecutar

1. Asegúrate de tener Go instalado
2. Instala las dependencias: `go mod tidy`
3. Configura la base de datos en `config/config.go`
4. Ejecuta: `go run main.go`

## Beneficios de esta Estructura

- **Mantenibilidad**: Cada entidad tiene sus propios archivos
- **Escalabilidad**: Fácil agregar nuevas entidades
- **Testabilidad**: Cada capa puede ser testeada independientemente
- **Separación de responsabilidades**: Clara división entre capas
- **Reutilización**: Los servicios pueden ser reutilizados por diferentes controladores
- **Organización por dominio**: Cada entidad tiene su propia carpeta lógica
- **Facilidad de navegación**: Fácil encontrar código relacionado con una entidad específica
