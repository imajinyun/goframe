package util

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func GetExecDir() string {
	file, err := os.Getwd()
	if err == nil {
		return file + "/"
	}

	return ""
}

func GetRootDir() (string, error) {
	exe, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(exe)
	for {
		if _, err := os.Stat(filepath.Join(dir, ".go-root")); err == nil {
			return dir, nil
		}

		pdir := filepath.Dir(dir)
		if pdir == dir {
			break
		}
		dir = pdir
	}

	return "", fmt.Errorf("unable to find project root")
}

func IsProcessExist(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	if err := process.Signal(syscall.Signal(0)); err != nil {
		return false
	}

	return true
}

func KillProcess(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}
