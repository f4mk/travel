package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	goroutines = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "goroutines",
		Help: "Current number of goroutines.",
	})
	requests = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "requests_total",
		Help: "Total number of requests processed.",
	})
	errors = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "errors_total",
		Help: "Total number of errors encountered.",
	})
	panics = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "panics_total",
		Help: "Total number of panics encountered.",
	})

	requestTimeBuckets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "request_time_buckets",
			Help: "Request time counts per bucket.",
		},
		[]string{"bucket"},
	)
)

func init() {
	prometheus.MustRegister(goroutines, requests, errors, panics, requestTimeBuckets)
}

func SetGoroutines(g float64) {
	goroutines.Set(g)
}

func SetRequests(r float64) {
	requests.Set(r)
}

func SetErrors(e float64) {
	errors.Set(e)
}

func SetPanics(p float64) {
	panics.Set(p)
}

func SetBucket(b string, v float64) {
	requestTimeBuckets.WithLabelValues(b).Set(v)
}
