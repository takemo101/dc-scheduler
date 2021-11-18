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
	common "github.com/takemo101/dc-scheduler/pkg/application/common"
)

// UserController 管理者コントローラ
type UserController struct {
	value         support.ContextValue
	toastr        support.ToastrMessage
	sessionStore  core.SessionStore
	searchUseCase application.UserSearchUseCase
	detailUseCase application.UserDetailUseCase
	storeUseCase  application.UserStoreUseCase
	updateUseCase application.UserUpdateUseCase
	deleteUseCase application.UserDeleteUseCase
}

// NewUserController コンストラクタ
func NewUserController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	sessionStore core.SessionStore,
	searchUseCase application.UserSearchUseCase,
	detailUseCase application.UserDetailUseCase,
	storeUseCase application.UserStoreUseCase,
	updateUseCase application.UserUpdateUseCase,
	deleteUseCase application.UserDeleteUseCase,
) UserController {
	return UserController{
		value,
		toastr,
		sessionStore,
		searchUseCase,
		detailUseCase,
		storeUseCase,
		updateUseCase,
		deleteUseCase,
	}
}

// Index 一覧表示
func (ctl UserController) Index(c *fiber.Ctx) (err error) {
	var form form.UserSearch
	response := ctl.value.GetResponseHelper(c)

	if err := c.QueryParser(&form); err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.searchUseCase.Execute(
		application.UserSearchInput{
			Page:  form.Page,
			Limit: 20,
		},
	)
	if appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	dto.Pagination.SetURL(c.BaseURL() + c.OriginalURL())

	return response.View("admin/user/index", helper.DataMap(vm.ToUserIndexMap(dto)))
}

// Create 追加フォーム
func (ctl UserController) Create(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)

	return response.View("admin/user/create", helper.DataMap{
		"content_footer": true,
	})
}

// Store 追加処理
func (ctl UserController) Store(c *fiber.Ctx) (err error) {
	var form form.UserCreate
	response := ctl.value.GetResponseHelper(c)

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

	_, appError := ctl.storeUseCase.Execute(application.UserStoreInput{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
		Active:   form.ActiveToBool(),
	})
	if appError != nil && appError.HasError() {
		if appError.HaveType(application.UserDuplicateError) {
			ctl.sessionStore.SetErrorResource(
				c,
				helper.ErrorToMap("email", appError),
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
	return response.Redirect(c, "system/user")
}

// Edit 編集フォーム
func (ctl UserController) Edit(c *fiber.Ctx) (err error) {
	response := ctl.value.GetResponseHelper(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.detailUseCase.Execute(uint(id))
	if appError != nil {
		if appError.HaveType(common.NotFoundDataError) {
			return response.Error(appError, fiber.StatusNotFound)
		}
		return response.Error(appError)
	}

	return response.View("admin/user/edit", helper.DataMap{
		"content_footer": true,
		"user":           vm.ToUserDetailMap(dto),
	})
}

// Update 更新処理
func (ctl UserController) Update(c *fiber.Ctx) (err error) {
	var form form.UserUpdate
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

	if appError := ctl.updateUseCase.Execute(uint(id), application.UserUpdateInput{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
		Active:   form.ActiveToBool(),
	}); appError != nil && appError.HasError() {
		if appError.HaveType(application.UserDuplicateError) {
			ctl.sessionStore.SetErrorResource(
				c,
				helper.ErrorToMap("email", appError),
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
func (ctl UserController) Delete(c *fiber.Ctx) (err error) {
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
	return response.Redirect(c, "system/user")
}
