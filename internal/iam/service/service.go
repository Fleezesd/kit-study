package service

import (
	"context"
	"fmt"
	"kit-study/internal/iam/service/dto"
	"kit-study/internal/pkg/errno"
	"kit-study/internal/pkg/log"
	"kit-study/pkg/auth"
	"kit-study/pkg/token"
)

type Service interface {
	Health(ctx context.Context) (string, error)
	Login(ctx context.Context, req dto.LoginRequest) (rsp dto.LoginResponse, err error)
}

func NewService() Service {
	var server Service
	server = &baseServer{}
	server = NewLogMiddlewareServer()(server)
	return server
}

type baseServer struct {
}

// 确保 baseServer 实现了接口
var _ Service = (*baseServer)(nil)

func (s baseServer) Health(ctx context.Context) (string, error) {
	log.Debugw(fmt.Sprint(ctx.Value(ContextReqUUid), "service", "health"))
	return fmt.Sprintln("health"), nil
}

func (s baseServer) Login(ctx context.Context, req dto.LoginRequest) (rsp dto.LoginResponse, err error) {
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
		return rsp, errno.ErrSignToken
	}
	return
}