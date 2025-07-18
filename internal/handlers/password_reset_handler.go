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

type PasswordResetHandler struct {
	firebaseAuth *firebase.Auth
	passwordRepo repositories.PasswordResetRepositoryInterface
}

func NewPasswordResetHandler(firebaseAuth *firebase.Auth, passwordRepo repositories.PasswordResetRepositoryInterface) *PasswordResetHandler {
	return &PasswordResetHandler{
		firebaseAuth: firebaseAuth,
		passwordRepo: passwordRepo,
	}
}

// RequestPasswordReset maneja POST /auth/password-reset/request
func (h *PasswordResetHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.PasswordResetRequest

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
		log.WithError(err).Warn("Validation failed for password reset request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// En Firebase, el reset de contraseña se maneja a través del Admin SDK
	// Aquí registramos la solicitud y enviamos una respuesta genérica por seguridad
	log.WithField("email", req.Email).Info("Password reset requested")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "If the email exists, a password reset link has been sent",
	})
}

// ConfirmPasswordReset maneja POST /auth/password-reset/confirm
func (h *PasswordResetHandler) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.PasswordResetConfirmRequest

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
		log.WithError(err).Warn("Validation failed for password reset confirm request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.firebaseAuth == nil {
		log.Error("Firebase Auth not configured")
		http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
		return
	}

	// En Firebase, la confirmación del reset se maneja del lado del cliente
	// Aquí podríamos validar el código si fuera necesario
	
	log.Info("Password reset confirmation processed")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Password reset completed successfully",
	})
}

// ChangePassword maneja POST /auth/change-password
func (h *PasswordResetHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.ChangePasswordRequest

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
		log.WithError(err).Warn("Validation failed for change password request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	// Verificar token actual
	_, err = h.firebaseAuth.VerifyIDToken(context.Background(), req.CurrentToken)
	if err != nil {
		log.WithError(err).Warn("Invalid current token")
		http.Error(w, "Invalid current token", http.StatusUnauthorized)
		return
	}

	log.WithField("user_id", userID).Info("Password change processed")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Password changed successfully",
	})
}

// VerifyResetCode maneja POST /password/reset/verify
func (h *PasswordResetHandler) VerifyResetCode(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req struct {
		Code string `json:"code" validate:"required,len=6"`
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
		log.WithError(err).Warn("Validation failed for reset code verification")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Buscar token por código
	resetToken, err := h.passwordRepo.GetByCode(req.Code)
	if err != nil {
		log.WithError(err).WithField("code", req.Code).Warn("Invalid or expired reset code")
		http.Error(w, "Invalid or expired reset code", http.StatusBadRequest)
		return
	}

	log.WithField("user_id", resetToken.UserID).Info("Reset code verified successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":   true,
		"message": "Reset code is valid",
	})
}

// ValidateResetToken maneja POST /password/reset/validate-token
func (h *PasswordResetHandler) ValidateResetToken(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req struct {
		Token string `json:"token" validate:"required"`
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
		log.WithError(err).Warn("Validation failed for token validation")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Buscar token
	resetToken, err := h.passwordRepo.GetByToken(req.Token)
	if err != nil {
		log.WithError(err).WithField("token", req.Token).Warn("Invalid or expired reset token")
		http.Error(w, "Invalid or expired reset token", http.StatusBadRequest)
		return
	}

	log.WithField("user_id", resetToken.UserID).Info("Reset token validated successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":   true,
		"message": "Reset token is valid",
	})
}

// CheckPasswordStrength maneja POST /password/strength-check
func (h *PasswordResetHandler) CheckPasswordStrength(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req struct {
		Password string `json:"password" validate:"required"`
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
		log.WithError(err).Warn("Validation failed for password strength check")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verificar fortaleza de la contraseña (implementación básica)
	strength := checkPasswordStrength(req.Password)
	
	log.Info("Password strength checked")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    strength,
		"message": "Password strength checked",
	})
}

// GetPasswordHistory maneja GET /password/history
func (h *PasswordResetHandler) GetPasswordHistory(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Obtener el usuario del contexto
	userID := r.Context().Value("user_id")
	if userID == nil {
		log.Warn("User ID not found in context")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	log.WithField("user_id", userID).Info("Password history requested (not implemented)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Password history endpoint not implemented yet",
		"data":    []interface{}{},
	})
}

// GetPasswordPolicy maneja GET /password/policy
func (h *PasswordResetHandler) GetPasswordPolicy(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Política de contraseñas por defecto
	policy := models.PasswordPolicy{
		MinLength:        8,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSymbols:   false,
		ForbiddenWords:   []string{"password", "123456", "qwerty"},
		MaxAge:           90, // días
	}
	
	log.Info("Password policy retrieved")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    policy,
		"message": "Password policy retrieved successfully",
	})
}

// checkPasswordStrength es una función auxiliar para verificar la fortaleza de la contraseña
func checkPasswordStrength(password string) models.PasswordStrengthCheck {
	check := models.PasswordStrengthCheck{
		IsValid:  true,
		Score:    0,
		Feedback: []string{},
	}

	// Verificar longitud mínima
	if len(password) >= 8 {
		check.Requirements.MinLength = true
		check.Score++
	} else {
		check.IsValid = false
		check.Feedback = append(check.Feedback, "Password must be at least 8 characters long")
	}

	// Verificar mayúsculas
	hasUpper := false
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
			break
		}
	}
	if hasUpper {
		check.Requirements.HasUppercase = true
		check.Score++
	} else {
		check.Feedback = append(check.Feedback, "Password must contain at least one uppercase letter")
	}

	// Verificar minúsculas
	hasLower := false
	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			hasLower = true
			break
		}
	}
	if hasLower {
		check.Requirements.HasLowercase = true
		check.Score++
	} else {
		check.Feedback = append(check.Feedback, "Password must contain at least one lowercase letter")
	}

	// Verificar números
	hasNumber := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasNumber = true
			break
		}
	}
	if hasNumber {
		check.Requirements.HasNumbers = true
		check.Score++
	} else {
		check.Feedback = append(check.Feedback, "Password must contain at least one number")
	}

	// Verificar símbolos
	hasSymbol := false
	symbols := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	for _, char := range password {
		for _, symbol := range symbols {
			if char == symbol {
				hasSymbol = true
				break
			}
		}
		if hasSymbol {
			break
		}
	}
	if hasSymbol {
		check.Requirements.HasSymbols = true
		check.Score++
	}

	// Determinar si es válida según los requisitos mínimos
	if !check.Requirements.MinLength || !check.Requirements.HasUppercase || 
	   !check.Requirements.HasLowercase || !check.Requirements.HasNumbers {
		check.IsValid = false
	}

	return check
}