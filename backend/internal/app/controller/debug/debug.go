// Package debug provides handler support for the debugging endpoints.
package debug

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/f4mk/api/internal/app/service/check"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Config struct {
	Build string
	Log   *zerolog.Logger
	DB    *sqlx.DB
}

type mux struct {
	*http.ServeMux
}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	m.ServeMux.ServeHTTP(w, r)
}

// Mux registers all the debug routes from the standard library into a new mux
// bypassing the use of the DefaultServerMux. Using the DefaultServerMux would
// be a security risk since a dependency could inject a handler into our service
// without us knowing it.
func New(cfg Config) *mux {

	mux := &mux{http.NewServeMux()}

	mux.HandleFunc("/debug/pprof", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	check := check.NewService(cfg.Build, cfg.Log, cfg.DB)

	mux.HandleFunc("/debug/readiness", check.Readiness)
	mux.HandleFunc("/debug/liveness", check.Liveness)

	return mux
}
