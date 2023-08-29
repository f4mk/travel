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

func New(shutdown chan os.Signal, timeout time.Duration, mw ...Middleware) *App {

	mux := httptreemux.NewContextMux()
	app := App{
		mux:        mux,
		otmux:      otelhttp.NewHandler(mux, "request"),
		shutdown:   shutdown,
		timeout:    timeout,
		middleware: mw,
	}

	return &app
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.otmux.ServeHTTP(w, r)
}

func (a App) Handle(method string, path string, handler Handler, mw ...Middleware) {

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), a.timeout)
		defer cancel()

		ctx, span := a.startSpan(ctx, w, r)
		defer span.End()

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

func (a App) SignalShutdown() {

	a.shutdown <- syscall.SIGTERM
}

func (a *App) startSpan(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, trace.Span) {

	// There are times when the handler is called without a tracer, such
	// as with tests. We need a span for the trace id.
	span := trace.SpanFromContext(ctx)

	// If a tracer exists, then replace the span for the one currently
	// found in the context. This may have come from over the wire.
	if a.tracer != nil {
		ctx, span = a.tracer.Start(ctx, "pkg.web.handle")
		span.SetAttributes(attribute.String("endpoint", r.RequestURI))
	}

	// Inject the trace information into the response.
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(w.Header()))

	return ctx, span
}
