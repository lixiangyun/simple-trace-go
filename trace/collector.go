package trace

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

const defaultBatchInterval = 1
const defaultBatchSize = 100
const defaultFileName = "collector.txt"

type Collector struct {
	batchInterval time.Duration
	batchsize     int
	recordfile    *os.File
	buffer        []*SpanRecord
	spanc         chan *SpanRecord
	stop          chan struct{}
	shutdown      chan struct{}
}

func NewCollector() *Collector {
	c := new(Collector)

	c.batchsize = defaultBatchSize
	c.batchInterval = defaultBatchInterval * time.Second

	c.spanc = make(chan *SpanRecord, defaultBatchSize)
	c.buffer = make([]*SpanRecord, 0)
	c.stop = make(chan struct{}, 1)
	c.shutdown = make(chan struct{}, 1)

	file, err := os.Create(defaultFileName)
	if err != nil {
		fmt.Println("create file error!")
		return nil
	}

	fileinfo, err := file.Stat()
	if err == nil {
		file.Seek(fileinfo.Size(), 0)
	}

	c.recordfile = file

	go loop(c)

	return c
}

func loop(c *Collector) {
	timeout := time.NewTimer(c.batchInterval)

	for {
		var bSend bool

		select {
		case span := <-c.spanc:
			c.buffer = append(c.buffer, span)
			if len(c.buffer) >= c.batchsize {
				bSend = true
			}
		case <-timeout.C:
			bSend = true
		case <-c.stop:
			c.send()
			close(c.shutdown)
			log.Println("collector close.")
			return
		}

		if bSend {
			c.send()
			timeout.Reset(c.batchInterval)
		}
	}
}

func (c *Collector) send() error {

	for _, span := range c.buffer {
		buf, err := json.Marshal(span)
		if err != nil {
			log.Println("send failed!", err.Error())
		} else {
			c.recordfile.Write(buf)
		}
	}

	c.buffer = make([]*SpanRecord, 0)

	return nil
}

func (c *Collector) Record(s *SpanRecord) {
	c.spanc <- s
}

func (c *Collector) Stop() {
	close(c.stop)
	<-c.shutdown
	c.recordfile.Close()
}
