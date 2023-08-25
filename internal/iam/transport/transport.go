package transport

import (
	"context"
	"kit-study/pkg/token"
	"net/http"

	e "kit-study/internal/iam/endpoint"

	"kit-study/internal/iam/service"
	"kit-study/internal/pkg/log"

	httptransport "github.com/go-kit/kit/transport/http"
	uuid "github.com/satori/go.uuid"
)

// http Handler
func NewHttpHandler(endpoint e.EndPointServer) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder), //程序中的全部报错都会走这里面
		httptransport.ServerBefore(func(ctx context.Context, request *http.Request) context.Context { // 添加middleware 增加请求的uuid
			UUID := uuid.NewV5(uuid.NewV4(), "req_uuid").String()
			ctx = context.WithValue(ctx, service.ContextReqUUid, UUID)
			log.Debugw("给请求添加uuid", "UUID", UUID)
			ctx = context.WithValue(ctx, token.JWT_CONTEXT_KEY, request.Header.Get("Authorization"))
			log.Debugw("把请求中的token发送到context中", "Token", request.Header.Get("Authorization"))
			return ctx
		}),
	}

	m := http.NewServeMux()
	m.Handle("/health", httptransport.NewServer(
		endpoint.HealthEndPoint,
		decodeHTTPHealthRequest,   //解析请求值
		encodeHTTPGenericResponse, //返回值
		options...,
	))
	m.Handle("/login", httptransport.NewServer(
		endpoint.LoginEndPoint,
		decodeHTTPLoginRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	return m
}
