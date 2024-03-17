package uid

import (
	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/contract"
)

type UidProvider struct{}

func (p *UidProvider) Name() string {
	return contract.UidKey
}

func (p *UidProvider) Boot(container goframe.IContainer) error {
	return nil
}

func (p *UidProvider) Params(container goframe.IContainer) []any {
	return []any{}
}

func (p *UidProvider) Inject(container goframe.IContainer) goframe.Handler {
	return NewUidService
}

func (p *UidProvider) IsDefer() bool {
	return false
}
