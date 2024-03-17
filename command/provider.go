package command

import (
	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/cobra"
)

var providerCommand = &cobra.Command{
	Use:   "provider",
	Short: "",
}

var providerListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all providers",
	Long:  "list all providers",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer().(goframe.IContainer)
		list := (container.(*goframe.Container)).Providers()

		for _, item := range list {
			println(item)
		}

		return nil
	},
}

var providerNewCommand = &cobra.Command{
	Use:     "new",
	Short:   "new a provider",
	Aliases: []string{"add", "create", "init"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func initProviderCommand() *cobra.Command {
	providerCommand.AddCommand(providerListCommand)

	return providerCommand
}
