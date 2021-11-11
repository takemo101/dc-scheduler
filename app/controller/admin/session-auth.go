package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/takemo101/dc-scheduler/app/form"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
	application "github.com/takemo101/dc-scheduler/pkg/application/admin"
)

// SessionAuthController is session auth
type SessionAuthController struct {
	loginUseCase  application.AdminLoginUseCase
	logoutUseCase application.AdminLogoutUseCase
	sessionStore  core.SessionStore
	value         support.ContextValue
}

// NewSessionAuthController is create auth controller
func NewSessionAuthController(
	loginUseCase application.AdminLoginUseCase,
	logoutUseCase application.AdminLogoutUseCase,
	value support.ContextValue,
	sessionStore core.SessionStore,
) SessionAuthController {
	return SessionAuthController{
		loginUseCase,
		logoutUseCase,
		sessionStore,
		value,
	}
}

// LoginForm render login form
func (ctl SessionAuthController) LoginForm(c *fiber.Ctx) error {
	response := ctl.value.GetResponseHelper(c)
	return response.View("admin/auth/login", helper.DataMap{})
}

// Login ログイン
func (ctl SessionAuthController) Login(c *fiber.Ctx) (err error) {
	var form form.Login
	response := ctl.value.GetResponseHelper(c)

	if err := c.BodyParser(&form); err != nil {
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

	// AdminAuthContextを生成
	auth, err := support.CreateSessionAdminAuthContext(
		c,
		ctl.sessionStore,
	)
	if err != nil {
		return err
	}

	// ログイン実行
	appError := ctl.loginUseCase.Execute(
		auth,
		application.AdminLoginInput{
			Email:    form.Email,
			Password: form.Password,
		},
	)
	if appError != nil && appError.HasError() {
		if appError.HaveType(application.AdminNotFoundAccountError) {
			ctl.sessionStore.SetErrorResource(
				c,
				helper.ErrorToMap("email", appError),
				helper.StructToFormMap(&form),
			)
			return response.Back(c)
		}

		return response.Error(appError)
	}

	// ダッシュボードへ
	return response.Redirect(c, "system")
}

// Logout ログアウト
func (ctl SessionAuthController) Logout(c *fiber.Ctx) (err error) {
	response := ctl.value.GetResponseHelper(c)

	// AdminAuthContextを生成
	auth, err := support.CreateSessionAdminAuthContext(
		c,
		ctl.sessionStore,
	)
	if err != nil {
		return err
	}

	if appError := ctl.logoutUseCase.Execute(auth); appError != nil && appError.HasError() {
		return response.Error(appError)
	}

	return response.Redirect(c, "system/auth/login")
}
