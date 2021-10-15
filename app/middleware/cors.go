package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/takemo101/dc-scheduler/core"
)

// Cors is struct
type Cors struct {
	logger core.Logger
	app    core.Application
	config core.Config
}

// NewCors is create middleware
func NewCors(
	logger core.Logger,
	app core.Application,
	config core.Config,
) Cors {
	return Cors{
		logger: logger,
		app:    app,
		config: config,
	}
}

// Setup cors control middleware
func (m Cors) Setup() {
	m.logger.Info("setup cors")
	m.app.App.Use(m.CreateHandler())
}

// CreateHandler is create middleware handler
func (m Cors) CreateHandler() fiber.Handler {
	config := cors.Config{
		AllowMethods: strings.Join([]string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
			"HEAD",
		}, ", "),
		AllowHeaders: strings.Join([]string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"Origin",
			"Cache-Control",
			"X-Requested-With",
		}, ", "),
		AllowCredentials: true,
		MaxAge:           int(m.config.Cors.MaxAge * time.Hour),
		AllowOrigins:     strings.Join(m.config.Cors.Origins, ", "),
	}

	return cors.New(config)
}
