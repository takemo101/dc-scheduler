package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// DashboardController is home dashboard
type DashboardController struct {
	config core.Config
	value  support.ContextValue
}

// NewDashboardController is create dashboard
func NewDashboardController(
	config core.Config,
	value support.ContextValue,
) DashboardController {
	return DashboardController{
		config,
		value,
	}
}

// Dashboard render home
func (ctl DashboardController) Dashboard(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)

	return response.View("home", helper.DataMap{
		"config": ctl.config,
	})
}