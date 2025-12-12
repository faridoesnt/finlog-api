package handlers

import (
	. "finlog-api/api/contracts"
)

var app *App

func Init(a *App) {
	app = a
}
