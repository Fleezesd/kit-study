package service

import (
	"context"
	"fmt"
)

type Service interface {
	Health(ctx context.Context) (string, error)
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
	return fmt.Sprintln("health"), nil
}
