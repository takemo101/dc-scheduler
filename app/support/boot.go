package support

import "go.uber.org/fx"

// Module export
var Module = fx.Options(
	fx.Provide(NewRequestValue),
)
