package transport

import (
	"net/http"

	e "kit-study/endpoint"

	httptransport "github.com/go-kit/kit/transport/http"
)


// http Handler
func NewHttpHandler(endpoint e.EndPointServer) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder), //程序中的全部报错都会走这里面
	}

	m := http.NewServeMux()
	m.Handle("/health", httptransport.NewServer(
		endpoint.HealthEndPoint,
		decodeHTTPHealthRequest,   //解析请求值
		encodeHTTPGenericResponse, //返回值
		options...,
	))
	return m
}
