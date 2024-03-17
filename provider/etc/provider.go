package etc

import (
	"path/filepath"

	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/contract"
)

type EtcProvider struct{}

func (p *EtcProvider) Name() string {
	return contract.EtcKey
}

func (p *EtcProvider) Boot(container goframe.IContainer) error {
	return nil
}

func (p *EtcProvider) Params(container goframe.IContainer) []any {
	appsvc := container.MustMake(contract.AppKey).(contract.IApp)
	envsvc := container.MustMake(contract.EnvKey).(contract.IEnv)

	env := envsvc.Env()
	etcdir := appsvc.EtcDir()
	envdir := filepath.Join(etcdir, env)

	return []any{container, envdir, envsvc.All()}
}

func (p *EtcProvider) Inject(container goframe.IContainer) goframe.Handler {
	return NewEtcService
}

func (p *EtcProvider) IsDefer() bool {
	return false
}
