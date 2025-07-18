package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"it-app_user/internal/config"
	"it-app_user/internal/logger"
	"it-app_user/internal/models"
	"it-app_user/internal/routes"
	"it-app_user/pkg/firebase"
)

type Server struct {
	config       config.Config
	router       *mux.Router
	firebaseAuth *firebase.Auth
}

func NewServer(cfg config.Config) (*Server, error) {
	// Inicializar logger
	logger.Init()
	log := logger.GetLogger()

	// Conectar a la base de datos
	models.ConnectDB()
	
	// Ejecutar migraciones
	models.MigrateDB()

	// Inicializar Firebase Auth (opcional)
	var firebaseAuth *firebase.Auth
	var err error
	if cfg.FirebaseProjectID != "" {
		firebaseAuth, err = firebase.NewAuth("firebase-service-account.json")
		if err != nil {
			log.WithError(err).Warn("Failed to initialize Firebase Auth, continuing without it")
		}
	}

	// Crear servidor
	server := &Server{
		config:       cfg,
		firebaseAuth: firebaseAuth,
	}

	// Configurar rutas
	server.setupRoutes()

	return server, nil
}

func (s *Server) setupRoutes() {
	// Usar el router de routes.go
	s.router = routes.SetupRoutes(s.firebaseAuth, s.config.RateLimitRPS, s.config.RateLimitBurst)
}

func (s *Server) Start() error {
	log := logger.GetLogger()
	
	server := &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.WithField("port", s.config.Port).Info("Server starting")
	return server.ListenAndServe()
}