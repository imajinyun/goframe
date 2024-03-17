package orm

import (
	"context"

	"github.com/imajinyun/goframe/contract"

	"github.com/imajinyun/goframe"
)

func GetConfig(container goframe.IContainer) *contract.DbConfig {
	etcsvc := container.MustMake(contract.EtcKey).(contract.IEtc)
	logsvc := container.MustMake(contract.LogKey).(contract.ILog)
	config := &contract.DbConfig{}
	if err := etcsvc.Load("database.default", config); err != nil {
		logsvc.Error(context.Background(), "parse database config error", nil)
		return nil
	}

	return config
}

func WithConfigPath(path string) contract.DbOption {
	return func(container goframe.IContainer, config *contract.DbConfig) error {
		etcsvc := container.MustMake(contract.EtcKey).(contract.IEtc)
		if err := etcsvc.Load(path, config); err != nil {
			return err
		}
		return nil
	}
}

func WithGormConfig(fn func(options *contract.DbConfig)) contract.DbOption {
	return func(container goframe.IContainer, config *contract.DbConfig) error {
		fn(config)
		return nil
	}
}

func WithDryRun() contract.DbOption {
	return func(container goframe.IContainer, config *contract.DbConfig) error {
		config.DryRun = true
		return nil
	}
}

func WithFullSaveAssociations() contract.DbOption {
	return func(container goframe.IContainer, config *contract.DbConfig) error {
		config.FullSaveAssociations = true
		return nil
	}
}
