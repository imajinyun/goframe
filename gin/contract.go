package gin

import "github.com/imajinyun/goframe/contract"

func (ctx *Context) MustMakeApp() contract.IApp {
	return ctx.MustMake(contract.AppKey).(contract.IApp)
}
