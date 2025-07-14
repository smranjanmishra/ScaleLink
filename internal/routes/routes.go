package routes

import (
	"linksprint/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, urlHandler *handlers.URLHandler, analyticsHandler *handlers.AnalyticsHandler) {
	// API v1 group
	api := app.Group("/api/v1")

	// URL shortening endpoints
	urls := api.Group("/urls")
	urls.Post("/shorten", urlHandler.CreateShortURL)
	urls.Get("/", urlHandler.ListURLs)
	urls.Get("/:shortCode/stats", urlHandler.GetURLStats)
	urls.Delete("/:shortCode", urlHandler.DeleteURL)

	// Analytics endpoints
	analytics := api.Group("/analytics")
	analytics.Get("/:shortCode", analyticsHandler.GetAnalytics)
	analytics.Get("/global", analyticsHandler.GetGlobalAnalytics)
	analytics.Post("/track", analyticsHandler.TrackClick)

	// Redirect endpoint (must be last to avoid conflicts)
	app.Get("/:shortCode", urlHandler.RedirectToOriginal)

	// API documentation endpoint
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":        "LinkSprint API",
			"version":     "1.0.0",
			"description": "Distributed URL Shortener & Analytics API",
			"endpoints": fiber.Map{
				"urls": fiber.Map{
					"POST /api/v1/urls/shorten":         "Create a short URL",
					"GET /api/v1/urls":                  "List all URLs",
					"GET /api/v1/urls/:shortCode/stats": "Get URL statistics",
					"DELETE /api/v1/urls/:shortCode":    "Delete a URL",
				},
				"analytics": fiber.Map{
					"GET /api/v1/analytics/:shortCode": "Get analytics for a URL",
					"GET /api/v1/analytics/global":     "Get global analytics",
					"POST /api/v1/analytics/track":     "Track a click event",
				},
				"redirect": fiber.Map{
					"GET /:shortCode": "Redirect to original URL",
				},
			},
		})
	})
}
