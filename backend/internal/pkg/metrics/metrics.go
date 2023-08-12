package metrics

import (
	"context"
	"expvar"
	"runtime"
)

type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

var m *metrics

type ctxKey int

const key ctxKey = 1

func init() {
	m = &metrics{
		goroutines: expvar.NewInt("gorountines"),
		requests:   expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
}

// Set sets the metrics data into the context.
func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

// AddGoroutines refreshes the goroutine metric every 100 requests.
func AddGoroutines(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		if v.requests.Value()%100 == 0 {
			v.goroutines.Set(int64(runtime.NumGoroutine()))
		}
	}
}

// AddRequests increments the request metric by 1.
func AddRequests(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.requests.Add(1)
	}
}

// AddErrors increments the errors metric by 1.
func AddErrors(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Add(1)
	}
}

// AddPanics increments the panics metric by 1.
func AddPanics(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Add(1)
	}
}
