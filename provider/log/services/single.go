package services

import (
	"os"
	"path/filepath"

	"github.com/imajinyun/goframe"

	"github.com/pkg/errors"

	"github.com/imajinyun/goframe/contract"
	"github.com/imajinyun/goframe/util"
)

type SingleLog struct {
	Log

	dir    string
	file   string
	writer *os.File
}

func NewSingleLog(params ...any) (any, error) {
	container := params[0].(goframe.IContainer)
	level := params[1].(contract.LogLevel)
	handler := params[2].(contract.Handler)
	formatter := params[3].(contract.Formatter)

	appsvc := container.MustMake(contract.AppKey).(contract.IApp)
	etcsvc := container.MustMake(contract.EtcKey).(contract.IEtc)

	log := &SingleLog{}
	log.SetLevel(level)
	log.SetHandler(handler)
	log.SetFormatter(formatter)

	dir := appsvc.LogDir()
	if etcsvc.Exist("log.dir") {
		dir = etcsvc.GetString("log.dir")
	}
	log.dir = dir
	if !util.Exist(dir) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}

	log.file = "gogin.log"
	if etcsvc.Exist("log.file") {
		log.file = etcsvc.GetString("log.file")
	}

	writer, err := os.OpenFile(filepath.Join(log.dir, log.file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return nil, errors.Wrap(err, "open log file error")
	}

	log.SetWriter(writer)
	log.container = container

	return log, nil
}
