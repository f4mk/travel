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
				err := errors.New("authenticate: failed: unexpected auth header format; expected: Bearer <token>")
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			token := parts[1]
			claims, err := a.ValidateToken(token)

			if err != nil {
				return web.NewRequestError(
					errors.New("authenticate: failed: token is invalid"),
					http.StatusUnauthorized,
				)
			}

			// Check if auth token is about to expire
			if time.Until(claims.ExpiresAt.Time) < time.Duration(10)*time.Minute { // assuming 10 minutes here
				refreshToken, err := r.Cookie("refresh_token")
				if err != nil {
					// Missing refresh token, this could mean potential session hijack or cookie deletion
					// Redirect to login
					http.Redirect(w, r, "/login", http.StatusFound)
					return nil
				}

				newClaims, err := a.ValidateRefreshToken(ctx, refreshToken.Value)
				if err != nil {
					// Invalid refresh token
					// Redirect to login
					http.Redirect(w, r, "/login", http.StatusFound)
					return nil
				}

				// Generate new tokens
				newAuthToken, newRefreshToken, err := a.GenerateTokens(newClaims)
				if err != nil {
					return web.NewRequestError(
						errors.New("authenticate: failed to generate new tokens"),
						http.StatusInternalServerError,
					)
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
