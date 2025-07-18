# 🔧 API Reference - User Service

## Índice
- [🌐 Base URL](#-base-url)
- [🔐 Autenticación](#-autenticación)
- [📊 Códigos de Respuesta](#-códigos-de-respuesta)
- [👤 Usuarios](#-usuarios)
- [🔐 Autenticación](#-autenticación-1)
- [🎫 Tokens](#-tokens)
- [🔑 Password Reset](#-password-reset)
- [📧 Email Verification](#-email-verification)
- [🔓 Login & Sessions](#-login--sessions)
- [📝 Ejemplos de Uso](#-ejemplos-de-uso)
- [🚨 Manejo de Errores](#-manejo-de-errores)
- [📈 Rate Limiting](#-rate-limiting)

## 🌐 Base URL

```
http://localhost:8081
```

**Producción**: `https://api.innovatech.com/user-service`

## 🔐 Autenticación

### Bearer Token (Firebase)
Para endpoints protegidos, incluir el token de Firebase en el header:

```http
Authorization: Bearer <firebase-id-token>
```

### Obtener Token (Cliente)
```javascript
// En tu aplicación cliente
import { getAuth, signInWithEmailAndPassword } from 'firebase/auth';

const auth = getAuth();
const userCredential = await signInWithEmailAndPassword(auth, email, password);
const idToken = await userCredential.user.getIdToken();

// Usar el idToken en las requests
fetch('/api/protected-endpoint', {
  headers: {
    'Authorization': `Bearer ${idToken}`,
    'Content-Type': 'application/json'
  }
});
```

## 📊 Códigos de Respuesta

| Código | Descripción | Uso |
|--------|-------------|-----|
| `200` | OK | Operación exitosa |
| `201` | Created | Recurso creado exitosamente |
| `400` | Bad Request | Datos de entrada inválidos |
| `401` | Unauthorized | Token faltante o inválido |
| `403` | Forbidden | Sin permisos para la operación |
| `404` | Not Found | Recurso no encontrado |
| `409` | Conflict | Conflicto (ej: email duplicado) |
| `429` | Too Many Requests | Rate limit excedido |
| `500` | Internal Server Error | Error interno del servidor |

## 👤 Usuarios

### Obtener Todos los Usuarios
```http
GET /users
```

**Query Parameters:**
- `limit` (opcional): Número máximo de usuarios (default: 50, max: 100)
- `offset` (opcional): Número de usuarios a saltar (default: 0)

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "firebase_id": "firebase_user_123",
      "email": "usuario@ejemplo.com",
      "email_verified": true,
      "username": "usuario123",
      "first_name": "Juan",
      "last_name": "Pérez",
      "provider": "firebase",
      "provider_id": "firebase_user_123",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "login_count": 5,
      "last_login_at": "2024-01-01T12:00:00Z",
      "last_login_ip": "192.168.1.1",
      "last_login_device": "Chrome/Windows",
      "disabled": false,
      "status": "active"
    }
  ],
  "count": 1,
  "limit": 50,
  "offset": 0,
  "message": "Users retrieved successfully"
}
```

### Obtener Usuario por ID
```http
GET /users/{id}
```

**Path Parameters:**
- `id`: ID del usuario (integer)

**Response:**
```json
{
  "data": {
    "id": 1,
    "firebase_id": "firebase_user_123",
    "email": "usuario@ejemplo.com",
    "username": "usuario123",
    "first_name": "Juan",
    "last_name": "Pérez",
    "status": "active"
  },
  "message": "User retrieved successfully"
}
```

### Obtener Usuario por Firebase ID
```http
GET /users/firebase/{firebase_id}
```

**Path Parameters:**
- `firebase_id`: Firebase ID del usuario

### Obtener Usuario por Username
```http
GET /users/username/{username}
```

**Path Parameters:**
- `username`: Username del usuario

### Obtener Usuario por Email
```http
GET /users/email/{email}
```

**Path Parameters:**
- `email`: Email del usuario

### Buscar Usuarios
```http
GET /users/search?q={query}
```

**Query Parameters:**
- `q`: Término de búsqueda (busca en nombre, email, username)
- `limit` (opcional): Número máximo de resultados (default: 20, max: 50)
- `offset` (opcional): Número de resultados a saltar (default: 0)

**Response:**
```json
{
  "data": [...],
  "count": 5,
  "query": "juan",
  "limit": 20,
  "offset": 0,
  "message": "Search completed successfully"
}
```

### Contar Usuarios
```http
GET /users/count
```

**Response:**
```json
{
  "total": 1250,
  "message": "Users counted successfully"
}
```

### Crear Usuario 🔒
```http
POST /users/create
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "firebase_id": "firebase_user_123",
  "email": "usuario@ejemplo.com",
  "username": "usuario123",
  "first_name": "Juan",
  "last_name": "Pérez",
  "provider": "firebase",
  "provider_id": "firebase_user_123",
  "status": "active"
}
```

**Required Fields:**
- `firebase_id`
- `email`
- `username`

**Response:**
```json
{
  "data": {
    "id": 1,
    "firebase_id": "firebase_user_123",
    "email": "usuario@ejemplo.com",
    "username": "usuario123",
    "first_name": "Juan",
    "last_name": "Pérez",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "message": "User created successfully"
}
```

### Actualizar Usuario 🔒
```http
PUT /users/{id}
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "username": "nuevo_usuario",
  "first_name": "Juan Carlos",
  "last_name": "Pérez García",
  "status": "active",
  "email_verified": true,
  "disabled": false
}
```

### Eliminar Usuario 🔒
```http
DELETE /users/{id}
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "User deleted successfully"
}
```

### Obtener Usuarios Activos 🔒
```http
GET /users/active
Authorization: Bearer <token>
```

### Actualizar Info de Login 🔒
```http
POST /users/{id}/login
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "login_ip": "192.168.1.1",
  "login_device": "Chrome/Windows"
}
```

## 🔐 Autenticación

### Login
```http
POST /auth/login
```

**Request Body:**
```json
{
  "id_token": "firebase-id-token-here"
}
```

**Response:**
```json
{
  "user": {
    "firebase_id": "firebase_user_123",
    "email": "usuario@ejemplo.com",
    "email_verified": true,
    "username": "usuario123",
    "status": "active"
  },
  "message": "Login successful"
}
```

### Logout
```http
POST /auth/logout
```

**Response:**
```json
{
  "message": "Logout successful"
}
```

### Verificar Estado de Autenticación
```http
GET /auth/status
Authorization: Bearer <token>
```

**Response:**
```json
{
  "authenticated": true,
  "user": {
    "firebase_id": "firebase_user_123",
    "email": "usuario@ejemplo.com",
    "email_verified": true,
    "display_name": "Juan Pérez"
  },
  "expires_at": 1640995200
}
```

### Refrescar Token
```http
POST /auth/refresh
```

**Response:**
```json
{
  "message": "Token refresh should be handled on the client side using Firebase SDK"
}
```

### Obtener Perfil 🔒
```http
GET /auth/profile
Authorization: Bearer <token>
```

### Actualizar Perfil 🔒
```http
PUT /auth/profile
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "display_name": "Juan Carlos Pérez",
  "photo_url": "https://example.com/photo.jpg"
}
```

### Cambiar Contraseña 🔒
```http
POST /auth/change-password
Authorization: Bearer <token>
```

### Revocar Todos los Tokens 🔒
```http
POST /auth/revoke-tokens
Authorization: Bearer <token>
```

### Obtener Sesiones Activas 🔒
```http
GET /auth/sessions
Authorization: Bearer <token>
```

### Revocar Sesión Específica 🔒
```http
DELETE /auth/sessions/{session_id}
Authorization: Bearer <token>
```

## 🎫 Tokens

### Verificar Token
```http
POST /tokens/verify
```

**Request Body:**
```json
{
  "id_token": "firebase-id-token-here"
}
```

**Response:**
```json
{
  "valid": true,
  "user": {
    "firebase_id": "firebase_user_123",
    "email": "usuario@ejemplo.com",
    "email_verified": true,
    "username": "usuario123"
  },
  "expires_at": 1640995200,
  "issued_at": 1640991600
}
```

### Refrescar Token
```http
POST /tokens/refresh
```

### Crear Token Personalizado
```http
POST /tokens/custom
```

**Request Body:**
```json
{
  "uid": "firebase_user_123",
  "claims": {
    "role": "admin",
    "permissions": ["read", "write"]
  }
}
```

### Revocar Token 🔒
```http
POST /tokens/revoke
Authorization: Bearer <token>
```

### Revocar Todos los Tokens 🔒
```http
POST /tokens/revoke-all
Authorization: Bearer <token>
```

### Obtener Info del Token 🔒
```http
GET /tokens/info
Authorization: Bearer <token>
```

### Validar Token 🔒
```http
POST /tokens/validate
Authorization: Bearer <token>
```

## 🔑 Password Reset

### Solicitar Reset de Contraseña
```http
POST /password/reset/request
```

**Request Body:**
```json
{
  "email": "usuario@ejemplo.com"
}
```

**Response:**
```json
{
  "message": "If the email exists, a password reset link has been sent"
}
```

### Verificar Código de Reset
```http
POST /password/reset/verify
```

**Request Body:**
```json
{
  "code": "123456"
}
```

### Confirmar Reset de Contraseña
```http
POST /password/reset/confirm
```

**Request Body:**
```json
{
  "code": "123456",
  "new_password": "nueva_contraseña_segura"
}
```

### Validar Token de Reset
```http
POST /password/reset/validate-token
```

**Request Body:**
```json
{
  "token": "reset-token-here"
}
```

### Cambiar Contraseña 🔒
```http
POST /password/change
Authorization: Bearer <token>
```

### Verificar Fortaleza de Contraseña 🔒
```http
POST /password/strength-check
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "password": "mi_contraseña_123"
}
```

**Response:**
```json
{
  "data": {
    "is_valid": true,
    "score": 4,
    "feedback": [],
    "requirements": {
      "min_length": true,
      "has_uppercase": true,
      "has_lowercase": true,
      "has_numbers": true,
      "has_symbols": false
    }
  },
  "message": "Password strength checked"
}
```

### Obtener Historial de Contraseñas 🔒
```http
GET /password/history
Authorization: Bearer <token>
```

### Obtener Política de Contraseñas 🔒
```http
GET /password/policy
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "min_length": 8,
    "require_uppercase": true,
    "require_lowercase": true,
    "require_numbers": true,
    "require_symbols": false,
    "forbidden_words": ["password", "123456", "qwerty"],
    "max_age_days": 90
  },
  "message": "Password policy retrieved successfully"
}
```

## 📧 Email Verification

### Enviar Email de Verificación
```http
POST /email/send-verification
```

**Request Body:**
```json
{
  "email": "usuario@ejemplo.com",
  "language": "es"
}
```

### Verificar Email
```http
POST /email/verify
```

**Request Body:**
```json
{
  "id_token": "firebase-id-token-here"
}
```

### Verificar Email con Código
```http
POST /email/verify-code
```

**Request Body:**
```json
{
  "email": "usuario@ejemplo.com",
  "verification_code": "123456"
}
```

### Reenviar Email de Verificación
```http
POST /email/resend
```

**Request Body:**
```json
{
  "email": "usuario@ejemplo.com"
}
```

### Verificar Estado de Verificación
```http
GET /email/status?email=usuario@ejemplo.com
```

**Response:**
```json
{
  "email_verified": true,
  "email": "usuario@ejemplo.com",
  "can_resend": false,
  "attempts_left": 3
}
```

### Obtener Mi Estado de Verificación 🔒
```http
GET /email/my-status
Authorization: Bearer <token>
```

### Actualizar Email 🔒
```http
POST /email/update-email
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "new_email": "nuevo@ejemplo.com"
}
```

### Obtener Historial de Verificaciones 🔒
```http
GET /email/verification-history
Authorization: Bearer <token>
```

### Obtener Configuraciones de Email 🔒
```http
GET /email/settings
Authorization: Bearer <token>
```

### Actualizar Configuraciones de Email 🔒
```http
PUT /email/settings
Authorization: Bearer <token>
```

## 🔓 Login & Sessions

### Trackear Login de Usuario
```http
POST /login/track
```

**Request Body:**
```json
{
  "user_id": 1,
  "firebase_id": "firebase_user_123",
  "login_ip": "192.168.1.1",
  "login_device": "Chrome/Windows",
  "login_method": "firebase",
  "user_agent": "Mozilla/5.0...",
  "success": true
}
```

### Obtener Historial de Login
```http
GET /login/history/{user_id}
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "login_time": "2024-01-01T10:00:00Z",
      "login_ip": "192.168.1.1",
      "login_device": "Chrome/Windows",
      "login_method": "firebase",
      "success": true
    }
  ],
  "count": 1,
  "message": "Login history retrieved successfully"
}
```

### Obtener Intentos de Login
```http
GET /login/attempts/{email}
```

### Verificación de Seguridad
```http
POST /login/security-check
```

**Request Body:**
```json
{
  "email": "usuario@ejemplo.com",
  "ip_address": "192.168.1.1",
  "user_agent": "Mozilla/5.0..."
}
```

**Response:**
```json
{
  "data": {
    "is_safe": true,
    "risk_level": "low",
    "blocked": false,
    "requires_2fa": false,
    "suspicious_activity": false,
    "recommendations": []
  },
  "message": "Security check completed"
}
```

### Obtener Mi Historial de Login 🔒
```http
GET /login/my-history?limit=20&offset=0
Authorization: Bearer <token>
```

### Obtener Sesiones Activas 🔒
```http
GET /login/active-sessions
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": [
    {
      "session_id": "session_123",
      "device": "Chrome/Windows",
      "ip_address": "192.168.1.1",
      "location": "New York, US",
      "login_time": "2024-01-01T10:00:00Z",
      "last_seen": "2024-01-01T12:00:00Z",
      "is_current": true
    }
  ],
  "count": 1,
  "message": "Active sessions retrieved successfully"
}
```

### Terminar Sesión Específica 🔒
```http
DELETE /login/terminate-session/{session_id}
Authorization: Bearer <token>
```

### Terminar Todas las Sesiones 🔒
```http
DELETE /login/terminate-all-sessions
Authorization: Bearer <token>
```

### Obtener Actividad Sospechosa 🔒
```http
GET /login/suspicious-activity
Authorization: Bearer <token>
```

## 📝 Ejemplos de Uso

### Flujo Completo de Registro y Login

#### 1. Crear Usuario
```bash
curl -X POST http://localhost:8081/users/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <firebase-token>" \
  -d '{
    "firebase_id": "firebase_user_123",
    "email": "juan@ejemplo.com",
    "username": "juan123",
    "first_name": "Juan",
    "last_name": "Pérez",
    "status": "active"
  }'
```

#### 2. Login
```bash
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "id_token": "<firebase-id-token>"
  }'
```

#### 3. Verificar Token
```bash
curl -X POST http://localhost:8081/tokens/verify \
  -H "Content-Type: application/json" \
  -d '{
    "id_token": "<firebase-id-token>"
  }'
```

#### 4. Obtener Perfil
```bash
curl -X GET http://localhost:8081/auth/profile \
  -H "Authorization: Bearer <firebase-token>"
```

### Flujo de Reset de Contraseña

#### 1. Solicitar Reset
```bash
curl -X POST http://localhost:8081/password/reset/request \
  -H "Content-Type: application/json" \
  -d '{
    "email": "juan@ejemplo.com"
  }'
```

#### 2. Verificar Código
```bash
curl -X POST http://localhost:8081/password/reset/verify \
  -H "Content-Type: application/json" \
  -d '{
    "code": "123456"
  }'
```

#### 3. Confirmar Reset
```bash
curl -X POST http://localhost:8081/password/reset/confirm \
  -H "Content-Type: application/json" \
  -d '{
    "code": "123456",
    "new_password": "nueva_contraseña_segura"
  }'
```

## 🚨 Manejo de Errores

### Formato de Error Estándar
```json
{
  "error": "Validation failed: Field 'email' failed validation: email",
  "code": "VALIDATION_ERROR",
  "details": {
    "field": "email",
    "value": "invalid-email",
    "constraint": "email format"
  },
  "timestamp": "2024-01-01T12:00:00Z",
  "path": "/users/create"
}
```

### Códigos de Error Comunes

| Código | Descripción | Solución |
|--------|-------------|----------|
| `VALIDATION_ERROR` | Datos de entrada inválidos | Verificar formato de datos |
| `UNAUTHORIZED` | Token faltante o inválido | Incluir token válido |
| `USER_NOT_FOUND` | Usuario no existe | Verificar ID de usuario |
| `EMAIL_ALREADY_EXISTS` | Email ya registrado | Usar email diferente |
| `USERNAME_TAKEN` | Username ya en uso | Elegir username diferente |
| `RATE_LIMIT_EXCEEDED` | Demasiadas requests | Esperar antes de reintentar |
| `FIREBASE_ERROR` | Error de Firebase | Verificar configuración |

### Ejemplos de Errores

#### Error de Validación
```json
{
  "error": "Validation failed: Field 'username' failed validation: min",
  "details": {
    "field": "username",
    "constraint": "minimum 3 characters"
  }
}
```

#### Error de Autenticación
```json
{
  "error": "Invalid token",
  "code": "UNAUTHORIZED"
}
```

#### Error de Conflicto
```json
{
  "error": "User already exists",
  "code": "EMAIL_ALREADY_EXISTS",
  "details": {
    "email": "usuario@ejemplo.com"
  }
}
```

## 📈 Rate Limiting

### Límites por Defecto
- **Requests por segundo**: 100
- **Burst capacity**: 200
- **Por IP**: Aplicado por dirección IP

### Headers de Rate Limiting
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

### Configuración
```bash
# Variables de entorno
RATE_LIMIT_RPS=100      # Requests por segundo
RATE_LIMIT_BURST=200    # Capacidad de burst
```

### Manejo de Rate Limiting
```javascript
// Ejemplo de manejo en cliente
async function makeRequest(url, options) {
  try {
    const response = await fetch(url, options);
    
    if (response.status === 429) {
      const retryAfter = response.headers.get('Retry-After');
      console.log(`Rate limited. Retry after ${retryAfter} seconds`);
      
      // Esperar y reintentar
      await new Promise(resolve => setTimeout(resolve, retryAfter * 1000));
      return makeRequest(url, options);
    }
    
    return response;
  } catch (error) {
    console.error('Request failed:', error);
    throw error;
  }
}
```

---

## 🔗 Enlaces Útiles

- [🏠 Inicio](../README.md)
- [📖 Guía Completa](GUIDE.md)
- [🏗️ Arquitectura](ARCHITECTURE.md)
- [🔥 Firebase Setup](FIREBASE.md)
- [🐳 Docker Guide](DOCKER.md)

---

**¿Necesitas ayuda?** Contacta al equipo de desarrollo o crea un issue en GitHub.