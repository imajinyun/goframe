package gin

import "github.com/imajinyun/goframe"

func (e *Engine) SetContainer(container goframe.IContainer) {
	e.container = container
}

func (e *Engine) GetContainer() goframe.IContainer {
	return e.container
}

func (e *Engine) Bind(provider goframe.IProvider) error {
	return e.container.Bind(provider)
}

func (e *Engine) IsBind(key string) bool {
	return e.container.IsBind(key)
}
