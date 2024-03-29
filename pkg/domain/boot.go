package domain

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewAdminService),
	fx.Provide(NewUserService),
	fx.Provide(NewBotService),
)
