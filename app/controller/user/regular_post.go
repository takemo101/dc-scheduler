package controller

import (
	"fmt"
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

// RegularPostController 定期配信コントローラ
type RegularPostController struct {
	value            support.ContextValue
	toastr           support.ToastrMessage
	sessionStore     core.SessionStore
	searchUseCase    application.RegularPostSearchUseCase
	storeUseCase     application.RegularPostStoreUseCase
	editFormUseCase  application.RegularPostEditFormUseCase
	updateUseCase    application.RegularPostUpdateUseCase
	addUseCase       application.RegularTimingAddUseCase
	removeUseCase    application.RegularTimingRemoveUseCase
	botDetailUseCase application.BotDetailUseCase
}

// NewRegularPostController コンストラクタ
func NewRegularPostController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	sessionStore core.SessionStore,
	searchUseCase application.RegularPostSearchUseCase,
	storeUseCase application.RegularPostStoreUseCase,
	editFormUseCase application.RegularPostEditFormUseCase,
	updateUseCase application.RegularPostUpdateUseCase,
	addUseCase application.RegularTimingAddUseCase,
	removeUseCase application.RegularTimingRemoveUseCase,
	botDetailUseCase application.BotDetailUseCase,
) RegularPostController {
	return RegularPostController{
		value,
		toastr,
		sessionStore,
		searchUseCase,
		storeUseCase,
		editFormUseCase,
		updateUseCase,
		addUseCase,
		removeUseCase,
		botDetailUseCase,
	}
}

// Index 一覧表示
func (ctl RegularPostController) Index(c *fiber.Ctx) (err error) {
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
		application.RegularPostSearchInput{
			Page:  form.Page,
			Limit: 20,
		},
	)
	if appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	dto.Pagination.SetURL(c.BaseURL() + c.OriginalURL())

	return response.View("user/message/regular_post/index", helper.DataMap(vm.ToRegularPostIndexMap(dto)))
}

// Create 追加フォーム
func (ctl RegularPostController) Create(c *fiber.Ctx) error {
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

	return response.View("user/message/regular_post/create", helper.DataMap{
		"content_footer": true,
		"bot":            vm.ToBotDetailMap(dto),
	})
}

// Store 追加処理
func (ctl RegularPostController) Store(c *fiber.Ctx) (err error) {
	var form form.RegularPostCreateAndUpdate
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

	storeID, appError := ctl.storeUseCase.Execute(
		context,
		uint(id),
		application.RegularPostStoreInput{
			Message: form.Message,
			Active:  form.ActiveToBool(),
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

	return response.Redirect(
		c,
		fmt.Sprintf("user/message/regular/%d/timing/edit", storeID),
	)
}

// Edit 編集フォーム
func (ctl RegularPostController) Edit(c *fiber.Ctx) (err error) {
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

	dto, appError := ctl.editFormUseCase.Execute(context, uint(id))
	if appError != nil {
		if appError.HaveType(common.NotFoundDataError) {
			return response.Error(appError, fiber.StatusNotFound)
		}
		return response.Error(appError)
	}

	return response.View("user/message/regular_post/edit", helper.DataMap{
		"content_footer": true,
		"regular_post":   vm.ToRegularPostDetailMap(dto),
	})
}

// Update 更新処理
func (ctl RegularPostController) Update(c *fiber.Ctx) (err error) {
	var form form.RegularPostCreateAndUpdate
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

	if appError := ctl.updateUseCase.Execute(
		context,
		uint(id),
		application.RegularPostUpdateInput{
			Message: form.Message,
			Active:  form.ActiveToBool(),
		},
	); appError != nil && appError.HasError() {
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
