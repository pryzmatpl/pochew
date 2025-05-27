package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// Service handles encryption and decryption operations
type Service struct {
	iterations int
}

// NewService creates a new encryption service
func NewService(iterations int) *Service {
	return &Service{
		iterations: iterations,
	}
}

// GenerateUserKey generates a unique encryption key for a user
func (s *Service) GenerateUserKey() (string, error) {
	key := make([]byte, 32) // 256 bits
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate user key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// DeriveKey derives an encryption key from a master key and salt
func (s *Service) DeriveKey(masterKey, salt string) ([]byte, error) {
	masterKeyBytes, err := base64.StdEncoding.DecodeString(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode master key: %w", err)
	}
	
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %w", err)
	}
	
	return pbkdf2.Key(masterKeyBytes, saltBytes, s.iterations, 32, sha256.New), nil
}

// Encrypt encrypts data using AES-256-GCM
func (s *Service) Encrypt(data, userKey string) (string, error) {
	// Generate a random salt for key derivation
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	
	// Derive encryption key
	key, err := s.DeriveKey(userKey, base64.StdEncoding.EncodeToString(salt))
	if err != nil {
		return "", fmt.Errorf("failed to derive key: %w", err)
	}
	
	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}
	
	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}
	
	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	
	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	
	// Prepend salt to ciphertext
	result := append(salt, ciphertext...)
	
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts data using AES-256-GCM
func (s *Service) Decrypt(encryptedData, userKey string) (string, error) {
	// Decode base64
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted data: %w", err)
	}
	
	if len(data) < 16 {
		return "", fmt.Errorf("encrypted data too short")
	}
	
	// Extract salt and ciphertext
	salt := data[:16]
	ciphertext := data[16:]
	
	// Derive encryption key
	key, err := s.DeriveKey(userKey, base64.StdEncoding.EncodeToString(salt))
	if err != nil {
		return "", fmt.Errorf("failed to derive key: %w", err)
	}
	
	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}
	
	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}
	
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	
	// Extract nonce and encrypted data
	nonce, encData := ciphertext[:nonceSize], ciphertext[nonceSize:]
	
	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, encData, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}
	
	return string(plaintext), nil
}

// EncryptFile encrypts file content for storage
func (s *Service) EncryptFile(content []byte, userKey string) ([]byte, error) {
	encrypted, err := s.Encrypt(string(content), userKey)
	if err != nil {
		return nil, err
	}
	return []byte(encrypted), nil
}

// DecryptFile decrypts file content from storage
func (s *Service) DecryptFile(encryptedContent []byte, userKey string) ([]byte, error) {
	decrypted, err := s.Decrypt(string(encryptedContent), userKey)
	if err != nil {
		return nil, err
	}
	return []byte(decrypted), nil
} 