package orm

import (
	"context"
	"time"

	"gorm.io/gorm/logger"

	"github.com/imajinyun/goframe/contract"
)

type OrmLogger struct {
	logger contract.ILog
}

func NewOrmLogger(logger contract.ILog) *OrmLogger {
	return &OrmLogger{logger: logger}
}

func (l *OrmLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *OrmLogger) Info(ctx context.Context, msg string, val ...any) {
	fields := map[string]any{"fields": val}
	l.logger.Info(ctx, msg, fields)
}

func (l *OrmLogger) Warn(ctx context.Context, msg string, val ...any) {
	fields := map[string]any{"fields": val}
	l.logger.Warn(ctx, msg, fields)
}

func (l *OrmLogger) Error(ctx context.Context, msg string, val ...any) {
	fields := map[string]any{"fields": val}
	l.logger.Error(ctx, msg, fields)
}

func (l *OrmLogger) Trace(ctx context.Context, begin time.Time, fn func() (string, int64), err error) {
	sql, rows := fn()
	elapsed := time.Since(begin)
	fields := map[string]any{
		"begin": begin,
		"error": err,
		"sql":   sql,
		"rows":  rows,
		"time":  elapsed,
	}

	s := "orm trace sql"
	l.logger.Trace(ctx, s, fields)
}
