package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/readitlater/backend/internal/config"
	"github.com/readitlater/backend/internal/models"
)

// Service handles authentication operations
type Service struct {
	config *config.Config
}

// NewService creates a new authentication service
func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

// HashPassword hashes a password using bcrypt
func (s *Service) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.PasswordSaltRounds)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword checks if a password matches the hash
func (s *Service) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a JWT token for a user
func (s *Service) GenerateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24 hours

	claims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "readitlater",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token and returns the claims
func (s *Service) ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired")
	}

	return claims, nil
}

// GenerateRefreshToken generates a secure random refresh token
func (s *Service) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateSessionToken generates a secure session token
func (s *Service) GenerateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// HashSessionToken hashes a session token for storage
func (s *Service) HashSessionToken(token string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash session token: %w", err)
	}
	return string(bytes), nil
}

// CheckSessionToken checks if a session token matches the hash
func (s *Service) CheckSessionToken(token, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}

// ValidatePassword validates password strength
func (s *Service) ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case isSpecialChar(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// ValidateEmail validates email format
func (s *Service) ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Simple email validation
	atIndex := -1
	dotIndex := -1
	
	for i, char := range email {
		if char == '@' {
			if atIndex != -1 {
				return fmt.Errorf("invalid email format: multiple @ symbols")
			}
			atIndex = i
		} else if char == '.' && atIndex != -1 {
			dotIndex = i
		}
	}

	if atIndex == -1 {
		return fmt.Errorf("invalid email format: missing @ symbol")
	}
	if atIndex == 0 {
		return fmt.Errorf("invalid email format: @ symbol at beginning")
	}
	if atIndex == len(email)-1 {
		return fmt.Errorf("invalid email format: @ symbol at end")
	}
	if dotIndex == -1 || dotIndex <= atIndex {
		return fmt.Errorf("invalid email format: missing or misplaced domain")
	}
	if dotIndex == len(email)-1 {
		return fmt.Errorf("invalid email format: domain ends with dot")
	}

	return nil
}

// ValidateUsername validates username format
func (s *Service) ValidateUsername(username string) error {
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if len(username) > 50 {
		return fmt.Errorf("username must be at most 50 characters long")
	}

	// Check if username contains only allowed characters
	for _, char := range username {
		if !isAlphanumeric(char) && char != '_' && char != '-' {
			return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
		}
	}

	// Username cannot start or end with special characters
	if username[0] == '_' || username[0] == '-' {
		return fmt.Errorf("username cannot start with underscore or hyphen")
	}
	if username[len(username)-1] == '_' || username[len(username)-1] == '-' {
		return fmt.Errorf("username cannot end with underscore or hyphen")
	}

	return nil
}

// GenerateUserEncryptionKey generates a unique encryption key for a user
func (s *Service) GenerateUserEncryptionKey() (string, error) {
	bytes := make([]byte, 32) // 256-bit key
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate encryption key: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// Helper functions

func isSpecialChar(char rune) bool {
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	for _, special := range specialChars {
		if char == special {
			return true
		}
	}
	return false
}

func isAlphanumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		   (char >= 'A' && char <= 'Z') ||
		   (char >= '0' && char <= '9')
} 