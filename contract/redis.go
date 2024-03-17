package contract

import (
	"fmt"

	"github.com/imajinyun/goframe"
	"github.com/redis/go-redis/v9"
)

const RedisKey = "gogin:redis"

type RedisOption func(container goframe.IContainer, config *RedisConfig) error

type IRedis interface {
	GetClient(option ...RedisOption) (*redis.Client, error)
}

type RedisConfig struct {
	*redis.Options
}

func (c *RedisConfig) UniqueKey() string {
	return fmt.Sprintf("%s:%v:%s:%s", c.Addr, c.DB, c.Username, c.Network)
}
