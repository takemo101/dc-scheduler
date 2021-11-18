package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// SecurityKey is struct
type SecurityKey struct {
	logger core.Logger
	app    core.Application
	value  support.ContextValue
}

// NewSecurityKey is create middleware
func NewSecurityKey(
	logger core.Logger,
	app core.Application,
	value support.ContextValue,
) SecurityKey {
	return SecurityKey{
		logger: logger,
		app:    app,
		value:  value,
	}
}

// Setup security key middleware
func (m SecurityKey) Setup() {
	m.logger.Info("setup security key")
	m.app.App.Use(m.CreateHandler("key", []string{}))
}

// CreateHandler is create middleware handler
func (m SecurityKey) CreateHandler(key string, keys []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		response := m.value.GetResponseHelper(c)
		token := c.Query(key, "")

		for _, t := range keys {
			if token == t {
				return c.Next()
			}
		}

		return response.JsonError(c, errors.New("keys do not match"), fiber.StatusUnauthorized)
	}
}
