package command

import (
	"log"
	"path/filepath"

	"github.com/swaggo/swag/gen"

	"github.com/imajinyun/goframe/cobra"
	"github.com/imajinyun/goframe/contract"
)

var swaggerCommand = &cobra.Command{
	Use:     "swagger",
	Short:   "generate swagger file",
	Long:    "generate swagger file",
	Aliases: []string{"swg"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

var swaggerGenCommand = &cobra.Command{
	Use:   "gen",
	Short: "generate swagger file",
	Long:  "generate swagger file",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appsvc := container.MustMake(contract.AppKey).(contract.IApp)

		cfg := &gen.Config{
			SearchDir:          filepath.Join(appsvc.AppDir(), "http"),
			Excludes:           "",
			OutputDir:          filepath.Join(appsvc.WorkDir(), "docs/swagger"),
			OutputTypes:        []string{"go", "json", "yaml"},
			MainAPIFile:        "swagger.go",
			PropNamingStrategy: "camel",
			ParseVendor:        false,
			ParseDependency:    0,
			ParseInternal:      false,
			MarkdownFilesDir:   "",
			GeneratedTime:      true,
		}

		if err := gen.New().Build(cfg); err != nil {
			log.Printf("swagger generate error: %v", err)
			return err
		}

		return nil
	},
}

func initSwaggerCommand() *cobra.Command {
	swaggerCommand.AddCommand(swaggerGenCommand)

	return swaggerCommand
}
