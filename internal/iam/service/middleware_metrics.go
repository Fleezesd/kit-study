package service

import (
	"context"
	"time"

	pb "github.com/fleezesd/kit-study/pkg/proto/iam"
	"github.com/go-kit/kit/metrics"
	metricsprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

type metricsMiddlewareServer struct {
	next      Service
	counter   metrics.Counter   // 计数器
	histogram metrics.Histogram // 柱状图
}

func NewMetricsMiddlewareServer() NewMiddlewareServer {
	// 本质还是prometheus 初始化一个counter一样
	counter := metricsprometheus.NewCounterFrom(prometheus.CounterOpts{
		Subsystem:   "user_agent",
		Name:        "request_count",
		Help:        "Number of requests",
		ConstLabels: map[string]string{},
	}, []string{"method"})
	histogram := metricsprometheus.NewHistogramFrom(prometheus.HistogramOpts{
		Subsystem: "user_agent",
		Name:      "request_consume",
		Help:      "Request consumes time",
	}, []string{"method"})
	return func(service Service) Service {
		return &metricsMiddlewareServer{
			next:      service,
			counter:   counter,
			histogram: histogram,
		}
	}
}

func (m *metricsMiddlewareServer) Health(ctx context.Context, request interface{}) (rsp interface{}, err error) {
	rsp, err = m.next.Health(ctx, request)
	return
}
func (m *metricsMiddlewareServer) Login(ctx context.Context, req *pb.LoginRequest) (rsp *pb.LoginResponse, err error) {
	// 统计请求次数
	defer func(start time.Time) {
		method := []string{"method", "login"}
		// 统计请求个数
		m.counter.With(method...).Add(1)
		// 统计请求耗时
		m.histogram.With(method...).Observe(float64(time.Since(start).Seconds()))
	}(time.Now())
	rsp, err = m.next.Login(ctx, req)
	return
}
