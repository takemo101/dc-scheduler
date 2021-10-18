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
	// --- Message UseCase ---
)
