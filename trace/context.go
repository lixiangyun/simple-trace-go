package trace

type Context struct {
	TraceID  string `json:"x-trace-id"`
	SpanID   string `json:"x-span-id"`
	ParentID string `json:"x-parent-id"`
}

func NewContext(pctx *Context) *Context {
	if pctx != nil {
		return &Context{TraceID: pctx.TraceID, ParentID: pctx.SpanID, SpanID: GetSpanID()}
	} else {
		return &Context{TraceID: GetTraceID(), SpanID: GetSpanID()}
	}
}
