package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

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
			claims, err := a.ValidateToken(token)
			if err != nil {
				return web.NewRequestError(
					err,
					http.StatusUnauthorized,
				)
			}

			// Check if auth token is about to expire
			if time.Until(claims.ExpiresAt.Time) < time.Duration(10)*time.Minute {
				refreshToken, err := r.Cookie("refresh_token")
				if err != nil {
					// Missing refresh token, this could mean potential session hijack or cookie deletion
					return web.NewRequestError(
						auth.ErrMissingRefreshToken,
						http.StatusUnauthorized,
					)
				}
				newClaims, err := a.ValidateRefreshToken(ctx, refreshToken.Value)
				if err != nil {
					// Invalid refresh token
					return web.NewRequestError(
						err,
						http.StatusUnauthorized,
					)
				}
				// Generate new tokens
				newAuthToken, err := a.GenerateToken(newClaims, a.AuthDuration)
				if err != nil {
					return err
				}
				newRefreshToken, err := a.GenerateToken(newClaims, a.RefreshDuration)
				if err != nil {
					return err
				}
				// Set the new auth token in the response header and the new refresh token as a cookie
				w.Header().Set("Authorization", "Bearer "+newAuthToken)
				http.SetCookie(w, &http.Cookie{
					Name:     "refresh_token",
					Value:    newRefreshToken,
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
				})
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
