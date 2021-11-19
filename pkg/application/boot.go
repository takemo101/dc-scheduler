package application

import (
	admin "github.com/takemo101/dc-scheduler/pkg/application/admin"
	user "github.com/takemo101/dc-scheduler/pkg/application/user"
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

	// --- User UseCase by admin --
	fx.Provide(admin.NewUserStoreUseCase),
	fx.Provide(admin.NewUserUpdateUseCase),
	fx.Provide(admin.NewUserDeleteUseCase),
	fx.Provide(admin.NewUserDetailUseCase),
	fx.Provide(admin.NewUserSearchUseCase),

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
	fx.Provide(admin.NewApiPostSearchUseCase),
	fx.Provide(admin.NewApiPostSendUseCase),

	// --- User UseCase by user --
	fx.Provide(user.NewUserRegistUseCase),
	fx.Provide(user.NewUserActivationUseCase),
	fx.Provide(user.NewMyAccountDetailUseCase),
	fx.Provide(user.NewMyAccountUpdateUseCase),
	fx.Provide(user.NewUserLoginUseCase),
	fx.Provide(user.NewUserLogoutUseCase),

	// --- Bot UseCase by user ---
	fx.Provide(user.NewBotStoreUseCase),
	fx.Provide(user.NewBotUpdateUseCase),
	fx.Provide(user.NewBotDeleteUseCase),
	fx.Provide(user.NewBotDetailUseCase),
	fx.Provide(user.NewBotSearchUseCase),

	// --- Message UseCase by user ---
	fx.Provide(user.NewPostMessageDeleteUseCase),
	fx.Provide(user.NewSentMessageHistoryUseCase),
	fx.Provide(user.NewImmediatePostSearchUseCase),
	fx.Provide(user.NewImmediatePostStoreUseCase),
	fx.Provide(user.NewSchedulePostSearchUseCase),
	fx.Provide(user.NewSchedulePostStoreUseCase),
	fx.Provide(user.NewSchedulePostEditFormUseCase),
	fx.Provide(user.NewSchedulePostUpdateUseCase),
	fx.Provide(user.NewRegularPostSearchUseCase),
	fx.Provide(user.NewRegularPostStoreUseCase),
	fx.Provide(user.NewRegularPostEditFormUseCase),
	fx.Provide(user.NewRegularPostUpdateUseCase),
	fx.Provide(user.NewRegularTimingAddUseCase),
	fx.Provide(user.NewRegularTimingRemoveUseCase),
	fx.Provide(user.NewApiPostSearchUseCase),
	fx.Provide(user.NewApiPostStoreUseCase),
)
