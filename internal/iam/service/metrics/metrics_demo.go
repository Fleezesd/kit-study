package metrics

import (
	"context"
	"runtime"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/expvar"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/metrics/statsd"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// Counter
func CounterExample() {
	var myCount metrics.Counter
	// expvar 类型
	myCount = expvar.NewCounter("my_count")
	myCount.Add(1)
}

// Histogram
func APIReqTimeHistogram() {
	// 请求持续时间的直方图					 prometheus 的 summary 对象
	var dur metrics.Histogram = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "myservice",
		Subsystem: "api",
		Name:      "request_duration_seconds",
		Help:      "Total time spent serving requests.",
	}, []string{})
	go handleRequest(dur)
	// ...
}

func handleRequest(dur metrics.Histogram) {
	//
	defer func(begin time.Time) { dur.Observe(time.Since(begin).Seconds()) }(time.Now())
	// handle request
}

// StasD
func GoRoutineNumStatsD() {
	// 当前运行的 goroutine 数量的计量表，通过 StatsD 导出
	statsd := statsd.New("foo_svc.", log.NewNopLogger())
	report := time.NewTicker(5 * time.Second)
	defer report.Stop()
	go statsd.SendLoop(context.Background(), report.C, "tcp", "statsd.internal:8125")
	goroutines := statsd.NewGauge("goroutine_count")
	go exportGoroutines(goroutines)
	// ...
}

func exportGoroutines(g metrics.Gauge) {
	// Guage 测量器 递增递减均可
	for range time.Tick(time.Second) {
		g.Set(float64(runtime.NumGoroutine()))
	}
}
