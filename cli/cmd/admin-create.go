package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/application"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// AdminCreateCommand is struct
type AdminCreateCommand struct {
	logger  core.Logger
	root    RootCommand
	usecase application.AdminStoreUseCase
}

// Setup is setup route
func (c AdminCreateCommand) Setup() {
	c.logger.Info("setup admin_create-command")

	var name, email, pass string

	cmd := &cobra.Command{
		Use:   "admin:create",
		Short: "create admin",
		Run: func(cmd *cobra.Command, args []string) {
			input := application.AdminStoreInput{
				Name:     name,
				Email:    email,
				Password: pass,
				Role:     string(domain.AdminRoleSystem),
			}

			id, err := c.usecase.Execute(input)
			if err != nil {
				c.logger.Error(err)
				fmt.Println(err)
				return
			}
			fmt.Println(fmt.Sprintf("success create user is ID[%d]", id))
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "admin", "create admin name")
	cmd.Flags().StringVarP(&email, "email", "e", "admin@example.com", "create admin email")
	cmd.Flags().StringVarP(&pass, "pass", "p", "admin", "create admin pass")

	c.root.Cmd.AddCommand(cmd)
}

// NewAdminCreateCommand create new admin create command
func NewAdminCreateCommand(
	root RootCommand,
	logger core.Logger,
	usecase application.AdminStoreUseCase,
) AdminCreateCommand {
	return AdminCreateCommand{
		root:    root,
		logger:  logger,
		usecase: usecase,
	}
}
