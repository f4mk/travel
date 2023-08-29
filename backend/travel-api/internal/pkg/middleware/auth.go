package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/f4mk/travel/backend/pkg/web"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
)

func Authenticate(a *auth.Auth) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			token, err := a.ExtractTokenFromHeader(r.Header.Get("Authorization"))
			if err != nil {
				return web.NewRequestError(err, http.StatusUnauthorized)
			}
			claims, err := a.ValidateToken(ctx, token)
			if err != nil {
				return web.NewRequestError(
					err,
					http.StatusUnauthorized,
				)
			}

			ctx = auth.SetClaims(ctx, claims)

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

func Authorize(roles ...string) web.Middleware {

	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims, err := auth.GetClaims(ctx)

			if err != nil {
				return web.NewRequestError(
					fmt.Errorf("authorize: failed: %s", err),
					http.StatusForbidden,
				)
			}

			if !claims.Authorize(roles...) {
				return web.NewRequestError(
					errors.New("authorize: failed: not authorized"),
					http.StatusForbidden,
				)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
