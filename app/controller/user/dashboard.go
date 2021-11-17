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
	sessionStore     core.SessionStore
	sentMessageQuery query.SentMessageQuery
}

// NewDashboardController is create dashboard
func NewDashboardController(
	config core.Config,
	value support.ContextValue,
	sessionStore core.SessionStore,
	sentMessageQuery query.SentMessageQuery,
) DashboardController {
	return DashboardController{
		config,
		value,
		sessionStore,
		sentMessageQuery,
	}
}

// Dashboard render home
func (ctl DashboardController) Dashboard(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)

	// UserAuthContextを生成
	context, err := support.CreateSessionUserAuthContext(
		c,
		ctl.sessionStore,
	)
	if err != nil {
		return response.Error(err)
	}

	auth, err := context.UserAuth()
	if err != nil {
		return response.Error(err)
	}

	informations := ctl.config.LoadToValueArray("setting", "information", []helper.DataMap{})

	// 配信履歴直近10件
	list, err := ctl.sentMessageQuery.RecentlyListByUserID(auth.ID(), 10)
	if err != nil {
		return response.Error(err)
	}

	return response.View("user/home", helper.DataMap{
		"config":        ctl.config,
		"sent_messages": vm.ToSentMessagesMap(list),
		"informations":  informations,
	})
}
