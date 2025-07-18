package database

import (
	"fmt"
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB establece la conexión con la base de datos PostgreSQL
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
	
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	// Configurar pool de conexiones
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error al obtener la instancia de base de datos: %v", err)
	}

	// Configuraciones del pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Conexión a base de datos establecida exitosamente")
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}