package client

import (
	"context"
	"io"
	"time"

	"github.com/fleezesd/kit-study/internal/iam/client/etcd"

	ep "github.com/fleezesd/kit-study/internal/iam/endpoint"

	"github.com/go-kit/kit/sd"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/fleezesd/kit-study/internal/iam/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"

	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
)

// 为grpc客户端添加sd和lb能力

type UserAgent struct {
	instancer *etcdv3.Instancer
	logger    log.Logger
}

// sd 服务发现
func NewUserAgentClient(addr []string, logger log.Logger) (*UserAgent, error) {
	var (
		etcdAddrs = addr
		prefix    = "/registry/server"
		ttl       = 5 * time.Second
		ctx       = context.Background()
	)
	options := etcdv3.ClientOptions{
		DialTimeout:   ttl,
		DialKeepAlive: ttl,
	}
	// Build the client use kit etcd options
	etcdClient, err := etcdv3.NewClient(ctx, etcdAddrs, options)
	if err != nil {
		return nil, err
	}
	// 构建实例器 instancer 监听etcd的更改
	instancer, err := etcdv3.NewInstancer(etcdClient, prefix, logger)
	if err != nil {
		return nil, err
	}
	return &UserAgent{
		instancer: instancer,
		logger:    logger,
	}, err

}

// grpc 负载均衡
func (u *UserAgent) UserAgentClient() (service.Service, error) {
	var (
		endpoints ep.EndPointServer
	)
	{
		// 工厂方法 作为创建endpointer的入参
		factory := u.Factory(ep.MakeLoginEndPoint)
		// kit 底层 	service, closer, err := c.factory(instance) 来确定grpc service
		endpointer := sd.NewEndpointer(u.instancer, factory, u.logger)
		balancer := lb.NewRoundRobin(endpointer)
		// retry is endpoint It represents a single RPC method.对endpoint的包装
		retry := lb.Retry(1, 3*time.Second, balancer)
		endpoints.LoginEndPoint = retry
	}
	return endpoints, nil
}

// 工厂方法 创建endpoints的工厂方法
func (u *UserAgent) Factory(makeEndpoint func(service.Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		instance = etcd.GetAddr(instance)
		conn, err := grpc.Dial(instance, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, nil, err
		}
		srv := u.NewGRPCClient(conn)
		endpoints := makeEndpoint(srv)
		return endpoints, conn, nil
	}
}
