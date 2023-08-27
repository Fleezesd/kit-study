package main

import (
	"context"

	"net"

	"net/http"

	"github.com/fleezesd/kit-study/internal/iam/endpoint"
	"github.com/fleezesd/kit-study/internal/iam/service"
	"github.com/fleezesd/kit-study/internal/iam/transport"
	"github.com/fleezesd/kit-study/internal/pkg/log"
	pb "github.com/fleezesd/kit-study/pkg/proto/iam"
	"github.com/fleezesd/kit-study/pkg/rate"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

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

func HTTPServerRun(lc fx.Lifecycle, handler http.Handler) {
	srv := &http.Server{Addr: ":8080", Handler: handler}
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
	srv := grpc.NewServer(grpc.UnaryInterceptor(grpctransport.Interceptor))
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", ":9090")
			if err != nil {
				log.Fatalw("failed to listen", "err", err)
				return err
			}
			log.Infow("Starting GRPC server at", "addr", ":9090")
			pb.RegisterUserServer(srv, grpcServer)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.GracefulStop()
			return nil
		},
	})
}
