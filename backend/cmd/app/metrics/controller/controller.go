package app

import (
	"context"
	"expvar"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	"github.com/f4mk/api/cmd/app/metrics/service"
	"github.com/f4mk/api/pkg/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Log            *zerolog.Logger
	MetricsService *service.Metrics
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

	mux.HandleFunc("/debug/readiness", readiness)
	mux.HandleFunc("/debug/liveness", readiness)
	mux.HandleFunc("/metrics", cfg.MetricsService.Serve)

	return mux
}

func readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	res := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	if err := web.Respond(ctx, w, res, http.StatusOK); err != nil {
		log.Err(err).Msg("readiness: failed to respond:")
	}
}
