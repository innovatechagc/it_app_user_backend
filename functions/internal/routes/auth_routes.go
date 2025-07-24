package routes

import (
	"github.com/gorilla/mux"
	"it-app_user/internal/handlers"
	"it-app_user/internal/middleware"
)

// SetupAuthRoutes configura todas las rutas relacionadas con autenticación
func SetupAuthRoutes(router *mux.Router, authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	// Subrouter para autenticación
	authRouter := router.PathPrefix("/auth").Subrouter()
	
	// Rutas públicas de autenticación
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRouter.HandleFunc("/google", authHandler.GoogleLogin).Methods("POST")
	authRouter.HandleFunc("/facebook", authHandler.FacebookLogin).Methods("POST")
	authRouter.HandleFunc("/email", authHandler.EmailPasswordLogin).Methods("POST")
	authRouter.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	authRouter.HandleFunc("/status", authHandler.CheckAuthStatus).Methods("GET")
	authRouter.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
	
	// Rutas que requieren autenticación
	if authMiddleware != nil {
		protectedAuthRouter := authRouter.PathPrefix("").Subrouter()
		protectedAuthRouter.Use(authMiddleware.RequireAuth)
		
		protectedAuthRouter.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")
		protectedAuthRouter.HandleFunc("/profile", authHandler.UpdateProfile).Methods("PUT")
		protectedAuthRouter.HandleFunc("/change-password", authHandler.ChangePassword).Methods("POST")
		protectedAuthRouter.HandleFunc("/revoke-tokens", authHandler.RevokeAllTokens).Methods("POST")
		protectedAuthRouter.HandleFunc("/sessions", authHandler.GetActiveSessions).Methods("GET")
		protectedAuthRouter.HandleFunc("/sessions/{session_id}", authHandler.RevokeSession).Methods("DELETE")
	}
}