# Certificados SSL para Desarrollo

Esta carpeta contiene los certificados SSL autofirmados necesarios para ejecutar el servidor en modo HTTPS durante el desarrollo.

## Generar Certificados

Si necesitas regenerar los certificados (por ejemplo, si cambia tu IP local), ejecuta:

```bash
cd backend/certs
go run generate_cert.go
```

Los certificados generados son:
- `server.crt` - Certificado SSL
- `server.key` - Clave privada

## Nota

Estos certificados son **solo para desarrollo**. Los archivos `.crt` y `.key` están excluidos del control de versiones (`.gitignore`).

En producción, debes usar certificados válidos emitidos por una autoridad certificadora.

