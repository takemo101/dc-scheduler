package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// ViewRender is struct
type ViewRender struct {
	logger core.Logger
	value  support.ContextValue
}

// NewViewRender is create middleware
func NewViewRender(
	logger core.Logger,
	value support.ContextValue,
) ViewRender {
	return ViewRender{
		logger: logger,
		value:  value,
	}
}

func (v ViewRender) CreateHandler(handler func(*fiber.Ctx, *helper.ViewRender) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		render := v.value.GetViewRender(c)
		return render.HandleRender(c, handler)
	}
}
