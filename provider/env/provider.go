package env

import (
	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/contract"
)

type EnvProvider struct {
	Dir string
}

func (p *EnvProvider) Name() string {
	return contract.EnvKey
}

func (p *EnvProvider) Boot(container goframe.IContainer) error {
	app := container.MustMake(contract.AppKey).(contract.IApp)
	p.Dir = app.WorkDir()

	return nil
}

func (p *EnvProvider) Params(container goframe.IContainer) []any {
	return []any{p.Dir}
}

func (p *EnvProvider) Inject(container goframe.IContainer) goframe.Handler {
	return NewEnvService
}

func (p *EnvProvider) IsDefer() bool {
	return false
}
