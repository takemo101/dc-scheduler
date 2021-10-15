package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
)

// SessionAuthController is session auth
type SessionAuthController struct {
	value support.RequestValue
}

// NewSessionAuthController is create auth controller
func NewSessionAuthController(
	value support.RequestValue,
) SessionAuthController {
	return SessionAuthController{
		value: value,
	}
}

// LoginForm render login form
func (ctl SessionAuthController) LoginForm(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)
	return response.View("auth/login", helper.DataMap{})
}

// Login login auth process
func (ctl SessionAuthController) Login(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)
	return response.Redirect(c, "system")
}

// Logout logout auth process
func (ctl SessionAuthController) Logout(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)
	return response.Redirect(c, "system/auth/login")
}
