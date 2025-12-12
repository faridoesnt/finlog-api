package contracts

import (
	"context"
	"finlog-api/api/entities"
	"finlog-api/api/models/request"
	"time"
)

type EmailRepository interface {
	UpsertEmailEvent(ctx context.Context, p entities.UpsertEmailEventParams) error
	ApplyEmailStatus(ctx context.Context, resendID, toEmail, eventType string, at time.Time, lastErr string) error
}

type EmailService interface {
	SendEmail(to, subject, html string) error
	HandleWebhook(ctx context.Context, body request.ResendWebhookPayload, payload []byte) error
}