# Etapa 1: Build
FROM golang:1.23-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git ca-certificates tzdata

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o meditrack ./main.go

# Etapa 2: Runtime
FROM alpine:latest

# Instalar certificados CA, timezone data y wget para health checks
RUN apk --no-cache add ca-certificates tzdata wget

# Establecer zona horaria
ENV TZ=America/Santiago

WORKDIR /root/

# Copiar el binario compilado desde la etapa de build
COPY --from=builder /app/meditrack .

# Copiar certificados y plantillas de email si existen
COPY --from=builder /app/certs ./certs
COPY --from=builder /app/mailer/templates ./mailer/templates

# Exponer puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./meditrack"]

