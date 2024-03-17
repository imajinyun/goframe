package orm

import (
	"context"
	"sync"
	"time"

	"github.com/imajinyun/goframe"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/imajinyun/goframe/contract"
)

type OrmService struct {
	container goframe.IContainer
	gdbs      map[string]*gorm.DB
	lock      *sync.RWMutex
}

func NewOrmService(params ...any) (any, error) {
	return &OrmService{
		container: params[0].(goframe.IContainer),
		gdbs:      make(map[string]*gorm.DB),
		lock:      &sync.RWMutex{},
	}, nil
}

func (s *OrmService) Get(option ...contract.DbOption) (*gorm.DB, error) {
	logsvc := s.container.MustMake(contract.LogKey).(contract.ILog)

	config := GetConfig(s.container)
	config.Config = &gorm.Config{Logger: NewOrmLogger(logsvc)}

	for _, opt := range option {
		if err := opt(s.container, config); err != nil {
			return nil, err
		}
	}

	if config.Dsn == "" {
		dsn, err := config.FormatDsn()
		if err != nil {
			return nil, err
		}
		config.Dsn = dsn
	}

	s.lock.RLock()
	if db, ok := s.gdbs[config.Dsn]; ok {
		s.lock.RUnlock()
		return db, nil
	}
	s.lock.RUnlock()

	s.lock.Lock()
	defer s.lock.Unlock()

	var gdb *gorm.DB
	var err error
	switch config.Driver {
	case "mysql":
		gdb, err = gorm.Open(mysql.Open(config.Dsn), config)
	case "postgres":
		gdb, err = gorm.Open(postgres.Open(config.Dsn), config)
	case "sqlite":
		gdb, err = gorm.Open(sqlite.Open(config.Dsn), config)
	case "sqlserver":
		gdb, err = gorm.Open(sqlserver.Open(config.Dsn), config)
	case "clickhouse":
		gdb, err = gorm.Open(clickhouse.Open(config.Dsn), config)
	}

	sdb, err := gdb.DB()
	if err != nil {
		return gdb, err
	}

	if config.ConnMaxIdle > 0 {
		sdb.SetMaxIdleConns(config.ConnMaxIdle)
	}

	if config.ConnMaxOpen > 0 {
		sdb.SetMaxOpenConns(config.ConnMaxOpen)
	}

	if config.ConnMaxLifetime != "" {
		lifetime, err := time.ParseDuration(config.ConnMaxLifetime)
		if err != nil {
			logsvc.Error(context.Background(), "conn max lifetime error", map[string]any{"err": err})
		} else {
			sdb.SetConnMaxLifetime(lifetime)
		}
	}

	if config.ConnMaxIdletime != "" {
		idletime, err := time.ParseDuration(config.ConnMaxIdletime)
		if err != nil {
			logsvc.Error(context.Background(), "conn max idletime error", map[string]any{"err": err})
		} else {
			sdb.SetConnMaxIdleTime(idletime)
		}
	}

	if err != nil {
		s.gdbs[config.Dsn] = gdb
	}

	return gdb, err
}
