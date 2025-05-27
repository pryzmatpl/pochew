package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Port        string
	Environment string
	Debug       bool

	// Database configuration
	DatabaseURL string
	RedisURL    string

	// Authentication configuration
	JWTSecret           string
	EncryptionKey       string
	PasswordSaltRounds  int
	SessionSecret       string

	// Storage configuration
	LocalStoragePath string
	MaxStorageSize   string
	CleanupInterval  time.Duration

	// Cloud backup configuration
	EnableCloudBackup     bool
	CloudStorageProvider  string
	CloudStorageBucket    string
	CloudStorageRegion    string
	CloudAccessKey        string
	CloudSecretKey        string

	// Security configuration
	EncryptionAlgorithm       string
	KeyDerivationIterations   int
	AllowedOrigins           []string

	// Content capture configuration
	MaxContentSize       string
	AllowedContentTypes  []string
	CaptureTimeout       time.Duration
	UserAgent           string

	// Rate limiting
	RateLimitWindow      time.Duration
	RateLimitMaxRequests int

	// Logging
	LogLevel  string
	LogFormat string
	LogFile   string

	// Frontend URL
	FrontendURL string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Don't fail if .env doesn't exist
		fmt.Println("No .env file found, using environment variables")
	}

	config := &Config{
		Port:        getEnv("BACKEND_PORT", "8080"),
		Environment: getEnv("GO_ENV", "development"),
		Debug:       getEnvBool("DEBUG", true),

		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/readitlater?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),

		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		EncryptionKey:      getEnv("ENCRYPTION_KEY", "your-32-character-encryption-key"),
		PasswordSaltRounds: getEnvInt("PASSWORD_SALT_ROUNDS", 12),
		SessionSecret:      getEnv("SESSION_SECRET", "your-session-secret-key"),

		LocalStoragePath: getEnv("LOCAL_STORAGE_PATH", "./data/storage"),
		MaxStorageSize:   getEnv("MAX_STORAGE_SIZE", "10GB"),
		CleanupInterval:  getEnvDuration("CLEANUP_INTERVAL", 24*time.Hour),

		EnableCloudBackup:    getEnvBool("ENABLE_CLOUD_BACKUP", false),
		CloudStorageProvider: getEnv("CLOUD_STORAGE_PROVIDER", "local"),
		CloudStorageBucket:   getEnv("CLOUD_STORAGE_BUCKET", "readitlater-backup"),
		CloudStorageRegion:   getEnv("CLOUD_STORAGE_REGION", "us-east-1"),
		CloudAccessKey:       getEnv("CLOUD_ACCESS_KEY", ""),
		CloudSecretKey:       getEnv("CLOUD_SECRET_KEY", ""),

		EncryptionAlgorithm:     getEnv("ENCRYPTION_ALGORITHM", "AES-256-GCM"),
		KeyDerivationIterations: getEnvInt("KEY_DERIVATION_ITERATIONS", 100000),
		AllowedOrigins:         getEnvSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),

		MaxContentSize:      getEnv("MAX_CONTENT_SIZE", "50MB"),
		AllowedContentTypes: getEnvSlice("ALLOWED_CONTENT_TYPES", []string{"text/html", "application/pdf", "text/plain"}),
		CaptureTimeout:      getEnvDuration("CAPTURE_TIMEOUT", 30*time.Second),
		UserAgent:          getEnv("USER_AGENT", "ReadItLater/1.0"),

		RateLimitWindow:      getEnvDuration("RATE_LIMIT_WINDOW", 15*time.Minute),
		RateLimitMaxRequests: getEnvInt("RATE_LIMIT_MAX_REQUESTS", 100),

		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
		LogFile:   getEnv("LOG_FILE", "./logs/app.log"),

		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
	}

	// Validate required configuration
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// validate checks if all required configuration is present
func (c *Config) validate() error {
	if c.JWTSecret == "your-super-secret-jwt-key" && c.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}

	if c.EncryptionKey == "your-32-character-encryption-key" && c.Environment == "production" {
		return fmt.Errorf("ENCRYPTION_KEY must be set in production")
	}

	if len(c.EncryptionKey) < 32 {
		return fmt.Errorf("ENCRYPTION_KEY must be at least 32 characters long")
	}

	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	return nil
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated parsing
		result := []string{}
		for _, item := range splitAndTrim(value, ",") {
			if item != "" {
				result = append(result, item)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

func splitAndTrim(s, sep string) []string {
	parts := []string{}
	for _, part := range splitString(s, sep) {
		trimmed := trimString(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	
	parts := []string{}
	start := 0
	
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			parts = append(parts, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	parts = append(parts, s[start:])
	
	return parts
}

func trimString(s string) string {
	start := 0
	end := len(s)
	
	// Trim leading whitespace
	for start < end && isWhitespace(s[start]) {
		start++
	}
	
	// Trim trailing whitespace
	for end > start && isWhitespace(s[end-1]) {
		end--
	}
	
	return s[start:end]
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
} 