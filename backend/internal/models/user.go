package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	Email             string     `json:"email" db:"email" validate:"required,email"`
	Username          string     `json:"username" db:"username" validate:"required,min=3,max=50"`
	PasswordHash      string     `json:"-" db:"password_hash"`
	FirstName         string     `json:"first_name" db:"first_name" validate:"max=100"`
	LastName          string     `json:"last_name" db:"last_name" validate:"max=100"`
	IsActive          bool       `json:"is_active" db:"is_active"`
	IsEmailVerified   bool       `json:"is_email_verified" db:"is_email_verified"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at" db:"email_verified_at"`
	LastLoginAt       *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	
	// Cloud backup preferences
	EnableCloudBackup bool   `json:"enable_cloud_backup" db:"enable_cloud_backup"`
	EncryptionKey     string `json:"-" db:"encryption_key"` // User-specific encryption key
	
	// Storage statistics
	StorageUsed       int64 `json:"storage_used" db:"storage_used"`
	MaxStorageLimit   int64 `json:"max_storage_limit" db:"max_storage_limit"`
	
	// Preferences
	Theme             string `json:"theme" db:"theme"`
	Language          string `json:"language" db:"language"`
	Timezone          string `json:"timezone" db:"timezone"`
}

// UserRegistration represents the data needed to register a new user
type UserRegistration struct {
	Email           string `json:"email" validate:"required,email"`
	Username        string `json:"username" validate:"required,min=3,max=50"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	FirstName       string `json:"first_name" validate:"max=100"`
	LastName        string `json:"last_name" validate:"max=100"`
}

// UserLogin represents the data needed to login
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserUpdate represents the data that can be updated for a user
type UserUpdate struct {
	FirstName         *string `json:"first_name,omitempty" validate:"omitempty,max=100"`
	LastName          *string `json:"last_name,omitempty" validate:"omitempty,max=100"`
	EnableCloudBackup *bool   `json:"enable_cloud_backup,omitempty"`
	Theme             *string `json:"theme,omitempty" validate:"omitempty,oneof=light dark auto"`
	Language          *string `json:"language,omitempty" validate:"omitempty,len=2"`
	Timezone          *string `json:"timezone,omitempty"`
}

// PasswordChange represents the data needed to change password
type PasswordChange struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID                uuid.UUID  `json:"id"`
	Email             string     `json:"email"`
	Username          string     `json:"username"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	IsActive          bool       `json:"is_active"`
	IsEmailVerified   bool       `json:"is_email_verified"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at"`
	LastLoginAt       *time.Time `json:"last_login_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	EnableCloudBackup bool       `json:"enable_cloud_backup"`
	StorageUsed       int64      `json:"storage_used"`
	MaxStorageLimit   int64      `json:"max_storage_limit"`
	Theme             string     `json:"theme"`
	Language          string     `json:"language"`
	Timezone          string     `json:"timezone"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:                u.ID,
		Email:             u.Email,
		Username:          u.Username,
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		IsActive:          u.IsActive,
		IsEmailVerified:   u.IsEmailVerified,
		EmailVerifiedAt:   u.EmailVerifiedAt,
		LastLoginAt:       u.LastLoginAt,
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
		EnableCloudBackup: u.EnableCloudBackup,
		StorageUsed:       u.StorageUsed,
		MaxStorageLimit:   u.MaxStorageLimit,
		Theme:             u.Theme,
		Language:          u.Language,
		Timezone:          u.Timezone,
	}
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	if u.FirstName == "" {
		return u.LastName
	}
	if u.LastName == "" {
		return u.FirstName
	}
	return u.FirstName + " " + u.LastName
}

// IsStorageLimitExceeded checks if the user has exceeded their storage limit
func (u *User) IsStorageLimitExceeded() bool {
	return u.StorageUsed >= u.MaxStorageLimit
}

// GetStorageUsagePercentage returns the storage usage as a percentage
func (u *User) GetStorageUsagePercentage() float64 {
	if u.MaxStorageLimit == 0 {
		return 0
	}
	return float64(u.StorageUsed) / float64(u.MaxStorageLimit) * 100
} 