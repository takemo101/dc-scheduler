package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/form"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/app/vm"
	"github.com/takemo101/dc-scheduler/core"
	application "github.com/takemo101/dc-scheduler/pkg/application/user"
)

// AccountController アカウント
type AccountController struct {
	value         support.ContextValue
	toastr        support.ToastrMessage
	sessionStore  core.SessionStore
	detailUseCase application.MyAccountDetailUseCase
	updateUseCase application.MyAccountUpdateUseCase
}

// NewAccountController コンストラクタ
func NewAccountController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	sessionStore core.SessionStore,
	detailUseCase application.MyAccountDetailUseCase,
	updateUseCase application.MyAccountUpdateUseCase,
) AccountController {
	return AccountController{
		value,
		toastr,
		sessionStore,
		detailUseCase,
		updateUseCase,
	}
}

// Edit render admin edit form
func (ctl AccountController) Edit(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)

	// UserAuthContextを生成
	auth, err := support.CreateSessionUserAuthContext(
		c,
		ctl.sessionStore,
	)
	if err != nil {
		return response.Error(err)
	}

	dto, appError := ctl.detailUseCase.Execute(auth)
	if appError != nil {
		return response.Error(err)
	}

	return response.View("user/account/edit", helper.DataMap{
		"content_footer": true,
		"user":           vm.ToUserDetailMap(dto),
	})
}

// Update admin update process
func (ctl AccountController) Update(c *fiber.Ctx) error {
	var form form.AccountUpdate
	response := ctl.value.GetResponseHelper(c)

	// UserAuthContextを生成
	auth, err := support.CreateSessionUserAuthContext(
		c,
		ctl.sessionStore,
	)
	if err != nil {
		return err
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

	if appError := ctl.updateUseCase.Execute(auth, application.MyAccountUpdateInput{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
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
