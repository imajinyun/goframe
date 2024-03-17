package trace

import (
	"context"
	"net/http"
	"time"

	"github.com/imajinyun/goframe/contract"
	"github.com/imajinyun/goframe/gin"
	"github.com/imajinyun/goframe/util"

	"github.com/imajinyun/goframe"
)

type TraceKey string

var ContextKey = TraceKey("trace-key")

type TraceService struct {
	uid   contract.IUid
	span  contract.IUid
	trace contract.IUid
}

func NewTraceService(params ...any) (any, error) {
	container := params[0].(goframe.IContainer)
	uidsvc := container.MustMake(contract.UidKey).(contract.IUid)

	return &TraceService{uid: uidsvc}, nil
}

func (s *TraceService) GetTrace(ctx context.Context) *contract.TraceContext {
	if gctx, ok := ctx.(*gin.Context); ok {
		if val, ok := gctx.Get(string(ContextKey)); ok {
			return val.(*contract.TraceContext)
		}
	}

	if tc, ok := ctx.Value(ContextKey).(*contract.TraceContext); ok {
		return tc
	}

	return nil
}

func (s *TraceService) NewTrace() *contract.TraceContext {
	var tid, sid string

	tid = util.If(s.trace != nil, s.trace.NewUid(), s.uid.NewUid()).(string)
	sid = util.If(s.span != nil, s.span.NewUid(), s.uid.NewUid()).(string)
	tcx := &contract.TraceContext{
		TraceID:    tid,
		ParentID:   "",
		SpanID:     sid,
		SubSpanID:  "",
		Annotation: map[string]string{},
	}

	return tcx
}

func (s *TraceService) StartSpan(tcx *contract.TraceContext) *contract.TraceContext {
	var ssid string
	ssid = util.If(s.span != nil, s.span.NewUid(), s.uid.NewUid()).(string)
	span := &contract.TraceContext{
		TraceID:    tcx.TraceID,
		ParentID:   "",
		SpanID:     tcx.SpanID,
		SubSpanID:  ssid,
		Annotation: map[string]string{contract.TraceKeyTime: time.Now().String()},
	}

	return span
}

func (s *TraceService) WithTrace(ctx context.Context, tcx *contract.TraceContext) context.Context {
	if gctx, ok := ctx.(*gin.Context); ok {
		gctx.Set(string(ContextKey), tcx)
		return gctx
	} else {
		return context.WithValue(ctx, ContextKey, tcx)
	}
}

func (s *TraceService) ExtractHttp(req *http.Request) *contract.TraceContext {
	tcx := &contract.TraceContext{}
	tcx.TraceID = req.Header.Get(contract.TraceKeyTraceID)
	tcx.ParentID = req.Header.Get(contract.TraceKeyParentID)
	tcx.SpanID = req.Header.Get(contract.TraceKeySpanID)
	tcx.SubSpanID = ""

	if tcx.TraceID == "" {
		tcx.TraceID = s.uid.NewUid()
	}

	if tcx.SpanID == "" {
		tcx.SpanID = s.uid.NewUid()
	}

	return tcx
}

func (s *TraceService) InjectHttp(req *http.Request, tcx *contract.TraceContext) *http.Request {
	req.Header.Add(contract.TraceKeyTraceID, tcx.TraceID)
	req.Header.Add(contract.TraceKeySpanID, tcx.SpanID)
	req.Header.Add(contract.TraceKeySubSpanID, tcx.SubSpanID)
	req.Header.Add(contract.TraceKeyParentID, tcx.ParentID)

	return req
}

func (s *TraceService) ToMap(tcx *contract.TraceContext) map[string]string {
	m := map[string]string{}
	if tcx == nil {
		return m
	}
	m[contract.TraceKeyTraceID] = tcx.TraceID
	m[contract.TraceKeySpanID] = tcx.SpanID
	m[contract.TraceKeySubSpanID] = tcx.SubSpanID
	m[contract.TraceKeyParentID] = tcx.ParentID

	if tcx.Annotation != nil {
		for k, v := range tcx.Annotation {
			m[k] = v
		}
	}

	return m
}
