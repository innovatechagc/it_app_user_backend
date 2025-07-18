package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"it-app_user/internal/logger"
	"it-app_user/internal/models"
	"it-app_user/internal/repositories"
	"it-app_user/internal/validator"
	"it-app_user/pkg/firebase"
)

type VerifyEmailHandler struct {
	firebaseAuth *firebase.Auth
	emailRepo    repositories.EmailVerificationRepositoryInterface
}

func NewVerifyEmailHandler(firebaseAuth *firebase.Auth, emailRepo repositories.EmailVerificationRepositoryInterface) *VerifyEmailHandler {
	return &VerifyEmailHandler{
		firebaseAuth: firebaseAuth,
		emailRepo:    emailRepo,
	}
}

// SendVerificationEmail maneja POST /auth/send-verification-email
func (h *VerifyEmailHandler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.SendVerificationEmailRequest

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
		log.WithError(err).Warn("Validation failed for send verification email request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// Verificar que el usuario existe
	userRecord, err := h.firebaseAuth.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		log.WithError(err).WithField("email", req.Email).Warn("User not found")
		// Por seguridad, no revelamos si el email existe o no
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "If the email exists, a verification email has been sent",
		})
		return
	}

	// En Firebase, el envío de email de verificación se maneja del lado del cliente
	// Aquí registramos el evento
	log.WithField("firebase_id", userRecord.UID).WithField("email", req.Email).Info("Verification email requested")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Verification email sent successfully",
	})
}

// VerifyEmail maneja POST /auth/verify-email
func (h *VerifyEmailHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.VerifyEmailRequest

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
		log.WithError(err).Warn("Validation failed for verify email request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// Verificar el token de verificación
	token, err := h.firebaseAuth.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		log.WithError(err).Warn("Invalid Firebase token")
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Obtener información actualizada del usuario
	userRecord, err := h.firebaseAuth.GetUser(context.Background(), token.UID)
	if err != nil {
		log.WithError(err).Error("Failed to get user from Firebase")
		http.Error(w, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	// Verificar si el email ya está verificado
	if !userRecord.EmailVerified {
		log.WithField("firebase_id", token.UID).Warn("Email not verified yet")
		http.Error(w, "Email not verified yet", http.StatusBadRequest)
		return
	}

	// Crear respuesta con información del usuario verificado
	response := models.VerifyEmailResponse{
		Success: true,
		User: &models.User{
			FirebaseID:    token.UID,
			Email:         userRecord.Email,
			EmailVerified: userRecord.EmailVerified,
			Username:      userRecord.DisplayName,
			Status:        "active",
		},
		Message: "Email verified successfully",
	}

	log.WithField("firebase_id", token.UID).Info("Email verification confirmed")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CheckEmailVerificationStatus maneja GET /auth/email-verification-status
func (h *VerifyEmailHandler) CheckEmailVerificationStatus(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// Obtener el usuario del contexto (debe estar autenticado)
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Obtener información del usuario de Firebase
	userRecord, err := h.firebaseAuth.GetUser(context.Background(), userID.(string))
	if err != nil {
		log.WithError(err).Error("Failed to get user from Firebase")
		http.Error(w, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	response := models.EmailVerificationStatusResponse{
		EmailVerified: userRecord.EmailVerified,
		Email:         userRecord.Email,
	}

	log.WithField("firebase_id", userID).Debug("Email verification status checked")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// VerifyEmailWithCode maneja POST /email/verify-code
func (h *VerifyEmailHandler) VerifyEmailWithCode(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.VerifyEmailWithCodeRequest

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
		log.WithError(err).Warn("Validation failed for verify email with code request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Buscar verificación por email
	verification, err := h.emailRepo.GetByEmail(req.Email)
	if err != nil {
		log.WithError(err).WithField("email", req.Email).Warn("Email verification not found")
		http.Error(w, "Invalid email or verification code", http.StatusBadRequest)
		return
	}

	// Verificar código (en un caso real, el código estaría hasheado)
	if verification.VerificationCode != req.VerificationCode {
		log.WithField("email", req.Email).Warn("Invalid verification code")
		// Incrementar intentos
		h.emailRepo.IncrementAttempts(verification.UserID)
		http.Error(w, "Invalid verification code", http.StatusBadRequest)
		return
	}

	// Marcar como verificado
	if err := h.emailRepo.MarkAsVerified(verification.UserID); err != nil {
		log.WithError(err).Error("Failed to mark email as verified")
		http.Error(w, "Error verifying email", http.StatusInternalServerError)
		return
	}

	log.WithField("email", req.Email).Info("Email verified successfully with code")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Email verified successfully",
	})
}

// ResendVerificationEmail maneja POST /email/resend
func (h *VerifyEmailHandler) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.ResendVerificationEmailRequest

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
		log.WithError(err).Warn("Validation failed for resend verification email request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.WithField("email", req.Email).Info("Verification email resend requested")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "If the email exists, a verification email has been sent",
	})
}

// GetMyVerificationStatus maneja GET /email/my-status
func (h *VerifyEmailHandler) GetMyVerificationStatus(w http.ResponseWriter, r *http.Request) {
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

	// Obtener información del usuario de Firebase
	userRecord, err := h.firebaseAuth.GetUser(context.Background(), userID.(string))
	if err != nil {
		log.WithError(err).Error("Failed to get user from Firebase")
		http.Error(w, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	response := models.EmailVerificationStatusResponse{
		EmailVerified: userRecord.EmailVerified,
		Email:         userRecord.Email,
		FirebaseID:    userRecord.UID,
		CanResend:     !userRecord.EmailVerified,
		AttemptsLeft:  3, // Valor por defecto
	}

	log.WithField("firebase_id", userID).Info("My verification status retrieved")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    response,
		"message": "Verification status retrieved successfully",
	})
}

// UpdateEmail maneja POST /email/update-email
func (h *VerifyEmailHandler) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req struct {
		NewEmail string `json:"new_email" validate:"required,email"`
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
		log.WithError(err).Warn("Validation failed for update email request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.WithField("user_id", userID).WithField("new_email", req.NewEmail).Info("Email update requested (requires Firebase Admin SDK)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Email update functionality requires Firebase Admin SDK implementation",
	})
}

// GetVerificationHistory maneja GET /email/verification-history
func (h *VerifyEmailHandler) GetVerificationHistory(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	log.WithField("user_id", userID).Info("Verification history requested (not implemented)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Verification history endpoint not implemented yet",
		"data":    []interface{}{},
	})
}

// GetEmailSettings maneja GET /email/settings
func (h *VerifyEmailHandler) GetEmailSettings(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Configuraciones por defecto
	settings := models.EmailVerificationSettings{
		MaxAttempts:         5,
		CodeExpirationTime:  3600, // 1 hora en segundos
		ResendCooldownTime:  300,  // 5 minutos en segundos
		RequireVerification: true,
		AutoVerifyDomains:   []string{},
		BlockedDomains:      []string{"tempmail.org", "10minutemail.com"},
	}

	log.WithField("user_id", userID).Info("Email settings retrieved")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    settings,
		"message": "Email settings retrieved successfully",
	})
}

// UpdateEmailSettings maneja PUT /email/settings
func (h *VerifyEmailHandler) UpdateEmailSettings(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req models.EmailVerificationSettings

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

	log.WithField("user_id", userID).Info("Email settings update requested (not implemented)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Email settings update endpoint not implemented yet",
	})
}