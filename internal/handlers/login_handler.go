package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"it-app_user/internal/logger"
	"it-app_user/internal/validator"
	"it-app_user/pkg/firebase"
)

type LoginHandler struct {
	firebaseAuth *firebase.Auth
}

func NewLoginHandler(firebaseAuth *firebase.Auth) *LoginHandler {
	return &LoginHandler{
		firebaseAuth: firebaseAuth,
	}
}

// UpdateLoginInfo maneja POST /users/{id}/login
func (h *LoginHandler) UpdateLoginInfo(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid user ID provided for login update")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		LoginIP     string `json:"login_ip" validate:"omitempty,ip"`
		LoginDevice string `json:"login_device" validate:"max=255"`
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
		log.WithError(err).Warn("Validation failed for login info update")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mock login info update - en producción esto se enviaría a otro servicio
	log.WithField("user_id", id).WithField("login_ip", req.LoginIP).WithField("login_device", req.LoginDevice).Info("Login info updated successfully (mock)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login info updated successfully",
	})
}

// TrackUserLogin maneja POST /login/track
func (h *LoginHandler) TrackUserLogin(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	var req struct {
		UserID      int    `json:"user_id" validate:"required"`
		LoginIP     string `json:"login_ip" validate:"omitempty,ip"`
		LoginDevice string `json:"login_device" validate:"max=255"`
		LoginMethod string `json:"login_method" validate:"required,oneof=firebase google email"`
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
		log.WithError(err).Warn("Validation failed for track login request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mock login tracking - en producción esto se enviaría a otro servicio
	log.WithFields(map[string]interface{}{
		"user_id":      req.UserID,
		"login_ip":     req.LoginIP,
		"login_device": req.LoginDevice,
		"login_method": req.LoginMethod,
	}).Info("User login tracked successfully (mock)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login tracked successfully",
	})
}

// GetLoginHistory maneja GET /login/history/{user_id}
func (h *LoginHandler) GetLoginHistory(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.WithError(err).Warn("Invalid user ID provided for login history")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Mock login history - en producción esto vendría de otro servicio
	loginHistory := []map[string]interface{}{
		{
			"id":           1,
			"user_id":      userID,
			"login_time":   "2024-01-01T10:00:00Z",
			"login_ip":     "192.168.1.1",
			"login_device": "Chrome/Windows",
			"login_method": "firebase",
			"success":      true,
		},
		{
			"id":           2,
			"user_id":      userID,
			"login_time":   "2024-01-01T09:00:00Z",
			"login_ip":     "192.168.1.2",
			"login_device": "Safari/Mac",
			"login_method": "google",
			"success":      true,
		},
	}

	log.WithField("user_id", userID).Info("Login history retrieved successfully (mock)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    loginHistory,
		"count":   len(loginHistory),
		"message": "Login history retrieved successfully",
	})
}

// GetLoginAttempts maneja GET /login/attempts/{email}
func (h *LoginHandler) GetLoginAttempts(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	email := vars["email"]
	if email == "" {
		log.Warn("Email is required for login attempts")
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Mock login attempts data
	attempts := []map[string]interface{}{
		{
			"timestamp":    "2024-01-01T10:00:00Z",
			"ip_address":   "192.168.1.1",
			"user_agent":   "Chrome/Windows",
			"success":      true,
			"failure_reason": nil,
		},
		{
			"timestamp":    "2024-01-01T09:45:00Z",
			"ip_address":   "192.168.1.1",
			"user_agent":   "Chrome/Windows",
			"success":      false,
			"failure_reason": "Invalid password",
		},
	}

	log.WithField("email", email).Info("Login attempts retrieved successfully (mock)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    attempts,
		"count":   len(attempts),
		"email":   email,
		"message": "Login attempts retrieved successfully",
	})
}

// SecurityCheck maneja POST /login/security-check
func (h *LoginHandler) SecurityCheck(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	var req struct {
		Email     string `json:"email" validate:"required,email"`
		IPAddress string `json:"ip_address" validate:"omitempty,ip"`
		UserAgent string `json:"user_agent" validate:"max=500"`
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
		log.WithError(err).Warn("Validation failed for security check request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mock security check
	securityCheck := map[string]interface{}{
		"is_safe":           true,
		"risk_level":        "low",
		"blocked":           false,
		"requires_2fa":      false,
		"suspicious_activity": false,
		"recommendations":   []string{},
	}

	log.WithField("email", req.Email).Info("Security check completed")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    securityCheck,
		"message": "Security check completed",
	})
}

// GetMyLoginHistory maneja GET /login/my-history
func (h *LoginHandler) GetMyLoginHistory(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parámetros de paginación
	limit := 20 // default
	offset := 0 // default
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Mock login history
	loginHistory := []map[string]interface{}{
		{
			"id":           1,
			"login_time":   "2024-01-01T10:00:00Z",
			"login_ip":     "192.168.1.1",
			"login_device": "Chrome/Windows",
			"login_method": "firebase",
			"success":      true,
			"location":     "New York, US",
		},
		{
			"id":           2,
			"login_time":   "2024-01-01T09:00:00Z",
			"login_ip":     "192.168.1.2",
			"login_device": "Safari/Mac",
			"login_method": "google",
			"success":      true,
			"location":     "California, US",
		},
	}

	log.WithField("user_id", userID).Info("My login history retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    loginHistory,
		"count":   len(loginHistory),
		"limit":   limit,
		"offset":  offset,
		"message": "Login history retrieved successfully",
	})
}

// GetActiveSessions maneja GET /login/active-sessions
func (h *LoginHandler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Mock active sessions
	activeSessions := []map[string]interface{}{
		{
			"session_id":   "session_123",
			"device":       "Chrome/Windows",
			"ip_address":   "192.168.1.1",
			"location":     "New York, US",
			"login_time":   "2024-01-01T10:00:00Z",
			"last_seen":    "2024-01-01T12:00:00Z",
			"is_current":   true,
		},
		{
			"session_id":   "session_456",
			"device":       "Safari/iPhone",
			"ip_address":   "192.168.1.100",
			"location":     "New York, US",
			"login_time":   "2024-01-01T08:00:00Z",
			"last_seen":    "2024-01-01T11:30:00Z",
			"is_current":   false,
		},
	}

	log.WithField("user_id", userID).Info("Active sessions retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    activeSessions,
		"count":   len(activeSessions),
		"message": "Active sessions retrieved successfully",
	})
}

// TerminateSession maneja DELETE /login/terminate-session/{session_id}
func (h *LoginHandler) TerminateSession(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	sessionID := vars["session_id"]
	if sessionID == "" {
		log.Warn("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	log.WithFields(map[string]interface{}{
		"user_id":    userID,
		"session_id": sessionID,
	}).Info("Session terminated successfully (mock)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Session terminated successfully",
	})
}

// TerminateAllSessions maneja DELETE /login/terminate-all-sessions
func (h *LoginHandler) TerminateAllSessions(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if h.firebaseAuth != nil {
		// Revocar todos los tokens de Firebase
		err := h.firebaseAuth.RevokeRefreshTokens(r.Context(), userID.(string))
		if err != nil {
			log.WithError(err).Error("Failed to revoke Firebase tokens")
			http.Error(w, "Failed to terminate sessions", http.StatusInternalServerError)
			return
		}
	}

	log.WithField("user_id", userID).Info("All sessions terminated successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "All sessions terminated successfully",
	})
}

// GetSuspiciousActivity maneja GET /login/suspicious-activity
func (h *LoginHandler) GetSuspiciousActivity(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Mock suspicious activity data
	suspiciousActivity := []map[string]interface{}{
		{
			"id":          1,
			"timestamp":   "2024-01-01T03:00:00Z",
			"activity":    "Multiple failed login attempts",
			"ip_address":  "192.168.1.50",
			"location":    "Unknown",
			"risk_level":  "medium",
			"blocked":     true,
		},
		{
			"id":          2,
			"timestamp":   "2024-01-01T02:30:00Z",
			"activity":    "Login from new device",
			"ip_address":  "192.168.1.200",
			"location":    "California, US",
			"risk_level":  "low",
			"blocked":     false,
		},
	}

	log.WithField("user_id", userID).Info("Suspicious activity retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    suspiciousActivity,
		"count":   len(suspiciousActivity),
		"message": "Suspicious activity retrieved successfully",
	})
}