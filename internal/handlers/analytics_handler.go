package handlers

import (
	"linksprint/internal/database"
	"linksprint/internal/models"
	"linksprint/internal/redis"
	"linksprint/internal/services"

	"github.com/gofiber/fiber/v2"
)

// AnalyticsHandler handles analytics-related HTTP requests
type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(db *database.DB, redis *redis.Client) *AnalyticsHandler {
	analyticsService := services.NewAnalyticsService(db, redis)
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetAnalytics handles GET /api/v1/analytics/:shortCode
func (h *AnalyticsHandler) GetAnalytics(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	analytics, err := h.analyticsService.GetAnalytics(c.Context(), shortCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(analytics)
}

// GetGlobalAnalytics handles GET /api/v1/analytics/global
func (h *AnalyticsHandler) GetGlobalAnalytics(c *fiber.Ctx) error {
	analytics, err := h.analyticsService.GetGlobalAnalytics(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(analytics)
}

// TrackClick handles POST /api/v1/analytics/track
func (h *AnalyticsHandler) TrackClick(c *fiber.Ctx) error {
	var req models.AnalyticsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.ShortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "short_code is required",
		})
	}

	// Track the click
	err := h.analyticsService.TrackClick(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":    "Click tracked successfully",
		"short_code": req.ShortCode,
	})
}
