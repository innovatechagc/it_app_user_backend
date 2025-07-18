# API Endpoints - User Service

## Base URL
```
http://localhost:8081
```

## 🔧 Sistema

### Health Check
- **GET** `/health` - Verificar estado del servicio
- **GET** `/ping` - Ping simple (responde "pong")

---

## 👤 Usuarios (`/users`)

### Rutas Públicas
- **GET** `/users` - Obtener todos los usuarios
- **GET** `/users/{id}` - Obtener usuario por ID
- **GET** `/users/firebase/{firebase_id}` - Obtener usuario por Firebase ID
- **GET** `/users/username/{username}` - Obtener usuario por username
- **GET** `/users/email/{email}` - Obtener usuario por email
- **GET** `/users/search` - Buscar usuarios
- **GET** `/users/count` - Contar total de usuarios

### Rutas Protegidas (requieren autenticación)
- **POST** `/users/create` - Crear nuevo usuario
- **PUT** `/users/{id}` - Actualizar usuario
- **DELETE** `/users/{id}` - Eliminar usuario
- **POST** `/users/{id}/login` - Actualizar info de login
- **GET** `/users/active` - Obtener usuarios activos
- **GET** `/users/{id}/profile` - Obtener perfil de usuario
- **GET** `/users/{id}/settings` - Obtener configuraciones de usuario
- **GET** `/users/{id}/stats` - Obtener estadísticas de usuario

---

## 🔐 Autenticación (`/auth`)

### Rutas Públicas
- **POST** `/auth/login` - Iniciar sesión
- **POST** `/auth/logout` - Cerrar sesión
- **GET** `/auth/status` - Verificar estado de autenticación
- **POST** `/auth/refresh` - Refrescar token

### Rutas Protegidas
- **GET** `/auth/profile` - Obtener perfil del usuario autenticado
- **PUT** `/auth/profile` - Actualizar perfil del usuario autenticado
- **POST** `/auth/change-password` - Cambiar contraseña
- **POST** `/auth/revoke-tokens` - Revocar todos los tokens
- **GET** `/auth/sessions` - Obtener sesiones activas
- **DELETE** `/auth/sessions/{session_id}` - Revocar sesión específica

---

## 🎫 Tokens (`/tokens`)

### Rutas Públicas
- **POST** `/tokens/verify` - Verificar token
- **POST** `/tokens/refresh` - Refrescar token
- **POST** `/tokens/custom` - Crear token personalizado

### Rutas Protegidas
- **POST** `/tokens/revoke` - Revocar token
- **POST** `/tokens/revoke-all` - Revocar todos los tokens
- **GET** `/tokens/info` - Obtener información del token
- **POST** `/tokens/validate` - Validar token

---

## 🔑 Reset de Contraseña (`/password`)

### Rutas Públicas
- **POST** `/password/reset/request` - Solicitar reset de contraseña
- **POST** `/password/reset/verify` - Verificar código de reset
- **POST** `/password/reset/confirm` - Confirmar reset de contraseña
- **POST** `/password/reset/validate-token` - Validar token de reset

### Rutas Protegidas
- **POST** `/password/change` - Cambiar contraseña
- **POST** `/password/strength-check` - Verificar fortaleza de contraseña
- **GET** `/password/history` - Obtener historial de contraseñas
- **GET** `/password/policy` - Obtener política de contraseñas

---

## 📧 Verificación de Email (`/email`)

### Rutas Públicas
- **POST** `/email/send-verification` - Enviar email de verificación
- **POST** `/email/verify` - Verificar email
- **POST** `/email/verify-code` - Verificar email con código
- **POST** `/email/resend` - Reenviar email de verificación
- **GET** `/email/status` - Verificar estado de verificación

### Rutas Protegidas
- **GET** `/email/my-status` - Obtener mi estado de verificación
- **POST** `/email/update-email` - Actualizar email
- **GET** `/email/verification-history` - Obtener historial de verificaciones
- **GET** `/email/settings` - Obtener configuraciones de email
- **PUT** `/email/settings` - Actualizar configuraciones de email

---

## 🔓 Login y Sesiones (`/login`)

### Rutas Públicas
- **POST** `/login/track` - Trackear login de usuario
- **GET** `/login/history/{user_id}` - Obtener historial de login
- **GET** `/login/attempts/{email}` - Obtener intentos de login
- **POST** `/login/security-check` - Verificación de seguridad

### Rutas Protegidas
- **POST** `/login/update-info/{id}` - Actualizar información de login
- **GET** `/login/my-history` - Obtener mi historial de login
- **GET** `/login/active-sessions` - Obtener sesiones activas
- **DELETE** `/login/terminate-session/{session_id}` - Terminar sesión específica
- **DELETE** `/login/terminate-all-sessions` - Terminar todas las sesiones
- **GET** `/login/suspicious-activity` - Obtener actividad sospechosa

---

## 📝 Ejemplos de Uso

### Crear Usuario
```bash
curl -X POST http://localhost:8081/users/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <firebase-token>" \
  -d '{
    "firebase_id": "firebase_user_123",
    "email": "usuario@ejemplo.com",
    "username": "usuario123",
    "first_name": "Juan",
    "last_name": "Pérez",
    "status": "active"
  }'
```

### Login
```bash
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "id_token": "<firebase-id-token>"
  }'
```

### Verificar Token
```bash
curl -X POST http://localhost:8081/tokens/verify \
  -H "Content-Type: application/json" \
  -d '{
    "id_token": "<firebase-id-token>"
  }'
```

### Solicitar Reset de Contraseña
```bash
curl -X POST http://localhost:8081/password/reset/request \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@ejemplo.com"
  }'
```

### Enviar Verificación de Email
```bash
curl -X POST http://localhost:8081/email/send-verification \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@ejemplo.com"
  }'
```

---

## 🔒 Autenticación

Para endpoints protegidos, incluir el token de Firebase en el header:

```bash
Authorization: Bearer <firebase-id-token>
```

## 📊 Códigos de Respuesta

- **200** - OK
- **201** - Created
- **400** - Bad Request
- **401** - Unauthorized
- **403** - Forbidden
- **404** - Not Found
- **409** - Conflict
- **429** - Too Many Requests
- **500** - Internal Server Error

## 🚀 Rate Limiting

- **Límite por defecto**: 100 requests por segundo
- **Burst**: 200 requests
- **Configurable** via variables de entorno `RATE_LIMIT_RPS` y `RATE_LIMIT_BURST`