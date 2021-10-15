package controller

import (
	admin "github.com/takemo101/dc-scheduler/app/controller/admin"
	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	// admin controller
	fx.Provide(admin.NewDashboardController),
	fx.Provide(admin.NewSessionAuthController),
)
