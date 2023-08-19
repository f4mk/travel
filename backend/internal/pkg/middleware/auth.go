package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/f4mk/api/internal/pkg/auth"
	"github.com/f4mk/api/pkg/web"
)

func Authenticate(a *auth.Auth) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			authHeader := r.Header.Get("Authorization")
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := auth.ErrAuthHeaderFormat
				return web.NewRequestError(err, http.StatusUnauthorized)
			}
			token := parts[1]
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
