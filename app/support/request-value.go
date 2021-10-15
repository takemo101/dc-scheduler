package support

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/helper"
)

const (
	ViewRenderKey     string = "view-render"
	ResponseHelperKey string = "response-helper"
)

// RequestValue
type RequestValue struct {
	//
}

// NewRequestValue
func NewRequestValue() RequestValue {
	return RequestValue{}
}

// GetViewRender
func (r RequestValue) GetViewRender(c *fiber.Ctx) *helper.ViewRender {
	if render, ok := c.Locals(ViewRenderKey).(*helper.ViewRender); ok {
		return render
	}
	return nil
}

// SetViewRender
func (r RequestValue) SetViewRender(c *fiber.Ctx, render *helper.ViewRender) {
	c.Locals(ViewRenderKey, render)
}

// GetResponseHelper
func (r RequestValue) GetResponseHelper(c *fiber.Ctx) *helper.ResponseHelper {
	if response, ok := c.Locals(ResponseHelperKey).(*helper.ResponseHelper); ok {
		return response
	}
	return nil
}

// SetResponseHelper
func (r RequestValue) SetResponseHelper(c *fiber.Ctx, response *helper.ResponseHelper) {
	c.Locals(ResponseHelperKey, response)
}
