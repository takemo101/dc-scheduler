package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/form"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	application "github.com/takemo101/dc-scheduler/pkg/application/admin"
	common "github.com/takemo101/dc-scheduler/pkg/application/common"
)

// ApiPostApiController Api関連コントローラ
type ApiPostApiController struct {
	value       support.ContextValue
	sendUseCase application.ApiPostSendUseCase
}

// NewApiPostApiController コンストラクタ
func NewApiPostApiController(
	value support.ContextValue,
	sendUseCase application.ApiPostSendUseCase,
) ApiPostApiController {
	return ApiPostApiController{
		value,
		sendUseCase,
	}
}

// Send 配信
func (ctl ApiPostApiController) Send(c *fiber.Ctx) (err error) {
	var form form.ApiPostSend
	response := ctl.value.GetResponseHelper(c)

	key := c.Params("key")

	if err := c.BodyParser(&form); err != nil {
		return response.JsonError(c, err)
	}

	if err := form.Sanitize(); err != nil {
		return response.JsonError(c, err)
	}

	if err := form.Validate(); err != nil {
		return response.JsonErrorMessages(c, err, helper.ErrorsToMap(err))
	}

	if appError := ctl.sendUseCase.Execute(
		key,
		form.Message,
		time.Now(),
	); appError != nil && appError.HasError() {
		if appError.HaveType(common.NotFoundDataError) {
			return response.JsonError(c, appError, fiber.StatusNotFound)
		}

		return response.JsonError(c, appError)
	}

	return response.JsonSuccess(c, "message api send successfully")
}
