package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/abilioesteves/goh/gohserver"
	"github.com/abilioesteves/whisper/web/metrics"
)

// GetPrometheusMiddleware gets the middleware that feeds a prometheus registry with basic data over the http requests
func GetPrometheusMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := new(gohserver.StatusWriter).Init(w)
			next.ServeHTTP(sw, r)
			metrics.Latency.WithLabelValues(strconv.Itoa(sw.StatusCode), r.Method, r.URL.Path).Observe(time.Since(start).Seconds()) // feed prometheus
		})
	}
}
