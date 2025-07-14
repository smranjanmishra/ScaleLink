package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"linksprint/internal/database"
	"linksprint/internal/models"
	"linksprint/internal/redis"
)

// AnalyticsService handles analytics business logic
type AnalyticsService struct {
	db    *database.DB
	redis *redis.Client
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(db *database.DB, redis *redis.Client) *AnalyticsService {
	return &AnalyticsService{
		db:    db,
		redis: redis,
	}
}

// TrackClick tracks a click event
func (s *AnalyticsService) TrackClick(ctx context.Context, req *models.AnalyticsRequest) error {
	// Get URL ID from short code
	urlID, err := s.getURLIDByShortCode(ctx, req.ShortCode)
	if err != nil {
		return fmt.Errorf("URL not found: %w", err)
	}

	// Insert analytics record
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO analytics (url_id, short_code, ip_address, user_agent, referer, country, city)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, urlID, req.ShortCode, req.IPAddress, req.UserAgent, req.Referer, req.Country, req.City)
	if err != nil {
		return fmt.Errorf("failed to track click: %w", err)
	}

	// Increment click count in Redis
	s.redis.IncrementClickCount(ctx, req.ShortCode)

	return nil
}

// GetAnalytics gets analytics for a specific URL
func (s *AnalyticsService) GetAnalytics(ctx context.Context, shortCode string) (*models.AnalyticsSummary, error) {
	// Get URL info
	url, err := s.getURLByShortCode(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("URL not found: %w", err)
	}

	// Get total clicks
	var totalClicks int64
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM analytics WHERE short_code = $1
	`, shortCode).Scan(&totalClicks)
	if err != nil {
		return nil, fmt.Errorf("failed to get total clicks: %w", err)
	}

	// Get unique clicks (by IP)
	var uniqueClicks int64
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT ip_address) FROM analytics WHERE short_code = $1
	`, shortCode).Scan(&uniqueClicks)
	if err != nil {
		return nil, fmt.Errorf("failed to get unique clicks: %w", err)
	}

	// Get last clicked time
	var lastClickedAt *time.Time
	err = s.db.QueryRowContext(ctx, `
		SELECT MAX(clicked_at) FROM analytics WHERE short_code = $1
	`, shortCode).Scan(&lastClickedAt)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get last clicked time: %w", err)
	}

	// Get top countries
	topCountries, err := s.getTopCountries(ctx, shortCode)
	if err != nil {
		log.Printf("Warning: failed to get top countries: %v", err)
	}

	// Get top cities
	topCities, err := s.getTopCities(ctx, shortCode)
	if err != nil {
		log.Printf("Warning: failed to get top cities: %v", err)
	}

	// Get top referers
	topReferers, err := s.getTopReferers(ctx, shortCode)
	if err != nil {
		log.Printf("Warning: failed to get top referers: %v", err)
	}

	// Get click trend (last 7 days)
	clickTrend, err := s.getClickTrend(ctx, shortCode, 7)
	if err != nil {
		log.Printf("Warning: failed to get click trend: %v", err)
	}

	return &models.AnalyticsSummary{
		ShortCode:     shortCode,
		TotalClicks:   totalClicks,
		UniqueClicks:  uniqueClicks,
		TopCountries:  topCountries,
		TopCities:     topCities,
		TopReferers:   topReferers,
		ClickTrend:    clickTrend,
		LastClickedAt: lastClickedAt,
	}, nil
}

// GetGlobalAnalytics gets global analytics
func (s *AnalyticsService) GetGlobalAnalytics(ctx context.Context) (*models.GlobalAnalytics, error) {
	// Get total URLs
	var totalURLs int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM urls WHERE is_active = true").Scan(&totalURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to get total URLs: %w", err)
	}

	// Get total clicks
	var totalClicks int64
	err = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM analytics").Scan(&totalClicks)
	if err != nil {
		return nil, fmt.Errorf("failed to get total clicks: %w", err)
	}

	// Get active URLs (created in last 30 days)
	var activeURLs int64
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM urls 
		WHERE is_active = true AND created_at >= NOW() - INTERVAL '30 days'
	`).Scan(&activeURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to get active URLs: %w", err)
	}

	// Get today's clicks
	var todayClicks int64
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM analytics 
		WHERE DATE(clicked_at) = CURRENT_DATE
	`).Scan(&todayClicks)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's clicks: %w", err)
	}

	// Get this week's clicks
	var thisWeekClicks int64
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM analytics 
		WHERE clicked_at >= DATE_TRUNC('week', CURRENT_DATE)
	`).Scan(&thisWeekClicks)
	if err != nil {
		return nil, fmt.Errorf("failed to get this week's clicks: %w", err)
	}

	// Get this month's clicks
	var thisMonthClicks int64
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM analytics 
		WHERE clicked_at >= DATE_TRUNC('month', CURRENT_DATE)
	`).Scan(&thisMonthClicks)
	if err != nil {
		return nil, fmt.Errorf("failed to get this month's clicks: %w", err)
	}

	return &models.GlobalAnalytics{
		TotalURLs:       totalURLs,
		TotalClicks:     totalClicks,
		ActiveURLs:      activeURLs,
		TodayClicks:     todayClicks,
		ThisWeekClicks:  thisWeekClicks,
		ThisMonthClicks: thisMonthClicks,
	}, nil
}

// Helper methods

func (s *AnalyticsService) getURLIDByShortCode(ctx context.Context, shortCode string) (string, error) {
	var urlID string
	err := s.db.QueryRowContext(ctx, `
		SELECT id FROM urls WHERE short_code = $1 AND is_active = true
	`, shortCode).Scan(&urlID)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("URL not found")
	}
	return urlID, err
}

func (s *AnalyticsService) getURLByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
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

func (s *AnalyticsService) getTopCountries(ctx context.Context, shortCode string) ([]models.Country, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT country, COUNT(*) as count 
		FROM analytics 
		WHERE short_code = $1 AND country IS NOT NULL AND country != ''
		GROUP BY country 
		ORDER BY count DESC 
		LIMIT 5
	`, shortCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []models.Country
	for rows.Next() {
		var country models.Country
		err := rows.Scan(&country.Name, &country.Count)
		if err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}
	return countries, nil
}

func (s *AnalyticsService) getTopCities(ctx context.Context, shortCode string) ([]models.City, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT city, COUNT(*) as count 
		FROM analytics 
		WHERE short_code = $1 AND city IS NOT NULL AND city != ''
		GROUP BY city 
		ORDER BY count DESC 
		LIMIT 5
	`, shortCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []models.City
	for rows.Next() {
		var city models.City
		err := rows.Scan(&city.Name, &city.Count)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}
	return cities, nil
}

func (s *AnalyticsService) getTopReferers(ctx context.Context, shortCode string) ([]models.Referer, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT referer, COUNT(*) as count 
		FROM analytics 
		WHERE short_code = $1 AND referer IS NOT NULL AND referer != ''
		GROUP BY referer 
		ORDER BY count DESC 
		LIMIT 5
	`, shortCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var referers []models.Referer
	for rows.Next() {
		var referer models.Referer
		err := rows.Scan(&referer.URL, &referer.Count)
		if err != nil {
			return nil, err
		}
		referers = append(referers, referer)
	}
	return referers, nil
}

func (s *AnalyticsService) getClickTrend(ctx context.Context, shortCode string, days int) ([]models.ClickTrend, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT DATE(clicked_at) as date, COUNT(*) as count 
		FROM analytics 
		WHERE short_code = $1 AND clicked_at >= CURRENT_DATE - INTERVAL '$2 days'
		GROUP BY DATE(clicked_at) 
		ORDER BY date
	`, shortCode, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []models.ClickTrend
	for rows.Next() {
		var trend models.ClickTrend
		err := rows.Scan(&trend.Date, &trend.Count)
		if err != nil {
			return nil, err
		}
		trends = append(trends, trend)
	}
	return trends, nil
}
