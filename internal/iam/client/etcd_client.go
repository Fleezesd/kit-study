package client

import (
	"context"
	"time"

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
}
