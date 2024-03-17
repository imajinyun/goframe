package distributed

import (
	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/contract"
)

type DistributedProvider struct{}

func (p *DistributedProvider) Name() string {
	return contract.DistributedKey
}

func (p *DistributedProvider) Boot(container goframe.IContainer) error {
	return nil
}

func (p *DistributedProvider) Params(container goframe.IContainer) []any {
	return []any{container}
}

func (p *DistributedProvider) Inject(container goframe.IContainer) goframe.Handler {
	return NewDistributedService
}

func (p *DistributedProvider) IsDefer() bool {
	return false
}
