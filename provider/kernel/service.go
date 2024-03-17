package kernel

import (
	"net/http"

	"google.golang.org/grpc"

	"github.com/imajinyun/goframe/gin"
)

type KernelService struct {
	grpcEngine *grpc.Server
	httpEngine *gin.Engine
}

func NewKernelService(params ...any) (any, error) {
	httpEngine, grpcEngine := params[0].(*gin.Engine), params[1].(*grpc.Server)

	return &KernelService{httpEngine: httpEngine, grpcEngine: grpcEngine}, nil
}

func (s *KernelService) GrpcHandler() *grpc.Server {
	return s.grpcEngine
}

func (s *KernelService) HttpHandler() http.Handler {
	return s.httpEngine
}
