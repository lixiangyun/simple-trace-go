package trace

import "log"

const defaultBatchSize = 100 // the counts of span record

var spanRecordChan chan *SpanRecord

func init() {
	spanRecordChan = make(chan *SpanRecord, defaultBatchSize)
	go collectorLoop()
}

func collectorLoop() {

	for {
		span := <-spanRecordChan
		count := len(spanRecordChan) + 1

		buffer := make([]*SpanRecord, count)
		buffer[0] = span

		for i := 1; i < count; i++ {
			buffer[i] = <-spanRecordChan
		}

		err := PostSpan(buffer)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func Collector(s *SpanRecord) {
	spanRecordChan <- s
}
