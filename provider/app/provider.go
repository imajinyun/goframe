package app

import (
	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/contract"
)

type AppProvider struct {
	Base string
}

func (p *AppProvider) Name() string {
	return contract.AppKey
}

func (p *AppProvider) Boot(container goframe.IContainer) error {
	return nil
}

func (p *AppProvider) Params(container goframe.IContainer) []any {
	return []any{container, p.Base}
}

func (p *AppProvider) Inject(container goframe.IContainer) goframe.Handler {
	return NewAppService
}

func (p *AppProvider) IsDefer() bool {
	return false
}
