package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/app/vm"
	"github.com/takemo101/dc-scheduler/core"
	query "github.com/takemo101/dc-scheduler/pkg/application/query"
)

// DashboardController is home dashboard
type DashboardController struct {
	config           core.Config
	value            support.ContextValue
	sentMessageQuery query.SentMessageQuery
}

// NewDashboardController is create dashboard
func NewDashboardController(
	config core.Config,
	value support.ContextValue,
	sentMessageQuery query.SentMessageQuery,
) DashboardController {
	return DashboardController{
		config,
		value,
		sentMessageQuery,
	}
}

// Dashboard render home
func (ctl DashboardController) Dashboard(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)

	informations := ctl.config.LoadToValueArray("setting", "information", []helper.DataMap{})

	// 配信履歴直近10件
	list, err := ctl.sentMessageQuery.RecentlyList(10)
	if err != nil {
		return response.Error(err)
	}

	return response.View("admin/home", helper.DataMap{
		"config":        ctl.config,
		"sent_messages": vm.ToSentMessagesMap(list),
		"informations":  informations,
	})
}
