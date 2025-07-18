package routes

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	
	"it-app_user/internal/handlers"
	"it-app_user/internal/middleware"
	"it-app_user/internal/models"
	"it-app_user/internal/repositories"
	"it-app_user/pkg/firebase"
)

func SetupRoutes(firebaseAuth *firebase.Auth, rateLimitRPS int, rateLimitBurst int) *mux.Router {
	router := mux.NewRouter()
	
	// Crear repositorios
	db := models.GetDB()
	userRepo := repositories.NewUserRepository(db)
	emailRepo := repositories.NewEmailVerificationRepository(db)
	passwordRepo := repositories.NewPasswordResetRepository(db)
	
	// Crear handlers
	userHandler := handlers.NewUserHandler(userRepo)
	authHandler := handlers.NewAuthHandler(firebaseAuth)
	tokenHandler := handlers.NewTokenHandler(firebaseAuth)
	passwordResetHandler := handlers.NewPasswordResetHandler(firebaseAuth, passwordRepo)
	emailHandler := handlers.NewVerifyEmailHandler(firebaseAuth, emailRepo)
	loginHandler := handlers.NewLoginHandler(firebaseAuth)
	
	// Middleware global
	rateLimiter := middleware.NewRateLimiter(
		rate.Every(time.Second/time.Duration(rateLimitRPS)), 
		rateLimitBurst,
	)
	router.Use(middleware.LoggingMiddleware)
	router.Use(rateLimiter.Middleware)
	router.Use(middleware.CORSMiddleware)
	
	// Middleware de autenticación (opcional)
	var authMiddleware *middleware.AuthMiddleware
	if firebaseAuth != nil {
		authMiddleware = middleware.NewAuthMiddleware(firebaseAuth)
	}
	
	// Rutas de salud
	router.HandleFunc("/health", userHandler.HealthCheck).Methods("GET")
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}).Methods("GET")
	
	// Configurar todas las rutas por módulos
	SetupUserRoutes(router, userHandler, authMiddleware)
	SetupAuthRoutes(router, authHandler, authMiddleware)
	SetupTokenRoutes(router, tokenHandler, authMiddleware)
	SetupPasswordResetRoutes(router, passwordResetHandler, authMiddleware)
	SetupEmailVerificationRoutes(router, emailHandler, authMiddleware)
	SetupLoginRoutes(router, loginHandler, authMiddleware)
	
	return router
}