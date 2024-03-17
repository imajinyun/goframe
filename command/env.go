package command

import (
	"log"

	"github.com/imajinyun/goframe/cobra"
	"github.com/imajinyun/goframe/contract"
	"github.com/imajinyun/goframe/util"
)

var envCommand = &cobra.Command{
	Use:   "env",
	Short: "get env variable",
	Long:  "get env variable",
	Run: func(cmd *cobra.Command, args []string) {
		container := cmd.GetContainer()
		envsvc := container.MustMake(contract.EnvKey).(contract.IEnv)
		log.Println("environment:", envsvc.Env())
	},
	Args: cobra.NoArgs,
}

var envListCommand = &cobra.Command{
	Use:   "env",
	Short: "get env list",
	Long:  "get env list",
	Run: func(cmd *cobra.Command, args []string) {
		container := cmd.GetContainer()
		envsvc := container.MustMake(contract.EnvKey).(contract.IEnv)

		list := [][]string{}
		for k, v := range envsvc.All() {
			list = append(list, []string{k, v})
		}
		util.Pretty(list)
	},
}

func initEnvCommand() *cobra.Command {
	envCommand.AddCommand(envListCommand)

	return envCommand
}
