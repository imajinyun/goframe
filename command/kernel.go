package command

import (
	"github.com/robfig/cron/v3"

	"github.com/imajinyun/goframe/cobra"
)

func AddKernelCommand(root *cobra.Command) {
	setCronCommand(root)

	root.AddCommand(initAppCommand())
	root.AddCommand(initEnvCommand())
	root.AddCommand(initEtcCommand())
	root.AddCommand(initCronCommand())
	root.AddCommand(initNewCommand())
	root.AddCommand(initSwaggerCommand())
	root.AddCommand(initProviderCommand())
}

func setCronCommand(root *cobra.Command) {
	if root.Cron == nil {
		flag := cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor
		root.Cron = cron.New(cron.WithParser(cron.NewParser(flag)))
		root.CronSpecs = []cobra.CronSpec{}
	}
}
