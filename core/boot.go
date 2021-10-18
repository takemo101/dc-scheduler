package core

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewConfig),
	fx.Provide(NewLogger),
	fx.Provide(NewDatabase),
	fx.Provide(NewTemplateEngine),
	fx.Provide(NewPath),
	fx.Provide(NewMailFactory),
	fx.Provide(NewApplication),
	fx.Provide(NewSessionStore),
)
