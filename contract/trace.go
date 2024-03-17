package contract

import (
	"context"
	"net/http"
)

const TraceKey = "gogin:trace"

const (
	TraceKeyTraceID   = "trace_id"
	TraceKeyParentID  = "parent_id"
	TraceKeySpanID    = "span_id"
	TraceKeySubSpanID = "sub_span_id"
	TraceKeyMethod    = "method"
	TraceKeyCaller    = "caller"
	TraceKeyTime      = "time"
)

type ITrace interface {
	NewTrace() *TraceContext
	GetTrace(ctx context.Context) *TraceContext
	WithTrace(ctx context.Context, tcx *TraceContext) context.Context
	StartSpan(tcx *TraceContext) *TraceContext
	ToMap(tcx *TraceContext) map[string]string
	ExtractHttp(req *http.Request) *TraceContext
	InjectHttp(req *http.Request, tcx *TraceContext) *http.Request
}

type TraceContext struct {
	TraceID    string
	ParentID   string
	SpanID     string
	SubSpanID  string
	Annotation map[string]string
}
