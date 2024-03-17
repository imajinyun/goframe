package log

import (
	"io"
	"strings"

	"github.com/imajinyun/goframe"

	"github.com/imajinyun/goframe/contract"
	"github.com/imajinyun/goframe/provider/log/formatter"
	"github.com/imajinyun/goframe/provider/log/services"
)

type LogProvider struct {
	goframe.IProvider

	Driver    string
	Level     contract.LogLevel
	Writer    io.Writer
	Handler   contract.Handler
	Formatter contract.Formatter
}

func (p *LogProvider) Name() string {
	return contract.LogKey
}

func (p *LogProvider) Boot(container goframe.IContainer) error {
	return nil
}

func (p *LogProvider) Params(container goframe.IContainer) []any {
	etcsvc := container.MustMake(contract.EtcKey).(contract.IEtc)

	if p.Formatter == nil {
		p.Formatter = formatter.TextFormatter
		if etcsvc.Exist("log.formatter") {
			v := etcsvc.GetString("log.formatter")
			if v == "json" {
				p.Formatter = formatter.JsonFormatter
			} else if v == "text" {
				p.Formatter = formatter.TextFormatter
			}
		}
	}

	if p.Level == contract.UnknownLevel {
		p.Level = contract.InfoLevel
		if etcsvc.Exist("log.level") {
			p.Level = level(etcsvc.GetString("log.level"))
		}
	}

	return []any{container, p.Level, p.Handler, p.Formatter}
}

func (p *LogProvider) Inject(container goframe.IContainer) goframe.Handler {
	if p.Driver == "" {
		etc, err := container.Make(contract.EtcKey)
		if err != nil {
			return services.NewConsoleLog
		}

		etcsvc := etc.(contract.IEtc)
		p.Driver = strings.ToLower(etcsvc.GetString("log.driver"))
	}

	switch p.Driver {
	case "single":
		return services.NewSingleLog
	case "rotate":
		return services.NewRotateLog
	default:
		return services.NewConsoleLog
	}
}

func (p *LogProvider) IsDefer() bool {
	return false
}

func level(typo string) contract.LogLevel {
	switch strings.ToLower(typo) {
	case "panic":
		return contract.PanicLevel
	case "fatal":
		return contract.FatalLevel
	case "error":
		return contract.ErrorLevel
	case "warn":
		return contract.WarnLevel
	case "info":
		return contract.InfoLevel
	case "debug":
		return contract.DebugLevel
	case "trace":
		return contract.TraceLevel
	}

	return contract.UnknownLevel
}
