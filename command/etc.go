package command

import (
	"log"

	"github.com/kr/pretty"

	"github.com/imajinyun/goframe/cobra"
	"github.com/imajinyun/goframe/contract"
)

var etcCommand = &cobra.Command{
	Use:   "etc",
	Short: "get etc info",
	Long:  "get etc info",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

var etcGetCommand = &cobra.Command{
	Use:   "get",
	Short: "get given etc value",
	Long:  "get given etc value",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			log.Println("param is error")
			return nil
		}
		container := cmd.GetContainer()
		etcsvc := container.MustMake(contract.EtcKey).(contract.IEtc)
		key := args[0]
		val := etcsvc.Get(key)
		if val == nil {
			return nil
		}
		log.Printf("%# v\n", pretty.Formatter(val))

		return nil
	},
}

func initEtcCommand() *cobra.Command {
	etcCommand.AddCommand(etcGetCommand)

	return etcCommand
}
