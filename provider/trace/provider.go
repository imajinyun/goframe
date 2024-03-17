package trace

import (
	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/contract"
)

type TraceProvider struct {
	container goframe.IContainer
}

func (p *TraceProvider) Name() string {
	return contract.TraceKey
}

func (p *TraceProvider) Boot(container goframe.IContainer) error {
	p.container = container

	return nil
}

func (p *TraceProvider) Params(container goframe.IContainer) []any {
	return []any{p.container}
}

func (p *TraceProvider) Inject(container goframe.IContainer) goframe.Handler {
	return NewTraceService
}

func (p *TraceProvider) IsDefer() bool {
	return false
}
