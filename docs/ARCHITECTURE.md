# 🏗️ Arquitectura del Sistema - User Service

## Índice
- [🎯 Visión General](#-visión-general)
- [🏛️ Patrones de Arquitectura](#️-patrones-de-arquitectura)
- [📦 Estructura de Capas](#-estructura-de-capas)
- [🔄 Flujo de Datos](#-flujo-de-datos)
- [🗄️ Diseño de Base de Datos](#️-diseño-de-base-de-datos)
- [🔌 Integraciones Externas](#-integraciones-externas)
- [🛡️ Seguridad](#️-seguridad)
- [📈 Escalabilidad](#-escalabilidad)
- [🔧 Decisiones de Diseño](#-decisiones-de-diseño)

## 🎯 Visión General

El User Service está diseñado siguiendo principios de **Clean Architecture** y **Domain-Driven Design (DDD)**, proporcionando una base sólida para la gestión de usuarios en sistemas distribuidos.

### Principios de Diseño
- **Separación de Responsabilidades**: Cada capa tiene una responsabilidad específica
- **Inversión de Dependencias**: Las capas internas no dependen de las externas
- **Testabilidad**: Diseño que facilita testing unitario e integración
- **Escalabilidad**: Arquitectura preparada para crecimiento horizontal
- **Mantenibilidad**: Código organizado y fácil de mantener

## 🏛️ Patrones de Arquitectura

### Clean Architecture
```
┌─────────────────────────────────────────────────────────┐
│                 External Interfaces                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │     Web     │  │   Database  │  │  Firebase   │     │
│  │  Framework  │  │             │  │    Auth     │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│              Interface Adapters                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │  Controllers│  │ Repositories│  │  Presenters │     │
│  │ (Handlers)  │  │             │  │             │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                Application Business Rules               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   Use Cases │  │  Services   │  │ Validators  │     │
│  │             │  │             │  │             │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                Enterprise Business Rules                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │  Entities   │  │   Models    │  │  Interfaces │     │
│  │             │  │             │  │             │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
```

### Repository Pattern
```go
// Abstracción de la capa de datos
type UserRepositoryInterface interface {
    GetByID(id uint) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id uint) error
}

// Implementación concreta
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) GetByID(id uint) (*User, error) {
    // Implementación específica de GORM
}
```

### Dependency Injection
```go
// Los handlers reciben dependencias inyectadas
type UserHandler struct {
    userRepo repositories.UserRepositoryInterface
    logger   logger.Logger
}

func NewUserHandler(userRepo repositories.UserRepositoryInterface) *UserHandler {
    return &UserHandler{
        userRepo: userRepo,
        logger:   logger.GetLogger(),
    }
}
```

## 📦 Estructura de Capas

### 1. **Presentation Layer** (`internal/handlers/`, `internal/routes/`)
**Responsabilidad**: Manejo de HTTP requests/responses

```go
// Ejemplo de Handler
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    // 1. Extraer parámetros de la request
    // 2. Validar entrada
    // 3. Llamar a la capa de negocio
    // 4. Formatear respuesta
    // 5. Enviar respuesta HTTP
}
```

**Componentes**:
- **Handlers**: Controladores HTTP
- **Routes**: Definición de rutas
- **Middleware**: Interceptores de requests

### 2. **Application Layer** (`internal/services/`, `internal/validator/`)
**Responsabilidad**: Lógica de aplicación y orquestación

```go
// Ejemplo de Service
type UserService struct {
    userRepo repositories.UserRepositoryInterface
    validator validator.Validator
}

func (s *UserService) CreateUser(req *CreateUserRequest) (*User, error) {
    // 1. Validar datos de entrada
    // 2. Aplicar reglas de negocio
    // 3. Coordinar con repositorios
    // 4. Manejar transacciones
    return user, nil
}
```

**Componentes**:
- **Services**: Lógica de aplicación
- **Validators**: Validación de datos
- **DTOs**: Objetos de transferencia de datos

### 3. **Domain Layer** (`internal/models/`)
**Responsabilidad**: Entidades de dominio y reglas de negocio

```go
// Ejemplo de Entity
type User struct {
    ID            uint      `gorm:"primaryKey"`
    FirebaseID    string    `gorm:"uniqueIndex"`
    Email         string    `gorm:"uniqueIndex"`
    Username      string    `gorm:"uniqueIndex"`
    // ... otros campos
}

// Métodos de dominio
func (u *User) IsActive() bool {
    return u.Status == "active" && !u.Disabled
}
```

**Componentes**:
- **Entities**: Objetos de dominio
- **Value Objects**: Objetos inmutables
- **Domain Services**: Lógica de dominio compleja

### 4. **Infrastructure Layer** (`internal/database/`, `pkg/firebase/`)
**Responsabilidad**: Implementaciones técnicas y servicios externos

```go
// Ejemplo de Repository Implementation
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) GetByID(id uint) (*User, error) {
    var user User
    err := r.db.First(&user, id).Error
    return &user, err
}
```

**Componentes**:
- **Repositories**: Acceso a datos
- **External Services**: APIs externas
- **Configuration**: Configuración del sistema

## 🔄 Flujo de Datos

### Request Flow
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │───▶│ Middleware  │───▶│   Handler   │───▶│   Service   │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
                          │                   │                   │
                          ▼                   ▼                   ▼
                   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
                   │   Logging   │    │ Validation  │    │ Repository  │
                   └─────────────┘    └─────────────┘    └─────────────┘
                                                                │
                                                                ▼
                                                        ┌─────────────┐
                                                        │  Database   │
                                                        └─────────────┘
```

### Ejemplo de Flujo Completo
```go
// 1. Request llega al middleware
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        duration := time.Since(start)
        log.Info("Request processed", "duration", duration)
    })
}

// 2. Handler procesa la request
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // 3. Service maneja la lógica de negocio
    user, err := h.userService.CreateUser(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // 4. Response
    json.NewEncoder(w).Encode(user)
}

// 5. Service coordina operaciones
func (s *UserService) CreateUser(req *CreateUserRequest) (*User, error) {
    // Validación
    if err := s.validator.Validate(req); err != nil {
        return nil, err
    }
    
    // Crear entidad
    user := &User{
        Email:    req.Email,
        Username: req.Username,
    }
    
    // Persistir
    return s.userRepo.Create(user)
}
```

## 🗄️ Diseño de Base de Datos

### Modelo Entidad-Relación
```
┌─────────────────────────────────────────────────────────┐
│                        USERS                            │
├─────────────────────────────────────────────────────────┤
│ id (PK)                    │ SERIAL                     │
│ firebase_id (UK)           │ VARCHAR(128)               │
│ email (UK)                 │ VARCHAR(255)               │
│ email_verified             │ BOOLEAN                    │
│ username (UK)              │ VARCHAR(50)                │
│ first_name                 │ VARCHAR(100)               │
│ last_name                  │ VARCHAR(100)               │
│ provider                   │ VARCHAR(50)                │
│ provider_id                │ VARCHAR(128)               │
│ created_at                 │ TIMESTAMP WITH TIME ZONE   │
│ updated_at                 │ TIMESTAMP WITH TIME ZONE   │
│ login_count                │ INTEGER                    │
│ last_login_at              │ TIMESTAMP WITH TIME ZONE   │
│ last_login_ip              │ INET                       │
│ last_login_device          │ VARCHAR(255)               │
│ disabled                   │ BOOLEAN                    │
│ status                     │ VARCHAR(20)                │
│ pk_aut_use_id              │ VARCHAR(128)               │
└─────────────────────────────────────────────────────────┘
                                │
                                │ 1:1
                                ▼
┌─────────────────────────────────────────────────────────┐
│                   USER_PROFILES                         │
├─────────────────────────────────────────────────────────┤
│ id (PK)                    │ SERIAL                     │
│ user_id (FK)               │ INTEGER                    │
│ avatar                     │ VARCHAR(500)               │
│ bio                        │ TEXT                       │
│ website                    │ VARCHAR(255)               │
│ location                   │ VARCHAR(100)               │
│ birthday                   │ DATE                       │
│ gender                     │ VARCHAR(20)                │
│ phone                      │ VARCHAR(20)                │
│ preferences                │ JSONB                      │
│ privacy                    │ JSONB                      │
│ created_at                 │ TIMESTAMP WITH TIME ZONE   │
│ updated_at                 │ TIMESTAMP WITH TIME ZONE   │
└─────────────────────────────────────────────────────────┘
```

### Índices Optimizados
```sql
-- Índices principales para consultas frecuentes
CREATE INDEX CONCURRENTLY idx_users_firebase_id ON users(firebase_id);
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);
CREATE INDEX CONCURRENTLY idx_users_username ON users(username);
CREATE INDEX CONCURRENTLY idx_users_status ON users(status) WHERE status = 'active';
CREATE INDEX CONCURRENTLY idx_users_created_at ON users(created_at);

-- Índices compuestos para consultas complejas
CREATE INDEX CONCURRENTLY idx_users_status_created ON users(status, created_at);
CREATE INDEX CONCURRENTLY idx_users_provider_status ON users(provider, status);
```

### Estrategias de Particionamiento
```sql
-- Particionamiento por fecha para tablas de logs
CREATE TABLE user_activity_logs (
    id SERIAL,
    user_id INTEGER,
    action VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE
) PARTITION BY RANGE (created_at);

-- Particiones mensuales
CREATE TABLE user_activity_logs_2024_01 PARTITION OF user_activity_logs
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
```

## 🔌 Integraciones Externas

### Firebase Authentication
```go
// Integración con Firebase Auth
type FirebaseAuth struct {
    client *auth.Client
    config *firebase.Config
}

func (f *FirebaseAuth) VerifyToken(ctx context.Context, token string) (*auth.Token, error) {
    return f.client.VerifyIDToken(ctx, token)
}
```

**Flujo de Autenticación**:
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │───▶│   Service   │───▶│  Firebase   │───▶│  Database   │
│             │    │             │    │    Auth     │    │             │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
      │                   │                   │                   │
      │ 1. Send Token     │ 2. Verify Token  │ 3. Get User Info  │
      │                   │                   │                   │ 4. Store/Update
      │◀──────────────────│◀──────────────────│◀──────────────────│
      │ 5. Return User    │                   │                   │
```

### Base de Datos PostgreSQL
```go
// Configuración de conexión optimizada
func NewDatabase(config Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    // Configuración del pool de conexiones
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return db, err
}
```

## 🛡️ Seguridad

### Autenticación y Autorización
```go
// Middleware de autenticación
func AuthMiddleware(firebaseAuth *firebase.Auth) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := extractToken(r)
            if token == "" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            decodedToken, err := firebaseAuth.VerifyIDToken(r.Context(), token)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }
            
            ctx := context.WithValue(r.Context(), "user_id", decodedToken.UID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Rate Limiting
```go
// Rate limiter por IP
type RateLimiter struct {
    visitors map[string]*rate.Limiter
    mu       sync.RWMutex
    rate     rate.Limit
    burst    int
}

func (rl *RateLimiter) Allow(ip string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    limiter, exists := rl.visitors[ip]
    if !exists {
        limiter = rate.NewLimiter(rl.rate, rl.burst)
        rl.visitors[ip] = limiter
    }
    
    return limiter.Allow()
}
```

### Validación de Datos
```go
// Validación robusta con go-playground/validator
type CreateUserRequest struct {
    Email     string `json:"email" validate:"required,email,max=255"`
    Username  string `json:"username" validate:"required,min=3,max=50,alphanum"`
    FirstName string `json:"first_name" validate:"max=100"`
    LastName  string `json:"last_name" validate:"max=100"`
}

func ValidateStruct(s interface{}) error {
    validate := validator.New()
    return validate.Struct(s)
}
```

## 📈 Escalabilidad

### Escalabilidad Horizontal
```yaml
# Kubernetes Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
spec:
  replicas: 5  # Múltiples instancias
  selector:
    matchLabels:
      app: user-service
  template:
    spec:
      containers:
      - name: user-service
        image: user-service:latest
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
```

### Caching Strategy
```go
// Redis para caching
type CacheService struct {
    client *redis.Client
}

func (c *CacheService) GetUser(id string) (*User, error) {
    // 1. Intentar obtener del cache
    cached, err := c.client.Get(ctx, "user:"+id).Result()
    if err == nil {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }
    
    // 2. Si no está en cache, obtener de DB
    user, err := c.userRepo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // 3. Guardar en cache
    userJSON, _ := json.Marshal(user)
    c.client.Set(ctx, "user:"+id, userJSON, time.Hour)
    
    return user, nil
}
```

### Database Scaling
```go
// Read/Write splitting
type DatabaseCluster struct {
    writeDB *gorm.DB
    readDB  *gorm.DB
}

func (dc *DatabaseCluster) Create(user *User) error {
    return dc.writeDB.Create(user).Error
}

func (dc *DatabaseCluster) GetByID(id uint) (*User, error) {
    var user User
    err := dc.readDB.First(&user, id).Error
    return &user, err
}
```

## 🔧 Decisiones de Diseño

### ¿Por qué Go?
- **Performance**: Compilado, garbage collector eficiente
- **Concurrencia**: Goroutines y channels nativos
- **Simplicidad**: Sintaxis simple, fácil de mantener
- **Ecosistema**: Excelentes librerías para web services

### ¿Por qué GORM?
- **Productividad**: ORM maduro con muchas características
- **Migraciones**: Auto-migrate para desarrollo rápido
- **Hooks**: Callbacks para lógica de negocio
- **Associations**: Manejo de relaciones simplificado

### ¿Por qué Firebase Auth?
- **Seguridad**: Manejo seguro de autenticación
- **Escalabilidad**: Servicio gestionado, escala automáticamente
- **Integraciones**: Fácil integración con frontend
- **Características**: 2FA, social login, etc.

### ¿Por qué PostgreSQL?
- **ACID**: Transacciones confiables
- **JSON Support**: JSONB para datos semi-estructurados
- **Performance**: Excelente para cargas de trabajo complejas
- **Extensibilidad**: Muchas extensiones disponibles

### Trade-offs Considerados

#### Microservicio vs Monolito
**Elegido**: Microservicio
- ✅ **Pros**: Escalabilidad independiente, tecnologías específicas
- ❌ **Cons**: Complejidad de red, eventual consistency

#### ORM vs SQL Raw
**Elegido**: GORM (ORM)
- ✅ **Pros**: Productividad, type safety, migraciones
- ❌ **Cons**: Performance overhead, menos control

#### Síncrono vs Asíncrono
**Elegido**: Síncrono con capacidad asíncrona
- ✅ **Pros**: Simplicidad, debugging fácil
- ❌ **Cons**: Latencia en operaciones lentas

### Patrones Evitados

#### God Object
```go
// ❌ Evitado: Una estructura que hace todo
type UserManager struct {
    // Demasiadas responsabilidades
}

// ✅ Preferido: Separación de responsabilidades
type UserHandler struct{}
type UserRepository struct{}
type UserValidator struct{}
```

#### Tight Coupling
```go
// ❌ Evitado: Dependencia directa
type UserService struct {
    db *gorm.DB  // Acoplado a GORM
}

// ✅ Preferido: Dependency injection
type UserService struct {
    repo UserRepositoryInterface  // Abstracción
}
```

Esta arquitectura proporciona una base sólida, escalable y mantenible para el User Service, siguiendo las mejores prácticas de la industria y permitiendo evolución futura del sistema.