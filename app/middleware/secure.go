package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/core"
)

// Secure is struct
type Secure struct {
	logger core.Logger
	app    core.Application
}

// NewSecure is create middleware
func NewSecure(
	logger core.Logger,
	app core.Application,
) Secure {
	return Secure{
		logger: logger,
		app:    app,
	}
}

// Setup secure control middleware
func (m Secure) Setup() {
	m.logger.Info("setup secure")

	m.app.App.Use(func(c *fiber.Ctx) error {
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Download-Options", "noopen")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Strict-Transport-Security", "max-age=31536000")
		c.Set("X-DNS-Prefetch-Control", "off")

		return c.Next()
	})
}
