@echo off
REM Script de inicio rápido para MediTrack con Docker (Windows)

echo 🚀 Iniciando MediTrack con Docker...
echo.

REM Verificar si existe .env
if not exist .env (
    echo ⚠️  Archivo .env no encontrado.
    if exist .env.example (
        echo    Creando desde .env.example...
        copy .env.example .env >nul
        echo ✅ Archivo .env creado. Por favor, edita las variables según tu entorno.
        echo.
        set /p CONTINUE="¿Deseas continuar con la configuración por defecto? (s/n): "
        if /i not "%CONTINUE%"=="s" (
            echo ❌ Abortado. Por favor, edita .env y vuelve a ejecutar este script.
            exit /b 1
        )
    ) else (
        echo ⚠️  No se encontró .env.example
        echo    Continuando sin archivo .env (usando valores por defecto de docker-compose.yml)
        echo    Para personalizar, crea un archivo .env con tus variables de entorno.
        echo.
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

REM Verificar certificados SSL (ya generados)
echo.
echo Certificados SSL verificados (backend\certs\server.crt)

echo.
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
echo 🌐 Frontend HTTP:  http://localhost:3000
echo 🔒 Frontend HTTPS: https://localhost:3443  (para acceso desde celular con cámara)
echo 🔌 Backend API HTTP:  http://localhost:8080
echo 🔒 Backend API HTTPS:  https://localhost:8443
echo ❤️  Health Check: http://localhost:8080/health
echo.
echo 📱 Para acceder desde tu celular:
echo    1. Asegúrate de que tu celular esté en la misma red WiFi
echo    2. Obtén tu IP local (ver mensaje de generación de certificados)
echo    3. Accede desde el celular: https://TU_IP:3443
echo    4. Acepta el certificado autofirmado cuando se solicite
echo.
echo 📝 Para ver los logs:
echo    docker compose logs -f
echo.
echo 🛑 Para detener los servicios:
echo    docker compose down
echo.

pause

