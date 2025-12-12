CREATE TABLE IF NOT EXISTS email_messages (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  resend_id VARCHAR(64) NOT NULL,
  to_email VARCHAR(255) NOT NULL,
  subject VARCHAR(255) NULL,

  status ENUM('sent','delivered','bounced','complained','failed') NOT NULL DEFAULT 'sent',
  last_error TEXT NULL,
  last_event_at DATETIME(3) NULL,

  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),

  UNIQUE KEY uq_email_messages_resend_id (resend_id),
  KEY idx_email_messages_to_email (to_email),
  KEY idx_email_messages_status (status)
);

CREATE TABLE IF NOT EXISTS email_events (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  resend_id VARCHAR(64) NOT NULL,
  event_type VARCHAR(64) NOT NULL,
  to_email VARCHAR(255) NOT NULL,

  error TEXT NULL,
  occurred_at DATETIME(3) NOT NULL,
  received_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

  raw_payload JSON NULL,

  UNIQUE KEY uq_email_events_resend_event (resend_id, event_type),
  KEY idx_email_events_to_email (to_email),
  KEY idx_email_events_type_time (event_type, occurred_at)
);

CREATE TABLE IF NOT EXISTS email_suppressions (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  email VARCHAR(255) NOT NULL,
  reason ENUM('complained','bounced') NOT NULL,
  resend_id VARCHAR(64) NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

  UNIQUE KEY uq_email_suppressions_email_reason (email, reason),
  KEY idx_email_suppressions_email (email)
);
