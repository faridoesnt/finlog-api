package category

const (
	listCategories = `
		SELECT id, user_id, name, is_expense, icon_key, is_active, created_at, updated_at
		FROM categories
		WHERE user_id = ?
		  AND is_active = 1
	`

	findCategoryByID = `
		SELECT id, user_id, name, is_expense, icon_key, is_active, created_at, updated_at
		FROM categories
		WHERE id = ? AND user_id = ?
		LIMIT 1
	`

	findCategoryByName = `
		SELECT id, user_id, name, is_expense, icon_key, is_active, created_at, updated_at
		FROM categories
		WHERE user_id = ? AND LOWER(name) = LOWER(?) AND is_expense = ?
		LIMIT 1
	`

	insertCategory = `
		INSERT INTO categories (user_id, name, is_expense, icon_key, is_active)
		VALUES (?, ?, ?, ?, 1)
	`

	updateCategory = `
		UPDATE categories
		SET name = ?, is_expense = ?, icon_key = ?, is_active = ?, updated_at = NOW()
		WHERE id = ? AND user_id = ?
	`

	deleteCategory = `
		UPDATE categories
		SET is_active = 0, updated_at = NOW()
		WHERE id = ? AND user_id = ?
	`
)
