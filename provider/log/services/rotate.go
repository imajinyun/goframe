package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/imajinyun/goframe"

	rotatelog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"

	"github.com/imajinyun/goframe/contract"
	"github.com/imajinyun/goframe/util"
)

type RotateLog struct {
	Log

	path string
	file string
}

func NewRotateLog(params ...any) (any, error) {
	container := params[0].(goframe.IContainer)
	level := params[1].(contract.LogLevel)
	handler := params[2].(contract.Handler)
	formatter := params[2].(contract.Formatter)

	appsvc := container.MustMake(contract.AppKey).(contract.IApp)
	etcsvc := container.MustMake(contract.EtcKey).(contract.IEtc)

	path := appsvc.LogDir()
	if etcsvc.Exist("log.folder") {
		path = etcsvc.GetString("log.folder")
	}

	if !util.Exist(path) {
		os.MkdirAll(path, os.ModePerm)
	}

	file := "rotate.log"
	if etcsvc.Exist("log.file") {
		file = etcsvc.GetString("log.file")
	}

	dfmt := "%Y%m%d%H"
	if etcsvc.Exist("log.date_format") {
		dfmt = etcsvc.GetString("log.date_format")
	}

	name := rotatelog.WithLinkName(filepath.Join(path, file))
	options := []rotatelog.Option{name}

	if etcsvc.Exist("log.rotate_time") {
		if rotateTimeParse, err := time.ParseDuration(etcsvc.GetString("log.rotate_time")); err == nil {
			options = append(options, rotatelog.WithRotationTime(rotateTimeParse))
		}
	}

	log := &RotateLog{}
	log.SetLevel(level)
	log.SetHandler(handler)
	log.SetFormatter(formatter)
	log.path = path
	log.file = file

	writer, err := rotatelog.New(fmt.Sprintf("%s.%s", filepath.Join(log.path, log.file), dfmt), options...)
	if err != nil {
		return nil, errors.Wrap(err, "new rotatelogs error")
	}
	log.SetWriter(writer)
	log.container = container

	return log, nil
}
