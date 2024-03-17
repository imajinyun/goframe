package redis

import (
	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/contract"
)

type RedisProvider struct{}

func (p *RedisProvider) Name() string {
	return contract.RedisKey
}

func (p *RedisProvider) Boot(container goframe.IContainer) error {
	return nil
}

func (p *RedisProvider) Params(container goframe.IContainer) []any {
	return []any{container}
}

func (p *RedisProvider) Inject(container goframe.IContainer) goframe.Handler {
	return NewRedisService
}

func (p *RedisProvider) IsDefer() bool {
	return true
}
