package pkg

import (
	"github.com/takemo101/dc-scheduler/pkg/application"
	"github.com/takemo101/dc-scheduler/pkg/infrastructure"
	"go.uber.org/fx"
)

var Module = fx.Options(
	application.Module,
	infrastructure.Module,
)
