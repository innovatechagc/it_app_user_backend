# 🗄️ Database Documentation - User Service

## Índice
- [🎯 Visión General](#-visión-general)
- [📊 Esquema de Base de Datos](#-esquema-de-base-de-datos)
- [🔧 Configuración](#-configuración)
- [📈 Migraciones](#-migraciones)
- [🔍 Índices y Optimización](#-índices-y-optimización)
- [💾 Backup y Restauración](#-backup-y-restauración)
- [📊 Monitoreo](#-monitoreo)
- [🔧 Troubleshooting](#-troubleshooting)

## 🎯 Visión General

El User Service utiliza **PostgreSQL** como base de datos principal con **GORM** como ORM. La base de datos está diseñada para ser escalable, eficiente y mantener la integridad de los datos.

### Características
- **PostgreSQL 12+**: Base de datos relacional robusta
- **GORM**: ORM para Go con auto-migraciones
- **JSONB**: Soporte para datos semi-estructurados
- **Índices Optimizados**: Para consultas frecuentes
- **Constraints**: Integridad referencial y validaciones

## 📊 Esquema de Base de Datos

### Tabla Principal: `users`

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
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'pending')),
    pk_aut_use_id VARCHAR(128)
);
```

#### Campos Principales

| Campo | Tipo | Descripción | Constraints |
|-------|------|-------------|-------------|
| `id` | SERIAL | ID único autoincremental | PRIMARY KEY |
| `firebase_id` | VARCHAR(128) | ID de Firebase Auth | UNIQUE, NOT NULL |
| `email` | VARCHAR(255) | Email del usuario | UNIQUE, NOT NULL |
| `email_verified` | BOOLEAN | Estado de verificación | DEFAULT FALSE |
| `username` | VARCHAR(50) | Nombre de usuario único | UNIQUE, NOT NULL |
| `first_name` | VARCHAR(100) | Nombre | - |
| `last_name` | VARCHAR(100) | Apellido | - |
| `provider` | VARCHAR(50) | Proveedor de auth | - |
| `provider_id` | VARCHAR(128) | ID del proveedor | - |
| `created_at` | TIMESTAMP | Fecha de creación | DEFAULT NOW() |
| `updated_at` | TIMESTAMP | Última actualización | DEFAULT NOW() |
| `login_count` | INTEGER | Contador de logins | DEFAULT 0 |
| `last_login_at` | TIMESTAMP | Último login | - |
| `last_login_ip` | INET | IP del último login | - |
| `last_login_device` | VARCHAR(255) | Dispositivo del último login | - |
| `disabled` | BOOLEAN | Usuario deshabilitado | DEFAULT FALSE |
| `status` | VARCHAR(20) | Estado del usuario | CHECK constraint |
| `pk_aut_use_id` | VARCHAR(128) | ID de autorización | - |

### Tablas Relacionadas

#### `user_profiles`
```sql
CREATE TABLE user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    avatar VARCHAR(500),
    bio TEXT,
    website VARCHAR(255),
    location VARCHAR(100),
    birthday DATE,
    gender VARCHAR(20),
    phone VARCHAR(20),
    preferences JSONB,
    privacy JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id)
);
```

#### `user_settings`
```sql
CREATE TABLE user_settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'UTC',
    theme VARCHAR(20) DEFAULT 'light',
    notifications JSONB,
    privacy JSONB,
    security JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id)
);
```

#### `user_stats`
```sql
CREATE TABLE user_stats (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    login_count INTEGER DEFAULT 0,
    last_login_at TIMESTAMP WITH TIME ZONE,
    profile_views INTEGER DEFAULT 0,
    account_age_days INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    last_active_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id)
);
```

#### `email_verifications`
```sql
CREATE TABLE email_verifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    firebase_id VARCHAR(128) NOT NULL,
    email VARCHAR(255) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP WITH TIME ZONE,
    verification_code VARCHAR(6),
    code_expires_at TIMESTAMP WITH TIME ZONE,
    attempts_count INTEGER DEFAULT 0,
    last_attempt_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### `password_reset_tokens`
```sql
CREATE TABLE password_reset_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    firebase_id VARCHAR(128) NOT NULL,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) UNIQUE NOT NULL,
    code VARCHAR(6),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## 🔧 Configuración

### Variables de Entorno
```bash
# Base de Datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=itapp
```

### Configuración de Conexión (GORM)
```go
// internal/database/database.go
func ConnectDB() {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    
    // Configurar pool de conexiones
    sqlDB, _ := DB.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
}
```

### Pool de Conexiones
```go
// Configuración optimizada para producción
sqlDB.SetMaxIdleConns(10)        // Conexiones idle
sqlDB.SetMaxOpenConns(100)       // Conexiones máximas
sqlDB.SetConnMaxLifetime(time.Hour) // Tiempo de vida
```

## 📈 Migraciones

### Auto-Migraciones (GORM)
```go
// internal/models/connection.go
func MigrateDB() {
    db := database.GetDB()
    err := db.AutoMigrate(
        &User{},
        &EmailVerification{},
        &PasswordResetToken{},
        &UserProfile{},
        &UserSettings{},
        &UserStats{},
    )
    
    if err != nil {
        log.Fatalf("Error al ejecutar migraciones: %v", err)
    }
}
```

### Migraciones Manuales
```sql
-- migrations/001_create_users_table.sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    firebase_id VARCHAR(128) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    -- ... resto de campos
);

-- migrations/002_add_indexes.sql
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_firebase_id ON users(firebase_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_username ON users(username);
```

### Script de Migración
```bash
#!/bin/bash
# scripts/migrate.sh

echo "Running database migrations..."

# Ejecutar migraciones SQL
for file in migrations/*.sql; do
    echo "Executing $file..."
    psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f "$file"
done

echo "Migrations completed!"
```

## 🔍 Índices y Optimización

### Índices Principales
```sql
-- Índices únicos (automáticos)
CREATE UNIQUE INDEX users_firebase_id_key ON users(firebase_id);
CREATE UNIQUE INDEX users_email_key ON users(email);
CREATE UNIQUE INDEX users_username_key ON users(username);

-- Índices de búsqueda frecuente
CREATE INDEX CONCURRENTLY idx_users_status ON users(status) WHERE status = 'active';
CREATE INDEX CONCURRENTLY idx_users_created_at ON users(created_at);
CREATE INDEX CONCURRENTLY idx_users_last_login ON users(last_login_at);

-- Índices compuestos
CREATE INDEX CONCURRENTLY idx_users_status_created ON users(status, created_at);
CREATE INDEX CONCURRENTLY idx_users_provider_status ON users(provider, status);
CREATE INDEX CONCURRENTLY idx_users_email_verified ON users(email_verified, status);
```

### Índices para Tablas Relacionadas
```sql
-- user_profiles
CREATE INDEX CONCURRENTLY idx_user_profiles_user_id ON user_profiles(user_id);

-- user_settings  
CREATE INDEX CONCURRENTLY idx_user_settings_user_id ON user_settings(user_id);

-- email_verifications
CREATE INDEX CONCURRENTLY idx_email_verifications_user_id ON email_verifications(user_id);
CREATE INDEX CONCURRENTLY idx_email_verifications_email ON email_verifications(email);
CREATE INDEX CONCURRENTLY idx_email_verifications_code ON email_verifications(verification_code);

-- password_reset_tokens
CREATE INDEX CONCURRENTLY idx_password_reset_tokens_token ON password_reset_tokens(token);
CREATE INDEX CONCURRENTLY idx_password_reset_tokens_email ON password_reset_tokens(email);
CREATE INDEX CONCURRENTLY idx_password_reset_tokens_expires ON password_reset_tokens(expires_at);
```

### Análisis de Performance
```sql
-- Verificar uso de índices
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'user@example.com';

-- Estadísticas de tablas
SELECT 
    schemaname,
    tablename,
    n_tup_ins as inserts,
    n_tup_upd as updates,
    n_tup_del as deletes,
    n_live_tup as live_tuples,
    n_dead_tup as dead_tuples
FROM pg_stat_user_tables 
WHERE tablename = 'users';

-- Índices no utilizados
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes 
WHERE idx_scan = 0;
```

## 💾 Backup y Restauración

### Backup Automático
```bash
#!/bin/bash
# scripts/backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"
DB_NAME="itapp"

# Crear backup completo
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME > "$BACKUP_DIR/backup_$DATE.sql"

# Comprimir backup
gzip "$BACKUP_DIR/backup_$DATE.sql"

# Limpiar backups antiguos (mantener últimos 7 días)
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete

echo "Backup completed: backup_$DATE.sql.gz"
```

### Backup Solo de Datos
```bash
# Solo datos (sin esquema)
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME --data-only > data_backup.sql

# Solo esquema (sin datos)
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME --schema-only > schema_backup.sql

# Tabla específica
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME -t users > users_backup.sql
```

### Restauración
```bash
# Restaurar backup completo
psql -h $DB_HOST -U $DB_USER -d $DB_NAME < backup_20240101_120000.sql

# Restaurar con recreación de BD
dropdb -h $DB_HOST -U $DB_USER $DB_NAME
createdb -h $DB_HOST -U $DB_USER $DB_NAME
psql -h $DB_HOST -U $DB_USER -d $DB_NAME < backup_20240101_120000.sql
```

### Cron Job para Backups
```bash
# Agregar a crontab
crontab -e

# Backup diario a las 2 AM
0 2 * * * /path/to/scripts/backup.sh >> /var/log/backup.log 2>&1

# Backup semanal completo los domingos
0 1 * * 0 /path/to/scripts/full_backup.sh >> /var/log/backup.log 2>&1
```

## 📊 Monitoreo

### Métricas Importantes
```sql
-- Tamaño de la base de datos
SELECT 
    pg_database.datname,
    pg_size_pretty(pg_database_size(pg_database.datname)) AS size
FROM pg_database
WHERE datname = 'itapp';

-- Tamaño de tablas
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size(tablename::regclass)) AS size,
    pg_size_pretty(pg_relation_size(tablename::regclass)) AS table_size,
    pg_size_pretty(pg_total_relation_size(tablename::regclass) - pg_relation_size(tablename::regclass)) AS index_size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(tablename::regclass) DESC;

-- Conexiones activas
SELECT 
    count(*) as total_connections,
    count(*) FILTER (WHERE state = 'active') as active_connections,
    count(*) FILTER (WHERE state = 'idle') as idle_connections
FROM pg_stat_activity 
WHERE datname = 'itapp';

-- Queries lentas
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    rows
FROM pg_stat_statements 
WHERE query NOT LIKE '%pg_stat_statements%'
ORDER BY mean_time DESC 
LIMIT 10;
```

### Alertas y Monitoreo
```sql
-- Crear función para alertas
CREATE OR REPLACE FUNCTION check_db_health()
RETURNS TABLE(
    metric VARCHAR,
    value NUMERIC,
    status VARCHAR
) AS $$
BEGIN
    -- Verificar conexiones
    RETURN QUERY
    SELECT 
        'active_connections'::VARCHAR,
        count(*)::NUMERIC,
        CASE 
            WHEN count(*) > 80 THEN 'CRITICAL'
            WHEN count(*) > 60 THEN 'WARNING'
            ELSE 'OK'
        END::VARCHAR
    FROM pg_stat_activity 
    WHERE datname = 'itapp' AND state = 'active';
    
    -- Verificar tamaño de BD
    RETURN QUERY
    SELECT 
        'database_size_mb'::VARCHAR,
        (pg_database_size('itapp') / 1024 / 1024)::NUMERIC,
        CASE 
            WHEN pg_database_size('itapp') > 10737418240 THEN 'WARNING' -- 10GB
            ELSE 'OK'
        END::VARCHAR;
END;
$$ LANGUAGE plpgsql;

-- Ejecutar verificación
SELECT * FROM check_db_health();
```

### Dashboard de Métricas
```sql
-- Vista para dashboard
CREATE VIEW db_dashboard AS
SELECT 
    'Users Total' as metric,
    count(*)::text as value
FROM users
UNION ALL
SELECT 
    'Active Users',
    count(*)::text
FROM users 
WHERE status = 'active' AND disabled = false
UNION ALL
SELECT 
    'New Users Today',
    count(*)::text
FROM users 
WHERE created_at >= CURRENT_DATE
UNION ALL
SELECT 
    'Database Size',
    pg_size_pretty(pg_database_size('itapp'))
UNION ALL
SELECT 
    'Active Connections',
    count(*)::text
FROM pg_stat_activity 
WHERE datname = 'itapp' AND state = 'active';

-- Consultar dashboard
SELECT * FROM db_dashboard;
```

## 🔧 Troubleshooting

### Problemas Comunes

#### 1. Conexión Rechazada
```bash
# Verificar que PostgreSQL esté corriendo
sudo systemctl status postgresql

# Verificar puerto
netstat -tlnp | grep 5432

# Verificar configuración
sudo -u postgres psql -c "SHOW listen_addresses;"
```

#### 2. Demasiadas Conexiones
```sql
-- Ver conexiones actuales
SELECT count(*) FROM pg_stat_activity WHERE datname = 'itapp';

-- Ver límite de conexiones
SHOW max_connections;

-- Terminar conexiones idle
SELECT pg_terminate_backend(pid) 
FROM pg_stat_activity 
WHERE datname = 'itapp' 
  AND state = 'idle' 
  AND state_change < now() - interval '1 hour';
```

#### 3. Queries Lentas
```sql
-- Habilitar log de queries lentas
ALTER SYSTEM SET log_min_duration_statement = 1000; -- 1 segundo
SELECT pg_reload_conf();

-- Ver queries activas
SELECT 
    pid,
    now() - pg_stat_activity.query_start AS duration,
    query 
FROM pg_stat_activity 
WHERE (now() - pg_stat_activity.query_start) > interval '5 minutes';
```

#### 4. Espacio en Disco
```sql
-- Verificar espacio usado
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size(tablename::regclass)) AS size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(tablename::regclass) DESC;

-- Limpiar datos antiguos
DELETE FROM password_reset_tokens WHERE expires_at < NOW() - INTERVAL '1 day';
DELETE FROM email_verifications WHERE created_at < NOW() - INTERVAL '30 days' AND is_verified = true;

-- VACUUM para recuperar espacio
VACUUM ANALYZE;
```

### Mantenimiento Regular

#### Script de Mantenimiento
```bash
#!/bin/bash
# scripts/maintenance.sh

echo "Starting database maintenance..."

# Limpiar tokens expirados
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "
DELETE FROM password_reset_tokens WHERE expires_at < NOW();
DELETE FROM email_verifications WHERE code_expires_at < NOW() AND is_verified = false;
"

# Actualizar estadísticas
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "ANALYZE;"

# VACUUM ligero
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "VACUUM;"

echo "Maintenance completed!"
```

#### Cron Job de Mantenimiento
```bash
# Mantenimiento diario a las 3 AM
0 3 * * * /path/to/scripts/maintenance.sh >> /var/log/db_maintenance.log 2>&1

# VACUUM FULL semanal (más intensivo)
0 4 * * 0 psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "VACUUM FULL;" >> /var/log/db_maintenance.log 2>&1
```

---

## 📚 Recursos Adicionales

- [📖 PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [🔧 GORM Documentation](https://gorm.io/docs/)
- [📊 PostgreSQL Performance](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [🏠 Volver al README](../README.md)

---

**¿Problemas con la base de datos?** Revisa los logs, verifica las conexiones y consulta las métricas de performance.