package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HttpRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Количество HTTP-запросов",
		})

	HttpErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Ошибки HTTP",
		})

	HttpDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Время ответа API",
			Buckets: prometheus.DefBuckets,
		})
)

func Init() {
	prometheus.MustRegister(
		HttpRequestsTotal, HttpErrorsTotal, HttpDuration,
	)
}
