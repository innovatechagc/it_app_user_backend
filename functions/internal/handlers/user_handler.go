package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"it-app_user/internal/logger"
	"it-app_user/internal/models"
	"it-app_user/internal/repositories"
	"it-app_user/internal/validator"
)

type UserHandler struct {
	userRepo repositories.UserRepositoryInterface
}

func NewUserHandler(userRepo repositories.UserRepositoryInterface) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// HealthCheck maneja GET /health
func (h *UserHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "service": "user-service"}`))
	
	log.Info("Health check requested")
}

// GetAllUsers maneja GET /users
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	// Parámetros de paginación
	limit := 50 // default
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
	
	users, err := h.userRepo.GetAll(limit, offset)
	if err != nil {
		log.WithError(err).Error("Failed to fetch users")
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	log.WithField("count", len(users)).Info("Users retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    users,
		"count":   len(users),
		"limit":   limit,
		"offset":  offset,
		"message": "Users retrieved successfully",
	})
}

// GetUserByID maneja GET /users/{id}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByID(uint(id))
	if err != nil {
		log.WithError(err).WithField("user_id", id).Error("Failed to fetch user")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	log.WithField("user_id", id).Info("User retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    user,
		"message": "User retrieved successfully",
	})
}

// CreateUser maneja POST /users/create
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var req models.CreateUserRequest

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

	// Establecer valores por defecto
	if req.Status == "" {
		req.Status = "active"
	}

	// Validar estructura
	if err := validator.ValidateStruct(&req); err != nil {
		log.WithError(err).Warn("Validation failed for create user request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Crear usuario usando el repositorio
	user := &models.User{
		FirebaseID:    req.FirebaseID,
		Email:         req.Email,
		EmailVerified: req.EmailVerified,
		Username:      req.Username,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Provider:      req.Provider,
		ProviderID:    req.ProviderID,
		Status:        req.Status,
	}

	if err := h.userRepo.Create(user); err != nil {
		log.WithError(err).Error("Failed to create user")
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	log.WithField("user_id", user.ID).Info("User created successfully")
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    user,
		"message": "User created successfully",
	})
}

// UpdateUser maneja PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateUserRequest

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
		log.WithError(err).Warn("Validation failed for update user request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener usuario existente
	user, err := h.userRepo.GetByID(uint(id))
	if err != nil {
		log.WithError(err).WithField("user_id", id).Error("User not found for update")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Actualizar campos
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.EmailVerified != nil {
		user.EmailVerified = *req.EmailVerified
	}
	if req.Disabled != nil {
		user.Disabled = *req.Disabled
	}

	// Guardar cambios
	if err := h.userRepo.Update(user); err != nil {
		log.WithError(err).WithField("user_id", id).Error("Failed to update user")
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	log.WithField("user_id", id).Info("User updated successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    user,
		"message": "User updated successfully",
	})
}

// DeleteUser maneja DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).Warn("Invalid user ID provided")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Verificar que el usuario existe antes de eliminarlo
	_, err = h.userRepo.GetByID(uint(id))
	if err != nil {
		log.WithError(err).WithField("user_id", id).Error("User not found for deletion")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Eliminar usuario
	if err := h.userRepo.Delete(uint(id)); err != nil {
		log.WithError(err).WithField("user_id", id).Error("Failed to delete user")
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	log.WithField("user_id", id).Info("User deleted successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User deleted successfully",
	})
}

// GetUserByFirebaseID maneja GET /users/firebase/{firebase_id}
func (h *UserHandler) GetUserByFirebaseID(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	firebaseID := vars["firebase_id"]
	if firebaseID == "" {
		log.Warn("Firebase ID is required but not provided")
		http.Error(w, "Firebase ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByFirebaseID(firebaseID)
	if err != nil {
		log.WithError(err).WithField("firebase_id", firebaseID).Error("Failed to fetch user by Firebase ID")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	log.WithField("firebase_id", firebaseID).Info("User retrieved successfully by Firebase ID")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    user,
		"message": "User retrieved successfully",
	})
}

// GetUserByUsername maneja GET /users/username/{username}
func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		log.Warn("Username is required but not provided")
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByUsername(username)
	if err != nil {
		log.WithError(err).WithField("username", username).Error("Failed to fetch user by username")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	log.WithField("username", username).Info("User retrieved successfully by username")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    user,
		"message": "User retrieved successfully",
	})
}

// GetUserByEmail maneja GET /users/email/{email}
func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	email := vars["email"]
	if email == "" {
		log.Warn("Email is required but not provided")
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByEmail(email)
	if err != nil {
		log.WithError(err).WithField("email", email).Error("Failed to fetch user by email")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	log.WithField("email", email).Info("User retrieved successfully by email")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    user,
		"message": "User retrieved successfully",
	})
}

// SearchUsers maneja GET /users/search
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	query := r.URL.Query().Get("q")
	if query == "" {
		log.Warn("Search query is required")
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	// Parámetros de paginación
	limit := 20 // default para búsquedas
	offset := 0 // default
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}
	
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.userRepo.SearchUsers(query, limit, offset)
	if err != nil {
		log.WithError(err).WithField("query", query).Error("Failed to search users")
		http.Error(w, "Error searching users", http.StatusInternalServerError)
		return
	}

	log.WithFields(map[string]interface{}{
		"query": query,
		"count": len(users),
	}).Info("Users search completed")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    users,
		"count":   len(users),
		"query":   query,
		"limit":   limit,
		"offset":  offset,
		"message": "Search completed successfully",
	})
}

// CountUsers maneja GET /users/count
func (h *UserHandler) CountUsers(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	count, err := h.userRepo.CountUsers()
	if err != nil {
		log.WithError(err).Error("Failed to count users")
		http.Error(w, "Error counting users", http.StatusInternalServerError)
		return
	}

	log.WithField("total_users", count).Info("Users counted successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":   count,
		"message": "Users counted successfully",
	})
}

// GetActiveUsers maneja GET /users/active
func (h *UserHandler) GetActiveUsers(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	
	users, err := h.userRepo.GetActiveUsers()
	if err != nil {
		log.WithError(err).Error("Failed to fetch active users")
		http.Error(w, "Error fetching active users", http.StatusInternalServerError)
		return
	}

	log.WithField("count", len(users)).Info("Active users retrieved successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    users,
		"count":   len(users),
		"message": "Active users retrieved successfully",
	})
}

// UpdateLoginInfo maneja POST /users/{id}/login
func (h *UserHandler) UpdateLoginInfo(w http.ResponseWriter, r *http.Request) {
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

	// Actualizar información de login
	if err := h.userRepo.UpdateLoginInfo(uint(id), req.LoginIP, req.LoginDevice); err != nil {
		log.WithError(err).WithField("user_id", id).Error("Failed to update login info")
		http.Error(w, "Error updating login info", http.StatusInternalServerError)
		return
	}

	log.WithField("user_id", id).Info("Login info updated successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login info updated successfully",
	})
}

// GetUserProfile maneja GET /users/{id}/profile - Placeholder
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	log.WithField("user_id", idStr).Info("User profile requested (not implemented)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User profile endpoint not implemented yet",
		"user_id": idStr,
	})
}

// GetUserSettings maneja GET /users/{id}/settings - Placeholder
func (h *UserHandler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	log.WithField("user_id", idStr).Info("User settings requested (not implemented)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User settings endpoint not implemented yet",
		"user_id": idStr,
	})
}

// GetUserStats maneja GET /users/{id}/stats - Placeholder
func (h *UserHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	vars := mux.Vars(r)
	idStr := vars["id"]
	
	log.WithField("user_id", idStr).Info("User stats requested (not implemented)")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User stats endpoint not implemented yet",
		"user_id": idStr,
	})
}


