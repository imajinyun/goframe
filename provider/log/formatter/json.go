package formatter

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/imajinyun/goframe/contract"
)

func JsonFormatter(level contract.LogLevel, t time.Time, msg string, fields map[string]any) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	if fields == nil {
		fields = make(map[string]any)
	}
	fields["msg"] = msg
	fields["level"] = level
	fields["time"] = t.Format(time.RFC3339)
	byt, err := json.Marshal(fields)
	if err != nil {
		return buf.Bytes(), errors.Wrap(err, "json format error")
	}

	buf.Write(byt)

	return buf.Bytes(), nil
}
