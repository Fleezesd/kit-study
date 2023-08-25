package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"kit-study/internal/iam/service"
	"kit-study/internal/iam/service/dto"
	"kit-study/internal/pkg/log"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

// 解析请求HTTP & 数据转换 dto & 处理endpoint层的返回值

// 解析 HTTP 请求
func decodeHTTPHealthRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	// 此处可解析HTTP请求
	return nil, nil
}

func decodeHTTPLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var login dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&login) // 从body体中拿取
	if err != nil {
		return nil, err
	}
	log.Debugw(fmt.Sprint(ctx.Value(service.ContextReqUUid)), "解析请求数据", login)
	return login, nil
}

// 处理返回数据
func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	log.Debugw(fmt.Sprint(ctx.Value(service.ContextReqUUid)), "请求结束返回值", response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorWrapper struct {
	Error string `json:"errors"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
