package redis

import (
	"sync"

	"github.com/imajinyun/goframe"

	"github.com/redis/go-redis/v9"

	"github.com/imajinyun/goframe/contract"
)

type RedisService struct {
	container goframe.IContainer
	clients   map[string]*redis.Client
	lock      *sync.RWMutex
}

func NewRedisService(params ...any) (any, error) {
	return &RedisService{
		container: params[0].(goframe.IContainer),
		clients:   make(map[string]*redis.Client),
		lock:      &sync.RWMutex{},
	}, nil
}

func (s *RedisService) GetClient(option ...contract.RedisOption) (*redis.Client, error) {
	config := GetConfig(s.container)

	for _, opt := range option {
		if err := opt(s.container, config); err != nil {
			return nil, err
		}
	}

	key := config.UniqueKey()

	s.lock.RLock()
	if db, ok := s.clients[key]; ok {
		s.lock.RUnlock()
		return db, nil
	}
	s.lock.RUnlock()

	s.lock.RLock()
	defer s.lock.RUnlock()

	client := redis.NewClient(config.Options)
	s.clients[key] = client

	return client, nil
}
