package request

type ResendWebhookPayload struct {
	Type string `json:"type"`
	Data struct {
		ID    string `json:"id"`
		To    string `json:"to"`
		From  string `json:"from,omitempty"`
		Subject string `json:"subject,omitempty"`

		Error string `json:"error,omitempty"`
	} `json:"data"`
	CreatedAt string `json:"created_at,omitempty"`
}
