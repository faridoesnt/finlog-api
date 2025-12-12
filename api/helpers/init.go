package helpers

import "finlog-api/api/contracts"

var app *contracts.App

func Init(a *contracts.App) {
	app = a
}
