package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler handles application errors
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default error
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Return JSON error response
	return c.Status(code).JSON(fiber.Map{
		"error":      message,
		"code":       code,
		"path":       c.Path(),
		"method":     c.Method(),
		"request_id": c.Get("X-Request-ID", "unknown"),
	})
}
