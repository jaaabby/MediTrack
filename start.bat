@echo off
REM Script de inicio rápido para MediTrack con Docker (Windows)

echo 🚀 Iniciando MediTrack con Docker...
echo.

REM Verificar si existe .env
if not exist .env (
    echo ⚠️  Archivo .env no encontrado. Creando desde .env.example...
    if exist .env.example (
        copy .env.example .env
        echo ✅ Archivo .env creado. Por favor, edita las variables según tu entorno.
        echo.
        set /p CONTINUE="¿Deseas continuar con la configuración por defecto? (s/n): "
        if /i not "%CONTINUE%"=="s" (
            echo ❌ Abortado. Por favor, edita .env y vuelve a ejecutar este script.
            exit /b 1
        )
    ) else (
        echo ❌ Error: No se encontró .env.example
        exit /b 1
    )
)

REM Verificar que Docker esté disponible
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Error: Docker no está disponible en el PATH.
    echo.
    echo Por favor:
    echo 1. Cierra y vuelve a abrir esta ventana de CMD
    echo 2. O reinicia tu sesión de Windows
    echo 3. Asegúrate de que Docker Desktop esté corriendo
    echo.
    pause
    exit /b 1
)

echo 📦 Construyendo imágenes Docker...
docker compose build

echo.
echo 🔧 Iniciando servicios...
docker compose up -d

echo.
echo ⏳ Esperando a que los servicios estén listos...
timeout /t 10 /nobreak > nul

echo.
echo 📊 Verificando estado de los servicios...
docker compose ps

echo.
echo ✅ MediTrack está corriendo!
echo.
echo 🌐 Frontend: http://localhost:3000
echo 🔌 Backend API: http://localhost:8080
echo ❤️  Health Check: http://localhost:8080/health
echo.
echo 📝 Para ver los logs:
echo    docker compose logs -f
echo.
echo 🛑 Para detener los servicios:
echo    docker compose down
echo.

pause

