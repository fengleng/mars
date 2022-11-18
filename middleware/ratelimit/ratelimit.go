package ratelimit

import (
	"context"

	"github.com/gososy/sorpc/aegis/ratelimit"
	"github.com/gososy/sorpc/aegis/ratelimit/bbr"

	"github.com/gososy/sorpc/errors"
	"github.com/gososy/sorpc/middleware"
)

// ErrLimitExceed is service unavailable due to rate limit exceeded.
var ErrLimitExceed = errors.New(429, "RATELIMIT", "service unavailable due to rate limit exceeded")

// Option is ratelimit option.
type Option func(*options)

// WithLimiter set Limiter implementation,
// default is bbr limiter
func WithLimiter(limiter ratelimit.Limiter) Option {
	return func(o *options) {
		o.limiter = limiter
	}
}

type options struct {
	limiter ratelimit.Limiter
}

// Server ratelimiter middleware
func Server(opts ...Option) middleware.Middleware {
	options := &options{
		limiter: bbr.NewLimiter(),
	}
	for _, o := range opts {
		o(options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			done, e := options.limiter.Allow()
			if e != nil {
				// rejected
				return nil, ErrLimitExceed
			}
			// allowed
			reply, err = handler(ctx, req)
			done(ratelimit.DoneInfo{Err: err})
			return
		}
	}
}
