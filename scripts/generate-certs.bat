@echo off
setlocal enabledelayedexpansion
REM Script para generar certificados SSL para desarrollo con Docker (Windows)

echo 🔐 Generando certificados SSL para MediTrack...
echo.

REM Directorios
set BACKEND_CERTS_DIR=backend\certs
set FRONTEND_CERTS_DIR=frontend\certs

REM Crear directorios si no existen
if not exist "%BACKEND_CERTS_DIR%" mkdir "%BACKEND_CERTS_DIR%"
if not exist "%FRONTEND_CERTS_DIR%" mkdir "%FRONTEND_CERTS_DIR%"

REM Obtener IP local
set LOCAL_IP=127.0.0.1
for /f "tokens=2 delims=:" %%a in ('ipconfig ^| findstr /i "IPv4"') do (
    set "LOCAL_IP=%%a"
    set "LOCAL_IP=!LOCAL_IP: =!"
    if not "!LOCAL_IP!"=="" (
        goto :found_ip
    )
)
:found_ip

echo 📱 IP local detectada: %LOCAL_IP%
echo.

REM Verificar si OpenSSL está disponible
where openssl >nul 2>&1
if %errorlevel% equ 0 (
    echo 🔧 Generando certificados para backend...
    cd %BACKEND_CERTS_DIR%
    
    openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -days 365 -nodes ^
        -subj "/C=CL/ST=Region/L=City/O=MediTrack Dev/CN=localhost" ^
        -addext "subjectAltName=IP:127.0.0.1,IP:%LOCAL_IP%,IP:172.17.0.1,DNS:localhost,DNS:backend,DNS:*.localhost"
    
    cd ..\..
    
    echo 🔧 Generando certificados para frontend...
    cd %FRONTEND_CERTS_DIR%
    
    openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -days 365 -nodes ^
        -subj "/C=CL/ST=Region/L=City/O=MediTrack Dev/CN=localhost" ^
        -addext "subjectAltName=IP:127.0.0.1,IP:%LOCAL_IP%,IP:172.17.0.1,DNS:localhost,DNS:frontend,DNS:*.localhost"
    
    cd ..\..
    
    echo.
    echo ✅ Certificados SSL generados exitosamente con OpenSSL!
) else (
    echo ⚠️  OpenSSL no encontrado. Intentando usar Go...
    
    REM Verificar si Go está disponible
    where go >nul 2>&1
    if %errorlevel% equ 0 (
        echo 🔧 Generando certificados para backend...
        cd %BACKEND_CERTS_DIR%
        go run generate_cert.go
        cd ..\..
        
        echo 🔧 Copiando certificados para frontend...
        cd %FRONTEND_CERTS_DIR%
        copy ..\%BACKEND_CERTS_DIR%\server.crt server.crt
        copy ..\%BACKEND_CERTS_DIR%\server.key server.key
        cd ..\..
        
        echo.
        echo ✅ Certificados SSL generados exitosamente con Go!
    ) else (
        echo.
        echo ⚠️  OpenSSL y Go no encontrados. Intentando usar Docker...
        REM Intentar usar Docker para generar certificados con Go
        docker --version >nul 2>&1
        if !errorlevel! equ 0 (
            echo 🔧 Generando certificados usando Docker con Go...
            if exist "%BACKEND_CERTS_DIR%\generate_cert.go" (
                REM Construir ruta absoluta para Docker
                set "ABS_CERT_PATH=%CD%\%BACKEND_CERTS_DIR%"
                
                echo    Ejecutando Docker...
                docker run --rm -v "!ABS_CERT_PATH!:/certs" -w /certs golang:1.23-alpine go run generate_cert.go
                if !errorlevel! equ 0 (
                    if exist "%BACKEND_CERTS_DIR%\server.crt" (
                        echo ✅ Certificados backend generados con Docker/Go
                        
                        REM Copiar certificados al frontend
                        if not exist "%FRONTEND_CERTS_DIR%" mkdir "%FRONTEND_CERTS_DIR%"
                        copy "%BACKEND_CERTS_DIR%\server.crt" "%FRONTEND_CERTS_DIR%\server.crt" >nul 2>&1
                        copy "%BACKEND_CERTS_DIR%\server.key" "%FRONTEND_CERTS_DIR%\server.key" >nul 2>&1
                        echo ✅ Certificados frontend copiados desde backend
                        goto :certs_generated
                    ) else (
                        echo ❌ Error: Los certificados no se generaron correctamente
                        goto :error_no_tools
                    )
                ) else (
                    echo ❌ Error al generar certificados con Docker
                    echo    Verifica que Docker Desktop esté corriendo
                    goto :error_no_tools
                )
            ) else (
                echo ❌ No se encontró generate_cert.go en %BACKEND_CERTS_DIR%
                goto :error_no_tools
            )
        ) else (
            :error_no_tools
            echo.
            echo ❌ Error: Se requiere OpenSSL, Go o Docker para generar certificados
            echo.
            echo Opciones:
            echo   1. Instala OpenSSL: https://slproweb.com/products/Win32OpenSSL.html
            echo   2. Instala Go: https://golang.org/dl/
            echo   3. Asegúrate de que Docker Desktop esté corriendo
            echo.
            echo O genera los certificados manualmente usando el script de Go:
            echo   cd backend\certs
            echo   go run generate_cert.go
            echo   copy server.crt ..\..\frontend\certs\
            echo   copy server.key ..\..\frontend\certs\
            pause
            exit /b 1
        )
    )
)

:certs_generated

echo.
echo 📋 Ubicación de certificados:
echo    Backend:  %BACKEND_CERTS_DIR%\
echo    Frontend: %FRONTEND_CERTS_DIR%\
echo.
echo 📱 Para acceder desde tu celular:
echo    Frontend HTTPS: https://%LOCAL_IP%:3443
echo    Backend HTTPS:  https://%LOCAL_IP%:8443
echo.
echo ⚠️  Nota: Debes aceptar el certificado autofirmado en tu navegador/celular
echo.
pause

