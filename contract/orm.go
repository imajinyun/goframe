package contract

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/imajinyun/goframe"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const OrmKey = "gogin:orm"

type IOrm interface {
	Get(option ...DbOption) (*gorm.DB, error)
	GetTables(ctx context.Context, db *gorm.DB) ([]string, error)
	HasTable(ctx context.Context, db *gorm.DB, table string) (bool, error)
	GetTableColumns(ctx context.Context, db *gorm.DB, table string) ([]TableColumn, error)
}

type TableColumn struct {
	Key     string `gorm:"column:Key"`
	Type    string `gorm:"column:Type"`
	Null    string `gorm:"column:Null"`
	Field   string `gorm:"column:Field"`
	Extra   string `gorm:"column:Extra"`
	Default string `gorm:"column:Default"`
}

type DbOption func(container goframe.IContainer, config *DbConfig) error

type DbConfig struct {
	WriterTimeout string `yaml:"writer_timeout"`
	ReaderTimeout string `yaml:"reader_timeout"`
	ParseTime     bool   `yaml:"parse_time"`
	Loc           string `yaml:"loc"`
	Dsn           string `yaml:"dsn"`

	Host      string `ymal:"host"`
	Port      int    `yaml:"port"`
	Driver    string `yaml:"driver"`
	Charset   string `yaml:"charset"`
	Timeout   string `yaml:"timeout"`
	Protocol  string `yaml:"protocol"`
	Database  string `yaml:"database"`
	Collation string `yaml:"collation"`

	AllowNativePasswords bool `yaml:"allow_native_passwords"`

	Username string `yaml:"username"`
	Password string `yaml:"password"`

	ConnMaxIdle     int    `ymal:"conn_max_idle"`
	ConnMaxOpen     int    `ymal:"conn_max_open"`
	ConnMaxIdletime string `yaml:"conn_max_idletime"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime"`

	*gorm.Config
}

func (c *DbConfig) FormatDsn() (string, error) {
	port := strconv.Itoa(c.Port)
	timeout, err := time.ParseDuration(c.Timeout)
	if err != nil {
		return "", err
	}

	rto, err := time.ParseDuration(c.ReaderTimeout)
	if err != nil {
		return "", nil
	}

	wto, err := time.ParseDuration(c.WriterTimeout)
	if err != nil {
		return "", nil
	}

	loc, err := time.LoadLocation(c.Loc)
	if err != nil {
		return "", nil
	}

	driver := &mysql.Config{
		User:                 c.Username,
		Passwd:               c.Password,
		DBName:               c.Database,
		Net:                  c.Protocol,
		Loc:                  loc,
		Addr:                 net.JoinHostPort(c.Host, port),
		Timeout:              timeout,
		Collation:            c.Collation,
		ParseTime:            c.ParseTime,
		ReadTimeout:          rto,
		WriteTimeout:         wto,
		AllowNativePasswords: c.AllowNativePasswords,
	}

	return driver.FormatDSN(), nil
}
