package cmd

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/takemo101/dc-scheduler/core/contract"
	"go.uber.org/fx"
)

// Module export
var Module = fx.Options(
	fx.Provide(NewStorageLinkCommand),
	fx.Provide(NewMigrateCommand),
	fx.Provide(NewRollbackCommand),
	fx.Provide(NewAutoMigrateCommand),
	fx.Provide(NewCommandRoot),
	fx.Provide(NewMailCommand),
	fx.Provide(NewCommand),
)

// CommandOptions is base commands options
type CommandOptions struct {
	Migrations []*gormigrate.Migration
	Models     []interface{}
}

// Commands is commands slice
type Commands []contract.Command

// NewCommand is setup command
func NewCommand(
	storageLinkCommand StorageLinkCommand,
	migrateCommand MigrateCommand,
	rollbackCommand RollbackCommand,
	autoMigrateCommand AutoMigrateCommand,
	mail MailCommand,
) Commands {
	return Commands{
		storageLinkCommand,
		migrateCommand,
		rollbackCommand,
		autoMigrateCommand,
		mail,
	}
}

// Setup all the command
func (commands Commands) Setup() {
	for _, cmd := range commands {
		cmd.Setup()
	}
}
