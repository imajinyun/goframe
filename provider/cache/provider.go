package cache

import (
	"strings"

	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/contract"
	"github.com/imajinyun/goframe/provider/cache/services"
)

type CacheProvider struct {
	goframe.IProvider

	Driver string
}

func (p *CacheProvider) Name() string {
	return contract.CacheKey
}

func (p *CacheProvider) Boot(container goframe.IContainer) error {
	return nil
}

func (p *CacheProvider) Params(container goframe.IContainer) []any {
	return []any{container}
}

func (p *CacheProvider) Inject(container goframe.IContainer) goframe.Handler {
	if p.Driver == "" {
		etc, err := container.Make(contract.EtcKey)
		if err != nil {
			return services.NewMemoryCache
		}

		etcsvc := etc.(contract.IEtc)
		p.Driver = strings.ToLower(etcsvc.GetString("cache.driver"))
	}

	switch p.Driver {
	case "redis":
		return services.NewRedisCache
	case "memory":
		return services.NewMemoryCache
	default:
		return services.NewMemoryCache
	}
}

func (p *CacheProvider) IsDefer() bool {
	return true
}
