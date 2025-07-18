# 🔥 Firebase Setup - Configuración Completa

## Índice
- [🎯 Visión General](#-visión-general)
- [🚀 Setup Inicial](#-setup-inicial)
- [🔧 Configuración del Servidor](#-configuración-del-servidor)
- [💻 Configuración del Cliente](#-configuración-del-cliente)
- [🔐 Autenticación](#-autenticación)
- [🧪 Testing](#-testing)
- [🚀 Producción](#-producción)
- [🔧 Troubleshooting](#-troubleshooting)

## 🎯 Visión General

Firebase Authentication proporciona autenticación robusta y escalable para el User Service. Esta guía cubre la configuración completa tanto del lado del servidor (Go) como del cliente (JavaScript).

### Tu Configuración Actual
```javascript
const firebaseConfig = {
  apiKey: "AIzaSyCbAXjfv-f0n7b91CrQ6nkn2Pt1TNSunFw",
  authDomain: "innovatech-app.firebaseapp.com",
  projectId: "innovatech-app",
  storageBucket: "innovatech-app.firebasestorage.app",
  messagingSenderId: "68143034843",
  appId: "1:68143034843:web:56d4938a0e629d76ae77fe",
  measurementId: "G-MRSYLS5KWK"
};
```

## 🚀 Setup Inicial

### 1. Obtener Service Account de Firebase

#### Paso a Paso:
1. Ve a [Firebase Console](https://console.firebase.google.com/)
2. Selecciona tu proyecto: **innovatech-app**
3. Ve a **Configuración del proyecto** (⚙️)
4. Pestaña **"Cuentas de servicio"**
5. Haz clic en **"Generar nueva clave privada"**
6. Descarga el archivo JSON
7. Renómbralo a `firebase-service-account.json`
8. Colócalo en la raíz del proyecto

### 2. Estructura del Service Account
```json
{
  "type": "service_account",
  "project_id": "innovatech-app",
  "private_key_id": "abc123...",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
  "client_email": "firebase-adminsdk-xyz@innovatech-app.iam.gserviceaccount.com",
  "client_id": "123456789",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-xyz%40innovatech-app.iam.gserviceaccount.com"
}
```

## 🔧 Configuración del Servidor

### Variables de Entorno (Ya configuradas)
```bash
# .env
FIREBASE_PROJECT_ID=innovatech-app
```

### Verificar Configuración
```bash
# 1. Verificar archivo service account
ls -la firebase-service-account.json

# 2. Verificar variables de entorno
echo $FIREBASE_PROJECT_ID

# 3. Ejecutar el servicio
go run cmd/main.go
```

**Output esperado:**
```
Firebase Auth initialized successfully for project: innovatech-app
Server starting on port 8081
```

## 💻 Configuración del Cliente

### Implementación Completa
```javascript
// firebase-config.js
import { initializeApp } from 'firebase/app';
import { getAuth } from 'firebase/auth';
import { getAnalytics } from 'firebase/analytics';

const firebaseConfig = {
  apiKey: "AIzaSyCbAXjfv-f0n7b91CrQ6nkn2Pt1TNSunFw",
  authDomain: "innovatech-app.firebaseapp.com",
  projectId: "innovatech-app",
  storageBucket: "innovatech-app.firebasestorage.app",
  messagingSenderId: "68143034843",
  appId: "1:68143034843:web:56d4938a0e629d76ae77fe",
  measurementId: "G-MRSYLS5KWK"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
const auth = getAuth(app);
const analytics = getAnalytics(app);

export { auth, analytics };
```

### Servicio de Autenticación
```javascript
// auth-service.js
import { 
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
  signOut,
  onAuthStateChanged
} from 'firebase/auth';
import { auth } from './firebase-config.js';

class AuthService {
  // Login y comunicación con tu backend
  async login(email, password) {
    try {
      const userCredential = await signInWithEmailAndPassword(auth, email, password);
      const idToken = await userCredential.user.getIdToken();
      
      // Enviar token a tu User Service
      const response = await fetch('http://localhost:8081/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id_token: idToken })
      });
      
      return await response.json();
    } catch (error) {
      throw new Error(`Login failed: ${error.message}`);
    }
  }

  // Obtener token para requests autenticadas
  async getCurrentToken() {
    const user = auth.currentUser;
    return user ? await user.getIdToken() : null;
  }
}

export default new AuthService();
```

## 🔐 Endpoints Disponibles

### Autenticación
```bash
# Login con Firebase token
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"id_token": "tu-firebase-token"}'

# Verificar token
curl -X POST http://localhost:8081/tokens/verify \
  -H "Content-Type: application/json" \
  -d '{"id_token": "tu-firebase-token"}'

# Estado de autenticación
curl -X GET http://localhost:8081/auth/status \
  -H "Authorization: Bearer tu-firebase-token"
```

### Usuarios Protegidos
```bash
# Crear usuario (requiere auth)
curl -X POST http://localhost:8081/users/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer tu-firebase-token" \
  -d '{
    "firebase_id": "firebase_user_123",
    "email": "usuario@ejemplo.com",
    "username": "usuario123"
  }'
```

## 🧪 Testing

### Test Rápido
```bash
# 1. Verificar que el servicio está corriendo
curl http://localhost:8081/health

# 2. Probar endpoint público
curl http://localhost:8081/users

# 3. Con Firebase token (obtener desde cliente)
curl -X GET http://localhost:8081/auth/profile \
  -H "Authorization: Bearer tu-token-aqui"
```

### Flujo Completo de Testing
```javascript
// En tu aplicación cliente
import AuthService from './auth-service.js';

// 1. Login
const result = await AuthService.login('test@example.com', 'password');
console.log('Login result:', result);

// 2. Obtener token
const token = await AuthService.getCurrentToken();
console.log('Current token:', token);

// 3. Hacer request autenticada
const response = await fetch('http://localhost:8081/users/create', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    firebase_id: 'test_user_123',
    email: 'test@example.com',
    username: 'testuser'
  })
});
```

## 🚀 Producción

### Seguridad
```bash
# Variables de entorno seguras
FIREBASE_PROJECT_ID=innovatech-app
ENVIRONMENT=production

# Service account como secreto
kubectl create secret generic firebase-secret \
  --from-file=firebase-service-account.json
```

### CORS Configurado
El servidor ya está configurado para manejar CORS apropiadamente para tu dominio.

## 🔧 Troubleshooting

### Problemas Comunes

#### 1. "Firebase service account file not found"
```bash
# Solución: Descargar el archivo desde Firebase Console
# y colocarlo como firebase-service-account.json en la raíz
```

#### 2. "Invalid token"
```bash
# Los tokens expiran en 1 hora
# Obtener un nuevo token desde el cliente
```

#### 3. CORS errors
```bash
# Verificar que tu dominio está en Firebase Console
# Authentication → Settings → Authorized domains
```

### Verificación Rápida
```bash
# Verificar configuración
echo "Project ID: $FIREBASE_PROJECT_ID"
ls -la firebase-service-account.json
go run cmd/main.go
```

---

## 📋 Checklist de Configuración

- [ ] ✅ Proyecto Firebase creado: **innovatech-app**
- [ ] ⏳ Service account descargado y colocado
- [ ] ✅ Variables de entorno configuradas
- [ ] ✅ Cliente JavaScript configurado
- [ ] ⏳ Authentication habilitado en Firebase Console
- [ ] ⏳ Dominios autorizados configurados

## 🔗 Enlaces Útiles

- [🏠 README Principal](../README.md)
- [📖 Guía Completa](GUIDE.md)
- [🔧 API Documentation](API.md)
- [🏗️ Arquitectura](ARCHITECTURE.md)

---

**¡Solo necesitas descargar el service account de Firebase Console y ya estará todo listo!** 🚀