package endpoint

import (
	"context"
	"fmt"
	"kit-study/internal/pkg/log"
	"time"

	"kit-study/internal/iam/service"

	"github.com/go-kit/kit/endpoint"
)

func logMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// 记录请求耗时
			defer func(begin time.Time) {
				log.Debugw(fmt.Sprint(ctx.Value(service.ContextReqUUid)), "调用endpoint层logMiddleware", "处理完请求", "耗时毫秒", time.Since(begin).Milliseconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}
