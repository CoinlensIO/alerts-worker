package events

import (
	"encoding/json"
	"time"
)

// EventData interface for different event types
type EventData interface{}

// Event wrapper for all events
type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Data      EventData `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

type BinanceMarkPriceEvent struct {
	Symbol          string  `json:"symbol"`
	Price           float64 `json:"price"`
	IndexPrice      float64 `json:"index_price"`
	FundingRate     float64 `json:"funding_rate"`
	NextFundingTime int64   `json:"next_funding_time"`
	Timestamp       int64   `json:"timestamp"`
}

func UnmarshalEvent(data []byte) (*Event, error) {
	// First unmarshal to get the type
	var baseEvent Event
	if err := json.Unmarshal(data, &baseEvent); err != nil {
		return nil, err
	}

	return &baseEvent, nil
}
