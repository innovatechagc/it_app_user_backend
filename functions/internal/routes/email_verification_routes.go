package routes

import (
	"github.com/gorilla/mux"
	"it-app_user/internal/handlers"
	"it-app_user/internal/middleware"
)

// SetupEmailVerificationRoutes configura todas las rutas relacionadas con verificación de email
func SetupEmailVerificationRoutes(router *mux.Router, emailHandler *handlers.VerifyEmailHandler, authMiddleware *middleware.AuthMiddleware) {
	// Subrouter para verificación de email
	emailRouter := router.PathPrefix("/email").Subrouter()
	
	// Rutas públicas de verificación de email
	emailRouter.HandleFunc("/send-verification", emailHandler.SendVerificationEmail).Methods("POST")
	emailRouter.HandleFunc("/verify", emailHandler.VerifyEmail).Methods("POST")
	emailRouter.HandleFunc("/verify-code", emailHandler.VerifyEmailWithCode).Methods("POST")
	emailRouter.HandleFunc("/resend", emailHandler.ResendVerificationEmail).Methods("POST")
	emailRouter.HandleFunc("/status", emailHandler.CheckEmailVerificationStatus).Methods("GET")
	
	// Rutas que requieren autenticación
	if authMiddleware != nil {
		protectedEmailRouter := emailRouter.PathPrefix("").Subrouter()
		protectedEmailRouter.Use(authMiddleware.RequireAuth)
		
		protectedEmailRouter.HandleFunc("/my-status", emailHandler.GetMyVerificationStatus).Methods("GET")
		protectedEmailRouter.HandleFunc("/update-email", emailHandler.UpdateEmail).Methods("POST")
		protectedEmailRouter.HandleFunc("/verification-history", emailHandler.GetVerificationHistory).Methods("GET")
		protectedEmailRouter.HandleFunc("/settings", emailHandler.GetEmailSettings).Methods("GET")
		protectedEmailRouter.HandleFunc("/settings", emailHandler.UpdateEmailSettings).Methods("PUT")
	}
}