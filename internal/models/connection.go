package models

import (
	"log"
	"gorm.io/gorm"
	"it-app_user/internal/database"
)

// ConnectDB establece la conexión con la base de datos PostgreSQL
func ConnectDB() {
	database.ConnectDB()
}

// MigrateDB ejecuta las migraciones automáticas
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
	
	log.Println("Migraciones ejecutadas exitosamente")
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return database.GetDB()
}