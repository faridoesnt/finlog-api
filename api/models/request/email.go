package request

import "time"

type ResendWebhookPayload struct {
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Data      struct {
		EmailID   string    `json:"email_id"`
		To        []string  `json:"to"`
		From      string    `json:"from,omitempty"`
		Subject   string    `json:"subject,omitempty"`
		Error     string    `json:"error,omitempty"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"data"`
}