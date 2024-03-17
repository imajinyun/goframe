package goframe

type Handler func(...any) (any, error)

type IProvider interface {
	Name() string
	Boot(IContainer) error
	Params(IContainer) []any
	Inject(IContainer) Handler
	IsDefer() bool
}
