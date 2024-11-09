package cobra

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/imajinyun/goframe/contract"
)

func (c *Command) AddDcsCronCommand(name string, spec string, cmd *Command, hold time.Duration) {
	root := c.Root()

	if root.Cron == nil {
		flag := cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor
		root.Cron = cron.New(cron.WithParser(cron.NewParser(flag)))
		root.CronSpecs = []CronSpec{}
	}

	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Cmd:  cmd,
		Type: "distributed-cron",
		Spec: spec,
		Name: name,
	})

	appsvc := root.GetContainer().MustMake(contract.AppKey).(contract.IApp)
	dstsvc := root.GetContainer().MustMake(contract.DcsKey).(contract.IDcs)
	appid := appsvc.AppID()

	var cronCmd Command
	ctx := root.Context()
	cronCmd = *cmd
	cronCmd.args = []string{}
	cronCmd.SetParentNull()

	root.Cron.AddFunc(spec, func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		sid, err := dstsvc.Select(name, appid, hold)
		if err != nil {
			return
		}

		if sid != appid {
			return
		}

		if err := cronCmd.ExecuteContext(ctx); err != nil {
			log.Println(err)
		}
	})
}
