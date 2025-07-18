package repositories

import (
	"time"
	"gorm.io/gorm"
	"it-app_user/internal/models"
)

type EmailVerificationRepository struct {
	db *gorm.DB
}

// NewEmailVerificationRepository crea una nueva instancia del repositorio de verificación de email
func NewEmailVerificationRepository(db *gorm.DB) EmailVerificationRepositoryInterface {
	return &EmailVerificationRepository{db: db}
}

// GetByUserID obtiene la verificación de email por ID de usuario
func (r *EmailVerificationRepository) GetByUserID(userID uint) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.db.Where("user_id = ?", userID).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// GetByEmail obtiene la verificación de email por email
func (r *EmailVerificationRepository) GetByEmail(email string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.db.Where("email = ?", email).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// GetByFirebaseID obtiene la verificación de email por Firebase ID
func (r *EmailVerificationRepository) GetByFirebaseID(firebaseID string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.db.Where("firebase_id = ?", firebaseID).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// Create crea una nueva verificación de email
func (r *EmailVerificationRepository) Create(verification *models.EmailVerification) error {
	return r.db.Create(verification).Error
}

// Update actualiza una verificación de email existente
func (r *EmailVerificationRepository) Update(verification *models.EmailVerification) error {
	return r.db.Save(verification).Error
}

// Delete elimina una verificación de email por su ID
func (r *EmailVerificationRepository) Delete(id uint) error {
	return r.db.Delete(&models.EmailVerification{}, id).Error
}

// MarkAsVerified marca un email como verificado
func (r *EmailVerificationRepository) MarkAsVerified(userID uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"is_verified":  true,
		"verified_at":  &now,
		"updated_at":   now,
	}
	
	return r.db.Model(&models.EmailVerification{}).Where("user_id = ?", userID).Updates(updates).Error
}

// IncrementAttempts incrementa el contador de intentos de verificación
func (r *EmailVerificationRepository) IncrementAttempts(userID uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"attempts_count":  gorm.Expr("attempts_count + 1"),
		"last_attempt_at": &now,
		"updated_at":      now,
	}
	
	return r.db.Model(&models.EmailVerification{}).Where("user_id = ?", userID).Updates(updates).Error
}

// GetPendingVerifications obtiene todas las verificaciones pendientes
func (r *EmailVerificationRepository) GetPendingVerifications() ([]models.EmailVerification, error) {
	var verifications []models.EmailVerification
	err := r.db.Where("is_verified = ?", false).Find(&verifications).Error
	return verifications, err
}