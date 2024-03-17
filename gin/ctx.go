package gin

import (
	"context"
)

func (ctx *Context) Context() context.Context {
	return ctx.Request.Context()
}

func (ctx *Context) Make(key string) (any, error) {
	return ctx.container.Make(key)
}

func (ctx *Context) MustMake(key string) any {
	return ctx.container.MustMake(key)
}

func (ctx *Context) MakeNew(key string, params []any) (any, error) {
	return ctx.container.NewMake(key, params)
}

func (ctx *Context) IsBind(key string) bool {
	return ctx.container.IsBind(key)
}
