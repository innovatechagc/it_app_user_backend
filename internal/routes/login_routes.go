package routes

import (
	"github.com/gorilla/mux"
	"it-app_user/internal/handlers"
	"it-app_user/internal/middleware"
)

// SetupLoginRoutes configura todas las rutas relacionadas con login y sesiones
func SetupLoginRoutes(router *mux.Router, loginHandler *handlers.LoginHandler, authMiddleware *middleware.AuthMiddleware) {
	// Subrouter para login
	loginRouter := router.PathPrefix("/login").Subrouter()
	
	// Rutas públicas de login
	loginRouter.HandleFunc("/track", loginHandler.TrackUserLogin).Methods("POST")
	loginRouter.HandleFunc("/history/{user_id:[0-9]+}", loginHandler.GetLoginHistory).Methods("GET")
	loginRouter.HandleFunc("/attempts/{email}", loginHandler.GetLoginAttempts).Methods("GET")
	loginRouter.HandleFunc("/security-check", loginHandler.SecurityCheck).Methods("POST")
	
	// Rutas que requieren autenticación
	if authMiddleware != nil {
		protectedLoginRouter := loginRouter.PathPrefix("").Subrouter()
		protectedLoginRouter.Use(authMiddleware.RequireAuth)
		
		protectedLoginRouter.HandleFunc("/update-info/{id:[0-9]+}", loginHandler.UpdateLoginInfo).Methods("POST")
		protectedLoginRouter.HandleFunc("/my-history", loginHandler.GetMyLoginHistory).Methods("GET")
		protectedLoginRouter.HandleFunc("/active-sessions", loginHandler.GetActiveSessions).Methods("GET")
		protectedLoginRouter.HandleFunc("/terminate-session/{session_id}", loginHandler.TerminateSession).Methods("DELETE")
		protectedLoginRouter.HandleFunc("/terminate-all-sessions", loginHandler.TerminateAllSessions).Methods("DELETE")
		protectedLoginRouter.HandleFunc("/suspicious-activity", loginHandler.GetSuspiciousActivity).Methods("GET")
	}
}