package seeds

// CategorySeed defines a default category blueprint.
type CategorySeed struct {
	Name      string
	IsExpense bool
	IconKey   string
}

// DefaultCategories returns a set of income and expense categories for new users.
func DefaultCategories() []CategorySeed {
	return []CategorySeed{
		// Income
		{Name: "Gaji", IsExpense: false, IconKey: "salary"},
		{Name: "Bonus", IsExpense: false, IconKey: "bonus"},
		{Name: "Investasi", IsExpense: false, IconKey: "investment"},
		{Name: "Freelance", IsExpense: false, IconKey: "freelance"},
		{Name: "Hadiah", IsExpense: false, IconKey: "gift"},
		{Name: "Pendapatan Lainnya", IsExpense: false, IconKey: "category"},
		// Expense
		{Name: "Transportasi", IsExpense: true, IconKey: "transport"},
		{Name: "Makanan & Minuman", IsExpense: true, IconKey: "food"},
		{Name: "Belanja", IsExpense: true, IconKey: "shopping"},
		{Name: "Tagihan & Utilitas", IsExpense: true, IconKey: "bill"},
		{Name: "Kesehatan", IsExpense: true, IconKey: "health"},
		{Name: "Pendidikan", IsExpense: true, IconKey: "education"},
		{Name: "Hiburan", IsExpense: true, IconKey: "entertain"},
		{Name: "Lainnya", IsExpense: true, IconKey: "category"},
	}
}
