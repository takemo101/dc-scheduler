package route

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/middleware"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// ApiRoute is struct
type ApiRoute struct {
	logger    core.Logger
	app       core.Application
	cors      middleware.Cors
	basicAuth middleware.BasicAuth
	value     support.ContextValue
}

// Setup is setup route
func (r ApiRoute) Setup() {
	r.logger.Info("setup api-route")

	app := r.app.App

	systemApi := app.Group("/system-api", r.cors.CreateHandler())
	{
		systemApi.Get("/", func(c *fiber.Ctx) error {
			response := r.value.GetResponseHelper(c)
			return response.Json(c, fiber.Map{
				"message": "it's system-api",
			})
		})
		systemApi.Get("/success", func(c *fiber.Ctx) error {
			response := r.value.GetResponseHelper(c)
			return response.JsonSuccess(c, "success", fiber.Map{
				"data": "json data",
			})
		})
		systemApi.Get("/error", func(c *fiber.Ctx) error {
			response := r.value.GetResponseHelper(c)
			return response.JsonError(c, errors.New("error"))
		})
	}
}

// NewApiRoute create new web route
func NewApiRoute(
	logger core.Logger,
	app core.Application,
	cors middleware.Cors,
	basicAuth middleware.BasicAuth,
	value support.ContextValue,
) ApiRoute {
	return ApiRoute{
		logger:    logger,
		app:       app,
		cors:      cors,
		basicAuth: basicAuth,
		value:     value,
	}
}
