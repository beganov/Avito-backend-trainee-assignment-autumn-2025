package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	UsersCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "users_created_total",
			Help: "Общее количество созданных пользователей",
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
		UsersCreatedTotal, HttpDuration,
	)
}
