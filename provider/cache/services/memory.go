package services

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/imajinyun/goframe"

	"github.com/pkg/errors"
)

type MemoryCache struct {
	container goframe.IContainer
	data      map[string]*MemoryData
	lock      sync.RWMutex
}

type MemoryData struct {
	val   any
	ttl   time.Duration
	ctime time.Time
}

func NewMemoryCache(params ...any) (any, error) {
	return &MemoryCache{
		container: params[0].(goframe.IContainer),
		data:      make(map[string]*MemoryData),
		lock:      sync.RWMutex{},
	}, nil
}

func (c *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	var val string
	if err := c.GetObject(ctx, key, &val); err != nil {
		return "", err
	}

	return val, nil
}

func (c *MemoryCache) GetObject(ctx context.Context, key string, obj any) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if v, ok := c.data[key]; ok {
		if v.ttl != NoneDuration {
			if time.Now().Sub(v.ctime) > v.ttl {
				delete(c.data, key)
				return ErrKeyNotFound
			}
		}

		byt, _ := json.Marshal(v.val)
		if err := json.Unmarshal(byt, obj); err != nil {
			return err
		}

		return nil
	}

	return ErrKeyNotFound
}

func (c *MemoryCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	errs := make([]string, 0, len(keys))
	rets := make(map[string]string)

	for _, key := range keys {
		val, err := c.Get(ctx, key)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		rets[key] = val
	}

	if len(errs) == 0 {
		return rets, nil
	}

	return rets, errors.New(strings.Join(errs, "|"))
}

func (c *MemoryCache) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
	return c.SetObject(ctx, key, val, ttl)
}

func (c *MemoryCache) SetObject(ctx context.Context, key string, val any, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[key] = &MemoryData{
		val:   val,
		ttl:   ttl,
		ctime: time.Now(),
	}

	return nil
}

func (c *MemoryCache) SetMany(ctx context.Context, data map[string]string, ttl time.Duration) error {
	errs := []string{}
	for k, v := range data {
		if err := c.Set(ctx, k, v, ttl); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "|"))
	}

	return nil
}

func (c *MemoryCache) SetForever(ctx context.Context, key string, val string) error {
	return c.Set(ctx, key, val, NoneDuration)
}

func (c *MemoryCache) SetForeverObject(ctx context.Context, key string, val any) error {
	return c.SetObject(ctx, key, val, NoneDuration)
}

func (c *MemoryCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if v, ok := c.data[key]; ok {
		return v.ttl, nil
	}

	return NoneDuration, ErrKeyNotFound
}

func (c *MemoryCache) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if v, ok := c.data[key]; ok {
		v.ttl = ttl
		return nil
	}

	return ErrKeyNotFound
}

func (c *MemoryCache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	var val int64
	err := c.GetObject(ctx, key, &val)
	val = val + step
	if err == nil {
		c.data[key].val = val
		return val, nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return 0, err
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[key] = &MemoryData{
		val:   val,
		ttl:   NoneDuration,
		ctime: time.Now(),
	}

	return val, nil
}

func (c *MemoryCache) Increment(ctx context.Context, key string) (int64, error) {
	return c.Calc(ctx, key, 1)
}

func (c *MemoryCache) Decrement(ctx context.Context, key string) (int64, error) {
	return c.Calc(ctx, key, -1)
}

func (c *MemoryCache) Del(ctx context.Context, key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, key)

	return nil
}

func (c *MemoryCache) DelMany(ctx context.Context, keys []string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, key := range keys {
		delete(c.data, key)
	}

	return nil
}
