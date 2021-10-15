package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/takemo101/dc-scheduler/core"
)

const CSRFContextKey string = "csrf_token"

// Csrf is struct
type Csrf struct {
	logger core.Logger
	app    core.Application
}

// NewCsrf is create middleware
func NewCsrf(
	logger core.Logger,
	app core.Application,
) Csrf {
	return Csrf{
		logger: logger,
		app:    app,
	}
}

// Setup csrf control middleware
func (m Csrf) Setup() {
	m.app.App.Use(m.CreateHandler("header:X-CSRF-Token"))
}

// CreateHandler is create middleware handler
func (m Csrf) CreateHandler(keyLookup string) fiber.Handler {
	m.logger.Info("setup csrf")

	return csrf.New(csrf.Config{
		KeyLookup:      keyLookup,
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUID,
		ContextKey:     CSRFContextKey,
	})
}

// GetCSRFToken to token
func GetCSRFToken(c *fiber.Ctx) string {
	token := c.Locals(CSRFContextKey).(string)
	return token
}
