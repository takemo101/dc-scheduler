package app

import (
	"github.com/takemo101/dc-scheduler/app/controller"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/middleware"
	"github.com/takemo101/dc-scheduler/app/route"
	"github.com/takemo101/dc-scheduler/app/support"
	"go.uber.org/fx"
)

var Module = fx.Options(
	middleware.Module,
	helper.Module,
	support.Module,
	route.Module,
	controller.Module,
)
