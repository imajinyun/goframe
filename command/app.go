package command

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/imajinyun/goframe/cobra"
	"github.com/imajinyun/goframe/contract"
	"github.com/imajinyun/goframe/util"
)

var appAddr string = ""
var appDaemon bool = false

var appCommand = &cobra.Command{
	Use:   "app",
	Short: "App contains some useful commands",
	Long:  "App contains some useful commands",
	RunE: func(c *cobra.Command, args []string) error {
		c.Help()

		return nil
	},
}

var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "Start a golang web server",
	Long:  "Start a golang web server",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		kernelsvc := container.MustMake(contract.KernelKey).(contract.IKernel)
		handler := kernelsvc.HttpHandler()

		if appAddr == "" {
			envsvc := container.MustMake(contract.EnvKey).(contract.IEnv)
			if len(envsvc.Get("ADDRESS")) > 0 {
				appAddr = envsvc.Get("ADDRESS")
			} else {
				etcsvc := container.MustMake(contract.EtcKey).(contract.IEtc)
				if etcsvc.Exist("app.address") {
					appAddr = etcsvc.GetString("app.address")
				} else {
					appAddr = ":8080"
				}
			}
		}
		log.Printf("Server started at: http://127.0.0.1%v\n", appAddr)

		srv := &http.Server{Addr: appAddr, Handler: handler}
		go func() {
			srv.ListenAndServe()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server shutdown: ", err)
		}

		return nil
	},
}

var appStateCommand = &cobra.Command{
	Use:   "state",
	Short: "Get started app process id",
	Long:  "Get started app process id",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		appsvc := container.MustMake(contract.AppKey).(contract.IApp)
		fpfile := filepath.Join(appsvc.RunDir(), "app.pid")
		content, err := os.ReadFile(fpfile)
		if err != nil {
			return err
		}

		if content != nil && len(content) > 0 {
			pid, err := strconv.Atoi(string(content))
			if err != nil {
				return err
			}

			if util.IsProcessExist(pid) {
				log.Println("app pid:", pid)
				return nil
			}
		}

		return nil
	},
}

func initAppCommand() *cobra.Command {
	appCommand.AddCommand(appStartCommand)
	appCommand.AddCommand(appStateCommand)

	return appCommand
}
