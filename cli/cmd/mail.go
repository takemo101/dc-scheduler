package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/core"
)

// MailCommand is struct
type MailCommand struct {
	logger       core.Logger
	root         RootCommand
	mailTemplate helper.MailTemplate
}

// Setup is setup command
func (c MailCommand) Setup() {
	c.logger.Info("setup mail_send-command")

	var to, message string

	cmd := &cobra.Command{
		Use:   "mail:send",
		Short: "mail send",
		Run: func(cmd *cobra.Command, args []string) {
			mail, _, creatorErr := c.mailTemplate.GetTextMailCreatorByKey("test", core.BindData{
				"message": message,
			})

			if creatorErr != nil {
				fmt.Println(creatorErr)
			} else {
				mail.To(to)
				if err := mail.Send(); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(fmt.Sprintf("success send mail to %s", to))
				}
			}
		},
	}

	cmd.Flags().StringVarP(&to, "to", "t", "xxx@xxx.com", "create mail to address")
	cmd.Flags().StringVarP(&message, "message", "m", "message", "create mail message")

	c.root.Cmd.AddCommand(cmd)
}

// NewMailCommand create new mail send command
func NewMailCommand(
	root RootCommand,
	logger core.Logger,
	mailTemplate helper.MailTemplate,
) MailCommand {
	return MailCommand{
		root:         root,
		logger:       logger,
		mailTemplate: mailTemplate,
	}
}
