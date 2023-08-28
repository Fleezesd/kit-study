package etcd

import (
	"context"
	"io"
	"time"

	"github.com/go-kit/kit/sd/lb"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"

	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/log"
)

func NewUserAgentClient(addr []string) {
	var (
		etcdAddrs = addr
		prefix    = "svc.user.Addr"
		instance  = "localhost:9090" // grpc addr
		key       = prefix + instance
		value     = "http://" + instance // based on our transport
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
		panic(err)
	}
	// Build the registrar.
	registrar := etcdv3.NewRegistrar(etcdClient, etcdv3.Service{
		Key:   key,
		Value: value,
	}, log.NewNopLogger())
	// 本地grpc服务节点注册到服务器 etcdAddrs
	registrar.Register()

	defer registrar.Deregister()

	// 我们可能还想连接到其他服务并使用他们的方法。
	// 我们可以构建一个实例器 instancer来监听来自etcd的更改
	// 创建端点endpointer,用负载平衡器balancer包装它来选择一个端点，最后用重试策略包装它，以获得可以直接用作端点。
	barPrefix := "/services/barsvc"
	logger := log.NewNopLogger()
	// 构建实例器 instancer 监听etcd的更改
	instancer, err := etcdv3.NewInstancer(etcdClient, barPrefix, logger)
	if err != nil {
		panic(err)
	}
	// NewEndpointer creates an Endpointer that subscribes to updates(订阅更新) from Instancer src and uses factory f to create Endpoints.
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	// NewRoundRobin returns a load balancer that returns services in sequence (按顺序返回服务).
	balancer := lb.NewRoundRobin(endpointer)
	// retry is endpoint It represents a single RPC method.
	retry := lb.Retry(3, 3*time.Second, balancer)

	req := struct{}{}
	if _, err = retry(ctx, req); err != nil {
		panic(err)
	}
}

// 工厂方法 创建endpointer
func factory(instance string) (endpoint.Endpoint, io.Closer, error) {
	return endpoint.Nop, nil, nil
}
