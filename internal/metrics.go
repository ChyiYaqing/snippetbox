package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"net/http"
)

const serverNamespace = "srv_instance"

// statusRecorder to record the status code from the http.ResponseWriter
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Flush() {
	rw.ResponseWriter.(http.Flusher).Flush()
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: serverNamespace,
		Subsystem: "api_requests",
		Name:      "total",
		Help:      "当前服务副本的接口总请求数",
	},
	[]string{"path", "method"},
)

var successRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: serverNamespace,
		Subsystem: "api_requests",
		Name:      "success_total",
		Help:      "当前服务副本的接口总请求数",
	}, []string{"path", "method"})

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: serverNamespace,
		Name:      "response_status",
		Help:      "Status of HTTP response",
	},
	[]string{"path", "method", "code"},
)

var responseTimeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: serverNamespace,
	Subsystem: "api_requests",
	Name:      "duration_second",
	Help:      "当前接口响应时间ms",
	Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}, //
}, []string{"path", "method", "code"})

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(successRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(responseTimeHistogram)
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		path = strings.ReplaceAll(path, "//", "/")
		rw := NewResponseWriter(w)

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		fmt.Printf("path=%s, method=%s, code=%d, duration=%f \n", path, r.Method, rw.statusCode, float64(duration))
		statusCode := rw.statusCode
		responseStatus.WithLabelValues(path, r.Method, strconv.Itoa(statusCode)).Inc()
		responseTimeHistogram.WithLabelValues(path, r.Method, strconv.Itoa(statusCode)).Observe(float64(duration))
		totalRequests.WithLabelValues(path, r.Method).Inc()
		if statusCode == 200 {
			successRequests.WithLabelValues(path, r.Method).Inc()
		}
	})
}
