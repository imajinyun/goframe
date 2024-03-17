package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/imajinyun/goframe"

	"github.com/redis/go-redis/v9"

	"github.com/imajinyun/goframe/contract"
)

func GetConfig(container goframe.IContainer) *contract.RedisConfig {
	logsvc := container.MustMake(contract.LogKey).(contract.ILog)
	config := &contract.RedisConfig{Options: &redis.Options{}}
	opt := WithConfigPath("redis")
	err := opt(container, config)
	if err != nil {
		logsvc.Error(context.Background(), "parse redis config error", nil)
		return nil
	}

	return config
}

func WithConfigPath(path string) contract.RedisOption {
	return func(container goframe.IContainer, rdc *contract.RedisConfig) error {
		etcsvc := container.MustMake(contract.EtcKey).(contract.IEtc)
		config := etcsvc.GetStringMapString(path)

		if host, ok := config["host"]; ok {
			if port, ok := config["port"]; ok {
				rdc.Addr = host + ":" + port
			}
		}

		if db, ok := config["db"]; ok {
			tmp, err := strconv.Atoi(db)
			if err != nil {
				return err
			}
			rdc.DB = tmp
		}

		if username, ok := config["username"]; ok {
			rdc.Username = username
		}

		if password, ok := config["password"]; ok {
			rdc.Password = password
		}

		if timeout, ok := config["timeout"]; ok {
			tmp, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			rdc.DialTimeout = tmp
		}

		if min, ok := config["conn_min_idle"]; ok {
			tmp, err := strconv.Atoi(min)
			if err != nil {
				return err
			}
			rdc.MinIdleConns = tmp
		}

		if max, ok := config["conn_max_idle"]; ok {
			tmp, err := strconv.Atoi(max)
			if err != nil {
				return err
			}
			rdc.PoolSize = tmp
		}

		if timeout, ok := config["conn_max_lifetime"]; ok {
			tmp, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			rdc.ConnMaxLifetime = tmp
		}

		if idletime, ok := config["conn_max_idletime"]; ok {
			t, err := time.ParseDuration(idletime)
			if err != nil {
				return err
			}
			rdc.ConnMaxIdleTime = t
		}

		return nil
	}
}

func WithRedisConfig(fn func(options *contract.RedisConfig)) contract.RedisOption {
	return func(container goframe.IContainer, rdc *contract.RedisConfig) error {
		fn(rdc)
		return nil
	}
}
