package main

import (
	"os"
	"path"

	"github.com/takemo101/dc-scheduler/app"
	"github.com/takemo101/dc-scheduler/cli/cmd"
	"github.com/takemo101/dc-scheduler/cli/kernel"
	"github.com/takemo101/dc-scheduler/core/contract"
	"github.com/takemo101/dc-scheduler/database"
	"github.com/takemo101/dc-scheduler/pkg"
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
	current := os.Getenv("APP_CURRENT")
	if current == "" {
		current, _ = os.Getwd()
	}

	// boot cobra application
	kernel.Run(
		kernel.CLIOptions{
			ConfigPath:       path.Join(current, "config.yml"),
			CurrentDirectory: current,
			CommandOptions: cmd.CommandOptions{
				Models:     database.Models,
				Migrations: database.Migrations,
			},
			CLIBooterConstructor: NewCLIBooter,
			FXOption: fx.Options(
				app.Module,
				pkg.Module,
				cmd.Module,
			),
		},
	)
}
