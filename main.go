package main

import (
	"context"
	"kit-study/internal/iam/endpoint"
	"kit-study/internal/iam/service"
	"kit-study/internal/iam/transport"
	"kit-study/pkg/rate"
	"net"


	"net/http"

	"kit-study/internal/pkg/log"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			rate.NewRateLimiter,
			service.NewService,
			endpoint.NewEndPointServer,
			transport.NewHttpHandler,
		),
		fx.Invoke(
			NewHTTPServer,
		),
	).Run()
}

func NewHTTPServer(lc fx.Lifecycle, handler http.Handler) *http.Server {
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
	return srv
}
