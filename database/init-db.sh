#!/bin/bash
set -e

echo "🚀 Iniciando inicialización de la base de datos..."

# Esperar a que PostgreSQL esté listo
until pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"; do
  echo "⏳ Esperando a que PostgreSQL esté listo..."
  sleep 2
done

echo "✅ PostgreSQL está listo"

# Ejecutar migraciones
if [ -d "/docker-entrypoint-initdb.d/migrations" ]; then
    echo "📦 Ejecutando migraciones..."
    for migration in /docker-entrypoint-initdb.d/migrations/*.up.sql; do
        if [ -f "$migration" ]; then
            echo "  ➡️  Ejecutando: $(basename $migration)"
            psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$migration"
        fi
    done
    echo "✅ Migraciones completadas"
fi

# Ejecutar script de poblado
if [ -f "/docker-entrypoint-initdb.d/01-script.sql" ]; then
    echo "📊 Ejecutando script de poblado..."
    psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/01-script.sql
    echo "✅ Script de poblado completado"
fi

echo "🎉 Inicialización de la base de datos completada"

