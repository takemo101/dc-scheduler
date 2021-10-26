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

// RegularPostController 即時配信コントローラ
type RegularPostController struct {
	value             support.ContextValue
	toastr            support.ToastrMessage
	sessionStore      core.SessionStore
	searchUseCase     application.RegularPostSearchUseCase
	createFormUseCase application.PostMessageCreateFormUseCase
	storeUseCase      application.RegularPostStoreUseCase
	editFormUseCase   application.RegularPostEditFormUseCase
	updateUseCase     application.RegularPostUpdateUseCase
	addUseCase        application.RegularTimingAddUseCase
	removeUseCase     application.RegularTimingRemoveUseCase
}

// NewRegularPostController コンストラクタ
func NewRegularPostController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	sessionStore core.SessionStore,
	searchUseCase application.RegularPostSearchUseCase,
	createFormUseCase application.PostMessageCreateFormUseCase,
	storeUseCase application.RegularPostStoreUseCase,
	editFormUseCase application.RegularPostEditFormUseCase,
	updateUseCase application.RegularPostUpdateUseCase,
	addUseCase application.RegularTimingAddUseCase,
	removeUseCase application.RegularTimingRemoveUseCase,
) RegularPostController {
	return RegularPostController{
		value,
		toastr,
		sessionStore,
		searchUseCase,
		createFormUseCase,
		storeUseCase,
		editFormUseCase,
		updateUseCase,
		addUseCase,
		removeUseCase,
	}
}

// Index 一覧表示
func (ctl RegularPostController) Index(c *fiber.Ctx) (err error) {
	var form form.PostMessageSearch
	response := ctl.value.GetResponseHelper(c)

	if err := c.QueryParser(&form); err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.searchUseCase.Execute(
		application.RegularPostSearchInput{
			Page:  form.Page,
			Limit: 20,
		},
	)
	if appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	dto.Pagination.SetURL(c.BaseURL() + c.OriginalURL())

	return response.View("message/regular_post/index", helper.DataMap(vm.ToRegularPostIndexMap(dto)))
}

// Create 追加フォーム
func (ctl RegularPostController) Create(c *fiber.Ctx) error {
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

	return response.View("message/regular_post/create", helper.DataMap{
		"content_footer": true,
		"bot":            vm.ToBotDetailMap(dto),
	})
}

// Store 追加処理
func (ctl RegularPostController) Store(c *fiber.Ctx) (err error) {
	var form form.RegularPostCreateAndUpdate
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

	_, appError := ctl.storeUseCase.Execute(uint(id), application.RegularPostStoreInput{
		Message: form.Message,
		Active:  form.ActiveToBool(),
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
	return response.Redirect(c, "system/message/regular")
}

// Edit 編集フォーム
func (ctl RegularPostController) Edit(c *fiber.Ctx) (err error) {
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

	return response.View("message/regular_post/edit", helper.DataMap{
		"content_footer": true,
		"regular_post":   vm.ToRegularPostDetailMap(dto),
	})
}

// Update 更新処理
func (ctl RegularPostController) Update(c *fiber.Ctx) (err error) {
	var form form.RegularPostCreateAndUpdate
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

	if appError := ctl.updateUseCase.Execute(uint(id), application.RegularPostUpdateInput{
		Message: form.Message,
		Active:  form.ActiveToBool(),
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
