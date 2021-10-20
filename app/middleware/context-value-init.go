package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// ContextValue is struct
type ContextValueInit struct {
	logger core.Logger
	app    core.Application
	value  support.ContextValue
	// value dependency
	config core.Config
	path   core.Path
}

// NewContextValueInit is create middleware
func NewContextValueInit(
	logger core.Logger,
	app core.Application,
	// value dependency
	config core.Config,
	path core.Path,
) ContextValueInit {
	return ContextValueInit{
		logger: logger,
		app:    app,
		// value dependency
		config: config,
		path:   path,
	}
}

// Setup user-value control middleware
func (m ContextValueInit) Setup() {
	m.logger.Info("setup context-value")
	m.app.App.Use(m.CreateHandler())
}

// CreateHandler is create middleware handler
func (m ContextValueInit) CreateHandler() fiber.Handler {
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
