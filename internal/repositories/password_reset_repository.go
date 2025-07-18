package repositories

import (
	"time"
	"gorm.io/gorm"
	"it-app_user/internal/models"
)

type PasswordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository crea una nueva instancia del repositorio de reset de contraseña
func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepositoryInterface {
	return &PasswordResetRepository{db: db}
}

// GetByToken obtiene un token de reset por su token
func (r *PasswordResetRepository) GetByToken(token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.Where("token = ? AND is_used = ? AND expires_at > ?", token, false, time.Now()).First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

// GetByCode obtiene un token de reset por su código
func (r *PasswordResetRepository) GetByCode(code string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.Where("code = ? AND is_used = ? AND expires_at > ?", code, false, time.Now()).First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

// GetByUserID obtiene el token de reset más reciente de un usuario
func (r *PasswordResetRepository) GetByUserID(userID uint) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

// Create crea un nuevo token de reset de contraseña
func (r *PasswordResetRepository) Create(resetToken *models.PasswordResetToken) error {
	return r.db.Create(resetToken).Error
}

// Update actualiza un token de reset existente
func (r *PasswordResetRepository) Update(resetToken *models.PasswordResetToken) error {
	return r.db.Save(resetToken).Error
}

// Delete elimina un token de reset por su ID
func (r *PasswordResetRepository) Delete(id uint) error {
	return r.db.Delete(&models.PasswordResetToken{}, id).Error
}

// MarkAsUsed marca un token como usado
func (r *PasswordResetRepository) MarkAsUsed(id uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"is_used":    true,
		"used_at":    &now,
		"updated_at": now,
	}
	
	return r.db.Model(&models.PasswordResetToken{}).Where("id = ?", id).Updates(updates).Error
}

// CleanExpiredTokens elimina todos los tokens expirados
func (r *PasswordResetRepository) CleanExpiredTokens() error {
	return r.db.Where("expires_at < ? OR is_used = ?", time.Now(), true).Delete(&models.PasswordResetToken{}).Error
}