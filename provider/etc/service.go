package etc

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/imajinyun/goframe"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"

	"github.com/imajinyun/goframe/contract"
)

var errEtcServiceParams = errors.New("etc service params error")

type EtcService struct {
	container goframe.IContainer
	dir       string
	sep       string
	envs      map[string]string
	etcs      map[string]any
	raws      map[string][]byte
	lock      sync.RWMutex
}

func NewEtcService(params ...any) (any, error) {
	if len(params) != 3 {
		return nil, errEtcServiceParams
	}
	container := params[0].(goframe.IContainer)
	envdir := params[1].(string)
	envs := params[2].(map[string]string)

	cs := &EtcService{
		container: container,
		dir:       envdir,
		sep:       ".",
		envs:      envs,
		etcs:      map[string]any{},
		raws:      map[string][]byte{},
		lock:      sync.RWMutex{},
	}

	if _, err := os.Stat(envdir); os.IsNotExist(err) {
		return cs, nil
	}

	files, err := os.ReadDir(envdir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, file := range files {
		name := file.Name()
		if err := cs.load(envdir, name); err != nil {
			log.Println(err)
			continue
		}
	}

	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := watch.Add(envdir); err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		for {
			select {
			case evt := <-watch.Events:
				{
					path, _ := filepath.Abs(evt.Name)
					idx := strings.LastIndex(path, string(os.PathSeparator))
					dir, name := path[:idx], path[idx+1:]

					if evt.Op&fsnotify.Create == fsnotify.Create {
						log.Println("create:", evt.Name)
						cs.load(dir, name)
					}

					if evt.Op&fsnotify.Write == fsnotify.Write {
						log.Println("write:", evt.Name)
						cs.load(dir, name)
					}

					if evt.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("remove:", evt.Name)
						cs.remove(dir, name)
					}
				}
			case err := <-watch.Errors:
				{
					log.Println(err)
					return
				}
			}
		}
	}()

	return cs, nil
}

func (s *EtcService) Get(key string) any {
	return s.find(key)
}

func (s *EtcService) GetBool(key string) bool {
	return cast.ToBool(s.find(key))
}

func (s *EtcService) GetInt(key string) int {
	return cast.ToInt(s.find(key))
}

func (s *EtcService) GetInt64(key string) int64 {
	return cast.ToInt64(s.find(key))
}

func (s *EtcService) GetFloat32(key string) float32 {
	return cast.ToFloat32(s.find(key))
}

func (s *EtcService) GetFloat64(key string) float64 {
	return cast.ToFloat64(s.find(key))
}

func (s *EtcService) GetString(key string) string {
	return cast.ToString(s.find(key))
}

func (s *EtcService) GetTime(key string) time.Time {
	return cast.ToTime(s.find(key))
}

func (s *EtcService) GetIntSlice(key string) []int {
	return cast.ToIntSlice(s.find(key))
}

func (s *EtcService) GetStringSlice(key string) []string {
	return cast.ToStringSlice(s.find(key))
}

func (s *EtcService) GetStringMap(key string) map[string]any {
	return cast.ToStringMap(s.find(key))
}

func (s *EtcService) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(s.find(key))
}

func (s *EtcService) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(s.find(key))
}

func (s *EtcService) Load(key string, val any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  val,
	})
	if err != nil {
		return err
	}

	return decoder.Decode(s.find(key))
}

func (s *EtcService) Exist(key string) bool {
	return s.find(key) != nil
}

func (s *EtcService) find(key string) any {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return search(s.etcs, strings.Split(key, s.sep))
}

func (s *EtcService) load(dir string, file string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	t := strings.Split(file, ".")
	if len(t) == 2 || (t[1] == "yaml" || t[1] == "yml") {
		name := t[0]

		buf, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil {
			return err
		}

		buf = replace(buf, s.envs)
		dic := map[string]any{}
		if err := yaml.Unmarshal(buf, &dic); err != nil {
			return err
		}
		s.etcs[name] = dic
		s.raws[name] = buf

		if name == "app" && s.container.IsBind(contract.AppKey) {
			if p, ok := dic["path"]; ok {
				appsvc := s.container.MustMake(contract.AppKey).(contract.IApp)
				appsvc.Load(cast.ToStringMapString(p))
			}
		}
	}

	return nil
}

func (s *EtcService) remove(dir string, file string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	t := strings.Split(file, ".")
	if len(t) == 2 && (t[1] == "yaml" || t[1] == "yml") {
		name := t[0]
		delete(s.raws, name)
		delete(s.raws, name)
	}

	return nil
}

func replace(content []byte, dict map[string]string) []byte {
	if dict == nil {
		return content
	}

	for k, val := range dict {
		key := fmt.Sprintf("env(%s)", k)
		content = bytes.ReplaceAll(content, []byte(key), []byte(val))
	}

	return content
}

func search(source map[string]any, path []string) any {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		if len(path) == 1 {
			return next
		}

		switch next.(type) {
		case map[any]any:
			return search(cast.ToStringMap(next), path[1:])
		case map[string]any:
			return search(next.(map[string]any), path[1:])
		default:
			return nil
		}
	}

	return nil
}
