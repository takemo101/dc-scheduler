package middleware

import (
	"github.com/takemo101/dc-scheduler/core/contract"
	"go.uber.org/fx"
)

// Module export
var Module = fx.Options(
	fx.Provide(NewMethodOverride),
	fx.Provide(NewSecure),
	fx.Provide(NewCsrf),
	fx.Provide(NewCors),
	fx.Provide(NewBasicAuth),
	fx.Provide(NewSessionAdminAuth),
	fx.Provide(NewSessionAdminRole),
	fx.Provide(NewSessionUserAuth),
	fx.Provide(NewViewRender),
	fx.Provide(NewContextValueInit),
	fx.Provide(NewSecurityKey),
	fx.Provide(NewMiddleware),
)

// Middlewares is slice
type Middlewares []contract.Middleware

// NewMiddleware is setup new middlewares
func NewMiddleware(
	methodOverride MethodOverride,
	// secure Secure,
	value ContextValueInit,

) Middlewares {
	return Middlewares{
		methodOverride,
		// secure,
		value,
	}
}

// Setup all the middleware
func (middlewares Middlewares) Setup() {
	for _, middleware := range middlewares {
		middleware.Setup()
	}
}
