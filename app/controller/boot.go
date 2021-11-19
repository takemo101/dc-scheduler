package controller

import (
	admin "github.com/takemo101/dc-scheduler/app/controller/admin"
	api "github.com/takemo101/dc-scheduler/app/controller/api"
	user "github.com/takemo101/dc-scheduler/app/controller/user"
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
	fx.Provide(admin.NewApiPostController),

	// user controller
	fx.Provide(user.NewDashboardController),
	fx.Provide(user.NewSessionAuthController),
	fx.Provide(user.NewUserController),
	fx.Provide(user.NewAccountController),
	fx.Provide(user.NewBotController),
	fx.Provide(user.NewPostMessageController),
	fx.Provide(user.NewImmediatePostController),
	fx.Provide(user.NewSchedulePostController),
	fx.Provide(user.NewRegularPostController),
	fx.Provide(user.NewRegularTimingController),
	fx.Provide(user.NewApiPostController),

	// api controller
	fx.Provide(api.NewPostMessageApiController),
	fx.Provide(api.NewApiPostApiController),
)
