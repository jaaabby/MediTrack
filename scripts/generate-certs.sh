#!/bin/bash

# Script para generar certificados SSL para desarrollo con Docker
# Genera certificados para backend y frontend

set -e

echo "🔐 Generando certificados SSL para MediTrack..."

# Colores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Directorios
BACKEND_CERTS_DIR="backend/certs"
FRONTEND_CERTS_DIR="frontend/certs"

# Crear directorios si no existen
mkdir -p "$BACKEND_CERTS_DIR"
mkdir -p "$FRONTEND_CERTS_DIR"

# Obtener IP local (para acceso desde celular en la misma red)
get_local_ip() {
    # Windows
    if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
        ipconfig | grep -i "IPv4" | head -1 | awk '{print $NF}' | tr -d '\r'
    # Linux
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        hostname -I | awk '{print $1}'
    # macOS
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}' | head -1
    else
        echo "127.0.0.1"
    fi
}

LOCAL_IP=$(get_local_ip)
echo -e "${YELLOW}📱 IP local detectada: $LOCAL_IP${NC}"

# Generar certificados para backend
echo "🔧 Generando certificados para backend..."
cd "$BACKEND_CERTS_DIR"

# Usar OpenSSL si está disponible, sino usar Go
if command -v openssl &> /dev/null; then
    openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -days 365 -nodes \
        -subj "/C=CL/ST=Region/L=City/O=MediTrack Dev/CN=localhost" \
        -addext "subjectAltName=IP:127.0.0.1,IP:${LOCAL_IP},IP:172.17.0.1,DNS:localhost,DNS:backend,DNS:*.localhost"
    
    echo -e "${GREEN}✅ Certificados backend generados con OpenSSL${NC}"
else
    # Usar Go como fallback
    if command -v go &> /dev/null; then
        echo "Usando Go para generar certificados..."
        go run generate_cert.go
        echo -e "${GREEN}✅ Certificados backend generados con Go${NC}"
    else
        echo "❌ Error: Se requiere OpenSSL o Go para generar certificados"
        exit 1
    fi
fi

cd - > /dev/null

# Generar certificados para frontend
echo "🔧 Generando certificados para frontend..."
cd "$FRONTEND_CERTS_DIR"

if command -v openssl &> /dev/null; then
    openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -days 365 -nodes \
        -subj "/C=CL/ST=Region/L=City/O=MediTrack Dev/CN=localhost" \
        -addext "subjectAltName=IP:127.0.0.1,IP:${LOCAL_IP},IP:172.17.0.1,DNS:localhost,DNS:frontend,DNS:*.localhost"
    
    echo -e "${GREEN}✅ Certificados frontend generados con OpenSSL${NC}"
else
    # Copiar certificados del backend si no hay OpenSSL
    cp "../$BACKEND_CERTS_DIR/server.crt" server.crt
    cp "../$BACKEND_CERTS_DIR/server.key" server.key
    echo -e "${YELLOW}⚠️  Certificados frontend copiados del backend (OpenSSL no disponible)${NC}"
fi

cd - > /dev/null

echo ""
echo -e "${GREEN}✅ Certificados SSL generados exitosamente!${NC}"
echo ""
echo "📋 Ubicación de certificados:"
echo "   Backend:  $BACKEND_CERTS_DIR/"
echo "   Frontend: $FRONTEND_CERTS_DIR/"
echo ""
echo "📱 Para acceder desde tu celular:"
echo "   Frontend HTTPS: https://$LOCAL_IP:3443"
echo "   Backend HTTPS:  https://$LOCAL_IP:8443"
echo ""
echo "⚠️  Nota: Debes aceptar el certificado autofirmado en tu navegador/celular"

