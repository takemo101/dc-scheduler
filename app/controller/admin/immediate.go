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

// ImmediatePostController 即時配信コントローラ
type ImmediatePostController struct {
	value             support.ContextValue
	toastr            support.ToastrMessage
	config            core.Config
	sessionStore      core.SessionStore
	createFormUseCase application.PostMessageCreateFormUseCase
	storeUseCase      application.ImmediatePostStoreUseCase
}

// NewImmediatePostController コンストラクタ
func NewImmediatePostController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	config core.Config,
	sessionStore core.SessionStore,
	createFormUseCase application.PostMessageCreateFormUseCase,
	storeUseCase application.ImmediatePostStoreUseCase,
) ImmediatePostController {
	return ImmediatePostController{
		value,
		toastr,
		config,
		sessionStore,
		createFormUseCase,
		storeUseCase,
	}
}

// Create 追加フォーム
func (ctl ImmediatePostController) Create(c *fiber.Ctx) error {
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

	return response.View("message/immediate_post/create", helper.DataMap{
		"content_footer": true,
		"bot":            vm.ToBotDetailMap(dto),
	})
}

// Store 追加処理
func (ctl ImmediatePostController) Store(c *fiber.Ctx) (err error) {
	var form form.ImmediatePostCreate
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

	if err := form.Validate(c); err != nil {
		ctl.sessionStore.SetErrorResource(
			c,
			helper.ErrorsToMap(err),
			helper.StructToFormMap(&form),
		)
		return response.Back(c)
	}

	_, appError := ctl.storeUseCase.Execute(uint(id), application.ImmediatePostStoreInput{
		Message: form.Message,
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
	return response.Redirect(c, "system/message")
}
