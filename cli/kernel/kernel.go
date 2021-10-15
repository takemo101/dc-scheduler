package kernel

import (
	"github.com/takemo101/dc-scheduler/cli/cmd"
	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/core/contract"
	"go.uber.org/fx"
)

// AppOptions app boot options
type CLIOptions struct {
	ConfigPath           contract.ConfigPath
	CommandOptions       cmd.CommandOptions
	CLIBooterConstructor interface{}
	FXOption             fx.Option
}

// boot is initialize cli
func boot(
	logger core.Logger,
	database core.Database,
	config core.Config,
	root cmd.RootCommand,
	booter contract.CLIBooter,
) {
	sql, err := database.DB()
	if err != nil {
		logger.Info("database connection sql failed : %v", err)
	}

	defer sql.Close()

	logger.Info("-- start cli --")

	sql.SetMaxOpenConns(config.DB.Connection.Max)

	booter.CLIBoot()
	root.Cmd.Execute()

	logger.Info("-- stop cli --")
}

// app run
func Run(options CLIOptions) {
	opts := []fx.Option{
		options.FXOption,
		core.Module,
		fx.Provide(
			func() contract.ConfigPath {
				return options.ConfigPath
			},
			func() cmd.CommandOptions {
				return options.CommandOptions
			},
			options.CLIBooterConstructor,
		),
		fx.Invoke(boot),
	}
	fx.New(opts...).Done()
}
