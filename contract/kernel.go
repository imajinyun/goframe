package contract

import (
	"net/http"

	"google.golang.org/grpc"
)

const KernelKey = "gogin:kernel"

type IKernel interface {
	GrpcHandler() *grpc.Server
	HttpHandler() http.Handler
}
