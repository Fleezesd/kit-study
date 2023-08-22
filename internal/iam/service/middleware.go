package service

import (
	"context"
	"fmt"
	"kit-study/internal/pkg/log"
)

// 抽象:对应 Service 安装中间件 (serivce加一层装饰)
// 1. 日志中间件

const ContextReqUUid = "req_uuid"

type NewMiddleware func(Service) Service

type logMiddleware struct {
	next Service
}

func NewLogMiddlewareServer() NewMiddleware {
	return func(service Service) Service {
		return &logMiddleware{
			next: service,
		}
	}
}

func (l *logMiddleware) Health(ctx context.Context) (out string, err error) {
	// log 装饰  记录调用结果
	defer func() {
		log.Debugw(fmt.Sprint(ctx.Value(ContextReqUUid)), "调用service层logMiddleware", "service:health", "res", out)
	}()
	out, err = l.next.Health(ctx)
	if err != nil {
		return "", err
	}
	return out, nil
}
