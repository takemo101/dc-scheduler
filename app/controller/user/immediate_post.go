package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/form"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/app/vm"
	"github.com/takemo101/dc-scheduler/core"
	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	application "github.com/takemo101/dc-scheduler/pkg/application/user"
)

// ImmediatePostController 即時配信コントローラ
type ImmediatePostController struct {
	value            support.ContextValue
	toastr           support.ToastrMessage
	sessionStore     core.SessionStore
	searchUseCase    application.ImmediatePostSearchUseCase
	storeUseCase     application.ImmediatePostStoreUseCase
	botDetailUseCase application.BotDetailUseCase
}

// NewImmediatePostController コンストラクタ
func NewImmediatePostController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	sessionStore core.SessionStore,
	searchUseCase application.ImmediatePostSearchUseCase,
	storeUseCase application.ImmediatePostStoreUseCase,
	botDetailUseCase application.BotDetailUseCase,
) ImmediatePostController {
	return ImmediatePostController{
		value,
		toastr,
		sessionStore,
		searchUseCase,
		storeUseCase,
		botDetailUseCase,
	}
}

// Index 一覧表示
func (ctl ImmediatePostController) Index(c *fiber.Ctx) (err error) {
	var form form.PostMessageSearch
	response := ctl.value.GetResponseHelper(c)

	// UserAuthContextを生成
	context, err := support.CreateSessionUserAuthContext(
		c,
		ctl.sessionStore,
	)
	if err != nil {
		return response.Error(err)
	}

	if err := c.QueryParser(&form); err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.searchUseCase.Execute(
		context,
		application.ImmediatePostSearchInput{
			Page:  form.Page,
			Limit: 20,
		},
	)
	if appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	dto.Pagination.SetURL(c.BaseURL() + c.OriginalURL())

	return response.View("admin/message/immediate_post/index", helper.DataMap(vm.ToImmediatePostIndexMap(dto)))
}

// Create 追加フォーム
func (ctl ImmediatePostController) Create(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)

	// UserAuthContextを生成
	context, err := support.CreateSessionUserAuthContext(
		c,
		ctl.sessionStore,
	)
	if err != nil {
		return response.Error(err)
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.botDetailUseCase.Execute(context, uint(id))
	if appError != nil {
		if appError.HaveType(common.NotFoundDataError) {
			return response.Error(appError, fiber.StatusNotFound)
		}
		return response.Error(appError)
	}

	return response.View("user/message/immediate_post/create", helper.DataMap{
		"content_footer": true,
		"bot":            vm.ToBotDetailMap(dto),
	})
}

// Store 追加処理
func (ctl ImmediatePostController) Store(c *fiber.Ctx) (err error) {
	var form form.ImmediatePostCreate
	response := ctl.value.GetResponseHelper(c)

	// UserAuthContextを生成
	context, err := support.CreateSessionUserAuthContext(
		c,
		ctl.sessionStore,
	)
	if err != nil {
		return response.Error(err)
	}

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

	_, appError := ctl.storeUseCase.Execute(
		context,
		uint(id),
		application.ImmediatePostStoreInput{
			Message: form.Message,
		},
	)
	if appError != nil && appError.HasError() {
		if appError.HaveType(common.NotFoundDataError) {
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
	return response.Redirect(c, "user/message/immediate")
}
