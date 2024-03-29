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

// BotController ボットコントローラ
type BotController struct {
	value         support.ContextValue
	toastr        support.ToastrMessage
	config        core.Config
	sessionStore  core.SessionStore
	searchUseCase application.BotSearchUseCase
	detailUseCase application.BotDetailUseCase
	storeUseCase  application.BotStoreUseCase
	updateUseCase application.BotUpdateUseCase
	deleteUseCase application.BotDeleteUseCase
}

// NewBotController コンストラクタ
func NewBotController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	config core.Config,
	sessionStore core.SessionStore,
	searchUseCase application.BotSearchUseCase,
	detailUseCase application.BotDetailUseCase,
	storeUseCase application.BotStoreUseCase,
	updateUseCase application.BotUpdateUseCase,
	deleteUseCase application.BotDeleteUseCase,
) BotController {
	return BotController{
		value,
		toastr,
		config,
		sessionStore,
		searchUseCase,
		detailUseCase,
		storeUseCase,
		updateUseCase,
		deleteUseCase,
	}
}

// Index 一覧表示
func (ctl BotController) Index(c *fiber.Ctx) (err error) {
	var form form.BotSearch
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
		application.BotSearchInput{
			Page:  form.Page,
			Limit: 20,
		},
	)
	if appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	dto.Pagination.SetURL(c.BaseURL() + c.OriginalURL())

	return response.View("user/bot/index", helper.DataMap(vm.ToBotIndexMap(dto)))
}

// Create 追加フォーム
func (ctl BotController) Create(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)

	return response.View("user/bot/create", helper.DataMap{
		"content_footer": true,
	})
}

// Store 追加処理
func (ctl BotController) Store(c *fiber.Ctx) (err error) {
	var form form.BotCreateAndUpdate
	response := ctl.value.GetResponseHelper(c)

	// UserAuthContextを生成
	context, err := support.CreateSessionUserAuthContext(
		c,
		ctl.sessionStore,
	)
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

	// アバターファイル取得
	file, _ := c.FormFile("avatar")

	// ディレクトリー設定取得
	directory := ctl.config.LoadToValueString(
		"setting",
		"directory.bot_avatar",
		"",
	)

	_, appError := ctl.storeUseCase.Execute(
		context,
		application.BotStoreInput{
			Name:            form.Name,
			AtatarFile:      file,
			AtatarDirectory: directory,
			Webhook:         form.Webhook,
			Active:          form.ActiveToBool(),
		},
	)
	if appError != nil && appError.HasError() {
		if appError.HaveType(application.BotDuplicateError) || appError.HaveType(application.BotWebhookInvalidError) {
			ctl.sessionStore.SetErrorResource(
				c,
				helper.ErrorToMap("webhook", appError),
				helper.StructToFormMap(&form),
			)
			return response.Back(c)
		}

		return response.Error(appError)
	}

	ctl.toastr.SetToastr(
		c,
		support.ToastrStore,
		support.ToastrStore.Message(),
		support.Messages{},
	)
	return response.Redirect(c, "user/bot")
}

// Edit 編集フォーム
func (ctl BotController) Edit(c *fiber.Ctx) (err error) {
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

	dto, appError := ctl.detailUseCase.Execute(context, uint(id))
	if appError != nil {
		if appError.HaveType(common.NotFoundDataError) {
			return response.Error(appError, fiber.StatusNotFound)
		}
		return response.Error(appError)
	}

	return response.View("user/bot/edit", helper.DataMap{
		"content_footer": true,
		"bot":            vm.ToBotDetailMap(dto),
	})
}

// Update 更新処理
func (ctl BotController) Update(c *fiber.Ctx) (err error) {
	var form form.BotCreateAndUpdate
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

	if err := form.Validate(c); err != nil {
		ctl.sessionStore.SetErrorResource(
			c,
			helper.ErrorsToMap(err),
			helper.StructToFormMap(&form),
		)
		return response.Back(c)
	}

	// アバターファイル取得
	file, err := c.FormFile("avatar")
	if err != nil {
		file = nil
	}

	// ディレクトリー設定取得
	directory := ctl.config.LoadToValueString(
		"setting",
		"directory.bot_avatar",
		"",
	)

	if appError := ctl.updateUseCase.Execute(
		context,
		uint(id),
		application.BotUpdateInput{
			Name:            form.Name,
			AtatarFile:      file,
			AtatarDirectory: directory,
			Webhook:         form.Webhook,
			Active:          form.ActiveToBool(),
		},
	); appError != nil && appError.HasError() {
		if appError.HaveType(application.BotDuplicateError) || appError.HaveType(application.BotWebhookInvalidError) {
			ctl.sessionStore.SetErrorResource(
				c,
				helper.ErrorToMap("webhook", appError),
				helper.StructToFormMap(&form),
			)
			return response.Back(c)
		}

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

// Delete 削除処理
func (ctl BotController) Delete(c *fiber.Ctx) (err error) {
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
	return response.Redirect(c, "user/bot")
}
