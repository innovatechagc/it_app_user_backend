package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"it-app_user/internal/logger"
	"it-app_user/internal/models"
	"it-app_user/internal/validator"
	"it-app_user/pkg/firebase"
)

type AuthHandler struct {
	firebaseAuth *firebase.Auth
}

func NewAuthHandler(firebaseAuth *firebase.Auth) *AuthHandler {
	return &AuthHandler{
		firebaseAuth: firebaseAuth,
	}
}

// Login maneja POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.LoginRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read request body")
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.WithError(err).Error("Failed to unmarshal JSON")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validar estructura
	if err := validator.ValidateStruct(&req); err != nil {
		log.WithError(err).Warn("Validation failed for login request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// Verificar token de Firebase
	token, err := h.firebaseAuth.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		log.WithError(err).Warn("Invalid Firebase token")
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Obtener información del usuario de Firebase
	userRecord, err := h.firebaseAuth.GetUser(context.Background(), token.UID)
	if err != nil {
		log.WithError(err).Error("Failed to get user from Firebase")
		http.Error(w, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	// Determinar el proveedor basado en los datos del usuario
	provider := h.determineProvider(userRecord, req.Provider)
	
	// Extraer nombres del display name si están disponibles
	firstName, lastName := h.extractNames(userRecord.DisplayName)

	// Crear respuesta de login
	response := models.LoginResponse{
		User: &models.User{
			FirebaseID:    token.UID,
			Email:         userRecord.Email,
			EmailVerified: userRecord.EmailVerified,
			Username:      h.generateUsername(userRecord),
			FirstName:     firstName,
			LastName:      lastName,
			PhotoURL:      userRecord.PhotoURL,
			Provider:      provider,
			Status:        "active",
		},
		Message: "Login successful",
	}

	log.WithFields(map[string]interface{}{
		"firebase_id": token.UID,
		"provider":    provider,
		"email":       userRecord.Email,
	}).Info("User logged in successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper methods
func (h *AuthHandler) determineProvider(userRecord *auth.UserRecord, requestProvider string) string {
	// Si hay información de proveedores en Firebase, usar esa
	if len(userRecord.ProviderUserInfo) > 0 {
		for _, provider := range userRecord.ProviderUserInfo {
			if provider.ProviderID == "google.com" {
				return "google.com"
			}
			if provider.ProviderID == "facebook.com" {
				return "facebook.com"
			}
		}
	}
	
	// Si no, usar el proveedor del request o por defecto password
	if requestProvider != "" {
		return requestProvider
	}
	
	return "password"
}

func (h *AuthHandler) extractNames(displayName string) (string, string) {
	if displayName == "" {
		return "", ""
	}
	
	parts := strings.Fields(displayName)
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	
	firstName := parts[0]
	lastName := strings.Join(parts[1:], " ")
	return firstName, lastName
}

func (h *AuthHandler) generateUsername(userRecord *auth.UserRecord) string {
	if userRecord.DisplayName != "" {
		// Limpiar el display name para usarlo como username
		username := strings.ToLower(strings.ReplaceAll(userRecord.DisplayName, " ", ""))
		// Remover caracteres especiales básicos
		username = strings.ReplaceAll(username, ".", "")
		username = strings.ReplaceAll(username, "-", "")
		if len(username) >= 3 && len(username) <= 50 {
			return username
		}
	}
	
	// Si no hay display name o no es válido, usar parte del email
	if userRecord.Email != "" {
		emailParts := strings.Split(userRecord.Email, "@")
		if len(emailParts) > 0 {
			username := emailParts[0]
			if len(username) >= 3 && len(username) <= 50 {
				return username
			}
		}
	}
	
	// Fallback: generar username basado en UID
	return "user" + userRecord.UID[:8]
}

// Logout maneja POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// En Firebase, el logout se maneja del lado del cliente
	// Aquí podríamos registrar el evento de logout
	
	log.Info("User logout requested")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Logout successful",
	})
}

// CheckAuthStatus maneja GET /auth/status
func (h *AuthHandler) CheckAuthStatus(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener token del header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
			"message":       "No token provided",
		})
		return
	}

	// Verificar formato Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
			"message":       "Invalid token format",
		})
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	token := parts[1]
	decodedToken, err := h.firebaseAuth.VerifyIDToken(context.Background(), token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
			"message":       "Invalid token",
		})
		return
	}

	// Obtener información del usuario
	userRecord, err := h.firebaseAuth.GetUser(context.Background(), decodedToken.UID)
	if err != nil {
		log.WithError(err).Error("Failed to get user from Firebase")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": true,
			"user_id":       decodedToken.UID,
			"expires_at":    decodedToken.Expires,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"firebase_id":    decodedToken.UID,
			"email":          userRecord.Email,
			"email_verified": userRecord.EmailVerified,
			"display_name":   userRecord.DisplayName,
		},
		"expires_at": decodedToken.Expires,
	})
}

// RefreshToken maneja POST /auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// En Firebase, el refresh de tokens se maneja del lado del cliente
	log.Info("Token refresh requested")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Token refresh should be handled on the client side using Firebase SDK",
	})
}

// GetProfile maneja GET /auth/profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto (debe estar autenticado)
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// Obtener información del usuario de Firebase
	userRecord, err := h.firebaseAuth.GetUser(context.Background(), userID.(string))
	if err != nil {
		log.WithError(err).Error("Failed to get user profile from Firebase")
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	profile := map[string]interface{}{
		"firebase_id":    userRecord.UID,
		"email":          userRecord.Email,
		"email_verified": userRecord.EmailVerified,
		"display_name":   userRecord.DisplayName,
		"photo_url":      userRecord.PhotoURL,
		"phone_number":   userRecord.PhoneNumber,
		"disabled":       userRecord.Disabled,
		"created_at":     userRecord.UserMetadata.CreationTimestamp,
		"last_login_at":  userRecord.UserMetadata.LastLogInTimestamp,
	}

	log.WithField("user_id", userID).Info("User profile retrieved")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    profile,
		"message": "Profile retrieved successfully",
	})
}

// UpdateProfile maneja PUT /auth/profile
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req struct {
		DisplayName string `json:"display_name" validate:"max=100"`
		PhotoURL    string `json:"photo_url" validate:"omitempty,url"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read request body")
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.WithError(err).Error("Failed to unmarshal JSON")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validar estructura
	if err := validator.ValidateStruct(&req); err != nil {
		log.WithError(err).Warn("Validation failed for profile update")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.WithField("user_id", userID).Info("Profile update requested (Firebase Admin SDK required)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Profile update functionality requires Firebase Admin SDK implementation",
	})
}

// ChangePassword maneja POST /auth/change-password
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	log.WithField("user_id", userID).Info("Password change requested")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Password change should be handled on the client side using Firebase SDK",
	})
}

// RevokeAllTokens maneja POST /auth/revoke-tokens
func (h *AuthHandler) RevokeAllTokens(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// Revocar todos los tokens del usuario
	err := h.firebaseAuth.RevokeRefreshTokens(context.Background(), userID.(string))
	if err != nil {
		log.WithError(err).Error("Failed to revoke tokens")
		http.Error(w, "Failed to revoke tokens", http.StatusInternalServerError)
		return
	}

	log.WithField("user_id", userID).Info("All tokens revoked successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "All tokens revoked successfully",
	})
}

// GetActiveSessions maneja GET /auth/sessions
func (h *AuthHandler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	log.WithField("user_id", userID).Info("Active sessions requested (not implemented)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Active sessions endpoint not implemented yet",
		"data":    []interface{}{},
	})
}

// RevokeSession maneja DELETE /auth/sessions/{session_id}
func (h *AuthHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	log.WithField("user_id", userID).Info("Session revocation requested (not implemented)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Session revocation endpoint not implemented yet",
	})
}

// GoogleLogin maneja POST /auth/google
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req struct {
		IDToken string `json:"id_token" validate:"required"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read request body")
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.WithError(err).Error("Failed to unmarshal JSON")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		log.WithError(err).Warn("Validation failed for Google login")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// Verificar token de Google a través de Firebase
	token, err := h.firebaseAuth.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		log.WithError(err).Warn("Invalid Google token")
		http.Error(w, "Invalid Google token", http.StatusUnauthorized)
		return
	}

	// Obtener información del usuario
	userRecord, err := h.firebaseAuth.GetUser(context.Background(), token.UID)
	if err != nil {
		log.WithError(err).Error("Failed to get user from Firebase")
		http.Error(w, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	// Verificar que realmente sea un login de Google
	isGoogleProvider := false
	for _, provider := range userRecord.ProviderUserInfo {
		if provider.ProviderID == "google.com" {
			isGoogleProvider = true
			break
		}
	}

	if !isGoogleProvider {
		log.Warn("Token is not from Google provider")
		http.Error(w, "Invalid Google authentication", http.StatusBadRequest)
		return
	}

	firstName, lastName := h.extractNames(userRecord.DisplayName)

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"firebase_id":    token.UID,
			"email":          userRecord.Email,
			"email_verified": userRecord.EmailVerified,
			"username":       h.generateUsername(userRecord),
			"first_name":     firstName,
			"last_name":      lastName,
			"photo_url":      userRecord.PhotoURL,
			"provider":       "google.com",
			"status":         "active",
		},
		"provider": "google.com",
		"message":  "Google login successful",
	}

	log.WithFields(map[string]interface{}{
		"firebase_id": token.UID,
		"provider":    "google.com",
		"email":       userRecord.Email,
	}).Info("User logged in with Google successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// FacebookLogin maneja POST /auth/facebook
func (h *AuthHandler) FacebookLogin(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req struct {
		IDToken string `json:"id_token" validate:"required"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read request body")
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.WithError(err).Error("Failed to unmarshal JSON")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		log.WithError(err).Warn("Validation failed for Facebook login")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// Verificar token de Facebook a través de Firebase
	token, err := h.firebaseAuth.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		log.WithError(err).Warn("Invalid Facebook token")
		http.Error(w, "Invalid Facebook token", http.StatusUnauthorized)
		return
	}

	// Obtener información del usuario
	userRecord, err := h.firebaseAuth.GetUser(context.Background(), token.UID)
	if err != nil {
		log.WithError(err).Error("Failed to get user from Firebase")
		http.Error(w, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	// Verificar que realmente sea un login de Facebook
	isFacebookProvider := false
	for _, provider := range userRecord.ProviderUserInfo {
		if provider.ProviderID == "facebook.com" {
			isFacebookProvider = true
			break
		}
	}

	if !isFacebookProvider {
		log.Warn("Token is not from Facebook provider")
		http.Error(w, "Invalid Facebook authentication", http.StatusBadRequest)
		return
	}

	firstName, lastName := h.extractNames(userRecord.DisplayName)

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"firebase_id":    token.UID,
			"email":          userRecord.Email,
			"email_verified": userRecord.EmailVerified,
			"username":       h.generateUsername(userRecord),
			"first_name":     firstName,
			"last_name":      lastName,
			"photo_url":      userRecord.PhotoURL,
			"provider":       "facebook.com",
			"status":         "active",
		},
		"provider": "facebook.com",
		"message":  "Facebook login successful",
	}

	log.WithFields(map[string]interface{}{
		"firebase_id": token.UID,
		"provider":    "facebook.com",
		"email":       userRecord.Email,
	}).Info("User logged in with Facebook successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// EmailPasswordLogin maneja POST /auth/email
func (h *AuthHandler) EmailPasswordLogin(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req struct {
		IDToken string `json:"id_token" validate:"required"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read request body")
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.WithError(err).Error("Failed to unmarshal JSON")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		log.WithError(err).Warn("Validation failed for email login")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// Verificar token de Firebase
	token, err := h.firebaseAuth.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		log.WithError(err).Warn("Invalid Firebase token")
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Obtener información del usuario
	userRecord, err := h.firebaseAuth.GetUser(context.Background(), token.UID)
	if err != nil {
		log.WithError(err).Error("Failed to get user from Firebase")
		http.Error(w, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	// Verificar que sea autenticación por email/password
	isEmailProvider := false
	for _, provider := range userRecord.ProviderUserInfo {
		if provider.ProviderID == "password" {
			isEmailProvider = true
			break
		}
	}

	if !isEmailProvider {
		log.Warn("Token is not from email/password provider")
		http.Error(w, "Invalid email/password authentication", http.StatusBadRequest)
		return
	}

	firstName, lastName := h.extractNames(userRecord.DisplayName)

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"firebase_id":    token.UID,
			"email":          userRecord.Email,
			"email_verified": userRecord.EmailVerified,
			"username":       h.generateUsername(userRecord),
			"first_name":     firstName,
			"last_name":      lastName,
			"photo_url":      userRecord.PhotoURL,
			"provider":       "password",
			"status":         "active",
		},
		"provider": "password",
		"message":  "Email login successful",
	}

	log.WithFields(map[string]interface{}{
		"firebase_id": token.UID,
		"provider":    "password",
		"email":       userRecord.Email,
	}).Info("User logged in with email/password successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}