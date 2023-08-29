package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/f4mk/travel/backend/pkg/web"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/metrics"
)

func Metrics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = metrics.Set(ctx)
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("value missing from context")
			}

			err = handler(ctx, w, r)

			d := time.Since(v.Now).Milliseconds()
			metrics.AddRequestsTime(ctx, d)
			metrics.AddRequests(ctx)
			metrics.AddGoroutines(ctx)

			if err != nil {
				metrics.AddErrors(ctx)
			}
			return err
		}
		return h
	}
	return m
}
