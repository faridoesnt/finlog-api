ALTER TABLE users
  ADD COLUMN is_verified TINYINT(1) NOT NULL DEFAULT 0,
  ADD COLUMN verification_token VARCHAR(128),
  ADD COLUMN verification_expires_at DATETIME;
