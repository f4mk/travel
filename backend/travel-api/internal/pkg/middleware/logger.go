package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
)

func Logger(log *zerolog.Logger) web.Middleware {

	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v, err := web.GetValues(ctx)

			if err != nil {
				log.Err(err).Msg("value missing from context")
				return web.NewShutdownError("value missing from context")
			}

			log.Info().Msgf(
				"%s : started   : %s %s -> %s",
				v.TraceID,
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
			)

			err = handler(ctx, w, r)

			log.Info().Msgf(
				"%s : completed : %s %s -> %s [%d] (%s)",
				v.TraceID,
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				v.StatusCode,
				time.Since(v.Now),
			)

			return err
		}

		return h
	}

	return m
}
