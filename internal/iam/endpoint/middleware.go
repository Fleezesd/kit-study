package endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/fleezesd/kit-study/internal/pkg/errno"
	"github.com/fleezesd/kit-study/pkg/token"

	"go.uber.org/ratelimit"

	"github.com/fleezesd/kit-study/internal/iam/service"
	"github.com/fleezesd/kit-study/internal/pkg/log"

	"github.com/go-kit/kit/endpoint"
)

func logMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// 记录请求耗时
			defer func(begin time.Time) {
				log.Debugw(fmt.Sprint(ctx.Value(service.ContextReqUUid)), "logMiddleware", "endpoint", "耗时毫秒", time.Since(begin).Milliseconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}

func AuthMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			jwtToken := fmt.Sprint(ctx.Value(token.JWT_CONTEXT_KEY))
			if jwtToken == "" {
				log.Debugw(fmt.Sprint(ctx.Value(service.ContextReqUUid)), "authMiddleware", "endpoint", "error", errno.ErrTokenEmpty.Message)
				return "", errno.ErrTokenEmpty
			}
			jwtInfo, err := token.ParseToken(jwtToken)
			if err != nil {
				log.Debugw(fmt.Sprint(ctx.Value(service.ContextReqUUid)), "authMiddleware", "endpoint", "error", errno.ErrTokenInvalid.Message)
			}
			if v, ok := jwtInfo["Name"]; ok {
				ctx = context.WithValue(ctx, "name", v)
			}
			return next(ctx, request)
		}

	}
}

func UberRateMiddleware(limit ratelimit.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			limit.Take()
			return next(ctx, request)
		}
	}
}
