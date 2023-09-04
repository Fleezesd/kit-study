package main

import (
	"context"
	"flag"
	"time"

	"net"

	"net/http"

	"github.com/fleezesd/kit-study/internal/iam/client/etcd"
	"github.com/fleezesd/kit-study/internal/iam/endpoint"
	"github.com/fleezesd/kit-study/internal/iam/service"
	"github.com/fleezesd/kit-study/internal/iam/transport"
	"github.com/fleezesd/kit-study/internal/pkg/log"
	"github.com/fleezesd/kit-study/pkg/opentelemetry/trace"
	pb "github.com/fleezesd/kit-study/pkg/proto/iam"
	"github.com/fleezesd/kit-study/pkg/rate"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var (
	httpAddr = flag.String("h", "127.0.0.1:8080", "httpAddr")
	grpcAddr = flag.String("g", "127.0.0.1:8081", "grpcAddr")
)

func HTTPServerRun(lc fx.Lifecycle, handler http.Handler) {
	flag.Parse()
	srv := &http.Server{Addr: *httpAddr, Handler: handler}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Infow("Starting HTTP server at", "addr", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return
}

func GRPCServerRun(lc fx.Lifecycle, grpcServer pb.UserServer) {
	flag.Parse()
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpctransport.Interceptor))
	// etcd 注册
	var (
		registry etcd.Registry
	)
	// opentelemetry trace provider
	tp, err := trace.NewTraceProvider()
	if err != nil {
		log.Fatalw("make trace provider err", "err", err)
		return
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", *grpcAddr)
			if err != nil {
				log.Fatalw("failed to listen", "err", err)
				return err
			}
			log.Infow("Starting GRPC server at", "addr", *grpcAddr)
			pb.RegisterUserServer(srv, grpcServer)
			// grpc 服务启动
			go srv.Serve(ln)
			// etcd 注册
			registry, err = etcd.RegistryEtcd(*grpcAddr, 10*time.Second)
			if err != nil {
				log.Fatalw("failed to regist etcd", "err", err)
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			defer func() {
				if err := tp.Shutdown(context.Background()); err != nil {
					log.Debugw("Error shutting down tracer provider", "err", err)
				}
			}()
			srv.GracefulStop()
			registry.UnRegistry()
			return nil
		},
	})
}

func main() {
	fx.New(
		fx.Provide(
			log.NewOptions,
			log.NewLogger,
			rate.NewRateLimiter,
			service.NewService,
			endpoint.NewEndPointServer,
			transport.NewHttpHandler,
			transport.NewGRPCServer,
		),
		fx.Invoke(
			HTTPServerRun,
			GRPCServerRun,
		),
	).Run()
}
