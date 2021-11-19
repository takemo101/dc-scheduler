package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/form"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/app/vm"
	"github.com/takemo101/dc-scheduler/core"
	application "github.com/takemo101/dc-scheduler/pkg/application/admin"
)

// PostMessageController 配信関連コントローラ
type PostMessageController struct {
	value          support.ContextValue
	toastr         support.ToastrMessage
	config         core.Config
	sessionStore   core.SessionStore
	searchUseCase  application.PostMessageSearchUseCase
	deleteUseCase  application.PostMessageDeleteUseCase
	historyUseCase application.SentMessageHistoryUseCase
}

// NewPostMessageController コンストラクタ
func NewPostMessageController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	config core.Config,
	sessionStore core.SessionStore,
	searchUseCase application.PostMessageSearchUseCase,
	deleteUseCase application.PostMessageDeleteUseCase,
	historyUseCase application.SentMessageHistoryUseCase,
) PostMessageController {
	return PostMessageController{
		value,
		toastr,
		config,
		sessionStore,
		searchUseCase,
		deleteUseCase,
		historyUseCase,
	}
}

// Index 一覧表示
func (ctl PostMessageController) Index(c *fiber.Ctx) (err error) {
	var form form.PostMessageSearch
	response := ctl.value.GetResponseHelper(c)

	if err := c.QueryParser(&form); err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.searchUseCase.Execute(
		application.PostMessageSearchInput{
			Page:  form.Page,
			Limit: 20,
		},
	)
	if appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	dto.Pagination.SetURL(c.BaseURL() + c.OriginalURL())

	return response.View("admin/message/index", helper.DataMap(vm.ToPostMessageIndexMap(dto)))
}

// Delete 削除処理
func (ctl PostMessageController) Delete(c *fiber.Ctx) (err error) {
	response := ctl.value.GetResponseHelper(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(err)
	}

	if appError := ctl.deleteUseCase.Execute(uint(id)); appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	ctl.toastr.SetToastr(
		c,
		support.ToastrDelete,
		support.ToastrDelete.Message(),
		support.Messages{},
	)
	return response.Back(c)
}

// History 削除処理
func (ctl PostMessageController) History(c *fiber.Ctx) (err error) {
	var form form.SentMessageHistory
	response := ctl.value.GetResponseHelper(c)

	if err := c.QueryParser(&form); err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.historyUseCase.Execute(
		application.SentMessageHistoryInput{
			Page:  form.Page,
			Limit: 20,
		},
	)
	if appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	dto.Pagination.SetURL(c.BaseURL() + c.OriginalURL())

	return response.View("admin/message/history", helper.DataMap(vm.ToSentMessageHistoryMap(dto)))
}
