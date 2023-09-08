package web

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ctxKey int

const key ctxKey = 1

type Values struct {
	TraceID    string
	Tracer     trace.Tracer
	Now        time.Time
	StatusCode int
}

func GetValues(ctx context.Context) (*Values, error) {
	v, ok := ctx.Value(key).(*Values)

	if !ok {
		return nil, errors.New("value missing from context")
	}
	return v, nil
}

func SetValues(ctx context.Context, v *Values) context.Context {
	return context.WithValue(ctx, key, v)
}

func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*Values)

	if !ok {
		return "00000000000000000000000000000000"
	}
	return v.TraceID
}

func AddSpan(ctx context.Context, spanName string, keyValues ...attribute.KeyValue) (context.Context, trace.Span) {
	v, ok := ctx.Value(key).(*Values)
	if !ok || v.Tracer == nil {
		return ctx, trace.SpanFromContext(ctx)
	}

	ctx, span := v.Tracer.Start(ctx, spanName)
	for _, kv := range keyValues {
		span.SetAttributes(kv)
	}

	return ctx, span
}

func SetStatusCode(ctx context.Context, s int) error {

	v, ok := ctx.Value(key).(*Values)

	if !ok {
		return errors.New("value missing from context")
	}

	v.StatusCode = s

	return nil
}
