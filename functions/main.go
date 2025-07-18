package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/gorilla/mux"
	
	"functions/internal/config"
	"functions/internal/database"
	"functions/internal/handlers"
	"functions/internal/middleware"
	"functions/internal/routes"
	"functions/pkg/firebase"
)

func init() {
	// Registrar la función HTTP
	funcframework.RegisterHTTPFunction("/", handleHTTP)
}

func main() {
	// Usar PORT del entorno o 8080 por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Iniciar el servidor del framework de funciones
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS para todas las respuestas
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	
	// Manejar preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Inicializar base de datos
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Printf("Error initializing database: %v", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	// Inicializar Firebase Auth
	firebaseAuth, err := firebase.NewAuth()
	if err != nil {
		log.Printf("Error initializing Firebase Auth: %v", err)
		http.Error(w, "Firebase initialization error", http.StatusInternalServerError)
		return
	}

	// Crear router
	router := mux.NewRouter()
	
	// Aplicar middlewares globales
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.RateLimitMiddleware)

	// Configurar rutas
	routes.SetupRoutes(router, db, firebaseAuth, cfg)

	// Servir la request
	router.ServeHTTP(w, r)
}