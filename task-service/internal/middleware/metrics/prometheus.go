package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//TODO добавить handler label для отслеживания обработчиков

type Middleware struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestSize     *prometheus.SummaryVec
	responseSize    *prometheus.SummaryVec
}

func New(registry prometheus.Registerer, buckets []float64, prefix string) *Middleware {
	if buckets == nil {
		buckets = prometheus.ExponentialBuckets(0.1, 1.5, 5)
	}
	prefix = prefix + "_"
	m := &Middleware{
		requestsTotal: promauto.With(registry).NewCounterVec(
			prometheus.CounterOpts{
				Name: prefix + "requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "code"},
		),
		requestDuration: promauto.With(registry).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    prefix + "request_duration_seconds",
				Help:    "HTTP request duration.",
				Buckets: buckets,
			},
			[]string{"method", "code"},
		),
		requestSize: promauto.With(registry).NewSummaryVec(
			prometheus.SummaryOpts{
				Name: prefix + "request_size_bytes",
				Help: "HTTP request size.",
			},
			[]string{"method", "code"},
		),
		responseSize: promauto.With(registry).NewSummaryVec(
			prometheus.SummaryOpts{
				Name: prefix + "response_size_bytes",
				Help: "HTTP response size.",
			},
			[]string{"method", "code"},
		),
	}
	return m
}

func (m *Middleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := NewResponseWriter(w)

			defer func() {
				duration := time.Since(start).Seconds()

				statusCode := rw.Status()

				m.requestsTotal.WithLabelValues(
					r.Method,
					statusCode,
				).Inc()

				m.requestDuration.WithLabelValues(
					r.Method,
					statusCode,
				).Observe(duration)

				m.requestSize.WithLabelValues(
					r.Method,
					statusCode,
				).Observe(float64(r.ContentLength))

				m.responseSize.WithLabelValues(
					r.Method,
					statusCode,
				).Observe(float64(rw.Size()))
			}()

			next.ServeHTTP(rw, r)
		})
	}
}

type ResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK, 0}
}

func (rw *ResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *ResponseWriter) Status() string {
	return http.StatusText(rw.status)
}

func (rw *ResponseWriter) Size() int {
	return rw.size
}
