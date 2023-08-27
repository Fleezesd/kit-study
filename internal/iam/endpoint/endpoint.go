package endpoint

import (
	"context"

	"github.com/fleezesd/kit-study/internal/iam/service/dto"

	"go.uber.org/ratelimit"

	"github.com/fleezesd/kit-study/internal/iam/service"

	"github.com/go-kit/kit/endpoint"

	pb "github.com/fleezesd/kit-study/pkg/proto/iam"
)

// service层方法 封装到endpoint

type EndPointServer struct {
	HealthEndPoint endpoint.Endpoint
	LoginEndPoint  endpoint.Endpoint
}

func NewEndPointServer(svc service.Service, limiter ratelimit.Limiter) EndPointServer {
	var healthEndPoint endpoint.Endpoint
	{
		healthEndPoint = MakeHealthEndPoint(svc)
		healthEndPoint = logMiddleware()(healthEndPoint)
		healthEndPoint = AuthMiddleware()(healthEndPoint)
		healthEndPoint = UberRateMiddleware(limiter)(healthEndPoint)
	}
	var loginEndPoint endpoint.Endpoint
	{
		loginEndPoint = MakeLoginEndPoint(svc)
		loginEndPoint = logMiddleware()(loginEndPoint)
		loginEndPoint = UberRateMiddleware(limiter)(loginEndPoint)
	}
	return EndPointServer{
		HealthEndPoint: healthEndPoint,
		LoginEndPoint:  loginEndPoint,
	}
}

func MakeHealthEndPoint(s service.Service) endpoint.Endpoint {
	//封装 svc.Health
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := s.Health(ctx, request)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func (s EndPointServer) Health(ctx context.Context, request interface{}) (res interface{}, err error) {
	// service层方法 直接上抛到endpoint层
	res, err = s.HealthEndPoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func MakeLoginEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(dto.LoginRequest); ok {
			req = request.(dto.LoginRequest)
			return s.Login(ctx, req)
		} else {
			// 为了直接使用 service 层的方法 需要一层grpc结构转换
			grpcReq := request.(*pb.LoginRequest)
			req = dto.LoginRequest{
				Username: grpcReq.Username,
				Password: grpcReq.Password,
			}
			return s.Login(ctx, req)
		}
	}
}

func (s EndPointServer) Login(ctx context.Context, req dto.LoginRequest) (rsp dto.LoginResponse, err error) {
	res, err := s.LoginEndPoint(ctx, req)
	rsp = res.(dto.LoginResponse)
	return
}
