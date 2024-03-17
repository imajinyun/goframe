package contract

import (
	"context"
	"time"

	"github.com/imajinyun/goframe"
)

const CacheKey = "gogin:cache"

type RememberFunc func(ctx context.Context, container goframe.IContainer) (any, error)

type ICache interface {
	Set(ctx context.Context, key string, val string, timeout time.Duration) error
	SetObject(ctx context.Context, key string, val any, timeout time.Duration) error
	SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error
	SetForever(ctx context.Context, key string, val string) error
	SetForeverObject(ctx context.Context, key string, val any) error

	Get(ctx context.Context, key string) (string, error)
	GetObject(ctx context.Context, key string, model any) error
	GetMany(ctx context.Context, keys []string) (map[string]string, error)

	SetTTL(ctx context.Context, key string, timeout time.Duration) error
	GetTTL(ctx context.Context, key string) error

	Remember(ctx context.Context, key string, timeout time.Duration, fn RememberFunc, model any) error

	Calc(ctx context.Context, key string, step int64) (int64, error)
	Increment(ctx context.Context, key string) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)

	Del(ctx context.Context, key string) error
	DelMany(ctx context.Context, keys []string) error
}
