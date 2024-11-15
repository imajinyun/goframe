package dcs

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/imajinyun/goframe"
	"github.com/imajinyun/goframe/contract"
)

var errDcsServiceParams = errors.New("distributed service params error")

type DcsService struct {
	container goframe.IContainer
}

func NewDcsService(params ...any) (any, error) {
	if len(params) != 1 {
		return nil, errDcsServiceParams
	}

	container := params[0].(goframe.IContainer)

	return &DcsService{container: container}, nil
}

func (s *DcsService) Select(name string, id string, hold time.Duration) (string, error) {
	appsvc := s.container.MustMake(contract.AppKey).(contract.IApp)
	dir := appsvc.RunDir()
	file := filepath.Join(dir, "dcs_"+name)

	lock, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		return "", err
	}

	if err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		byt, err := io.ReadAll(lock)
		if err != nil {
			return "", err
		}
		return string(byt), nil
	}

	go func() {
		defer func() {
			syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			lock.Close()
			os.Remove(file)
		}()

		timer := time.NewTimer(hold)
		<-timer.C
	}()

	if _, err := lock.WriteString(id); err != nil {
		return "", err
	}

	return id, nil
}
