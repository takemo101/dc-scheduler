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

// ApiPostController Api配信コントローラ
type ApiPostController struct {
	value         support.ContextValue
	toastr        support.ToastrMessage
	sessionStore  core.SessionStore
	searchUseCase application.ApiPostSearchUseCase
}

// NewApiPostController コンストラクタ
func NewApiPostController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	sessionStore core.SessionStore,
	searchUseCase application.ApiPostSearchUseCase,
) ApiPostController {
	return ApiPostController{
		value,
		toastr,
		sessionStore,
		searchUseCase,
	}
}

// Index 一覧表示
func (ctl ApiPostController) Index(c *fiber.Ctx) (err error) {
	var form form.PostMessageSearch
	response := ctl.value.GetResponseHelper(c)

	if err := c.QueryParser(&form); err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.searchUseCase.Execute(
		application.ApiPostSearchInput{
			Page:  form.Page,
			Limit: 20,
		},
	)
	if appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	dto.Pagination.SetURL(c.BaseURL() + c.OriginalURL())

	return response.View("admin/message/api_post/index", helper.DataMap(vm.ToApiPostIndexMap(dto)))
}
