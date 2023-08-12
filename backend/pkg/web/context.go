package web

import (
	"context"
	"errors"
	"time"
)

type ctxKey int

const keyValues ctxKey = 1

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

func GetValues(ctx context.Context) (*Values, error) {
	v, ok := ctx.Value(keyValues).(*Values)

	if !ok {
		return nil, errors.New("value missing from context")
	}
	return v, nil
}

func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(keyValues).(*Values)

	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}
	return v.TraceID
}

func SetStatusCode(ctx context.Context, s int) error {

	v, ok := ctx.Value(keyValues).(*Values)

	if !ok {
		return errors.New("value missing from context")
	}

	v.StatusCode = s

	return nil
}
