package worker

import (
	"alerts-worker/internal/events"
	"alerts-worker/pkg/metrics"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type WorkerOptions struct {
	WorkerCount int
}

type Worker struct {
	client      *redis.Client
	key         string
	handler     func(context.Context, *events.Event) error
	logger      zerolog.Logger
	wg          sync.WaitGroup
	workerCount int
	metrics     *metrics.WorkerMetrics
	running     atomic.Bool
	shutdownCtx context.Context
	cancelFunc  context.CancelFunc
}

func NewWorker(client *redis.Client, key string, handler func(context.Context, *events.Event) error, opts *WorkerOptions, metrics *metrics.WorkerMetrics) *Worker {
	workerCount := opts.WorkerCount
	if workerCount <= 0 {
		workerCount = 1
	}
	ctx, cancel := context.WithCancel(context.Background())

	return &Worker{
		client:      client,
		key:         key,
		handler:     handler,
		workerCount: workerCount,
		metrics:     metrics,
		shutdownCtx: ctx,
		cancelFunc:  cancel,
	}
}

func (w *Worker) process(ctx context.Context, workerID int) error {
	logger := w.logger.With().Int("worker_id", workerID).Logger()
	workerIDStr := fmt.Sprintf("%d", workerID)

	// Set worker as active
	w.metrics.WorkerStatus.WithLabelValues(w.key, workerIDStr).Set(1)
	defer w.metrics.WorkerStatus.WithLabelValues(w.key, workerIDStr).Set(0)

	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("worker shutting down")
			return nil
		default:
			// Update queue size and log it for visibility
			if size, err := w.client.LLen(ctx, w.key).Result(); err == nil {
				w.metrics.QueueSize.WithLabelValues(w.key).Set(float64(size))
				if size > 0 {
					logger.Debug().Int64("queue_size", size).Msg("current queue size")
				}
			}

			// Set worker as idle
			w.metrics.WorkerBusy.WithLabelValues(w.key, workerIDStr).Set(0)

			// Create a separate context with timeout for BRPop
			brpopCtx, brpopCancel := context.WithTimeout(ctx, 10*time.Second)

			// Try to get an event from the queue
			result, err := w.client.BRPop(brpopCtx, 5*time.Second, w.key).Result()
			brpopCancel() // Cancel the BRPop context immediately after the operation

			if err != nil {
				if errors.Is(err, redis.Nil) {
					// This is normal when the queue is empty
					logger.Debug().Msg("queue empty, waiting for items")
					continue
				}
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					// Check if the parent context was canceled
					select {
					case <-ctx.Done():
						logger.Info().Msg("parent context canceled")
						return nil
					default:
						// It was just the BRPop context timing out
						logger.Debug().Msg("brpop context timed out")
						continue
					}
				}
				w.metrics.EventProcessingErrors.WithLabelValues(w.key, workerIDStr, "", "redis_error").Inc()
				logger.Error().Err(err).Msg("error getting event from queue")
				time.Sleep(1 * time.Second) // Add backoff on Redis errors
				continue
			}

			// Log successful pull for debugging
			logger.Debug().Str("payload_size", fmt.Sprintf("%d bytes", len(result[1]))).Msg("pulled item from queue")

			// Set worker as busy
			w.metrics.WorkerBusy.WithLabelValues(w.key, workerIDStr).Set(1)

			// First unmarshal to get the basic event details
			var baseEvent events.Event
			if err := json.Unmarshal([]byte(result[1]), &baseEvent); err != nil {
				w.metrics.EventProcessingErrors.WithLabelValues(w.key, workerIDStr, "", "unmarshal_error").Inc()
				logger.Error().Err(err).Str("payload", result[1]).Msg("error unmarshaling event")
				continue
			}

			// Log event information
			logger.Info().
				Str("event_type", baseEvent.Type).
				Time("created_at", baseEvent.CreatedAt).
				Msg("processing event")

			// Record event age
			eventAge := time.Since(baseEvent.CreatedAt).Seconds()
			w.metrics.EventAgeSeconds.WithLabelValues(w.key, baseEvent.Type).Observe(eventAge)

			// Apply a timeout to the handler
			processCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

			// Use the custom unmarshaler to get the appropriate event type
			event, err := events.UnmarshalEvent([]byte(result[1]))
			if err != nil {
				cancel()
				w.metrics.EventProcessingErrors.WithLabelValues(w.key, workerIDStr, baseEvent.Type, "unmarshal_error").Inc()
				logger.Error().Err(err).Str("payload", result[1]).Msg("error unmarshaling specific event type")
				continue
			}

			// Process the event and measure duration
			processStart := time.Now()
			if err := w.handler(processCtx, event); err != nil {
				cancel() // Cancel the context immediately
				w.metrics.EventProcessingErrors.WithLabelValues(w.key, workerIDStr, baseEvent.Type, "handler_error").Inc()
				w.metrics.WorkerLastError.WithLabelValues(w.key, workerIDStr).Set(float64(time.Now().Unix()))
				logger.Error().Err(err).Msg("error handling event")
				continue
			}

			// Record success metrics
			duration := time.Since(processStart).Seconds()
			w.metrics.EventProcessingDuration.WithLabelValues(w.key, workerIDStr, baseEvent.Type).Observe(duration)
			w.metrics.EventsProcessedTotal.WithLabelValues(w.key, workerIDStr, baseEvent.Type, "success").Inc()
			w.metrics.WorkerLastSuccess.WithLabelValues(w.key, workerIDStr).Set(float64(time.Now().Unix()))
			logger.Info().
				Str("event_type", baseEvent.Type).
				Float64("duration_seconds", duration).
				Msg("event processed successfully")

			cancel() // Cancel the context after successful processing
		}
	}
}

func (w *Worker) Start(ctx context.Context) error {
	if !w.running.CompareAndSwap(false, true) {
		return errors.New("already running")
	}

	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go func(workerID int) {
			defer w.wg.Done()
			_ = w.process(ctx, workerID)
		}(i)
	}

	return nil
}

func (w *Worker) Stop(timeout time.Duration) {
	w.cancelFunc()
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return
	case <-time.After(timeout):
		w.logger.Warn().Msg("workers forced to stop due to timeout")
	}
}
