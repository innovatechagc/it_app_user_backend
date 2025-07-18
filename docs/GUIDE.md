# 📖 Guía Completa - User Service

## Índice
- [🎯 Propósito del Proyecto](#-propósito-del-proyecto)
- [🏗️ Arquitectura General](#️-arquitectura-general)
- [🚀 Instalación Detallada](#-instalación-detallada)
- [⚙️ Configuración](#️-configuración)
- [🔥 Firebase Setup](#-firebase-setup)
- [🗄️ Base de Datos](#️-base-de-datos)
- [🔧 Desarrollo](#-desarrollo)
- [🧪 Testing](#-testing)
- [📊 Monitoreo](#-monitoreo)
- [🚀 Despliegue](#-despliegue)
- [🔧 Troubleshooting](#-troubleshooting)

## 🎯 Propósito del Proyecto

El **User Service** es un microservicio diseñado para gestionar usuarios de manera robusta y escalable. Proporciona:

### Funcionalidades Principales
- **Gestión de Usuarios**: CRUD completo de usuarios
- **Autenticación**: Integración con Firebase Auth
- **Autorización**: Control de acceso basado en roles
- **Validación**: Validación robusta de datos de entrada
- **Logging**: Sistema de logs estructurado
- **Rate Limiting**: Protección contra abuso de API
- **Monitoreo**: Métricas y health checks

### Casos de Uso
- **Aplicaciones Web**: Backend para aplicaciones SPA
- **Aplicaciones Móviles**: API para apps iOS/Android
- **Microservicios**: Parte de una arquitectura de microservicios
- **Sistemas Empresariales**: Gestión de usuarios corporativos

## 🏗️ Arquitectura General

### Patrón de Arquitectura
El proyecto sigue el patrón de **Arquitectura Limpia** (Clean Architecture):

```
┌─────────────────────────────────────────────────────────┐
│                    Presentation Layer                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   Routes    │  │  Handlers   │  │ Middleware  │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                    Business Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │  Services   │  │ Validators  │  │   Models    │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                     Data Layer                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │Repositories │  │  Database   │  │  Firebase   │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
```

### Componentes Principales

#### 1. **Handlers** (`internal/handlers/`)
- Manejan las peticiones HTTP
- Validan entrada y formatean salida
- Coordinan entre servicios y repositorios

#### 2. **Repositories** (`internal/repositories/`)
- Abstracción de la capa de datos
- Implementan interfaces para testabilidad
- Manejan operaciones CRUD

#### 3. **Models** (`internal/models/`)
- Estructuras de datos
- Validaciones de campo
- Mapeo de base de datos

#### 4. **Middleware** (`internal/middleware/`)
- Autenticación y autorización
- Logging de requests
- Rate limiting
- CORS

#### 5. **Services** (`pkg/`)
- Lógica de negocio
- Integración con servicios externos
- Firebase Auth

## 🚀 Instalación Detallada

### Prerrequisitos

#### 1. **Go 1.21+**
```bash
# Verificar versión
go version

# Si no tienes Go instalado:
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# macOS
brew install go

# Windows
# Descargar desde https://golang.org/dl/
```

#### 2. **PostgreSQL 12+**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib

# macOS
brew install postgresql
brew services start postgresql

# Windows
# Descargar desde https://www.postgresql.org/download/
```

#### 3. **Docker & Docker Compose**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install docker.io docker-compose

# macOS
brew install docker docker-compose

# Windows
# Descargar Docker Desktop
```

### Instalación del Proyecto

#### Opción 1: Desarrollo Local
```bash
# 1. Clonar repositorio
git clone https://github.com/your-org/it-app_user.git
cd it-app_user

# 2. Instalar dependencias
go mod tidy

# 3. Configurar variables de entorno
cp .env.example .env

# 4. Editar configuración
nano .env  # o tu editor preferido

# 5. Configurar base de datos
createdb itapp

# 6. Ejecutar migraciones (automáticas al iniciar)
go run cmd/main.go
```

#### Opción 2: Docker (Recomendado)
```bash
# 1. Clonar repositorio
git clone https://github.com/your-org/it-app_user.git
cd it-app_user

# 2. Configurar variables de entorno
cp .env.example .env

# 3. Levantar servicios
docker-compose up --build
```

## ⚙️ Configuración

### Variables de Entorno

#### Base de Datos
```bash
DB_HOST=localhost          # Host de PostgreSQL
DB_PORT=5432              # Puerto de PostgreSQL
DB_USER=postgres          # Usuario de base de datos
DB_PASSWORD=postgres      # Contraseña de base de datos
DB_NAME=itapp            # Nombre de la base de datos
```

#### Servidor
```bash
PORT=8081                 # Puerto del servicio
ENVIRONMENT=development   # Entorno (development/production)
LOG_LEVEL=info           # Nivel de logs (debug/info/warn/error)
```

#### Rate Limiting
```bash
RATE_LIMIT_RPS=100       # Requests por segundo
RATE_LIMIT_BURST=200     # Capacidad de burst
```

#### Firebase
```bash
FIREBASE_PROJECT_ID=innovatech-app  # ID del proyecto Firebase
```

### Configuración por Entorno

#### Desarrollo
```bash
ENVIRONMENT=development
LOG_LEVEL=debug
RATE_LIMIT_RPS=1000
DB_HOST=localhost
```

#### Producción
```bash
ENVIRONMENT=production
LOG_LEVEL=warn
RATE_LIMIT_RPS=50
DB_HOST=your-production-db
```

## 🔥 Firebase Setup

### 1. Crear Proyecto Firebase
1. Ve a [Firebase Console](https://console.firebase.google.com/)
2. Crea un nuevo proyecto o selecciona `innovatech-app`
3. Habilita Authentication
4. Configura métodos de autenticación (Email/Password, Google, etc.)

### 2. Obtener Service Account
1. Ve a **Configuración del proyecto** → **Cuentas de servicio**
2. Haz clic en **Generar nueva clave privada**
3. Descarga el archivo JSON
4. Guárdalo como `firebase-service-account.json` en la raíz del proyecto

### 3. Configurar Cliente
```javascript
// Tu configuración del cliente (ya tienes esto)
const firebaseConfig = {
  apiKey: "AIzaSyCbAXjfv-f0n7b91CrQ6nkn2Pt1TNSunFw",
  authDomain: "innovatech-app.firebaseapp.com",
  projectId: "innovatech-app",
  // ... resto de la configuración
};
```

## 🗄️ Base de Datos

### Esquema Principal

#### Tabla `users`
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    firebase_id VARCHAR(128) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    provider VARCHAR(50),
    provider_id VARCHAR(128),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    login_count INTEGER DEFAULT 0,
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_login_ip INET,
    last_login_device VARCHAR(255),
    disabled BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'active',
    pk_aut_use_id VARCHAR(128)
);
```

### Migraciones
Las migraciones se ejecutan automáticamente al iniciar el servicio usando GORM AutoMigrate.

### Índices
```sql
-- Índices para optimizar consultas
CREATE INDEX idx_users_firebase_id ON users(firebase_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_status ON users(status);
```

## 🔧 Desarrollo

### Estructura de Desarrollo
```bash
# Ejecutar en modo desarrollo
go run cmd/main.go

# Con hot reload (usando air)
go install github.com/cosmtrek/air@latest
air

# Ejecutar tests
go test ./...

# Ejecutar tests con cobertura
go test -cover ./...
```

### Agregar Nuevos Endpoints

#### 1. Crear Handler
```go
// internal/handlers/new_handler.go
func (h *Handler) NewEndpoint(w http.ResponseWriter, r *http.Request) {
    // Implementación
}
```

#### 2. Agregar Ruta
```go
// internal/routes/routes.go
router.HandleFunc("/new-endpoint", handler.NewEndpoint).Methods("GET")
```

#### 3. Agregar Tests
```go
// internal/handlers/new_handler_test.go
func TestNewEndpoint(t *testing.T) {
    // Tests
}
```

### Convenciones de Código

#### Naming
- **Packages**: lowercase, single word
- **Functions**: CamelCase, exported start with uppercase
- **Variables**: camelCase
- **Constants**: UPPER_CASE

#### Error Handling
```go
// Siempre manejar errores
if err != nil {
    log.WithError(err).Error("Description")
    return err
}
```

#### Logging
```go
// Usar logging estructurado
log.WithFields(map[string]interface{}{
    "user_id": userID,
    "action": "create_user",
}).Info("User created successfully")
```

## 🧪 Testing

### Tipos de Tests

#### Unit Tests
```bash
# Ejecutar tests unitarios
go test ./internal/handlers -v
go test ./internal/repositories -v
```

#### Integration Tests
```bash
# Tests de integración con base de datos
go test ./internal/integration -v
```

#### API Tests
```bash
# Tests de API completos
go test ./tests/api -v
```

### Cobertura de Tests
```bash
# Generar reporte de cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Mocking
```go
// Usar interfaces para mocking
type UserRepository interface {
    GetByID(id uint) (*User, error)
}

// Mock en tests
type MockUserRepository struct{}
func (m *MockUserRepository) GetByID(id uint) (*User, error) {
    return &User{ID: id}, nil
}
```

## 📊 Monitoreo

### Health Checks
```bash
# Health check básico
curl http://localhost:8081/health

# Response
{
  "status": "ok",
  "service": "user-service"
}
```

### Métricas Disponibles
- **Request Rate**: Requests por segundo
- **Response Time**: Latencia promedio
- **Error Rate**: Porcentaje de errores
- **Database Connections**: Conexiones activas

### Logs Estructurados
```json
{
  "level": "info",
  "msg": "HTTP Request",
  "method": "GET",
  "path": "/users",
  "status_code": 200,
  "duration": 45,
  "remote_addr": "127.0.0.1:54321",
  "user_agent": "curl/7.68.0",
  "time": "2024-01-01T12:00:00Z"
}
```

## 🚀 Despliegue

### Docker Production
```bash
# Build imagen
docker build -t user-service:v1.0.0 .

# Run container
docker run -d \
  --name user-service \
  -p 8080:8080 \
  -e DB_HOST=prod-db \
  -e ENVIRONMENT=production \
  user-service:v1.0.0
```

### Kubernetes
```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: user-service:v1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: "postgres-service"
```

### Variables de Producción
```bash
# Configuración de producción
ENVIRONMENT=production
LOG_LEVEL=warn
RATE_LIMIT_RPS=50
DB_MAX_CONNECTIONS=100
DB_TIMEOUT=30s
```

## 🔧 Troubleshooting

### Problemas Comunes

#### 1. Error de Conexión a Base de Datos
```bash
# Verificar que PostgreSQL esté corriendo
sudo systemctl status postgresql

# Verificar conexión
psql -h localhost -U postgres -d itapp
```

#### 2. Firebase Auth No Funciona
```bash
# Verificar archivo service account
ls -la firebase-service-account.json

# Verificar variables de entorno
echo $FIREBASE_PROJECT_ID
```

#### 3. Rate Limiting Muy Restrictivo
```bash
# Ajustar en .env
RATE_LIMIT_RPS=1000
RATE_LIMIT_BURST=2000
```

#### 4. Logs No Aparecen
```bash
# Cambiar nivel de log
LOG_LEVEL=debug
```

### Debugging

#### Logs Detallados
```bash
# Ejecutar con logs debug
LOG_LEVEL=debug go run cmd/main.go
```

#### Profiling
```go
// Agregar pprof para profiling
import _ "net/http/pprof"

// En main.go
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

#### Database Debugging
```bash
# Habilitar logs de SQL en GORM
DB_LOG_LEVEL=info
```

### Contacto de Soporte
- **Email**: soporte@innovatech.com
- **Slack**: #user-service-support
- **Issues**: GitHub Issues
- **Docs**: Esta documentación