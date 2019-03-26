package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/abilioesteves/goh/gohserver"
	"github.com/abilioesteves/whisper/web/metrics"
)

// PrometheusMiddleware feeds the prometheus endpoint with basic data over the request performance
func PrometheusMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := new(gohserver.StatusWriter).Init(w)
		next.ServeHTTP(sw, r)
		metrics.Latency.WithLabelValues(strconv.Itoa(sw.StatusCode), r.Method, r.URL.Path).Observe(time.Since(start).Seconds()) // feed prometheus
	})

}
