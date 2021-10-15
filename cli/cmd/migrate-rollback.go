package cmd

import (
	"fmt"
	"strings"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/spf13/cobra"
	"github.com/takemo101/dc-scheduler/core"
)

// RollbackCommand is struct
type RollbackCommand struct {
	logger     core.Logger
	root       RootCommand
	db         core.Database
	migrations []*gormigrate.Migration
}

// Setup is setup command
func (c RollbackCommand) Setup() {
	c.logger.Info("setup migrate:rollback-command")

	var process string

	cmd := &cobra.Command{
		Use:   "migrate:rollback",
		Short: "migration migrate down",
		Run: func(cmd *cobra.Command, args []string) {

			m := gormigrate.New(c.db.GormDB, gormigrate.DefaultOptions, c.migrations)

			switch n := strings.ToLower(process); n {
			case "step":
				if err := m.RollbackLast(); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("finish rollback step")
				}
			case "all":
				for {
					if err := m.RollbackLast(); err != nil {
						break
					}
				}
				fmt.Println("finish rollback all")
			default:
				if len(n) > 0 {
					if err := m.RollbackTo(n); err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("finish rollback to id:" + n)
					}
				} else {

				}
			}
		},
	}

	cmd.Flags().StringVarP(&process, "process", "p", "step", "migrate process name step or all or id")

	c.root.Cmd.AddCommand(cmd)
}

// NewRollbackCommand create migrate command
func NewRollbackCommand(
	root RootCommand,
	logger core.Logger,
	db core.Database,
	options CommandOptions,
) RollbackCommand {
	return RollbackCommand{
		root:       root,
		logger:     logger,
		db:         db,
		migrations: options.Migrations,
	}
}
