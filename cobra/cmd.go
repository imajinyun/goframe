package cobra

import (
	"log"

	"github.com/imajinyun/goframe"
)

type CronSpec struct {
	Cmd  *Command
	Type string
	Spec string
	Name string
}

func (c *Command) SetContainer(container goframe.IContainer) {
	c.container = container
}

func (c *Command) GetContainer() goframe.IContainer {
	return c.Root().container
}

func (c *Command) SetParentNull() {
	c.parent = nil
}

func (c *Command) AddCronCommand(spec string, cmd *Command) {
	root := c.Root()
	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Cmd:  cmd,
		Type: "normal-cron",
		Spec: spec,
	})
	root.Cron.AddFunc(spec, func() {
		var cronCmd Command
		ctx := root.Context()
		cronCmd = *cmd
		cronCmd.args = []string{}
		cronCmd.SetParentNull()
		cronCmd.SetContainer(root.GetContainer())

		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		if err := cronCmd.ExecuteContext(ctx); err != nil {
			log.Println(err)
		}
	})
}
