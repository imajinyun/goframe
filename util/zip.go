package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(src string, dst string) ([]string, error) {
	var names []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return names, err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dst, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return names, fmt.Errorf("%s: illegal file path", fpath)
		}

		names = append(names, fpath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return names, err
		}

		out, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return names, err
		}

		p, err := f.Open()
		if err != nil {
			return names, err
		}

		_, err = io.Copy(out, p)
		out.Close()

		if err != nil {
			return names, err
		}
	}

	return names, nil
}
