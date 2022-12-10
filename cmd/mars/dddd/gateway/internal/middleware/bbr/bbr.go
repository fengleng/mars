package bbr

import (
	"bytes"
	"github.com/fengleng/dddd/gateway/internal/middleware"
	"github.com/fengleng/dddd/gateway/proto/client"
	"io"
	"net/http"

	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"
)

var _nopBody = io.NopCloser(&bytes.Buffer{})

func init() {
	middleware.Register("bbr", Middleware)
}

func Middleware(c *client.Middleware) (middleware.Middleware, error) {
	limiter := bbr.NewLimiter() //use default settings
	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			done, err := limiter.Allow()
			if err != nil {
				return &http.Response{
					Status:     http.StatusText(http.StatusTooManyRequests),
					StatusCode: http.StatusTooManyRequests,
					Body:       _nopBody,
				}, nil
			}
			resp, err := next.RoundTrip(req)
			done(ratelimit.DoneInfo{Err: err})
			return resp, err
		})
	}, nil
}
