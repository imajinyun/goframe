package util

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
)

func Exist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}

	return true
}

func IsHideDir(path string) bool {
	return len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".")
}

func SubDir(dir string) ([]string, error) {
	subs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, sub := range subs {
		if sub.IsDir() {
			res = append(res, sub.Name())
		}
	}

	return res, nil
}

func CopyDir(src string, dst string) error {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		name := strings.Replace(path, src, "", 1)
		if name == "" {
			return nil
		}

		if info.IsDir() {
			return os.Mkdir(filepath.Join(dst, name), 0o755)
		} else {
			data, err := os.ReadFile(filepath.Join(src, name))
			if err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(dst, name), data, os.ModePerm)
		}
	})

	return err
}

func CopyFile(src string, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, os.ModePerm)
}

func Download(path string, url string) error {
	log.Printf("path: %v, url: %v\n", path, url)
	rsp, err := resty.New().R().Get(url)
	if err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, rsp.RawBody())

	return err
}

func DownloadFile(path string, url string) error {
	log.Printf("path: %v, url: %v\n", path, url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
