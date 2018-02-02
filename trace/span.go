package trace

const (
	CLIENT = "CLIENT"
	SERVER = "SERVER"
)

type Context struct {
	TraceID  string `json:"x-trace-id"`
	SpanID   string `json:"x-span-id"`
	ParentID string `json:"x-parent-id"`
}

type Span struct {
	kind      string
	name      string
	ctx       Context
	edp       Endpoint
	tags      interface{}
	timestamp int64
	event     []Annotations
}

type Annotations struct {
	Timestamp int64  `json:"timestamp"` // UNIX时间，单位毫秒
	Value     string `json:"value"`     // 事件内容
}

type Endpoint struct {
	Name string `json:"serviceName"`
	IPv4 string `json:"ipv4"`
	Port int    `json:"port"`
}

type Tag struct {
	Kv interface{}
}

type SpanRecord struct {
	TraceID   string `json:"traceId"`   // 调用ID
	SpanID    string `json:"id"`        // 当前spanID
	ParentID  string `json:"parentId"`  // 父spanID
	Name      string `json:"name"`      // 跟踪名称: http method: post/get...
	Kind      string `json:"kind"`      // 类型: CLIENT/SERVER
	Timestamp int64  `json:"timestamp"` // UNIX时间，单位毫秒
	Duration  int64  `json:"duration"`  // 时间间隔，单位毫秒

	LocalEndpoint interface{} `json:"localEndpoint"`
}

func NewEndPoint(name, ip string, port int) Endpoint {
	return Endpoint{Name: name, IPv4: ip, Port: port}
}

func NewSpan(pctx *Context, kind, method string, edp Endpoint) *Span {
	if kind != CLIENT && kind != SERVER {
		return nil
	}
	if pctx == nil || pctx.TraceID == "" {
		return nil
	}
	return &Span{kind: kind, name: method, ctx: *pctx, edp: edp}
}

func (s *Span) Begin() {
	s.timestamp = GetTimeStamp()
}

func (s *Span) Add(value string) {
	if s.event == nil {
		s.event = make([]Annotations, 0)
	}
	s.event = append(s.event, Annotations{Timestamp: GetTimeStamp(), Value: value})
}

func (s *Span) GetContext() *Context {
	return &s.ctx
}

func (s *Span) Tags(kv interface{}) {
	s.tags = kv
}

func (s *Span) End() {

	timestamp := GetTimeStamp()

	span := new(SpanRecord)

	span.TraceID = s.ctx.TraceID
	span.SpanID = s.ctx.SpanID
	span.ParentID = s.ctx.ParentID

	span.Name = s.name
	span.Kind = s.kind
	span.Timestamp = s.timestamp
	span.Duration = timestamp - s.timestamp

	span.LocalEndpoint = s.edp

	Collector(span)
}

func NewContext(pctx *Context) *Context {
	if pctx != nil {
		return &Context{
			TraceID:  pctx.TraceID,
			ParentID: pctx.SpanID,
			SpanID:   GetSpanID()}
	} else {
		return &Context{
			TraceID: GetTraceID(),
			SpanID:  GetSpanID()}
	}
}
