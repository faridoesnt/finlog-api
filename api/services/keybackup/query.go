package keybackup

const (
	selectActiveKeyQuery = `
		SELECT 
			id,
			user_id,
			encrypted_data_key,
			salt,
			is_active,
			rotated_at,
			deleted_at,
			created_at,
			updated_at
		FROM user_encrypted_data_keys
		WHERE user_id = ? AND is_active = 1
		ORDER BY created_at DESC
		LIMIT 1
	`

	insertKeyQuery = `
		INSERT INTO user_encrypted_data_keys (user_id, encrypted_data_key, salt, is_active)
		VALUES (?, ?, ?, ?)
	`

	deactivateActiveKeyQuery = `
		UPDATE user_encrypted_data_keys
		SET is_active = 0, rotated_at = ?, deleted_at = ?, updated_at = NOW()
		WHERE user_id = ? AND is_active = 1
	`

	rotationSummaryQuery = `
		SELECT 
			COUNT(1) AS rotation_count,
			MAX(rotated_at) AS last_rotated_at
		FROM user_encrypted_data_keys
		WHERE user_id = ? AND is_active = 0
	`
)
