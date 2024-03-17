package command

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sevlyar/go-daemon"

	"github.com/imajinyun/goframe/cobra"
	"github.com/imajinyun/goframe/contract"
	"github.com/imajinyun/goframe/util"
)

var cronDaemon = false

var cronCommand = &cobra.Command{
	Use:   "cron",
	Short: "cron command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}

		return nil
	},
}

var cronListCommand = &cobra.Command{
	Use:   "list",
	Short: "cron list command",
	RunE: func(cmd *cobra.Command, args []string) error {
		cs := cmd.Root().CronSpecs
		ps := [][]string{}
		for _, v := range cs {
			ps = append(ps, []string{v.Type, v.Spec, v.Name, v.Cmd.Short, v.Cmd.Short})
		}

		return nil
	},
}

var cronStartCommand = &cobra.Command{
	Use:   "start",
	Short: "start cron daemon short description",
	Long:  "start cron daemon long description",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appsvc := container.MustMake(contract.AppKey).(contract.IApp)
		workDir := appsvc.WorkDir()
		pidFileName := filepath.Join(appsvc.RunDir(), "cron.pid")
		logFileName := filepath.Join(appsvc.LogDir(), "cron.log")

		if cronDaemon {
			ctx := &daemon.Context{
				PidFileName: pidFileName,
				PidFilePerm: 0o664,
				LogFileName: logFileName,
				LogFilePerm: 0o640,
				WorkDir:     workDir,
				Umask:       0o27,
				Args:        []string{"", "cron", "start", "--daemon=true"},
			}
			p, err := ctx.Reborn()
			if err != nil {
				return err
			}

			if p != nil {
				log.Println("cron serve started, pid:", p.Pid)
				log.Println("log file:", logFileName)
				return nil
			}

			defer ctx.Release()
			log.Println("daemon started")
			// gspt.SetProcTitle("gogin cron")
			cmd.Root().Cron.Run()
			return nil
		}

		log.Println("start cron job")
		content := strconv.Itoa(os.Getegid())
		log.Println("[PID]", content)
		if err := os.WriteFile(pidFileName, []byte(content), 0o664); err != nil {
			return err
		}

		// gspt.SetProcTitle("gogin cron")
		cmd.Root().Cron.Run()
		return nil
	},
}

var cronRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "restart cron daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appsvc := container.MustMake(contract.AppKey).(contract.IApp)
		pidFileName := filepath.Join(appsvc.RunDir(), "cron.pid")
		content, err := os.ReadFile(pidFileName)
		if err != nil {
			return err
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}

			if util.IsProcessExist(pid) {
				if err := util.KillProcess(pid); err != nil {
					return err
				}
				for i := 0; i < 10; i++ {
					if util.IsProcessExist(pid) == false {
						break
					}
					time.Sleep(1 * time.Second)
				}
				log.Println("kill process:", pid)
			}
		}
		cronDaemon = false

		return cronStartCommand.RunE(cmd, args)
	},
}

func initCronCommand() *cobra.Command {
	cronStartCommand.Flags().BoolVarP(&cronDaemon, "daemon", "d", false, "start serve daemon")
	cronCommand.AddCommand(cronStartCommand)
	cronCommand.AddCommand(cronRestartCommand)

	return cronCommand
}
