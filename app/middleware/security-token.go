package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// SecurityToken is struct
type SecurityToken struct {
	logger core.Logger
	app    core.Application
	value  support.ContextValue
}

// NewSecurityToken is create middleware
func NewSecurityToken(
	logger core.Logger,
	app core.Application,
	value support.ContextValue,
) SecurityToken {
	return SecurityToken{
		logger: logger,
		app:    app,
		value:  value,
	}
}

// Setup security token middleware
func (m SecurityToken) Setup() {
	m.logger.Info("setup security token")
	m.app.App.Use(m.CreateHandler([]string{}))
}

// CreateHandler is create middleware handler
func (m SecurityToken) CreateHandler(tokens []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		response := m.value.GetResponseHelper(c)
		token := c.Query("token", "")

		for _, t := range tokens {
			if token == t {
				return c.Next()
			}
		}

		return response.JsonError(c, errors.New("tokens do not match"), fiber.StatusUnauthorized)
	}
}
