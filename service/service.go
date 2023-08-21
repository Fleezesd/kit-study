package service

import (
	"context"
	"fmt"
)

type Service interface {
	Health(ctx context.Context) (string, error)
}

func NewService() Service {
	return &serviceServer{}
}

type serviceServer struct {
}

// 确保 baseServer 实现了接口
var _ Service = (*serviceServer)(nil)

func (s serviceServer) Health(ctx context.Context) (string, error) {
	return fmt.Sprintln("health"), nil
}
