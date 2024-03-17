package middleware

import (
	"context"
	"log"
	"time"

	"github.com/imajinyun/goframe/gin"
)

var defaultTimeout = Timeout{d: 10 * time.Second}

type Timeout struct {
	d time.Duration
}

type TimeoutOption func(t *Timeout)

func NewTimeout(opts ...TimeoutOption) *Timeout {
	dt := defaultTimeout
	for _, opt := range opts {
		opt(&dt)
	}

	return &dt
}

func (t *Timeout) Func() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		finished := make(chan struct{}, 1)
		panicked := make(chan interface{}, 1)

		timeoutCtx, cancelFunc := context.WithTimeout(ctx.Context(), t.d)
		defer cancelFunc()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicked <- p
				}
			}()

			ctx.Next()

			finished <- struct{}{}
		}()

		select {
		case p := <-panicked:
			ctx.ToSetStatus(500).ToJson("timeout")
			log.Println(p)
		case <-finished:
			log.Println("finished")
		case <-timeoutCtx.Done():
			ctx.ToSetStatus(500).ToJson("timeout")
		}
	}
}

func WithTimeout(d time.Duration) TimeoutOption {
	return func(t *Timeout) {
		t.d = d
	}
}
