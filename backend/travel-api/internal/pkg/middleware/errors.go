package middleware

import (
	"context"
	"net/http"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
)

func Errors(log *zerolog.Logger) web.Middleware {

	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v, err := web.GetValues(ctx)

			if err != nil {
				return web.NewShutdownError("value missing from context")
			}

			if err := handler(ctx, w, r); err != nil {
				log.Err(err).Msgf("%s : ERROR     :", v.TraceID)

				if err := web.RespondError(ctx, w, err); err != nil {
					return err
				}

				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
