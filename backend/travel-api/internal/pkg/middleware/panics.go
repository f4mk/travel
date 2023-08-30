package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/metrics"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
)

func Panics(log *zerolog.Logger) web.Middleware {

	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			v, err := web.GetValues(ctx)

			if err != nil {
				return web.NewShutdownError("value missing from context")
			}

			defer func() {

				if r := recover(); r != nil {
					err = fmt.Errorf("panic: %v", r)
					metrics.AddPanics(ctx)
					log.Error().Msgf("%s : PANIC     :\n%s", v.TraceID, debug.Stack())
				}
			}()

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
