package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Article represents a saved webpage/article
type Article struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	URL             string     `json:"url" db:"url" validate:"required,url"`
	Title           string     `json:"title" db:"title"`
	Description     string     `json:"description" db:"description"`
	Content         string     `json:"content" db:"content"`
	ContentText     string     `json:"content_text" db:"content_text"` // Plain text version for search
	Author          string     `json:"author" db:"author"`
	SiteName        string     `json:"site_name" db:"site_name"`
	ImageURL        string     `json:"image_url" db:"image_url"`
	PublishedAt     *time.Time `json:"published_at" db:"published_at"`
	WordCount       int        `json:"word_count" db:"word_count"`
	ReadingTime     int        `json:"reading_time" db:"reading_time"` // in minutes
	Language        string     `json:"language" db:"language"`
	
	// Storage and encryption
	IsEncrypted     bool   `json:"is_encrypted" db:"is_encrypted"`
	StorageSize     int64  `json:"storage_size" db:"storage_size"`
	LocalPath       string `json:"local_path" db:"local_path"`
	CloudPath       string `json:"cloud_path" db:"cloud_path"`
	
	// Status and metadata
	Status          string    `json:"status" db:"status"` // pending, processed, failed
	IsRead          bool      `json:"is_read" db:"is_read"`
	IsFavorite      bool      `json:"is_favorite" db:"is_favorite"`
	IsArchived      bool      `json:"is_archived" db:"is_archived"`
	ReadAt          *time.Time `json:"read_at" db:"read_at"`
	
	// Tags and categories
	Tags            []string  `json:"tags" db:"tags"`
	Category        string    `json:"category" db:"category"`
	
	// Timestamps
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	
	// Capture metadata
	CaptureMethod   string    `json:"capture_method" db:"capture_method"` // extension, api, manual
	UserAgent       string    `json:"user_agent" db:"user_agent"`
	CapturedAt      time.Time `json:"captured_at" db:"captured_at"`
}

// ArticleCreate represents the data needed to create a new article
type ArticleCreate struct {
	URL           string            `json:"url" validate:"required,url"`
	Title         string            `json:"title,omitempty"`
	Description   string            `json:"description,omitempty"`
	Content       string            `json:"content,omitempty"`
	Summary       string            `json:"summary,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	Category      string            `json:"category,omitempty"`
	CaptureMethod string            `json:"capture_method,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	UserKey       string            `json:"user_key,omitempty"`
}

// ArticleUpdate represents the data that can be updated for an article
type ArticleUpdate struct {
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
	Category    *string   `json:"category,omitempty"`
	IsRead      *bool     `json:"is_read,omitempty"`
	IsFavorite  *bool     `json:"is_favorite,omitempty"`
	IsArchived  *bool     `json:"is_archived,omitempty"`
}

// ArticleResponse represents the article data returned in API responses
type ArticleResponse struct {
	ID              uuid.UUID  `json:"id"`
	URL             string     `json:"url"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Author          string     `json:"author"`
	SiteName        string     `json:"site_name"`
	ImageURL        string     `json:"image_url"`
	PublishedAt     *time.Time `json:"published_at"`
	WordCount       int        `json:"word_count"`
	ReadingTime     int        `json:"reading_time"`
	Language        string     `json:"language"`
	StorageSize     int64      `json:"storage_size"`
	Status          string     `json:"status"`
	IsRead          bool       `json:"is_read"`
	IsFavorite      bool       `json:"is_favorite"`
	IsArchived      bool       `json:"is_archived"`
	ReadAt          *time.Time `json:"read_at"`
	Tags            []string   `json:"tags"`
	Category        string     `json:"category"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CaptureMethod   string     `json:"capture_method"`
	CapturedAt      time.Time  `json:"captured_at"`
}

// ArticleContent represents the full content of an article
type ArticleContent struct {
	Article  Article           `json:"article"`
	ID       uuid.UUID         `json:"id"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
}

// ArticleFilter represents filters for querying articles
type ArticleFilter struct {
	UserID     uuid.UUID `json:"user_id"`
	Status     string    `json:"status,omitempty"`
	IsRead     *bool     `json:"is_read,omitempty"`
	IsFavorite *bool     `json:"is_favorite,omitempty"`
	IsArchived *bool     `json:"is_archived,omitempty"`
	Category   string    `json:"category,omitempty"`
	Tags       []string  `json:"tags,omitempty"`
	Search     string    `json:"search,omitempty"`
	DateFrom   *time.Time `json:"date_from,omitempty"`
	DateTo     *time.Time `json:"date_to,omitempty"`
	Limit      int       `json:"limit,omitempty"`
	Offset     int       `json:"offset,omitempty"`
	SortBy     string    `json:"sort_by,omitempty"` // created_at, updated_at, title, reading_time
	SortOrder  string    `json:"sort_order,omitempty"` // asc, desc
}

// ArticleStats represents statistics about user's articles
type ArticleStats struct {
	TotalArticles    int   `json:"total_articles"`
	ReadArticles     int   `json:"read_articles"`
	UnreadArticles   int   `json:"unread_articles"`
	FavoriteArticles int   `json:"favorite_articles"`
	ArchivedArticles int   `json:"archived_articles"`
	TotalWordCount   int   `json:"total_word_count"`
	TotalReadingTime int   `json:"total_reading_time"`
	StorageUsed      int64 `json:"storage_used"`
}

// ToResponse converts an Article to ArticleResponse
func (a *Article) ToResponse() *ArticleResponse {
	return &ArticleResponse{
		ID:            a.ID,
		URL:           a.URL,
		Title:         a.Title,
		Description:   a.Description,
		Author:        a.Author,
		SiteName:      a.SiteName,
		ImageURL:      a.ImageURL,
		PublishedAt:   a.PublishedAt,
		WordCount:     a.WordCount,
		ReadingTime:   a.ReadingTime,
		Language:      a.Language,
		StorageSize:   a.StorageSize,
		Status:        a.Status,
		IsRead:        a.IsRead,
		IsFavorite:    a.IsFavorite,
		IsArchived:    a.IsArchived,
		ReadAt:        a.ReadAt,
		Tags:          a.Tags,
		Category:      a.Category,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
		CaptureMethod: a.CaptureMethod,
		CapturedAt:    a.CapturedAt,
	}
}

// ToContent converts an Article to ArticleContent
func (a *Article) ToContent() *ArticleContent {
	return &ArticleContent{
		Article: *a,
		ID:      a.ID,
		Content: a.Content,
		Metadata: map[string]string{
			"title":       a.Title,
			"description": a.Description,
			"author":      a.Author,
			"site_name":   a.SiteName,
			"image_url":   a.ImageURL,
			"language":    a.Language,
			"word_count":  fmt.Sprintf("%d", a.WordCount),
			"reading_time": fmt.Sprintf("%d", a.ReadingTime),
			"status":      a.Status,
			"is_read":     fmt.Sprintf("%v", a.IsRead),
			"is_favorite": fmt.Sprintf("%v", a.IsFavorite),
			"is_archived": fmt.Sprintf("%v", a.IsArchived),
			"created_at":  a.CreatedAt.Format(time.RFC3339),
			"updated_at":  a.UpdatedAt.Format(time.RFC3339),
			"capture_method": a.CaptureMethod,
			"captured_at": a.CapturedAt.Format(time.RFC3339),
		},
	}
}

// MarkAsRead marks the article as read
func (a *Article) MarkAsRead() {
	a.IsRead = true
	now := time.Now()
	a.ReadAt = &now
	a.UpdatedAt = now
}

// ToggleFavorite toggles the favorite status
func (a *Article) ToggleFavorite() {
	a.IsFavorite = !a.IsFavorite
	a.UpdatedAt = time.Now()
}

// Archive archives the article
func (a *Article) Archive() {
	a.IsArchived = true
	a.UpdatedAt = time.Now()
}

// Unarchive unarchives the article
func (a *Article) Unarchive() {
	a.IsArchived = false
	a.UpdatedAt = time.Now()
}

// CalculateReadingTime calculates reading time based on word count
func (a *Article) CalculateReadingTime() {
	// Average reading speed is 200-250 words per minute
	// We'll use 225 as a middle ground
	if a.WordCount > 0 {
		a.ReadingTime = (a.WordCount + 224) / 225 // Round up
		if a.ReadingTime < 1 {
			a.ReadingTime = 1
		}
	}
}

// GetDomain extracts domain from URL
func (a *Article) GetDomain() string {
	if a.URL == "" {
		return ""
	}
	
	// Simple domain extraction
	start := 0
	if len(a.URL) > 8 && a.URL[:8] == "https://" {
		start = 8
	} else if len(a.URL) > 7 && a.URL[:7] == "http://" {
		start = 7
	}
	
	end := len(a.URL)
	for i := start; i < len(a.URL); i++ {
		if a.URL[i] == '/' || a.URL[i] == '?' || a.URL[i] == '#' {
			end = i
			break
		}
	}
	
	domain := a.URL[start:end]
	
	// Remove www. prefix
	if len(domain) > 4 && domain[:4] == "www." {
		domain = domain[4:]
	}
	
	return domain
} 