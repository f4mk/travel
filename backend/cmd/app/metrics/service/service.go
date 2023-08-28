package service

import (
	"net/http"

	metrics "github.com/f4mk/api/cmd/app/metrics/usecase"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

type Metrics struct {
	log  *zerolog.Logger
	core *metrics.Core
}

func New(l *zerolog.Logger, c *metrics.Core) *Metrics {

	return &Metrics{
		log:  l,
		core: c,
	}
}

func (s *Metrics) Serve(w http.ResponseWriter, r *http.Request) {
	m, err := s.core.CollectMetrics()
	if err != nil {
		s.log.Err(err).Msg("error collecting metrics")
		if err := writeResponse(w, err.Error(), http.StatusInternalServerError); err != nil {
			s.log.Err(err).Msg("error sending err response")
		}
		return
	}
	goroutines.Set(m.Goroutines)
	requests.Set(m.Requests)
	errors.Set(m.Errors)
	panics.Set(m.Panics)
	for bucketLabel, count := range m.RequestTimes {
		requestTimeBuckets.WithLabelValues(bucketLabel).Set(float64(count))
	}
	promhttp.Handler().ServeHTTP(w, r)
}

func writeResponse(w http.ResponseWriter, data string, statusCode int) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)

	if _, err := w.Write([]byte(data)); err != nil {
		return err
	}
	return nil
}
