package email

import (
	"context"
	"encoding/json"
	"finlog-api/api/constants"
	"finlog-api/api/contracts"
	"finlog-api/api/entities"
	"finlog-api/api/models/request"
	"fmt"
	"time"

	"github.com/resend/resend-go/v3"
)

type Service struct {
	app    *contracts.App
	repo contracts.EmailRepository
	client *resend.Client
	from   string
}

func Init(app *contracts.App) contracts.EmailService {
	apiKey := app.Config[constants.ResendAPIKey]
	if apiKey == "" {
		app.Logger.Panic().Msg("RESEND_API_KEY is not set")
	}

	from := app.Config[constants.EmailFrom]
	if from == "" {
		app.Logger.Panic().Msg("EMAIL_FROM is not set")
	}

	repo := initRepository(app)

	return &Service{
		app:    app,
		repo: repo,
		client: resend.NewClient(apiKey),
		from:   from,
	}
}

func (s *Service) SendEmail(to, subject, html string) error {
	const (
		maxRetries = 3
		timeout    = 10 * time.Second
	)

	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		resp, err := s.client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
			From:    s.from,
			To:      []string{to},
			Subject: subject,
			Html:    html,
		})

		cancel()

		if err == nil {
			s.app.Logger.Info().
				Str("email", to).
				Str("resend_id", resp.Id).
				Int("attempt", attempt).
				Msg("email_sent")

			return nil
		}

		lastErr = err

		s.app.Logger.Warn().
			Err(err).
			Str("email", to).
			Int("attempt", attempt).
			Msg("email_send_failed")

		// exponential backoff: 1s, 2s, 3s
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	return fmt.Errorf("email failed after %d retries: %w", maxRetries, lastErr)
}

func (s *Service) HandleWebhook(ctx context.Context, payload request.ResendWebhookPayload, raw []byte) error {
	switch payload.Type {
	case "email.delivered", "email.bounced", "email.complained", "email.failed":
	default:
		s.app.Logger.Warn().Str("type", payload.Type).Msg("unknown_email_webhook")
		return nil
	}

	occurredAt := time.Now()
	if payload.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, payload.CreatedAt); err == nil {
			occurredAt = t
		}
	}

	rawJSON := json.RawMessage(raw)

	if err := s.repo.UpsertEmailEvent(ctx, entities.UpsertEmailEventParams{
		ResendID:    payload.Data.ID,
		EventType:   payload.Type,
		ToEmail:     payload.Data.To,
		Error:       payload.Data.Error,
		OccurredAt:  occurredAt,
		RawPayload:  rawJSON,
	}); err != nil {
		return fmt.Errorf("upsert email event: %w", err)
	}

	if err := s.repo.ApplyEmailStatus(ctx, payload.Data.ID, payload.Data.To, payload.Type, occurredAt, payload.Data.Error); err != nil {
		return fmt.Errorf("apply email status: %w", err)
	}

	s.app.Logger.Info().
		Str("type", payload.Type).
		Str("email", payload.Data.To).
		Str("resend_id", payload.Data.ID).
		Msg("email_webhook_received")

	return nil
}