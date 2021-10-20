package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
	"github.com/thoas/go-funk"
)

// BasicAuth is struct
type BasicAuth struct {
	logger core.Logger
	config core.Config
	value  support.ContextValue
}

// NewBasicAuth is create middleware
func NewBasicAuth(
	logger core.Logger,
	config core.Config,
	value support.ContextValue,
) BasicAuth {
	return BasicAuth{
		logger: logger,
		config: config,
		value:  value,
	}
}

// CreateHandler is create middleware handler
func (m BasicAuth) CreateHandler() fiber.Handler {
	m.logger.Info("setup basic-auth")

	mapData, _ := m.config.Load("basic-auth")

	r := funk.Map(mapData, func(k string, v interface{}) (string, string) {
		return k, v.(string)
	})

	return basicauth.New(basicauth.Config{
		Users: r.(map[string]string),
		Unauthorized: func(c *fiber.Ctx) error {
			response := m.value.GetResponseHelper(c)

			return response.JsonError(c, fiber.ErrUnauthorized, fiber.StatusUnauthorized)
		},
	})
}
