package transaction

const (
	listTransactions = `
		SELECT 
			t.id,
			t.user_id,
			t.category_id,
			c.name AS category_name,
			t.payload_ciphertext,
			t.payload_nonce,
			t.payload_tag,
			t.occurred_at,
			t.is_expense,
			t.created_at,
			t.updated_at
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = ? AND YEAR(t.occurred_at) = ? AND MONTH(t.occurred_at) = ?
		ORDER BY t.occurred_at DESC, t.id DESC
	`

	findTransactionByID = `
		SELECT 
			t.id,
			t.user_id,
			t.category_id,
			c.name AS category_name,
			t.payload_ciphertext,
			t.payload_nonce,
			t.payload_tag,
			t.occurred_at,
			t.is_expense,
			t.created_at,
			t.updated_at
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		WHERE t.id = ? AND t.user_id = ?
		LIMIT 1
	`

	insertTransaction = `
	INSERT INTO transactions (user_id, category_id, payload_ciphertext, payload_nonce, payload_tag, occurred_at, is_expense, batch_id)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

	updateTransaction = `
		UPDATE transactions
		SET payload_ciphertext = ?, payload_nonce = ?, payload_tag = ?, occurred_at = ?, is_expense = ?, category_id = ?, updated_at = NOW()
		WHERE id = ? AND user_id = ?
	`

	deleteTransaction = `
		DELETE FROM transactions WHERE id = ? AND user_id = ?
	`
)
