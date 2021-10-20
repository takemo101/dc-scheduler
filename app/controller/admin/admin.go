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
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// AdminController 管理者コントローラ
type AdminController struct {
	value         support.ContextValue
	toastr        support.ToastrMessage
	sessionStore  core.SessionStore
	searchUseCase application.AdminSearchUseCase
	detailUseCase application.AdminDetailUseCase
	storeUseCase  application.AdminStoreUseCase
	updateUseCase application.AdminUpdateUseCase
	deleteUseCase application.AdminDeleteUseCase
}

// NewAdminController コンストラクタ
func NewAdminController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	sessionStore core.SessionStore,
	searchUseCase application.AdminSearchUseCase,
	detailUseCase application.AdminDetailUseCase,
	storeUseCase application.AdminStoreUseCase,
	updateUseCase application.AdminUpdateUseCase,
	deleteUseCase application.AdminDeleteUseCase,
) AdminController {
	return AdminController{
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
func (ctl AdminController) Index(c *fiber.Ctx) (err error) {
	var form form.AdminSearch
	response := ctl.value.GetResponseHelper(c)

	if err := c.QueryParser(&form); err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.searchUseCase.Execute(
		application.AdminSearchInput{
			Page:  form.Page,
			Limit: 20,
		},
	)
	if appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	dto.Pagination.SetURL(c.BaseURL() + c.OriginalURL())

	return response.View("admin/index", helper.DataMap(vm.ToAdminIndexMap(dto)))
}

// Create 追加フォーム
func (ctl AdminController) Create(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)

	return response.View("admin/create", helper.DataMap{
		"content_footer": true,
		"roles":          vm.ToKeyValueMap(domain.AdminRoleToArray()),
	})
}

// Store 追加処理
func (ctl AdminController) Store(c *fiber.Ctx) (err error) {
	var form form.AdminCreate
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

	_, appError := ctl.storeUseCase.Execute(application.AdminStoreInput{
		Name:     form.Name,
		Email:    form.Email,
		Role:     form.Role,
		Password: form.Password,
	})
	if appError != nil && appError.HasError() {
		if appError.HaveType(application.AdminDuplicateError) {
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
	return response.Redirect(c, "system/admin")
}

// Edit 編集フォーム
func (ctl AdminController) Edit(c *fiber.Ctx) (err error) {
	response := ctl.value.GetResponseHelper(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.detailUseCase.Execute(uint(id))
	if appError != nil {
		if appError.HaveType(application.NotFoundDataError) {
			return response.Error(appError, fiber.StatusNotFound)
		}
		return response.Error(appError)
	}

	return response.View("admin/edit", helper.DataMap{
		"content_footer": true,
		"admin":          vm.ToAdminDetailMap(dto),
		"roles":          vm.ToKeyValueMap(domain.AdminRoleToArray()),
	})
}

// Update 更新処理
func (ctl AdminController) Update(c *fiber.Ctx) (err error) {
	var form form.AdminUpdate
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

	if appError := ctl.updateUseCase.Execute(uint(id), application.AdminUpdateInput{
		Name:     form.Name,
		Email:    form.Email,
		Role:     form.Role,
		Password: form.Password,
	}); appError != nil && appError.HasError() {
		if appError.HaveType(application.AdminDuplicateError) {
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
func (ctl AdminController) Delete(c *fiber.Ctx) (err error) {
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
	return response.Redirect(c, "system/admin")
}
