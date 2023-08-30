package middleware

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

func RateLimit(log *zerolog.Logger, rps int) web.Middleware {
	var mu sync.Mutex
	ipLimiters := make(map[string]*rate.Limiter)

	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v, err := web.GetValues(ctx)

			if err != nil {
				return web.NewShutdownError("value missing from context")
			}
			ip := r.Header.Get("X-Forwarded-For")

			mu.Lock()
			limiter, ok := ipLimiters[ip]

			if !ok {
				limiter = rate.NewLimiter(rate.Limit(rps), rps)
				ipLimiters[ip] = limiter
			}
			mu.Unlock()

			err = limiter.Wait(ctx)

			// Context was done either by cancelation or timeout
			if err != nil {
				log.Warn().Msgf("user reached rate limit: %s", v.TraceID)
				return web.NewRequestError(errors.New("too many requests"), http.StatusTooManyRequests)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m

}
