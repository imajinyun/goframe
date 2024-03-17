package services

import (
	"context"
	"errors"
	"sync"
	"time"

	redisv9 "github.com/redis/go-redis/v9"

	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/contract"
	"github.com/imajinyun/goframe/provider/redis"
)

type RedisCache struct {
	container goframe.IContainer
	client    *redisv9.Client
	lock      sync.RWMutex
}

func NewRedisCache(params ...any) (any, error) {
	container := params[0].(goframe.IContainer)
	if !container.IsBind(contract.RedisKey) {
		if err := container.Bind(&redis.RedisProvider{}); err != nil {
			return nil, err
		}
	}

	redissvc := container.MustMake(contract.RedisKey).(contract.IRedis)
	client, err := redissvc.GetClient(redis.WithConfigPath("cache"))
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		container: container,
		client:    client,
		lock:      sync.RWMutex{},
	}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	cmd := c.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redisv9.Nil) {
		return "", ErrKeyNotFound
	}

	return cmd.Result()
}

func (c *RedisCache) GetObject(ctx context.Context, key string, obj any) error {
	cmd := c.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redisv9.Nil) {
		return ErrKeyNotFound
	}

	if err := cmd.Scan(obj); err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	pipe := c.client.Pipeline()
	vals := make(map[string]string)
	cmds := make([]*redisv9.StringCmd, 0, len(keys))

	for _, key := range keys {
		cmds = append(cmds, pipe.Get(ctx, key))
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, err
	}

	errs := make([]string, 0, len(keys))
	for _, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		vals[cmd.Args()[0].(string)] = val
	}

	return vals, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
	return c.client.Set(ctx, key, val, ttl).Err()
}

func (c *RedisCache) SetObject(ctx context.Context, key string, val any, ttl time.Duration) error {
	return c.client.Set(ctx, key, val, ttl).Err()
}

func (c *RedisCache) SetMany(ctx context.Context, data map[string]string, ttl time.Duration) error {
	pipe := c.client.Pipeline()
	cmds := make([]*redisv9.StatusCmd, 0, len(data))

	for key, val := range data {
		cmds = append(cmds, pipe.Set(ctx, key, val, ttl))
	}
	_, err := pipe.Exec(ctx)

	return err
}

func (c *RedisCache) SetForever(ctx context.Context, key string, val string) error {
	return c.client.Set(ctx, key, val, NoneDuration).Err()
}

func (c *RedisCache) SetForeverObject(ctx context.Context, key string, val any) error {
	return c.client.Set(ctx, key, val, NoneDuration).Err()
}

func (c *RedisCache) GetTTL(ctx context.Context, key string, val any) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

func (c *RedisCache) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, key, ttl).Err()
}

func (c *RedisCache) Remember(ctx context.Context, key string, ttl time.Duration, fn contract.RememberFunc, val any) error {
	err := c.GetObject(ctx, key, val)
	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return err
	}

	obj, err := fn(ctx, c.container)
	if err != nil {
		return err
	}

	if err := c.SetObject(ctx, key, obj, ttl); err != nil {
		return err
	}

	if err := c.GetObject(ctx, key, val); err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	return c.client.IncrBy(ctx, key, step).Result()
}

func (c *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.IncrBy(ctx, key, 1).Result()
}

func (c *RedisCache) Decrement(ctx context.Context, key string) (int64, error) {
	return c.client.IncrBy(ctx, key, -1).Result()
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) DelMany(ctx context.Context, keys []string) error {
	pipe := c.client.Pipeline()
	cmds := make([]*redisv9.IntCmd, 0, len(keys))
	for _, key := range keys {
		cmds = append(cmds, pipe.Del(ctx, key))
	}
	_, err := pipe.Exec(ctx)

	return err
}
