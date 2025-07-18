# 🚀 User Service - Microservicio de Usuarios

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-12+-316192?style=for-the-badge&logo=postgresql)
![Firebase](https://img.shields.io/badge/Firebase-Auth-FFCA28?style=for-the-badge&logo=firebase)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)

**Microservicio robusto y escalable para gestión de usuarios con autenticación Firebase**

[📚 Documentación](#-documentación) • [🚀 Inicio Rápido](#-inicio-rápido) • [🔧 API](#-api-endpoints) • [🐳 Docker](#-docker)

</div>

---

## 🌟 Características Principales

- **🏗️ Arquitectura Limpia**: Patrón Repository, separación de responsabilidades
- **🔥 Firebase Auth**: Integración completa con autenticación Firebase
- **📊 Logging Estructurado**: Logs JSON con Logrus, múltiples niveles
- **🛡️ Rate Limiting**: Protección contra abuso de API configurable
- **✅ Validaciones Robustas**: Validación completa de datos con go-playground/validator
- **🌐 CORS Configurado**: Soporte completo para aplicaciones web
- **🐳 Docker Ready**: Configuración completa con Docker Compose
- **⚡ Middleware Avanzado**: Auth, logging, CORS, rate limiting
- **📈 Escalable**: Diseñado para alta concurrencia y rendimiento

## 📋 Requisitos del Sistema

| Componente | Versión Mínima | Recomendada |
|------------|----------------|-------------|
| **Go** | 1.21+ | 1.21+ |
| **PostgreSQL** | 12+ | 15+ |
| **Docker** | 20.10+ | 24.0+ |
| **Docker Compose** | 2.0+ | 2.20+ |

## 🚀 Inicio Rápido

### Opción 1: Desarrollo Local

```bash
# 1. Clonar el repositorio
git clone <repository-url>
cd it-app_user

# 2. Configurar variables de entorno
cp .env.example .env
# Editar .env con tus configuraciones

# 3. Instalar dependencias
go mod tidy

# 4. Levantar base de datos
docker-compose up postgres -d

# 5. Ejecutar el servicio
go run cmd/main.go
```

### Opción 2: Docker (Recomendado)

```bash
# Todo en uno - Base de datos + Servicio
docker-compose up --build
```

## 📚 Documentación

| Documento | Descripción |
|-----------|-------------|
| [📖 Guía Completa](docs/GUIDE.md) | Documentación detallada del proyecto |
| [🏗️ Arquitectura](docs/ARCHITECTURE.md) | Diseño y estructura del sistema |
| [🔧 API Reference](docs/API.md) | Documentación completa de endpoints |
| [�  Firebase Setup](docs/FIREBASE.md) | Configuración de Firebase |
| [🐳 Docker Guide](docs/DOCKER.md) | Guía de Docker y despliegue |
| [🗄️ Database](docs/DATABASE.md) | Esquema y migraciones de BD |
| [⚙️ Configuration](docs/CONFIG.md) | Variables de entorno y configuración |

## 🔧 API Endpoints

### 🌐 Públicos
- `GET /health` - Health check
- `GET /users` - Listar usuarios
- `POST /users/create` - Crear usuario
- `POST /auth/login` - Iniciar sesión

### 🔒 Protegidos (Firebase Auth)
- `PUT /users/{id}` - Actualizar usuario
- `DELETE /users/{id}` - Eliminar usuario
- `POST /tokens/verify` - Verificar token
- `GET /auth/profile` - Perfil del usuario

[📋 Ver todos los endpoints](ENDPOINTS.md)

## ⚙️ Variables de Entorno

```bash
# Base de Datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=itapp

# Servidor
PORT=8081
ENVIRONMENT=development
LOG_LEVEL=info

# Rate Limiting
RATE_LIMIT_RPS=100
RATE_LIMIT_BURST=200

# Firebase
FIREBASE_PROJECT_ID=innovatech-app
```

## 🏗️ Estructura del Proyecto

```
it-app_user/
├── 📁 cmd/                    # Punto de entrada
│   └── main.go
├── 📁 internal/               # Código interno
│   ├── 📁 config/            # Configuración
│   ├── 📁 database/          # Conexión BD
│   ├── 📁 handlers/          # HTTP handlers
│   ├── 📁 middleware/        # Middleware HTTP
│   ├── 📁 models/            # Modelos de datos
│   ├── 📁 repositories/      # Capa de datos
│   ├── 📁 routes/            # Definición de rutas
│   ├── 📁 server/            # Servidor HTTP
│   ├── 📁 logger/            # Sistema de logging
│   └── 📁 validator/         # Validaciones
├── 📁 pkg/                   # Paquetes públicos
│   └── 📁 firebase/          # Cliente Firebase
├── 📁 docs/                  # Documentación
└── 📁 scripts/               # Scripts de utilidad
```

## 🐳 Docker

### Desarrollo
```bash
# Solo base de datos
docker-compose up postgres -d

# Servicio completo
docker-compose up --build
```

### Producción
```bash
# Construir imagen
docker build -t user-service:latest .

# Ejecutar
docker run -p 8080:8080 \
  -e DB_HOST=your-db \
  -e FIREBASE_PROJECT_ID=your-project \
  user-service:latest
```

## 🧪 Testing

```bash
# Todos los tests
go test ./...

# Con cobertura
go test -cover ./...

# Tests específicos
go test ./internal/handlers -v
```

## 📊 Monitoreo y Logs

### Logs Estructurados (JSON)
```json
{
  "level": "info",
  "msg": "HTTP Request",
  "method": "GET",
  "path": "/users",
  "status_code": 200,
  "duration": 45,
  "time": "2024-01-01T12:00:00Z"
}
```

### Métricas Disponibles
- Requests por segundo
- Latencia de respuesta
- Errores por endpoint
- Uso de memoria y CPU

## 🔒 Seguridad

- ✅ **Rate Limiting** configurable
- ✅ **Validación de entrada** robusta
- ✅ **CORS** configurado
- ✅ **Firebase Auth** integrado
- ✅ **Logging de seguridad**
- ✅ **Headers de seguridad**

## 🚀 Despliegue

### Entornos Soportados
- **Local Development**
- **Docker Containers**
- **Kubernetes**
- **Cloud Platforms** (AWS, GCP, Azure)

### Variables de Producción
```bash
ENVIRONMENT=production
LOG_LEVEL=warn
RATE_LIMIT_RPS=50
DB_MAX_CONNECTIONS=100
```

## 🤝 Contribución

1. **Fork** el proyecto
2. **Crear** rama feature (`git checkout -b feature/nueva-funcionalidad`)
3. **Commit** cambios (`git commit -am 'Add: nueva funcionalidad'`)
4. **Push** a la rama (`git push origin feature/nueva-funcionalidad`)
5. **Crear** Pull Request

### Estándares de Código
- Seguir [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Tests para nuevas funcionalidades
- Documentación actualizada
- Commits descriptivos

## 📄 Licencia

Este proyecto está bajo la **Licencia MIT**. Ver [LICENSE](LICENSE) para más detalles.

## 🆘 Soporte

- 📧 **Email**: soporte@innovatech.com
- 🐛 **Issues**: [GitHub Issues](https://github.com/your-repo/issues)
- �  **Docs**: [Documentación Completa](docs/)
- 💬 **Discord**: [Servidor de Discord](#)

---

<div align="center">

**Hecho con ❤️ por el equipo de InnovaTech**

[⬆️ Volver arriba](#-user-service---microservicio-de-usuarios)

</div>