package budget

const (
	getMonthlyBudget = `
		SELECT 
			0 AS income,
			0 AS expense,
			MAX(updated_at) AS last_updated
		FROM transactions
		WHERE user_id = ? AND YEAR(occurred_at) = ? AND MONTH(occurred_at) = ?
	`
)
