package controller

import (
	admin "github.com/takemo101/dc-scheduler/app/controller/admin"
	api "github.com/takemo101/dc-scheduler/app/controller/api"
	"go.uber.org/fx"
)

// Module コントローラモジュールDI
var Module = fx.Options(
	// admin controller
	fx.Provide(admin.NewDashboardController),
	fx.Provide(admin.NewSessionAuthController),
	fx.Provide(admin.NewAdminController),
	fx.Provide(admin.NewAccountController),
	fx.Provide(admin.NewUserController),
	fx.Provide(admin.NewBotController),
	fx.Provide(admin.NewPostMessageController),
	fx.Provide(admin.NewImmediatePostController),
	fx.Provide(admin.NewSchedulePostController),
	fx.Provide(admin.NewRegularPostController),
	fx.Provide(admin.NewRegularTimingController),

	// api controller
	fx.Provide(api.NewPostMessageApiController),
)
