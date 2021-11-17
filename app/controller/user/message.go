package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/form"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/app/vm"
	"github.com/takemo101/dc-scheduler/core"
	application "github.com/takemo101/dc-scheduler/pkg/application/user"
)

// PostMessageController 配信関連コントローラ
type PostMessageController struct {
	value          support.ContextValue
	toastr         support.ToastrMessage
	config         core.Config
	sessionStore   core.SessionStore
	deleteUseCase  application.PostMessageDeleteUseCase
	historyUseCase application.SentMessageHistoryUseCase
}

// NewPostMessageController コンストラクタ
func NewPostMessageController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	config core.Config,
	sessionStore core.SessionStore,
	deleteUseCase application.PostMessageDeleteUseCase,
	historyUseCase application.SentMessageHistoryUseCase,
) PostMessageController {
	return PostMessageController{
		value,
		toastr,
		config,
		sessionStore,
		deleteUseCase,
		historyUseCase,
	}
}

// Delete 削除処理
func (ctl PostMessageController) Delete(c *fiber.Ctx) (err error) {
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

	if appError := ctl.deleteUseCase.Execute(context, uint(id)); appError != nil && appError.HasError() {
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

	dto, appError := ctl.historyUseCase.Execute(
		context,
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
