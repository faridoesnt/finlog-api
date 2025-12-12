package email

const (
	upsertEmailEvent = `
		INSERT INTO email_events (resend_id, event_type, to_email, error, occurred_at, raw_payload)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		occurred_at = VALUES(occurred_at),
		error = VALUES(error),
		raw_payload = VALUES(raw_payload)
	`

	insertEmailMessage = `
		INSERT INTO email_messages (resend_id, to_email, status, last_event_at)
		VALUES (?, ?, 'sent', ?)
		ON DUPLICATE KEY UPDATE to_email = VALUES(to_email)
	`

	updateEmailMessage = `
		UPDATE email_messages
		SET
		status = IF(
			? >
			CASE status
			WHEN 'sent' THEN 0
			WHEN 'delivered' THEN 1
			WHEN 'failed' THEN 2
			WHEN 'bounced' THEN 3
			WHEN 'complained' THEN 4
			ELSE -1
			END,
			?, status
		),
		last_event_at = GREATEST(IFNULL(last_event_at, '1970-01-01'), ?),
		last_error = CASE
			WHEN ? IS NOT NULL AND ? != '' THEN ?
			ELSE last_error
		END
		WHERE resend_id = ?
	`

	insertEmailSuppression = `
		INSERT INTO email_suppressions (email, reason, resend_id)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE resend_id = VALUES(resend_id)
	`
)