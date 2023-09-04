package client

import (
	"context"

	ep "github.com/fleezesd/kit-study/internal/iam/endpoint"
	"github.com/fleezesd/kit-study/internal/iam/service"
	"github.com/fleezesd/kit-study/internal/iam/transport"
	"github.com/fleezesd/kit-study/internal/pkg/log"
	pb "github.com/fleezesd/kit-study/pkg/proto/iam"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NewGRPCClient returns a Service backed by a gRPC server at the other end of the conn
func (u *UserAgent) NewGRPCClient(conn *grpc.ClientConn) service.Service {
	// global client middlewares  client的全局中间件
	options := []grpctransport.ClientOption{
		grpctransport.ClientBefore(func(ctx context.Context, m *metadata.MD) context.Context {
			UUID := uuid.NewV5(uuid.Must(uuid.NewV4(), nil), service.ContextReqUUid).String()
			log.Debugw("跟请求添加uuid", "UUID", UUID)
			m.Set(service.ContextReqUUid, UUID)
			//  将 UUID 传递给服务器端
			ctx = metadata.NewOutgoingContext(context.Background(), *m)
			return ctx
		}),
	}

	var loginEndPoint endpoint.Endpoint
	{
		/*
			cc *grpc.ClientConn,
			serviceName string,
			method string,
			enc EncodeRequestFunc,
			dec DecodeResponseFunc,
			grpcReply interface{},
			options ...ClientOption,
		*/
		loginEndPoint = grpctransport.NewClient(
			conn,
			"pb.User",
			"RpcUserLogin",
			transport.GRPCLoginRequest,
			transport.GRPCLoginResponse,
			pb.LoginResponse{},
			options...,
		).Endpoint()

	}
	return ep.EndPointServer{
		LoginEndPoint: loginEndPoint,
	}
}
