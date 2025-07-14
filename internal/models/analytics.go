package models

import (
	"time"
)

// Analytics represents a click event
type Analytics struct {
	ID        string    `json:"id" db:"id"`
	URLID     string    `json:"url_id" db:"url_id"`
	ShortCode string    `json:"short_code" db:"short_code"`
	IPAddress string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent string    `json:"user_agent,omitempty" db:"user_agent"`
	Referer   string    `json:"referer,omitempty" db:"referer"`
	Country   string    `json:"country,omitempty" db:"country"`
	City      string    `json:"city,omitempty" db:"city"`
	ClickedAt time.Time `json:"clicked_at" db:"clicked_at"`
}

// AnalyticsSummary represents aggregated analytics data
type AnalyticsSummary struct {
	ShortCode     string    `json:"short_code"`
	TotalClicks   int64     `json:"total_clicks"`
	UniqueClicks  int64     `json:"unique_clicks"`
	TopCountries  []Country `json:"top_countries"`
	TopCities     []City    `json:"top_cities"`
	TopReferers   []Referer `json:"top_referers"`
	ClickTrend    []ClickTrend `json:"click_trend"`
	LastClickedAt *time.Time `json:"last_clicked_at,omitempty"`
}

// Country represents country analytics
type Country struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// City represents city analytics
type City struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// Referer represents referer analytics
type Referer struct {
	URL   string `json:"url"`
	Count int64  `json:"count"`
}

// ClickTrend represents click trends over time
type ClickTrend struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// GlobalAnalytics represents global statistics
type GlobalAnalytics struct {
	TotalURLs      int64 `json:"total_urls"`
	TotalClicks    int64 `json:"total_clicks"`
	ActiveURLs     int64 `json:"active_urls"`
	TodayClicks    int64 `json:"today_clicks"`
	ThisWeekClicks int64 `json:"this_week_clicks"`
	ThisMonthClicks int64 `json:"this_month_clicks"`
}

// AnalyticsRequest represents the request to track analytics
type AnalyticsRequest struct {
	ShortCode string `json:"short_code" validate:"required"`
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Referer   string `json:"referer,omitempty"`
	Country   string `json:"country,omitempty"`
	City      string `json:"city,omitempty"`
} 