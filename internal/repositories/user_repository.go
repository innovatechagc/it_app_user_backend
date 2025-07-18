package repositories

import (
	"time"
	"gorm.io/gorm"
	"it-app_user/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia del repositorio de usuarios
func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db: db}
}

// GetByID obtiene un usuario por su ID
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail obtiene un usuario por su email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByFirebaseID obtiene un usuario por su Firebase ID
func (r *UserRepository) GetByFirebaseID(firebaseID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("firebase_id = ?", firebaseID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername obtiene un usuario por su username
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll obtiene todos los usuarios con paginación
func (r *UserRepository) GetAll(limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// Create crea un nuevo usuario
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update actualiza un usuario existente
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete elimina un usuario por su ID
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// UpdateLoginInfo actualiza la información de login del usuario
func (r *UserRepository) UpdateLoginInfo(id uint, loginIP, loginDevice string) error {
	updates := map[string]interface{}{
		"login_count":        gorm.Expr("login_count + 1"),
		"last_login_at":      time.Now(),
		"last_login_ip":      &loginIP,
		"last_login_device":  &loginDevice,
		"updated_at":         time.Now(),
	}
	
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// GetActiveUsers obtiene todos los usuarios activos
func (r *UserRepository) GetActiveUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("status = ? AND disabled = ?", "active", false).Find(&users).Error
	return users, err
}

// SearchUsers busca usuarios por nombre, email o username
func (r *UserRepository) SearchUsers(query string, limit, offset int) ([]models.User, error) {
	var users []models.User
	searchPattern := "%" + query + "%"
	
	err := r.db.Where(
		"first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR username ILIKE ?",
		searchPattern, searchPattern, searchPattern, searchPattern,
	).Limit(limit).Offset(offset).Find(&users).Error
	
	return users, err
}

// CountUsers cuenta el total de usuarios
func (r *UserRepository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}