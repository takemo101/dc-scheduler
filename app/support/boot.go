package support

import "go.uber.org/fx"

// Module サポートモジュールDI
var Module = fx.Options(
	fx.Provide(NewContextValue),
	fx.Provide(NewToastrMessage),
)
