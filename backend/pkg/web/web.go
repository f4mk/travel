package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
)

type App struct {
	*httptreemux.ContextMux
	shutdown   chan os.Signal
	middleware []Middleware
}
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func New(shutdown chan os.Signal, mw ...Middleware) *App {

	app := App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		middleware: mw,
	}

	return &app
}

func (a App) Handle(method string, path string, handler Handler, mw ...Middleware) {

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now().UTC(),
		}
		ctx = context.WithValue(ctx, keyValues, &v)

		//First wrap specific mw
		hh := wrapMiddleware(mw, handler)

		//Second wrap common mw
		hh = wrapMiddleware(a.middleware, hh)

		if err := hh(ctx, w, r); err != nil {
			a.SignalShutdown()

			return
		}
	}

	a.ContextMux.Handle(method, path, h)
}

func (a App) SignalShutdown() {

	a.shutdown <- syscall.SIGTERM
}
