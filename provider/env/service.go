package env

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"strings"

	"github.com/imajinyun/goframe/contract"
)

var errEnvServiceParams = errors.New("env service params error")

type EnvService struct {
	dir  string
	dict map[string]string
}

func NewEnvService(params ...any) (any, error) {
	if len(params) != 1 {
		return nil, errEnvServiceParams
	}

	dir := params[0].(string)
	env := &EnvService{
		dir:  dir,
		dict: map[string]string{"APP_ENV": contract.EnvDevelopment},
	}
	name := path.Join(dir, ".env")
	file, err := os.Open(name)
	if err == nil {
		defer file.Close()

		buf := bufio.NewReader(file)
		for {
			line, _, err := buf.ReadLine()
			if err == io.EOF {
				break
			}

			s := bytes.SplitN(line, []byte{'='}, 2)
			if len(s) < 2 {
				continue
			}
			key, val := string(s[0]), string(s[1])
			env.dict[key] = val
		}
	}

	for _, v := range os.Environ() {
		pair := strings.SplitN(v, "=", 2)
		if len(pair) < 2 {
			continue
		}
		env.dict[pair[0]] = pair[1]
	}

	return env, nil
}

func (s *EnvService) Env() string {
	return s.Get("APP_ENV")
}

func (s *EnvService) Get(key string) string {
	if val, ok := s.dict[key]; ok {
		return val
	}

	return ""
}

func (s *EnvService) All() map[string]string {
	return s.dict
}

func (s *EnvService) Exist(key string) bool {
	_, ok := s.dict[key]

	return ok
}
