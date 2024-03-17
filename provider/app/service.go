package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/imajinyun/goframe"

	"github.com/google/uuid"
)

var errAppServiceParams = errors.New("app service params error")

type AppService struct {
	container goframe.IContainer
	workspace string
	id        string
	cfgs      map[string]string
	envs      map[string]string
	args      map[string]string
}

func NewAppService(params ...any) (any, error) {
	if len(params) != 2 {
		return nil, errAppServiceParams
	}

	container := params[0].(goframe.IContainer)
	workspace := params[1].(string)

	id := uuid.New().String()
	cfgs := make(map[string]string)
	appsvc := &AppService{container: container, workspace: workspace, id: id, cfgs: cfgs}
	_ = appsvc.loadEnvs()
	_ = appsvc.loadArgs()

	return appsvc, nil
}

func (s *AppService) AppID() string {
	return s.id
}

func (s *AppService) Version() string {
	return GoginVersion
}

func (s *AppService) WorkDir() string {
	if s.workspace != "" {
		return s.workspace
	}
	workspace := s.getValueByKey("base_dir", "BASE_DIR", "app.path.base_dir")
	if workspace != "" {
		return workspace
	}

	return ""
}

func (s *AppService) AppDir() string {
	val := s.getValueByKey("app_dir", "APP_DIR", "app.path.app_dir")
	if val != "" {
		return val
	}

	return filepath.Join(s.WorkDir(), "app")
}

func (s *AppService) EtcDir() string {
	val := s.getValueByKey("etc_dir", "ETC_DIR", "app.path.etc_dir")
	if val != "" {
		return val
	}

	return filepath.Join(s.WorkDir(), "etc")
}

func (s *AppService) LogDir() string {
	val := s.getValueByKey("log_dir", "LOG_DIR", "app.path.log_dir")
	if val != "" {
		return val
	}

	return filepath.Join(s.WorkDir(), "log")
}

func (s *AppService) RunDir() string {
	val := s.getValueByKey("run_dir", "RUN_DIR", "app.path.run_dir")
	if val != "" {
		return val
	}

	return filepath.Join(s.StorageDir(), "run")
}

func (s *AppService) HttpDir() string {
	val := s.getValueByKey("console_dir", "CONSOLE_DIR", "app.path.console_dir")
	if val != "" {
		return val
	}

	return filepath.Join(s.WorkDir(), "app", "http")
}

func (s *AppService) TestDir() string {
	val := s.getValueByKey("test_dir", "TEST_DIR", "app.path.test_dir")
	if val != "" {
		return val
	}

	return filepath.Join(s.WorkDir(), "test")
}

func (s *AppService) ConsoleDir() string {
	val := s.getValueByKey("console_dir", "CONSOLE_DIR", "app.path.console_dir")
	if val != "" {
		return val
	}

	return filepath.Join(s.WorkDir(), "app")
}

func (s *AppService) CommandDir() string {
	val := s.getValueByKey("command_dir", "COMMAND_DIR", "app.path.command_dir")
	if val != "" {
		return val
	}

	return filepath.Join(s.ConsoleDir(), "command")
}

func (s *AppService) ProviderDir() string {
	val := s.getValueByKey("provider_dir", "PROVIDER_DIR", "app.path.provider_dir")
	if val != "" {
		return val
	}

	return filepath.Join(s.WorkDir(), "app", "provider")
}

func (s *AppService) MiddlewareDir() string {
	val := s.getValueByKey("middleware_dir", "MIDDLEWARE_DIR", "app.path.middleware_dir")
	if val != "" {
		return val
	}

	return filepath.Join(s.HttpDir(), "middleware")
}

func (s *AppService) StorageDir() string {
	val := s.getValueByKey("storage_dir", "STORAGE_DIR", "app.path.storage_dir")
	if val != "" {
		return val
	}

	return filepath.Join(s.WorkDir(), "storage")
}

func (s *AppService) DeployDir() string {
	val := s.getValueByKey("deploy_dir", "DEPLOY_DIR", "app.path.deploy_dir")
	if val != "" {
		return val
	}

	return filepath.Join(s.WorkDir(), "deploy")
}

func (s *AppService) Load(kv map[string]string) {
	for k, v := range kv {
		s.cfgs[k] = v
	}
}

func (s *AppService) loadEnvs() error {
	if s.envs == nil {
		s.envs = map[string]string{}
	}

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		s.envs[pair[0]] = pair[1]
	}

	return nil
}

func (s *AppService) loadArgs() error {
	if s.args == nil {
		s.args = map[string]string{}
	}

	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "--") {
			pair := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)
			s.args[pair[0]] = pair[1]
		}
	}

	return nil
}

func (s *AppService) getValueByKey(argKey string, envKey string, cfgKey string) string {
	if s.args != nil {
		if v, ok := s.args[argKey]; ok {
			return v
		}
	}

	if s.envs != nil {
		if v, ok := s.envs[envKey]; ok {
			return v
		}
	}

	if s.cfgs != nil {
		if v, ok := s.cfgs[cfgKey]; ok {
			return v
		}
	}

	return ""
}
