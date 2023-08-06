package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/f4mk/api/internal/pkg/auth"
	"github.com/f4mk/api/pkg/web"
)

func Authenticate(a *auth.Auth) web.Middleware {
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// TODO: may be this way
			// authHeader := r.Header.Get("Authorization")
			// parts := strings.Split(authHeader, " ")
			// if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			// 	err := errors.New("authenticate: failed: unexpected auth header format; expected: Bearer <token>")

			cookie, err := r.Cookie("auth_token")

			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					return web.NewRequestError(errors.New("authenticate: failed: auth_token cookie not found"), http.StatusUnauthorized)
				}
				err := errors.New("authenticate: failed: cannot parse auth_token cookie")
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			token := cookie.Value
			claims, err := a.ValidateToken(token)

			if err != nil {

				return web.NewRequestError(
					errors.New("authenticate: failed: token is invalid"),
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
