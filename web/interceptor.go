package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/abilioesteves/goh/gohserver"
	"github.com/abilioesteves/whisper/web/metrics"
	"github.com/sirupsen/logrus"
)

// Interceptor defines the structure of an Interceptor
type Interceptor struct {
	Start   time.Time
	Writer  *gohserver.StatusWriter
	Request *http.Request
}

// Intercept intercepts a func handler and performs pre and post operations
func Intercept(f func(_ http.ResponseWriter, _ *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := new(Interceptor).Init(w, r)
		i.Before()
		defer i.After()
		f(i.Writer, i.Request) // execs the original func handler
	})
}

// Init initializes the request interceptor
func (i *Interceptor) Init(w http.ResponseWriter, r *http.Request) *Interceptor {
	i.Start = time.Now()
	i.Writer = new(gohserver.StatusWriter).Init(w)
	i.Request = r
	return i
}

// Before defines the operations to be executed before the original request handler
func (i *Interceptor) Before() {
	logrus.Infof("Validating token for request to '%v'", i.Request.RequestURI)
	// TODO token validation
}

// After defines the operations to be executed after the original request handler
func (i *Interceptor) After() {
	gohserver.HandleError(i.Writer)                                                                                                                 // gets the status
	metrics.Latency.WithLabelValues(strconv.Itoa(i.Writer.StatusCode), i.Request.Method, i.Request.URL.Path).Observe(time.Since(i.Start).Seconds()) // feed prometheus
	logrus.Infof("Done processing request to '%v' with status %v", i.Request.RequestURI, i.Writer.StatusCode)                                       // logs
}
