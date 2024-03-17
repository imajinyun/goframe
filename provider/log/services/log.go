package services

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/contract"
	"github.com/imajinyun/goframe/provider/log/formatter"
)

type Log struct {
	level     contract.LogLevel
	writer    io.Writer
	handler   contract.Handler
	formatter contract.Formatter
	container goframe.IContainer
}

func (l *Log) IsLevelEnable(level contract.LogLevel) bool {
	return level <= l.level
}

func (l *Log) SetLevel(level contract.LogLevel) {
	l.level = level
}

func (l *Log) SetWriter(writer io.Writer) {
	l.writer = writer
}

func (l *Log) SetHandler(handler contract.Handler) {
	l.handler = handler
}

func (l *Log) SetFormatter(formatter contract.Formatter) {
	l.formatter = formatter
}

func (l *Log) Panic(ctx context.Context, msg string, data map[string]any) {
	l.logf(ctx, contract.PanicLevel, msg, data)
}

func (l *Log) Fatal(ctx context.Context, msg string, data map[string]any) {
	l.logf(ctx, contract.FatalLevel, msg, data)
}

func (l *Log) Error(ctx context.Context, msg string, data map[string]any) {
	l.logf(ctx, contract.ErrorLevel, msg, data)
}

func (l *Log) Warn(ctx context.Context, msg string, data map[string]any) {
	l.logf(ctx, contract.WarnLevel, msg, data)
}

func (l *Log) Info(ctx context.Context, msg string, data map[string]any) {
	l.logf(ctx, contract.InfoLevel, msg, data)
}

func (l *Log) Debug(ctx context.Context, msg string, data map[string]any) {
	l.logf(ctx, contract.DebugLevel, msg, data)
}

func (l *Log) Trace(ctx context.Context, msg string, data map[string]any) {
	l.logf(ctx, contract.TraceLevel, msg, data)
}

func (l *Log) logf(ctx context.Context, level contract.LogLevel, msg string, data map[string]any) error {
	if !l.IsLevelEnable(level) {
		return nil
	}

	if data == nil {
		data = make(map[string]any)
	}

	item := data
	if l.handler != nil {
		if t := l.handler(ctx); t != nil {
			for k, v := range t {
				item[k] = v
			}
		}
	}

	if l.container != nil && l.container.IsBind(contract.TraceKey) {
		tacsvc := l.container.MustMake(contract.TraceKey).(contract.ITrace)
		if trace := tacsvc.GetTrace(ctx); trace != nil {
			dict := tacsvc.ToMap(trace)
			for k, v := range dict {
				item[k] = v
			}
		}
	}

	if l.formatter == nil {
		l.formatter = formatter.TextFormatter
	}

	byt, err := l.formatter(level, time.Now(), msg, item)
	if err != nil {
		return err
	}

	if level == contract.PanicLevel {
		log.Panicln(string(byt))
		return nil
	}
	l.writer.Write(byt)
	l.writer.Write([]byte("\r\n"))

	return nil
}
