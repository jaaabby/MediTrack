#!/bin/bash

# Script de inicio rápido para Meditrack con Docker

echo "🚀 Iniciando Meditrack con Docker..."
echo ""

# Verificar si existe .env
if [ ! -f .env ]; then
    echo "⚠️  Archivo .env no encontrado. Creando desde .env.example..."
    if [ -f .env.example ]; then
        cp .env.example .env
        echo "✅ Archivo .env creado. Por favor, edita las variables según tu entorno."
        echo ""
        read -p "¿Deseas continuar con la configuración por defecto? (s/n): " -n 1 -r
        echo ""
        if [[ ! $REPLY =~ ^[Ss]$ ]]; then
            echo "❌ Abortado. Por favor, edita .env y vuelve a ejecutar este script."
            exit 1
        fi
    else
        echo "❌ Error: No se encontró .env.example"
        exit 1
    fi
fi

# Verificar que Docker esté disponible
if ! command -v docker &> /dev/null; then
    echo "❌ Error: Docker no está disponible en el PATH."
    echo ""
    echo "Por favor, asegúrate de que Docker esté instalado y corriendo."
    exit 1
fi

# Generar certificados SSL si no existen
echo ""
echo "🔐 Verificando certificados SSL..."
if [ ! -f "backend/certs/server.crt" ]; then
    echo "⚠️  Certificados SSL no encontrados. Generando..."
    if [ -f "scripts/generate-certs.sh" ]; then
        chmod +x scripts/generate-certs.sh
        ./scripts/generate-certs.sh
    else
        echo "⚠️  Script de generación de certificados no encontrado."
        echo "   Continuando sin HTTPS (solo HTTP disponible)."
    fi
else
    echo "✅ Certificados SSL encontrados."
fi

echo ""
echo "📦 Construyendo imágenes Docker..."
docker compose build

echo ""
echo "🔧 Iniciando servicios..."
docker compose up -d

echo ""
echo "⏳ Esperando a que los servicios estén listos..."
sleep 10

echo ""
echo "📊 Verificando estado de los servicios..."
docker compose ps

echo ""
echo "✅ Meditrack está corriendo!"
echo ""
echo "🌐 Frontend HTTP:  http://localhost:3000"
echo "🔒 Frontend HTTPS: https://localhost:3443  (para acceso desde celular con cámara)"
echo "🔌 Backend API HTTP:  http://localhost:8080"
echo "🔒 Backend API HTTPS:  https://localhost:8443"
echo "❤️  Health Check: http://localhost:8080/health"
echo ""
echo "📱 Para acceder desde tu celular:"
echo "   1. Asegúrate de que tu celular esté en la misma red WiFi"
echo "   2. Obtén tu IP local (ver mensaje de generación de certificados)"
echo "   3. Accede desde el celular: https://TU_IP:3443"
echo "   4. Acepta el certificado autofirmado cuando se solicite"
echo ""
echo "📝 Para ver los logs:"
echo "   docker compose logs -f"
echo ""
echo "🛑 Para detener los servicios:"
echo "   docker compose down"
echo ""

