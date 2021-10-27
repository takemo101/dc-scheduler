package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/form"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/app/vm"
	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/application"
)

// SchedulePostController 即時配信コントローラ
type SchedulePostController struct {
	value             support.ContextValue
	toastr            support.ToastrMessage
	sessionStore      core.SessionStore
	searchUseCase     application.SchedulePostSearchUseCase
	createFormUseCase application.PostMessageCreateFormUseCase
	storeUseCase      application.SchedulePostStoreUseCase
	editFormUseCase   application.SchedulePostEditFormUseCase
	updateUseCase     application.SchedulePostUpdateUseCase
}

// NewSchedulePostController コンストラクタ
func NewSchedulePostController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	sessionStore core.SessionStore,
	searchUseCase application.SchedulePostSearchUseCase,
	createFormUseCase application.PostMessageCreateFormUseCase,
	storeUseCase application.SchedulePostStoreUseCase,
	editFormUseCase application.SchedulePostEditFormUseCase,
	updateUseCase application.SchedulePostUpdateUseCase,
) SchedulePostController {
	return SchedulePostController{
		value,
		toastr,
		sessionStore,
		searchUseCase,
		createFormUseCase,
		storeUseCase,
		editFormUseCase,
		updateUseCase,
	}
}

// Index 一覧表示
func (ctl SchedulePostController) Index(c *fiber.Ctx) (err error) {
	var form form.PostMessageSearch
	response := ctl.value.GetResponseHelper(c)

	if err := c.QueryParser(&form); err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.searchUseCase.Execute(
		application.SchedulePostSearchInput{
			Page:  form.Page,
			Limit: 20,
		},
	)
	if appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	dto.Pagination.SetURL(c.BaseURL() + c.OriginalURL())

	return response.View("message/schedule_post/index", helper.DataMap(vm.ToSchedulePostIndexMap(dto)))
}

// Create 追加フォーム
func (ctl SchedulePostController) Create(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.createFormUseCase.Execute(uint(id))
	if appError != nil {
		if appError.HaveType(application.NotFoundDataError) {
			return response.Error(appError, fiber.StatusNotFound)
		}
		return response.Error(appError)
	}

	return response.View("message/schedule_post/create", helper.DataMap{
		"content_footer": true,
		"bot":            vm.ToBotDetailMap(dto),
	})
}

// Store 追加処理
func (ctl SchedulePostController) Store(c *fiber.Ctx) (err error) {
	var form form.SchedulePostCreateAndUpdate
	response := ctl.value.GetResponseHelper(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(err)
	}

	if err := c.BodyParser(&form); err != nil {
		return response.Error(err)
	}

	if err := form.Sanitize(); err != nil {
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

	_, appError := ctl.storeUseCase.Execute(uint(id), application.SchedulePostStoreInput{
		Message:       form.Message,
		ReservationAt: form.ReservationAtToTime(),
	})
	if appError != nil && appError.HasError() {
		if appError.HaveType(application.NotFoundDataError) {
			return response.Error(appError, fiber.StatusNotFound)
		}
		return response.Error(appError)
	}

	ctl.toastr.SetToastr(
		c,
		support.ToastrStore,
		support.ToastrStore.Message(),
		support.Messages{},
	)
	return response.Redirect(c, "system/message/schedule")
}

// Edit 編集フォーム
func (ctl SchedulePostController) Edit(c *fiber.Ctx) (err error) {
	response := ctl.value.GetResponseHelper(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.editFormUseCase.Execute(uint(id))
	if appError != nil {
		if appError.HaveType(application.NotFoundDataError) {
			return response.Error(appError, fiber.StatusNotFound)
		}
		return response.Error(appError)
	}

	return response.View("message/schedule_post/edit", helper.DataMap{
		"content_footer": true,
		"schedule_post":  vm.ToSchedulePostDetailMap(dto),
	})
}

// Update 更新処理
func (ctl SchedulePostController) Update(c *fiber.Ctx) (err error) {
	var form form.SchedulePostCreateAndUpdate
	response := ctl.value.GetResponseHelper(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(err)
	}

	if err := c.BodyParser(&form); err != nil {
		return response.Error(err)
	}

	if err := form.Sanitize(); err != nil {
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

	if appError := ctl.updateUseCase.Execute(uint(id), application.SchedulePostUpdateInput{
		Message:       form.Message,
		ReservationAt: form.ReservationAtToTime(),
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
