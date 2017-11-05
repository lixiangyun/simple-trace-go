package trace

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	tm := time.Now()
	rand.Seed(int64(tm.Nanosecond()))
}

const (
	default_traceid_len = 128
	default_spanid_len  = 64
)

func getuuid(length int) string {
	length = length / (8 * 4)
	var output string
	for i := 0; i < length; i++ {
		output += fmt.Sprintf("%08x", rand.Uint32())
	}
	return output
}

func getTraceID() string {
	return getuuid(default_traceid_len)
}

func getSpanID() string {
	return getuuid(default_spanid_len)
}
