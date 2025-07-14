package models

import (
	"time"
)

// URL represents a shortened URL
type URL struct {
	ID          string    `json:"id" db:"id"`
	ShortCode   string    `json:"short_code" db:"short_code"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	Title       string    `json:"title,omitempty" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy   string    `json:"created_by,omitempty" db:"created_by"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	ClickCount  int64     `json:"click_count,omitempty"`
}

// CreateURLRequest represents the request to create a new URL
type CreateURLRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	CustomCode  string `json:"custom_code,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// CreateURLResponse represents the response when creating a URL
type CreateURLResponse struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// URLListResponse represents the response for listing URLs
type URLListResponse struct {
	URLs        []URL `json:"urls"`
	Total       int   `json:"total"`
	Page        int   `json:"page"`
	PerPage     int   `json:"per_page"`
	TotalPages  int   `json:"total_pages"`
}

// URLStats represents statistics for a URL
type URLStats struct {
	ShortCode     string    `json:"short_code"`
	OriginalURL   string    `json:"original_url"`
	TotalClicks   int64     `json:"total_clicks"`
	UniqueClicks  int64     `json:"unique_clicks"`
	CreatedAt     time.Time `json:"created_at"`
	LastClickedAt *time.Time `json:"last_clicked_at,omitempty"`
}

// IsExpired checks if the URL has expired
func (u *URL) IsExpired() bool {
	if u.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*u.ExpiresAt)
}

// IsActive checks if the URL is active and not expired
func (u *URL) IsActive() bool {
	return u.IsActive && !u.IsExpired()
} 