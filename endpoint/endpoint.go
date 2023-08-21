package endpoint

import (
	"context"

	"kit-study/service"

	"github.com/go-kit/kit/endpoint"
)

// service层方法 封装到endpoint

type EndPointServer struct {
	HealthEndPoint endpoint.Endpoint
}

func NewEndPointServer(svc service.Service) EndPointServer {
	var healthEndPoint endpoint.Endpoint
	{
		healthEndPoint = MakeHealthEndPoint(svc)
	}
	return EndPointServer{
		HealthEndPoint: healthEndPoint,
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
