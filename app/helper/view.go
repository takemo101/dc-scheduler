package helper

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/imdario/mergo"
)

type DataMap map[string]interface{}

// ViewRender is render manage
type ViewRender struct {
	data *ViewData
}

// ViewData is render data
type ViewData struct {
	data DataMap
	js   DataMap
	name string
	err  fiber.Error
}

func NewViewRender() *ViewRender {
	return &ViewRender{
		data: new(ViewData),
	}
}

func (v *ViewRender) Render(name string, data DataMap) error {
	v.SetName(name)
	v.SetData(data)
	return nil
}

func (v *ViewRender) SetData(data DataMap) {
	mergo.Merge(&v.data.data, data)
}

func (v *ViewRender) SetJS(data DataMap) {
	mergo.Merge(&v.data.js, data)
}

func (v *ViewRender) SetName(name string) {
	v.data.name = name
}

func (v *ViewRender) Error(message string, code int) error {
	v.data.err = fiber.Error{
		Message: message,
		Code:    code,
	}
	v.SetName("error/error")
	return nil
}

func (v *ViewRender) HandleRender(c *fiber.Ctx, handler func(*fiber.Ctx, *ViewRender)) error {
	err := c.Next()

	if err == nil && len(v.data.name) > 0 {
		handler(c, v)

		js, err := json.Marshal(v.data.js)
		if err != nil {
			return err
		}

		v.data.data["js"] = string(js)

		if v.data.err.Message == "" {
			return c.Render(v.data.name, fiber.Map(v.data.data))
		} else {
			v.data.data["error_message"] = v.data.err.Message
			c.Status(v.data.err.Code)
			return c.Render(v.data.name, fiber.Map(v.data.data))
		}
	}
	return err
}
