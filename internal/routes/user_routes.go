package routes

import (
	"github.com/gorilla/mux"
	"it-app_user/internal/handlers"
	"it-app_user/internal/middleware"
)

// SetupUserRoutes configura todas las rutas relacionadas con usuarios
func SetupUserRoutes(router *mux.Router, userHandler *handlers.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	// Subrouter para usuarios
	userRouter := router.PathPrefix("/users").Subrouter()
	
	// Rutas públicas de usuarios
	userRouter.HandleFunc("", userHandler.GetAllUsers).Methods("GET")
	userRouter.HandleFunc("/{id:[0-9]+}", userHandler.GetUserByID).Methods("GET")
	userRouter.HandleFunc("/firebase/{firebase_id}", userHandler.GetUserByFirebaseID).Methods("GET")
	userRouter.HandleFunc("/username/{username}", userHandler.GetUserByUsername).Methods("GET")
	userRouter.HandleFunc("/email/{email}", userHandler.GetUserByEmail).Methods("GET")
	userRouter.HandleFunc("/search", userHandler.SearchUsers).Methods("GET")
	userRouter.HandleFunc("/count", userHandler.CountUsers).Methods("GET")
	
	// Rutas que requieren autenticación
	if authMiddleware != nil {
		protectedUserRouter := userRouter.PathPrefix("").Subrouter()
		protectedUserRouter.Use(authMiddleware.RequireAuth)
		
		// CRUD protegido
		protectedUserRouter.HandleFunc("/create", userHandler.CreateUser).Methods("POST")
		protectedUserRouter.HandleFunc("/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
		protectedUserRouter.HandleFunc("/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")
		
		// Operaciones específicas
		protectedUserRouter.HandleFunc("/{id:[0-9]+}/login", userHandler.UpdateLoginInfo).Methods("POST")
		protectedUserRouter.HandleFunc("/active", userHandler.GetActiveUsers).Methods("GET")
		protectedUserRouter.HandleFunc("/{id:[0-9]+}/profile", userHandler.GetUserProfile).Methods("GET")
		protectedUserRouter.HandleFunc("/{id:[0-9]+}/settings", userHandler.GetUserSettings).Methods("GET")
		protectedUserRouter.HandleFunc("/{id:[0-9]+}/stats", userHandler.GetUserStats).Methods("GET")
	} else {
		// Si no hay autenticación, todas las rutas son públicas (desarrollo)
		userRouter.HandleFunc("/create", userHandler.CreateUser).Methods("POST")
		userRouter.HandleFunc("/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
		userRouter.HandleFunc("/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")
		userRouter.HandleFunc("/{id:[0-9]+}/login", userHandler.UpdateLoginInfo).Methods("POST")
		userRouter.HandleFunc("/active", userHandler.GetActiveUsers).Methods("GET")
	}
}