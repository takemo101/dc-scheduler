package controller

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/form"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/app/vm"
	"github.com/takemo101/dc-scheduler/core"
	application "github.com/takemo101/dc-scheduler/pkg/application/admin"
	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// RegularTimingController 即時配信コントローラ
type RegularTimingController struct {
	value           support.ContextValue
	toastr          support.ToastrMessage
	sessionStore    core.SessionStore
	editFormUseCase application.RegularPostEditFormUseCase
	addUseCase      application.RegularTimingAddUseCase
	removeUseCase   application.RegularTimingRemoveUseCase
}

// NewRegularTimingController コンストラクタ
func NewRegularTimingController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	sessionStore core.SessionStore,
	editFormUseCase application.RegularPostEditFormUseCase,
	addUseCase application.RegularTimingAddUseCase,
	removeUseCase application.RegularTimingRemoveUseCase,
) RegularTimingController {
	return RegularTimingController{
		value,
		toastr,
		sessionStore,
		editFormUseCase,
		addUseCase,
		removeUseCase,
	}
}

// Edit 編集フォーム
func (ctl RegularTimingController) Edit(c *fiber.Ctx) (err error) {
	response := ctl.value.GetResponseHelper(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.editFormUseCase.Execute(uint(id))
	if appError != nil {
		if appError.HaveType(common.NotFoundDataError) {
			return response.Error(appError, fiber.StatusNotFound)
		}
		return response.Error(appError)
	}

	return response.View("message/regular_post/timing", helper.DataMap{
		"content_footer": true,
		"day_of_weeks":   vm.ToKeyValueMap(domain.DayOfWeekToArray()),
		"regular_post":   vm.ToRegularPostDetailMap(dto),
	})
}

// Store 追加処理
func (ctl RegularTimingController) Add(c *fiber.Ctx) (err error) {
	var form form.RegularTimingAdd
	response := ctl.value.GetResponseHelper(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(err)
	}

	if err := c.BodyParser(&form); err != nil {
		return response.Error(err)
	}

	if err := form.Validate(); err != nil {
		ctl.sessionStore.SetErrorResource(
			c,
			helper.ErrorsToMap(err),
			helper.StructToFormMap(&form),
		)
		return response.Back(c)
	}

	appError := ctl.addUseCase.Execute(uint(id), application.RegularTimingInput{
		DayOfWeek: form.DayOfWeek,
		HourTime:  form.HourTimeToTime(),
	})
	if appError != nil && appError.HasError() {
		if appError.HaveType(common.NotFoundDataError) {
			return response.Error(appError, fiber.StatusNotFound)
		} else if appError.HaveType(application.RegularTimingDuplicateError) {
			ctl.sessionStore.SetErrorResource(
				c,
				helper.ErrorToMap("hour_time", appError),
				helper.StructToFormMap(&form),
			)
			return response.Back(c)
		}
		return response.Error(appError)
	}

	ctl.toastr.SetToastr(
		c,
		support.ToastrStore,
		support.ToastrStore.Message(),
		support.Messages{},
	)
	return response.Back(c)
}

// Remove 削除処理
func (ctl RegularTimingController) Remove(c *fiber.Ctx) (err error) {
	response := ctl.value.GetResponseHelper(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(err)
	}

	dayOfWeek := c.Params("day_of_week", "")
	hourTime := c.Params("hour_time", "")

	// hour_time を time.Timeに変換
	ht, err := time.ParseInLocation("15:04", hourTime, time.Local)
	if err != nil {
		return response.Error(err)
	}

	if appError := ctl.removeUseCase.Execute(uint(id), application.RegularTimingInput{
		DayOfWeek: dayOfWeek,
		HourTime:  ht,
	}); appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	ctl.toastr.SetToastr(
		c,
		support.ToastrUpdate,
		support.ToastrUpdate.Message(),
		support.Messages{},
	)
	return response.Back(c)
}
