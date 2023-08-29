package metrics

type Metrics struct {
	Goroutines   float64
	Requests     float64
	Errors       float64
	Panics       float64
	RequestTimes map[string]int64
}
