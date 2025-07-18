package repositories

import "it-app_user/internal/models"

// UserRepositoryInterface define los métodos para el repositorio de usuarios
type UserRepositoryInterface interface {
	// CRUD básico
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByFirebaseID(firebaseID string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetAll(limit, offset int) ([]models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	
	// Métodos específicos
	UpdateLoginInfo(id uint, loginIP, loginDevice string) error
	GetActiveUsers() ([]models.User, error)
	SearchUsers(query string, limit, offset int) ([]models.User, error)
	CountUsers() (int64, error)
}

// EmailVerificationRepositoryInterface define los métodos para verificación de email
type EmailVerificationRepositoryInterface interface {
	GetByUserID(userID uint) (*models.EmailVerification, error)
	GetByEmail(email string) (*models.EmailVerification, error)
	GetByFirebaseID(firebaseID string) (*models.EmailVerification, error)
	Create(verification *models.EmailVerification) error
	Update(verification *models.EmailVerification) error
	Delete(id uint) error
	MarkAsVerified(userID uint) error
	IncrementAttempts(userID uint) error
	GetPendingVerifications() ([]models.EmailVerification, error)
}

// PasswordResetRepositoryInterface define los métodos para reset de contraseña
type PasswordResetRepositoryInterface interface {
	GetByToken(token string) (*models.PasswordResetToken, error)
	GetByCode(code string) (*models.PasswordResetToken, error)
	GetByUserID(userID uint) (*models.PasswordResetToken, error)
	Create(resetToken *models.PasswordResetToken) error
	Update(resetToken *models.PasswordResetToken) error
	Delete(id uint) error
	MarkAsUsed(id uint) error
	CleanExpiredTokens() error
}

// UserProfileRepositoryInterface define los métodos para perfiles de usuario
type UserProfileRepositoryInterface interface {
	GetByUserID(userID uint) (*models.UserProfile, error)
	Create(profile *models.UserProfile) error
	Update(profile *models.UserProfile) error
	Delete(userID uint) error
	UpdateAvatar(userID uint, avatarURL string) error
}

// UserSettingsRepositoryInterface define los métodos para configuraciones de usuario
type UserSettingsRepositoryInterface interface {
	GetByUserID(userID uint) (*models.UserSettings, error)
	Create(settings *models.UserSettings) error
	Update(settings *models.UserSettings) error
	Delete(userID uint) error
	UpdateLanguage(userID uint, language string) error
	UpdateTheme(userID uint, theme string) error
}

// UserStatsRepositoryInterface define los métodos para estadísticas de usuario
type UserStatsRepositoryInterface interface {
	GetByUserID(userID uint) (*models.UserStats, error)
	Create(stats *models.UserStats) error
	Update(stats *models.UserStats) error
	Delete(userID uint) error
	IncrementLoginCount(userID uint) error
	IncrementProfileViews(userID uint) error
	UpdateLastActive(userID uint) error
}