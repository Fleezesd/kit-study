package main

import (
	"context"
	"fmt"
	"kit-study/endpoint"
	"kit-study/service"
	"kit-study/transport"
	"net"

	"net/http"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			endpoint.NewEndPointServer,
			service.NewService,
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
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
