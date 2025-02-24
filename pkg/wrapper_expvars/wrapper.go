package wrapperexpvars

import (
	"expvar"
	"math/rand"
	"runtime"
	"time"
)

var (
	startTime = time.Now().UTC()
)

func goroutines() interface{} {
	return runtime.NumGoroutine()
}

func uptime() interface{} {
	uptime := time.Since(startTime)
	return int64(uptime)
}

func responseTime() interface{} {
	resp := time.Duration(rand.Intn(1000)) * time.Millisecond
	return int64(resp)
}

func init() {
	expvar.Publish("Goroutines", expvar.Func(goroutines))
	expvar.Publish("Uptime", expvar.Func(uptime))
	expvar.Publish("MeanResponse", expvar.Func(responseTime))
}
