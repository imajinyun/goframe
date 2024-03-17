package formatter

import (
	"bytes"
	"fmt"
	"time"

	"github.com/imajinyun/goframe/contract"
)

func TextFormatter(level contract.LogLevel, t time.Time, msg string, fields map[string]any) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	sep := "\t"

	buf.WriteString(Prefix(level))
	buf.WriteString(sep)
	buf.WriteString(t.Format(time.RFC3339))
	buf.WriteString(sep)

	buf.WriteString("\"")
	buf.WriteString(msg)
	buf.WriteString("\"")
	buf.WriteString(sep)

	buf.WriteString(fmt.Sprint(fields))

	return buf.Bytes(), nil
}
