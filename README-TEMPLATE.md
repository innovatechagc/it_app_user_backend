# 🚀 Go Firebase Functions Template

Template estándar para desplegar aplicaciones Go como Firebase Functions con GitHub Actions.

## 📁 Estructura del Proyecto

```
├── functions/                 # Código de Firebase Functions
│   ├── internal/             # Lógica interna de la aplicación
│   │   ├── config/          # Configuración
│   │   ├── database/        # Conexión a base de datos
│   │   ├── handlers/        # Controladores HTTP
│   │   ├── middleware/      # Middlewares
│   │   ├── models/          # Modelos de datos
│   │   ├── repositories/    # Capa de datos
│   │   ├── routes/          # Definición de rutas
│   │   └── validator/       # Validaciones
│   ├── pkg/                 # Paquetes reutilizables
│   │   └── firebase/        # Cliente Firebase
│   ├── .env.example         # Variables de entorno ejemplo
│   ├── go.mod              # Dependencias Go
│   └── main.go             # Punto de entrada
├── .github/workflows/       # GitHub Actions
│   └── deploy-functions.yml # Pipeline de despliegue
├── firebase.json           # Configuración Firebase
├── .firebaserc            # Proyecto Firebase
└── README.md              # Documentación
```

## 🛠️ Configuración Inicial

### 1. Firebase Setup

```bash
# Instalar Firebase CLI
npm install -g firebase-tools

# Login a Firebase
firebase login

# Inicializar proyecto (si es nuevo)
firebase init functions

# Configurar proyecto
firebase use --add your-project-id
```

### 2. Variables de Entorno

Copia `.env.example` a `.env` en la carpeta `functions/`:

```bash
cp functions/.env.example functions/.env
```

Configura las variables necesarias:
- `FIREBASE_PROJECT_ID`: ID de tu proyecto Firebase
- `DB_*`: Configuración de base de datos
- Otras variables según tu aplicación

### 3. GitHub Secrets

Configura estos secrets en tu repositorio GitHub:

- `FIREBASE_TOKEN`: Token de Firebase CLI
  ```bash
  firebase login:ci
  ```

## 🚀 Despliegue

### Automático (GitHub Actions)
El despliegue se ejecuta automáticamente al hacer push a `main`.

### Manual
```bash
# Desde la raíz del proyecto
firebase deploy --only functions
```

## 🧪 Desarrollo Local

```bash
# Instalar dependencias
cd functions
go mod tidy

# Ejecutar localmente
go run main.go

# Con emulador Firebase
firebase emulators:start --only functions
```

## 📝 Uso como Template

1. **Clona este repositorio**
2. **Actualiza firebase.json** con tu configuración
3. **Modifica .firebaserc** con tu project ID
4. **Configura GitHub Secrets**
5. **Personaliza el código** en `functions/`
6. **Push a main** para desplegar

## 🔧 Personalización

### Agregar nuevos endpoints
1. Crear handler en `functions/internal/handlers/`
2. Agregar ruta en `functions/internal/routes/`
3. Actualizar documentación

### Configurar base de datos
1. Actualizar `functions/internal/database/`
2. Configurar variables de entorno
3. Agregar modelos en `functions/internal/models/`

### Middlewares personalizados
1. Crear en `functions/internal/middleware/`
2. Aplicar en `functions/main.go`

## 📚 Documentación Adicional

- [Firebase Functions Go](https://firebase.google.com/docs/functions/go)
- [Functions Framework Go](https://github.com/GoogleCloudPlatform/functions-framework-go)
- [GitHub Actions](https://docs.github.com/en/actions)

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una rama feature
3. Commit tus cambios
4. Push a la rama
5. Abre un Pull Request