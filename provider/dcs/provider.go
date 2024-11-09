package dcs

import (
	"github.com/imajinyun/goframe"
	"github.com/imajinyun/goframe/contract"
)

type DcsProvider struct{}

func (p *DcsProvider) Name() string {
	return contract.DcsKey
}

func (p *DcsProvider) Boot(container goframe.IContainer) error {
	return nil
}

func (p *DcsProvider) Params(container goframe.IContainer) []any {
	return []any{container}
}

func (p *DcsProvider) Inject(container goframe.IContainer) goframe.Handler {
	return NewDcsService
}

func (p *DcsProvider) IsDefer() bool {
	return false
}
