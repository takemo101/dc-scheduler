package controller

import (
	admin "github.com/takemo101/dc-scheduler/app/controller/admin"
	"go.uber.org/fx"
)

// Module コントローラモジュールDI
var Module = fx.Options(
	// admin controller
	fx.Provide(admin.NewDashboardController),
	fx.Provide(admin.NewSessionAuthController),
	fx.Provide(admin.NewAdminController),
	fx.Provide(admin.NewAccountController),
	fx.Provide(admin.NewBotController),
	fx.Provide(admin.NewPostMessageController),
	fx.Provide(admin.NewImmediatePostController),
)