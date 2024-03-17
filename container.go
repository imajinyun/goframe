package goframe

import (
	"errors"
	"fmt"
	"sync"
)

type IContainer interface {
	Bind(provider IProvider) error
	IsBind(key string) bool

	Make(key string) (any, error)
	MustMake(key string) any
	NewMake(key string, params []any) (any, error)
}

type Container struct {
	IContainer

	providers map[string]IProvider
	instances map[string]any

	lock sync.RWMutex
}

func NewContainer() *Container {
	return &Container{
		providers: map[string]IProvider{},
		instances: map[string]any{},
		lock:      sync.RWMutex{},
	}
}

func (c *Container) Bind(provider IProvider) error {
	c.lock.Lock()
	key := provider.Name()
	c.providers[key] = provider
	c.lock.Unlock()

	if provider.IsDefer() == false {
		if err := provider.Boot(c); err != nil {
			return err
		}

		params, handler := provider.Params(c), provider.Inject(c)
		instance, err := handler(params...)
		if err != nil {
			return errors.New(err.Error())
		}
		c.instances[key] = instance
	}

	return nil
}

func (c *Container) IsBind(key string) bool {
	return c.findProvider(key) != nil
}

func (c *Container) Make(key string) (any, error) {
	return c.make(key, nil, false)
}

func (c *Container) MustMake(key string) any {
	srv, err := c.make(key, nil, false)
	if err != nil {
		panic(err)
	}

	return srv
}

func (c *Container) NewMake(key string, params []any) (any, error) {
	return c.make(key, params, true)
}

func (c *Container) Providers() []string {
	list := []string{}
	for _, provider := range c.providers {
		name := provider.Name()
		list = append(list, fmt.Sprint(name))
	}

	return list
}

func (c *Container) make(key string, params []any, forced bool) (any, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	p := c.findProvider(key)
	if p == nil {
		return nil, errors.New(fmt.Sprintf("contract %s not bound", key))
	}

	if forced {
		return c.newInstance(p, params)
	}

	if ins, ok := c.instances[key]; ok {
		return ins, nil
	}

	ins, err := c.newInstance(p, nil)
	if err != nil {
		return nil, err
	}

	c.instances[key] = ins

	return ins, nil
}

func (c *Container) findProvider(key string) IProvider {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if p, ok := c.providers[key]; ok {
		return p
	}

	return nil
}

func (c *Container) newInstance(p IProvider, params []any) (any, error) {
	if err := p.Boot(c); err != nil {
		return nil, err
	}

	if params == nil {
		params = p.Params(c)
	}

	handler := p.Inject(c)
	ins, err := handler(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return ins, nil
}
