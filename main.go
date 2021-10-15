package main

import (
	"github.com/takemo101/dc-scheduler/app"
	"github.com/takemo101/dc-scheduler/app/middleware"
	"github.com/takemo101/dc-scheduler/app/route"
	"github.com/takemo101/dc-scheduler/boot"
	"github.com/takemo101/dc-scheduler/core/contract"
	"go.uber.org/fx"
)

// AppBooter is module root struct
type AppBooter struct {
	routes      route.Routes
	middlewares middleware.Middlewares
}

// AppBoot all setup
func (booter AppBooter) AppBoot() {
	booter.middlewares.Setup()
	booter.routes.Setup()
}

// NewAppBooter app create
func NewAppBooter(
	routes route.Routes,
	middlewares middleware.Middlewares,
) contract.AppBooter {
	return AppBooter{
		routes:      routes,
		middlewares: middlewares,
	}
}

func main() {

	// boot gin application
	boot.Run(
		boot.AppOptions{
			ConfigPath:           "config.yml",
			AppBooterConstructor: NewAppBooter,
			FXOption: fx.Options(
				app.Module,
			),
		},
	)
}
