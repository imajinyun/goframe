package services

import (
	"os"

	"github.com/imajinyun/goframe/contract"

	"github.com/imajinyun/goframe"
)

type ConsoleLog struct {
	Log
}

func NewConsoleLog(params ...any) (any, error) {
	container := params[0].(goframe.IContainer)
	level := params[1].(contract.LogLevel)
	handler := params[2].(contract.Handler)
	formatter := params[3].(contract.Formatter)

	log := &ConsoleLog{}
	log.SetLevel(level)
	log.SetWriter(os.Stdout)
	log.SetHandler(handler)
	log.SetFormatter(formatter)
	log.container = container

	return log, nil
}
