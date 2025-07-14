package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"linksprint/internal/database"
	"linksprint/internal/models"
	"linksprint/internal/redis"
)

// URLService handles URL shortening business logic
type URLService struct {
	db    *database.DB
	redis *redis.Client
}

// NewURLService creates a new URL service
func NewURLService(db *database.DB, redis *redis.Client) *URLService {
	return &URLService{
		db:    db,
		redis: redis,
	}
}

// CreateShortURL creates a new shortened URL
func (s *URLService) CreateShortURL(ctx context.Context, req *models.CreateURLRequest) (*models.CreateURLResponse, error) {
	// Validate original URL
	if err := s.validateURL(req.OriginalURL); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Generate short code
	shortCode := req.CustomCode
	if shortCode == "" {
		var err error
		shortCode, err = s.generateShortCode()
		if err != nil {
			return nil, fmt.Errorf("failed to generate short code: %w", err)
		}
	} else {
		// Validate custom code
		if err := s.validateShortCode(shortCode); err != nil {
			return nil, fmt.Errorf("invalid custom code: %w", err)
		}
	}

	// Check if short code already exists
	exists, err := s.shortCodeExists(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("failed to check short code: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("short code already exists")
	}

	// Create URL in database
	urlID, err := s.createURLInDB(ctx, req, shortCode)
	if err != nil {
		return nil, fmt.Errorf("failed to create URL in database: %w", err)
	}

	// Cache the URL in Redis
	if err := s.redis.SetURL(ctx, shortCode, req.OriginalURL); err != nil {
		log.Printf("Warning: failed to cache URL in Redis: %v", err)
	}

	// Build short URL
	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortCode)

	return &models.CreateURLResponse{
		ShortCode:   shortCode,
		OriginalURL: req.OriginalURL,
		ShortURL:    shortURL,
		CreatedAt:   time.Now(),
	}, nil
}

// GetOriginalURL retrieves the original URL for a short code
func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	// Try to get from cache first
	if originalURL, err := s.redis.GetURL(ctx, shortCode); err == nil {
		// Increment click count in Redis
		s.redis.IncrementClickCount(ctx, shortCode)
		return originalURL, nil
	}

	// If not in cache, get from database
	originalURL, err := s.getURLFromDB(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("URL not found: %w", err)
	}

	// Cache the URL for future requests
	if err := s.redis.SetURL(ctx, shortCode, originalURL); err != nil {
		log.Printf("Warning: failed to cache URL in Redis: %v", err)
	}

	// Increment click count
	s.redis.IncrementClickCount(ctx, shortCode)

	return originalURL, nil
}

// GetURLStats gets statistics for a URL
func (s *URLService) GetURLStats(ctx context.Context, shortCode string) (*models.URLStats, error) {
	// Get URL from database
	url, err := s.getURLByShortCode(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("URL not found: %w", err)
	}

	// Get click count from Redis
	clickCount, err := s.redis.GetClickCount(ctx, shortCode)
	if err != nil {
		clickCount = 0 // Default to 0 if not found
	}

	return &models.URLStats{
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
		TotalClicks: clickCount,
		CreatedAt:   url.CreatedAt,
	}, nil
}

// ListURLs lists all URLs with pagination
func (s *URLService) ListURLs(ctx context.Context, page, perPage int) (*models.URLListResponse, error) {
	offset := (page - 1) * perPage

	// Get total count
	var total int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM urls WHERE is_active = true").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get URLs
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, short_code, original_url, title, description, created_at, updated_at, created_by, is_active, expires_at
		FROM urls 
		WHERE is_active = true 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query URLs: %w", err)
	}
	defer rows.Close()

	var urls []models.URL
	for rows.Next() {
		var url models.URL
		err := rows.Scan(
			&url.ID,
			&url.ShortCode,
			&url.OriginalURL,
			&url.Title,
			&url.Description,
			&url.CreatedAt,
			&url.UpdatedAt,
			&url.CreatedBy,
			&url.IsActive,
			&url.ExpiresAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan URL: %w", err)
		}
		urls = append(urls, url)
	}

	totalPages := (total + perPage - 1) / perPage

	return &models.URLListResponse{
		URLs:       urls,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// Helper methods

func (s *URLService) validateURL(originalURL string) error {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("URL must have scheme and host")
	}
	return nil
}

func (s *URLService) validateShortCode(shortCode string) error {
	if len(shortCode) < 3 || len(shortCode) > 10 {
		return fmt.Errorf("short code must be between 3 and 10 characters")
	}
	// Check for valid characters (alphanumeric and hyphens)
	for _, char := range shortCode {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-') {
			return fmt.Errorf("short code can only contain letters, numbers, and hyphens")
		}
	}
	return nil
}

func (s *URLService) generateShortCode() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}

	return string(bytes), nil
}

func (s *URLService) shortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)", shortCode).Scan(&exists)
	return exists, err
}

func (s *URLService) createURLInDB(ctx context.Context, req *models.CreateURLRequest, shortCode string) (string, error) {
	var urlID string
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO urls (short_code, original_url, title, description, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, shortCode, req.OriginalURL, req.Title, req.Description, req.ExpiresAt).Scan(&urlID)
	return urlID, err
}

func (s *URLService) getURLFromDB(ctx context.Context, shortCode string) (string, error) {
	var originalURL string
	err := s.db.QueryRowContext(ctx, `
		SELECT original_url FROM urls 
		WHERE short_code = $1 AND is_active = true 
		AND (expires_at IS NULL OR expires_at > NOW())
	`, shortCode).Scan(&originalURL)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("URL not found or expired")
	}
	return originalURL, err
}

func (s *URLService) getURLByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	var url models.URL
	err := s.db.QueryRowContext(ctx, `
		SELECT id, short_code, original_url, title, description, created_at, updated_at, created_by, is_active, expires_at
		FROM urls WHERE short_code = $1 AND is_active = true
	`, shortCode).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.Title,
		&url.Description,
		&url.CreatedAt,
		&url.UpdatedAt,
		&url.CreatedBy,
		&url.IsActive,
		&url.ExpiresAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("URL not found")
	}
	return &url, err
} 