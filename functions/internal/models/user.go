package models

import (
	"time"
)

type User struct {
	ID              uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	FirebaseID      string     `json:"firebase_id" gorm:"uniqueIndex;size:128;not null"`
	Email           string     `json:"email" gorm:"uniqueIndex;size:255;not null"`
	EmailVerified   bool       `json:"email_verified" gorm:"default:false"`
	Username        string     `json:"username" gorm:"uniqueIndex;size:50;not null"`
	FirstName       string     `json:"first_name" gorm:"size:100"`
	LastName        string     `json:"last_name" gorm:"size:100"`
	Provider        string     `json:"provider" gorm:"size:50"`
	ProviderID      string     `json:"provider_id" gorm:"size:128"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	LoginCount      int        `json:"login_count" gorm:"default:0"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	LastLoginIP     *string    `json:"last_login_ip" gorm:"size:45"`
	LastLoginDevice *string    `json:"last_login_device" gorm:"size:255"`
	Disabled        bool       `json:"disabled" gorm:"default:false"`
	Status          string     `json:"status" gorm:"size:20;default:'active';check:status IN ('active','inactive','pending')"`
	PkAutUseID      string     `json:"pk_aut_use_id" gorm:"size:128"`
}

type CreateUserRequest struct {
	FirebaseID    string `json:"firebase_id" validate:"required,min=1,max=128"`
	Email         string `json:"email" validate:"required,email,max=255"`
	EmailVerified bool   `json:"email_verified"`
	Username      string `json:"username" validate:"required,min=3,max=50,alphanum"`
	FirstName     string `json:"first_name" validate:"max=100"`
	LastName      string `json:"last_name" validate:"max=100"`
	Provider      string `json:"provider" validate:"max=50"`
	ProviderID    string `json:"provider_id" validate:"max=128"`
	Status        string `json:"status" validate:"oneof=active inactive pending"`
}

type UpdateUserRequest struct {
	EmailVerified   *bool   `json:"email_verified"`
	Username        string  `json:"username" validate:"omitempty,min=3,max=50,alphanum"`
	FirstName       string  `json:"first_name" validate:"max=100"`
	LastName        string  `json:"last_name" validate:"max=100"`
	Provider        string  `json:"provider" validate:"max=50"`
	ProviderID      string  `json:"provider_id" validate:"max=128"`
	LastLoginIP     string  `json:"last_login_ip" validate:"omitempty,ip"`
	LastLoginDevice string  `json:"last_login_device" validate:"max=255"`
	Disabled        *bool   `json:"disabled"`
	Status          string  `json:"status" validate:"omitempty,oneof=active inactive pending"`
}

// User Profile models - Modelos relacionados con el perfil del usuario
type UserProfile struct {
	ID          uint                   `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      uint                   `json:"user_id" gorm:"not null;uniqueIndex"`
	Avatar      string                 `json:"avatar,omitempty" gorm:"size:500"`
	Bio         string                 `json:"bio,omitempty" gorm:"type:text"`
	Website     string                 `json:"website,omitempty" gorm:"size:255"`
	Location    string                 `json:"location,omitempty" gorm:"size:100"`
	Birthday    *time.Time             `json:"birthday,omitempty"`
	Gender      string                 `json:"gender,omitempty" gorm:"size:20"`
	Phone       string                 `json:"phone,omitempty" gorm:"size:20"`
	Preferences string                 `json:"preferences,omitempty" gorm:"type:jsonb"` // JSON en PostgreSQL
	Privacy     string                 `json:"privacy,omitempty" gorm:"type:jsonb"`     // JSON en PostgreSQL
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relación con User
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// User Settings models - Modelos relacionados con configuraciones del usuario
type UserSettings struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        uint      `json:"user_id" gorm:"not null;uniqueIndex"`
	Language      string    `json:"language" gorm:"size:10;default:'en'"`
	Timezone      string    `json:"timezone" gorm:"size:50;default:'UTC'"`
	Theme         string    `json:"theme" gorm:"size:20;default:'light'"`
	Notifications string    `json:"notifications" gorm:"type:jsonb"` // JSON en PostgreSQL
	Privacy       string    `json:"privacy" gorm:"type:jsonb"`       // JSON en PostgreSQL
	Security      string    `json:"security" gorm:"type:jsonb"`      // JSON en PostgreSQL
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relación con User
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// User Statistics models - Modelos relacionados con estadísticas del usuario
type UserStats struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       uint       `json:"user_id" gorm:"not null;uniqueIndex"`
	LoginCount   int        `json:"login_count" gorm:"default:0"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	ProfileViews int        `json:"profile_views" gorm:"default:0"`
	AccountAge   int        `json:"account_age_days" gorm:"default:0"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relación con User
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
