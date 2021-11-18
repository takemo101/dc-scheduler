package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/form"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
	application "github.com/takemo101/dc-scheduler/pkg/application/user"
)

// UserController アカウント
type UserController struct {
	value             support.ContextValue
	toastr            support.ToastrMessage
	sessionStore      core.SessionStore
	registUseCase     application.UserRegistUseCase
	activationUseCase application.UserActivationUseCase
}

// NewUserController コンストラクタ
func NewUserController(
	value support.ContextValue,
	toastr support.ToastrMessage,
	sessionStore core.SessionStore,
	registUseCase application.UserRegistUseCase,
	activationUseCase application.UserActivationUseCase,
) UserController {
	return UserController{
		value,
		toastr,
		sessionStore,
		registUseCase,
		activationUseCase,
	}
}

// RegistForm render user regist form
func (ctl UserController) RegistFrom(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)

	return response.View("user/auth/regist", helper.DataMap{})
}

// Regist user regist process
func (ctl UserController) Regist(c *fiber.Ctx) error {
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

	if _, appError := ctl.registUseCase.Execute(application.UserRegistInput{
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

	return response.View("user/auth/registered", helper.DataMap{})
}

// Activation user activation process
func (ctl UserController) Activation(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)

	signature := c.Params("+")

	if appError := ctl.activationUseCase.Execute(signature); appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	return response.View("user/auth/activated", helper.DataMap{})
}
