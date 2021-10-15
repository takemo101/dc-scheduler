package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// RequestValue is struct
type RequestValueInit struct {
	logger core.Logger
	app    core.Application
	value  support.RequestValue
	// value dependency
	config core.Config
	path   core.Path
}

// NewRequestValueInit is create middleware
func NewRequestValueInit(
	logger core.Logger,
	app core.Application,
	// value dependency
	config core.Config,
	path core.Path,
) RequestValueInit {
	return RequestValueInit{
		logger: logger,
		app:    app,
		// value dependency
		config: config,
		path:   path,
	}
}

// Setup user-value control middleware
func (m RequestValueInit) Setup() {
	m.logger.Info("setup request-value init")
	m.app.App.Use(m.CreateHandler())
}

// CreateHandler is create middleware handler
func (m RequestValueInit) CreateHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// set ViewRender
		render := helper.NewViewRender()
		m.value.SetViewRender(c, render)
		// set ResponseHelper
		m.value.SetResponseHelper(c, helper.NewResponseHelper(
			m.logger,
			m.config,
			m.path,
			render,
		))
		return c.Next()
	}
}
