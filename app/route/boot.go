package route

import (
	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/core/contract"
	"go.uber.org/fx"
)

// Module export
var Module = fx.Options(
	fx.Provide(NewApiRoute),
	fx.Provide(NewAdminRoute),
	fx.Provide(NewRoute),
)

// Routes is slice
type Routes []contract.Route

// NewRoute is setup routes
func NewRoute(
	config core.Config,
	api ApiRoute,
	admin AdminRoute,
) Routes {
	return Routes{
		api,
		admin,
	}
}

// Setup all the route
func (routes Routes) Setup() {
	for _, route := range routes {
		route.Setup()
	}
}
