package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	registry         prometheus.Registerer
	defaultSkipPaths = map[string]struct{}{"/": {}, "/metrics": {}, "/favicon.ico": {}}
	defaultBuckets   = []float64{
		0.00001, // 10µs
		0.00005, // 50µs
		0.0001,  // 100µs
		0.0005,  // 500µs
		0.001,   // 1ms
		0.005,   // 5ms
		0.01,    // 10ms
		0.05,    // 50ms
		0.1,     // 100 ms
		0.5,     // 500 ms
		1.0,     // 1s
		5.0,     // 5s
		10.0,    // 10s
		25.0,    // 25s
	}
)

func getMetrics(registrer prometheus.Registerer, serviceName string, labels map[string]string) (
	total *prometheus.CounterVec, duration *prometheus.HistogramVec, inProgress *prometheus.GaugeVec) {

	registry = registrer
	constLabels := make(prometheus.Labels)
	if serviceName != "" {
		constLabels["service"] = serviceName
	}
	for label, value := range labels {
		constLabels[label] = value
	}

	total = promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Name:        "http_requests_total",
			Help:        "Count all http requests by status code, method and path.",
			ConstLabels: constLabels,
		},
		[]string{"status", "method", "path"},
	)

	duration = promauto.With(registry).NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "http_requests_duration",
			Help:        "Duration of all HTTP requests by status code, method and path.",
			ConstLabels: constLabels,
			Buckets:     defaultBuckets,
		},
		[]string{"status", "method", "path"},
	)

	inProgress = promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "http_requests_in_progress_total",
			Help:        "All the requests in progress by method and path",
			ConstLabels: constLabels,
		}, []string{"method", "path"},
	)

	return total, duration, inProgress
}
