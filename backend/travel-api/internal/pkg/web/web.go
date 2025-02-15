package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	mux        *httptreemux.ContextMux
	otmux      http.Handler
	tracer     trace.Tracer
	shutdown   chan os.Signal
	timeout    time.Duration
	middleware []Middleware
}
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func New(shutdown chan os.Signal, timeout time.Duration, tracer trace.Tracer, mw ...Middleware) *App {

	mux := httptreemux.NewContextMux()
	app := App{
		mux:        mux,
		otmux:      otelhttp.NewHandler(mux, "request"),
		tracer:     tracer,
		shutdown:   shutdown,
		timeout:    timeout,
		middleware: mw,
	}

	return &app
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.otmux.ServeHTTP(w, r)
}

func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx, span := a.startSpan(w, r)
		defer span.End()

		ctx, cancel := context.WithTimeout(ctx, a.timeout)
		defer cancel()

		v := Values{
			TraceID: span.SpanContext().TraceID().String(),
			Tracer:  a.tracer,
			Now:     time.Now().UTC(),
		}
		ctx = SetValues(ctx, &v)

		//First wrap specific mw
		hh := wrapMiddleware(mw, handler)

		//Second wrap common mw
		hh = wrapMiddleware(a.middleware, hh)

		if err := hh(ctx, w, r); err != nil {
			if validateShutdown(err) {
				a.SignalShutdown()
			}

			return
		}
	}

	a.mux.Handle(method, path, h)
}

func (a *App) startSpan(w http.ResponseWriter, r *http.Request) (context.Context, trace.Span) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	if a.tracer != nil {
		ctx, span = a.tracer.Start(ctx, "handler")
		span.SetAttributes(attribute.String("endpoint", r.RequestURI))
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(w.Header()))

	return ctx, span
}

func (a *App) SignalShutdown() {

	a.shutdown <- syscall.SIGTERM
}
