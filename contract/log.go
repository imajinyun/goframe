package contract

import (
	"context"
	"io"
	"time"
)

type LogLevel uint32

const LogKey = "gogin:log"

const (
	UnknownLevel LogLevel = iota
	PanicLevel
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type Handler func(ctx context.Context) map[string]any

type Formatter func(level LogLevel, t time.Time, msg string, data map[string]any) ([]byte, error)

type ILog interface {
	Panic(ctx context.Context, msg string, data map[string]any)
	Fatal(ctx context.Context, msg string, data map[string]any)
	Error(ctx context.Context, msg string, data map[string]any)
	Warn(ctx context.Context, msg string, data map[string]any)
	Info(ctx context.Context, msg string, data map[string]any)
	Debug(ctx context.Context, msg string, data map[string]any)
	Trace(ctx context.Context, msg string, data map[string]any)

	SetLevel(level LogLevel)
	SetWriter(writer io.Writer)
	SetHandler(handler Handler)
	SetFormatter(formatter Formatter)
}
