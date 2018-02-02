package main

import (
	"time"

	"github.com/lixiangyun/trace-go/trace"
)

func fun2(ctx *trace.Context, a, b int) {
	ep := trace.NewEndPoint("srv2", "192.168.0.2", 1002)

	sp := trace.NewSpan(ctx, trace.SERVER, "fun2", ep)
	if sp != nil {
		sp.Begin()
	}

	time.Sleep(100 * time.Millisecond)

	if sp != nil {
		sp.End()
	}

}

func fun1(ctx *trace.Context, a, b int) {
	ep := trace.NewEndPoint("srv1", "192.168.0.1", 1001)

	sp := trace.NewSpan(ctx, trace.SERVER, "fun1", ep)
	if sp != nil {
		sp.Begin()
	}

	time.Sleep(50 * time.Millisecond)

	ctx2 := trace.NewContext(ctx)

	sp1 := trace.NewSpan(ctx2, trace.CLIENT, "fun2", ep)
	if sp1 != nil {
		sp1.Begin()
	}

	fun2(ctx2, a, b)

	if sp1 != nil {
		sp1.End()
	}

	time.Sleep(50 * time.Millisecond)

	if sp != nil {
		sp.End()
	}
}

func main() {

	trace.ZipKinEndpointSet("127.0.0.1:9411")

	ep := trace.NewEndPoint("cli", "192.168.0.1", 1001)

	ctx := trace.NewContext(nil)

	sp := trace.NewSpan(ctx, trace.CLIENT, "fun1", ep)
	if sp != nil {
		sp.Begin()
	}

	fun1(ctx, 1, 2)

	if sp != nil {
		sp.End()
	}

	time.Sleep(time.Second * 1)
}
