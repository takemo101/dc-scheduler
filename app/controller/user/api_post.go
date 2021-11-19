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

// ApiPostController Api配信コントローラ
type ApiPostController struct {
	value            support.ContextValue
	toastr           support.ToastrMessage
	sessionStore     core.SessionStore
	searchUseCase    application.ApiPostSearchUseCase
	storeUseCase     application.ApiPostStoreUseCase
	botDetailUseCase application.BotDetailUseCase
}

// NewApiPostController コンストラクタ
func NewApiPostController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	sessionStore core.SessionStore,
	searchUseCase application.ApiPostSearchUseCase,
	storeUseCase application.ApiPostStoreUseCase,
	botDetailUseCase application.BotDetailUseCase,
) ApiPostController {
	return ApiPostController{
		value,
		toastr,
		sessionStore,
		searchUseCase,
		storeUseCase,
		botDetailUseCase,
	}
}

// Index 一覧表示
func (ctl ApiPostController) Index(c *fiber.Ctx) (err error) {
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
		application.ApiPostSearchInput{
			Page:  form.Page,
			Limit: 20,
		},
	)
	if appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	dto.Pagination.SetURL(c.BaseURL() + c.OriginalURL())

	return response.View("user/message/api_post/index", helper.DataMap(vm.ToApiPostIndexMap(dto)))
}

// Store 追加処理
func (ctl ApiPostController) Store(c *fiber.Ctx) (err error) {
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

	_, appError := ctl.storeUseCase.Execute(
		context,
		uint(id),
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
	return response.Redirect(c, "user/message/api")
}
