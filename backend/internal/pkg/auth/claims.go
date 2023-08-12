package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey int

const key ctxKey = 1

const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

func (c Claims) Authorize(roles ...string) bool {

	for _, has := range c.Roles {
		for _, want := range roles {
			if has == want {
				return true
			}
		}
	}

	return false
}

// SetClaims stores the claims in the context.
func SetClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, key, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) (Claims, error) {
	v, ok := ctx.Value(key).(Claims)
	if !ok {
		return Claims{}, errors.New("claims not found")
	}
	return v, nil
}
