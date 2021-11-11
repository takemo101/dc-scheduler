package application

import (
	admin "github.com/takemo101/dc-scheduler/pkg/application/admin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	// --- Admin UseCase by admin --
	fx.Provide(admin.NewAdminStoreUseCase),
	fx.Provide(admin.NewAdminUpdateUseCase),
	fx.Provide(admin.NewAdminDeleteUseCase),
	fx.Provide(admin.NewAdminDetailUseCase),
	fx.Provide(admin.NewAdminSearchUseCase),
	fx.Provide(admin.NewMyAccountDetailUseCase),
	fx.Provide(admin.NewMyAccountUpdateUseCase),
	fx.Provide(admin.NewAdminLoginUseCase),
	fx.Provide(admin.NewAdminLogoutUseCase),

	// --- Bot UseCase by admin ---
	fx.Provide(admin.NewBotUpdateUseCase),
	fx.Provide(admin.NewBotDeleteUseCase),
	fx.Provide(admin.NewBotDetailUseCase),
	fx.Provide(admin.NewBotSearchUseCase),

	// --- Message UseCase by admin ---
	fx.Provide(admin.NewPostMessageSearchUseCase),
	fx.Provide(admin.NewPostMessageDeleteUseCase),
	fx.Provide(admin.NewSentMessageHistoryUseCase),
	fx.Provide(admin.NewImmediatePostSearchUseCase),
	fx.Provide(admin.NewSchedulePostSearchUseCase),
	fx.Provide(admin.NewSchedulePostEditFormUseCase),
	fx.Provide(admin.NewSchedulePostUpdateUseCase),
	fx.Provide(admin.NewSchedulePostSendUseCase),
	fx.Provide(admin.NewRegularPostSearchUseCase),
	fx.Provide(admin.NewRegularPostEditFormUseCase),
	fx.Provide(admin.NewRegularPostUpdateUseCase),
	fx.Provide(admin.NewRegularTimingAddUseCase),
	fx.Provide(admin.NewRegularTimingRemoveUseCase),
	fx.Provide(admin.NewRegularPostSendUseCase),
)
