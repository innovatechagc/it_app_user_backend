package routes

import (
	"github.com/gorilla/mux"
	"it-app_user/internal/handlers"
	"it-app_user/internal/middleware"
)

// SetupPasswordResetRoutes configura todas las rutas relacionadas con reset de contraseña
func SetupPasswordResetRoutes(router *mux.Router, passwordResetHandler *handlers.PasswordResetHandler, authMiddleware *middleware.AuthMiddleware) {
	// Subrouter para reset de contraseña
	passwordRouter := router.PathPrefix("/password").Subrouter()
	
	// Rutas públicas de reset de contraseña
	passwordRouter.HandleFunc("/reset/request", passwordResetHandler.RequestPasswordReset).Methods("POST")
	passwordRouter.HandleFunc("/reset/verify", passwordResetHandler.VerifyResetCode).Methods("POST")
	passwordRouter.HandleFunc("/reset/confirm", passwordResetHandler.ConfirmPasswordReset).Methods("POST")
	passwordRouter.HandleFunc("/reset/validate-token", passwordResetHandler.ValidateResetToken).Methods("POST")
	
	// Rutas que requieren autenticación
	if authMiddleware != nil {
		protectedPasswordRouter := passwordRouter.PathPrefix("").Subrouter()
		protectedPasswordRouter.Use(authMiddleware.RequireAuth)
		
		protectedPasswordRouter.HandleFunc("/change", passwordResetHandler.ChangePassword).Methods("POST")
		protectedPasswordRouter.HandleFunc("/strength-check", passwordResetHandler.CheckPasswordStrength).Methods("POST")
		protectedPasswordRouter.HandleFunc("/history", passwordResetHandler.GetPasswordHistory).Methods("GET")
		protectedPasswordRouter.HandleFunc("/policy", passwordResetHandler.GetPasswordPolicy).Methods("GET")
	}
}