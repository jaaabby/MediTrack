#!/bin/bash

# Script de inicio rápido para MediTrack con Docker

echo "🚀 Iniciando MediTrack con Docker..."
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
echo "✅ MediTrack está corriendo!"
echo ""
echo "🌐 Frontend: http://localhost:3000"
echo "🔌 Backend API: http://localhost:8080"
echo "❤️  Health Check: http://localhost:8080/health"
echo ""
echo "📝 Para ver los logs:"
echo "   docker compose logs -f"
echo ""
echo "🛑 Para detener los servicios:"
echo "   docker compose down"
echo ""

