package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
	application "github.com/takemo101/dc-scheduler/pkg/application/admin"
	common "github.com/takemo101/dc-scheduler/pkg/application/common"
)

// PostMessageApiController 配信関連コントローラ
type PostMessageApiController struct {
	value           support.ContextValue
	config          core.Config
	scheduleUseCase application.SchedulePostSendUseCase
	regularUseCase  application.RegularPostSendUseCase
}

// NewPostMessageApiController コンストラクタ
func NewPostMessageApiController(
	value support.ContextValue,
	config core.Config,
	scheduleUseCase application.SchedulePostSendUseCase,
	regularUseCase application.RegularPostSendUseCase,
) PostMessageApiController {
	return PostMessageApiController{
		value,
		config,
		scheduleUseCase,
		regularUseCase,
	}
}

// Send 配信
func (ctl PostMessageApiController) Send(c *fiber.Ctx) (err error) {
	response := ctl.value.GetResponseHelper(c)

	sends := make(chan common.AppError, 2)
	defer close(sends)

	go func() {
		sends <- ctl.scheduleUseCase.Execute(time.Now())
	}()
	go func() {
		sends <- ctl.regularUseCase.Execute(time.Now())
	}()

	counter := 0

	for appError := range sends {
		counter++
		if appError != nil {
			return response.JsonError(c, appError)
		} else if counter == 2 {
			break
		}
	}

	return response.JsonSuccess(c, "message schedule send successfully")
}
