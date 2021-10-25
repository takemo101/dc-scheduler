package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/application"
)

// SchedulePostApiController 予約配信関連コントローラ
type SchedulePostApiController struct {
	value       support.ContextValue
	config      core.Config
	sendUseCase application.SchedulePostSendUseCase
}

// NewSchedulePostApiController コンストラクタ
func NewSchedulePostApiController(
	value support.ContextValue,
	config core.Config,
	sendUseCase application.SchedulePostSendUseCase,
) SchedulePostApiController {
	return SchedulePostApiController{
		value,
		config,
		sendUseCase,
	}
}

// Send 配信
func (ctl SchedulePostApiController) Send(c *fiber.Ctx) (err error) {
	response := ctl.value.GetResponseHelper(c)

	appError := ctl.sendUseCase.Execute(time.Now())
	if appError != nil && appError.HasError() {
		return response.JsonError(c, appError)
	}

	return response.JsonSuccess(c, "message schedule send successfully")
}
