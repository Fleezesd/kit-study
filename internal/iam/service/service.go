package service

import (
	"context"
	"fmt"

	"github.com/fleezesd/kit-study/internal/pkg/errno"
	"github.com/fleezesd/kit-study/internal/pkg/log"
	"github.com/fleezesd/kit-study/pkg/auth"
	pb "github.com/fleezesd/kit-study/pkg/proto/iam"
	"github.com/fleezesd/kit-study/pkg/token"
	"go.opentelemetry.io/contrib/instrumentation/github.com/go-kit/kit/otelkit"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Service interface {
	Health(ctx context.Context, request interface{}) (rsp interface{}, err error)
	Login(ctx context.Context, req *pb.LoginRequest) (rsp *pb.LoginResponse, err error)
}

func NewService() Service {
	var server Service
	server = &baseServer{}
	server = NewLogMiddlewareServer()(server)
	server = NewMetricsMiddlewareServer()(server)
	return server
}

type baseServer struct {
}

// 确保 baseServer 实现了接口
var _ Service = (*baseServer)(nil)

func (s baseServer) Health(ctx context.Context, request interface{}) (res interface{}, err error) {
	log.Debugw(fmt.Sprint(ctx.Value(ContextReqUUid), "service", "health"))
	return fmt.Sprintln("health"), nil
}

func (s baseServer) Login(ctx context.Context, req *pb.LoginRequest) (rsp *pb.LoginResponse, err error) {
	tracer := ctx.Value("tracer").(trace.Tracer)
	otelkit.WithOperation("login")
	_, span := tracer.Start(ctx, "login", oteltrace.WithAttributes(attribute.String("username", req.GetUsername())))
	defer span.End()
	log.Debugw("loginEndpoint span", "span", span)
	rsp = &pb.LoginResponse{}
	log.Debugw(fmt.Sprint(ctx.Value(ContextReqUUid)), "service", "login")
	if req.Username != "admin" {
		return rsp, errno.ErrUserNotFound
	}
	// bcrypt-hash加密
	if err := auth.Compare("$2a$10$RJrPY12wL8uTl.o7gdQtburp5Y9VxFODYiwhaBrQJxbC7jhEfYjbC", req.Password); err != nil {
		return rsp, errno.ErrPasswordIncorrect
	}
	rsp.Token, err = token.Sign(req.Username)
	if err != nil {
		return nil, errno.ErrSignToken
	}
	return
}
