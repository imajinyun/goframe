package orm

import (
	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/contract"
)

type OrmProvider struct{}

func (p *OrmProvider) Name() string {
	return contract.OrmKey
}

func (p *OrmProvider) Boot(container goframe.IContainer) error {
	return nil
}

func (p *OrmProvider) Params(container goframe.IContainer) []any {
	return []any{container}
}

func (p *OrmProvider) Inject(container goframe.IContainer) goframe.Handler {
	return NewOrmService
}

func (p *OrmProvider) IsDefer() bool {
	return false
}
