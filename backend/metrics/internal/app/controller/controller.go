package app

import (
	"encoding/json"
	"expvar"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/f4mk/travel/backend/metrics/internal/app/service"
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

func readiness(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonData, err := json.Marshal(res)

	if err != nil {
		log.Err(err).Msg("readiness: failed to marshal")
		return
	}

	if _, err := w.Write(jsonData); err != nil {
		log.Err(err).Msg("readiness: failed to write response")
		return
	}

}
