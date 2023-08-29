package service

import (
	"math"

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
	panics.Set()
}

func SetBucket(b string, v float64) {
	requestTimeBuckets.WithLabelValues(b).Set(float64(v))
}

func bucketToSeconds(bucket string) float64 {
	switch bucket {
	case "<20ms":
		return 0.02
	case "20ms-50ms":
		return 0.05
	case "50ms-100ms":
		return 0.1
	case "100ms-200ms":
		return 0.2
	case "200ms-500ms":
		return 0.5
	default:
		return math.Inf(+1)
	}
}
