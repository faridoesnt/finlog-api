package services

import (
	"finlog-api/api/contracts"
	"finlog-api/api/services/auth"
	"finlog-api/api/services/budget"
	"finlog-api/api/services/category"
	"finlog-api/api/services/email"
	"finlog-api/api/services/importbatch"
	"finlog-api/api/services/keybackup"
	"finlog-api/api/services/transaction"
)

func Init(app *contracts.App) *contracts.Services {
	srv := &contracts.Services{
		Auth:         auth.Init(app),
		Categories:   category.Init(app),
		Transactions: transaction.Init(app),
		Budget:       budget.Init(app),
		KeyBackup:    keybackup.Init(app),
		Import:       importbatch.Init(app),
		Email:        email.Init(app),
	}

	app.Logger.Log().Msg("Initializing Services: Pass")

	return srv
}
