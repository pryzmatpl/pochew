package storage

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/readitlater/backend/internal/config"
	"github.com/readitlater/backend/internal/encryption"
)

// Service handles file storage operations
type Service struct {
	config     *config.Config
	logger     *logrus.Logger
	encryption *encryption.Service
}

// ArticleContent represents the structure of stored article content
type ArticleContent struct {
	ID          string            `json:"id"`
	UserID      string            `json:"user_id"`
	Title       string            `json:"title"`
	URL         string            `json:"url"`
	Content     string            `json:"content"`
	Summary     string            `json:"summary"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	IsEncrypted bool              `json:"is_encrypted"`
}

// NewService creates a new storage service
func NewService(cfg *config.Config, logger *logrus.Logger) *Service {
	encService := encryption.NewService(100000) // 100k iterations for PBKDF2
	return &Service{
		config:     cfg,
		logger:     logger,
		encryption: encService,
	}
}

// SaveContent saves content to local storage
func (s *Service) SaveContent(userID, articleID, content string) (string, error) {
	// Parse content as ArticleContent
	var articleContent ArticleContent
	if err := json.Unmarshal([]byte(content), &articleContent); err != nil {
		return "", fmt.Errorf("failed to parse article content: %w", err)
	}

	// Ensure user directory exists
	userDir := s.getUserStoragePath(userID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create user directory: %w", err)
	}

	// Set metadata
	articleContent.UserID = userID
	articleContent.ID = articleID
	articleContent.UpdatedAt = time.Now()
	if articleContent.CreatedAt.IsZero() {
		articleContent.CreatedAt = time.Now()
	}

	// Serialize content
	contentBytes, err := json.Marshal(articleContent)
	if err != nil {
		return "", fmt.Errorf("failed to serialize content: %w", err)
	}

	// Save to file
	filePath := filepath.Join(userDir, fmt.Sprintf("%s.json", articleID))
	if err := os.WriteFile(filePath, contentBytes, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"article_id": articleID,
		"file_path":  filePath,
	}).Info("Content saved successfully")

	return filePath, nil
}

// SaveEncryptedContent saves encrypted content to local storage
func (s *Service) SaveEncryptedContent(userID, articleID, content, userKey string) (string, error) {
	// Parse content as ArticleContent
	var articleContent ArticleContent
	if err := json.Unmarshal([]byte(content), &articleContent); err != nil {
		return "", fmt.Errorf("failed to parse article content: %w", err)
	}

	// Encrypt the content field
	encryptedContent, err := s.encryption.Encrypt(articleContent.Content, userKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt content: %w", err)
	}

	// Encrypt the summary if present
	if articleContent.Summary != "" {
		encryptedSummary, err := s.encryption.Encrypt(articleContent.Summary, userKey)
		if err != nil {
			return "", fmt.Errorf("failed to encrypt summary: %w", err)
		}
		articleContent.Summary = encryptedSummary
	}

	articleContent.Content = encryptedContent
	articleContent.IsEncrypted = true

	// Convert back to JSON and save
	contentBytes, err := json.Marshal(articleContent)
	if err != nil {
		return "", fmt.Errorf("failed to serialize encrypted content: %w", err)
	}

	return s.SaveContent(userID, articleID, string(contentBytes))
}

// GetContent retrieves content from storage
func (s *Service) GetContent(userID, articleID string) (string, error) {
	filePath := filepath.Join(s.getUserStoragePath(userID), fmt.Sprintf("%s.json", articleID))
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("article not found: %s", articleID)
	}

	// Read file
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"article_id": articleID,
		"file_path":  filePath,
	}).Info("Content retrieved successfully")

	return string(contentBytes), nil
}

// GetDecryptedContent retrieves and decrypts content from storage
func (s *Service) GetDecryptedContent(userID, articleID, userKey string) (string, error) {
	content, err := s.GetContent(userID, articleID)
	if err != nil {
		return "", err
	}

	var articleContent ArticleContent
	if err := json.Unmarshal([]byte(content), &articleContent); err != nil {
		return "", fmt.Errorf("failed to parse article content: %w", err)
	}

	// If not encrypted, return as is
	if !articleContent.IsEncrypted {
		return content, nil
	}

	// Decrypt content
	decryptedContent, err := s.encryption.Decrypt(articleContent.Content, userKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt content: %w", err)
	}
	articleContent.Content = decryptedContent

	// Decrypt summary if present
	if articleContent.Summary != "" {
		decryptedSummary, err := s.encryption.Decrypt(articleContent.Summary, userKey)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt summary: %w", err)
		}
		articleContent.Summary = decryptedSummary
	}

	articleContent.IsEncrypted = false

	// Convert back to JSON
	decryptedBytes, err := json.Marshal(articleContent)
	if err != nil {
		return "", fmt.Errorf("failed to serialize decrypted content: %w", err)
	}

	return string(decryptedBytes), nil
}

// DeleteContent deletes content from storage
func (s *Service) DeleteContent(userID, articleID string) error {
	filePath := filepath.Join(s.getUserStoragePath(userID), fmt.Sprintf("%s.json", articleID))
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("article not found: %s", articleID)
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"article_id": articleID,
		"file_path":  filePath,
	}).Info("Content deleted successfully")

	return nil
}

// ListUserContent lists all articles for a user
func (s *Service) ListUserContent(userID string) ([]ArticleContent, error) {
	userDir := s.getUserStoragePath(userID)
	
	// Check if user directory exists
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		return []ArticleContent{}, nil
	}

	var articles []ArticleContent
	
	err := filepath.WalkDir(userDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			s.logger.WithError(err).WithField("file_path", path).Warn("Failed to read article file")
			return nil
		}

		var article ArticleContent
		if err := json.Unmarshal(contentBytes, &article); err != nil {
			s.logger.WithError(err).WithField("file_path", path).Warn("Failed to parse article file")
			return nil
		}

		articles = append(articles, article)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list user content: %w", err)
	}

	return articles, nil
}

// BackupToCloud backs up content to cloud storage
func (s *Service) BackupToCloud(userID, articleID string) error {
	// Get content
	content, err := s.GetContent(userID, articleID)
	if err != nil {
		return fmt.Errorf("failed to get content for backup: %w", err)
	}

	// For now, this is a placeholder for cloud backup
	// In a real implementation, this would upload to AWS S3, Google Cloud Storage, etc.
	s.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"article_id": articleID,
		"content_size": len(content),
	}).Info("Cloud backup initiated (placeholder)")

	return nil
}

// RestoreFromCloud restores content from cloud storage
func (s *Service) RestoreFromCloud(userID, articleID string) error {
	// For now, this is a placeholder for cloud restore
	// In a real implementation, this would download from cloud storage
	s.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"article_id": articleID,
	}).Info("Cloud restore initiated (placeholder)")

	return nil
}

// GetStorageStats returns storage statistics for a user
func (s *Service) GetStorageStats(userID string) (map[string]interface{}, error) {
	userDir := s.getUserStoragePath(userID)
	
	stats := map[string]interface{}{
		"total_articles": 0,
		"total_size":     int64(0),
		"encrypted_articles": 0,
	}

	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		return stats, nil
	}

	err := filepath.WalkDir(userDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		stats["total_articles"] = stats["total_articles"].(int) + 1
		stats["total_size"] = stats["total_size"].(int64) + info.Size()

		// Check if encrypted
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		var article ArticleContent
		if err := json.Unmarshal(contentBytes, &article); err != nil {
			return nil
		}

		if article.IsEncrypted {
			stats["encrypted_articles"] = stats["encrypted_articles"].(int) + 1
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to calculate storage stats: %w", err)
	}

	return stats, nil
}

// getUserStoragePath returns the storage path for a user
func (s *Service) getUserStoragePath(userID string) string {
	return filepath.Join(s.config.LocalStoragePath, "users", userID)
} 