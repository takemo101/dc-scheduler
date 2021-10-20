package support

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/helper"
)

const (
	ViewRenderKey     string = "view-render"
	ResponseHelperKey string = "response-helper"
)

// ContextValue
type ContextValue struct {
	//
}

// NewContextValue
func NewContextValue() ContextValue {
	return ContextValue{}
}

// --- ViewRender Context ---

// GetViewRender
func (r ContextValue) GetViewRender(c *fiber.Ctx) *helper.ViewRender {
	if render, ok := c.Locals(ViewRenderKey).(*helper.ViewRender); ok {
		return render
	}
	return nil
}

// SetViewRender
func (r ContextValue) SetViewRender(c *fiber.Ctx, render *helper.ViewRender) {
	c.Locals(ViewRenderKey, render)
}

// --- ResponseHelper Context ---

// GetResponseHelper
func (r ContextValue) GetResponseHelper(c *fiber.Ctx) *helper.ResponseHelper {
	if response, ok := c.Locals(ResponseHelperKey).(*helper.ResponseHelper); ok {
		return response
	}
	return nil
}

// SetResponseHelper
func (r ContextValue) SetResponseHelper(c *fiber.Ctx, response *helper.ResponseHelper) {
	c.Locals(ResponseHelperKey, response)
}
