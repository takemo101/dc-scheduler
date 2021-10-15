package boot

import (
	"testing"

	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/core/contract"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// TestOptions app boot options
type TestOptions struct {
	ConfigPath contract.ConfigPath
	FXOption   fx.Option
}

// Testing test func
func Testing(t *testing.T, options TestOptions, tests ...interface{}) {
	fxtest.New(
		t,
		fx.Provide(func() contract.ConfigPath {
			return options.ConfigPath
		}),
		options.FXOption,
		core.Module,
		fx.Invoke(tests...),
	).Done()
}
