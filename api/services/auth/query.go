package auth

const (
	findByEmail = `
		SELECT * FROM users WHERE email = ? LIMIT 1
	`

	findByID = `
		SELECT * FROM users WHERE id = ? LIMIT 1
	`

	findByVerificationToken = `
		SELECT * FROM users WHERE verification_token = ? LIMIT 1
	`

	insertUser = `
		INSERT INTO users (email, name, role, password, is_verified, verification_token, verification_expires_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	updateVerificationToken = `
		UPDATE users
		SET verification_token = ?, verification_expires_at = ?
		WHERE id = ?
	`

	markUserVerified = `
		UPDATE users
		SET is_verified = 1, verification_token = NULL, verification_expires_at = NULL
		WHERE id = ?
	`
)
