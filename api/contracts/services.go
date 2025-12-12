package contracts

type Services struct {
	Auth         AuthService
	Categories   CategoryService
	Transactions TransactionService
	Budget       BudgetService
	KeyBackup    KeyBackupService
	Import       ImportService
	Email        EmailService
}
