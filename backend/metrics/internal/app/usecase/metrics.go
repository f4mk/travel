package metrics

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Collector interface {
	Collect() (map[string]interface{}, error)
}

type Core struct {
	collector Collector
	log       *zerolog.Logger
}

func New(l *zerolog.Logger, c Collector) *Core {
	return &Core{
		log:       l,
		collector: c,
	}
}

func (c *Core) CollectMetrics() (Metrics, error) {
	rawMetrics, err := c.collector.Collect()
	if err != nil {
		c.log.Err(err).Msg("failed to get metrics from collector")
		return Metrics{}, err
	}
	return convertToPrometheusFormat(rawMetrics)
}

func convertToPrometheusFormat(rm map[string]any) (Metrics, error) {

	var metrics Metrics

	goroutines, err := extractMetric(rm, "goroutines")
	if err != nil {
		return Metrics{}, err
	}
	metrics.Goroutines = goroutines

	requests, err := extractMetric(rm, "requests")
	if err != nil {
		return Metrics{}, err
	}
	metrics.Requests = requests

	errors, err := extractMetric(rm, "errors")
	if err != nil {
		return Metrics{}, err
	}
	metrics.Errors = errors

	panics, err := extractMetric(rm, "panics")
	if err != nil {
		return Metrics{}, err
	}
	metrics.Panics = panics

	requestTimes, ok := rm["requestTimes"].(map[string]interface{})
	if !ok {
		return Metrics{}, fmt.Errorf("requestTimes is missing or is of wrong type")
	}

	metrics.RequestTimes = make(map[string]int64)
	for bucket, value := range requestTimes {
		count, ok := value.(float64)
		if !ok {
			return Metrics{}, fmt.Errorf("invalid count for bucket: %s", bucket)
		}
		metrics.RequestTimes[bucket] = int64(count)
	}

	return metrics, nil
}

func extractMetric(rm map[string]any, key string) (float64, error) {
	val, ok := rm[key]
	if !ok {
		return 0, fmt.Errorf("key %s not found", key)
	}

	floatVal, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("unexpected type for key %s, expected float64", key)
	}

	return floatVal, nil
}
