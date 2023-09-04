package transport

import (
	"context"

	ep "github.com/fleezesd/kit-study/internal/iam/endpoint"
	"github.com/fleezesd/kit-study/internal/iam/service"
	"github.com/fleezesd/kit-study/internal/pkg/log"
	pb "github.com/fleezesd/kit-study/pkg/proto/iam"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/metadata"
)

var tracer = otel.Tracer("kit-server")

type grpcServer struct {
	login grpctransport.Handler
	// 带上这个 要不grpcserver interface对应接口设定规则不通过
	pb.UnimplementedUserServer
}

func (s *grpcServer) RpcUserLogin(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, rsp, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rsp.(*pb.LoginResponse), nil
}

func GRPCLoginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LoginRequest)
	log.Debugw("解析请求数据", "loginRequest", req)
	return &pb.LoginRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}, nil
}

func GRPCLoginResponse(_ context.Context, grpcRsp interface{}) (interface{}, error) {
	rsp := grpcRsp.(*pb.LoginResponse)
	log.Debugw("请求结束返回值", "loginResponse", rsp)
	return rsp, nil
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func NewGRPCServer(endpoint ep.EndPointServer, logger *log.ZapLogger) pb.UserServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerBefore(func(ctx context.Context, m metadata.MD) context.Context {
			ctx = context.WithValue(ctx, service.ContextReqUUid, m.Get(service.ContextReqUUid))
			// tracer 放入context中
			ctx = context.WithValue(context.Background(), "tracer", tracer)
			return ctx
		}),
		grpctransport.ServerErrorHandler(log.NewLogErrorHandler(logger)),
	}
	return &grpcServer{
		login: grpctransport.NewServer(
			endpoint.LoginEndPoint,
			GRPCLoginRequest,
			GRPCLoginResponse,
			options...,
		),
	}
}
