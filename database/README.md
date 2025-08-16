# Base de Datos MediTrack

Este directorio contiene todos los archivos relacionados con la base de datos del sistema MediTrack.

## Estructura de Archivos

```
database/
├── migrations/
│   ├── 001_initial_schema.up.sql      # Esquema inicial de la base de datos
│   ├── 001_initial_schema.down.sql    # Revertir esquema inicial
│   ├── 002_seed_data.sql              # Datos de prueba
│   └── 002_seed_data.down.sql        # Revertir datos de prueba
├── run_migrations.ps1                 # Script completo de migraciones
├── run_seed_data.ps1                  # Script simple para datos de prueba
└── README.md                          # Este archivo
```

## Requisitos Previos

1. **PostgreSQL**: Debe estar instalado y ejecutándose
2. **psql**: Cliente de línea de comandos de PostgreSQL debe estar en el PATH
3. **Docker**: Si usas Docker Compose para ejecutar la base de datos

## Configuración de la Base de Datos

### Usando Docker Compose (Recomendado)

```bash
# Desde el directorio raíz del proyecto
docker-compose up -d postgres
```

### Configuración Manual

Si prefieres usar PostgreSQL local:

```bash
# Crear base de datos
createdb meditrack

# Crear usuario
createuser meditrack_user

# Establecer contraseña
psql -d meditrack -c "ALTER USER meditrack_user PASSWORD 'meditrack_password';"

# Otorgar permisos
psql -d meditrack -c "GRANT ALL PRIVILEGES ON DATABASE meditrack TO meditrack_user;"
```

## Uso de los Scripts

### 1. Script Simple para Datos de Prueba

```powershell
# Navegar al directorio database
cd database

# Ejecutar datos de prueba
.\run_seed_data.ps1
```

### 2. Script Completo de Migraciones

```powershell
# Navegar al directorio database
cd database

# Ejecutar migración hacia arriba (crear esquema)
.\run_migrations.ps1 -Action up

# Insertar datos de prueba
.\run_migrations.ps1 -Action seed

# Revertir datos de prueba
.\run_migrations.ps1 -Action down

# Resetear completamente la base de datos
.\run_migrations.ps1 -Action reset
```

### 3. Ejecución Manual con psql

```bash
# Crear esquema inicial
psql -h localhost -p 5432 -U meditrack_user -d meditrack -f migrations/001_initial_schema.up.sql

# Insertar datos de prueba
psql -h localhost -p 5432 -U meditrack_user -d meditrack -f migrations/002_seed_data.sql

# Revertir datos de prueba
psql -h localhost -p 5432 -U meditrack_user -d meditrack -f migrations/002_seed_data.down.sql
```

## Datos de Prueba Incluidos

El script `002_seed_data.sql` inserta los siguientes datos de ejemplo:

### Usuarios (8)
- **Doctores**: Dr. Juan Pérez (Cirugía General), Dr. María García (Cardiología), Dr. Carlos López (Ortopedia), Dr. Laura Fernández (Neurología), Dr. Roberto Jiménez (Traumatología)
- **Personal**: Enfermera Ana Martínez, Técnico Roberto Silva
- **Pacientes**: 3 pacientes de prueba

### Salas de Operación (4)
- Sala de Operaciones 1 (Cirugías generales)
- Sala de Operaciones 2 (Cirugías cardíacas)
- Sala de Operaciones 3 (Cirugías ortopédicas)
- Sala de Emergencias

### Ubicaciones (7)
- Almacén Central
- Salas de operación
- UCI
- Emergencias
- En Tránsito

### Insumos Médicos (10)
- Jeringas, catéteres, guantes, mascarillas, suturas, bisturíes, gasas, tubos endotraqueales, cánulas nasales, vendajes

### Rutas de Suministro (3)
- 3 rutas de suministro para diferentes pacientes con sus respectivos doctores e insumos

### Movimientos de Insumos (10)
- Diferentes estados y ubicaciones de los insumos médicos

## Solución de Problemas

### Error: "psql no está disponible"
- Asegúrate de que PostgreSQL esté instalado
- Verifica que `psql` esté en el PATH del sistema
- En Windows, puede ser necesario reiniciar PowerShell después de instalar PostgreSQL

### Error: "No se puede conectar a PostgreSQL"
- Verifica que PostgreSQL esté ejecutándose
- Confirma las credenciales en el script
- Verifica que el puerto 5432 esté disponible
- Si usas Docker, asegúrate de que el contenedor esté ejecutándose

### Error: "Permiso denegado"
- Verifica que el usuario tenga permisos en la base de datos
- Confirma que la base de datos exista
- Verifica que las credenciales sean correctas

## Notas Importantes

1. **UUIDs**: Todos los IDs son UUIDs predefinidos para facilitar las pruebas
2. **Relaciones**: Los datos respetan las restricciones de clave foránea
3. **Timestamps**: Los movimientos tienen timestamps realistas (hace X horas)
4. **Reversibilidad**: Todos los datos pueden ser revertidos usando el script `.down.sql`

## Desarrollo

Para agregar más datos de prueba:

1. Edita `002_seed_data.sql`
2. Agrega los nuevos INSERT statements
3. Actualiza `002_seed_data.down.sql` para incluir los DELETE correspondientes
4. Mantén la consistencia de las relaciones entre tablas

## Seguridad

⚠️ **Advertencia**: Estos scripts están diseñados para desarrollo y pruebas. No los uses en producción.

- Las contraseñas están hardcodeadas en los scripts
- Los datos son de ejemplo y no representan información real
- No hay validación de entrada de datos

