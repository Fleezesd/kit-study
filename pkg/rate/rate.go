package rate

import "go.uber.org/ratelimit"

func NewRateLimiter() ratelimit.Limiter {
	return ratelimit.New(1) // 一秒请求1次
}
