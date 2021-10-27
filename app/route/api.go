package route

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	controller "github.com/takemo101/dc-scheduler/app/controller/api"
	"github.com/takemo101/dc-scheduler/app/middleware"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// ApiRoute is struct
type ApiRoute struct {
	logger            core.Logger
	app               core.Application
	cors              middleware.Cors
	basicAuth         middleware.BasicAuth
	value             support.ContextValue
	security          middleware.SecurityToken
	messageController controller.PostMessageApiController
}

// NewApiRoute create new web route
func NewApiRoute(
	logger core.Logger,
	app core.Application,
	cors middleware.Cors,
	security middleware.SecurityToken,
	basicAuth middleware.BasicAuth,
	value support.ContextValue,
	messageController controller.PostMessageApiController,
) ApiRoute {
	return ApiRoute{
		logger:            logger,
		app:               app,
		cors:              cors,
		basicAuth:         basicAuth,
		value:             value,
		messageController: messageController,
		security:          security,
	}
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

	api := app.Group("/api")
	{
		// 設定からトークン取得
		values := r.app.Config.LoadToValueArray("setting", "security_token", []string{})
		tokens := make([]string, len(values))
		for i, t := range values {
			tokens[i] = t.(string)
		}

		message := api.Group("/message", r.security.CreateHandler("token", tokens))
		{
			message.Get("/send", r.messageController.Send)
		}
	}
}
