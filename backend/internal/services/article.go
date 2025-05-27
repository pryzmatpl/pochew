package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/readitlater/backend/internal/database"
	"github.com/readitlater/backend/internal/models"
	"github.com/readitlater/backend/internal/storage"
)

// ArticleService handles article-related operations
type ArticleService struct {
	db             *database.DB
	storageService *storage.Service
	logger         *logrus.Logger
}

// NewArticleService creates a new article service
func NewArticleService(db *database.DB, storageService *storage.Service, logger *logrus.Logger) *ArticleService {
	return &ArticleService{
		db:             db,
		storageService: storageService,
		logger:         logger,
	}
}

// CreateArticle creates a new article
func (s *ArticleService) CreateArticle(userID string, create *models.ArticleCreate) (*models.Article, error) {
	// Generate article ID
	articleID := uuid.New()
	
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	
	// Create article model
	article := &models.Article{
		ID:          articleID,
		UserID:      userUUID,
		Title:       create.Title,
		URL:         create.URL,
		Description: create.Description,
		Tags:        create.Tags,
		Category:    create.Category,
		IsRead:      false,
		IsFavorite:  false,
		IsArchived:  false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	query := `
		INSERT INTO articles (id, user_id, title, url, description, tags, category, is_read, is_favorite, is_archived, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	
	tagsJSON, _ := json.Marshal(create.Tags)
	_, err = s.db.Exec(query, article.ID, article.UserID, article.Title, article.URL, 
		article.Description, string(tagsJSON), article.Category, article.IsRead, article.IsFavorite, 
		article.IsArchived, article.CreatedAt, article.UpdatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create article in database: %w", err)
	}

	// Save content to storage if provided
	if create.Content != "" {
		storageContent := storage.ArticleContent{
			ID:        articleID.String(),
			UserID:    userID,
			Title:     create.Title,
			URL:       create.URL,
			Content:   create.Content,
			Summary:   create.Summary,
			Tags:      create.Tags,
			Metadata:  create.Metadata,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		}

		contentJSON, err := json.Marshal(storageContent)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize content: %w", err)
		}

		// Save encrypted if user has encryption key
		if create.UserKey != "" {
			_, err = s.storageService.SaveEncryptedContent(userID, articleID.String(), string(contentJSON), create.UserKey)
		} else {
			_, err = s.storageService.SaveContent(userID, articleID.String(), string(contentJSON))
		}

		if err != nil {
			// Rollback database entry
			s.db.Exec("DELETE FROM articles WHERE id = $1", articleID)
			return nil, fmt.Errorf("failed to save content to storage: %w", err)
		}
	}

	s.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"article_id": articleID.String(),
		"title":      article.Title,
	}).Info("Article created successfully")

	return article, nil
}

// GetArticle retrieves an article by ID
func (s *ArticleService) GetArticle(userID, articleID string) (*models.Article, error) {
	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	
	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		return nil, fmt.Errorf("invalid article ID: %w", err)
	}
	
	query := `
		SELECT id, user_id, title, url, description, tags, category, is_read, is_favorite, is_archived, created_at, updated_at
		FROM articles 
		WHERE id = $1 AND user_id = $2
	`
	
	var article models.Article
	var tagsJSON string
	
	err = s.db.QueryRow(query, articleUUID, userUUID).Scan(
		&article.ID, &article.UserID, &article.Title, &article.URL, &article.Description,
		&tagsJSON, &article.Category, &article.IsRead, &article.IsFavorite, &article.IsArchived,
		&article.CreatedAt, &article.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	// Parse tags
	if tagsJSON != "" {
		json.Unmarshal([]byte(tagsJSON), &article.Tags)
	}

	return &article, nil
}

// GetArticles retrieves articles with filtering
func (s *ArticleService) GetArticles(filter *models.ArticleFilter) ([]*models.Article, error) {
	query := `
		SELECT id, user_id, title, url, description, tags, category, is_read, is_favorite, is_archived, created_at, updated_at
		FROM articles 
		WHERE user_id = $1
	`
	args := []interface{}{filter.UserID}
	argIndex := 2

	// Apply filters
	if filter.IsRead != nil {
		query += fmt.Sprintf(" AND is_read = $%d", argIndex)
		args = append(args, *filter.IsRead)
		argIndex++
	}

	if filter.IsFavorite != nil {
		query += fmt.Sprintf(" AND is_favorite = $%d", argIndex)
		args = append(args, *filter.IsFavorite)
		argIndex++
	}

	if filter.IsArchived != nil {
		query += fmt.Sprintf(" AND is_archived = $%d", argIndex)
		args = append(args, *filter.IsArchived)
		argIndex++
	}

	if len(filter.Tags) > 0 {
		// Search for any of the provided tags
		for _, tag := range filter.Tags {
			query += fmt.Sprintf(" AND tags LIKE $%d", argIndex)
			args = append(args, "%"+tag+"%")
			argIndex++
		}
	}

	// Add ordering
	query += " ORDER BY created_at DESC"

	// Add pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query articles: %w", err)
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		var article models.Article
		var tagsJSON string

		err := rows.Scan(
			&article.ID, &article.UserID, &article.Title, &article.URL, &article.Description,
			&tagsJSON, &article.Category, &article.IsRead, &article.IsFavorite, &article.IsArchived,
			&article.CreatedAt, &article.UpdatedAt,
		)
		if err != nil {
			continue
		}

		// Parse tags
		if tagsJSON != "" {
			json.Unmarshal([]byte(tagsJSON), &article.Tags)
		}

		articles = append(articles, &article)
	}

	return articles, nil
}

// UpdateArticle updates an article
func (s *ArticleService) UpdateArticle(userID, articleID string, update *models.ArticleUpdate) (*models.Article, error) {
	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	
	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		return nil, fmt.Errorf("invalid article ID: %w", err)
	}
	
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if update.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *update.Title)
		argIndex++
	}

	if update.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *update.Description)
		argIndex++
	}

	if update.Category != nil {
		setParts = append(setParts, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, *update.Category)
		argIndex++
	}

	if update.Tags != nil {
		tagsJSON, _ := json.Marshal(*update.Tags)
		setParts = append(setParts, fmt.Sprintf("tags = $%d", argIndex))
		args = append(args, string(tagsJSON))
		argIndex++
	}

	if update.IsRead != nil {
		setParts = append(setParts, fmt.Sprintf("is_read = $%d", argIndex))
		args = append(args, *update.IsRead)
		argIndex++
	}

	if update.IsFavorite != nil {
		setParts = append(setParts, fmt.Sprintf("is_favorite = $%d", argIndex))
		args = append(args, *update.IsFavorite)
		argIndex++
	}

	if update.IsArchived != nil {
		setParts = append(setParts, fmt.Sprintf("is_archived = $%d", argIndex))
		args = append(args, *update.IsArchived)
		argIndex++
	}

	if len(setParts) == 0 {
		return s.GetArticle(userID, articleID)
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	query := fmt.Sprintf(`
		UPDATE articles 
		SET %s
		WHERE id = $%d AND user_id = $%d
	`, strings.Join(setParts, ", "), argIndex, argIndex+1)
	
	args = append(args, articleUUID, userUUID)

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	return s.GetArticle(userID, articleID)
}

// DeleteArticle deletes an article
func (s *ArticleService) DeleteArticle(userID, articleID string) error {
	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	
	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		return fmt.Errorf("invalid article ID: %w", err)
	}
	
	// Delete from storage first
	err = s.storageService.DeleteContent(userID, articleID)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to delete content from storage")
	}

	// Delete from database
	query := "DELETE FROM articles WHERE id = $1 AND user_id = $2"
	result, err := s.db.Exec(query, articleUUID, userUUID)
	if err != nil {
		return fmt.Errorf("failed to delete article from database: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	s.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"article_id": articleID,
	}).Info("Article deleted successfully")

	return nil
}

// GetArticleContent retrieves the full content of an article
func (s *ArticleService) GetArticleContent(userID, articleID string) (*models.ArticleContent, error) {
	// Get article metadata from database
	article, err := s.GetArticle(userID, articleID)
	if err != nil {
		return nil, err
	}

	// Get content from storage
	content, err := s.storageService.GetContent(userID, articleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content from storage: %w", err)
	}

	var storageContent storage.ArticleContent
	if err := json.Unmarshal([]byte(content), &storageContent); err != nil {
		return nil, fmt.Errorf("failed to parse stored content: %w", err)
	}

	return &models.ArticleContent{
		Article: *article,
		Content: storageContent.Content,
		Metadata: storageContent.Metadata,
	}, nil
}

// GetDecryptedArticleContent retrieves and decrypts the full content of an article
func (s *ArticleService) GetDecryptedArticleContent(userID, articleID, userKey string) (*models.ArticleContent, error) {
	// Get article metadata from database
	article, err := s.GetArticle(userID, articleID)
	if err != nil {
		return nil, err
	}

	// Get decrypted content from storage
	content, err := s.storageService.GetDecryptedContent(userID, articleID, userKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get decrypted content from storage: %w", err)
	}

	var storageContent storage.ArticleContent
	if err := json.Unmarshal([]byte(content), &storageContent); err != nil {
		return nil, fmt.Errorf("failed to parse stored content: %w", err)
	}

	return &models.ArticleContent{
		Article: *article,
		Content: storageContent.Content,
		Metadata: storageContent.Metadata,
	}, nil
}

// MarkAsRead marks an article as read
func (s *ArticleService) MarkAsRead(userID, articleID string) error {
	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	
	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		return fmt.Errorf("invalid article ID: %w", err)
	}
	
	query := "UPDATE articles SET is_read = true, updated_at = $1 WHERE id = $2 AND user_id = $3"
	_, err = s.db.Exec(query, time.Now(), articleUUID, userUUID)
	if err != nil {
		return fmt.Errorf("failed to mark article as read: %w", err)
	}
	return nil
}

// ToggleFavorite toggles the favorite status of an article
func (s *ArticleService) ToggleFavorite(userID, articleID string) error {
	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	
	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		return fmt.Errorf("invalid article ID: %w", err)
	}
	
	query := "UPDATE articles SET is_favorite = NOT is_favorite, updated_at = $1 WHERE id = $2 AND user_id = $3"
	_, err = s.db.Exec(query, time.Now(), articleUUID, userUUID)
	if err != nil {
		return fmt.Errorf("failed to toggle favorite status: %w", err)
	}
	return nil
}

// ToggleArchive toggles the archive status of an article
func (s *ArticleService) ToggleArchive(userID, articleID string) error {
	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	
	articleUUID, err := uuid.Parse(articleID)
	if err != nil {
		return fmt.Errorf("invalid article ID: %w", err)
	}
	
	query := "UPDATE articles SET is_archived = NOT is_archived, updated_at = $1 WHERE id = $2 AND user_id = $3"
	_, err = s.db.Exec(query, time.Now(), articleUUID, userUUID)
	if err != nil {
		return fmt.Errorf("failed to toggle archive status: %w", err)
	}
	return nil
}

// SearchArticles searches articles by text
func (s *ArticleService) SearchArticles(userID, query string, filter *models.ArticleFilter) ([]*models.Article, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	
	searchQuery := `
		SELECT id, user_id, title, url, description, tags, category, is_read, is_favorite, is_archived, created_at, updated_at
		FROM articles 
		WHERE user_id = $1 AND (
			title ILIKE $2 OR 
			description ILIKE $2 OR 
			url ILIKE $2 OR
			tags ILIKE $2
		)
	`
	args := []interface{}{userUUID, "%" + query + "%"}
	argIndex := 3

	// Apply additional filters
	if filter != nil {
		if filter.IsRead != nil {
			searchQuery += fmt.Sprintf(" AND is_read = $%d", argIndex)
			args = append(args, *filter.IsRead)
			argIndex++
		}

		if filter.IsFavorite != nil {
			searchQuery += fmt.Sprintf(" AND is_favorite = $%d", argIndex)
			args = append(args, *filter.IsFavorite)
			argIndex++
		}

		if filter.IsArchived != nil {
			searchQuery += fmt.Sprintf(" AND is_archived = $%d", argIndex)
			args = append(args, *filter.IsArchived)
			argIndex++
		}
	}

	searchQuery += " ORDER BY created_at DESC"

	rows, err := s.db.Query(searchQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search articles: %w", err)
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		var article models.Article
		var tagsJSON string

		err := rows.Scan(
			&article.ID, &article.UserID, &article.Title, &article.URL, &article.Description,
			&tagsJSON, &article.Category, &article.IsRead, &article.IsFavorite, &article.IsArchived,
			&article.CreatedAt, &article.UpdatedAt,
		)
		if err != nil {
			continue
		}

		// Parse tags
		if tagsJSON != "" {
			json.Unmarshal([]byte(tagsJSON), &article.Tags)
		}

		articles = append(articles, &article)
	}

	return articles, nil
}

// GetStats retrieves article statistics for a user
func (s *ArticleService) GetStats(userID string) (*models.ArticleStats, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN is_read = true THEN 1 END) as read,
			COUNT(CASE WHEN is_favorite = true THEN 1 END) as favorites,
			COUNT(CASE WHEN is_archived = true THEN 1 END) as archived
		FROM articles 
		WHERE user_id = $1
	`

	var stats models.ArticleStats
	err = s.db.QueryRow(query, userUUID).Scan(
		&stats.TotalArticles, &stats.ReadArticles, &stats.FavoriteArticles, &stats.ArchivedArticles,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get article stats: %w", err)
	}

	stats.UnreadArticles = stats.TotalArticles - stats.ReadArticles

	// Get storage stats
	storageStats, err := s.storageService.GetStorageStats(userID)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to get storage stats")
	} else {
		stats.StorageUsed = storageStats["total_size"].(int64)
	}

	return &stats, nil
} 