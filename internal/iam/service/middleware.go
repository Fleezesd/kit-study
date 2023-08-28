package service

import (
	"context"
	"fmt"

	"github.com/fleezesd/kit-study/internal/pkg/log"
	pb "github.com/fleezesd/kit-study/pkg/proto/iam"
)

// 抽象:对应 Service 安装中间件 (serivce加一层装饰)

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

func (l *logMiddleware) Health(ctx context.Context, request interface{}) (rsp interface{}, err error) {
	// log 装饰  记录调用结果
	defer func() {
		log.Debugw(fmt.Sprint(ctx.Value(ContextReqUUid)), "logmiddleware", "service-health", "res", rsp)
	}()
	rsp, err = l.next.Health(ctx, request)
	if err != nil {
		return "", err
	}
	return rsp, nil
}

func (l *logMiddleware) Login(ctx context.Context, req *pb.LoginRequest) (rsp *pb.LoginResponse, err error) {
	defer func() {
		log.Debugw(fmt.Sprint(ctx.Value(ContextReqUUid)), "logmiddleware", "service-login", "res", rsp)
	}()
	rsp, err = l.next.Login(ctx, req)
	return
}
