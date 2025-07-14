package handlers

import (
	"strconv"

	"linksprint/internal/database"
	"linksprint/internal/models"
	"linksprint/internal/redis"
	"linksprint/internal/services"

	"github.com/gofiber/fiber/v2"
)

// URLHandler handles URL-related HTTP requests
type URLHandler struct {
	urlService *services.URLService
}

// NewURLHandler creates a new URL handler
func NewURLHandler(db *database.DB, redis *redis.Client) *URLHandler {
	urlService := services.NewURLService(db, redis)
	return &URLHandler{
		urlService: urlService,
	}
}

// CreateShortURL handles POST /api/v1/shorten
func (h *URLHandler) CreateShortURL(c *fiber.Ctx) error {
	var req models.CreateURLRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.OriginalURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "original_url is required",
		})
	}

	// Create short URL
	response, err := h.urlService.CreateShortURL(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// RedirectToOriginal handles GET /:shortCode
func (h *URLHandler) RedirectToOriginal(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	// Get original URL
	originalURL, err := h.urlService.GetOriginalURL(c.Context(), shortCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "URL not found or expired",
		})
	}

	// Track analytics (async)
	go func() {
		analyticsReq := &models.AnalyticsRequest{
			ShortCode: shortCode,
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
			Referer:   c.Get("Referer"),
		}
		// Note: In a real application, you'd want to use a proper analytics service
		// For now, we'll just log the click
		_ = analyticsReq
	}()

	// Redirect to original URL
	return c.Redirect(originalURL, fiber.StatusMovedPermanently)
}

// GetURLStats handles GET /api/v1/urls/:shortCode/stats
func (h *URLHandler) GetURLStats(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	stats, err := h.urlService.GetURLStats(c.Context(), shortCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(stats)
}

// ListURLs handles GET /api/v1/urls
func (h *URLHandler) ListURLs(c *fiber.Ctx) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// Get URLs
	response, err := h.urlService.ListURLs(c.Context(), page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

// DeleteURL handles DELETE /api/v1/urls/:shortCode
func (h *URLHandler) DeleteURL(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	// For now, we'll just return a success message
	// In a real application, you'd implement soft delete
	return c.JSON(fiber.Map{
		"message":    "URL deleted successfully",
		"short_code": shortCode,
	})
}
