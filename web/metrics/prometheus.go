package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Latency is a prometheus register for latency histogram data
var Latency *prometheus.HistogramVec

func init() {
	Latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "http_request_duration_seconds",
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path",
		ConstLabels: prometheus.Labels{"service": "whisper"},
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(Latency)
}
