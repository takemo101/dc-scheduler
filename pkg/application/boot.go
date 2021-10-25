package application

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	// --- Admin UseCase --
	fx.Provide(NewAdminStoreUseCase),
	fx.Provide(NewAdminUpdateUseCase),
	fx.Provide(NewAdminDeleteUseCase),
	fx.Provide(NewAdminDetailUseCase),
	fx.Provide(NewAdminSearchUseCase),
	fx.Provide(NewMyAccountDetailUseCase),
	fx.Provide(NewMyAccountUpdateUseCase),
	fx.Provide(NewAdminLoginUseCase),
	fx.Provide(NewAdminLogoutUseCase),
	// --- Bot UseCase ---
	fx.Provide(NewBotStoreUseCase),
	fx.Provide(NewBotUpdateUseCase),
	fx.Provide(NewBotDeleteUseCase),
	fx.Provide(NewBotDetailUseCase),
	fx.Provide(NewBotSearchUseCase),
	// --- Message UseCase ---
	fx.Provide(NewPostMessageSearchUseCase),
	fx.Provide(NewPostMessageCreateFormUseCase),
	fx.Provide(NewPostMessageDeleteUseCase),
	fx.Provide(NewSentMessageHistoryUseCase),
	fx.Provide(NewImmediatePostSearchUseCase),
	fx.Provide(NewImmediatePostStoreUseCase),
	fx.Provide(NewSchedulePostSearchUseCase),
	fx.Provide(NewSchedulePostStoreUseCase),
	fx.Provide(NewSchedulePostEditFormUseCase),
	fx.Provide(NewSchedulePostUpdateUseCase),
	fx.Provide(NewSchedulePostSendUseCase),
)
