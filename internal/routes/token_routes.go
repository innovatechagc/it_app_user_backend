package routes

import (
	"github.com/gorilla/mux"
	"it-app_user/internal/handlers"
	"it-app_user/internal/middleware"
)

// SetupTokenRoutes configura todas las rutas relacionadas con tokens
func SetupTokenRoutes(router *mux.Router, tokenHandler *handlers.TokenHandler, authMiddleware *middleware.AuthMiddleware) {
	// Subrouter para tokens
	tokenRouter := router.PathPrefix("/tokens").Subrouter()
	
	// Rutas públicas de tokens
	tokenRouter.HandleFunc("/verify", tokenHandler.VerifyToken).Methods("POST")
	tokenRouter.HandleFunc("/refresh", tokenHandler.RefreshToken).Methods("POST")
	tokenRouter.HandleFunc("/custom", tokenHandler.CreateCustomToken).Methods("POST")
	
	// Rutas que requieren autenticación
	if authMiddleware != nil {
		protectedTokenRouter := tokenRouter.PathPrefix("").Subrouter()
		protectedTokenRouter.Use(authMiddleware.RequireAuth)
		
		protectedTokenRouter.HandleFunc("/revoke", tokenHandler.RevokeToken).Methods("POST")
		protectedTokenRouter.HandleFunc("/revoke-all", tokenHandler.RevokeAllTokens).Methods("POST")
		protectedTokenRouter.HandleFunc("/info", tokenHandler.GetTokenInfo).Methods("GET")
		protectedTokenRouter.HandleFunc("/validate", tokenHandler.ValidateToken).Methods("POST")
	}
}