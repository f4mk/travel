// Package debug provides handler support for the debugging endpoints.
package debug

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/f4mk/travel/backend/travel-api/internal/app/service/check"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Config struct {
	Build   string
	Log     *zerolog.Logger
	DB      *sqlx.DB
	Service *check.Service
}

type Mux struct {
	*http.ServeMux
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	m.ServeMux.ServeHTTP(w, r)
}

func New(cfg Config) *Mux {

	mux := &Mux{http.NewServeMux()}

	mux.HandleFunc("/debug/pprof", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	mux.HandleFunc("/debug/readiness", cfg.Service.Readiness)
	mux.HandleFunc("/debug/liveness", cfg.Service.Liveness)

	return mux
}
