package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//TODO добавить handler label, translator label для отслеживания обработчиков, трансляторов

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
			rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			//rw := NewResponseWriter(w)

			defer func() {
				duration := time.Since(start).Seconds()

				statusCode := strconv.Itoa(rw.Status())

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
				).Observe(float64(rw.BytesWritten()))
			}()

			next.ServeHTTP(rw, r)
		})
	}
}
