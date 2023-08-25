package endpoint

import (
	"context"
	"kit-study/internal/iam/service/dto"

	"kit-study/internal/iam/service"

	"github.com/go-kit/kit/endpoint"
)

// service层方法 封装到endpoint

type EndPointServer struct {
	HealthEndPoint endpoint.Endpoint
	LoginEndPoint  endpoint.Endpoint
}

func NewEndPointServer(svc service.Service) EndPointServer {
	var healthEndPoint endpoint.Endpoint
	{
		healthEndPoint = MakeHealthEndPoint(svc)
		healthEndPoint = logMiddleware()(healthEndPoint)
	}
	var loginEndPoint endpoint.Endpoint
	{
		loginEndPoint = MakeLoginEndPoint(svc)
		loginEndPoint = logMiddleware()(loginEndPoint)
	}
	return EndPointServer{
		HealthEndPoint: healthEndPoint,
		LoginEndPoint:  loginEndPoint,
	}
}

func MakeHealthEndPoint(s service.Service) endpoint.Endpoint {
	//封装 svc.Health
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := s.Health(ctx)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func (s EndPointServer) Health(ctx context.Context, request interface{}) (res interface{}) {
	// service层方法 直接上抛到endpoint层
	res, _ = s.HealthEndPoint(ctx, request)
	return res
}

func MakeLoginEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(dto.LoginRequest)
		return s.Login(ctx, req)
	}
}

func (s EndPointServer) Login(ctx context.Context, req dto.LoginRequest) (rsp dto.LoginResponse, err error) {
	res, err := s.LoginEndPoint(ctx, req)
	rsp = res.(dto.LoginResponse)
	return
}
