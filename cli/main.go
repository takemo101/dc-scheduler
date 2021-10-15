package main

import (
	"github.com/takemo101/dc-scheduler/app"
	"github.com/takemo101/dc-scheduler/cli/cmd"
	"github.com/takemo101/dc-scheduler/cli/kernel"
	"github.com/takemo101/dc-scheduler/core/contract"
	"github.com/takemo101/dc-scheduler/database"
	"go.uber.org/fx"
)

// CLIBooter is command module root struct
type CLIBooter struct {
	commands cmd.Commands
}

// CLIBooter all setup
func (booter CLIBooter) CLIBoot() {
	booter.commands.Setup()
}

// NewAppBooter app create
func NewCLIBooter(
	commands cmd.Commands,
) contract.CLIBooter {
	return CLIBooter{
		commands,
	}
}

func main() {
	// boot cobra application
	kernel.Run(
		kernel.CLIOptions{
			ConfigPath: "config.yml",
			CommandOptions: cmd.CommandOptions{
				Models:     database.Models,
				Migrations: database.Migrations,
			},
			CLIBooterConstructor: NewCLIBooter,
			FXOption: fx.Options(
				app.Module,
				cmd.Module,
			),
		},
	)
}
