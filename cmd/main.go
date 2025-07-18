package main

import (
	"os"

	"github.com/joho/godotenv"
	"it-app_user/internal/config"
	"it-app_user/internal/logger"
	"it-app_user/internal/server"
)

func main() {
	// Cargar archivo .env si existe
	if err := godotenv.Load(); err != nil {
		// No es un error fatal si no existe el archivo .env
		// Las variables de entorno del sistema tendrán precedencia
	}

	// Inicializar logger
	logger.Init()
	log := logger.GetLogger()

	// Cargar configuración
	cfg := config.LoadConfig()
	
	log.WithField("environment", cfg.Environment).WithField("port", cfg.Port).Info("Starting User Service")

	// Crear servidor
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.WithError(err).Fatal("Error creating server")
		os.Exit(1)
	}

	// Iniciar servidor
	log.WithField("port", cfg.Port).Info("Server initialized successfully")
	if err := srv.Start(); err != nil {
		log.WithError(err).Fatal("Error starting server")
		os.Exit(1)
	}
}
