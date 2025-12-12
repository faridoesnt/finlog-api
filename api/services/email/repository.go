package email

import (
	"context"
	"finlog-api/api/contracts"
	"finlog-api/api/datasources"
	"finlog-api/api/entities"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type EmailStatus string

const (
	StatusSent       EmailStatus = "sent"
	StatusDelivered  EmailStatus = "delivered"
	StatusFailed     EmailStatus = "failed"
	StatusBounced    EmailStatus = "bounced"
	StatusComplained EmailStatus = "complained"
)

var statusPriority = map[EmailStatus]int{
	StatusSent:       0,
	StatusDelivered:  1,
	StatusFailed:     2,
	StatusBounced:    3,
	StatusComplained: 4,
}

type Repository struct {
	app  *contracts.App
	stmt Statement
}

type Statement struct {
	upsertEmailEvent *sqlx.Stmt
}

func initRepository(app *contracts.App) contracts.EmailRepository {
	stmts := Statement{
		upsertEmailEvent: datasources.Prepare(app.Ds.WriterDB, upsertEmailEvent),
	}

	r := Repository{
		app:  app,
		stmt: stmts,
	}

	return &r
}

func (r *Repository) UpsertEmailEvent(ctx context.Context, p entities.UpsertEmailEventParams) error {
	_, err := r.stmt.upsertEmailEvent.ExecContext(ctx, p.ResendID, p.EventType, p.ToEmail, nullIfEmpty(p.Error), p.OccurredAt, nullJSON(p.RawPayload))
	return err
}

func (r *Repository) ApplyEmailStatus(ctx context.Context, resendID string, toEmail string, eventType string, at time.Time, lastErr string) error {
	r.app.Logger.Debug().
		Str("resend_id", resendID).
		Str("event", eventType).
		Msg("apply_email_status_start")

	status, ok := mapEventToStatus(eventType)
	if !ok {
		return nil
	}

	nextPriority := statusPriority[status]

	tx, err := r.app.Ds.WriterDB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			r.app.Logger.Debug().
				Str("resend_id", resendID).
				Msg("tx_closed_without_error")
		}
	}()


	// ensure row exists
	_, err = tx.ExecContext(ctx, insertEmailMessage, resendID, toEmail, at)
	if err != nil {
		r.app.Logger.Error().
			Err(err).
			Str("resend_id", resendID).
			Msg("insert_email_message_failed")

		return err
	}

	// update only if higher priority
	_, err = tx.ExecContext(ctx, updateEmailMessage,
		nextPriority,
		status,
		at,
		lastErr,
		lastErr,
		lastErr,
		resendID,
	)
	if err != nil {
		r.app.Logger.Error().
			Err(err).
			Str("resend_id", resendID).
			Str("status", string(status)).
			Int("priority", nextPriority).
			Msg("update_email_message_failed")

		return err
	}

	// suppression rules
	if status == StatusComplained || status == StatusBounced {
		_, _ = tx.ExecContext(ctx, insertEmailSuppression, toEmail, status, resendID)
	}

	r.app.Logger.Info().
		Str("resend_id", resendID).
		Str("status", string(status)).
		Msg("email_status_tx_commit")

	err = tx.Commit()
	if err != nil {
		r.app.Logger.Error().
			Err(err).
			Str("resend_id", resendID).
			Msg("email_status_tx_commit_failed")
		return err
	}

	r.app.Logger.Info().
		Str("resend_id", resendID).
		Msg("email_status_tx_commit")

	return nil
}

func mapEventToStatus(event string) (EmailStatus, bool) {
	switch event {
	case "email.delivered":
		return StatusDelivered, true
	case "email.failed":
		return StatusFailed, true
	case "email.bounced":
		return StatusBounced, true
	case "email.complained":
		return StatusComplained, true
	default:
		return "", false
	}
}

func nullIfEmpty(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func nullJSON(b []byte) interface{} {
	if len(b) == 0 {
		return nil
	}
	return string(b)
}
