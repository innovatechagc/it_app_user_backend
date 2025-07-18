# 🔥 Configuración de Firebase

## Pasos para configurar Firebase en el proyecto

### 1. **Obtener el archivo Service Account**

1. Ve a [Firebase Console](https://console.firebase.google.com/)
2. Selecciona tu proyecto: **innovatech-app**
3. Ve a **Configuración del proyecto** (ícono de engranaje)
4. Ve a la pestaña **Cuentas de servicio**
5. Haz clic en **Generar nueva clave privada**
6. Descarga el archivo JSON
7. Renómbralo a `firebase-service-account.json`
8. Colócalo en la raíz del proyecto

### 2. **Estructura del archivo Service Account**

El archivo debe tener esta estructura:

```json
{
  "type": "service_account",
  "project_id": "innovatech-app",
  "private_key_id": "...",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
  "client_email": "firebase-adminsdk-xxxxx@innovatech-app.iam.gserviceaccount.com",
  "client_id": "...",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "..."
}
```

### 3. **Configuración del cliente (JavaScript)**

Tu configuración del cliente ya está lista:

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

### 4. **Variables de entorno**

Asegúrate de que tu archivo `.env` tenga:

```bash
FIREBASE_PROJECT_ID=innovatech-app
```

### 5. **Probar la configuración**

Una vez que tengas el archivo `firebase-service-account.json`:

```bash
# Ejecutar el servicio
go run cmd/main.go

# O con Docker
docker-compose up --build
```

### 6. **Endpoints de Firebase disponibles**

- `POST /auth/login` - Login con Firebase ID Token
- `POST /tokens/verify` - Verificar Firebase ID Token
- `GET /auth/status` - Estado de autenticación
- `POST /auth/revoke-tokens` - Revocar tokens de Firebase

### 7. **Ejemplo de uso**

```bash
# Login con Firebase token
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "id_token": "tu-firebase-id-token-aqui"
  }'
```

## ⚠️ Importante

- **NUNCA** subas el archivo `firebase-service-account.json` a Git
- El archivo ya está en `.gitignore`
- Usa variables de entorno en producción