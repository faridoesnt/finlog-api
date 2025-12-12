package entities

import (
	"encoding/json"
	"time"
)

type UpsertEmailEventParams struct {
	ResendID   string
	EventType  string
	ToEmail    string
	Error      string
	OccurredAt time.Time
	RawPayload json.RawMessage
}