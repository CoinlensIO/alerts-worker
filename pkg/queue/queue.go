package queue

import (
	"alerts-worker/internal/events"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Queue implementation with zerolog
type Queue struct {
	client *redis.Client
	key    string
}

// NewQueue creates a new Queue instance
func NewQueue(client *redis.Client, key string) *Queue {
	return &Queue{
		client: client,
		key:    key,
	}
}

// Publish adds an event to the queue
func (q *Queue) Publish(ctx context.Context, event interface{}) error {
	return q.publishEvent(ctx, event)
}

// publishEvent handles the actual serialization and publishing to Redis
func (q *Queue) publishEvent(ctx context.Context, event interface{}) error {
	// Marshal the event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Push to Redis
	err = q.client.LPush(ctx, q.key, data).Err()
	if err != nil {
		return fmt.Errorf("failed to push event: %w", err)
	}

	// Log success with minimal details
	log.Debug().
		Str("queue", q.key).
		Str("event_type", getEventType(event)).
		Int("data_size", len(data)).
		Msg("event published to queue")

	return nil
}

// getEventType extracts the event type for logging
func getEventType(event interface{}) string {
	switch e := event.(type) {
	case *events.Event:
		return e.Type
	default:
		return "unknown"
	}
}
