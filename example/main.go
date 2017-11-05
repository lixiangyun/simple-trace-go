package main

import (
	"log"
	"time"

	"github.com/lixiangyun/simple-trace-go/trace"
)

func fun2(s *trace.Span, a, b int) {
	ep := trace.NewEndPoint("srv2", "192.168.0.2", 1002)

	sp := trace.RecvSpan(s.GetContext())

	sp.Begin(ep)

	sp.AddKV("aaa", "bbb", ep)

	sp.End(ep)

}

func fun1(s *trace.Span, a, b int) {
	ep := trace.NewEndPoint("srv1", "192.168.0.1", 1001)

	sp := trace.NewSpan(s.GetContext())

	sp.Begin(ep)

	sp.AddKV("aa", "bb", ep)

	fun2(sp, a, b)

	sp.End(ep)
}

func main() {

	ep := trace.NewEndPoint("cli", "192.168.0.1", 1001)

	collector := trace.NewCollector()
	if collector == nil {
		log.Println("new trace collector failed!")
		return
	}

	trace.SetCollector(collector)

	sp := trace.Start("test1")

	sp.Begin(ep)

	fun1(sp, 1, 2)

	sp.End(ep)

	time.Sleep(time.Second * 2)

	collector.Stop()

	return
}
