package cobra

import "github.com/imajinyun/goframe/contract"

func (c *Command) MustMakeApp() contract.IApp {
	return c.GetContainer().MustMake(contract.AppKey).(contract.IApp)
}

func (c *Command) MustMakeKernel() contract.IKernel {
	return c.GetContainer().MustMake(contract.KernelKey).(contract.IKernel)
}
