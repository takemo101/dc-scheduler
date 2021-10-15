package cmd

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/spf13/cobra"
	"github.com/takemo101/dc-scheduler/core"
)

type MigrationList []interface{}

// MigrateCommand is struct
type MigrateCommand struct {
	logger     core.Logger
	root       RootCommand
	db         core.Database
	migrations []*gormigrate.Migration
}

// Setup is setup command
func (c MigrateCommand) Setup() {
	c.logger.Info("setup migrate-command")

	var process string

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "migration migrate up",
		Run: func(cmd *cobra.Command, args []string) {

			m := gormigrate.New(c.db.GormDB, gormigrate.DefaultOptions, c.migrations)

			if len(process) > 0 {
				if err := m.MigrateTo(process); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("finish migrate to id:" + process)
				}
			} else {
				if err := m.Migrate(); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("finish migrate")
				}
			}
		},
	}

	cmd.Flags().StringVarP(&process, "process", "p", "", "migrate process name id or ...")

	c.root.Cmd.AddCommand(cmd)
}

// NewMigrateCommand create migrate command
func NewMigrateCommand(
	root RootCommand,
	logger core.Logger,
	db core.Database,
	options CommandOptions,
) MigrateCommand {
	return MigrateCommand{
		root:       root,
		logger:     logger,
		db:         db,
		migrations: options.Migrations,
	}
}
