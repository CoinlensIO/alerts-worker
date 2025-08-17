package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type WorkerMetrics struct {
	// Event processing metrics
	EventProcessingDuration *prometheus.HistogramVec
	EventsProcessedTotal    *prometheus.CounterVec
	EventProcessingErrors   *prometheus.CounterVec

	// Queue metrics
	QueueSize       *prometheus.GaugeVec
	QueueLatency    *prometheus.HistogramVec
	EventAgeSeconds *prometheus.HistogramVec

	// Worker metrics
	WorkerStatus      *prometheus.GaugeVec
	WorkerBusy        *prometheus.GaugeVec
	WorkerLastSuccess *prometheus.GaugeVec
	WorkerLastError   *prometheus.GaugeVec

	// Rate limiting metrics
	RateLimitWaitDuration *prometheus.HistogramVec
	RateLimitExceeded     *prometheus.CounterVec

	// Memory metrics
	MemoryUsage *prometheus.GaugeVec
}

func InitWorkerMetrics() *WorkerMetrics {
	return &WorkerMetrics{
		EventProcessingDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "worker_event_processing_duration_seconds",
				Help:    "Time taken to process each event",
				Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"queue", "worker_id", "event_type"},
		),

		EventsProcessedTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "worker_events_processed_total",
				Help: "Total number of processed events",
			},
			[]string{"queue", "worker_id", "event_type", "status"},
		),

		EventProcessingErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "worker_processing_errors_total",
				Help: "Total number of event processing errors",
			},
			[]string{"queue", "worker_id", "event_type", "error_type"},
		),

		QueueSize: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "worker_queue_size",
				Help: "Current number of events in queue",
			},
			[]string{"queue"},
		),

		QueueLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "worker_queue_latency_seconds",
				Help:    "Time events spend in queue before processing",
				Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"queue"},
		),

		EventAgeSeconds: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "worker_event_age_seconds",
				Help:    "Age of events when processed",
				Buckets: []float64{1, 5, 10, 30, 60, 300, 600, 1800, 3600},
			},
			[]string{"queue", "event_type"},
		),

		WorkerStatus: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "worker_status",
				Help: "Worker status (1=active, 0=inactive)",
			},
			[]string{"queue", "worker_id"},
		),

		WorkerBusy: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "worker_busy",
				Help: "Whether worker is currently processing (1=busy, 0=idle)",
			},
			[]string{"queue", "worker_id"},
		),

		WorkerLastSuccess: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "worker_last_success_timestamp",
				Help: "Timestamp of last successful event processing",
			},
			[]string{"queue", "worker_id"},
		),

		WorkerLastError: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "worker_last_error_timestamp",
				Help: "Timestamp of last error",
			},
			[]string{"queue", "worker_id"},
		),

		RateLimitWaitDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "worker_rate_limit_wait_duration_seconds",
				Help:    "Time spent waiting for rate limiter",
				Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5},
			},
			[]string{"queue", "worker_id"},
		),

		RateLimitExceeded: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "worker_rate_limit_exceeded_total",
				Help: "Number of times rate limit was exceeded",
			},
			[]string{"queue", "worker_id"},
		),

		MemoryUsage: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "worker_memory_bytes",
				Help: "Memory usage by worker",
			},
			[]string{"queue", "worker_id"},
		),
	}
}
