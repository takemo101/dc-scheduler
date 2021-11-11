package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/form"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/app/vm"
	"github.com/takemo101/dc-scheduler/core"
	application "github.com/takemo101/dc-scheduler/pkg/application/admin"
)

// ImmediatePostController 即時配信コントローラ
type ImmediatePostController struct {
	value         support.ContextValue
	toastr        support.ToastrMessage
	sessionStore  core.SessionStore
	searchUseCase application.ImmediatePostSearchUseCase
}

// NewImmediatePostController コンストラクタ
func NewImmediatePostController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	sessionStore core.SessionStore,
	searchUseCase application.ImmediatePostSearchUseCase,
) ImmediatePostController {
	return ImmediatePostController{
		value,
		toastr,
		sessionStore,
		searchUseCase,
	}
}

// Index 一覧表示
func (ctl ImmediatePostController) Index(c *fiber.Ctx) (err error) {
	var form form.PostMessageSearch
	response := ctl.value.GetResponseHelper(c)

	if err := c.QueryParser(&form); err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.searchUseCase.Execute(
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
