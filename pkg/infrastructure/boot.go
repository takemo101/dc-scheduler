package infrastructure

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewAdminAuthContext),
	fx.Provide(NewAdminRepository),
	fx.Provide(NewAdminQuery),
)
