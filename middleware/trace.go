package middleware

import (
	"github.com/imajinyun/goframe/contract"
	"github.com/imajinyun/goframe/gin"
)

type Trace struct{}

func NewTrace() *Trace {
	return &Trace{}
}

func (t *Trace) Func() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		trace := ctx.MustMake(contract.TraceKey).(contract.ITrace)
		tctx := trace.ExtractHttp(ctx.Request)
		trace.WithTrace(ctx, tctx)

		ctx.Next()
	}
}
