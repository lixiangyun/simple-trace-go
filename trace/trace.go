package trace

func Start(name string) *Span {

	ctx := Context{TraceID: getTraceID(),
		TraceName: name,
		SpanID:    getSpanID(),
		ParentID:  ""}

	return RecvSpan(ctx)
}

func SetCollector(c *Collector) {
	globalCollector = c
}
