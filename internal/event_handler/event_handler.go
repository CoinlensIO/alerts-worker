package event_handler

import (
	"alerts-worker/internal/events"
	"alerts-worker/internal/service"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	BackoffFactor  float64
}

type EventHandler struct {
	alertService service.AlertService
	logger       *zerolog.Logger
	stopChan     chan struct{}
	retryConfig  RetryConfig
}

type EventHandlerOptions struct {
	RetryConfig *RetryConfig
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 500 * time.Millisecond,
		MaxBackoff:     10 * time.Second,
		BackoffFactor:  2.0,
	}
}

func NewEventHandler(
	alertService service.AlertService,
	logger *zerolog.Logger,
	opts *EventHandlerOptions) *EventHandler {

	logger.Debug().
		Interface("options", opts).
		Msg("initializing event handler")

	if opts == nil {
		logger.Debug().Msg("no options provided, using defaults")
		opts = &EventHandlerOptions{
			RetryConfig: &RetryConfig{},
		}
	}

	if opts.RetryConfig == nil {
		logger.Debug().Msg("no retry config provided, using defaults")
		defaultConfig := DefaultRetryConfig()
		opts.RetryConfig = &defaultConfig
	}

	handler := &EventHandler{
		alertService: alertService,
		logger:       logger,
		stopChan:     make(chan struct{}),
		retryConfig:  *opts.RetryConfig,
	}

	return handler
}

// HandleEvent processes an event received from the worker
func (h *EventHandler) HandleEvent(ctx context.Context, event *events.Event) error {
	return h.handleSyncEvent(ctx, event)
}

// handleSyncEvent processes a sync event
func (h *EventHandler) handleSyncEvent(ctx context.Context, event *events.Event) error {
	logger := h.logger.With().
		Str("event_type", event.Type).
		Time("event_created_at", event.CreatedAt).
		Logger()

	logger.Debug().
		Interface("event_data", event).
		Msg("processing sync event")

	var lastErr error
	backoff := h.retryConfig.InitialBackoff

	for attempt := 0; attempt <= h.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			logger.Debug().
				Int("attempt", attempt).
				Dur("backoff", backoff).
				Err(lastErr).
				Msg("preparing retry attempt")

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				// Continue after backoff period
			}

			backoff = time.Duration(float64(backoff) * h.retryConfig.BackoffFactor)
			if backoff > h.retryConfig.MaxBackoff {
				backoff = h.retryConfig.MaxBackoff
			}
		}

		switch event.Type {
		case events.EventTypeBinanceMarkPrice:

		}

		return nil
	}

	return fmt.Errorf("failed after %d retries: %v", h.retryConfig.MaxRetries, lastErr)
}

func (h *EventHandler) Stop() {
	h.logger.Debug().Msg("initiating shutdown sequence")
	close(h.stopChan)
	h.logger.Info().Msg("event handler shutdown complete")
}
