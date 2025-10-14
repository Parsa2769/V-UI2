package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	TrafficBytesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "traffic_bytes_total",
			Help: "Total traffic in bytes",
		},
		[]string{"direction", "client_id"},
	)

	ActiveClients = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_clients",
			Help: "Number of active clients",
		},
	)

	NodesOnline = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "nodes_online",
			Help: "Number of online nodes",
		},
	)
)

func InitMetrics() {
	// Metrics are automatically registered via promauto
}

