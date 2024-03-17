package kernel

import (
	"google.golang.org/grpc"

	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/contract"
	"github.com/imajinyun/goframe/gin"
)

type KernelProvider struct {
	HttpEngine *gin.Engine
	GrpcEngine *grpc.Server
}

func (p *KernelProvider) Name() string {
	return contract.KernelKey
}

func (p *KernelProvider) Boot(c goframe.IContainer) error {
	if p.HttpEngine == nil {
		p.HttpEngine = gin.Default()
	}

	if p.GrpcEngine == nil {
		p.GrpcEngine = grpc.NewServer()
	}

	p.HttpEngine.SetContainer(c)

	return nil
}

func (p *KernelProvider) Params(c goframe.IContainer) []any {
	return []any{p.HttpEngine, p.GrpcEngine}
}

func (p *KernelProvider) Inject(c goframe.IContainer) goframe.Handler {
	return NewKernelService
}

func (p *KernelProvider) IsDefer() bool {
	return false
}
