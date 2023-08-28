package metrics

import (
	"context"
	"expvar"
	"runtime"
)

type metrics struct {
	goroutines   *expvar.Int
	requests     *expvar.Int
	errors       *expvar.Int
	panics       *expvar.Int
	requestTimes expvar.Map
}

var m *metrics

type ctxKey int

const key ctxKey = 1

func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
	m.requestTimes.Init()
	m.requestTimes.Set("<20ms", new(expvar.Int))
	m.requestTimes.Set("20ms-50ms", new(expvar.Int))
	m.requestTimes.Set("50ms-100ms", new(expvar.Int))
	m.requestTimes.Set("100ms-200ms", new(expvar.Int))
	m.requestTimes.Set("200ms-500ms", new(expvar.Int))
	m.requestTimes.Set(">500ms", new(expvar.Int))

	expvar.Publish("requestTimes", &m.requestTimes)
}

func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

func AddGoroutines(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		if v.requests.Value()%100 == 0 {
			v.goroutines.Set(int64(runtime.NumGoroutine()))
		}
	}
}

func AddRequests(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.requests.Add(1)
	}
}

func AddRequestsTime(ctx context.Context, d int64) {
	bucket := determineBucket(d)
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.requestTimes.Get(bucket).(*expvar.Int).Add(1)
	}

}

func AddErrors(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Add(1)
	}
}

func AddPanics(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Add(1)
	}
}

func determineBucket(d int64) string {
	switch {
	case d < 20:
		return "<20ms"
	case d < 50:
		return "20ms-50ms"
	case d < 100:
		return "50ms-100ms"
	case d < 200:
		return "100ms-200ms"
	case d < 500:
		return "200ms-500ms"
	default:
		return ">500ms"
	}
}
