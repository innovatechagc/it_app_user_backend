# 🐳 Docker Guide - User Service

## Índice
- [🎯 Visión General](#-visión-general)
- [🚀 Inicio Rápido](#-inicio-rápido)
- [🔧 Configuración](#-configuración)
- [📦 Construcción](#-construcción)
- [🏃 Ejecución](#-ejecución)
- [🔄 Docker Compose](#-docker-compose)
- [🚀 Producción](#-producción)
- [🔧 Troubleshooting](#-troubleshooting)

## 🎯 Visión General

El User Service está completamente dockerizado para facilitar el desarrollo, testing y despliegue. Incluye configuración para desarrollo local y producción.

### Componentes Docker
- **Aplicación Go**: Microservicio principal
- **PostgreSQL**: Base de datos
- **Volúmenes**: Persistencia de datos
- **Networks**: Comunicación entre servicios

## 🚀 Inicio Rápido

### Opción 1: Todo en Uno
```bash
# Clonar y ejecutar
git clone <repository-url>
cd it-app_user
docker-compose up --build
```

### Opción 2: Solo Base de Datos
```bash
# Solo PostgreSQL para desarrollo local
docker-compose up postgres -d

# Ejecutar aplicación localmente
go run cmd/main.go
```

### Verificar que Funciona
```bash
# Health check
curl http://localhost:8081/health

# Listar usuarios
curl http://localhost:8081/users
```

## 🔧 Configuración

### Dockerfile
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
```

### Docker Compose
```yaml
services:
  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    container_name: it-app_postgres
    environment:
      POSTGRES_DB: itapp
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Go Microservice
  app:
    build: .
    container_name: it-app_user_service
    environment:
      # Database
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: itapp
      
      # Server
      PORT: 8080
      ENVIRONMENT: development
      
      # Logging
      LOG_LEVEL: info
      
      # Rate Limiting
      RATE_LIMIT_RPS: 100
      RATE_LIMIT_BURST: 200
      
      # Firebase
      FIREBASE_PROJECT_ID: innovatech-app
    ports:
      - "8081:8080"
    volumes:
      - ./firebase-service-account.json:/root/firebase-service-account.json:ro
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

volumes:
  postgres_data:
```

## 📦 Construcción

### Construir Imagen
```bash
# Construir imagen de la aplicación
docker build -t user-service:latest .

# Construir con tag específico
docker build -t user-service:v1.0.0 .

# Construir sin cache
docker build --no-cache -t user-service:latest .
```

### Multi-stage Build
```dockerfile
# Optimizado para producción
FROM golang:1.21-alpine AS builder

# Instalar dependencias de build
RUN apk add --no-cache git ca-certificates tzdata

# Crear usuario no-root
RUN adduser -D -g '' appuser

WORKDIR /build

# Copiar y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copiar código fuente
COPY . .

# Build optimizado para producción
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o app cmd/main.go

# Final stage - imagen mínima
FROM scratch

# Copiar certificados SSL
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiar timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copiar usuario
COPY --from=builder /etc/passwd /etc/passwd

# Copiar binario
COPY --from=builder /build/app /app

# Usar usuario no-root
USER appuser

# Exponer puerto
EXPOSE 8080

# Ejecutar aplicación
ENTRYPOINT ["/app"]
```

### Optimizaciones de Build
```bash
# Usar .dockerignore para excluir archivos innecesarios
echo "
.git
.gitignore
README.md
docs/
*.md
.env
.vscode/
.idea/
" > .dockerignore

# Build con BuildKit (más rápido)
DOCKER_BUILDKIT=1 docker build -t user-service:latest .
```

## 🏃 Ejecución

### Comandos Básicos
```bash
# Ejecutar contenedor
docker run -d \
  --name user-service \
  -p 8081:8080 \
  -e DB_HOST=host.docker.internal \
  -e FIREBASE_PROJECT_ID=innovatech-app \
  user-service:latest

# Ver logs
docker logs user-service

# Logs en tiempo real
docker logs -f user-service

# Ejecutar comando dentro del contenedor
docker exec -it user-service sh

# Parar contenedor
docker stop user-service

# Eliminar contenedor
docker rm user-service
```

### Variables de Entorno
```bash
# Archivo .env para Docker
cat > .env.docker << EOF
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=itapp

# Server
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info

# Rate Limiting
RATE_LIMIT_RPS=100
RATE_LIMIT_BURST=200

# Firebase
FIREBASE_PROJECT_ID=innovatech-app
EOF

# Usar archivo .env
docker run --env-file .env.docker user-service:latest
```

## 🔄 Docker Compose

### Comandos Principales
```bash
# Iniciar todos los servicios
docker-compose up

# Iniciar en background
docker-compose up -d

# Construir y iniciar
docker-compose up --build

# Solo un servicio
docker-compose up postgres

# Ver logs
docker-compose logs

# Logs de un servicio específico
docker-compose logs app

# Parar servicios
docker-compose down

# Parar y eliminar volúmenes
docker-compose down -v

# Reiniciar un servicio
docker-compose restart app
```

### Desarrollo con Hot Reload
```yaml
# docker-compose.dev.yml
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    volumes:
      - .:/app
      - /app/vendor
    environment:
      - GO_ENV=development
    command: air -c .air.toml
```

```dockerfile
# Dockerfile.dev
FROM golang:1.21-alpine

WORKDIR /app

# Instalar air para hot reload
RUN go install github.com/cosmtrek/air@latest

# Copiar archivos de configuración
COPY go.mod go.sum ./
RUN go mod download

# Exponer puerto
EXPOSE 8080

# Comando por defecto
CMD ["air", "-c", ".air.toml"]
```

### Configuración de Air
```toml
# .air.toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main cmd/main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false
```

## 🚀 Producción

### Imagen de Producción
```dockerfile
# Dockerfile.prod
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

# Crear usuario no-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

USER appuser

EXPOSE 8080
CMD ["./main"]
```

### Docker Compose para Producción
```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.prod
    container_name: user-service-prod
    environment:
      DB_HOST: ${DB_HOST}
      DB_PASSWORD: ${DB_PASSWORD}
      FIREBASE_PROJECT_ID: ${FIREBASE_PROJECT_ID}
      ENVIRONMENT: production
      LOG_LEVEL: warn
    ports:
      - "8080:8080"
    volumes:
      - ./firebase-service-account.json:/root/firebase-service-account.json:ro
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  postgres:
    image: postgres:15-alpine
    container_name: postgres-prod
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_prod_data:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  postgres_prod_data:
```

### Despliegue con Docker Swarm
```bash
# Inicializar swarm
docker swarm init

# Crear stack
docker stack deploy -c docker-compose.prod.yml user-service-stack

# Ver servicios
docker service ls

# Escalar servicio
docker service scale user-service-stack_app=3

# Ver logs
docker service logs user-service-stack_app

# Actualizar servicio
docker service update --image user-service:v2.0.0 user-service-stack_app
```

### Registry y CI/CD
```bash
# Tag para registry
docker tag user-service:latest registry.company.com/user-service:v1.0.0

# Push a registry
docker push registry.company.com/user-service:v1.0.0

# Pull desde registry
docker pull registry.company.com/user-service:v1.0.0
```

## 🔧 Troubleshooting

### Problemas Comunes

#### 1. Contenedor no Inicia
```bash
# Ver logs detallados
docker logs user-service

# Ejecutar interactivamente para debug
docker run -it --rm user-service:latest sh

# Verificar variables de entorno
docker exec user-service env
```

#### 2. No se Conecta a la Base de Datos
```bash
# Verificar que postgres esté corriendo
docker-compose ps

# Verificar conectividad de red
docker exec user-service ping postgres

# Verificar variables de entorno de BD
docker exec user-service env | grep DB_
```

#### 3. Puerto ya en Uso
```bash
# Ver qué está usando el puerto
netstat -tlnp | grep 8081

# Cambiar puerto en docker-compose.yml
ports:
  - "8082:8080"  # Usar puerto 8082 en lugar de 8081
```

#### 4. Problemas de Permisos
```bash
# Verificar permisos del archivo Firebase
ls -la firebase-service-account.json

# Cambiar permisos si es necesario
chmod 644 firebase-service-account.json

# Verificar usuario dentro del contenedor
docker exec user-service whoami
```

### Debugging

#### Logs Detallados
```bash
# Logs con timestamps
docker-compose logs -t

# Seguir logs en tiempo real
docker-compose logs -f app

# Logs de los últimos 100 líneas
docker-compose logs --tail=100 app
```

#### Inspeccionar Contenedores
```bash
# Información detallada del contenedor
docker inspect user-service

# Estadísticas de uso
docker stats user-service

# Procesos corriendo
docker exec user-service ps aux
```

#### Network Debugging
```bash
# Listar redes
docker network ls

# Inspeccionar red
docker network inspect it-app_user_default

# Probar conectividad
docker exec user-service nslookup postgres
```

### Optimización

#### Reducir Tamaño de Imagen
```dockerfile
# Usar imagen base más pequeña
FROM alpine:latest

# Limpiar cache después de instalar paquetes
RUN apk add --no-cache ca-certificates && \
    rm -rf /var/cache/apk/*

# Usar multi-stage build
FROM golang:alpine AS builder
# ... build steps ...
FROM scratch
COPY --from=builder /app/main .
```

#### Mejorar Performance
```yaml
# docker-compose.yml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

### Monitoreo

#### Health Checks
```yaml
# En docker-compose.yml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

#### Métricas
```bash
# Estadísticas en tiempo real
docker stats

# Uso de espacio
docker system df

# Limpiar recursos no utilizados
docker system prune -a
```

---

## 📚 Comandos de Referencia Rápida

```bash
# Desarrollo
docker-compose up --build
docker-compose logs -f app
docker-compose down

# Producción
docker build -t user-service:prod -f Dockerfile.prod .
docker run -d --name user-service-prod user-service:prod

# Debugging
docker exec -it user-service sh
docker logs user-service
docker inspect user-service

# Limpieza
docker system prune
docker volume prune
docker image prune
```

---

## 🔗 Enlaces Útiles

- [🏠 README Principal](../README.md)
- [📖 Guía Completa](GUIDE.md)
- [🗄️ Base de Datos](DATABASE.md)
- [🔥 Firebase Setup](FIREBASE.md)

---

**¿Problemas con Docker?** Revisa los logs, verifica la configuración y asegúrate de que todos los servicios estén corriendo correctamente.