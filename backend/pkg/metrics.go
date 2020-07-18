package pkg

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	http2 "net/http"
)

type Metrics struct {
	requestCounter  prometheus.Counter
	requestDuration prometheus.Histogram
}

func (m *Metrics) RequestCounter() prometheus.Counter {
	return m.requestCounter
}

func NewMetrics() *Metrics {
	requestCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "request_count",
		Help: "Request count",
	})
	requestDuration := promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "backend_http_request_duration_ms",
		Help:    "Request duration",
		Buckets: prometheus.LinearBuckets(0, 0.5, 30),
	})
	return &Metrics{requestCounter: requestCounter, requestDuration: requestDuration}
}

func (m *Metrics) GetMetrics() http2.Handler {
	return m
}

func (m *Metrics) Inc() {
	m.requestCounter.Inc()
}

func (m *Metrics) ServeHTTP(w http2.ResponseWriter, r *http2.Request) {
	fmt.Println("reading metrics")
	promhttp.Handler().ServeHTTP(w, r)
}

func (m *Metrics) RequestDuration() prometheus.Histogram {
	return m.requestDuration
}
