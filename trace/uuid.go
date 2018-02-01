package trace

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(int64(time.Now().UnixNano()))
}

const (
	default_traceid_len = 128
	default_spanid_len  = 64
)

func GetTimeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}

func getuuid(length int) string {
	length = length / (8 * 4)
	var output string
	for i := 0; i < length; i++ {
		output += fmt.Sprintf("%08x", rand.Uint32())
	}
	return output
}

func GetTraceID() string {
	return getuuid(default_traceid_len)
}

func GetSpanID() string {
	return getuuid(default_spanid_len)
}
